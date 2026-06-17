package dashboard

import (
	"FinancialTracker/internal/models"
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

	accounts, err := s.store.GetPlaidAccountsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}

	transactions, _, err := s.store.GetTransactions(userID, models.TransactionFilters{
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

	spending, err := s.store.GetMonthlySpending(userID, 6)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch spending analytics: %w", err)
	}
	ApplySpendingAnalytics(data, spending)

	monthStart, monthEnd := CurrentMonthBounds()
	cashflow, err := s.store.GetMonthCashflow(userID, monthStart, monthEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cashflow analytics: %w", err)
	}
	spendingByTag, err := s.store.GetSpendingByTag(userID, monthStart, monthEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch spending by tag: %w", err)
	}
	incomeByTag, err := s.store.GetIncomeByTag(userID, monthStart, monthEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch income by tag: %w", err)
	}
	ApplyMonthAnalytics(data, cashflow, spendingByTag, incomeByTag)

	data.EditMode = editMode
	return data, nil
}

// SaveLayout validates and persists the dashboard widget configuration.
func (s *DashboardService) SaveLayout(userID int64, rawLayout string) (*models.DashboardData, error) {
	if rawLayout == "" {
		return nil, fmt.Errorf("no layout data provided")
	}

	layout, err := ParseLayoutJSON(rawLayout)
	if err != nil {
		return nil, fmt.Errorf("invalid layout format: %w", err)
	}

	if err := s.store.UpsertDashboardLayout(userID, &layout); err != nil {
		return nil, fmt.Errorf("failed to save dashboard layout: %w", err)
	}

	return s.GetDashboardData(userID, false)
}
