package queries

import (
	"FinancialTracker/internal/models"
	"database/sql"
	"time"
)

func analyticsAccountJoin(provider string) (table, alias string) {
	if provider == "stripe" {
		return "stripe_fc_account", "s"
	}
	return "plaid_account", "p"
}

func analyticsAccountFilter(alias string) string {
	return alias + ".user_id = ? AND " + alias + ".is_hidden = 0"
}

// GetMonthlySpending returns positive outflows grouped by month for the last N months.
func GetMonthlySpending(db *sql.DB, userID int64, provider string, months int) ([]models.MonthlySpend, error) {
	if months < 1 {
		months = 6
	}
	start := time.Now().AddDate(0, -months, 0).Unix()
	accountTable, accountAlias := analyticsAccountJoin(provider)
	query := `
		SELECT strftime('%Y-%m', t.date, 'unixepoch') AS month_key,
		       SUM(t.amount) AS total
		FROM transactions t
		JOIN ` + accountTable + ` ` + accountAlias + ` ON t.plaid_id = ` + accountAlias + `.id
		WHERE ` + analyticsAccountFilter(accountAlias) + `
		  AND t.provider = ?
		  AND t.amount > 0
		  AND t.date >= ?
		GROUP BY month_key
		ORDER BY month_key ASC`

	rows, err := db.Query(query, userID, provider, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.MonthlySpend
	for rows.Next() {
		var m models.MonthlySpend
		if err := rows.Scan(&m.Month, &m.Total); err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
}

// GetMonthCashflow returns total spend and income for a date range (unix seconds).
func GetMonthCashflow(db *sql.DB, userID int64, provider string, monthStart, monthEnd int64) (models.MonthCashflow, error) {
	var cf models.MonthCashflow
	accountTable, accountAlias := analyticsAccountJoin(provider)
	accountFilter := analyticsAccountFilter(accountAlias)

	spendQuery := `
		SELECT COALESCE(SUM(t.amount), 0)
		FROM transactions t
		JOIN ` + accountTable + ` ` + accountAlias + ` ON t.plaid_id = ` + accountAlias + `.id
		WHERE ` + accountFilter + `
		  AND t.provider = ?
		  AND t.amount > 0
		  AND t.date >= ? AND t.date < ?`
	if err := db.QueryRow(spendQuery, userID, provider, monthStart, monthEnd).Scan(&cf.Spend); err != nil {
		return cf, err
	}
	incomeQuery := `
		SELECT COALESCE(SUM(-t.amount), 0)
		FROM transactions t
		JOIN ` + accountTable + ` ` + accountAlias + ` ON t.plaid_id = ` + accountAlias + `.id
		WHERE ` + accountFilter + `
		  AND t.provider = ?
		  AND t.amount < 0
		  AND t.date >= ? AND t.date < ?`
	if err := db.QueryRow(incomeQuery, userID, provider, monthStart, monthEnd).Scan(&cf.Income); err != nil {
		return cf, err
	}
	return cf, nil
}

// GetSpendingByTag returns this month's spending split evenly across tags per transaction.
func GetSpendingByTag(db *sql.DB, userID int64, provider string, monthStart, monthEnd int64) ([]models.TagBreakdown, error) {
	return getTaggedBreakdown(db, userID, provider, monthStart, monthEnd, true)
}

// GetIncomeByTag returns this month's income split evenly across tags per transaction.
func GetIncomeByTag(db *sql.DB, userID int64, provider string, monthStart, monthEnd int64) ([]models.TagBreakdown, error) {
	return getTaggedBreakdown(db, userID, provider, monthStart, monthEnd, false)
}

func getTaggedBreakdown(db *sql.DB, userID int64, provider string, monthStart, monthEnd int64, spending bool) ([]models.TagBreakdown, error) {
	amountExpr := "t.amount / tag_count.cnt"
	amountFilter := "t.amount > 0"
	if !spending {
		amountExpr = "(-t.amount) / tag_count.cnt"
		amountFilter = "t.amount < 0"
	}

	accountTable, accountAlias := analyticsAccountJoin(provider)
	accountFilter := analyticsAccountFilter(accountAlias)

	taggedQuery := `
		SELECT tg.id, tg.name, tg.color, SUM(` + amountExpr + `) AS total
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
		WHERE ` + accountFilter + `
		  AND t.provider = ?
		  AND c.user_id = ?
		  AND ` + amountFilter + `
		  AND t.date >= ? AND t.date < ?
		GROUP BY tg.id, tg.name, tg.color
		HAVING total > 0.001
		ORDER BY total DESC`

	rows, err := db.Query(taggedQuery, userID, provider, userID, monthStart, monthEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.TagBreakdown
	for rows.Next() {
		var b models.TagBreakdown
		if err := rows.Scan(&b.TagID, &b.TagName, &b.Color, &b.Total); err != nil {
			return nil, err
		}
		result = append(result, b)
	}

	untaggedSum := "t.amount"
	if !spending {
		untaggedSum = "-t.amount"
	}
	untaggedQuery := `
		SELECT COALESCE(SUM(` + untaggedSum + `), 0)
		FROM transactions t
		JOIN ` + accountTable + ` ` + accountAlias + ` ON t.plaid_id = ` + accountAlias + `.id
		LEFT JOIN transaction_tags tt ON t.id = tt.transaction_id
		WHERE ` + accountFilter + `
		  AND t.provider = ?
		  AND ` + amountFilter + `
		  AND t.date >= ? AND t.date < ?
		  AND tt.transaction_id IS NULL`

	var untagged float64
	if err := db.QueryRow(untaggedQuery, userID, provider, monthStart, monthEnd).Scan(&untagged); err != nil {
		return nil, err
	}
	if untagged > 0.001 {
		result = append(result, models.TagBreakdown{
			TagID:   0,
			TagName: "Uncategorized",
			Color:   "slate",
			Total:   untagged,
		})
	}
	return result, nil
}
