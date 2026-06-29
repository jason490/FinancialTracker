package queries

import (
    "FinancialTracker/internal/models"
    "FinancialTracker/internal/utils"
    "database/sql"
    "github.com/google/uuid"
)

// CreatePlaidAccount inserts a new plaid account record
func CreatePlaidAccount(db *sql.DB, a *models.Account) error {
    if a.Status == "" {
        a.Status = "active"
    }
	if a.RowID == "" {
		a.RowID = uuid.New().String()
	}
    query := `INSERT INTO plaid_account (user_id, row_id, plaid_account_id, plaid_item_id, name, mask, type, subtype, balance, available_balance, currency, status, is_hidden, monthly_payment) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
    result, err := db.Exec(query, a.UserID, a.RowID, a.PlaidAccountID, a.PlaidItemID, a.Name, a.Mask, a.Type, a.Subtype, a.Balance, a.AvailableBalance, a.Currency, a.Status, a.IsHidden, a.MonthlyPayment)
    if err != nil {
        return err
    }
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    a.ID = id
    return nil
}

// UpdatePlaidAccount updates an existing plaid account record
func UpdatePlaidAccount(db *sql.DB, a *models.Account) error {
    if a.Status == "" {
        a.Status = "active"
    }
    query := `UPDATE plaid_account 
              SET name = ?, mask = ?, type = ?, subtype = ?, balance = ?, available_balance = ?, currency = ?, status = ?, is_hidden = ?
              WHERE plaid_account_id = ? AND user_id = ?`
    _, err := db.Exec(query, a.Name, a.Mask, a.Type, a.Subtype, a.Balance, a.AvailableBalance, a.Currency, a.Status, a.IsHidden, a.PlaidAccountID, a.UserID)
    return err
}

// CreatePlaidItem inserts a new plaid item record
func CreatePlaidItem(db *sql.DB, item *models.PlaidItem) error {
    if item.RowID == "" {
        item.RowID = uuid.New().String()
    }
    if item.Status == "" {
        item.Status = "active"
    }

    encryptedToken, err := utils.Encrypt(item.AccessToken)
    if err != nil {
        return err
    }

    query := `INSERT INTO plaid_items (row_id, user_id, plaid_item_id, access_token, institution_id, institution_name, sync_cursor, status, error_code) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
    result, err := db.Exec(query, item.RowID, item.UserID, item.PlaidItemID, encryptedToken, item.InstitutionID, item.InstitutionName, item.SyncCursor, item.Status, item.ErrorCode)
    if err != nil {
        return err
    }
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    item.ID = id
    return nil
}

// GetPlaidItemByItemID retrieves a plaid item by its Plaid item ID
func GetPlaidItemByItemID(db *sql.DB, itemID string) (*models.PlaidItem, error) {
    query := `SELECT id, row_id, user_id, plaid_item_id, access_token, institution_id, institution_name, sync_cursor, status, error_code, last_synced, created_at 
              FROM plaid_items WHERE plaid_item_id = ?`
    var i models.PlaidItem
    var encryptedToken string
    err := db.QueryRow(query, itemID).Scan(&i.ID, &i.RowID, &i.UserID, &i.PlaidItemID, &encryptedToken, &i.InstitutionID, &i.InstitutionName, &i.SyncCursor, &i.Status, &i.ErrorCode, &i.LastSynced, &i.CreatedAt)
    if err != nil {
        return nil, err
    }

    decryptedToken, err := utils.Decrypt(encryptedToken)
    if err != nil {
        return nil, err
    }
    i.AccessToken = decryptedToken

    return &i, nil
}

// UpdatePlaidItem updates an existing plaid item's access token and metadata
func UpdatePlaidItem(db *sql.DB, item *models.PlaidItem) error {
    if item.Status == "" {
        item.Status = "active"
    }

    encryptedToken, err := utils.Encrypt(item.AccessToken)
    if err != nil {
        return err
    }

    query := `UPDATE plaid_items 
              SET access_token = ?, institution_id = ?, institution_name = ?, status = ?, error_code = ? 
              WHERE plaid_item_id = ? AND user_id = ?`
    _, err = db.Exec(query, encryptedToken, item.InstitutionID, item.InstitutionName, item.Status, item.ErrorCode, item.PlaidItemID, item.UserID)
    return err
}

