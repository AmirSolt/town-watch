package server

import (
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

type CheckoutID string

const (
	MonthlyCheckoutID CheckoutID = "monthly"
	YearlyCheckoutID  CheckoutID = "yearly"
)

func (server *Server) GetCheckoutUrl(checkoutID CheckoutID) {
	stripe.Key = "sk_test_Hrs6SAopgFPF0bZXSN3f6ELN"

	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("https://example.com/success"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:    stripe.String("USD"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{},
					Recurring: &stripe.CheckoutSessionLineItemPriceDataRecurringParams{
						Interval:      stripe.String("month"),
						IntervalCount: stripe.Int64(1),
					},
					UnitAmount: stripe.Int64(1),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}
	result, err := session.New(params)
}
