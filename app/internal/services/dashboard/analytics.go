package dashboard

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/utils"
)

// ApplySpendingAnalytics fills spending trend and derived summary fields on dashboard data.
func ApplySpendingAnalytics(data *models.DashboardData, trend []models.MonthlySpend) {
	data.SpendingTrend = trend
	if len(trend) == 0 {
		return
	}
	var total float64
	for _, m := range trend {
		total += m.Total
	}
	data.Summary.AvgMonthlySpend = total / float64(len(trend))

	liquid := data.Summary.Cash + data.Summary.Savings
	if data.Summary.AvgMonthlySpend > 0 && liquid > 0 {
		data.Summary.MonthsToZero = liquid / data.Summary.AvgMonthlySpend
	}
}

// ApplyMonthAnalytics attaches current-month cashflow and tag donut data.
func ApplyMonthAnalytics(data *models.DashboardData, cashflow models.MonthCashflow, spending, income []models.TagBreakdown) {
	data.MonthCashflow = cashflow
	data.SpendingByTag = normalizeBreakdownColors(spending)
	data.IncomeByTag = normalizeBreakdownColors(income)
}

func normalizeBreakdownColors(rows []models.TagBreakdown) []models.TagBreakdown {
	for i := range rows {
		rows[i].Color = utils.NormalizeTagColor(rows[i].Color)
	}
	return rows
}
