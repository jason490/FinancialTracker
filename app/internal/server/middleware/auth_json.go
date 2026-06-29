package middleware

import (
	"FinancialTracker/internal/server/handler"
	"net/http"

	"github.com/labstack/echo/v5"
)

// AuthMiddlewareJSON verifies the session and returns JSON errors when unauthenticated.
func (m *AuthMiddleware) AuthMiddlewareJSON(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		session, userID, err := m.loadSession(c)
		if err != nil || session == nil {
			if err == http.ErrNoCookie {
				return c.JSON(http.StatusUnauthorized, handler.ErrorResponse("unauthorized", "Session expired"))
			}
			return c.JSON(http.StatusUnauthorized, handler.ErrorResponse("unauthorized", "Authentication required"))
		}

		c.Set("user_id", userID)
		c.Set("session", session)

		m.maybeSyncStaleConnections(c)
		return next(c)
	}
}
