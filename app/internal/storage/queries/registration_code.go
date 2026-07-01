package queries

import (
	"FinancialTracker/internal/models"
	"database/sql"
	"time"
)

// CreateRegistrationCode stores a hashed invite code.
func CreateRegistrationCode(db *sql.DB, createdByUserID *int64, codeHash string, expiresAt int64) error {
	_, err := db.Exec(
		`INSERT INTO registration_codes (code_hash, created_by_user_id, expires_at) VALUES (?, ?, ?)`,
		codeHash, createdByUserID, expiresAt,
	)
	return err
}

// ListActiveRegistrationCodes returns unused, unexpired invite codes for validation.
func ListActiveRegistrationCodes(db *sql.DB) ([]*models.RegistrationCode, error) {
	rows, err := db.Query(
		`SELECT id, code_hash, created_by_user_id, expires_at, attempts, used_at, used_by_user_id, created_at
		 FROM registration_codes
		 WHERE used_at IS NULL AND expires_at > ?
		 ORDER BY created_at DESC`,
		time.Now().Unix(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []*models.RegistrationCode
	for rows.Next() {
		row, err := scanRegistrationCode(rows)
		if err != nil {
			return nil, err
		}
		codes = append(codes, row)
	}
	return codes, rows.Err()
}

// IncrementRegistrationCodeAttempts bumps the failed-attempt counter for an invite code.
func IncrementRegistrationCodeAttempts(db *sql.DB, id int64) error {
	_, err := db.Exec(`UPDATE registration_codes SET attempts = attempts + 1 WHERE id = ?`, id)
	return err
}

// MarkRegistrationCodeUsed marks an invite code as consumed by a new user.
func MarkRegistrationCodeUsed(db *sql.DB, id, usedByUserID int64) error {
	_, err := db.Exec(
		`UPDATE registration_codes SET used_at = ?, used_by_user_id = ? WHERE id = ?`,
		time.Now().Unix(), usedByUserID, id,
	)
	return err
}

func scanRegistrationCode(scanner interface {
	Scan(dest ...any) error
}) (*models.RegistrationCode, error) {
	row := &models.RegistrationCode{}
	var createdBy sql.NullInt64
	var usedAt sql.NullInt64
	var usedBy sql.NullInt64
	err := scanner.Scan(
		&row.ID, &row.CodeHash, &createdBy, &row.ExpiresAt,
		&row.Attempts, &usedAt, &usedBy, &row.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if createdBy.Valid {
		v := createdBy.Int64
		row.CreatedByUserID = &v
	}
	if usedAt.Valid {
		v := usedAt.Int64
		row.UsedAt = &v
	}
	if usedBy.Valid {
		v := usedBy.Int64
		row.UsedByUserID = &v
	}
	return row, nil
}
