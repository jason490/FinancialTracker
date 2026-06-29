package financial

import "FinancialTracker/internal/config"

const (
	ProviderStripe = config.ProviderStripe
	ProviderPlaid  = config.ProviderPlaid
)

// ActiveProvider returns the configured bank connection provider.
func ActiveProvider() string {
	return config.FinancialProvider()
}
