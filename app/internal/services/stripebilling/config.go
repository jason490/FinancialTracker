package stripebilling

import (
	"os"
	"strings"

	"FinancialTracker/internal/config"
	"FinancialTracker/internal/services/financial"
)

// webhookSecret resolves the Stripe webhook signing secret for the active environment.
func webhookSecret() string {
	if config.IsDevelopment() {
		if secret := strings.TrimSpace(os.Getenv("STRIPE_SANDBOX_WEBHOOK_SECRET")); secret != "" {
			return secret
		}
	}
	return strings.TrimSpace(os.Getenv("STRIPE_WEBHOOK_SECRET"))
}

// WebhookConfigured reports whether webhook verification can run.
func WebhookConfigured() bool {
	return webhookSecret() != ""
}

// BillingReady reports whether checkout can be offered (keys, prices, and webhook).
func BillingReady() bool {
	if !config.SubscriptionsEnabled() {
		return false
	}
	return financial.StripeConfigured() && WebhookConfigured() && billingPricesConfigured()
}

func billingPricesConfigured() bool {
	if config.IsDevelopment() {
		if strings.TrimSpace(os.Getenv("STRIPE_SANDBOX_PRICE_PLUS")) != "" &&
			strings.TrimSpace(os.Getenv("STRIPE_SANDBOX_PRICE_PREMIUM")) != "" {
			return true
		}
	}
	return strings.TrimSpace(os.Getenv("STRIPE_PRICE_PLUS")) != "" &&
		strings.TrimSpace(os.Getenv("STRIPE_PRICE_PREMIUM")) != ""
}
