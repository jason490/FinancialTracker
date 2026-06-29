package financial

import (
	"context"

	"FinancialTracker/internal/models/external"

	"github.com/labstack/echo/v5"
)

// Provider defines bank connection operations for Plaid and Stripe Financial Connections.
type Provider interface {
	Name() string
	GetManagementData(userID int64) (*external.ConnectionsPayload, error)
	CreateSession(c *echo.Context, userID int64) (*external.CreateSessionResponse, error)
	CreateUpdateSession(c *echo.Context, userID int64, rowID string) (*external.CreateSessionResponse, error)
	CompleteConnection(c *echo.Context, userID int64, req *external.CompleteConnectionRequest) error
	SyncUser(c *echo.Context, userID int64) error
	SyncItem(c *echo.Context, userID int64, rowID string) error
	Disconnect(c *echo.Context, userID int64, rowID string) error
	ToggleAccountVisibility(userID int64, accountID string) (bool, error)
	RemoveAccount(userID int64, accountID string) error
	UserHasStaleItems(userID int64) bool
	SyncStaleItems(ctx context.Context, userID int64)
}
