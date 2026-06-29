package transactions

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/services/financial"
	"FinancialTracker/internal/services/subscription"
	"FinancialTracker/internal/storage"
	"archive/zip"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"
)

// ErrExportAPILimitExceeded is returned when the user has no API calls left in
// the current billing period to spend on an export.
var ErrExportAPILimitExceeded = errors.New("monthly API limit reached for your plan")

type TransactionService struct {
	store        *storage.Storage
	subscription *subscription.Service
}

func NewTransactionService(store *storage.Storage, sub *subscription.Service) *TransactionService {
	return &TransactionService{
		store:        store,
		subscription: sub,
	}
}

// GetTransactionsPage prepares data for the transactions page including filtered transactions, tags, and categories
func (s *TransactionService) GetTransactionsPage(userID int64, filters models.TransactionFilters, page int, pageSize int) (models.TransactionPageData, error) {
	filters.Limit = pageSize
	filters.Offset = (page - 1) * pageSize

	transactions, totalCount, err := s.store.GetTransactions(userID, financial.ActiveProvider(), filters)
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

// ExportTransactionsZip streams every visible transaction for a user as a ZIP
// archive containing a single CSV file (csvName). Each export consumes one API
// call from the user's monthly quota (shared with Plaid syncs); when the quota
// is exhausted it returns ErrExportAPILimitExceeded before any bytes are
// written so the caller can surface a 429 to the client.
func (s *TransactionService) ExportTransactionsZip(userID int64, csvName string, w io.Writer) (int, error) {
	if s.subscription != nil {
		if err := s.subscription.ReserveAPICall(userID); err != nil {
			if errors.Is(err, subscription.ErrAPILimitExceeded) {
				return 0, ErrExportAPILimitExceeded
			}
			return 0, fmt.Errorf("failed to reserve api call: %w", err)
		}
	}

	zw := zip.NewWriter(w)
	defer zw.Close()

	entry, err := zw.Create(csvName)
	if err != nil {
		return 0, fmt.Errorf("failed to create zip entry: %w", err)
	}

	rows, err := s.writeTransactionsCSV(userID, entry)
	if err != nil {
		return rows, err
	}

	if err := zw.Close(); err != nil {
		return rows, fmt.Errorf("failed to finalize zip archive: %w", err)
	}
	return rows, nil
}

// writeTransactionsCSV writes the user's full transaction history as CSV to w
// and returns the number of data rows emitted (excluding the header).
func (s *TransactionService) writeTransactionsCSV(userID int64, w io.Writer) (int, error) {
	transactions, _, err := s.store.GetTransactions(userID, financial.ActiveProvider(), models.TransactionFilters{
		SortBy:  "date",
		SortDir: "desc",
		Limit:   math.MaxInt32,
		Offset:  0,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	categories, err := s.store.GetCategoriesByUserID(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch categories: %w", err)
	}
	categoryNames := make(map[int64]string, len(categories))
	for i := range categories {
		categoryNames[categories[i].ID] = categories[i].Name
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := []string{
		"date",
		"name",
		"merchant_name",
		"amount",
		"currency_sign",
		"pending",
		"provider",
		"provider_category",
		"tags",
		"tag_categories",
	}
	if err := writer.Write(header); err != nil {
		return 0, fmt.Errorf("failed to write csv header: %w", err)
	}

	rows := 0
	for i := range transactions {
		t := &transactions[i]

		tagNames := make([]string, 0, len(t.Tags))
		tagCats := make([]string, 0, len(t.Tags))
		for j := range t.Tags {
			tagNames = append(tagNames, t.Tags[j].Name)
			if name, ok := categoryNames[t.Tags[j].CategoryID]; ok {
				tagCats = append(tagCats, name)
			} else {
				tagCats = append(tagCats, "")
			}
		}

		// Stored amounts follow Plaid/Stripe sign convention (positive = expense,
		// negative = income). We invert so spreadsheets show debits as negative
		// and credits as positive, matching what the UI displays.
		displayAmount := -t.Amount

		record := []string{
			time.Unix(t.Date, 0).UTC().Format("2006-01-02"),
			t.Name,
			t.MerchantName,
			strconv.FormatFloat(displayAmount, 'f', 2, 64),
			signLabel(displayAmount),
			strconv.FormatBool(t.Pending),
			t.Provider,
			t.PlaidCategory,
			strings.Join(tagNames, "; "),
			strings.Join(tagCats, "; "),
		}

		if err := writer.Write(record); err != nil {
			return rows, fmt.Errorf("failed to write csv row %d: %w", i, err)
		}
		rows++
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return rows, fmt.Errorf("csv writer error: %w", err)
	}
	return rows, nil
}

// signLabel returns a human-readable sign descriptor for a transaction amount.
func signLabel(amount float64) string {
	switch {
	case amount > 0:
		return "credit"
	case amount < 0:
		return "debit"
	default:
		return "zero"
	}
}
