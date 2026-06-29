package subscription

import (
	"os"
	"testing"
)

func TestTierForPriceID(t *testing.T) {
	t.Setenv("ENV", "production")
	t.Setenv("STRIPE_PRICE_PLUS", "price_plus")
	t.Setenv("STRIPE_PRICE_PREMIUM", "price_premium")

	if got := TierForPriceID("price_plus"); got != "plus" {
		t.Fatalf("expected plus, got %q", got)
	}
	if got := TierForPriceID("price_premium"); got != "premium" {
		t.Fatalf("expected premium, got %q", got)
	}
	if got := TierForPriceID("price_unknown"); got != "" {
		t.Fatalf("expected empty tier, got %q", got)
	}
}

func TestBillingEnabledRequiresPrices(t *testing.T) {
	t.Setenv("ENV", "production")
	t.Setenv("STRIPE_PRICE_PLUS", "")
	t.Setenv("STRIPE_PRICE_PREMIUM", "")
	if BillingEnabled() {
		t.Fatal("expected billing disabled without prices")
	}

	t.Setenv("STRIPE_PRICE_PLUS", "price_plus")
	t.Setenv("STRIPE_PRICE_PREMIUM", "price_premium")
	if !BillingEnabled() {
		t.Fatal("expected billing enabled with prices")
	}
}

func TestPriceIDForTierSandboxOverride(t *testing.T) {
	t.Setenv("ENV", "development")
	t.Setenv("STRIPE_SANDBOX_PRICE_PLUS", "price_sandbox_plus")
	t.Setenv("STRIPE_PRICE_PLUS", "price_prod_plus")

	if got := PriceIDForTier("plus"); got != "price_sandbox_plus" {
		t.Fatalf("expected sandbox price, got %q", got)
	}
}

func TestParseUserIDOverridesEnv(t *testing.T) {
	// Smoke test env parsing without DB writes in this unit test.
	t.Setenv("SUBSCRIPTION_OVERRIDE_USER_IDS", "")
	ParseUserIDOverrides(nil)
	_ = os.Getenv("SUBSCRIPTION_OVERRIDE_USER_IDS")
}
