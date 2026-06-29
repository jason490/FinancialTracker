package subscription

import (
	"errors"
	"time"

	"FinancialTracker/internal/config"
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/storage"
)

var (
	ErrInvalidTier       = errors.New("invalid subscription tier")
	ErrStripeRequired    = errors.New("plan changes require Stripe checkout in production")
	ErrAlreadyOnTier     = errors.New("already on this plan")
	ErrAPILimitExceeded       = errors.New("monthly API limit reached for your plan")
	ErrUserNotFound           = errors.New("user not found")
	ErrSubscriptionsDisabled  = errors.New("subscriptions are disabled")
)

// TierCatalog lists available plans with skeleton pricing until Stripe is wired up.
var TierCatalog = []models.TierPlan{
	{
		ID:   models.SubscriptionTierFree,
		Name: "Free",
		Limits: models.PlaidLimits{
			MaxItems:         2,
			MaxAPICallsMonth: 20,
		},
		PriceMonthlyCents: 0,
	},
	{
		ID:   models.SubscriptionTierPlus,
		Name: "Plus",
		Limits: models.PlaidLimits{
			MaxItems:         5,
			MaxAPICallsMonth: 50,
		},
		PriceMonthlyCents: 299,
	},
	{
		ID:   models.SubscriptionTierPremium,
		Name: "Premium",
		Limits: models.PlaidLimits{
			MaxItems:         15,
			MaxAPICallsMonth: 100,
		},
		PriceMonthlyCents: 599,
	},
}

// Service manages subscription tiers and billing periods before Stripe is integrated.
type Service struct {
	store            *storage.Storage
	allowDirectChange bool
}

// NewService creates a subscription service. allowDirectChange enables skeleton tier changes (development).
func NewService(store *storage.Storage, allowDirectChange bool) *Service {
	return &Service{
		store:            store,
		allowDirectChange: allowDirectChange,
	}
}

// LimitsForTier returns Plaid quotas for a subscription tier, defaulting to free.
func LimitsForTier(tier string) models.PlaidLimits {
	for _, plan := range TierCatalog {
		if plan.ID == tier {
			return plan.Limits
		}
	}
	return TierCatalog[0].Limits
}

// BillingAnchor returns the Unix timestamp that monthly usage resets are based on.
// Free users anchor to signup; paid users anchor to when the subscription started.
func (s *Service) BillingAnchor(user *models.User) int64 {
	if user == nil {
		return time.Now().Unix()
	}
	if user.SubscriptionTier != models.SubscriptionTierFree && user.SubscriptionStartedAt > 0 {
		return user.SubscriptionStartedAt
	}
	if user.CreatedAt > 0 {
		return user.CreatedAt
	}
	return time.Now().Unix()
}

// CurrentBillingPeriod returns the active monthly usage window for a user.
func (s *Service) CurrentBillingPeriod(userID int64) (*models.BillingPeriod, error) {
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	anchor := s.BillingAnchor(user)
	now := time.Now().Unix()
	start := BillingPeriodStart(anchor, now)
	end := BillingPeriodEnd(start, anchor)

	return &models.BillingPeriod{
		Anchor:      anchor,
		PeriodStart: start,
		PeriodEnd:   end,
	}, nil
}

// GetUserSubscription loads subscription fields for a user.
func (s *Service) GetUserSubscription(userID int64) (*models.UserSubscription, error) {
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	tier := user.SubscriptionTier
	if tier == "" {
		tier = models.SubscriptionTierFree
	}

	return &models.UserSubscription{
		Tier:                  tier,
		SubscriptionStartedAt: user.SubscriptionStartedAt,
		StripeCustomerID:      user.StripeCustomerID,
		StripeSubscriptionID:  user.StripeSubscriptionID,
		CreatedAt:             user.CreatedAt,
	}, nil
}

// ChangeTier updates the user's plan. In production this will be replaced by Stripe webhooks.
func (s *Service) ChangeTier(userID int64, tier string) error {
	if !config.SubscriptionsEnabled() {
		return ErrSubscriptionsDisabled
	}
	if !isValidTier(tier) {
		return ErrInvalidTier
	}
	if !s.allowDirectChange {
		return ErrStripeRequired
	}

	user, err := s.store.GetUserByID(userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}

	currentTier := user.SubscriptionTier
	if currentTier == "" {
		currentTier = models.SubscriptionTierFree
	}
	if currentTier == tier {
		return ErrAlreadyOnTier
	}

	var startedAt *int64
	if tier != models.SubscriptionTierFree {
		now := time.Now().Unix()
		if currentTier == models.SubscriptionTierFree {
			startedAt = &now
		} else if user.SubscriptionStartedAt > 0 {
			existing := user.SubscriptionStartedAt
			startedAt = &existing
		} else {
			startedAt = &now
		}
	}

	return s.store.UpdateUserSubscription(userID, tier, startedAt, nil, nil)
}

// CanChangePlanDirectly reports whether skeleton tier changes are enabled.
func (s *Service) CanChangePlanDirectly() bool {
	return config.SubscriptionsEnabled() && s.allowDirectChange
}

// ReserveAPICall validates the user has at least one API call available in the
// current billing period and records its consumption. Returns
// ErrAPILimitExceeded when the per-tier monthly cap has been reached.
//
// All metered server features (Plaid syncs, CSV exports, etc.) share the same
// monthly quota recorded in plaid_api_usage so users see a single unified
// number on the Plan tab.
func (s *Service) ReserveAPICall(userID int64) error {
	if !config.SubscriptionsEnabled() {
		return nil
	}

	sub, err := s.GetUserSubscription(userID)
	if err != nil {
		return err
	}
	period, err := s.CurrentBillingPeriod(userID)
	if err != nil {
		return err
	}

	limits := LimitsForTier(sub.Tier)
	unlimited, err := s.HasUnlimitedLimits(userID)
	if err != nil {
		return err
	}
	if unlimited {
		return nil
	}

	used, err := s.store.GetPlaidAPICallCount(userID, period.PeriodStart)
	if err != nil {
		return err
	}
	if used >= limits.MaxAPICallsMonth {
		return ErrAPILimitExceeded
	}
	return s.store.IncrementPlaidAPICallCount(userID, period.PeriodStart)
}

func isValidTier(tier string) bool {
	switch tier {
	case models.SubscriptionTierFree, models.SubscriptionTierPlus, models.SubscriptionTierPremium:
		return true
	default:
		return false
	}
}