// GetPlaidItemByRowID retrieves a plaid item by its Row ID and ensures it belongs to the user
func GetPlaidItemByRowID(db *sql.DB, rowID string, userID int64) (*models.PlaidItem, error) {
    query := `SELECT id, row_id, user_id, plaid_item_id, access_token, institution_id, institution_name, sync_cursor, status, error_code, last_synced, created_at 
              FROM plaid_items WHERE row_id = ? AND user_id = ?`
    var i models.PlaidItem
    var encryptedToken string
    err := db.QueryRow(query, rowID, userID).Scan(&i.ID, &i.RowID, &i.UserID, &i.PlaidItemID, &encryptedToken, &i.InstitutionID, &i.InstitutionName, &i.SyncCursor, &i.Status, &i.ErrorCode, &i.LastSynced, &i.CreatedAt)
    if err != nil {
        return nil, err
    }

    decryptedToken, err := utils.Decrypt(encryptedToken)
    if err != nil {
        return nil, err
    }
    i.AccessToken = decryptedToken

    return &i, nil
}

// GetPlaidItemsByUserID retrieves all plaid items for a specific user
func GetPlaidItemsByUserID(db *sql.DB, userID int64) ([]models.PlaidItem, error) {
    query := `SELECT id, row_id, user_id, plaid_item_id, access_token, institution_id, institution_name, sync_cursor, status, error_code, last_synced, created_at 
              FROM plaid_items WHERE user_id = ?`
    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.PlaidItem
    for rows.Next() {
        var i models.PlaidItem
        var encryptedToken string
        if err := rows.Scan(&i.ID, &i.RowID, &i.UserID, &i.PlaidItemID, &encryptedToken, &i.InstitutionID, &i.InstitutionName, &i.SyncCursor, &i.Status, &i.ErrorCode, &i.LastSynced, &i.CreatedAt); err != nil {
            return nil, err
        }

        decryptedToken, err := utils.Decrypt(encryptedToken)
        if err != nil {
            return nil, err
        }
        i.AccessToken = decryptedToken

        items = append(items, i)
    }
    return items, nil
}

// UpdatePlaidItemStatus updates the connection status and last Plaid error code for an item
func UpdatePlaidItemStatus(db *sql.DB, itemID string, status string, errorCode string) error {
    query := `UPDATE plaid_items SET status = ?, error_code = ? WHERE plaid_item_id = ?`
    _, err := db.Exec(query, status, errorCode, itemID)
    return err
}

// MarkPlaidAccountsDisconnectedByItemID marks all accounts for an item as disconnected
func MarkPlaidAccountsDisconnectedByItemID(db *sql.DB, itemID string) error {
    query := `UPDATE plaid_account SET status = 'disconnected' WHERE plaid_item_id = ? AND status != 'disconnected'`
    _, err := db.Exec(query, itemID)
    return err
}

// UpdatePlaidItemLastSynced records when an item was last successfully synced
func UpdatePlaidItemLastSynced(db *sql.DB, itemID string, syncedAt int64) error {
    query := `UPDATE plaid_items SET last_synced = ? WHERE plaid_item_id = ?`
    _, err := db.Exec(query, syncedAt, itemID)
    return err
}

// UpdatePlaidItemCursor updates the sync cursor for a plaid item
func UpdatePlaidItemCursor(db *sql.DB, itemID string, cursor string) error {
    query := `UPDATE plaid_items SET sync_cursor = ? WHERE plaid_item_id = ?`
    _, err := db.Exec(query, cursor, itemID)
    return err
}

// GetPlaidAccountsByUserID retrieves all plaid accounts for a specific user
func GetPlaidAccountsByUserID(db *sql.DB, userID int64) ([]models.Account, error) {
    query := `SELECT id, row_id, user_id, plaid_account_id, plaid_item_id, name, mask, type, subtype, balance, available_balance, currency, status, is_hidden, monthly_payment, created_at FROM plaid_account WHERE user_id = ?`
    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var plaidAccounts []models.Account
    for rows.Next() {
        var a models.Account
        if err := rows.Scan(&a.ID, &a.RowID, &a.UserID, &a.PlaidAccountID, &a.PlaidItemID, &a.Name, &a.Mask, &a.Type, &a.Subtype, &a.Balance, &a.AvailableBalance, &a.Currency, &a.Status, &a.IsHidden, &a.MonthlyPayment, &a.CreatedAt); err != nil {
            return nil, err
        }
        plaidAccounts = append(plaidAccounts, a)
    }
    return plaidAccounts, nil
}

