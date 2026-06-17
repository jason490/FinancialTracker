package models

// User represents a user in the system
type User struct {
	ID              int64     `json:"id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	PasswordHash    string    `json:"-"`
	ThemePreference string    `json:"theme_preference"`
	CreatedAt       int64     `json:"created_at"`
	SSOs            []UserSSO `json:"sso_logins,omitempty"`
}

// UserSSO represents an SSO login linked to a user
type UserSSO struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Provider  string `json:"provider"`
	SSOID     string `json:"sso_id"`
	CreatedAt int64  `json:"created_at"`
}
