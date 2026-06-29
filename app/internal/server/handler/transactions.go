package handler

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/models/external"
	"FinancialTracker/internal/services/transactions"
	"bytes"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
)

// TransactionHandler serves JSON transaction endpoints for the SPA.
type TransactionHandler struct {
	transService *transactions.TransactionService
}

// NewTransactionHandler creates a TransactionHandler.
func NewTransactionHandler(transService *transactions.TransactionService) *TransactionHandler {
	return &TransactionHandler{transService: transService}
}

// HandleGetTransactions returns a paginated, filterable list of transactions.
func (h *TransactionHandler) HandleGetTransactions(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 25
	}

	filters := models.TransactionFilters{
		Search:  c.QueryParam("search"),
		SortBy:  c.QueryParam("sort_by"),
		SortDir: c.QueryParam("sort_dir"),
	}

	if filters.SortBy == "" {
		filters.SortBy = "date"
	}
	if filters.SortDir == "" {
		filters.SortDir = "desc"
	}

	if v := c.QueryParam("min_amount"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && !math.IsNaN(f) && !math.IsInf(f, 0) {
			abs := math.Abs(f)
			filters.MinAmount = &abs
		}
	}
	if v := c.QueryParam("max_amount"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && !math.IsNaN(f) && !math.IsInf(f, 0) {
			abs := math.Abs(f)
			filters.MaxAmount = &abs
		}
	}
	if v := c.QueryParam("start_date"); v != "" {
		if d, err := strconv.ParseInt(v, 10, 64); err == nil {
			filters.StartDate = &d
		}
	}
	if v := c.QueryParam("end_date"); v != "" {
		if d, err := strconv.ParseInt(v, 10, 64); err == nil {
			filters.EndDate = &d
		}
	}
	if v := c.QueryParam("category_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			filters.CategoryID = &id
		}
	}
	if v := c.QueryParam("tags"); v != "" {
		for _, s := range strings.Split(v, ",") {
			if id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64); err == nil {
				filters.Tags = append(filters.Tags, id)
			}
		}
	}

	data, err := h.transService.GetTransactionsPage(userID, filters, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("transactions_error", "Failed to load transactions"))
	}

	return c.JSON(http.StatusOK, transactions.BuildListPayload(&data))
}

// HandleBulkAddTag adds a tag to multiple transactions at once
func (h *TransactionHandler) HandleBulkAddTag(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	var req external.BulkTagRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	if err := h.transService.BulkAddTag(userID, req.TransactionIDs, req.TagID); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("bulk_tag_error", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// HandleExportTransactions streams every visible transaction for the
// authenticated user as a ZIP archive containing a single CSV file. The export
// consumes one API call from the user's monthly quota; when the quota is
// exhausted the response is HTTP 429 with code `export_api_limit`.
func (h *TransactionHandler) HandleExportTransactions(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	stamp := time.Now().UTC().Format("2006-01-02")
	csvName := fmt.Sprintf("financial-tracker-transactions-%s.csv", stamp)
	zipName := fmt.Sprintf("financial-tracker-transactions-%s.zip", stamp)

	// Buffer the archive so we only commit response headers after we know the
	// quota check passed and the archive built successfully. Exports are small
	// (rarely > a few MB even for power users) so the memory cost is fine and
	// it avoids streaming a partial file when an error surfaces mid-write.
	var buf bytes.Buffer
	if _, err := h.transService.ExportTransactionsZip(userID, csvName, &buf); err != nil {
		if errors.Is(err, transactions.ErrExportAPILimitExceeded) {
			return c.JSON(http.StatusTooManyRequests, ErrorResponse(
				"export_api_limit",
				"Monthly API limit reached for your plan. Upgrade or wait for the next billing period to export again.",
			))
		}
		c.Logger().Error("transaction export failed", "user_id", userID, "error", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse("export_failed", "Failed to build export"))
	}

	res := c.Response()
	res.Header().Set(echo.HeaderContentType, "application/zip")
	res.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, zipName))
	res.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	res.Header().Set("Cache-Control", "no-store")
	res.Header().Set("X-Content-Type-Options", "nosniff")
	res.WriteHeader(http.StatusOK)

	if _, err := buf.WriteTo(res); err != nil {
		c.Logger().Error("transaction export write failed", "user_id", userID, "error", err)
	}
	return nil
}

// HandleBulkRemoveTag removes a tag from multiple transactions at once.
func (h *TransactionHandler) HandleBulkRemoveTag(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	var req external.BulkTagRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	if err := h.transService.BulkRemoveTag(userID, req.TransactionIDs, req.TagID); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("bulk_tag_error", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
