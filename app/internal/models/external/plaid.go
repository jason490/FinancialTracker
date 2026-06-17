package external

import "FinancialTracker/internal/models"

// PlaidAccountView is a sanitized bank account for connection management.
type PlaidAccountView struct {
	PlaidAccountID string  `json:"plaid_account_id"`
	Name           string  `json:"name"`
	Mask           string  `json:"mask"`
	Subtype        string  `json:"subtype"`
	Balance        float64 `json:"balance"`
	Currency       string  `json:"currency"`
	Status         string  `json:"status"`
	IsHidden       bool    `json:"is_hidden"`
}

// PlaidConnectionView is a Plaid institution connection with its accounts.
type PlaidConnectionView struct {
	RowID           string             `json:"row_id"`
	InstitutionName string             `json:"institution_name"`
	Status          string             `json:"status"`
	CreatedAt       int64              `json:"created_at"`
	LastSynced      int64              `json:"last_synced"`
	Accounts        []PlaidAccountView `json:"accounts"`
}

// PlaidConnectionsPayload lists all Plaid connections for the user.
type PlaidConnectionsPayload struct {
	Connections []PlaidConnectionView `json:"connections"`
}

// PlaidLinkTokenResponse returns a Plaid Link token to the client.
type PlaidLinkTokenResponse struct {
	LinkToken string `json:"link_token"`
}

// PlaidExchangeRequest carries the public token from Plaid Link.
type PlaidExchangeRequest struct {
	PublicToken string `json:"public_token"`
}

// ToPlaidConnectionsPayload converts internal Plaid data to external views.
func ToPlaidConnectionsPayload(items []models.PlaidItemWithAccounts) *PlaidConnectionsPayload {
	connections := make([]PlaidConnectionView, 0, len(items))
	for i := range items {
		item := &items[i]
		accounts := make([]PlaidAccountView, 0, len(item.Accounts))
		for j := range item.Accounts {
			acc := &item.Accounts[j]
			accounts = append(accounts, PlaidAccountView{
				PlaidAccountID: acc.PlaidAccountID,
				Name:           acc.Name,
				Mask:           acc.Mask,
				Subtype:        acc.Subtype,
				Balance:        acc.Balance,
				Currency:       acc.Currency,
				Status:         acc.Status,
				IsHidden:       acc.IsHidden,
			})
		}
		connections = append(connections, PlaidConnectionView{
			RowID:           item.RowID,
			InstitutionName: item.InstitutionName,
			Status:          item.Status,
			CreatedAt:       item.CreatedAt,
			LastSynced:      item.LastSynced,
			Accounts:        accounts,
		})
	}
	return &PlaidConnectionsPayload{Connections: connections}
}
