package plaid

import (
	"errors"

	"FinancialTracker/internal/models"
	subscriptionService "FinancialTracker/internal/services/subscription"
)

// GetPlaidUsage returns current Plaid consumption and limits for a user.
func (p *PlaidService) GetPlaidUsage(userID int64) (*models.PlaidUsage, error) {
	sub, err := p.subscription.GetUserSubscription(userID)
	if err != nil {
		return nil, err
	}

	period, err := p.subscription.CurrentBillingPeriod(userID)
	if err != nil {
		return nil, err
	}

	limits, err := p.subscription.EffectiveLimits(userID)
	if err != nil {
		return nil, err
	}

	activeItems, err := p.store.CountActivePlaidItems(userID)
	if err != nil {
		return nil, err
	}

	apiCalls, err := p.store.GetPlaidAPICallCount(userID, period.PeriodStart)
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
func (p *PlaidService) ensureItemLimitAvailable(userID int64) error {
	usage, err := p.GetPlaidUsage(userID)
	if err != nil {
		return err
	}
	if usage.ActiveItems >= usage.Limits.MaxItems {
		return ErrPlaidItemLimitExceeded
	}
	return nil
}

// reservePlaidAPICall checks billing-period usage and records one call before
// hitting Plaid. Delegates to the shared subscription quota helper so Plaid
// syncs and other metered features (e.g. CSV exports) share the same counter.
func (p *PlaidService) reservePlaidAPICall(userID int64) error {
	if err := p.subscription.ReserveAPICall(userID); err != nil {
		if errors.Is(err, subscriptionService.ErrAPILimitExceeded) {
			return ErrPlaidAPILimitExceeded
		}
		return err
	}
	return nil
}
