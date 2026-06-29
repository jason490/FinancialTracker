package queries

import (
	"FinancialTracker/internal/models"
	"database/sql"
)

// UpdateUserSubscription updates tier and billing fields. Pass nil pointers to leave a field unchanged.
func UpdateUserSubscription(db *sql.DB, userID int64, tier string, startedAt *int64, stripeCustomerID, stripeSubscriptionID *string) error {
	if startedAt == nil && tier == models.SubscriptionTierFree {
		zero := int64(0)
		startedAt = &zero
	}

	query := `UPDATE users
	          SET subscription_tier = ?,
	              subscription_started_at = CASE WHEN ? = 'free' THEN NULL ELSE ? END,
	              stripe_customer_id = CASE WHEN ? = 'free' THEN NULL ELSE COALESCE(?, stripe_customer_id) END,
	              stripe_subscription_id = CASE WHEN ? = 'free' THEN NULL ELSE COALESCE(?, stripe_subscription_id) END
	          WHERE id = ?`

	var startedValue any
	if tier == models.SubscriptionTierFree {
		startedValue = nil
	} else if startedAt != nil && *startedAt > 0 {
		startedValue = *startedAt
	} else {
		startedValue = nil
	}

	_, err := db.Exec(query, tier, tier, startedValue, tier, stripeCustomerID, tier, stripeSubscriptionID, userID)
	return err
}

// ApplyStripeSubscription updates billing fields from Stripe webhooks without clearing the customer ID.
func ApplyStripeSubscription(db *sql.DB, userID int64, tier string, startedAt *int64, subscriptionID *string) error {
	var startedValue any
	if tier == models.SubscriptionTierFree {
		startedValue = nil
	} else if startedAt != nil && *startedAt > 0 {
		startedValue = *startedAt
	} else {
		startedValue = nil
	}

	query := `UPDATE users
	          SET subscription_tier = ?,
	              subscription_started_at = ?,
	              stripe_subscription_id = CASE WHEN ? = 'free' THEN NULL ELSE COALESCE(?, stripe_subscription_id) END
	          WHERE id = ?`

	var subID any
	if subscriptionID != nil {
		subID = *subscriptionID
	}
	_, err := db.Exec(query, tier, startedValue, tier, subID, userID)
	return err
}
