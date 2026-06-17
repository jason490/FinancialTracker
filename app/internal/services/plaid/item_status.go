package plaid

import (
	"context"
	"errors"

	"github.com/labstack/gommon/log"
	"github.com/plaid/plaid-go/v42/plaid"
)

// applyItemStatusFromPlaidError persists item status from a Plaid API error and marks accounts disconnected when terminal.
func (p *PlaidService) applyItemStatusFromPlaidError(itemID string, err error) {
	code, ok := parsePlaidError(err)
	if !ok {
		return
	}
	status := mapPlaidErrorToStatus(code, "")
	if err := p.store.UpdatePlaidItemStatus(itemID, status, code); err != nil {
		log.Errorf("Failed to update item status for %s: %v", itemID, err)
	}
	if isTerminalPlaidError(code) {
		if err := p.store.MarkPlaidAccountsDisconnectedByItemID(itemID); err != nil {
			log.Errorf("Failed to mark accounts disconnected for item %s: %v", itemID, err)
		}
	}
}

// applyItemStatusFromItemError persists status from a successful ItemGet response that includes item.error.
func (p *PlaidService) applyItemStatusFromItemError(item plaid.ItemWithConsentFields, itemID string) string {
	code, reason, hasError := itemErrorFromGet(item)
	if !hasError {
		if err := p.store.UpdatePlaidItemStatus(itemID, ItemStatusActive, ""); err != nil {
			log.Errorf("Failed to clear item error status for %s: %v", itemID, err)
		}
		return ItemStatusActive
	}
	status := mapPlaidErrorToStatus(code, reason)
	if err := p.store.UpdatePlaidItemStatus(itemID, status, code); err != nil {
		log.Errorf("Failed to update item status for %s: %v", itemID, err)
	}
	if isTerminalPlaidError(code) || status == ItemStatusDisconnected {
		if err := p.store.MarkPlaidAccountsDisconnectedByItemID(itemID); err != nil {
			log.Errorf("Failed to mark accounts disconnected for item %s: %v", itemID, err)
		}
	}
	return status
}

// fetchInstitutionDetails loads institution metadata from Plaid for an access token.
func (p *PlaidService) fetchInstitutionDetails(ctx *context.Context, accessToken string) (string, string, string, error) {
	itemRequest := plaid.NewItemGetRequest(accessToken)
	itemResp, _, err := p.client.PlaidApi.ItemGet(*ctx).ItemGetRequest(*itemRequest).Execute()
	if err != nil {
		log.Error(err)
		return "", "", "", err
	}

	institutionID := ""
	institutionName := ""
	if val, ok := itemResp.Item.GetInstitutionIdOk(); ok && val != nil {
		institutionID = *val
		instReq := plaid.NewInstitutionsGetByIdRequest(institutionID, []plaid.CountryCode{plaid.COUNTRYCODE_US})
		instResp, _, err := p.client.PlaidApi.InstitutionsGetById(*ctx).InstitutionsGetByIdRequest(*instReq).Execute()
		if err == nil {
			institutionName = instResp.Institution.Name
		}
	}

	status := p.applyItemStatusFromItemError(itemResp.Item, itemResp.Item.GetItemId())
	return institutionID, institutionName, status, nil
}

// RemovePlaidItemAtInstitution calls Plaid /item/remove to invalidate the access token.
func (p *PlaidService) RemovePlaidItemAtInstitution(ctx *context.Context, accessToken string) error {
	request := plaid.NewItemRemoveRequest(accessToken)
	_, _, err := p.client.PlaidApi.ItemRemove(*ctx).ItemRemoveRequest(*request).Execute()
	if err != nil {
		if code, ok := parsePlaidError(err); ok && code == "ITEM_NOT_FOUND" {
			return nil
		}
		return err
	}
	return nil
}

// DisconnectItem removes the Item at Plaid and deletes the local connection record.
func (p *PlaidService) DisconnectItem(ctx *context.Context, rowID string, userID int64) error {
	item, err := p.store.GetPlaidItemByRowID(rowID, userID)
	if err != nil {
		return errors.New("failed to find bank connection")
	}
	if err := p.RemovePlaidItemAtInstitution(ctx, item.AccessToken); err != nil {
		log.Errorf("Plaid item/remove failed for %s: %v", item.PlaidItemID, err)
	}
	return p.store.DeletePlaidItem(rowID, userID)
}
