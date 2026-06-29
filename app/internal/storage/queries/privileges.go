package queries

import (
	"FinancialTracker/internal/models"
	"database/sql"
)

// GetUserPrivileges loads per-user subscription overrides, if any.
func GetUserPrivileges(db *sql.DB, userID int64) (*models.UserPrivileges, error) {
	query := `SELECT user_id, unlimited_limits, stripe_coupon_id, notes
	          FROM user_privileges WHERE user_id = ?`
	var couponID, notes sql.NullString
	var unlimited int
	err := db.QueryRow(query, userID).Scan(&userID, &unlimited, &couponID, &notes)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	priv := &models.UserPrivileges{
		UserID:          userID,
		UnlimitedLimits: unlimited != 0,
	}
	if couponID.Valid {
		priv.StripeCouponID = couponID.String
	}
	if notes.Valid {
		priv.Notes = notes.String
	}
	return priv, nil
}

// UpsertUserPrivileges stores or updates per-user subscription overrides.
func UpsertUserPrivileges(db *sql.DB, priv *models.UserPrivileges) error {
	unlimited := 0
	if priv.UnlimitedLimits {
		unlimited = 1
	}
	query := `INSERT INTO user_privileges (user_id, unlimited_limits, stripe_coupon_id, notes)
	          VALUES (?, ?, ?, ?)
	          ON CONFLICT(user_id) DO UPDATE SET
	            unlimited_limits = excluded.unlimited_limits,
	            stripe_coupon_id = excluded.stripe_coupon_id,
	            notes = excluded.notes`
	_, err := db.Exec(query, priv.UserID, unlimited, nullString(priv.StripeCouponID), nullString(priv.Notes))
	return err
}

// GetUserByStripeCustomerID finds a user by their Stripe customer ID.
func GetUserByStripeCustomerID(db *sql.DB, customerID string) (*models.User, error) {
	var userID int64
	query := `SELECT id FROM users WHERE stripe_customer_id = ?`
	err := db.QueryRow(query, customerID).Scan(&userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return GetUserByID(db, userID)
}

func nullString(value string) any {
	if value == "" {
		return nil
	}
	return value
}
