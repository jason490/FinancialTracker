package auth

import "errors"

var (
	// ErrInvalidExchangeToken is returned when an SSO exchange token is invalid or expired.
	ErrInvalidExchangeToken = errors.New("invalid or expired exchange token")
	// ErrInvalidOAuthState is returned when an OAuth state parameter cannot be parsed.
	ErrInvalidOAuthState = errors.New("invalid oauth state")
)
