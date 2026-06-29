package connections

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/models/external"
)

const (
	ProviderPlaid  = "plaid"
	ProviderStripe = "stripe"
)

// BuildPayloadFromPlaid maps Plaid management data into the unified connections payload.
func BuildPayloadFromPlaid(items []models.PlaidItemWithAccounts, usage *models.PlaidUsage) *external.ConnectionsPayload {
	connections := make([]external.ConnectionView, 0, len(items))
	for i := range items {
		item := &items[i]
		accounts := make([]external.ConnectionAccountView, 0, len(item.Accounts))
		for j := range item.Accounts {
			acc := &item.Accounts[j]
			accounts = append(accounts, external.ConnectionAccountView{
				AccountID: acc.RowID,
				Name:      acc.Name,
				Mask:      acc.Mask,
				Subtype:   acc.Subtype,
				Balance:   acc.Balance,
				Currency:  acc.Currency,
				Status:    acc.Status,
				IsHidden:  acc.IsHidden,
			})
		}
		connections = append(connections, connectionView(
			item.RowID,
			item.InstitutionName,
			item.Status,
			item.CreatedAt,
			item.LastSynced,
			accounts,
		))
	}
	return buildPayload(ProviderPlaid, connections, usage)
}

// BuildPayloadFromStripe maps Stripe FC management data into the unified connections payload.
func BuildPayloadFromStripe(items []models.StripeFCItemWithAccounts, usage *models.PlaidUsage) *external.ConnectionsPayload {
	connections := make([]external.ConnectionView, 0, len(items))
	for i := range items {
		item := &items[i]
		accounts := make([]external.ConnectionAccountView, 0, len(item.Accounts))
		for j := range item.Accounts {
			acc := &item.Accounts[j]
			accounts = append(accounts, external.ConnectionAccountView{
				AccountID: acc.RowID,
				Name:      acc.Name,
				Mask:      acc.Mask,
				Subtype:   acc.Subtype,
				Balance:   acc.Balance,
				Currency:  acc.Currency,
				Status:    acc.Status,
				IsHidden:  acc.IsHidden,
			})
		}
		connections = append(connections, connectionView(
			item.RowID,
			item.InstitutionName,
			item.Status,
			item.CreatedAt,
			item.LastSynced,
			accounts,
		))
	}
	return buildPayload(ProviderStripe, connections, usage)
}

func connectionView(
	rowID, institutionName, status string,
	createdAt, lastSynced int64,
	accounts []external.ConnectionAccountView,
) external.ConnectionView {
	return external.ConnectionView{
		RowID:           rowID,
		InstitutionName: institutionName,
		Status:          status,
		CreatedAt:       createdAt,
		LastSynced:      lastSynced,
		Accounts:        accounts,
	}
}

func buildPayload(provider string, connections []external.ConnectionView, usage *models.PlaidUsage) *external.ConnectionsPayload {
	payload := &external.ConnectionsPayload{
		Provider:    provider,
		Connections: connections,
	}
	if usage != nil {
		payload.Usage = *usage
	}
	return payload
}
