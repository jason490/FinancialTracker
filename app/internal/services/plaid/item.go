package plaid

import (
	"FinancialTracker/internal/models"
	"context"
	"errors"

	"github.com/labstack/gommon/log"
)

// CreatePlaidItem fetches institution info from Plaid and saves a new item connection.
func (p *PlaidService) CreatePlaidItem(ctx *context.Context, userID int64, itemID string, accessToken string) error {
	institutionID, institutionName, status, err := p.fetchInstitutionDetails(ctx, userID, accessToken)
	if err != nil {
		p.applyItemStatusFromPlaidError(itemID, err)
		if code, ok := parsePlaidError(err); ok && isTerminalPlaidError(code) {
			return errors.New("this bank connection is no longer available")
		}
		return errors.New("failed to fetch account connection details")
	}

	plaidItem := &models.PlaidItem{
		UserID:          userID,
		PlaidItemID:     itemID,
		AccessToken:     accessToken,
		InstitutionID:   institutionID,
		InstitutionName: institutionName,
		Status:          status,
	}
	if err := p.store.CreatePlaidItem(plaidItem); err != nil {
		log.Error(err)
		return errors.New("failed to save bank connection")
	}
	return nil
}

// UpdatePlaidItem refreshes an existing Plaid Item's metadata, access token, and connection status.
func (p *PlaidService) UpdatePlaidItem(ctx *context.Context, userID int64, itemID string, accessToken string) error {
	institutionID, institutionName, status, err := p.fetchInstitutionDetails(ctx, userID, accessToken)
	if err != nil {
		p.applyItemStatusFromPlaidError(itemID, err)
		if code, ok := parsePlaidError(err); ok && isTerminalPlaidError(code) {
			return errors.New("this bank connection is no longer available")
		}
		return errors.New("failed to fetch updated account connection details")
	}

	plaidItem := &models.PlaidItem{
		UserID:          userID,
		PlaidItemID:     itemID,
		AccessToken:     accessToken,
		InstitutionID:   institutionID,
		InstitutionName: institutionName,
		Status:          status,
	}
	if err := p.store.UpdatePlaidItem(plaidItem); err != nil {
		log.Error(err)
		return errors.New("failed to update bank connection")
	}
	return nil
}

// GetItemByRowID loads a Plaid item by its public row identifier.
func (p *PlaidService) GetItemByRowID(rowID string, userID int64) (*models.PlaidItem, error) {
	return p.store.GetPlaidItemByRowID(rowID, userID)
}

// UpdateItemLastSynced records the last successful sync timestamp for an item.
func (p *PlaidService) UpdateItemLastSynced(itemID string, syncedAt int64) error {
	return p.store.UpdatePlaidItemLastSynced(itemID, syncedAt)
}
