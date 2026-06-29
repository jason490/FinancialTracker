package middleware

import (
	"FinancialTracker/internal/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

// loadSession reads and validates the session cookie, returning the session and user ID.
func (m *AuthMiddleware) loadSession(c *echo.Context) (*models.Session, int64, error) {
	cookie, err := c.Cookie("Session")
	if err != nil {
		return nil, 0, err
	}

	session, err := m.store.GetSession(cookie.Value)
	if err != nil || session == nil {
		return nil, 0, err
	}

	if session.ExpiresAt < time.Now().Unix() {
		_ = m.store.DeleteSession(cookie.Value)
		return nil, 0, http.ErrNoCookie
	}

	return session, session.UserID, nil
}
