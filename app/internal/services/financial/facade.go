package financial

import (
	"FinancialTracker/internal/services/plaid"
	"FinancialTracker/internal/services/stripefc"
	"FinancialTracker/internal/services/subscription"
	"FinancialTracker/internal/services/tags"
	"FinancialTracker/internal/storage"
)

// Facade delegates bank connection operations to the active provider.
type Facade struct {
	active Provider
}

// NewFacade wires Plaid and Stripe providers and selects the active one.
func NewFacade(store *storage.Storage, tagService *tags.TaggingService, sub *subscription.Service) *Facade {
	plaidSvc := plaid.NewPlaidService(store, tagService, sub)
	stripeSvc := stripefc.NewService(store, tagService, sub)

	var active Provider
	switch ActiveProvider() {
	case ProviderPlaid:
		active = NewPlaidAdapter(plaidSvc)
	default:
		active = stripeSvc
	}

	return &Facade{active: active}
}

// Active returns the configured provider implementation.
func (f *Facade) Active() Provider {
	return f.active
}
