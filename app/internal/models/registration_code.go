package models

// RegistrationCode is a single-use invite code for gated registration.
type RegistrationCode struct {
	ID              int64  `json:"id"`
	CodeHash        string `json:"-"`
	CreatedByUserID *int64 `json:"created_by_user_id,omitempty"`
	ExpiresAt       int64  `json:"expires_at"`
	Attempts        int    `json:"attempts"`
	UsedAt          *int64 `json:"used_at,omitempty"`
	UsedByUserID    *int64 `json:"used_by_user_id,omitempty"`
	CreatedAt       int64  `json:"created_at"`
}
