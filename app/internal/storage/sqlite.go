package storage

import (
    "FinancialTracker/internal/models"
    "FinancialTracker/internal/storage/queries"
    "database/sql"
)

// Storage handles the database connection lifecycle
type Storage struct {
    db *sql.DB
}

// NewSqliteStorage initializes a new Storage instance
func NewSqliteStorage(db *sql.DB) *Storage {
    return &Storage{
        db: db,
    }
}

// GetUserByEmail retrieves a user by their email address
func (s *Storage) GetUserByEmail(email string) (*models.User, error) {
    return queries.GetUserByEmail(s.db, email)
}

// GetUserByID retrieves a user by their ID
func (s *Storage) GetUserByID(id int64) (*models.User, error) {
    return queries.GetUserByID(s.db, id)
}

// GetUserBySSO retrieves a user by their SSO provider and ID
func (s *Storage) GetUserBySSO(provider, ssoID string) (*models.User, error) {
    return queries.GetUserBySSO(s.db, provider, ssoID)
}

// CreateUser creates a new user in the database and sets the userID in the struct
func (s *Storage) CreateUser(user *models.User) error {
    return queries.CreateUser(s.db, user)
}

// LinkSSO links an SSO provider to a user
func (s *Storage) LinkSSO(userID int64, provider, ssoID string) error {
    return queries.LinkSSO(s.db, userID, provider, ssoID)
}

// UnlinkSSO removes an SSO provider from a user
func (s *Storage) UnlinkSSO(userID int64, provider string) error {
    return queries.UnlinkSSO(s.db, userID, provider)
}

// UpdateUserTheme updates the theme preference
func (s *Storage) UpdateUserTheme(userID int64, theme string) error {
    return queries.UpdateUserTheme(s.db, userID, theme)
}

// UpdateUserInfo updates the basic information for a user
func (s *Storage) UpdateUserInfo(userID int64, firstName, lastName, email string) error {
    return queries.UpdateUserInfo(s.db, userID, firstName, lastName, email)
}

// UpdateUserPassword updates the password hash for a user
func (s *Storage) UpdateUserPassword(userID int64, passwordHash string) error {
    return queries.UpdateUserPassword(s.db, userID, passwordHash)
}

// CompleteOnboarding marks the user's onboarding flow as finished.
func (s *Storage) CompleteOnboarding(userID int64) error {
	return queries.CompleteOnboarding(s.db, userID)
}

// DeleteUser deletes a user
func (s *Storage) DeleteUser(userID int64) error {
    return queries.DeleteUser(s.db, userID)
}

// GetUserSubscriptionTier returns the subscription tier for a user.
func (s *Storage) GetUserSubscriptionTier(userID int64) (string, error) {
	return queries.GetUserSubscriptionTier(s.db, userID)
}

// CountActivePlaidItems returns how many non-disconnected Plaid items a user has.
func (s *Storage) CountActivePlaidItems(userID int64) (int, error) {
	return queries.CountActivePlaidItems(s.db, userID)
}

// UpdateUserSubscription updates subscription tier and billing metadata.
func (s *Storage) UpdateUserSubscription(userID int64, tier string, startedAt *int64, stripeCustomerID, stripeSubscriptionID *string) error {
	return queries.UpdateUserSubscription(s.db, userID, tier, startedAt, stripeCustomerID, stripeSubscriptionID)
}

// ApplyStripeSubscription updates tier and subscription ID from Stripe webhooks.
func (s *Storage) ApplyStripeSubscription(userID int64, tier string, startedAt *int64, subscriptionID *string) error {
	return queries.ApplyStripeSubscription(s.db, userID, tier, startedAt, subscriptionID)
}

// GetUserPrivileges returns per-user subscription overrides.
func (s *Storage) GetUserPrivileges(userID int64) (*models.UserPrivileges, error) {
	return queries.GetUserPrivileges(s.db, userID)
}

// UpsertUserPrivileges stores per-user subscription overrides.
func (s *Storage) UpsertUserPrivileges(priv *models.UserPrivileges) error {
	return queries.UpsertUserPrivileges(s.db, priv)
}

