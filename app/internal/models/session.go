package models

// Session represents an active user session
type Session struct {
	ID                string `json:"id"`
	UserID            int64  `json:"user_id"`
	ExpiresAt         int64  `json:"expires_at"`
	ReauthenticatedAt int64  `json:"reauthenticated_at"`
	CreatedAt         int64  `json:"created_at"`
}
