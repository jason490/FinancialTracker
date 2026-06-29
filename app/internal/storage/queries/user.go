package queries

import (
    "FinancialTracker/internal/models"
    "database/sql"
)

// GetUserByEmail retrieves a user by their email address
func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
    user := &models.User{}
    query := `SELECT id, email, first_name, last_name, password_hash, theme_preference,
              subscription_tier, subscription_started_at, stripe_customer_id, stripe_subscription_id,
              onboarding_completed, created_at
              FROM users WHERE email = ?`
    var stripeCustomerID, stripeSubscriptionID sql.NullString
    var subscriptionStartedAt sql.NullInt64
    var onboardingCompleted int
    err := db.QueryRow(query, email).Scan(
        &user.ID, &user.Email, &user.FirstName, &user.LastName, &user.PasswordHash, &user.ThemePreference,
        &user.SubscriptionTier, &subscriptionStartedAt, &stripeCustomerID, &stripeSubscriptionID,
        &onboardingCompleted, &user.CreatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    if subscriptionStartedAt.Valid {
        user.SubscriptionStartedAt = subscriptionStartedAt.Int64
    }
    if stripeCustomerID.Valid {
        user.StripeCustomerID = stripeCustomerID.String
    }
    if stripeSubscriptionID.Valid {
        user.StripeSubscriptionID = stripeSubscriptionID.String
    }
    user.OnboardingCompleted = onboardingCompleted != 0
    
    ssoQuery := `SELECT id, user_id, provider, sso_id, created_at FROM user_sso WHERE user_id = ?`
    rows, err := db.Query(ssoQuery, user.ID)
    if err == nil {
        defer rows.Close()
        for rows.Next() {
            var sso models.UserSSO
            if err := rows.Scan(&sso.ID, &sso.UserID, &sso.Provider, &sso.SSOID, &sso.CreatedAt); err == nil {
                user.SSOs = append(user.SSOs, sso)
            }
        }
    }

    return user, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(db *sql.DB, id int64) (*models.User, error) {
    user := &models.User{}
    query := `SELECT id, email, first_name, last_name, password_hash, theme_preference,
              subscription_tier, subscription_started_at, stripe_customer_id, stripe_subscription_id,
              onboarding_completed, created_at
              FROM users WHERE id = ?`
    var stripeCustomerID, stripeSubscriptionID sql.NullString
    var subscriptionStartedAt sql.NullInt64
    var onboardingCompleted int
    err := db.QueryRow(query, id).Scan(
        &user.ID, &user.Email, &user.FirstName, &user.LastName, &user.PasswordHash, &user.ThemePreference,
        &user.SubscriptionTier, &subscriptionStartedAt, &stripeCustomerID, &stripeSubscriptionID,
        &onboardingCompleted, &user.CreatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    if subscriptionStartedAt.Valid {
        user.SubscriptionStartedAt = subscriptionStartedAt.Int64
    }
    if stripeCustomerID.Valid {
        user.StripeCustomerID = stripeCustomerID.String
    }
    if stripeSubscriptionID.Valid {
        user.StripeSubscriptionID = stripeSubscriptionID.String
    }
    user.OnboardingCompleted = onboardingCompleted != 0

    ssoQuery := `SELECT id, user_id, provider, sso_id, created_at FROM user_sso WHERE user_id = ?`
    rows, err := db.Query(ssoQuery, user.ID)
    if err == nil {
        defer rows.Close()
        for rows.Next() {
            var sso models.UserSSO
            if err := rows.Scan(&sso.ID, &sso.UserID, &sso.Provider, &sso.SSOID, &sso.CreatedAt); err == nil {
                user.SSOs = append(user.SSOs, sso)
            }
        }
    }

    return user, nil
}

// GetUserBySSO retrieves a user by their SSO provider and ID
func GetUserBySSO(db *sql.DB, provider, ssoID string) (*models.User, error) {
    var userID int64
    query := `SELECT user_id FROM user_sso WHERE provider = ? AND sso_id = ?`
    err := db.QueryRow(query, provider, ssoID).Scan(&userID)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return GetUserByID(db, userID)
}

// GetUserSubscriptionTier returns the subscription tier for a user, defaulting to free.
func GetUserSubscriptionTier(db *sql.DB, userID int64) (string, error) {
	query := `SELECT COALESCE(subscription_tier, 'free') FROM users WHERE id = ?`
	var tier string
	err := db.QueryRow(query, userID).Scan(&tier)
	if err == sql.ErrNoRows {
		return models.SubscriptionTierFree, nil
	}
	if err != nil {
		return "", err
	}
	if tier == "" {
		return models.SubscriptionTierFree, nil
	}
	return tier, nil
}

// CreateUser creates a new user in the database
func CreateUser(db *sql.DB, user *models.User) error {
	if user.ThemePreference == "" {
		user.ThemePreference = "system"
	}
	if user.SubscriptionTier == "" {
		user.SubscriptionTier = models.SubscriptionTierFree
	}
	query := `INSERT INTO users (email, first_name, last_name, password_hash, theme_preference, subscription_tier, onboarding_completed) VALUES (?, ?, ?, ?, ?, ?, 0)`
	result, err := db.Exec(query, user.Email, user.FirstName, user.LastName, user.PasswordHash, user.ThemePreference, user.SubscriptionTier)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	user.OnboardingCompleted = false
	return nil
}

// CompleteOnboarding marks the user's onboarding flow as finished.
func CompleteOnboarding(db *sql.DB, userID int64) error {
	query := `UPDATE users SET onboarding_completed = 1 WHERE id = ?`
	_, err := db.Exec(query, userID)
	return err
}

// LinkSSO links an SSO provider to a user
func LinkSSO(db *sql.DB, userID int64, provider, ssoID string) error {
    query := `INSERT INTO user_sso (user_id, provider, sso_id) VALUES (?, ?, ?)`
    _, err := db.Exec(query, userID, provider, ssoID)
    return err
}

// UnlinkSSO removes an SSO provider from a user
func UnlinkSSO(db *sql.DB, userID int64, provider string) error {
    query := `DELETE FROM user_sso WHERE user_id = ? AND provider = ?`
    _, err := db.Exec(query, userID, provider)
    return err
}

// UpdateUserTheme updates the theme preference for a user
func UpdateUserTheme(db *sql.DB, userID int64, theme string) error {
    query := `UPDATE users SET theme_preference = ? WHERE id = ?`
    _, err := db.Exec(query, theme, userID)
    return err
}

// UpdateUserInfo updates the basic information for a user
func UpdateUserInfo(db *sql.DB, userID int64, firstName, lastName, email string) error {
    query := `UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?`
    _, err := db.Exec(query, firstName, lastName, email, userID)
    return err
}

// UpdateUserPassword updates the password hash for a user
func UpdateUserPassword(db *sql.DB, userID int64, passwordHash string) error {
    query := `UPDATE users SET password_hash = ? WHERE id = ?`
    _, err := db.Exec(query, passwordHash, userID)
    return err
}

// DeleteUser removes a user from the database
func DeleteUser(db *sql.DB, userID int64) error {
    query := `DELETE FROM users WHERE id = ?`
    _, err := db.Exec(query, userID)
    return err
}
