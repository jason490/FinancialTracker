package middleware

import (
	"context"

	"github.com/labstack/echo/v5"
)

// maybeSyncStalePlaid runs a background sync for the user when any item is past the stale interval.
func (m *AuthMiddleware) maybeSyncStalePlaid(c *echo.Context) {
	if m.plaid == nil {
		return
	}
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return
	}
	userID, ok := userIDVal.(int64)
	if !ok {
		return
	}
	if !m.plaid.UserHasStaleItems(userID) {
		return
	}
	ctx := context.Background()
	go m.plaid.SyncStaleItems(ctx, userID)
}
