package queries

import (
    "FinancialTracker/internal/models"
    "database/sql"
)

// CreateTransaction inserts a new transaction record
func CreateTransaction(db *sql.DB, t *models.Transaction) error {
	provider := t.Provider
	if provider == "" {
		provider = "plaid"
	}
	query := `INSERT INTO transactions (provider, plaid_id, plaid_transaction_id, date, amount, name, merchant_name, plaid_category, pending) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, provider, t.PlaidID, t.PlaidTransactionID, t.Date, t.Amount, t.Name, t.MerchantName, t.PlaidCategory, t.Pending)
    if err != nil {
        return err
    }
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    t.ID = id
    return nil
}

// GetTransactions retrieves transactions for a user with filtering and pagination.
func GetTransactions(db *sql.DB, userID int64, provider string, f models.TransactionFilters) ([]models.Transaction, int, error) {
	accountTable, accountAlias := accountJoinForProvider(provider)
	where := "WHERE " + accountAlias + ".user_id = ? AND " + accountAlias + ".is_hidden = 0 AND t.provider = ?"
	args := []interface{}{userID, provider}

	if f.Search != "" {
		where += " AND (t.name LIKE ? OR t.merchant_name LIKE ? OR t.plaid_category LIKE ?)"
		args = append(args, "%"+f.Search+"%", "%"+f.Search+"%", "%"+f.Search+"%")
	}
	// Filter by amount magnitude so users can specify min/max regardless of
	// transaction sign convention (Plaid/Stripe stores expenses as positive
	// and income/refunds as negative). The displayed amount in the UI is the
	// inverse of the stored sign, so filtering by ABS keeps the filter
	// intuitive: "max = 1" returns only transactions of magnitude <= $1,
	// "min = 50" returns only transactions of magnitude >= $50.
	if f.MinAmount != nil {
		where += " AND ABS(t.amount) >= ?"
		args = append(args, *f.MinAmount)
	}
	if f.MaxAmount != nil {
		where += " AND ABS(t.amount) <= ?"
		args = append(args, *f.MaxAmount)
	}
	if f.StartDate != nil {
		where += " AND t.date >= ?"
		args = append(args, *f.StartDate)
	}
	if f.EndDate != nil {
		where += " AND t.date <= ?"
		args = append(args, *f.EndDate)
	}

	// Filter by Category
	if f.CategoryID != nil {
		where += ` AND t.id IN (
			SELECT tt.transaction_id 
			FROM transaction_tags tt 
			JOIN tags tg ON tt.tag_id = tg.id 
			WHERE tg.category_id = ?
		)`
		args = append(args, *f.CategoryID)
	}

	// For tag filtering, we need a subquery or join
	if len(f.Tags) > 0 {
		where += " AND t.id IN (SELECT transaction_id FROM transaction_tags WHERE tag_id IN ("
		for i, tagID := range f.Tags {
			where += "?"
			args = append(args, tagID)
			if i < len(f.Tags)-1 {
				where += ","
			}
		}
		where += "))"
	}

	countQuery := `SELECT COUNT(*) FROM transactions t JOIN ` + accountTable + ` ` + accountAlias + ` ON t.plaid_id = ` + accountAlias + `.id ` + where
	var totalCount int
	err := db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	orderBy := "t.date"
	if f.SortBy != "" {
		switch f.SortBy {
		case "amount":
			orderBy = "t.amount"
		case "name":
			orderBy = "t.name"
		}
	}
	orderDir := "DESC"
	if f.SortDir == "asc" {
		orderDir = "ASC"
	}

	query := `SELECT t.id, t.provider, t.plaid_id, t.plaid_transaction_id, t.date, t.amount, t.name, t.merchant_name, t.plaid_category, t.pending, t.created_at 
              FROM transactions t 
              JOIN ` + accountTable + ` ` + accountAlias + ` ON t.plaid_id = ` + accountAlias + `.id 
              ` + where + ` ORDER BY ` + orderBy + ` ` + orderDir + ` LIMIT ? OFFSET ?`

	args = append(args, f.Limit, f.Offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.Provider, &t.PlaidID, &t.PlaidTransactionID, &t.Date, &t.Amount, &t.Name, &t.MerchantName, &t.PlaidCategory, &t.Pending, &t.CreatedAt); err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, t)
	}

	// Load tags for these transactions
	for i := range transactions {
		tags, err := GetTransactionTags(db, transactions[i].ID)
		if err != nil {
			return nil, 0, err
		}
		transactions[i].Tags = tags
	}

	return transactions, totalCount, nil
}

// GetTransactionByPlaidID retrieves a transaction by its external transaction ID.
func GetTransactionByPlaidID(db *sql.DB, plaidTransactionID string) (*models.Transaction, error) {
    query := `SELECT id, provider, plaid_id, plaid_transaction_id, date, amount, name, merchant_name, plaid_category, pending, created_at 
              FROM transactions WHERE plaid_transaction_id = ?`
    var t models.Transaction
    err := db.QueryRow(query, plaidTransactionID).Scan(&t.ID, &t.Provider, &t.PlaidID, &t.PlaidTransactionID, &t.Date, &t.Amount, &t.Name, &t.MerchantName, &t.PlaidCategory, &t.Pending, &t.CreatedAt)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &t, err
}

// UpdateTransaction updates an existing transaction record
func UpdateTransaction(db *sql.DB, t *models.Transaction) error {
    query := `UPDATE transactions SET amount = ?, name = ?, merchant_name = ?, plaid_category = ?, pending = ?, date = ? 
              WHERE plaid_transaction_id = ?`
    _, err := db.Exec(query, t.Amount, t.Name, t.MerchantName, t.PlaidCategory, t.Pending, t.Date, t.PlaidTransactionID)
    return err
}

// DeleteTransactionByPlaidID deletes a transaction by its Plaid transaction ID
func DeleteTransactionByPlaidID(db *sql.DB, plaidTransactionID string) error {
    query := `DELETE FROM transactions WHERE plaid_transaction_id = ?`
    _, err := db.Exec(query, plaidTransactionID)
    return err
}

// GetTransactionTags retrieves tags for a specific transaction
func GetTransactionTags(db *sql.DB, transactionID int64) ([]models.Tag, error) {
	query := `SELECT tg.id, tg.category_id, tg.name, tg.color, tg.created_at 
              FROM tags tg 
              JOIN transaction_tags tt ON tg.id = tt.tag_id 
              WHERE tt.transaction_id = ?`
	rows, err := db.Query(query, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.ID, &t.CategoryID, &t.Name, &t.Color, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
    return tags, nil
}

func accountJoinForProvider(provider string) (table, alias string) {
	if provider == "stripe" {
		return "stripe_fc_account", "s"
	}
	return "plaid_account", "p"
}
