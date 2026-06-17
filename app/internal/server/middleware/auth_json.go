package middleware

import (
	"FinancialTracker/internal/server/handler"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

// AuthMiddlewareJSON verifies the session and returns JSON errors when unauthenticated.
func (m *AuthMiddleware) AuthMiddlewareJSON(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		cookie, err := c.Cookie("Session")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, handler.ErrorResponse("unauthorized", "Authentication required"))
		}

		sessionID := cookie.Value
		session, err := m.store.GetSession(sessionID)
		if err != nil || session == nil {
			return c.JSON(http.StatusUnauthorized, handler.ErrorResponse("unauthorized", "Authentication required"))
		}

		if session.ExpiresAt < time.Now().Unix() {
			_ = m.store.DeleteSession(sessionID)
			return c.JSON(http.StatusUnauthorized, handler.ErrorResponse("unauthorized", "Session expired"))
		}

		c.Set("user_id", session.UserID)
		c.Set("session", session)

		m.maybeSyncStalePlaid(c)
		return next(c)
	}
}