// GetUserByStripeCustomerID finds a user by Stripe customer ID.
func (s *Storage) GetUserByStripeCustomerID(customerID string) (*models.User, error) {
	return queries.GetUserByStripeCustomerID(s.db, customerID)
}

// GetPlaidAPICallCount returns Plaid API calls used in the given billing period.
func (s *Storage) GetPlaidAPICallCount(userID int64, periodStart int64) (int, error) {
	return queries.GetPlaidAPICallCount(s.db, userID, periodStart)
}

// IncrementPlaidAPICallCount records one Plaid API call for the user in the given period.
func (s *Storage) IncrementPlaidAPICallCount(userID int64, periodStart int64) error {
	return queries.IncrementPlaidAPICallCount(s.db, userID, periodStart)
}

// CreatePlaidAccount creates a new plaid account
func (s *Storage) CreatePlaidAccount(account *models.Account) error {
    return queries.CreatePlaidAccount(s.db, account)
}

// UpdatePlaidAccount updates an existing plaid account
func (s *Storage) UpdatePlaidAccount(account *models.Account) error {
    return queries.UpdatePlaidAccount(s.db, account)
}

// CreatePlaidItem creates a new plaid item
func (s *Storage) CreatePlaidItem(item *models.PlaidItem) error {
    return queries.CreatePlaidItem(s.db, item)
}

// UpdatePlaidItem updates an existing plaid item
func (s *Storage) UpdatePlaidItem(item *models.PlaidItem) error {
    return queries.UpdatePlaidItem(s.db, item)
}

// GetPlaidItemByItemID retrieves a plaid item by its Plaid item ID
func (s *Storage) GetPlaidItemByItemID(itemID string) (*models.PlaidItem, error) {
    return queries.GetPlaidItemByItemID(s.db, itemID)
}

// GetPlaidItemByRowID retrieves a plaid item by its Row ID and ensures it belongs to the user
func (s *Storage) GetPlaidItemByRowID(rowID string, userID int64) (*models.PlaidItem, error) {
    return queries.GetPlaidItemByRowID(s.db, rowID, userID)
}

// GetPlaidItemsByUserID retrieves all plaid items for a specific user
func (s *Storage) GetPlaidItemsByUserID(userID int64) ([]models.PlaidItem, error) {
    return queries.GetPlaidItemsByUserID(s.db, userID)
}

// UpdatePlaidItemCursor updates the sync cursor for a plaid item
func (s *Storage) UpdatePlaidItemCursor(itemID string, cursor string) error {
    return queries.UpdatePlaidItemCursor(s.db, itemID, cursor)
}

// UpdatePlaidItemStatus updates the connection status and Plaid error code for an item
func (s *Storage) UpdatePlaidItemStatus(itemID string, status string, errorCode string) error {
    return queries.UpdatePlaidItemStatus(s.db, itemID, status, errorCode)
}

// UpdatePlaidItemLastSynced updates the last successful sync timestamp for an item
func (s *Storage) UpdatePlaidItemLastSynced(itemID string, syncedAt int64) error {
    return queries.UpdatePlaidItemLastSynced(s.db, itemID, syncedAt)
}

// MarkPlaidAccountsDisconnectedByItemID marks all accounts under an item as disconnected
func (s *Storage) MarkPlaidAccountsDisconnectedByItemID(itemID string) error {
    return queries.MarkPlaidAccountsDisconnectedByItemID(s.db, itemID)
}

// DeletePlaidItem deletes a plaid item by its Row ID and User ID for security
func (s *Storage) DeletePlaidItem(rowID string, userID int64) error {
	return queries.DeletePlaidItem(s.db, rowID, userID)
}

// DeletePlaidAccount deletes a specific plaid account and its transactions
func (s *Storage) DeletePlaidAccount(accountID string, userID int64) error {
	return queries.DeletePlaidAccount(s.db, accountID, userID)
}

// GetPlaidAccountsByUserID retrieves all plaid accounts for a specific user
func (s *Storage) GetPlaidAccountsByUserID(userID int64) ([]models.Account, error) {
    return queries.GetPlaidAccountsByUserID(s.db, userID)
}

