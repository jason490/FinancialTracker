package transactions

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/storage"
	"fmt"
	"strconv"
)

type TransactionService struct {
	store *storage.Storage
}

func NewTransactionService(store *storage.Storage) *TransactionService {
	return &TransactionService{
		store: store,
	}
}

// GetTransactionsPage prepares data for the transactions page including filtered transactions, tags, and categories
func (s *TransactionService) GetTransactionsPage(userID int64, filters models.TransactionFilters, page int, pageSize int) (models.TransactionPageData, error) {
	filters.Limit = pageSize
	filters.Offset = (page - 1) * pageSize

	transactions, totalCount, err := s.store.GetTransactions(userID, filters)
	if err != nil {
		return models.TransactionPageData{}, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	allTags, err := s.store.GetAllTagsByUserID(userID)
	if err != nil {
		// Log error but don't fail the whole request
		allTags = []models.Tag{}
	}

	categories, err := s.store.GetCategoriesByUserID(userID)
	if err != nil {
		// Log error but don't fail the whole request
		categories = []models.Category{}
	}

	return models.TransactionPageData{
		Transactions: transactions,
		TotalCount:   totalCount,
		CurrentPage:  page,
		PageSize:     pageSize,
		Filters:      filters,
		AllTags:      allTags,
		Categories:   categories,
	}, nil
}

// BulkAddTag adds a tag to multiple transactions for a user
func (s *TransactionService) BulkAddTag(userID int64, transactionIDs []string, tagID int64) error {
	if tagID == 0 || len(transactionIDs) == 0 {
		return fmt.Errorf("tag and transactions are required")
	}

	for _, idStr := range transactionIDs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}
		if err := s.store.AddTagToTransaction(userID, id, tagID); err != nil {
			return fmt.Errorf("failed to add tag %d to transaction %d: %w", tagID, id, err)
		}
	}

	return nil
}

// BulkRemoveTag removes a tag from multiple transactions for a user
func (s *TransactionService) BulkRemoveTag(userID int64, transactionIDs []string, tagID int64) error {
	if tagID == 0 || len(transactionIDs) == 0 {
		return fmt.Errorf("tag and transactions are required")
	}

	for _, idStr := range transactionIDs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}
		if err := s.store.RemoveTagFromTransaction(userID, id, tagID); err != nil {
			return fmt.Errorf("failed to remove tag %d from transaction %d: %w", tagID, id, err)
		}
	}

	return nil
}

// GetRecentTransactions fetches the most recent transactions for a user
func (s *TransactionService) GetRecentTransactions(userID int64, limit int) ([]models.Transaction, error) {
	transactions, _, err := s.store.GetTransactions(userID, models.TransactionFilters{
		Limit:  limit,
		Offset: 0,
	})
	return transactions, err
}
