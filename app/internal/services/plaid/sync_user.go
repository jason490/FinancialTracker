package plaid

import (
	"context"
	"errors"
	"sync"
	"time"

	"FinancialTracker/internal/models"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
)

// StaleSyncInterval is how long before an item is considered due for automatic sync.
const StaleSyncInterval = 24 * time.Hour

var userSyncInProgress sync.Map

// itemNeedsStaleSync reports whether an item should be included in automatic background sync.
func itemNeedsStaleSync(item models.PlaidItem, cutoff int64) bool {
	if item.Status == ItemStatusDisconnected {
		return false
	}
	return item.LastSynced == 0 || item.LastSynced < cutoff
}

// UserHasStaleItems reports whether the user has any non-disconnected items due for sync.
func (p *PlaidService) UserHasStaleItems(userID int64) bool {
	items, err := p.store.GetPlaidItemsByUserID(userID)
	if err != nil {
		return false
	}
	cutoff := time.Now().Add(-StaleSyncInterval).Unix()
	for _, item := range items {
		if itemNeedsStaleSync(item, cutoff) {
			return true
		}
	}
	return false
}

// SyncUser synchronizes accounts and transactions for all non-disconnected items (manual sync).
func (p *PlaidService) SyncUser(c *echo.Context, userID int64) error {
	if err := p.beginManualSync(userID); err != nil {
		return err
	}
	ctx := c.Request().Context()
	return p.syncItems(&ctx, userID, false)
}

// SyncStaleItems synchronizes only items whose last_synced is older than StaleSyncInterval.
func (p *PlaidService) SyncStaleItems(ctx context.Context, userID int64) {
	if _, loaded := userSyncInProgress.LoadOrStore(userID, true); loaded {
		return
	}
	defer userSyncInProgress.Delete(userID)

	if err := p.syncItems(&ctx, userID, true); err != nil {
		log.Errorf("Background stale sync failed for user %d: %v", userID, err)
	}
}

// syncItems runs account and transaction sync for a user's Plaid items.
func (p *PlaidService) syncItems(ctx *context.Context, userID int64, staleOnly bool) error {
	items, err := p.store.GetPlaidItemsByUserID(userID)
	if err != nil {
		log.Errorf("Failed to fetch Plaid items for user %d: %v", userID, err)
		return err
	}
	if len(items) <= 0 {
		return errors.New("No banks connected, please connect to a bank in the manage page")
	}

	cutoff := time.Now().Add(-StaleSyncInterval).Unix()
	now := time.Now().Unix()
	syncedAny := false

	for _, item := range items {
		if item.Status == ItemStatusDisconnected {
			continue
		}
		if staleOnly && !itemNeedsStaleSync(item, cutoff) {
			continue
		}

		if err := p.SyncPlaidAccounts(ctx, userID, item.PlaidItemID, item.AccessToken); err != nil {
			log.Errorf("Failed to sync accounts for item %s (user %d): %v", item.PlaidItemID, userID, err)
			continue
		}

		itemRow, err := p.store.GetPlaidItemByItemID(item.PlaidItemID)
		if err != nil {
			log.Errorf("Failed to refresh item %s after account sync: %v", item.PlaidItemID, err)
			continue
		}
		if itemRow.Status == ItemStatusDisconnected {
			continue
		}

		if err := p.syncItemTransactions(*ctx, userID, item.PlaidItemID, item.AccessToken, itemRow.SyncCursor); err != nil {
			log.Errorf("Failed to sync transactions for item %s (user %d): %v", item.PlaidItemID, userID, err)
			continue
		}

		if err := p.store.UpdatePlaidItemLastSynced(item.PlaidItemID, now); err != nil {
			log.Errorf("Failed to update last_synced for item %s: %v", item.PlaidItemID, err)
		}
		syncedAny = true
	}

	if !staleOnly && !syncedAny && len(items) > 0 {
		allDisconnected := true
		for _, item := range items {
			if item.Status != ItemStatusDisconnected {
				allDisconnected = false
				break
			}
		}
		if allDisconnected {
			return nil
		}
	}

	return nil
}
