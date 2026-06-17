package external

// SettingsProfile is the authenticated user's settings payload.
type SettingsProfile struct {
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	Email           string   `json:"email"`
	ThemePreference string   `json:"theme_preference"`
	HasPassword     bool     `json:"has_password"`
	SSOProviders    []string `json:"sso_providers"`
}

// UpdateProfileRequest updates the user's display name.
type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UpdatePasswordRequest changes the user's password.
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

// UpdateThemeRequest persists a theme preference.
type UpdateThemeRequest struct {
	Theme string `json:"theme"`
}

// DeleteAccountVerifyRequest verifies credentials before account deletion.
type DeleteAccountVerifyRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// DeleteAccountReauthStatusResponse reports whether re-auth is still valid.
type DeleteAccountReauthStatusResponse struct {
	ReauthVerified bool `json:"reauth_verified"`
}
