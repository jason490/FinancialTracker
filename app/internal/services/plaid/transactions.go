package plaid

import (
	"context"
	"FinancialTracker/internal/models"
	"errors"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/plaid/plaid-go/v42/plaid"
)

// SyncTransactions synchronizes transaction data from Plaid for all items associated with a user.
func (p *PlaidService) SyncTransactions(c *echo.Context, userID int64) error {
	items, err := p.store.GetPlaidItemsByUserID(userID)
	if err != nil {
		log.Errorf("Failed to fetch Plaid items for user %d: %v", userID, err)
		return errors.New("failed to fetch linked accounts")
	}

	if len(items) == 0 {
		return errors.New("you have no linked accounts")
	}

	ctx := c.Request().Context()
	for _, item := range items {
		if item.Status == ItemStatusDisconnected {
			continue
		}
		if err := p.syncItemTransactions(ctx, userID, item.PlaidItemID, item.AccessToken, item.SyncCursor); err != nil {
			log.Errorf("Failed to sync transactions for item %s: %v", item.PlaidItemID, err)
			continue
		}
	}
	return nil
}

// parsePlaidTransaction converts a Plaid API transaction into our internal models.Transaction
func (p *PlaidService) parsePlaidTransaction(t plaid.Transaction, userID int64) (*models.Transaction, error) {
	parsedDate, _ := time.Parse("2006-01-02", t.Date)
	unixDate := parsedDate.Unix()

	var merchantName, plaidCategory string
	if val, ok := t.GetMerchantNameOk(); ok {
		merchantName = *val
	}

	// Use Personal Finance Category (PFCv2) as prioritized in the pfc-taxonomy-all.csv
	if cats, ok := t.GetPersonalFinanceCategoryOk(); ok {
		plaidCategory = cats.GetPrimary()
	} else if cats := t.GetCategory(); len(cats) > 0 {
		plaidCategory = cats[0]
	}

	// Find internal account ID
	acc, err := p.store.GetAccountByPlaidAccountID(t.AccountId)
	if err != nil {
		return nil, err
	}

	return &models.Transaction{
		Provider:           "plaid",
		PlaidID:            acc.ID,
		PlaidTransactionID: t.TransactionId,
		Date:               unixDate,
		Amount:             t.Amount,
		Name:               t.Name,
		MerchantName:       merchantName,
		PlaidCategory:      plaidCategory,
		Pending:            t.Pending,
	}, nil
}

// SyncItemTransactions synchronizes transaction data for a specific Plaid item using the provided parameters.
func (p *PlaidService) SyncItemTransactions(c *echo.Context, userID int64, plaidItemID, accessToken, cursor string) error {
	return p.syncItemTransactions(c.Request().Context(), userID, plaidItemID, accessToken, cursor)
}

// syncItemTransactions synchronizes transaction data for a specific Plaid item.
func (p *PlaidService) syncItemTransactions(ctx context.Context, userID int64, plaidItemID, accessToken, cursor string) error {
	hasMore := true

	for hasMore {
		request := plaid.NewTransactionsSyncRequest(accessToken)
		if cursor != "" {
			request.SetCursor(cursor)
		}
		resp, _, err := p.client.PlaidApi.TransactionsSync(ctx).TransactionsSyncRequest(*request).Execute()
		if err != nil {
			p.applyItemStatusFromPlaidError(plaidItemID, err)
			if code, ok := parsePlaidError(err); ok {
				if isTerminalPlaidError(code) {
					return errors.New("bank connection is no longer available")
				}
				if code == "ITEM_LOGIN_REQUIRED" {
					return errors.New("bank connection requires re-authentication")
				}
			}
			return err
		}

		// Process added transactions
		for _, t := range resp.GetAdded() {
			trans, err := p.parsePlaidTransaction(t, userID)
			if err != nil {
				log.Errorf("Transaction %s belongs to unknown account %s", t.TransactionId, t.AccountId)
				continue
			}

			if err := p.store.CreateTransaction(trans); err == nil {
				// Automatically tag the new transaction
				if err := p.tagService.AutoTagTransaction(userID, trans); err != nil {
					log.Errorf("Failed to auto-tag transaction %s: %v", trans.PlaidTransactionID, err)
				}
			} else {
				log.Errorf("Failed to save transaction %s: %v", t.TransactionId, err)
			}
		}

		// Process modified transactions
		for _, t := range resp.GetModified() {
			trans, err := p.parsePlaidTransaction(t, userID)
			if err != nil {
				log.Errorf("Transaction %s belongs to unknown account %s", t.TransactionId, t.AccountId)
				continue
			}

			if err := p.store.UpdateTransaction(trans); err != nil {
				log.Errorf("Failed to update transaction %s: %v", t.TransactionId, err)
			} else {
				// Re-evaluate tags for modified transaction
				if err := p.tagService.AutoTagTransaction(userID, trans); err != nil {
					log.Errorf("Failed to re-tag modified transaction %s: %v", trans.PlaidTransactionID, err)
				}
			}
		}

		// Process removed transactions
		for _, t := range resp.GetRemoved() {
			if err := p.store.DeleteTransactionByPlaidID(t.GetTransactionId()); err != nil {
				log.Errorf("Failed to remove transaction %s: %v", t.GetTransactionId(), err)
			}
		}

		hasMore = resp.GetHasMore()
		cursor = resp.GetNextCursor()

		// Update cursor in DB to persist progress
		if err := p.store.UpdatePlaidItemCursor(plaidItemID, cursor); err != nil {
			log.Errorf("Failed to update sync cursor for item %s: %v", plaidItemID, err)
		}
	}
	if err := p.store.UpdatePlaidItemStatus(plaidItemID, ItemStatusActive, ""); err != nil {
		log.Errorf("Failed to clear item status after sync for %s: %v", plaidItemID, err)
	}
	return nil
}
