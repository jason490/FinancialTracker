package queries

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/utils"
	"database/sql"
	"time"
)

// GetFilteredCashflow returns total spend and income for transactions matching the filters.
// Transactions carrying a hidden tag are excluded so totals match dashboard analytics.
func GetFilteredCashflow(db *sql.DB, userID int64, provider string, f models.TransactionFilters) (models.MonthCashflow, error) {
	var cf models.MonthCashflow
	accountTable, accountAlias := accountJoinForProvider(provider)
	where, args := transactionFilterClause(accountAlias, userID, provider, f)
	from := `FROM transactions t JOIN ` + accountTable + ` ` + accountAlias + ` ON t.plaid_id = ` + accountAlias + `.id `

	spendQuery := `SELECT COALESCE(SUM(t.amount), 0) ` + from + where + ` AND t.amount > 0` + hiddenTagExclusion
	if err := db.QueryRow(spendQuery, args...).Scan(&cf.Spend); err != nil {
		return cf, err
	}

	incomeQuery := `SELECT COALESCE(SUM(-t.amount), 0) ` + from + where + ` AND t.amount < 0` + hiddenTagExclusion
	if err := db.QueryRow(incomeQuery, args...).Scan(&cf.Income); err != nil {
		return cf, err
	}
	return cf, nil
}

// GetFilteredSpendingByCategory returns spending split across tag categories for the filtered set.
// Multi-tag transactions split evenly across their tags (then rolled up by category).
func GetFilteredSpendingByCategory(db *sql.DB, userID int64, provider string, f models.TransactionFilters) ([]models.CategoryBreakdown, error) {
	accountTable, accountAlias := accountJoinForProvider(provider)
	where, args := transactionFilterClause(accountAlias, userID, provider, f)

	taggedQuery := `
		SELECT c.id, c.name, SUM(t.amount / tag_count.cnt) AS total
		FROM transactions t
		JOIN ` + accountTable + ` ` + accountAlias + ` ON t.plaid_id = ` + accountAlias + `.id
		JOIN transaction_tags tt ON t.id = tt.transaction_id
		JOIN tags tg ON tt.tag_id = tg.id
		JOIN categories c ON tg.category_id = c.id
		JOIN (
			SELECT transaction_id, COUNT(*) AS cnt
			FROM transaction_tags
			GROUP BY transaction_id
		) tag_count ON tag_count.transaction_id = t.id
		` + where + `
		  AND c.user_id = ?
		  AND t.amount > 0` + hiddenTagExclusion + `
		GROUP BY c.id, c.name
		HAVING total > 0.001
		ORDER BY total DESC`

	taggedArgs := append(append([]interface{}{}, args...), userID)
	rows, err := db.Query(taggedQuery, taggedArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.CategoryBreakdown
	for rows.Next() {
		var b models.CategoryBreakdown
		var categoryID int64
		if err := rows.Scan(&categoryID, &b.CategoryName, &b.Total); err != nil {
			return nil, err
		}
		b.Color = utils.PaletteColorForIndex(int(categoryID))
		result = append(result, b)
	}

	untaggedQuery := `
		SELECT COALESCE(SUM(t.amount), 0)
		FROM transactions t
		JOIN ` + accountTable + ` ` + accountAlias + ` ON t.plaid_id = ` + accountAlias + `.id
		LEFT JOIN transaction_tags tt ON t.id = tt.transaction_id
		` + where + `
		  AND t.amount > 0
		  AND tt.transaction_id IS NULL` + hiddenTagExclusion

	var untagged float64
	if err := db.QueryRow(untaggedQuery, args...).Scan(&untagged); err != nil {
		return nil, err
	}
	if untagged > 0.001 {
		result = append(result, models.CategoryBreakdown{
			CategoryName: "Uncategorized",
			Color:        "slate",
			Total:        untagged,
		})
	}
	return result, nil
}

// GetFilteredMonthlyTrend returns monthly income and spending for the filtered set.
// When no date filter is set, defaults to the last 6 months.
func GetFilteredMonthlyTrend(db *sql.DB, userID int64, provider string, f models.TransactionFilters) ([]models.MonthlySpend, error) {
	trendFilters := f
	if trendFilters.StartDate == nil && trendFilters.EndDate == nil {
		start := time.Now().AddDate(0, -6, 0).Unix()
		trendFilters.StartDate = &start
	}

	accountTable, accountAlias := accountJoinForProvider(provider)
	where, args := transactionFilterClause(accountAlias, userID, provider, trendFilters)

	query := `
		SELECT strftime('%Y-%m', t.date, 'unixepoch') AS month_key,
		       COALESCE(SUM(CASE WHEN t.amount > 0 THEN t.amount ELSE 0 END), 0) AS total,
		       COALESCE(SUM(CASE WHEN t.amount < 0 THEN -t.amount ELSE 0 END), 0) AS income
		FROM transactions t
		JOIN ` + accountTable + ` ` + accountAlias + ` ON t.plaid_id = ` + accountAlias + `.id
		` + where + hiddenTagExclusion + `
		GROUP BY month_key
		HAVING total > 0 OR income > 0
		ORDER BY month_key ASC`

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.MonthlySpend
	for rows.Next() {
		var m models.MonthlySpend
		if err := rows.Scan(&m.Month, &m.Total, &m.Income); err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
}
