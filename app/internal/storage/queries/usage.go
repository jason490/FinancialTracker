package queries

import (
	"database/sql"
)

// GetPlaidAPICallCount returns how many Plaid API calls a user has made in the given period.
func GetPlaidAPICallCount(db *sql.DB, userID int64, periodStart int64) (int, error) {
	query := `SELECT call_count FROM plaid_api_usage WHERE user_id = ? AND period_start = ?`
	var count int
	err := db.QueryRow(query, userID, periodStart).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return count, err
}

// IncrementPlaidAPICallCount atomically increments the user's Plaid API call counter for a period.
func IncrementPlaidAPICallCount(db *sql.DB, userID int64, periodStart int64) error {
	query := `INSERT INTO plaid_api_usage (user_id, period_start, call_count)
	          VALUES (?, ?, 1)
	          ON CONFLICT(user_id, period_start) DO UPDATE SET call_count = call_count + 1`
	_, err := db.Exec(query, userID, periodStart)
	return err
}

// GetStripeAPICallCount returns how many Stripe API calls a user has made in the given period.
func GetStripeAPICallCount(db *sql.DB, userID int64, periodStart int64) (int, error) {
	query := `SELECT call_count FROM stripe_api_usage WHERE user_id = ? AND period_start = ?`
	var count int
	err := db.QueryRow(query, userID, periodStart).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return count, err
}

// IncrementStripeAPICallCount atomically increments the user's Stripe API call counter for a period.
func IncrementStripeAPICallCount(db *sql.DB, userID int64, periodStart int64) error {
	query := `INSERT INTO stripe_api_usage (user_id, period_start, call_count)
	          VALUES (?, ?, 1)
	          ON CONFLICT(user_id, period_start) DO UPDATE SET call_count = call_count + 1`
	_, err := db.Exec(query, userID, periodStart)
	return err
}
