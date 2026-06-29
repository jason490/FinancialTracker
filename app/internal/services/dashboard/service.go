package dashboard

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/services/financial"
	"FinancialTracker/internal/services/tags"
	"FinancialTracker/internal/storage"
	"fmt"
)

type DashboardService struct {
	store      *storage.Storage
	tagService *tags.TaggingService
}

func NewDashboardService(store *storage.Storage, tagService *tags.TaggingService) *DashboardService {
	return &DashboardService{
		store:      store,
		tagService: tagService,
	}
}

// GetDashboardData builds the dashboard view model for a user.
func (s *DashboardService) GetDashboardData(userID int64, editMode bool) (*models.DashboardData, error) {
	s.tagService.SeedDefaults(userID)

	provider := financial.ActiveProvider()
	accounts, err := s.store.GetLinkedAccountsByUserID(userID, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}

	transactions, _, err := s.store.GetTransactions(userID, provider, models.TransactionFilters{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	layout, err := s.store.GetDashboardLayout(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch layout: %w", err)
	}
	if layout == nil {
		defaultLayout := DefaultLayout()
		layout = &defaultLayout
	}

	data := BuildDashboardData(accounts, transactions, *layout)

	spending, err := s.store.GetMonthlySpending(userID, provider, 6)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch spending analytics: %w", err)
	}
	ApplySpendingAnalytics(data, spending)

	monthStart, monthEnd := CurrentMonthBounds()
	cashflow, err := s.store.GetMonthCashflow(userID, provider, monthStart, monthEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cashflow analytics: %w", err)
	}
	spendingByTag, err := s.store.GetSpendingByTag(userID, provider, monthStart, monthEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch spending by tag: %w", err)
	}
	incomeByTag, err := s.store.GetIncomeByTag(userID, provider, monthStart, monthEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch income by tag: %w", err)
	}
	ApplyMonthAnalytics(data, cashflow, spendingByTag, incomeByTag)

	data.EditMode = editMode
	return data, nil
}
