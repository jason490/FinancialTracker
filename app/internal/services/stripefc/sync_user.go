package stripefc

import (
	"context"
	"errors"
	"sync"
	"time"

	"FinancialTracker/internal/models"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/stripe/stripe-go/v86"
)

var userSyncInProgress sync.Map

func itemNeedsStaleSync(item models.StripeFCItem, cutoff int64) bool {
	if item.Status == ItemStatusDisconnected {
		return false
	}
	return item.LastSynced == 0 || item.LastSynced < cutoff
}

// UserHasStaleItems reports whether the user has Stripe FC items due for sync.
func (s *Service) UserHasStaleItems(userID int64) bool {
	items, err := s.store.GetStripeFCItemsByUserID(userID)
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

// SyncStaleItems synchronizes only items older than StaleSyncInterval.
func (s *Service) SyncStaleItems(ctx context.Context, userID int64) {
	if _, loaded := userSyncInProgress.LoadOrStore(userID, true); loaded {
		return
	}
	defer userSyncInProgress.Delete(userID)

	if err := s.syncItems(ctx, userID, true); err != nil {
		log.Errorf("Background stale Stripe sync failed for user %d: %v", userID, err)
	}
}

// StaleSyncInterval is how long before a Stripe FC item is considered due for automatic sync.
const StaleSyncInterval = 24 * time.Hour

func (s *Service) syncItems(ctx context.Context, userID int64, staleOnly bool) error {
	items, err := s.store.GetStripeFCItemsByUserID(userID)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		if staleOnly {
			return nil
		}
		return errors.New("no banks connected, please connect to a bank in the manage page")
	}

	cutoff := time.Now().Add(-StaleSyncInterval).Unix()
	now := time.Now().Unix()
	syncedAny := false
	var limitErr error

	for _, item := range items {
		if item.Status == ItemStatusDisconnected {
			continue
		}
		if staleOnly && !itemNeedsStaleSync(item, cutoff) {
			continue
		}

		if err := s.syncStripeAccounts(ctx, userID, item.RowID); err != nil {
			log.Errorf("Failed to sync Stripe accounts for item %s: %v", item.RowID, err)
			if errors.Is(err, ErrStripeAPILimitExceeded) {
				limitErr = err
			}
			continue
		}

		if err := s.syncItemTransactions(ctx, userID, item.RowID); err != nil {
			log.Errorf("Failed to sync Stripe transactions for item %s: %v", item.RowID, err)
			if errors.Is(err, ErrStripeAPILimitExceeded) {
				limitErr = err
			}
			continue
		}

		if err := s.store.UpdateStripeFCItemLastSynced(item.RowID, now); err != nil {
			log.Errorf("Failed to update last_synced for Stripe item %s: %v", item.RowID, err)
		}
		syncedAny = true
	}

	if !staleOnly && limitErr != nil {
		return limitErr
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

// Disconnect removes a Stripe FC institution and disconnects its accounts at Stripe.
func (s *Service) Disconnect(c *echo.Context, userID int64, rowID string) error {
	item, err := s.store.GetStripeFCItemByRowID(rowID, userID)
	if err != nil {
		return errors.New("failed to find bank connection")
	}

	ctx := c.Request().Context()
	accounts, err := s.store.GetStripeFCAccountsByItemRowID(item.RowID)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if err := s.reserveStripeAPICall(userID); err != nil {
			log.Errorf("Stripe API limit reached while disconnecting account %s: %v", account.StripeAccountID, err)
			break
		}
		if _, err := s.client.V1FinancialConnectionsAccounts.Disconnect(ctx, account.StripeAccountID, &stripe.FinancialConnectionsAccountDisconnectParams{}); err != nil {
			log.Errorf("Stripe disconnect failed for account %s: %v", account.StripeAccountID, err)
		}
	}

	return s.store.DeleteStripeFCItem(rowID, userID)
}
