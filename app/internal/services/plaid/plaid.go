package plaid

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/services/tags"
	"FinancialTracker/internal/storage"
	"errors"
	"fmt"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"github.com/plaid/plaid-go/v42/plaid"
)

type PlaidService struct {
	client *plaid.APIClient
	store  *storage.Storage
	tagService *tags.TaggingService
}

func NewPlaidService(store *storage.Storage, tagService *tags.TaggingService) *PlaidService {
	clientID := os.Getenv("PLAID_CLIENT_ID")
	env := os.Getenv("PLAID_ENV")

	var secret string
	var environment plaid.Environment
	switch env {
	case "production":
		environment = plaid.Production
		secret = os.Getenv("PLAID_PROD_SECRET")
	case "development":
		environment = plaid.Sandbox
		secret = os.Getenv("PLAID_SANDBOX_SECRET")
	default:
		environment = plaid.Sandbox
		secret = os.Getenv("PLAID_SANDBOX_SECRET")
	}

	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", clientID)
	configuration.AddDefaultHeader("PLAID-SECRET", secret)
	configuration.UseEnvironment(environment)

	client := plaid.NewAPIClient(configuration)

	return &PlaidService{
		client: client,
		store:  store,
		tagService: tagService,
	}
}

// CreateLinkToken generates a new Plaid Link token for the user to initialize the bank connection flow.
func (p *PlaidService) CreateLinkToken(c *echo.Context, userId string) (string, error) {
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: userId,
	}
	request := plaid.NewLinkTokenCreateRequest(
		"Financial Tracker",
		"en",
		[]plaid.CountryCode{plaid.COUNTRYCODE_US},
	)
	request.SetUser(user)
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_TRANSACTIONS})

	resp, _, err := p.client.PlaidApi.LinkTokenCreate(c.Request().Context()).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", err
	}
	return resp.GetLinkToken(), nil
}

// CreateUpdateLinkToken generates a link token for an existing item to fix its connection or modify its shared accounts.
func (p *PlaidService) CreateUpdateLinkToken(c *echo.Context, userID string, accessToken string, itemStatus string) (string, error) {
	if itemStatus == ItemStatusDisconnected {
		return "", errors.New("this bank connection is no longer available; please disconnect and link again")
	}
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: userID,
	}
	request := plaid.NewLinkTokenCreateRequest(
		"Financial Tracker",
		"en",
		[]plaid.CountryCode{plaid.COUNTRYCODE_US},
	)
	request.SetUser(user)
	request.SetAccessToken(accessToken)

	// Enable account selection in update mode to allow adding new accounts or modifying shared ones
	update := plaid.NewLinkTokenCreateRequestUpdate()
	update.SetAccountSelectionEnabled(true)
	request.SetUpdate(*update)

	resp, _, err := p.client.PlaidApi.LinkTokenCreate(c.Request().Context()).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		log.Errorf("Failed to create update link token for user %s: %v", userID, err)
		return "", err
	}
	return resp.GetLinkToken(), nil
}

// ExchangeToken handles the public-to-access token exchange and initializes the connection.
// It also triggers the initial account and transaction sync.
func (p *PlaidService) ExchangeToken(c *echo.Context, userID int64, publicToken string) error {
	ctx := c.Request().Context()
	exchangeRequest := plaid.NewItemPublicTokenExchangeRequest(publicToken)
	exchangeResp, _, err := p.client.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(*exchangeRequest).Execute()
	if err != nil {
		log.Error(err)
		return errors.New("failed to exchange token with bank")
	}

	accessToken := exchangeResp.GetAccessToken()
	itemID := exchangeResp.GetItemId()

	// Check if the item already exists for this user to determine whether to update or create
	existing, err := p.store.GetPlaidItemByItemID(itemID)
	if err == nil && existing != nil {
		if existing.UserID != userID {
			return errors.New("this bank account is already linked to another user")
		}
		// Update existing item (e.g., after re-auth or adding accounts)
		if err := p.UpdatePlaidItem(&ctx, userID, itemID, accessToken); err != nil {
			return err
		}
	} else {
		// Create new item
		if err := p.CreatePlaidItem(&ctx, userID, itemID, accessToken); err != nil {
			return err
		}
	}

	return p.syncItems(&ctx, userID, true)
}

// GetManagementData fetches all connections and their accounts for a user
func (p *PlaidService) GetManagementData(userID int64) ([]models.PlaidItemWithAccounts, error) {
	items, err := p.store.GetPlaidItemsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plaid items: %w", err)
	}

	accounts, err := p.store.GetPlaidAccountsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plaid accounts: %w", err)
	}

	groupedItems := make([]models.PlaidItemWithAccounts, len(items))
	for i, item := range items {
		groupedItems[i] = models.PlaidItemWithAccounts{
			PlaidItem: item,
			Accounts:  []models.Account{},
		}
		for _, acc := range accounts {
			if acc.PlaidItemID == item.PlaidItemID {
				groupedItems[i].Accounts = append(groupedItems[i].Accounts, acc)
			}
		}
	}

	return groupedItems, nil
}

// DeletePlaidAccount removes a specific bank account if it is disconnected
func (p *PlaidService) DeletePlaidAccount(userID int64, accountID string) error {
	if accountID == "" {
		return errors.New("invalid account ID")
	}

	account, err := p.store.GetAccountByPlaidAccountID(accountID)
	if err != nil || account.UserID != userID {
		return errors.New("account not found")
	}

	if account.Status != "disconnected" {
		return errors.New("cannot delete an active account; please disconnect it first")
	}

	return p.store.DeletePlaidAccount(accountID, userID)
}

// ToggleAccountVisibility flips the hidden status of a bank account
func (p *PlaidService) ToggleAccountVisibility(userID int64, accountID string) (bool, error) {
	if accountID == "" {
		return false, errors.New("invalid account ID")
	}
	return p.store.ToggleAccountVisibility(accountID, userID)
}
