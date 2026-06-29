package config

import "testing"

func TestValidateRequiredEnv_PlaidDevelopment(t *testing.T) {
	t.Setenv("ENV", "development")
	t.Setenv("FINANCIAL_PROVIDER", "plaid")
	t.Setenv("ENCRYPTION_KEY", "12345678901234567890123456789012")
	t.Setenv("PLAID_CLIENT_ID", "client-id")
	t.Setenv("PLAID_SANDBOX_SECRET", "sandbox-secret")
	t.Setenv("PLAID_PROD_SECRET", "")

	if err := ValidateRequiredEnv(); err != nil {
		t.Fatalf("expected valid development plaid env, got %v", err)
	}
}

func TestValidateRequiredEnv_MissingEncryptionKey(t *testing.T) {
	t.Setenv("ENV", "development")
	t.Setenv("ENCRYPTION_KEY", "")
	t.Setenv("PLAID_CLIENT_ID", "client-id")
	t.Setenv("PLAID_SANDBOX_SECRET", "sandbox-secret")

	if err := ValidateRequiredEnv(); err == nil {
		t.Fatal("expected validation error for missing ENCRYPTION_KEY")
	}
}

func TestValidateRequiredEnv_ProductionRequiresMail(t *testing.T) {
	t.Setenv("ENV", "production")
	t.Setenv("FINANCIAL_PROVIDER", "plaid")
	t.Setenv("ENCRYPTION_KEY", "12345678901234567890123456789012")
	t.Setenv("PLAID_CLIENT_ID", "client-id")
	t.Setenv("PLAID_ENV", "production")
	t.Setenv("PLAID_PROD_SECRET", "prod-secret")
	t.Setenv("API_PUBLIC_URL", "https://example.com")
	t.Setenv("MAIL_SMTP_HOST", "")

	if err := ValidateRequiredEnv(); err == nil {
		t.Fatal("expected validation error for missing production mail settings")
	}
}

func TestSubscriptionsEnabled(t *testing.T) {
	t.Setenv("SUBSCRIPTIONS_ENABLED", "")
	if !SubscriptionsEnabled() {
		t.Fatal("expected subscriptions enabled by default")
	}

	t.Setenv("SUBSCRIPTIONS_ENABLED", "false")
	if SubscriptionsEnabled() {
		t.Fatal("expected subscriptions disabled")
	}
}
