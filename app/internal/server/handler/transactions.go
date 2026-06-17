package handler

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/models/external"
	"FinancialTracker/internal/services/transactions"
	"net/http"
	"strconv"
	"strings"

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
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
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
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filters.MinAmount = &f
		}
	}
	if v := c.QueryParam("max_amount"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filters.MaxAmount = &f
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

	return c.JSON(http.StatusOK, external.ToTransactionListPayload(&data))
}

// HandleBulkAddTag adds a tag to multiple transactions at once
func (h *TransactionHandler) HandleBulkAddTag(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
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

// HandleBulkRemoveTag removes a tag from multiple transactions at once.
func (h *TransactionHandler) HandleBulkRemoveTag(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
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
