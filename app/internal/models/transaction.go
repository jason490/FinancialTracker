package models

// Transaction represents a financial transaction
type Transaction struct {
	ID                 int64   `json:"id"`
	PlaidID            int64   `json:"plaid_id"`
	PlaidTransactionID string  `json:"plaid_transaction_id"`
	Date               int64   `json:"date"`
	Amount             float64 `json:"amount"`
	Name               string  `json:"name"`
	MerchantName       string  `json:"merchant_name"`
	PlaidCategory      string  `json:"plaid_category"`
	Pending            bool    `json:"pending"`
	CreatedAt          int64   `json:"created_at"`
	Tags               []Tag   `json:"tags,omitempty"`
}

// TransactionFilters represents the filters for querying transactions
type TransactionFilters struct {
	Search     string
	MinAmount  *float64
	MaxAmount  *float64
	StartDate  *int64
	EndDate    *int64
	CategoryID *int64
	Tags       []int64
	SortBy     string // "date", "amount", "name"
	SortDir    string // "asc", "desc"
	Limit      int
	Offset     int
}

// TransactionPageData contains data specifically for the transactions page
type TransactionPageData struct {
	Transactions []Transaction
	TotalCount   int
	CurrentPage  int
	PageSize     int
	Filters      TransactionFilters
	AllTags      []Tag
	Categories   []Category
}
