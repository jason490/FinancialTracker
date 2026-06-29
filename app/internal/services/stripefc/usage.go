package stripefc

import (
	"FinancialTracker/internal/models"
)

// GetUsage returns current Stripe API consumption and limits for a user.
func (s *Service) GetUsage(userID int64) (*models.PlaidUsage, error) {
	sub, err := s.subscription.GetUserSubscription(userID)
	if err != nil {
		return nil, err
	}

	period, err := s.subscription.CurrentBillingPeriod(userID)
	if err != nil {
		return nil, err
	}

	limits, err := s.subscription.EffectiveLimits(userID)
	if err != nil {
		return nil, err
	}

	activeItems, err := s.store.CountActiveStripeFCItems(userID)
	if err != nil {
		return nil, err
	}

	apiCalls, err := s.store.GetStripeAPICallCount(userID, period.PeriodStart)
	if err != nil {
		return nil, err
	}

	return &models.PlaidUsage{
		SubscriptionTier: sub.Tier,
		Limits:           limits,
		ActiveItems:      activeItems,
		APICallsUsed:     apiCalls,
		PeriodStart:      period.PeriodStart,
		PeriodEnd:        period.PeriodEnd,
	}, nil
}

// ensureItemLimitAvailable verifies the user can link another bank connection.
func (s *Service) ensureItemLimitAvailable(userID int64) error {
	usage, err := s.GetUsage(userID)
	if err != nil {
		return err
	}
	if usage.ActiveItems >= usage.Limits.MaxItems {
		return ErrStripeItemLimitExceeded
	}
	return nil
}

// reserveStripeAPICall checks billing-period usage and records one call before hitting Stripe.
func (s *Service) reserveStripeAPICall(userID int64) error {
	usage, err := s.GetUsage(userID)
	if err != nil {
		return err
	}
	unlimited, err := s.subscription.HasUnlimitedLimits(userID)
	if err != nil {
		return err
	}
	if unlimited {
		return s.store.IncrementStripeAPICallCount(userID, usage.PeriodStart)
	}
	if usage.APICallsUsed >= usage.Limits.MaxAPICallsMonth {
		return ErrStripeAPILimitExceeded
	}
	return s.store.IncrementStripeAPICallCount(userID, usage.PeriodStart)
}