// GetPlaidAccountsByItemID retrieves all plaid accounts for a specific plaid item ID
func (s *Storage) GetPlaidAccountsByItemID(itemID string) ([]models.Account, error) {
	return queries.GetPlaidAccountsByItemID(s.db, itemID)
}

// UpdatePlaidAccountStatus updates the status of a specific plaid account
func (s *Storage) UpdatePlaidAccountStatus(accountID string, status string) error {
	return queries.UpdatePlaidAccountStatus(s.db, accountID, status)
}

// ToggleAccountVisibility flips the is_hidden status of a plaid account and returns the new state
func (s *Storage) ToggleAccountVisibility(accountID string, userID int64) (bool, error) {
    return queries.ToggleAccountVisibility(s.db, accountID, userID)
}

// CreateTransaction creates a new transaction
func (s *Storage) CreateTransaction(t *models.Transaction) error {
    return queries.CreateTransaction(s.db, t)
}

// GetTransactionByPlaidID retrieves a transaction by its Plaid transaction ID
func (s *Storage) GetTransactionByPlaidID(plaidTransactionID string) (*models.Transaction, error) {
    return queries.GetTransactionByPlaidID(s.db, plaidTransactionID)
}

// UpdateTransaction updates an existing transaction record
func (s *Storage) UpdateTransaction(t *models.Transaction) error {
    return queries.UpdateTransaction(s.db, t)
}

// DeleteTransactionByPlaidID deletes a transaction by its Plaid transaction ID
func (s *Storage) DeleteTransactionByPlaidID(plaidTransactionID string) error {
    return queries.DeleteTransactionByPlaidID(s.db, plaidTransactionID)
}

// GetAccountByRowID retrieves an account by its row_id
func (s *Storage) GetAccountByRowID(rowID string) (*models.Account, error) {
    return queries.GetAccountByRowID(s.db, rowID)
}

// GetAccountByPlaidAccountID retrieves an account by its plaid_account_id
func (s *Storage) GetAccountByPlaidAccountID(plaidAccountID string) (*models.Account, error) {
    return queries.GetAccountByPlaidAccountID(s.db, plaidAccountID)
}

// GetTransactions retrieves transactions for a user with filtering and pagination.
func (s *Storage) GetTransactions(userID int64, provider string, f models.TransactionFilters) ([]models.Transaction, int, error) {
	return queries.GetTransactions(s.db, userID, provider, f)
}

// GetAllTagsByUserID retrieves all tags belonging to a user
func (s *Storage) GetAllTagsByUserID(userID int64) ([]models.Tag, error) {
	return queries.GetAllTagsByUserID(s.db, userID)
}

// GetTagByUserIDAndName retrieves a tag by its name for a specific user
func (s *Storage) GetTagByUserIDAndName(userID int64, name string) (*models.Tag, error) {
	return queries.GetTagByUserIDAndName(s.db, userID, name)
}

// GetCategoriesByUserID retrieves all categories for a user
func (s *Storage) GetCategoriesByUserID(userID int64) ([]models.Category, error) {
	return queries.GetCategoriesByUserID(s.db, userID)
}

// AddTagToTransaction associates a tag with a transaction
func (s *Storage) AddTagToTransaction(userID int64, transactionID, tagID int64) error {
	return queries.AddTagToTransaction(s.db, userID, transactionID, tagID)
}

// RemoveTagFromTransaction removes a tag association from a transaction
func (s *Storage) RemoveTagFromTransaction(userID int64, transactionID, tagID int64) error {
	return queries.RemoveTagFromTransaction(s.db, userID, transactionID, tagID)
}

// DeleteTag deletes a tag
func (s *Storage) DeleteTag(userID int64, tagID int64) error {
	return queries.DeleteTag(s.db, userID, tagID)
}

// UpdateTag updates a tag name and color
func (s *Storage) UpdateTag(userID int64, tagID int64, name string, color string) error {
	return queries.UpdateTag(s.db, userID, tagID, name, color)
}

// ApplyTagFiltersToPastTransactions applies all filters of a tag to a user's transactions
func (s *Storage) ApplyTagFiltersToPastTransactions(userID int64, tagID int64) error {
	return queries.ApplyTagFiltersToPastTransactions(s.db, userID, tagID)
}

