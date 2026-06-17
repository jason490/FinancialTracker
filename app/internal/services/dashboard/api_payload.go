package dashboard

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/models/external"
)

// BuildAPIPayload converts internal dashboard data into a lean API response.
func BuildAPIPayload(data *models.DashboardData) *external.DashboardPayload {
	if data == nil {
		return nil
	}

	groups := make(map[string][]external.DashboardAccountView, len(data.Groups))
	for bucket, accounts := range data.Groups {
		views := make([]external.DashboardAccountView, 0, len(accounts))
		for _, acc := range accounts {
			views = append(views, external.DashboardAccountView{
				Name:           acc.Name,
				Mask:           acc.Mask,
				Subtype:        acc.Subtype,
				Balance:        DisplayBalance(acc, ClassifyAccount(acc)),
				MonthlyPayment: acc.MonthlyPayment,
				Status:         acc.Status,
				IsHidden:       acc.IsHidden,
			})
		}
		groups[bucket] = views
	}

	transactions := make([]external.DashboardTransactionView, 0, len(data.Transactions))
	for _, tx := range data.Transactions {
		transactions = append(transactions, mapTransactionView(tx))
	}

	return &external.DashboardPayload{
		Summary:       data.Summary,
		Groups:        ensureGroups(groups),
		Transactions:  ensureTransactions(transactions),
		SpendingTrend: ensureMonthlySpend(data.SpendingTrend),
		MonthCashflow: data.MonthCashflow,
		SpendingByTag: mapTagSlices(data.SpendingByTag),
		IncomeByTag:   mapTagSlices(data.IncomeByTag),
		Layout:        ensureLayout(data.Layout),
		EditMode:      data.EditMode,
	}
}

// GetDashboardPayload loads dashboard data and returns the lean API payload.
func (s *DashboardService) GetDashboardPayload(userID int64, editMode bool) (*external.DashboardPayload, error) {
	data, err := s.GetDashboardData(userID, editMode)
	if err != nil {
		return nil, err
	}
	return BuildAPIPayload(data), nil
}

// SaveDeviceLayout validates and persists a dashboard layout for a specific device, returning the updated payload.
func (s *DashboardService) SaveDeviceLayout(userID int64, deviceType string, widgets []models.DashboardWidget) (*external.DashboardPayload, error) {
	existing, err := s.store.GetDashboardLayout(userID)
	if err != nil {
		return nil, err
	}

	layout := DefaultLayout()
	if existing != nil {
		layout = *existing
	}

	if deviceType == "mobile" {
		layout.Mobile = widgets
	} else {
		layout.Desktop = widgets
	}

	layout = NormalizeLayout(layout)
	if err := s.store.UpsertDashboardLayout(userID, &layout); err != nil {
		return nil, err
	}
	return s.GetDashboardPayload(userID, false)
}

func mapTransactionView(tx models.Transaction) external.DashboardTransactionView {
	view := external.DashboardTransactionView{
		Date:         tx.Date,
		Amount:       tx.Amount,
		Name:         tx.Name,
		MerchantName: tx.MerchantName,
		Pending:      tx.Pending,
	}
	if len(tx.Tags) > 0 {
		view.Tags = make([]external.DashboardTagView, len(tx.Tags))
		for i, tag := range tx.Tags {
			view.Tags[i] = external.DashboardTagView{
				Name:  tag.Name,
				Color: tag.Color,
			}
		}
	}
	return view
}

func mapTagSlices(rows []models.TagBreakdown) []external.TagSliceView {
	if len(rows) == 0 {
		return []external.TagSliceView{}
	}
	views := make([]external.TagSliceView, len(rows))
	for i, row := range rows {
		views[i] = external.TagSliceView{
			TagName: row.TagName,
			Color:   row.Color,
			Total:   row.Total,
		}
	}
	return views
}

func ensureMonthlySpend(rows []models.MonthlySpend) []models.MonthlySpend {
	if rows == nil {
		return []models.MonthlySpend{}
	}
	return rows
}

func ensureTransactions(rows []external.DashboardTransactionView) []external.DashboardTransactionView {
	if rows == nil {
		return []external.DashboardTransactionView{}
	}
	return rows
}

func ensureGroups(groups map[string][]external.DashboardAccountView) map[string][]external.DashboardAccountView {
	if groups == nil {
		groups = make(map[string][]external.DashboardAccountView, len(AllBuckets))
	}
	for _, bucket := range AllBuckets {
		if groups[bucket] == nil {
			groups[bucket] = []external.DashboardAccountView{}
		}
	}
	return groups
}

func ensureLayout(layout models.DashboardLayout) models.DashboardLayout {
	if layout.Desktop == nil {
		layout.Desktop = []models.DashboardWidget{}
	}
	if layout.Mobile == nil {
		layout.Mobile = []models.DashboardWidget{}
	}
	return layout
}
