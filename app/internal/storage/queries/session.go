package queries

import (
    "FinancialTracker/internal/models"
    "database/sql"
)

// CreateSession inserts a new session into the database
func CreateSession(db *sql.DB, session *models.Session) error {
    query := `INSERT INTO sessions (id, user_id, expires_at, reauthenticated_at) VALUES (?, ?, ?, ?)`
    _, err := db.Exec(query, session.ID, session.UserID, session.ExpiresAt, session.ReauthenticatedAt)
    return err
}

// GetSession retrieves a session by its ID
func GetSession(db *sql.DB, id string) (*models.Session, error) {
    session := &models.Session{}
    query := `SELECT id, user_id, expires_at, reauthenticated_at, created_at FROM sessions WHERE id = ?`
    err := db.QueryRow(query, id).Scan(&session.ID, &session.UserID, &session.ExpiresAt, &session.ReauthenticatedAt, &session.CreatedAt)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return session, nil
}

// UpdateSessionReauth updates the re-authentication timestamp
func UpdateSessionReauth(db *sql.DB, id string, timestamp int64) error {
    query := `UPDATE sessions SET reauthenticated_at = ? WHERE id = ?`
    _, err := db.Exec(query, timestamp, id)
    return err
}

// DeleteSession removes a session from the database
func DeleteSession(db *sql.DB, id string) error {
    query := `DELETE FROM sessions WHERE id = ?`
    _, err := db.Exec(query, id)
    return err
}

// DeleteUserSessions removes all sessions for a user.
func DeleteUserSessions(db *sql.DB, userID int64) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE user_id = ?`, userID)
	return err
}