func (s *Storage) CreateCategory(userID int64, name string) (int64, error) {
	return queries.CreateCategory(s.db, userID, name)
}

// CreateTag creates a new tag under a category
func (s *Storage) CreateTag(userID int64, categoryID int64, name string, color string) (int64, error) {
	return queries.CreateTag(s.db, userID, categoryID, name, color)
}

// CreateTagFilter creates a new auto-tagging rule
func (s *Storage) CreateTagFilter(userID int64, tagID int64, pattern string, filterType string) error {
	return queries.CreateTagFilter(s.db, userID, tagID, pattern, filterType)
}

// GetTagFiltersByUserID retrieves all tag filters for a user
func (s *Storage) GetTagFiltersByUserID(userID int64) ([]models.TagFilter, error) {
	return queries.GetTagFiltersByUserID(s.db, userID)
}

// GetCategoryByID retrieves a category by its ID
func (s *Storage) GetCategoryByID(categoryID int64) (*models.Category, error) {
	return queries.GetCategoryByID(s.db, categoryID)
}

// UpdateCategory updates a category name
func (s *Storage) UpdateCategory(userID int64, categoryID int64, name string) error {
	return queries.UpdateCategory(s.db, userID, categoryID, name)
}

// GetOrCreateMiscCategory ensures a "Misc" category exists for a user and returns its ID
func (s *Storage) GetOrCreateMiscCategory(userID int64) (int64, error) {
	return queries.GetOrCreateMiscCategory(s.db, userID)
}

// DeleteCategory deletes a category and handles its tags
func (s *Storage) DeleteCategory(userID int64, categoryID int64, moveTagsToCategoryID int64) error {
	return queries.DeleteCategory(s.db, userID, categoryID, moveTagsToCategoryID)
}

// MoveTagToCategory moves a tag to a different category
func (s *Storage) MoveTagToCategory(userID int64, tagID int64, categoryID int64) error {
	return queries.MoveTagToCategory(s.db, userID, tagID, categoryID)
}

// GetTagFiltersByTagID retrieves all filters for a specific tag
func (s *Storage) GetTagFiltersByTagID(userID int64, tagID int64) ([]models.TagFilter, error) {
	return queries.GetTagFiltersByTagID(s.db, userID, tagID)
}

// DeleteTagFiltersByTagID deletes all filters for a tag
func (s *Storage) DeleteTagFiltersByTagID(userID int64, tagID int64) error {
	return queries.DeleteTagFiltersByTagID(s.db, userID, tagID)
}

// BatchCreateTagFilters creates multiple tag filters
func (s *Storage) BatchCreateTagFilters(userID int64, tagID int64, filters []models.TagFilter) error {
	return queries.BatchCreateTagFilters(s.db, userID, tagID, filters)
}

// CreateSession creates a new session
func (s *Storage) CreateSession(session *models.Session) error {
    return queries.CreateSession(s.db, session)
}

// GetSession retrieves a session by ID
func (s *Storage) GetSession(id string) (*models.Session, error) {
    return queries.GetSession(s.db, id)
}

// UpdateSessionReauth updates the re-authentication timestamp
func (s *Storage) UpdateSessionReauth(id string, timestamp int64) error {
    return queries.UpdateSessionReauth(s.db, id, timestamp)
}

// DeleteSession deletes a session by ID
func (s *Storage) DeleteSession(id string) error {
    return queries.DeleteSession(s.db, id)
}

// DeleteUserSessions removes all sessions for a user.
func (s *Storage) DeleteUserSessions(userID int64) error {
	return queries.DeleteUserSessions(s.db, userID)
}

// InvalidatePasswordResetCodes removes unused reset codes for a user.
func (s *Storage) InvalidatePasswordResetCodes(userID int64) error {
	return queries.InvalidatePasswordResetCodes(s.db, userID)
}

// CreatePasswordResetCode stores a hashed reset code for a user.
func (s *Storage) CreatePasswordResetCode(userID int64, codeHash string, expiresAt int64) error {
	return queries.CreatePasswordResetCode(s.db, userID, codeHash, expiresAt)
}

