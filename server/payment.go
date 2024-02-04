package server

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/AmirSolt/town-watch/models"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/subscription"
	"github.com/stripe/stripe-go/v76/webhook"
	"github.com/stripe/stripe-go/v76/webhookendpoint"
)

type SubscriptionID string

const (
	MonthlySubscriptionID SubscriptionID = "monthly"
	YearlySubscriptionID  SubscriptionID = "yearly"
)

type CheckoutConfig struct {
	interval string
	amount   int64
}

func (server *Server) loadStripe() {
	stripe.Key = server.Env.STRIPE_PRIVATE_KEY

	params := &stripe.WebhookEndpointParams{
		EnabledEvents: []*string{
			stripe.String("customer.subscription.created"),
			stripe.String("customer.subscription.deleted"),
			stripe.String("customer.subscription.paused"),
			stripe.String("customer.subscription.resumed"),
			stripe.String("customer.subscription.updated"),
		},
		URL: stripe.String(fmt.Sprintf("%s/payment/webhook/events", server.Env.HOME_URL)),
	}
	_, err := webhookendpoint.New(params)
	if err != nil {
		log.Fatalln("Error: init stripe webhook events: %w", err)
	}
}

func (server *Server) GetCheckoutUrl(user *models.User, subscriptionID SubscriptionID) (*stripe.CheckoutSession, error) {

	var checkoutConfig CheckoutConfig
	if subscriptionID == MonthlySubscriptionID {
		checkoutConfig = CheckoutConfig{
			interval: "month",
			amount:   1000,
		}
	} else if subscriptionID == YearlySubscriptionID {
		checkoutConfig = CheckoutConfig{
			interval: "year",
			amount:   10000,
		}
	} else {
		return nil, fmt.Errorf("checkout session id not found")

	}

	return server.createCheckoutSession(user, checkoutConfig)
}

func (server *Server) createCheckoutSession(user *models.User, checkoutConfig CheckoutConfig) (*stripe.CheckoutSession, error) {

	customer, err := server.DB.queries.GetCustomerByUserID(context.Background(), user.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error customer find by userID: %w", err)
	}

	var stripeCustomerID *string
	if customer.StripeCustomerID == "" {
		stripeCustomerID = nil
	} else {
		stripeCustomerID = &customer.StripeCustomerID
	}

	params := &stripe.CheckoutSessionParams{
		ClientReferenceID: stripe.String(string(user.ID.Bytes[:])),
		Customer:          stripeCustomerID,
		Mode:              stripe.String("subscription"),
		CustomerEmail:     stripe.String(user.Email),
		ReturnURL:         stripe.String(server.Env.HOME_URL),
		SuccessURL:        stripe.String(fmt.Sprintf("%s/payment/success", server.Env.HOME_URL)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:    stripe.String("USD"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{},
					Recurring: &stripe.CheckoutSessionLineItemPriceDataRecurringParams{
						Interval:      stripe.String(checkoutConfig.interval),
						IntervalCount: stripe.Int64(1),
					},
					UnitAmount: stripe.Int64(checkoutConfig.amount),
				},
				Quantity: stripe.Int64(1),
			},
		},
	}

	result, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("checkout session creation failed: %w", err)
	}

	if stripeCustomerID == nil {
		_, err := server.DB.queries.CreateCustomer(context.Background(), models.CreateCustomerParams{
			StripeCustomerID: result.Customer.ID,
			UserID:           user.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("error user customerID could not be updated: %w", err)
		}
	}

	return result, nil
}

func (server *Server) CancelSubscription(subscriptionID SubscriptionID, comment string) error {
	// cancel stripe subsc
	_, errSub := subscription.Cancel(string(subscriptionID), &stripe.SubscriptionCancelParams{
		CancellationDetails: &stripe.SubscriptionCancelCancellationDetailsParams{
			Comment: stripe.String(comment),
		},
	})
	if errSub != nil {
		return fmt.Errorf("canceling stripe subscription: %w", errSub)
	}

	return nil
}

func (server *Server) HandleStripeWebhook(ginContext *gin.Context) {
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
		newSubsc, err := subscription.Get(event.Data.Object["subscription"].(string), nil)
		if err != nil {
			return fmt.Errorf("converting raw event to subscription object: %w", err)
		}

		// find customer
		customer, err := server.DB.queries.GetCustomerByStripeID(context.Background(), cust.ID)
		if err != nil {
			return fmt.Errorf("could not find customer by stripe id: %w", err)
		}

		// deactivate subsc in DB
		deSubscription, err := server.DB.queries.DeactivateSubscriptionByCustomerID(context.Background(), customer.ID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("deactivating customer active subscription: %w", err)
		}
		if err != sql.ErrNoRows {
			// cancel stripe subsc
			_, errSub := subscription.Cancel(deSubscription.StripeSubscriptionID, &stripe.SubscriptionCancelParams{
				CancellationDetails: &stripe.SubscriptionCancelCancellationDetailsParams{
					Feedback: stripe.String("switched_service"),
				},
			})
			if errSub != nil {
				return fmt.Errorf("canceling stripe subscription: %w", errSub)
			}
		}

		// create this subscription
		_, errSubsc := server.DB.queries.CreateSubscription(context.Background(), models.CreateSubscriptionParams{
			StripeSubscriptionID: newSubsc.ID,
			TierID:               "",
			IsActive:             true,
			CustomerID:           customer.ID,
		})
		if errSubsc != nil {
			return fmt.Errorf("creating new subscription: %w", errSubsc)
		}

		return nil

	}
	if event.Type == "customer.subscription.deleted" {
		subs, err := subscription.Get(event.Data.Object["subscription"].(string), nil)
		if err != nil {
			return fmt.Errorf("converting raw event to subscription object: %w", err)
		}

		errDe := server.DB.queries.DeactivateSubscriptionByStripeID(context.Background(), subs.ID)
		if errDe != nil {
			return fmt.Errorf("deactivating subscription by id: %w", err)
		}

		return nil
	}

	fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	return nil
}
