package plaid

import (
	"errors"
	"time"

	"github.com/labstack/echo/v5"
)

// SyncItemManual syncs accounts and transactions for one connection after a user action.
func (p *PlaidService) SyncItemManual(c *echo.Context, userID int64, rowID string) error {
	if err := p.beginManualSync(userID); err != nil {
		return err
	}

	item, err := p.store.GetPlaidItemByRowID(rowID, userID)
	if err != nil {
		return errors.New("connection not found")
	}
	if item.Status == ItemStatusDisconnected {
		return errors.New("this bank connection is no longer available; please disconnect and link again")
	}

	ctx := c.Request().Context()
	if err := p.SyncPlaidAccounts(&ctx, userID, item.PlaidItemID, item.AccessToken); err != nil {
		return err
	}
	if err := p.syncItemTransactions(ctx, userID, item.PlaidItemID, item.AccessToken, item.SyncCursor); err != nil {
		return err
	}

	if err := p.store.UpdatePlaidItemLastSynced(item.PlaidItemID, time.Now().Unix()); err != nil {
		return err
	}
	return nil
}