// GetActivePasswordResetCode returns the latest unused, unexpired reset code for a user.
func (s *Storage) GetActivePasswordResetCode(userID int64) (*models.PasswordResetCode, error) {
	return queries.GetActivePasswordResetCode(s.db, userID)
}

// IncrementPasswordResetAttempts bumps the failed-attempt counter for a reset code.
func (s *Storage) IncrementPasswordResetAttempts(id int64) error {
	return queries.IncrementPasswordResetAttempts(s.db, id)
}

// MarkPasswordResetCodeUsed marks a reset code as consumed.
func (s *Storage) MarkPasswordResetCodeUsed(id int64) error {
	return queries.MarkPasswordResetCodeUsed(s.db, id)
}

// GetLatestPasswordResetCreatedAt returns when the most recent reset code was created for a user.
func (s *Storage) GetLatestPasswordResetCreatedAt(userID int64) (int64, error) {
	return queries.GetLatestPasswordResetCreatedAt(s.db, userID)
}

// CountPasswordResetCodesSince returns how many reset codes were issued since a unix timestamp.
func (s *Storage) CountPasswordResetCodesSince(userID int64, since int64) (int, error) {
	return queries.CountPasswordResetCodesSince(s.db, userID, since)
}

// GetDashboardLayout retrieves a user's dashboard widget layout.
func (s *Storage) GetDashboardLayout(userID int64) (*models.DashboardLayout, error) {
	return queries.GetDashboardLayout(s.db, userID)
}

// UpsertDashboardLayout saves a user's dashboard widget layout.
func (s *Storage) UpsertDashboardLayout(userID int64, layout *models.DashboardLayout) error {
	return queries.UpsertDashboardLayout(s.db, userID, layout)
}

// GetMonthlySpending returns monthly expense totals for dashboard charts.
func (s *Storage) GetMonthlySpending(userID int64, provider string, months int) ([]models.MonthlySpend, error) {
	return queries.GetMonthlySpending(s.db, userID, provider, months)
}

// GetMonthCashflow returns spend and income totals for a date range.
func (s *Storage) GetMonthCashflow(userID int64, provider string, monthStart, monthEnd int64) (models.MonthCashflow, error) {
	return queries.GetMonthCashflow(s.db, userID, provider, monthStart, monthEnd)
}

// GetSpendingByTag returns spending breakdown by tag for a date range.
func (s *Storage) GetSpendingByTag(userID int64, provider string, monthStart, monthEnd int64) ([]models.TagBreakdown, error) {
	return queries.GetSpendingByTag(s.db, userID, provider, monthStart, monthEnd)
}

// GetIncomeByTag returns income breakdown by tag for a date range.
func (s *Storage) GetIncomeByTag(userID int64, provider string, monthStart, monthEnd int64) ([]models.TagBreakdown, error) {
	return queries.GetIncomeByTag(s.db, userID, provider, monthStart, monthEnd)
}

// UpdateAccountMonthlyPayment sets the monthly payment amount for a loan account.
func (s *Storage) UpdateAccountMonthlyPayment(plaidAccountID string, monthlyPayment float64) error {
	return queries.UpdateAccountMonthlyPayment(s.db, plaidAccountID, monthlyPayment)
}

func (s *Storage) CreateStripeFCItem(item *models.StripeFCItem) error {
	return queries.CreateStripeFCItem(s.db, item)
}

func (s *Storage) GetStripeFCItemByRowID(rowID string, userID int64) (*models.StripeFCItem, error) {
	return queries.GetStripeFCItemByRowID(s.db, rowID, userID)
}

func (s *Storage) GetStripeFCItemByInstitution(userID int64, institutionName string) (*models.StripeFCItem, error) {
	return queries.GetStripeFCItemByInstitution(s.db, userID, institutionName)
}

func (s *Storage) GetStripeFCItemsByUserID(userID int64) ([]models.StripeFCItem, error) {
	return queries.GetStripeFCItemsByUserID(s.db, userID)
}

func (s *Storage) UpdateStripeFCItemStatus(rowID, status, errorCode string) error {
	return queries.UpdateStripeFCItemStatus(s.db, rowID, status, errorCode)
}

