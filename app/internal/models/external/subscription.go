package external

import (
	"FinancialTracker/internal/models"
)

// SubscriptionPayload is the client-facing subscription and billing summary.
type SubscriptionPayload struct {
	Tier                  string                        `json:"tier"`
	Billing               models.BillingPeriod          `json:"billing"`
	Limits                models.PlaidLimits            `json:"limits"`
	Plans                 []models.TierPlan             `json:"plans"`
	StripeConfigured      bool                          `json:"stripe_configured"`
	BillingEnabled        bool                          `json:"billing_enabled"`
	SubscriptionsEnabled  bool                          `json:"subscriptions_enabled"`
	HasActiveSubscription bool                          `json:"has_active_subscription"`
	CanChangePlan         bool                          `json:"can_change_plan"`
	Privileges            models.SubscriptionPrivileges `json:"privileges"`
}

// CheckoutSessionRequest starts a Stripe Checkout flow for a paid tier.
type CheckoutSessionRequest struct {
	Tier string `json:"tier"`
}

// CheckoutSessionResponse returns the hosted Stripe Checkout URL.
type CheckoutSessionResponse struct {
	URL string `json:"url"`
}

// BillingPortalResponse returns the hosted Stripe Customer Portal URL.
type BillingPortalResponse struct {
	URL string `json:"url"`
}

// ChangeSubscriptionRequest switches the user's subscription tier.
type ChangeSubscriptionRequest struct {
	Tier string `json:"tier"`
}

// ChangeSubscriptionResponse confirms a tier change.
type ChangeSubscriptionResponse struct {
	Subscription SubscriptionPayload `json:"subscription"`
}

