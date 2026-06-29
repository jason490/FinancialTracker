package financial

import "os"

const (
	ProviderStripe = "stripe"
	ProviderPlaid  = "plaid"
)

// ActiveProvider returns the configured bank connection provider.
func ActiveProvider() string {
	switch os.Getenv("FINANCIAL_PROVIDER") {
	case ProviderStripe:
		return ProviderStripe
	default:
		return ProviderPlaid
	}
}
