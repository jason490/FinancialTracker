package middleware

import (
	"context"

	"github.com/labstack/echo/v5"
)

// maybeSyncStaleConnections runs a background sync when any connection is past the stale interval.
func (m *AuthMiddleware) maybeSyncStaleConnections(c *echo.Context) {
	if m.provider == nil {
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
	if !m.provider.UserHasStaleItems(userID) {
		return
	}
	ctx := context.Background()
	go m.provider.SyncStaleItems(ctx, userID)
}
