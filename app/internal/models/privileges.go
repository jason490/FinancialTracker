package models

// UserPrivileges holds per-user subscription overrides set by the operator.
type UserPrivileges struct {
	UserID           int64  `json:"user_id"`
	UnlimitedLimits  bool   `json:"unlimited_limits"`
	StripeCouponID   string `json:"stripe_coupon_id,omitempty"`
	Notes            string `json:"notes,omitempty"`
}

// SubscriptionPrivileges is the client-facing summary of operator overrides.
type SubscriptionPrivileges struct {
	UnlimitedLimits bool `json:"unlimited_limits"`
	HasDiscount     bool `json:"has_discount"`
}
