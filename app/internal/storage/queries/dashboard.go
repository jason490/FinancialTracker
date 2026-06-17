package queries

import (
	"FinancialTracker/internal/models"
	"database/sql"
	"encoding/json"
	"errors"
)

// GetDashboardLayout loads a user's saved dashboard layout, or nil if none exists.
func GetDashboardLayout(db *sql.DB, userID int64) (*models.DashboardLayout, error) {
	query := `SELECT layout_json FROM dashboard_layouts WHERE user_id = ?`
	var raw string
	err := db.QueryRow(query, userID).Scan(&raw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	var layout models.DashboardLayout
	if err := json.Unmarshal([]byte(raw), &layout); err != nil {
		return nil, err
	}
	return &layout, nil
}

// UpsertDashboardLayout saves or updates a user's dashboard layout.
func UpsertDashboardLayout(db *sql.DB, userID int64, layout *models.DashboardLayout) error {
	data, err := json.Marshal(layout)
	if err != nil {
		return err
	}
	query := `INSERT INTO dashboard_layouts (user_id, layout_json, updated_at)
              VALUES (?, ?, strftime('%s', 'now'))
              ON CONFLICT(user_id) DO UPDATE SET
                layout_json = excluded.layout_json,
                updated_at = strftime('%s', 'now')`
	_, err = db.Exec(query, userID, string(data))
	return err
}
