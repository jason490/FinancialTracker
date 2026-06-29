package models

// PasswordResetCode is a single-use temporary code for password recovery.
type PasswordResetCode struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	CodeHash  string `json:"-"`
	ExpiresAt int64  `json:"expires_at"`
	Attempts  int    `json:"attempts"`
	UsedAt    *int64 `json:"used_at,omitempty"`
	CreatedAt int64  `json:"created_at"`
}
