package middleware

import (
	financialService "FinancialTracker/internal/services/financial"
	"FinancialTracker/internal/storage"

	"github.com/labstack/echo/v5"
)

// AuthMiddleware handles authentication for routes
type AuthMiddleware struct {
	store    *storage.Storage
	provider financialService.Provider
}

// NewAuthMiddleware initializes a new AuthMiddleware
func NewAuthMiddleware(store *storage.Storage, facade *financialService.Facade) *AuthMiddleware {
	return &AuthMiddleware{
		store:    store,
		provider: facade.Active(),
	}
}

// SessionMiddleware attempts to verify the session but does NOT redirect on failure.
func (m *AuthMiddleware) SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		session, userID, err := m.loadSession(c)
		if err != nil || session == nil {
			return next(c)
		}

		c.Set("user_id", userID)
		c.Set("session", session)

		m.maybeSyncStaleConnections(c)
		return next(c)
	}
}
