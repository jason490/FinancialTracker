package config

import (
	"os"
	"strings"
)

const (
	ProviderPlaid  = "plaid"
	ProviderStripe = "stripe"
)

// IsDevelopment reports whether the app runs in local development mode.
func IsDevelopment() bool {
	return os.Getenv("ENV") == "development"
}

// EnvOr returns the environment variable value or a fallback when unset.
func EnvOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// EnvBool parses a boolean environment variable. Truthy values: 1, true, yes, on.
// When unset, default is returned.
func EnvBool(key string, defaultValue bool) bool {
	raw := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if raw == "" {
		return defaultValue
	}
	switch raw {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return defaultValue
	}
}

// SubscriptionsEnabled reports whether paid plans, Stripe billing, and usage limits are active.
func SubscriptionsEnabled() bool {
	return EnvBool("SUBSCRIPTIONS_ENABLED", true)
}

// RegistrationGateEnabled reports whether new accounts require an admin invite code.
func RegistrationGateEnabled() bool {
	return !SubscriptionsEnabled()
}

// RegistrationAdminEmails returns emails allowed to issue registration invite codes.
func RegistrationAdminEmails() []string {
	raw := os.Getenv("REGISTRATION_ADMIN_EMAILS")
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ";")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		email := strings.ToLower(strings.TrimSpace(part))
		if email != "" {
			out = append(out, email)
		}
	}
	return out
}

// FinancialProvider returns the active bank connection provider.
func FinancialProvider() string {
	switch strings.TrimSpace(strings.ToLower(os.Getenv("FINANCIAL_PROVIDER"))) {
	case ProviderStripe:
		return ProviderStripe
	default:
		return ProviderPlaid
	}
}

// PlaidEnvironment returns the configured Plaid environment name.
func PlaidEnvironment() string {
	env := strings.TrimSpace(strings.ToLower(os.Getenv("PLAID_ENV")))
	if env == "" {
		if IsDevelopment() {
			return "sandbox"
		}
		return "production"
	}
	return env
}
