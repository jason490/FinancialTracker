package middleware

import (
	plaidService "FinancialTracker/internal/services/plaid"
	"FinancialTracker/internal/storage"
	"github.com/labstack/echo/v5"
	"net/http"
	"time"
)

// AuthMiddleware handles authentication for routes
type AuthMiddleware struct {
	store *storage.Storage
	plaid *plaidService.PlaidService
}

// NewAuthMiddleware initializes a new AuthMiddleware
func NewAuthMiddleware(store *storage.Storage, plaid *plaidService.PlaidService) *AuthMiddleware {
	return &AuthMiddleware{
		store: store,
	}
}

// AuthMiddleware verifies the session from the database
func (m *AuthMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c *echo.Context) error {
        cookie, err := c.Cookie("Session")
        if err != nil {
            return c.Redirect(http.StatusFound, "/login")
        }

        sessionID := cookie.Value
        session, err := m.store.GetSession(sessionID)
        if err != nil || session == nil {
            return c.Redirect(http.StatusFound, "/login")
        }

        // Check if session is expired
        if session.ExpiresAt < time.Now().Unix() {
            _ = m.store.DeleteSession(sessionID)
            return c.Redirect(http.StatusFound, "/login")
        }

        // Set user ID in context for downstream handlers
        c.Set("user_id", session.UserID)
        c.Set("session", session)

        m.maybeSyncStalePlaid(c)
        return next(c)
    }
}

// SessionMiddleware attempts to verify the session but does NOT redirect on failure.
// Used for public pages that want to know if a user is logged in (e.g. Home).
func (m *AuthMiddleware) SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		cookie, err := c.Cookie("Session")
		if err != nil {
			return next(c)
		}

		sessionID := cookie.Value
		session, err := m.store.GetSession(sessionID)
		if err != nil || session == nil {
			return next(c)
		}

		// Check if session is expired
		if session.ExpiresAt < time.Now().Unix() {
			_ = m.store.DeleteSession(sessionID)
			return next(c)
		}

		// Set user ID in context for downstream handlers
		c.Set("user_id", session.UserID)
		c.Set("session", session)

		m.maybeSyncStalePlaid(c)
		return next(c)
	}
}
