package plaid

import (
	"FinancialTracker/internal/models"
	subscriptionService "FinancialTracker/internal/services/subscription"
)

// LimitsForTier returns Plaid quotas for a subscription tier.
func LimitsForTier(tier string) models.PlaidLimits {
	return subscriptionService.LimitsForTier(tier)
}
