package transactions

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/models/external"
)

// BuildListPayload converts transaction page data into the external API payload.
func BuildListPayload(data *models.TransactionPageData) *external.TransactionListPayload {
	txViews := make([]external.TransactionView, 0, len(data.Transactions))
	for i := range data.Transactions {
		t := &data.Transactions[i]
		tagViews := make([]external.TransactionTagView, 0, len(t.Tags))
		for j := range t.Tags {
			tagViews = append(tagViews, external.TransactionTagView{
				ID:    t.Tags[j].ID,
				Name:  t.Tags[j].Name,
				Color: t.Tags[j].Color,
			})
		}
		txViews = append(txViews, external.TransactionView{
			ID:           t.ID,
			Date:         t.Date,
			Amount:       t.Amount,
			Name:         t.Name,
			MerchantName: t.MerchantName,
			Pending:      t.Pending,
			Tags:         tagViews,
		})
	}

	tagOpts := make([]external.TransactionTagOption, 0, len(data.AllTags))
	for i := range data.AllTags {
		tagOpts = append(tagOpts, external.TransactionTagOption{
			ID:         data.AllTags[i].ID,
			Name:       data.AllTags[i].Name,
			Color:      data.AllTags[i].Color,
			CategoryID: data.AllTags[i].CategoryID,
		})
	}

	catViews := make([]external.TransactionCategoryView, 0, len(data.Categories))
	for i := range data.Categories {
		catViews = append(catViews, external.TransactionCategoryView{
			ID:   data.Categories[i].ID,
			Name: data.Categories[i].Name,
		})
	}

	totalPages := 0
	if data.PageSize > 0 {
		totalPages = (data.TotalCount + data.PageSize - 1) / data.PageSize
	}

	return &external.TransactionListPayload{
		Transactions: txViews,
		TotalCount:   data.TotalCount,
		Page:         data.CurrentPage,
		PageSize:     data.PageSize,
		TotalPages:   totalPages,
		Tags:         tagOpts,
		Categories:   catViews,
	}
}
