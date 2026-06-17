package models

// Account represents a bank account from Plaid
type Account struct {
	ID               int64   `json:"id"`
	UserID           int64   `json:"user_id"`
	PlaidAccountID   string  `json:"plaid_account_id"`
	PlaidItemID      string  `json:"plaid_item_id"`
	Name             string  `json:"name"`
	Mask             string  `json:"mask"`
	Type             string  `json:"type"`
	Subtype          string  `json:"subtype"`
	Balance          float64 `json:"balance"`
	AvailableBalance float64 `json:"available_balance"`
	Currency         string  `json:"currency"`
	Status           string  `json:"status"`
	IsHidden         bool    `json:"is_hidden"`
	MonthlyPayment   float64 `json:"monthly_payment"`
	CreatedAt        int64   `json:"created_at"`
}

// PlaidItem represents a connection to a financial institution via Plaid
type PlaidItem struct {
	ID              int64  `json:"id"`
	RowID           string `json:"row_id"`
	UserID          int64  `json:"user_id"`
	PlaidItemID     string `json:"plaid_item_id"`
	AccessToken     string `json:"-"`
	InstitutionID   string `json:"institution_id"`
	InstitutionName string `json:"institution_name"`
	SyncCursor      string `json:"sync_cursor"`
	Status          string `json:"status"`
	ErrorCode       string `json:"error_code,omitempty"`
	LastSynced      int64  `json:"last_synced"`
	CreatedAt       int64  `json:"created_at"`
}

// PlaidItemWithAccounts groups accounts under their respective Plaid item
type PlaidItemWithAccounts struct {
	PlaidItem
	Accounts []Account
}
