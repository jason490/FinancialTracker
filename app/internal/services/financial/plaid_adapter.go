package financial

import (
	"context"
	"errors"

	"FinancialTracker/internal/models/external"
	connectionService "FinancialTracker/internal/services/connections"
	plaidService "FinancialTracker/internal/services/plaid"

	"github.com/labstack/echo/v5"
)

// PlaidAdapter wraps PlaidService to implement Provider.
type PlaidAdapter struct {
	svc *plaidService.PlaidService
}

// NewPlaidAdapter creates a Provider backed by Plaid.
func NewPlaidAdapter(svc *plaidService.PlaidService) *PlaidAdapter {
	return &PlaidAdapter{svc: svc}
}

// Name returns the provider identifier.
func (a *PlaidAdapter) Name() string {
	return ProviderPlaid
}

// GetManagementData returns Plaid connections in the unified payload shape.
func (a *PlaidAdapter) GetManagementData(userID int64) (*external.ConnectionsPayload, error) {
	data, usage, err := a.svc.GetManagementData(userID)
	if err != nil {
		return nil, err
	}
	return connectionService.BuildPayloadFromPlaid(data, usage), nil
}

// CreateSession creates a Plaid Link token for a new connection.
func (a *PlaidAdapter) CreateSession(c *echo.Context, userID int64) (*external.CreateSessionResponse, error) {
	token, err := a.svc.CreateLinkToken(c, userID)
	if err != nil {
		return nil, err
	}
	return &external.CreateSessionResponse{LinkToken: token}, nil
}

// CreateUpdateSession creates a Plaid Link token in update mode.
func (a *PlaidAdapter) CreateUpdateSession(c *echo.Context, userID int64, rowID string) (*external.CreateSessionResponse, error) {
	item, err := a.svc.GetItemByRowID(rowID, userID)
	if err != nil {
		return nil, err
	}
	token, err := a.svc.CreateUpdateLinkToken(c, userID, item.AccessToken, item.Status)
	if err != nil {
		return nil, err
	}
	return &external.CreateSessionResponse{LinkToken: token}, nil
}

// CompleteConnection exchanges a Plaid public token.
func (a *PlaidAdapter) CompleteConnection(c *echo.Context, userID int64, req *external.CompleteConnectionRequest) error {
	if req == nil || req.PublicToken == "" {
		return errors.New("public token is required")
	}
	return a.svc.ExchangeToken(c, userID, req.PublicToken)
}

// SyncUser syncs all Plaid connections for the user.
func (a *PlaidAdapter) SyncUser(c *echo.Context, userID int64) error {
	return a.svc.SyncUser(c, userID)
}

// SyncItem syncs a single Plaid institution.
func (a *PlaidAdapter) SyncItem(c *echo.Context, userID int64, rowID string) error {
	return a.svc.SyncItemManual(c, userID, rowID)
}

// Disconnect removes a Plaid institution connection.
func (a *PlaidAdapter) Disconnect(c *echo.Context, userID int64, rowID string) error {
	ctx := c.Request().Context()
	return a.svc.DisconnectItem(&ctx, rowID, userID)
}

// ToggleAccountVisibility flips whether an account is hidden.
func (a *PlaidAdapter) ToggleAccountVisibility(userID int64, accountID string) (bool, error) {
	return a.svc.ToggleAccountVisibility(userID, accountID)
}

// RemoveAccount permanently deletes a disconnected account.
func (a *PlaidAdapter) RemoveAccount(userID int64, accountID string) error {
	return a.svc.DeletePlaidAccount(userID, accountID)
}

// UserHasStaleItems reports whether background sync should run.
func (a *PlaidAdapter) UserHasStaleItems(userID int64) bool {
	return a.svc.UserHasStaleItems(userID)
}

// SyncStaleItems runs background sync for stale Plaid items.
func (a *PlaidAdapter) SyncStaleItems(ctx context.Context, userID int64) {
	a.svc.SyncStaleItems(ctx, userID)
}
