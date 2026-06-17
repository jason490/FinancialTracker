package plaid

import (
	"FinancialTracker/internal/models"
	"context"
	"errors"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/plaid/plaid-go/v42/plaid"
)

// SyncPlaidAccounts fetches the latest account information from Plaid and updates the local database.
func (p *PlaidService) SyncPlaidAccounts(ctx *context.Context, userID int64, itemID string, accessToken string) error {
	accountsRequest := plaid.NewAccountsGetRequest(accessToken)
	accountsResp, _, err := p.client.PlaidApi.AccountsGet(*ctx).AccountsGetRequest(*accountsRequest).Execute()
	if err != nil {
		log.Errorf("Failed to fetch accounts from Plaid for user %d: %v", userID, err)
		p.applyItemStatusFromPlaidError(itemID, err)
		if code, ok := parsePlaidError(err); ok {
			if isTerminalPlaidError(code) {
				return errors.New("bank connection is no longer available")
			}
			if code == "ITEM_LOGIN_REQUIRED" {
				return errors.New("bank connection requires re-authentication")
			}
		}
		return errors.New("unable to fetch bank accounts")
	}

	if err := p.store.UpdatePlaidItemStatus(itemID, ItemStatusActive, ""); err != nil {
		log.Errorf("Failed to clear item status for %s: %v", itemID, err)
	}

	activeAccountIDs := make(map[string]bool)
	for _, pa := range accountsResp.GetAccounts() {
		activeAccountIDs[pa.AccountId] = true
		var mask, subtype, currency string
		var balance, available float64

		if val, ok := pa.GetMaskOk(); ok {
			mask = *val
		}
		if val, ok := pa.GetSubtypeOk(); ok {
			subtype = string(*val)
		}

		balances := pa.GetBalances()
		if val, ok := balances.GetIsoCurrencyCodeOk(); ok {
			currency = *val
		}
		if val, ok := balances.GetCurrentOk(); ok {
			balance = *val
		}
		if val, ok := balances.GetAvailableOk(); ok {
			available = *val
		}

		acc := &models.Account{
			UserID:           userID,
			PlaidAccountID:   pa.AccountId,
			PlaidItemID:      itemID,
			Name:             pa.Name,
			Mask:             mask,
			Type:             string(pa.Type),
			Subtype:          subtype,
			Balance:          balance,
			AvailableBalance: available,
			Currency:         currency,
			Status:           "active",
		}

		existing, err := p.store.GetAccountByPlaidAccountID(pa.AccountId)
		if err == nil && existing != nil {
			if err := p.store.UpdatePlaidAccount(acc); err != nil {
				log.Errorf("Failed to update account %s: %v", pa.AccountId, err)
			}
		} else {
			if err := p.store.CreatePlaidAccount(acc); err != nil {
				log.Errorf("Failed to create account %s: %v", pa.AccountId, err)
			}
		}
	}

	dbAccounts, err := p.store.GetPlaidAccountsByItemID(itemID)
	if err == nil {
		for _, dbAcc := range dbAccounts {
			if !activeAccountIDs[dbAcc.PlaidAccountID] && dbAcc.Status != "disconnected" {
				log.Infof("Marking account %s (%s) as disconnected for item %s", dbAcc.Name, dbAcc.PlaidAccountID, itemID)
				if err := p.store.UpdatePlaidAccountStatus(dbAcc.PlaidAccountID, "disconnected"); err != nil {
					log.Errorf("Failed to mark account %s as disconnected: %v", dbAcc.PlaidAccountID, err)
				}
			}
		}
	}

	// Disabling this feature for now
	// p.SyncLiabilities(ctx, accessToken)
	return nil
}

// SyncAccounts iterates through all connected items for a user and synchronizes their account metadata and balances.
func (p *PlaidService) SyncAllAccounts(c *echo.Context, userID int64) error {
	items, err := p.store.GetPlaidItemsByUserID(userID)
	if err != nil {
		log.Errorf("Failed to fetch Plaid items for user %d: %v", userID, err)
		return errors.New("failed to fetch linked accounts")
	}

	ctx := c.Request().Context()
	for _, item := range items {
		if item.Status == ItemStatusDisconnected {
			continue
		}
		if err := p.SyncPlaidAccounts(&ctx, userID, item.PlaidItemID, item.AccessToken); err != nil {
			log.Errorf("Failed to sync accounts for item %s (User %d): %v", item.PlaidItemID, userID, err)
		}
	}
	return nil
}
