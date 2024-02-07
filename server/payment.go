package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/AmirSolt/town-watch/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/paymentmethod"
	"github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/webhook"
	"github.com/stripe/stripe-go/v76/webhookendpoint"
)

type TierConfig struct {
	tier     models.Tier
	name     string
	interval string
	amount   int64
}

func (server *Server) loadPayment() {

	// stripe key
	stripe.Key = server.Env.STRIPE_PRIVATE_KEY

	// webhook setup
	params := &stripe.WebhookEndpointParams{
		EnabledEvents: []*string{
			// stripe.String("customer.subscription.updated"),
			stripe.String("customer.subscription.created"),
			stripe.String("customer.subscription.deleted"),
			// stripe.String("customer.subscription.resumed"),
			// stripe.String("customer.subscription.paused"),
			// stripe.String("payment_method.attached"),
			// stripe.String("payment_method.detached"),
		},
		URL: stripe.String(fmt.Sprintf("%s/payment/webhook/events", server.Env.HOME_URL)),
	}
	_, err := webhookendpoint.New(params)
	if err != nil {
		log.Fatalln("Error: init stripe webhook events: %w", err)
	}

	// TierConfig
	m := make(map[models.Tier]TierConfig)
	m[models.TierT0] = TierConfig{
		tier:     models.TierT0,
		name:     "Free",
		interval: "never",
		amount:   0,
	}
	m[models.TierT1] = TierConfig{
		tier:     models.TierT1,
		name:     "Monthly",
		interval: "month",
		amount:   1000,
	}
	m[models.TierT2] = TierConfig{
		tier:     models.TierT2,
		name:     "Yearly",
		interval: "year",
		amount:   10000,
	}
	server.TierConfigs = m
}

// ==================================================

func (server *Server) GetUserPaymentMethods(user *models.User) *customer.PaymentMethodIter {
	params := &stripe.CustomerListPaymentMethodsParams{
		Customer: stripe.String(user.StripeCustomerID.String),
	}
	params.Limit = stripe.Int64(5)
	return customer.ListPaymentMethods(params)
}
func (server *Server) DetachPaymentMethod(paymentMethodID string) error {
	_, err := paymentmethod.Detach(paymentMethodID, &stripe.PaymentMethodDetachParams{})
	if err != nil {
		return fmt.Errorf("canceling stripe subscription: %w", err)
	}
	return nil
}

func (server *Server) ChangeAutoPay(user *models.User, disable bool) error {
	params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(disable)}
	_, err := subscription.Update(user.StripeSubscriptionID.String, params)
	if err != nil {
		return fmt.Errorf("canceling stripe subscription: %w", err)
	}
	return nil
}

func (server *Server) CreateSubscriptionTier(user *models.User, tierConfig TierConfig) (*stripe.CheckoutSession, error) {
	return server.createCheckoutSession(user, tierConfig)
}
func (server *Server) ChangeSubscriptionTier(user *models.User, tierConfig TierConfig) (*stripe.CheckoutSession, error) {
	err := server.CancelSubscription(user)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return server.createCheckoutSession(user, tierConfig)
}
func (server *Server) CancelSubscription(user *models.User) error {
	_, errSub := subscription.Cancel(user.StripeSubscriptionID.String, &stripe.SubscriptionCancelParams{})
	if errSub != nil {
		return fmt.Errorf("canceling stripe subscription: %w", errSub)
	}
	return nil
}

