package stripefc

import (
	"context"
	"errors"
	"os"

	"FinancialTracker/internal/models"
	"FinancialTracker/internal/models/external"
	connectionService "FinancialTracker/internal/services/connections"

	"github.com/labstack/echo/v5"
	"github.com/stripe/stripe-go/v86"
)

// GetManagementData fetches all Stripe FC connections and accounts for a user.
func (s *Service) GetManagementData(userID int64) (*external.ConnectionsPayload, error) {
	items, err := s.store.GetStripeFCItemsByUserID(userID)
	if err != nil {
		return nil, err
	}

	accounts, err := s.store.GetStripeFCAccountsByUserID(userID)
	if err != nil {
		return nil, err
	}

	usage, err := s.GetUsage(userID)
	if err != nil {
		return nil, err
	}

	grouped := make([]models.StripeFCItemWithAccounts, len(items))
	for i, item := range items {
		grouped[i] = models.StripeFCItemWithAccounts{StripeFCItem: item}
		for _, acc := range accounts {
			if acc.StripeItemRowID == item.RowID {
				grouped[i].Accounts = append(grouped[i].Accounts, acc)
			}
		}
	}

	return connectionService.BuildPayloadFromStripe(grouped, usage), nil
}

// CreateSession creates a Stripe Financial Connections session for linking accounts.
func (s *Service) CreateSession(c *echo.Context, userID int64) (*external.CreateSessionResponse, error) {
	if err := s.ensureItemLimitAvailable(userID); err != nil {
		return nil, err
	}
	return s.createSession(c.Request().Context(), userID, "")
}

// CreateUpdateSession creates a session to relink accounts for an existing institution.
func (s *Service) CreateUpdateSession(c *echo.Context, userID int64, rowID string) (*external.CreateSessionResponse, error) {
	item, err := s.store.GetStripeFCItemByRowID(rowID, userID)
	if err != nil {
		return nil, errors.New("connection not found")
	}
	if item.Status == ItemStatusDisconnected {
		return nil, errors.New("this bank connection is no longer available; please disconnect and link again")
	}
	return s.createSession(c.Request().Context(), userID, rowID)
}

func (s *Service) createSession(ctx context.Context, userID int64, rowID string) (*external.CreateSessionResponse, error) {
	customerID, err := s.ensureStripeCustomer(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.reserveStripeAPICall(userID); err != nil {
		return nil, err
	}

	params := &stripe.FinancialConnectionsSessionCreateParams{
		AccountHolder: &stripe.FinancialConnectionsSessionCreateAccountHolderParams{
			Type:     stripe.String(string(stripe.FinancialConnectionsSessionAccountHolderTypeCustomer)),
			Customer: stripe.String(customerID),
		},
		Permissions: []*string{
			stripe.String(string(stripe.FinancialConnectionsSessionPermissionBalances)),
			stripe.String(string(stripe.FinancialConnectionsSessionPermissionTransactions)),
		},
		Filters: &stripe.FinancialConnectionsSessionCreateFiltersParams{
			Countries: []*string{stripe.String("US")},
		},
	}
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		params.ReturnURL = stripe.String(frontendURL + "/settings")
	}
	_ = rowID

	session, err := s.client.V1FinancialConnectionsSessions.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	return &external.CreateSessionResponse{ClientSecret: session.ClientSecret}, nil
}

// CompleteConnection stores accounts linked through a Financial Connections session.
func (s *Service) CompleteConnection(c *echo.Context, userID int64, req *external.CompleteConnectionRequest) error {
	if req == nil || req.SessionID == "" {
		return errors.New("session id is required")
	}

	ctx := c.Request().Context()
	if err := s.reserveStripeAPICall(userID); err != nil {
		return err
	}

	session, err := s.client.V1FinancialConnectionsSessions.Retrieve(ctx, req.SessionID, &stripe.FinancialConnectionsSessionRetrieveParams{})
	if err != nil {
		return err
	}

	listParams := &stripe.FinancialConnectionsAccountListParams{
		Session: stripe.String(session.ID),
	}
	iter := s.client.V1FinancialConnectionsAccounts.List(ctx, listParams)
	for account, err := range iter.All(ctx) {
		if err != nil {
			return err
		}
		if err := s.upsertAccountFromStripe(ctx, userID, account); err != nil {
			return err
		}
	}

	return s.syncItems(ctx, userID, false)
}
