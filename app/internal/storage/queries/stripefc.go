package queries

import (
	"FinancialTracker/internal/models"
	"database/sql"

	"github.com/google/uuid"
)

// CreateStripeFCItem inserts a Stripe FC institution grouping.
func CreateStripeFCItem(db *sql.DB, item *models.StripeFCItem) error {
	if item.RowID == "" {
		item.RowID = uuid.New().String()
	}
	if item.Status == "" {
		item.Status = "active"
	}
	query := `INSERT INTO stripe_fc_items (row_id, user_id, institution_name, status, error_code, transaction_refresh_id)
              VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, item.RowID, item.UserID, item.InstitutionName, item.Status, item.ErrorCode, item.TransactionRefreshID)
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

// GetStripeFCItemByRowID retrieves a Stripe FC item by row ID for a user.
func GetStripeFCItemByRowID(db *sql.DB, rowID string, userID int64) (*models.StripeFCItem, error) {
	query := `SELECT id, row_id, user_id, institution_name, status, error_code, last_synced, transaction_refresh_id, created_at
              FROM stripe_fc_items WHERE row_id = ? AND user_id = ?`
	var item models.StripeFCItem
	err := db.QueryRow(query, rowID, userID).Scan(
		&item.ID, &item.RowID, &item.UserID, &item.InstitutionName, &item.Status, &item.ErrorCode,
		&item.LastSynced, &item.TransactionRefreshID, &item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetStripeFCItemByInstitution retrieves a Stripe FC item by institution name for a user.
func GetStripeFCItemByInstitution(db *sql.DB, userID int64, institutionName string) (*models.StripeFCItem, error) {
	query := `SELECT id, row_id, user_id, institution_name, status, error_code, last_synced, transaction_refresh_id, created_at
              FROM stripe_fc_items WHERE user_id = ? AND institution_name = ?`
	var item models.StripeFCItem
	err := db.QueryRow(query, userID, institutionName).Scan(
		&item.ID, &item.RowID, &item.UserID, &item.InstitutionName, &item.Status, &item.ErrorCode,
		&item.LastSynced, &item.TransactionRefreshID, &item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetStripeFCItemsByUserID returns all Stripe FC items for a user.
func GetStripeFCItemsByUserID(db *sql.DB, userID int64) ([]models.StripeFCItem, error) {
	query := `SELECT id, row_id, user_id, institution_name, status, error_code, last_synced, transaction_refresh_id, created_at
              FROM stripe_fc_items WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.StripeFCItem
	for rows.Next() {
		var item models.StripeFCItem
		if err := rows.Scan(
			&item.ID, &item.RowID, &item.UserID, &item.InstitutionName, &item.Status, &item.ErrorCode,
			&item.LastSynced, &item.TransactionRefreshID, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// UpdateStripeFCItemStatus updates status and error code for a Stripe FC item.
func UpdateStripeFCItemStatus(db *sql.DB, rowID string, status, errorCode string) error {
	query := `UPDATE stripe_fc_items SET status = ?, error_code = ? WHERE row_id = ?`
	_, err := db.Exec(query, status, errorCode, rowID)
	return err
}

// UpdateStripeFCItemLastSynced updates the last sync timestamp for a Stripe FC item.
func UpdateStripeFCItemLastSynced(db *sql.DB, rowID string, syncedAt int64) error {
	query := `UPDATE stripe_fc_items SET last_synced = ? WHERE row_id = ?`
	_, err := db.Exec(query, syncedAt, rowID)
	return err
}

// UpdateStripeFCItemTransactionRefresh stores the latest transaction refresh token.
func UpdateStripeFCItemTransactionRefresh(db *sql.DB, rowID, refreshID string) error {
	query := `UPDATE stripe_fc_items SET transaction_refresh_id = ? WHERE row_id = ?`
	_, err := db.Exec(query, refreshID, rowID)
	return err
}

// DeleteStripeFCItem removes a Stripe FC item and cascades accounts.
func DeleteStripeFCItem(db *sql.DB, rowID string, userID int64) error {
	query := `DELETE FROM stripe_fc_items WHERE row_id = ? AND user_id = ?`
	_, err := db.Exec(query, rowID, userID)
	return err
}

// CountActiveStripeFCItems counts non-disconnected Stripe FC items for a user.
func CountActiveStripeFCItems(db *sql.DB, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM stripe_fc_items WHERE user_id = ? AND status != 'disconnected'`
	var count int
	err := db.QueryRow(query, userID).Scan(&count)
	return count, err
}

// CreateStripeFCAccount inserts a Stripe FC account record.
func CreateStripeFCAccount(db *sql.DB, account *models.StripeFCAccount) error {
	if account.Status == "" {
		account.Status = "active"
	}
	if account.RowID == "" {
		account.RowID = uuid.New().String()
	}
	query := `INSERT INTO stripe_fc_account (user_id, row_id, stripe_account_id, stripe_item_row_id, name, mask, type, subtype, balance, available_balance, currency, status, is_hidden)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query,
		account.UserID, account.RowID, account.StripeAccountID, account.StripeItemRowID, account.Name, account.Mask,
		account.Type, account.Subtype, account.Balance, account.AvailableBalance, account.Currency,
		account.Status, account.IsHidden,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	account.ID = id
	return nil
}

// UpdateStripeFCAccount updates an existing Stripe FC account record.
func UpdateStripeFCAccount(db *sql.DB, account *models.StripeFCAccount) error {
	if account.Status == "" {
		account.Status = "active"
	}
	query := `UPDATE stripe_fc_account
              SET name = ?, mask = ?, type = ?, subtype = ?, balance = ?, available_balance = ?, currency = ?, status = ?, is_hidden = ?
              WHERE stripe_account_id = ? AND user_id = ?`
	_, err := db.Exec(query,
		account.Name, account.Mask, account.Type, account.Subtype, account.Balance, account.AvailableBalance,
		account.Currency, account.Status, account.IsHidden, account.StripeAccountID, account.UserID,
	)
	return err
}

// GetStripeFCAccountByRowID retrieves a Stripe FC account by its row_id.
func GetStripeFCAccountByRowID(db *sql.DB, rowID string) (*models.StripeFCAccount, error) {
	query := `SELECT id, row_id, user_id, stripe_account_id, stripe_item_row_id, name, mask, type, subtype, balance, available_balance, currency, status, is_hidden, created_at
              FROM stripe_fc_account WHERE row_id = ?`
	var account models.StripeFCAccount
	err := db.QueryRow(query, rowID).Scan(
		&account.ID, &account.RowID, &account.UserID, &account.StripeAccountID, &account.StripeItemRowID, &account.Name,
		&account.Mask, &account.Type, &account.Subtype, &account.Balance, &account.AvailableBalance,
		&account.Currency, &account.Status, &account.IsHidden, &account.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetStripeFCAccountByStripeID retrieves a Stripe FC account by Stripe account ID.
func GetStripeFCAccountByStripeID(db *sql.DB, stripeAccountID string) (*models.StripeFCAccount, error) {
	query := `SELECT id, row_id, user_id, stripe_account_id, stripe_item_row_id, name, mask, type, subtype, balance, available_balance, currency, status, is_hidden, created_at
              FROM stripe_fc_account WHERE stripe_account_id = ?`
	var account models.StripeFCAccount
	err := db.QueryRow(query, stripeAccountID).Scan(
		&account.ID, &account.RowID, &account.UserID, &account.StripeAccountID, &account.StripeItemRowID, &account.Name,
		&account.Mask, &account.Type, &account.Subtype, &account.Balance, &account.AvailableBalance,
		&account.Currency, &account.Status, &account.IsHidden, &account.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetStripeFCAccountsByUserID returns visible Stripe FC accounts for dashboard and queries.
func GetStripeFCAccountsByUserID(db *sql.DB, userID int64) ([]models.StripeFCAccount, error) {
	query := `SELECT id, row_id, user_id, stripe_account_id, stripe_item_row_id, name, mask, type, subtype, balance, available_balance, currency, status, is_hidden, created_at
              FROM stripe_fc_account WHERE user_id = ? ORDER BY name ASC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.StripeFCAccount
	for rows.Next() {
		var account models.StripeFCAccount
		if err := rows.Scan(
			&account.ID, &account.RowID, &account.UserID, &account.StripeAccountID, &account.StripeItemRowID, &account.Name,
			&account.Mask, &account.Type, &account.Subtype, &account.Balance, &account.AvailableBalance,
			&account.Currency, &account.Status, &account.IsHidden, &account.CreatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// GetStripeFCAccountsByItemRowID returns accounts for a Stripe FC item.
func GetStripeFCAccountsByItemRowID(db *sql.DB, rowID string) ([]models.StripeFCAccount, error) {
	query := `SELECT id, row_id, user_id, stripe_account_id, stripe_item_row_id, name, mask, type, subtype, balance, available_balance, currency, status, is_hidden, created_at
              FROM stripe_fc_account WHERE stripe_item_row_id = ?`
	rows, err := db.Query(query, rowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.StripeFCAccount
	for rows.Next() {
		var account models.StripeFCAccount
		if err := rows.Scan(
			&account.ID, &account.RowID, &account.UserID, &account.StripeAccountID, &account.StripeItemRowID, &account.Name,
			&account.Mask, &account.Type, &account.Subtype, &account.Balance, &account.AvailableBalance,
			&account.Currency, &account.Status, &account.IsHidden, &account.CreatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// ToggleStripeFCAccountVisibility flips the hidden flag for a Stripe FC account.
func ToggleStripeFCAccountVisibility(db *sql.DB, stripeAccountID string, userID int64) (bool, error) {
	query := `UPDATE stripe_fc_account SET is_hidden = NOT is_hidden WHERE row_id = ? AND user_id = ?`
	result, err := db.Exec(query, stripeAccountID, userID)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return false, sql.ErrNoRows
	}
	var isHidden bool
	err = db.QueryRow(`SELECT is_hidden FROM stripe_fc_account WHERE row_id = ? AND user_id = ?`, stripeAccountID, userID).Scan(&isHidden)
	return isHidden, err
}

// DeleteStripeFCAccount removes a Stripe FC account record.
func DeleteStripeFCAccount(db *sql.DB, stripeAccountID string, userID int64) error {
	query := `DELETE FROM stripe_fc_account WHERE row_id = ? AND user_id = ?`
	_, err := db.Exec(query, stripeAccountID, userID)
	return err
}

// UpdateUserStripeCustomerID stores the Stripe customer ID on a user.
func UpdateUserStripeCustomerID(db *sql.DB, userID int64, customerID string) error {
	query := `UPDATE users SET stripe_customer_id = ? WHERE id = ?`
	_, err := db.Exec(query, customerID, userID)
	return err
}
