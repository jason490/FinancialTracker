package config

import (
	"fmt"
	"os"
	"strings"
)

// ValidateRequiredEnv ensures core environment variables are present before startup.
func ValidateRequiredEnv() error {
	var missing []string
	var problems []string


	switch FinancialProvider() {
	case ProviderPlaid:
		validatePlaidEnv(&missing)
	case ProviderStripe:
		validateStripeEnv(&missing)
	default:
		problems = append(problems, fmt.Sprintf("FINANCIAL_PROVIDER must be %q or %q", ProviderPlaid, ProviderStripe))
	}

	if !IsDevelopment() {
		// validateProductionEnv(&missing)
		if key := strings.TrimSpace(os.Getenv("ENCRYPTION_KEY")); key == "" {
			missing = append(missing, "ENCRYPTION_KEY")
		} else if len(key) != 32 {
			problems = append(problems, "ENCRYPTION_KEY must be exactly 32 bytes for AES-256")
		}
	}

	if len(missing) > 0 {
		problems = append(problems, fmt.Sprintf("missing required environment variables: %s", strings.Join(missing, ", ")))
	}
	if len(problems) == 0 {
		return nil
	}
	return fmt.Errorf("invalid environment configuration:\n- %s", strings.Join(problems, "\n- "))
}

func validatePlaidEnv(missing *[]string) {
	if strings.TrimSpace(os.Getenv("PLAID_CLIENT_ID")) == "" {
		*missing = append(*missing, "PLAID_CLIENT_ID")
	}

	env := PlaidEnvironment()
	if env == "production" && !IsDevelopment() {
		if strings.TrimSpace(os.Getenv("PLAID_PROD_SECRET")) == "" {
			*missing = append(*missing, "PLAID_PROD_SECRET")
		}
		return
	}

	if strings.TrimSpace(os.Getenv("PLAID_SANDBOX_SECRET")) == "" {
		*missing = append(*missing, "PLAID_SANDBOX_SECRET")
	}
}

func validateStripeEnv(missing *[]string) {
	if IsDevelopment() {
		if strings.TrimSpace(os.Getenv("STRIPE_SANDBOX_SECRET")) == "" {
			*missing = append(*missing, "STRIPE_SANDBOX_SECRET")
		}
		if strings.TrimSpace(os.Getenv("STRIPE_SANDBOX_PUBLISHABLE")) == "" {
			*missing = append(*missing, "STRIPE_SANDBOX_PUBLISHABLE")
		}
		return
	}

	if strings.TrimSpace(os.Getenv("STRIPE_PROD_SECRET")) == "" {
		*missing = append(*missing, "STRIPE_PROD_SECRET")
	}
	if strings.TrimSpace(os.Getenv("STRIPE_PUBLISHABLE")) == "" {
		*missing = append(*missing, "STRIPE_PUBLISHABLE")
	}
}

func validateProductionEnv(missing *[]string) {
	if strings.TrimSpace(os.Getenv("API_PUBLIC_URL")) == "" {
		*missing = append(*missing, "API_PUBLIC_URL")
	}

	for _, key := range []string{
		"MAIL_SMTP_HOST",
		"MAIL_SMTP_PORT",
		"MAIL_SMTP_USER",
		"MAIL_SMTP_PASSWORD",
		"MAIL_FROM",
	} {
		if strings.TrimSpace(os.Getenv(key)) == "" {
			*missing = append(*missing, key)
		}
	}
}
