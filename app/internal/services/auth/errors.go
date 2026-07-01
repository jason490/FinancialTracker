package auth

import "errors"

var (
	// ErrInvalidExchangeToken is returned when an SSO exchange token is invalid or expired.
	ErrInvalidExchangeToken = errors.New("invalid or expired exchange token")
	// ErrInvalidOAuthState is returned when an OAuth state parameter cannot be parsed.
	ErrInvalidOAuthState = errors.New("invalid oauth state")
	// ErrSSOAccountConflict is returned when a Google sign-in matches an existing
	// account that has not linked that SSO provider (e.g. a password-only account).
	ErrSSOAccountConflict = errors.New("an account with this email already exists; sign in and link Google from settings")
)