// GetAccountByRowID retrieves an account by its row_id
func GetAccountByRowID(db *sql.DB, rowID string) (*models.Account, error) {
    query := `SELECT id, row_id, user_id, plaid_account_id, plaid_item_id, name, mask, type, subtype, balance, available_balance, currency, status, is_hidden, monthly_payment, created_at FROM plaid_account WHERE row_id = ?`
    var a models.Account
    err := db.QueryRow(query, rowID).Scan(&a.ID, &a.RowID, &a.UserID, &a.PlaidAccountID, &a.PlaidItemID, &a.Name, &a.Mask, &a.Type, &a.Subtype, &a.Balance, &a.AvailableBalance, &a.Currency, &a.Status, &a.IsHidden, &a.MonthlyPayment, &a.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &a, nil
}

// GetAccountByPlaidAccountID retrieves an account by its plaid_account_id
func GetAccountByPlaidAccountID(db *sql.DB, plaidAccountID string) (*models.Account, error) {
    query := `SELECT id, row_id, user_id, plaid_account_id, plaid_item_id, name, mask, type, subtype, balance, available_balance, currency, status, is_hidden, monthly_payment, created_at FROM plaid_account WHERE plaid_account_id = ?`
    var a models.Account
    err := db.QueryRow(query, plaidAccountID).Scan(&a.ID, &a.RowID, &a.UserID, &a.PlaidAccountID, &a.PlaidItemID, &a.Name, &a.Mask, &a.Type, &a.Subtype, &a.Balance, &a.AvailableBalance, &a.Currency, &a.Status, &a.IsHidden, &a.MonthlyPayment, &a.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &a, nil
}

// GetPlaidAccountsByItemID retrieves all plaid accounts for a specific plaid item ID
func GetPlaidAccountsByItemID(db *sql.DB, itemID string) ([]models.Account, error) {
	query := `SELECT id, row_id, user_id, plaid_account_id, plaid_item_id, name, mask, type, subtype, balance, available_balance, currency, status, is_hidden, monthly_payment, created_at FROM plaid_account WHERE plaid_item_id = ?`
	rows, err := db.Query(query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plaidAccounts []models.Account
	for rows.Next() {
		var a models.Account
		if err := rows.Scan(&a.ID, &a.RowID, &a.UserID, &a.PlaidAccountID, &a.PlaidItemID, &a.Name, &a.Mask, &a.Type, &a.Subtype, &a.Balance, &a.AvailableBalance, &a.Currency, &a.Status, &a.IsHidden, &a.MonthlyPayment, &a.CreatedAt); err != nil {
			return nil, err
		}
		plaidAccounts = append(plaidAccounts, a)
	}
	return plaidAccounts, nil
}

// UpdateAccountMonthlyPayment sets the monthly payment for a loan account.
func UpdateAccountMonthlyPayment(db *sql.DB, plaidAccountID string, monthlyPayment float64) error {
	query := `UPDATE plaid_account SET monthly_payment = ? WHERE plaid_account_id = ?`
	_, err := db.Exec(query, monthlyPayment, plaidAccountID)
	return err
}

// UpdatePlaidAccountStatus updates the status of a specific plaid account
func UpdatePlaidAccountStatus(db *sql.DB, accountID string, status string) error {
	query := `UPDATE plaid_account SET status = ? WHERE plaid_account_id = ?`
	_, err := db.Exec(query, status, accountID)
	return err
}

// CountActivePlaidItems returns how many non-disconnected Plaid items a user has linked.
func CountActivePlaidItems(db *sql.DB, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM plaid_items WHERE user_id = ? AND status != 'disconnected'`
	var count int
	err := db.QueryRow(query, userID).Scan(&count)
	return count, err
}

// DeletePlaidItem deletes a plaid item by its Row ID and User ID for security
func DeletePlaidItem(db *sql.DB, rowID string, userID int64) error {
	query := `DELETE FROM plaid_items WHERE row_id = ? AND user_id = ?`
	_, err := db.Exec(query, rowID, userID)
	return err
}

// DeletePlaidAccount deletes a specific plaid account and its transactions
func DeletePlaidAccount(db *sql.DB, accountID string, userID int64) error {
	query := `DELETE FROM plaid_account WHERE row_id = ? AND user_id = ?`
	_, err := db.Exec(query, accountID, userID)
	return err
}

// ToggleAccountVisibility flips the is_hidden status of a plaid account and returns the new state
func ToggleAccountVisibility(db *sql.DB, accountID string, userID int64) (bool, error) {
    query := `UPDATE plaid_account SET is_hidden = NOT is_hidden WHERE row_id = ? AND user_id = ? RETURNING is_hidden`
    var isHidden bool
    err := db.QueryRow(query, accountID, userID).Scan(&isHidden)
    return isHidden, err
}
