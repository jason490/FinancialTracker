package models

// Subscription tier identifiers stored on users.subscription_tier.
const (
	SubscriptionTierFree    = "free"
	SubscriptionTierPlus    = "plus"
	SubscriptionTierPremium = "premium"
)

// PlaidLimits describes per-user Plaid quotas for a subscription tier.
type PlaidLimits struct {
	MaxItems         int `json:"max_items"`
	MaxAPICallsMonth int `json:"max_api_calls_month"`
}

// PlaidUsage reports current Plaid consumption for the active billing period.
type PlaidUsage struct {
	SubscriptionTier string      `json:"subscription_tier"`
	Limits           PlaidLimits `json:"limits"`
	ActiveItems      int         `json:"active_items"`
	APICallsUsed     int         `json:"api_calls_used"`
	PeriodStart      int64       `json:"period_start"`
	PeriodEnd        int64       `json:"period_end"`
}

// UserSubscription holds billing fields used for quota periods and future Stripe integration.
type UserSubscription struct {
	Tier                  string `json:"tier"`
	SubscriptionStartedAt int64  `json:"subscription_started_at,omitempty"`
	StripeCustomerID      string `json:"stripe_customer_id,omitempty"`
	StripeSubscriptionID  string `json:"stripe_subscription_id,omitempty"`
	CreatedAt             int64  `json:"created_at"`
}

// BillingPeriod is a monthly usage window anchored to signup or subscription purchase.
type BillingPeriod struct {
	Anchor      int64 `json:"anchor"`
	PeriodStart int64 `json:"period_start"`
	PeriodEnd   int64 `json:"period_end"`
}

// TierPlan describes a selectable subscription tier for the client.
type TierPlan struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Limits            PlaidLimits `json:"limits"`
	PriceMonthlyCents int         `json:"price_monthly_cents"`
}