func (server *Server) createCheckoutSession(user *models.User, tierConfig TierConfig) (*stripe.CheckoutSession, error) {

	var customerID *string = nil
	if user.StripeCustomerID.Valid {
		customerID = &user.StripeCustomerID.String
	}

	params := &stripe.CheckoutSessionParams{
		Customer:      customerID,
		Mode:          stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		CustomerEmail: stripe.String(user.Email),
		ReturnURL:     stripe.String(server.Env.HOME_URL),
		SuccessURL:    stripe.String(fmt.Sprintf("%s/user/wallet", server.Env.HOME_URL)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyUSD)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(tierConfig.name),
					},
					Recurring: &stripe.CheckoutSessionLineItemPriceDataRecurringParams{
						Interval:      stripe.String(tierConfig.interval),
						IntervalCount: stripe.Int64(1),
					},
					UnitAmount: stripe.Int64(getNewUnitAmount(server.TierConfigs[user.Tier], tierConfig)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: map[string]string{"tier": string(tierConfig.tier)},
	}

	result, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("checkout session creation failed: %w", err)
	}

	if customerID == nil {
		err := server.DB.queries.UpdateUserStripeCustomerID(context.Background(), models.UpdateUserStripeCustomerIDParams{
			StripeCustomerID: pgtype.Text{String: result.Customer.ID},
			ID:               user.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("error user customerID could not be updated: %w", err)
		}
	}

	return result, nil
}

func getNewUnitAmount(currentTierConfig TierConfig, targetTierConfig TierConfig) int64 {
	newCost := targetTierConfig.amount - currentTierConfig.amount
	if newCost < 0 {
		newCost = 0
	}
	return newCost
}

// ==================================================

func (server *Server) HandleStripeWebhook(ginContext *gin.Context) {
	// ==================================================================
	// The signature check is pulled directly from Stripe and it's not tested
	req := ginContext.Request
	w := ginContext.Writer

	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	endpointSecret := server.Env.STRIPE_WEBHOOK_KEY
	event, err := webhook.ConstructEvent(payload, req.Header.Get("Stripe-Signature"),
		endpointSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	// ==================================================================

	if err := server.handleStripeEvents(event); err != nil {
		fmt.Fprintf(os.Stderr, "Error handling event: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (server *Server) handleStripeEvents(event stripe.Event) error {
	if event.Type == "customer.subscription.created" {
		cust, err := customer.Get(event.Data.Object["customer"].(string), nil)
		if err != nil {
			return fmt.Errorf("converting raw event to customer object: %w", err)
		}
		subsc, err := subscription.Get(event.Data.Object["subscription"].(string), nil)
		if err != nil {
			return fmt.Errorf("converting raw event to subscription object: %w", err)
		}
		tier := event.Data.Object["metadata"].(string)

		user, errUser := server.DB.queries.GetUserByStripeCustomerID(context.Background(), pgtype.Text{String: cust.ID})
		if errUser != nil {
			return fmt.Errorf("could not find user by stripe id: %w", errUser)
		}

		errUpd := server.DB.queries.UpdateUserSubAndTier(context.Background(), models.UpdateUserSubAndTierParams{
			StripeSubscriptionID: pgtype.Text{String: subsc.ID},
			Tier:                 models.Tier(tier),
			ID:                   user.ID,
		})
		if errUpd != nil {
			return fmt.Errorf("could not update user UpdateUserSubAndTier: %w", errUpd)
		}

		return nil
	}

	if event.Type == "customer.subscription.deleted" {
		cust, err := customer.Get(event.Data.Object["customer"].(string), nil)
		if err != nil {
			return fmt.Errorf("converting raw event to customer object: %w", err)
		}
		subsc, err := subscription.Get(event.Data.Object["subscription"].(string), nil)
		if err != nil {
			return fmt.Errorf("converting raw event to subscription object: %w", err)
		}

		user, errUser := server.DB.queries.GetUserByStripeCustomerID(context.Background(), pgtype.Text{String: cust.ID})
		if errUser != nil {
			return fmt.Errorf("could not find user by stripe id: %w", errUser)
		}

		if user.StripeSubscriptionID.String == subsc.ID {
			errUpd := server.DB.queries.UpdateUserSubAndTier(context.Background(), models.UpdateUserSubAndTierParams{
				StripeSubscriptionID: pgtype.Text{String: "", Valid: false},
				Tier:                 models.TierT0,
				ID:                   user.ID,
			})
			if errUpd != nil {
				return fmt.Errorf("could not update user UpdateUserSubAndTier: %w", errUpd)
			}
		}
		return nil
	}

	fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	return nil
}
