package stripefc

import (
	"context"
	"errors"
	"time"

	"FinancialTracker/internal/models"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/stripe/stripe-go/v86"
)

// SyncUser synchronizes all Stripe FC connections for a user.
func (s *Service) SyncUser(c *echo.Context, userID int64) error {
	return s.syncItems(c.Request().Context(), userID, false)
}

// SyncItem synchronizes a single Stripe FC institution.
func (s *Service) SyncItem(c *echo.Context, userID int64, rowID string) error {
	item, err := s.store.GetStripeFCItemByRowID(rowID, userID)
	if err != nil {
		return errors.New("connection not found")
	}
	if item.Status == ItemStatusDisconnected {
		return errors.New("this bank connection is no longer available; please disconnect and link again")
	}

	ctx := c.Request().Context()
	if err := s.syncStripeAccounts(ctx, userID, rowID); err != nil {
		return err
	}
	if err := s.syncItemTransactions(ctx, userID, rowID); err != nil {
		return err
	}
	return s.store.UpdateStripeFCItemLastSynced(rowID, time.Now().Unix())
}

func (s *Service) syncItemTransactions(ctx context.Context, userID int64, itemRowID string) error {
	item, err := s.store.GetStripeFCItemByRowID(itemRowID, userID)
	if err != nil {
		return err
	}

	accounts, err := s.store.GetStripeFCAccountsByItemRowID(itemRowID)
	if err != nil {
		return err
	}

	var latestRefresh string
	for _, account := range accounts {
		if account.Status == ItemStatusDisconnected {
			continue
		}
		if err := s.refreshAccountTransactions(ctx, userID, account.StripeAccountID); err != nil {
			log.Errorf("Failed to refresh transactions for account %s: %v", account.StripeAccountID, err)
			continue
		}
		refreshID, err := s.ingestAccountTransactions(ctx, userID, account, item.TransactionRefreshID)
		if err != nil {
			log.Errorf("Failed to ingest transactions for account %s: %v", account.StripeAccountID, err)
			continue
		}
		if refreshID != "" {
			latestRefresh = refreshID
		}
	}

	if latestRefresh != "" && latestRefresh != item.TransactionRefreshID {
		if err := s.store.UpdateStripeFCItemTransactionRefresh(itemRowID, latestRefresh); err != nil {
			log.Errorf("Failed to update transaction refresh for item %s: %v", itemRowID, err)
		}
	}

	return s.store.UpdateStripeFCItemStatus(itemRowID, ItemStatusActive, "")
}

func (s *Service) refreshAccountTransactions(ctx context.Context, userID int64, accountID string) error {
	if err := s.reserveStripeAPICall(userID); err != nil {
		return err
	}
	_, err := s.client.V1FinancialConnectionsAccounts.Refresh(ctx, accountID, &stripe.FinancialConnectionsAccountRefreshParams{
		Features: []*string{stripe.String("transactions")},
	})
	return err
}

func (s *Service) ingestAccountTransactions(ctx context.Context, userID int64, account models.StripeFCAccount, afterRefresh string) (string, error) {
	params := &stripe.FinancialConnectionsTransactionListParams{
		Account: stripe.String(account.StripeAccountID),
	}
	params.Limit = stripe.Int64(100)
	if afterRefresh != "" {
		params.TransactionRefresh = &stripe.FinancialConnectionsTransactionListTransactionRefreshParams{
			After: stripe.String(afterRefresh),
		}
	}

	if err := s.reserveStripeAPICall(userID); err != nil {
		return "", err
	}

	var latestRefresh string
	iter := s.client.V1FinancialConnectionsTransactions.List(ctx, params)
	for txn, err := range iter.All(ctx) {
		if err != nil {
			return latestRefresh, err
		}
		if txn.TransactionRefresh != "" {
			latestRefresh = txn.TransactionRefresh
		}
		if err := s.saveTransaction(userID, account, txn); err != nil {
			log.Errorf("Failed to save Stripe transaction %s: %v", txn.ID, err)
		}
	}

	return latestRefresh, nil
}

func (s *Service) saveTransaction(userID int64, account models.StripeFCAccount, txn *stripe.FinancialConnectionsTransaction) error {
	pending := txn.Status == stripe.FinancialConnectionsTransactionStatusPending
	amount := float64(txn.Amount) / 100

	trans := &models.Transaction{
		Provider:           providerName,
		PlaidID:            account.ID,
		PlaidTransactionID: txn.ID,
		Date:               txn.TransactedAt,
		Amount:             amount,
		Name:               txn.Description,
		MerchantName:       txn.Description,
		PlaidCategory:      "",
		Pending:            pending,
	}

	existing, err := s.store.GetTransactionByPlaidID(txn.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		if err := s.store.CreateTransaction(trans); err != nil {
			return err
		}
		return s.tagService.AutoTagTransaction(userID, trans)
	}

	trans.ID = existing.ID
	if err := s.store.UpdateTransaction(trans); err != nil {
		return err
	}
	return s.tagService.AutoTagTransaction(userID, trans)
}
