package stripefc

import (
	"errors"
	"os"

	"FinancialTracker/internal/config"
	"FinancialTracker/internal/services/subscription"
	"FinancialTracker/internal/services/tags"
	"FinancialTracker/internal/storage"

	"github.com/stripe/stripe-go/v86"
)

// Item status values stored in stripe_fc_items.status.
const (
	ItemStatusActive       = "active"
	ItemStatusNeedsReauth  = "needs_reauth"
	ItemStatusDisconnected = "disconnected"
	ItemStatusError        = "error"
	providerName           = "stripe"
)

// Sentinel errors for subscription and usage limits.
var (
	ErrStripeAPILimitExceeded  = errors.New("monthly Stripe API limit reached for your plan")
	ErrStripeItemLimitExceeded = errors.New("bank connection limit reached for your plan")
)

// Service implements Stripe Financial Connections bank linking.
type Service struct {
	client       *stripe.Client
	store        *storage.Storage
	tagService   *tags.TaggingService
	subscription *subscription.Service
}

// NewService initializes a Stripe Financial Connections service.
func NewService(store *storage.Storage, tagService *tags.TaggingService, sub *subscription.Service) *Service {
	key := stripeSecretKey()
	return &Service{
		client:       stripe.NewClient(key),
		store:        store,
		tagService:   tagService,
		subscription: sub,
	}
}

func stripeSecretKey() string {
	if config.IsDevelopment() {
		return os.Getenv("STRIPE_SANDBOX_SECRET")
	}
	return os.Getenv("STRIPE_PROD_SECRET")
}

// Name returns the provider identifier.
func (s *Service) Name() string {
	return providerName
}
