package stripefc

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"FinancialTracker/internal/models"

	"github.com/labstack/gommon/log"
	"github.com/stripe/stripe-go/v86"
)

// upsertAccountFromStripe stores or updates a Stripe FC account grouped by institution.
func (s *Service) upsertAccountFromStripe(ctx context.Context, userID int64, account *stripe.FinancialConnectionsAccount) error {
	if account == nil || account.ID == "" {
		return nil
	}

	full, err := s.fetchAccount(ctx, userID, account.ID)
	if err != nil {
		return err
	}

	institutionName := full.InstitutionName
	if institutionName == "" {
		institutionName = "Unknown Institution"
	}

	item, err := s.store.GetStripeFCItemByInstitution(userID, institutionName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		item = &models.StripeFCItem{
			UserID:          userID,
			InstitutionName: institutionName,
			Status:          mapStripeItemStatus(full.Status),
		}
		if err := s.store.CreateStripeFCItem(item); err != nil {
			return err
		}
	}

	accType, subtype := mapStripeAccountType(full.Category, full.Subcategory)
	balance, available := extractBalances(full.Balance)
	currency := extractCurrency(full.Balance, "USD")

	stored := &models.StripeFCAccount{
		UserID:           userID,
		StripeAccountID:  full.ID,
		StripeItemRowID:  item.RowID,
		Name:             full.DisplayName,
		Mask:             full.Last4,
		Type:             accType,
		Subtype:          subtype,
		Balance:          balance,
		AvailableBalance: available,
		Currency:         currency,
		Status:           mapStripeAccountStatus(full.Status),
	}

	existing, err := s.store.GetStripeFCAccountByStripeID(full.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if err := s.store.CreateStripeFCAccount(stored); err != nil {
			return err
		}
	} else {
		stored.IsHidden = existing.IsHidden
		if err := s.store.UpdateStripeFCAccount(stored); err != nil {
			return err
		}
	}

	if err := s.store.UpdateStripeFCItemStatus(item.RowID, mapStripeItemStatus(full.Status), ""); err != nil {
		log.Errorf("Failed to update Stripe FC item status for %s: %v", item.RowID, err)
	}

	return s.ensureTransactionSubscription(ctx, userID, full.ID)
}

func (s *Service) fetchAccount(ctx context.Context, userID int64, accountID string) (*stripe.FinancialConnectionsAccount, error) {
	if err := s.reserveStripeAPICall(userID); err != nil {
		return nil, err
	}
	return s.client.V1FinancialConnectionsAccounts.GetByID(ctx, accountID, &stripe.FinancialConnectionsAccountRetrieveParams{})
}

func (s *Service) syncStripeAccounts(ctx context.Context, userID int64, itemRowID string) error {
	accounts, err := s.store.GetStripeFCAccountsByItemRowID(itemRowID)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if account.Status == ItemStatusDisconnected {
			continue
		}
		full, err := s.fetchAccount(ctx, userID, account.StripeAccountID)
		if err != nil {
			log.Errorf("Failed to refresh Stripe account %s: %v", account.StripeAccountID, err)
			continue
		}

		balance, available := extractBalances(full.Balance)
		currency := extractCurrency(full.Balance, account.Currency)

		updated := account
		updated.Name = full.DisplayName
		updated.Mask = full.Last4
		updated.Balance = balance
		updated.AvailableBalance = available
		updated.Currency = currency
		updated.Status = mapStripeAccountStatus(full.Status)

		if err := s.store.UpdateStripeFCAccount(&updated); err != nil {
			log.Errorf("Failed to update Stripe account %s: %v", account.StripeAccountID, err)
		}
	}

	return nil
}

func (s *Service) ensureTransactionSubscription(ctx context.Context, userID int64, accountID string) error {
	if err := s.reserveStripeAPICall(userID); err != nil {
		return err
	}
	_, err := s.client.V1FinancialConnectionsAccounts.Subscribe(ctx, accountID, &stripe.FinancialConnectionsAccountSubscribeParams{
		Features: []*string{stripe.String("transactions")},
	})
	return err
}

func extractBalances(balance *stripe.FinancialConnectionsAccountBalance) (current, available float64) {
	if balance == nil {
		return 0, 0
	}
	for _, value := range balance.Current {
		current = float64(value) / 100
		break
	}
	if balance.Cash != nil {
		for _, value := range balance.Cash.Available {
			available = float64(value) / 100
			break
		}
	}
	return current, available
}

func extractCurrency(balance *stripe.FinancialConnectionsAccountBalance, fallback string) string {
	if balance == nil {
		return fallback
	}
	for currency := range balance.Current {
		return strings.ToUpper(currency)
	}
	return fallback
}

func mapStripeAccountType(category stripe.FinancialConnectionsAccountCategory, subcategory stripe.FinancialConnectionsAccountSubcategory) (string, string) {
	switch category {
	case stripe.FinancialConnectionsAccountCategoryCredit:
		return "credit", string(subcategory)
	case stripe.FinancialConnectionsAccountCategoryInvestment:
		return "investment", string(subcategory)
	default:
		return "depository", string(subcategory)
	}
}

func mapStripeAccountStatus(status stripe.FinancialConnectionsAccountStatus) string {
	switch status {
	case stripe.FinancialConnectionsAccountStatusActive:
		return ItemStatusActive
	case stripe.FinancialConnectionsAccountStatusInactive, stripe.FinancialConnectionsAccountStatusDisconnected:
		return ItemStatusDisconnected
	default:
		return ItemStatusError
	}
}

func mapStripeItemStatus(status stripe.FinancialConnectionsAccountStatus) string {
	return mapStripeAccountStatus(status)
}

// ToggleAccountVisibility flips whether a Stripe FC account is hidden.
func (s *Service) ToggleAccountVisibility(userID int64, accountID string) (bool, error) {
	if accountID == "" {
		return false, errors.New("invalid account ID")
	}
	return s.store.ToggleStripeFCAccountVisibility(accountID, userID)
}

// RemoveAccount permanently deletes a disconnected Stripe FC account.
func (s *Service) RemoveAccount(userID int64, accountID string) error {
	if accountID == "" {
		return errors.New("invalid account ID")
	}

	account, err := s.store.GetStripeFCAccountByRowID(accountID)
	if err != nil || account.UserID != userID {
		return errors.New("account not found")
	}
	if account.Status != ItemStatusDisconnected {
		return errors.New("cannot delete an active account; please disconnect it first")
	}
	return s.store.DeleteStripeFCAccount(accountID, userID)
}
