package models

// StripeFCItem represents a Stripe Financial Connections institution grouping.
type StripeFCItem struct {
	ID                    int64  `json:"id"`
	RowID                 string `json:"row_id"`
	UserID                int64  `json:"user_id"`
	InstitutionName       string `json:"institution_name"`
	Status                string `json:"status"`
	ErrorCode             string `json:"error_code,omitempty"`
	LastSynced            int64  `json:"last_synced"`
	TransactionRefreshID  string `json:"transaction_refresh_id"`
	CreatedAt             int64  `json:"created_at"`
}

// StripeFCAccount represents a linked Stripe Financial Connections account.
type StripeFCAccount struct {
	ID               int64   `json:"id"`
	RowID            string  `json:"row_id"`
	UserID           int64   `json:"user_id"`
	StripeAccountID  string  `json:"stripe_account_id"`
	StripeItemRowID  string  `json:"stripe_item_row_id"`
	Name             string  `json:"name"`
	Mask             string  `json:"mask"`
	Type             string  `json:"type"`
	Subtype          string  `json:"subtype"`
	Balance          float64 `json:"balance"`
	AvailableBalance float64 `json:"available_balance"`
	Currency         string  `json:"currency"`
	Status           string  `json:"status"`
	IsHidden         bool    `json:"is_hidden"`
	CreatedAt        int64   `json:"created_at"`
}

// StripeFCItemWithAccounts groups accounts under a Stripe FC item.
type StripeFCItemWithAccounts struct {
	StripeFCItem
	Accounts []StripeFCAccount
}
