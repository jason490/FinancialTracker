package models

// DashboardWidget describes a single dashboard panel in the user's layout.
type DashboardWidget struct {
	ID      string `json:"id"`
	Visible bool   `json:"visible"`
	Order   int    `json:"order"`
}

// DashboardLayout is the persisted widget configuration for a user.
type DashboardLayout struct {
	Desktop []DashboardWidget `json:"desktop"`
	Mobile  []DashboardWidget `json:"mobile"`
}

// MonthlySpend is spending total for a calendar month (YYYY-MM).
type MonthlySpend struct {
	Month string  `json:"month"`
	Total float64 `json:"total"`
}

// MonthCashflow is spend and income for a calendar month.
type MonthCashflow struct {
	Spend  float64 `json:"spend"`
	Income float64 `json:"income"`
}

// DashboardSummary holds aggregated balances by account category.
type DashboardSummary struct {
	NetWorth            float64 `json:"net_worth"`
	Cash                float64 `json:"cash"`
	Savings             float64 `json:"savings"`
	CreditDebt          float64 `json:"credit_debt"`
	LoanDebt            float64 `json:"loan_debt"`
	LoanMonthlyPayments float64 `json:"loan_monthly_payments"`
	Investments         float64 `json:"investments"`
	AccountCount        int     `json:"account_count"`
	AvgMonthlySpend     float64 `json:"avg_monthly_spend"`
	MonthsToZero        float64 `json:"months_to_zero"`
}

// DashboardData is the view model for the dashboard page and widget partials.
type DashboardData struct {
	Summary       DashboardSummary     `json:"summary"`
	Groups        map[string][]Account `json:"groups"`
	Transactions  []Transaction        `json:"transactions"`
	SpendingTrend []MonthlySpend       `json:"spending_trend"`
	MonthCashflow MonthCashflow        `json:"month_cashflow"`
	SpendingByTag []TagBreakdown       `json:"spending_by_tag"`
	IncomeByTag   []TagBreakdown       `json:"income_by_tag"`
	Layout        DashboardLayout      `json:"layout"`
	EditMode      bool                 `json:"edit_mode"`
}
