package external

import "FinancialTracker/internal/models"

// SessionProfile is the minimal authenticated user payload exposed to clients.
type SessionProfile struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	ThemePreference string `json:"theme_preference"`
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

// DashboardPayload is the lean dashboard response for the SPA.
type DashboardPayload struct {
	Summary       models.DashboardSummary          `json:"summary"`
	Groups        map[string][]DashboardAccountView `json:"groups"`
	Transactions  []DashboardTransactionView       `json:"transactions"`
	SpendingTrend []models.MonthlySpend            `json:"spending_trend"`
	MonthCashflow models.MonthCashflow             `json:"month_cashflow"`
	SpendingByTag []TagSliceView                   `json:"spending_by_tag"`
	IncomeByTag   []TagSliceView                   `json:"income_by_tag"`
	Layout        models.DashboardLayout           `json:"layout"`
	EditMode      bool                             `json:"edit_mode"`
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

// ToTagsPayload converts internal category/tag data to the external API payload.
func ToTagsPayload(data []models.CategoryWithTags) *TagsPayload {
	cats := make([]CategoryWithTagsView, 0, len(data))
	for i := range data {
		cwt := &data[i]
		tags := make([]TagView, 0, len(cwt.Tags))
		for j := range cwt.Tags {
			tags = append(tags, TagView{
				ID:         cwt.Tags[j].ID,
				CategoryID: cwt.Tags[j].CategoryID,
				Name:       cwt.Tags[j].Name,
				Color:      cwt.Tags[j].Color,
			})
		}
		cats = append(cats, CategoryWithTagsView{
			ID:   cwt.ID,
			Name: cwt.Name,
			Tags: tags,
		})
	}
	return &TagsPayload{Categories: cats}
}

// ToTagFilterViews converts internal tag filters to the external API payload.
func ToTagFilterViews(filters []models.TagFilter) []TagFilterView {
	out := make([]TagFilterView, 0, len(filters))
	for i := range filters {
		out = append(out, TagFilterView{
			Pattern:    filters[i].Pattern,
			FilterType: filters[i].FilterType,
		})
	}
	return out
}

// ToTransactionListPayload converts internal TransactionPageData to the external API payload.
func ToTransactionListPayload(data *models.TransactionPageData) *TransactionListPayload {
	txViews := make([]TransactionView, 0, len(data.Transactions))
	for i := range data.Transactions {
		t := &data.Transactions[i]
		tagViews := make([]TransactionTagView, 0, len(t.Tags))
		for j := range t.Tags {
			tagViews = append(tagViews, TransactionTagView{
				ID:    t.Tags[j].ID,
				Name:  t.Tags[j].Name,
				Color: t.Tags[j].Color,
			})
		}
		txViews = append(txViews, TransactionView{
			ID:           t.ID,
			Date:         t.Date,
			Amount:       t.Amount,
			Name:         t.Name,
			MerchantName: t.MerchantName,
			Pending:      t.Pending,
			Tags:         tagViews,
		})
	}

	tagOpts := make([]TransactionTagOption, 0, len(data.AllTags))
	for i := range data.AllTags {
		tagOpts = append(tagOpts, TransactionTagOption{
			ID:         data.AllTags[i].ID,
			Name:       data.AllTags[i].Name,
			Color:      data.AllTags[i].Color,
			CategoryID: data.AllTags[i].CategoryID,
		})
	}

	catViews := make([]TransactionCategoryView, 0, len(data.Categories))
	for i := range data.Categories {
		catViews = append(catViews, TransactionCategoryView{
			ID:   data.Categories[i].ID,
			Name: data.Categories[i].Name,
		})
	}

	totalPages := 0
	if data.PageSize > 0 {
		totalPages = (data.TotalCount + data.PageSize - 1) / data.PageSize
	}

	return &TransactionListPayload{
		Transactions: txViews,
		TotalCount:   data.TotalCount,
		Page:         data.CurrentPage,
		PageSize:     data.PageSize,
		TotalPages:   totalPages,
		Tags:         tagOpts,
		Categories:   catViews,
	}
}
