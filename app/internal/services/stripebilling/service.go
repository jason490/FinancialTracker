package stripebilling

import (
	"FinancialTracker/internal/services/financial"
	"FinancialTracker/internal/services/subscription"
	"FinancialTracker/internal/storage"

	"github.com/stripe/stripe-go/v86"
)

// Service handles Stripe Checkout, Customer Portal, and billing webhooks.
type Service struct {
	client       *stripe.Client
	store        *storage.Storage
	subscription *subscription.Service
}

// NewService creates a Stripe billing service.
func NewService(store *storage.Storage, sub *subscription.Service) *Service {
	return &Service{
		client:       stripe.NewClient(financial.StripeSecretKey()),
		store:        store,
		subscription: sub,
	}
}
