package external

import "FinancialTracker/internal/models"

// ForgotPasswordResponse acknowledges a password reset code request.
type ForgotPasswordResponse struct {
	Message              string `json:"message"`
	CodeExpiresInSeconds int64  `json:"code_expires_in_seconds"`
}

// VerifyResetCodeResponse confirms a reset code is valid.
type VerifyResetCodeResponse struct {
	ExpiresAt int64 `json:"expires_at"`
}

// SessionProfile is the minimal authenticated user payload exposed to clients.
type SessionProfile struct {
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Email                string `json:"email"`
	ThemePreference      string `json:"theme_preference"`
	OnboardingCompleted  bool   `json:"onboarding_completed"`
}

// DashboardAccountView is a sanitized account row for dashboard widgets.
type DashboardAccountView struct {
	Name           string  `json:"name"`
	Mask           string  `json:"mask"`
	Subtype        string  `json:"subtype"`
	Balance        float64 `json:"balance"`
	MonthlyPayment float64 `json:"monthly_payment,omitempty"`
	Status         string  `json:"status"`
	IsHidden       bool    `json:"is_hidden"`
}

// DashboardTagView is a minimal tag label for transaction rows.
type DashboardTagView struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// DashboardTransactionView is a sanitized transaction row for dashboard widgets.
type DashboardTransactionView struct {
	Date         int64              `json:"date"`
	Amount       float64            `json:"amount"`
	Name         string             `json:"name"`
	MerchantName string             `json:"merchant_name"`
	Pending      bool               `json:"pending"`
	Tags         []DashboardTagView `json:"tags,omitempty"`
}

// TagSliceView is an aggregated amount for a tag without internal identifiers.
type TagSliceView struct {
	TagName string  `json:"tag_name"`
	Color   string  `json:"color"`
	Total   float64 `json:"total"`
}

// DashboardAccountGroups maps dashboard bucket keys to visible accounts.
type DashboardAccountGroups map[string][]DashboardAccountView

// DashboardPayload is the lean dashboard response for the SPA.
type DashboardPayload struct {
	Summary       models.DashboardSummary `json:"summary"`
	Groups        DashboardAccountGroups  `json:"groups"`
	Transactions  []DashboardTransactionView        `json:"transactions"`
	SpendingTrend []models.MonthlySpend             `json:"spending_trend"`
	MonthCashflow models.MonthCashflow              `json:"month_cashflow"`
	SpendingByTag []TagSliceView                    `json:"spending_by_tag"`
	IncomeByTag   []TagSliceView                    `json:"income_by_tag"`
	Layout        models.DashboardLayout            `json:"layout"`
	EditMode      bool                              `json:"edit_mode"`
}

// TransactionTagView is a tag attached to a transaction row, with id exposed for bulk operations.
type TransactionTagView struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// TransactionView is a sanitized transaction for the frontend, stripping Plaid internals.
type TransactionView struct {
	ID           int64                `json:"id"`
	Date         int64                `json:"date"`
	Amount       float64              `json:"amount"`
	Name         string               `json:"name"`
	MerchantName string               `json:"merchant_name"`
	Pending      bool                 `json:"pending"`
	Tags         []TransactionTagView `json:"tags,omitempty"`
}

// TransactionCategoryView is a sanitized category for filter dropdowns.
type TransactionCategoryView struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// TransactionTagOption is a tag option for filter dropdowns, with category_id for grouping.
type TransactionTagOption struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	CategoryID int64  `json:"category_id"`
}

// TransactionListPayload is the full paginated transaction list response.
type TransactionListPayload struct {
	Transactions []TransactionView         `json:"transactions"`
	TotalCount   int                       `json:"total_count"`
	Page         int                       `json:"page"`
	PageSize     int                       `json:"page_size"`
	TotalPages   int                       `json:"total_pages"`
	Tags         []TransactionTagOption    `json:"tags"`
	Categories   []TransactionCategoryView `json:"categories"`
}

// TagView is a tag row for the tags management page.
type TagView struct {
	ID         int64  `json:"id"`
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
}

// CategoryWithTagsView is a category and its tags for the tags management page.
type CategoryWithTagsView struct {
	ID   int64     `json:"id"`
	Name string    `json:"name"`
	Tags []TagView `json:"tags"`
}

// TagsPayload is the full tags and categories response for the SPA.
type TagsPayload struct {
	Categories []CategoryWithTagsView `json:"categories"`
}

// TagFilterView is a filter rule exposed for tag editing.
type TagFilterView struct {
	Pattern    string `json:"pattern"`
	FilterType string `json:"filter_type"`
}
