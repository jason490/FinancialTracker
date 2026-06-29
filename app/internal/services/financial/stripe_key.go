package financial

import (
	"FinancialTracker/internal/config"
	"os"
)

// StripeSecretKey resolves the Stripe secret key from environment.
func StripeSecretKey() string {
	if config.IsDevelopment() {
		return os.Getenv("STRIPE_SANDBOX_SECRET")
	}
	return os.Getenv("STRIPE_PROD_SECRET")
}

// StripePublishableKey resolves the Stripe publishable key from environment.
func StripePublishableKey() string {
	if config.IsDevelopment() {
		return os.Getenv("STRIPE_SANDBOX_PUBLISHABLE")
	}
	return os.Getenv("STRIPE_PUBLISHABLE")
}

// StripeConfigured reports whether a Stripe secret key is available.
func StripeConfigured() bool {
	return StripeSecretKey() != ""
}
