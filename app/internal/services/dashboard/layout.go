package dashboard

import (
	"FinancialTracker/internal/models"
	"encoding/json"
	"sort"
)

// Widget IDs for dashboard panels.
const (
	WidgetNetWorth           = "net_worth"
	WidgetCashAccounts       = "cash_accounts"
	WidgetSavingsAccounts    = "savings_accounts"
	WidgetCreditAccounts     = "credit_accounts"
	WidgetLoanAccounts       = "loan_accounts"
	WidgetInvestmentAccounts = "investment_accounts"
	WidgetRecentTransactions = "recent_transactions"
	WidgetQuickActions       = "quick_actions"
	WidgetSpendingTrend      = "spending_trend"
	WidgetMonthCashflow      = "month_cashflow"
	WidgetSpendingByTag      = "spending_by_tag"
	WidgetIncomeByTag        = "income_by_tag"
)

var validWidgetIDs = map[string]bool{
	WidgetNetWorth:           true,
	WidgetCashAccounts:       true,
	WidgetSavingsAccounts:    true,
	WidgetCreditAccounts:     true,
	WidgetLoanAccounts:       true,
	WidgetInvestmentAccounts: true,
	WidgetRecentTransactions: true,
	WidgetQuickActions:       true,
	WidgetSpendingTrend:      true,
	WidgetMonthCashflow:      true,
	WidgetSpendingByTag:      true,
	WidgetIncomeByTag:        true,
}

func defaultWidgets() []models.DashboardWidget {
	return []models.DashboardWidget{
		{ID: WidgetNetWorth, Visible: true, Order: 0},
		{ID: WidgetSpendingByTag, Visible: true, Order: 1},
		{ID: WidgetMonthCashflow, Visible: true, Order: 2},
		{ID: WidgetIncomeByTag, Visible: true, Order: 3},
		{ID: WidgetCashAccounts, Visible: true, Order: 4},
		{ID: WidgetSavingsAccounts, Visible: true, Order: 5},
		{ID: WidgetCreditAccounts, Visible: true, Order: 6},
		{ID: WidgetLoanAccounts, Visible: true, Order: 7},
		{ID: WidgetInvestmentAccounts, Visible: true, Order: 8},
		{ID: WidgetQuickActions, Visible: true, Order: 9},
		{ID: WidgetSpendingTrend, Visible: true, Order: 10},
		{ID: WidgetRecentTransactions, Visible: true, Order: 11},
	}
}

// DefaultLayout is the layout used for new users.
func DefaultLayout() models.DashboardLayout {
	return models.DashboardLayout{
		Desktop: defaultWidgets(),
		Mobile:  defaultWidgets(),
	}
}

func normalizeWidgets(widgets []models.DashboardWidget) []models.DashboardWidget {
	seen := make(map[string]bool)
	var filtered []models.DashboardWidget
	maxOrder := -1

	for _, w := range widgets {
		if !validWidgetIDs[w.ID] {
			continue
		}
		if seen[w.ID] {
			continue
		}
		seen[w.ID] = true
		filtered = append(filtered, w)
		if w.Order > maxOrder {
			maxOrder = w.Order
		}
	}

	defaults := defaultWidgets()
	for _, dw := range defaults {
		if !seen[dw.ID] {
			maxOrder++
			dw.Order = maxOrder
			filtered = append(filtered, dw)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Order < filtered[j].Order
	})

	return filtered
}

// NormalizeLayout merges a saved layout with defaults so new widgets appear and unknown IDs are dropped.
func NormalizeLayout(layout models.DashboardLayout) models.DashboardLayout {
	return models.DashboardLayout{
		Desktop: normalizeWidgets(layout.Desktop),
		Mobile:  normalizeWidgets(layout.Mobile),
	}
}

// ParseLayoutJSON decodes and validates layout JSON from the client.
func ParseLayoutJSON(raw string) (models.DashboardLayout, error) {
	var layout models.DashboardLayout
	if err := json.Unmarshal([]byte(raw), &layout); err != nil {
		return models.DashboardLayout{}, err
	}
	return NormalizeLayout(layout), nil
}

// WidgetsForRender returns widgets to display: all widgets in edit mode, otherwise only visible ones.
// Used primarily by the legacy Templ frontend, defaults to Desktop layout.
func WidgetsForRender(data *models.DashboardData) []models.DashboardWidget {
	if data == nil {
		return nil
	}
	layout := NormalizeLayout(data.Layout)
	if data.EditMode {
		widgets := make([]models.DashboardWidget, len(layout.Desktop))
		copy(widgets, layout.Desktop)
		sort.Slice(widgets, func(i, j int) bool {
			return widgets[i].Order < widgets[j].Order
		})
		return widgets
	}
	return VisibleWidgets(layout.Desktop)
}

// VisibleWidgets returns visible widgets sorted by order from the provided widget slice.
func VisibleWidgets(widgets []models.DashboardWidget) []models.DashboardWidget {
	var visible []models.DashboardWidget
	for _, w := range widgets {
		if w.Visible {
			visible = append(visible, w)
		}
	}
	sort.Slice(visible, func(i, j int) bool {
		return visible[i].Order < visible[j].Order
	})
	return visible
}