func (s *Storage) UpdateStripeFCItemLastSynced(rowID string, syncedAt int64) error {
	return queries.UpdateStripeFCItemLastSynced(s.db, rowID, syncedAt)
}

func (s *Storage) UpdateStripeFCItemTransactionRefresh(rowID, refreshID string) error {
	return queries.UpdateStripeFCItemTransactionRefresh(s.db, rowID, refreshID)
}

func (s *Storage) DeleteStripeFCItem(rowID string, userID int64) error {
	return queries.DeleteStripeFCItem(s.db, rowID, userID)
}

func (s *Storage) CountActiveStripeFCItems(userID int64) (int, error) {
	return queries.CountActiveStripeFCItems(s.db, userID)
}

func (s *Storage) CreateStripeFCAccount(account *models.StripeFCAccount) error {
	return queries.CreateStripeFCAccount(s.db, account)
}

func (s *Storage) UpdateStripeFCAccount(account *models.StripeFCAccount) error {
	return queries.UpdateStripeFCAccount(s.db, account)
}

// GetStripeFCAccountByRowID retrieves a Stripe FC account by its row_id.
func (s *Storage) GetStripeFCAccountByRowID(rowID string) (*models.StripeFCAccount, error) {
	return queries.GetStripeFCAccountByRowID(s.db, rowID)
}

// GetStripeFCAccountByStripeID retrieves a Stripe FC account by Stripe account ID.
func (s *Storage) GetStripeFCAccountByStripeID(stripeAccountID string) (*models.StripeFCAccount, error) {
	return queries.GetStripeFCAccountByStripeID(s.db, stripeAccountID)
}

func (s *Storage) GetStripeFCAccountsByUserID(userID int64) ([]models.StripeFCAccount, error) {
	return queries.GetStripeFCAccountsByUserID(s.db, userID)
}

func (s *Storage) GetStripeFCAccountsByItemRowID(rowID string) ([]models.StripeFCAccount, error) {
	return queries.GetStripeFCAccountsByItemRowID(s.db, rowID)
}

func (s *Storage) ToggleStripeFCAccountVisibility(stripeAccountID string, userID int64) (bool, error) {
	return queries.ToggleStripeFCAccountVisibility(s.db, stripeAccountID, userID)
}

func (s *Storage) DeleteStripeFCAccount(stripeAccountID string, userID int64) error {
	return queries.DeleteStripeFCAccount(s.db, stripeAccountID, userID)
}

func (s *Storage) UpdateUserStripeCustomerID(userID int64, customerID string) error {
	return queries.UpdateUserStripeCustomerID(s.db, userID, customerID)
}

func (s *Storage) GetStripeAPICallCount(userID int64, periodStart int64) (int, error) {
	return queries.GetStripeAPICallCount(s.db, userID, periodStart)
}

func (s *Storage) IncrementStripeAPICallCount(userID int64, periodStart int64) error {
	return queries.IncrementStripeAPICallCount(s.db, userID, periodStart)
}

// GetLinkedAccountsByUserID returns dashboard accounts for the active provider.
func (s *Storage) GetLinkedAccountsByUserID(userID int64, provider string) ([]models.Account, error) {
	if provider == "stripe" {
		stripeAccounts, err := queries.GetStripeFCAccountsByUserID(s.db, userID)
		if err != nil {
			return nil, err
		}
		accounts := make([]models.Account, 0, len(stripeAccounts))
		for _, acc := range stripeAccounts {
			accounts = append(accounts, models.Account{
				ID:               acc.ID,
				UserID:           acc.UserID,
				PlaidAccountID:   acc.StripeAccountID,
				PlaidItemID:      acc.StripeItemRowID,
				Name:             acc.Name,
				Mask:             acc.Mask,
				Type:             acc.Type,
				Subtype:          acc.Subtype,
				Balance:          acc.Balance,
				AvailableBalance: acc.AvailableBalance,
				Currency:         acc.Currency,
				Status:           acc.Status,
				IsHidden:         acc.IsHidden,
				CreatedAt:        acc.CreatedAt,
			})
		}
		return accounts, nil
	}
	return queries.GetPlaidAccountsByUserID(s.db, userID)
}
