package server

import (
	"context"
	"fmt"

	"github.com/AmirSolt/town-watch/models"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

type CheckoutID string

const (
	MonthlyCheckoutID CheckoutID = "monthly"
	YearlyCheckoutID  CheckoutID = "yearly"
)

type CheckoutConfig struct {
	interval string
	amount   int64
}

func (server *Server) GetCheckoutUrl(user *models.User, checkoutID CheckoutID) (*stripe.CheckoutSession, error) {
	stripe.Key = server.Env.STRIPE_PRIVATE_KEY

	var checkoutConfig CheckoutConfig
	if checkoutID == MonthlyCheckoutID {
		checkoutConfig = CheckoutConfig{
			interval: "month",
			amount:   1000,
		}
	} else if checkoutID == YearlyCheckoutID {
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

	params := &stripe.CheckoutSessionParams{
		ClientReferenceID: stripe.String(string(user.ID.Bytes[:])),
		Customer:          stripe.String(user.CustomerID.String),
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

	if user.CustomerID.String == "" {
		err := server.DB.queries.UpdateUserCustomerID(context.Background(), models.UpdateUserCustomerIDParams{
			CustomerID: pgtype.Text{String: result.Customer.ID},
			ID:         user.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("error user customerID could not be updated: %w", err)
		}
	}
	if user.CustomerID.String != result.Customer.ID {
		return nil, fmt.Errorf("error user customerID does not match stripe customerID: %w", err)
	}

	return result, nil
}
