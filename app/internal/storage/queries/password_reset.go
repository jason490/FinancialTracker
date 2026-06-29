package queries

import (
	"FinancialTracker/internal/models"
	"database/sql"
	"time"
)

// InvalidatePasswordResetCodes marks unused reset codes as consumed so issuance history is retained.
func InvalidatePasswordResetCodes(db *sql.DB, userID int64) error {
	now := time.Now().Unix()
	_, err := db.Exec(
		`UPDATE password_reset_codes SET used_at = ? WHERE user_id = ? AND used_at IS NULL`,
		now, userID,
	)
	return err
}

// CreatePasswordResetCode stores a hashed reset code for a user.
func CreatePasswordResetCode(db *sql.DB, userID int64, codeHash string, expiresAt int64) error {
	_, err := db.Exec(
		`INSERT INTO password_reset_codes (user_id, code_hash, expires_at) VALUES (?, ?, ?)`,
		userID, codeHash, expiresAt,
	)
	return err
}

// GetActivePasswordResetCode returns the latest unused, unexpired reset code for a user.
func GetActivePasswordResetCode(db *sql.DB, userID int64) (*models.PasswordResetCode, error) {
	row := &models.PasswordResetCode{}
	var usedAt sql.NullInt64
	err := db.QueryRow(
		`SELECT id, user_id, code_hash, expires_at, attempts, used_at, created_at
		 FROM password_reset_codes
		 WHERE user_id = ? AND used_at IS NULL AND expires_at > ?
		 ORDER BY created_at DESC
		 LIMIT 1`,
		userID, time.Now().Unix(),
	).Scan(&row.ID, &row.UserID, &row.CodeHash, &row.ExpiresAt, &row.Attempts, &usedAt, &row.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if usedAt.Valid {
		v := usedAt.Int64
		row.UsedAt = &v
	}
	return row, nil
}

// IncrementPasswordResetAttempts bumps the failed-attempt counter for a reset code.
func IncrementPasswordResetAttempts(db *sql.DB, id int64) error {
	_, err := db.Exec(`UPDATE password_reset_codes SET attempts = attempts + 1 WHERE id = ?`, id)
	return err
}

// MarkPasswordResetCodeUsed marks a reset code as consumed.
func MarkPasswordResetCodeUsed(db *sql.DB, id int64) error {
	_, err := db.Exec(`UPDATE password_reset_codes SET used_at = ? WHERE id = ?`, time.Now().Unix(), id)
	return err
}

// GetLatestPasswordResetCreatedAt returns when the most recent reset code was created for a user.
func GetLatestPasswordResetCreatedAt(db *sql.DB, userID int64) (int64, error) {
	var createdAt sql.NullInt64
	err := db.QueryRow(
		`SELECT created_at FROM password_reset_codes WHERE user_id = ? ORDER BY created_at DESC LIMIT 1`,
		userID,
	).Scan(&createdAt)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	if !createdAt.Valid {
		return 0, nil
	}
	return createdAt.Int64, nil
}

// CountPasswordResetCodesSince returns how many reset codes were issued since a unix timestamp.
func CountPasswordResetCodesSince(db *sql.DB, userID int64, since int64) (int, error) {
	var count int
	err := db.QueryRow(
		`SELECT COUNT(*) FROM password_reset_codes WHERE user_id = ? AND created_at >= ?`,
		userID, since,
	).Scan(&count)
	return count, err
}
