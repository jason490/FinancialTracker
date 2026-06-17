package legacyHandler

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/services/transactions"
	"FinancialTracker/internal/storage"
	"FinancialTracker/web/templ/components"
	"FinancialTracker/web/templ/pages"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
)

type TransactionHandler struct {
	store      *storage.Storage
	transService *transactions.TransactionService
}

func NewTransactionHandler(store *storage.Storage, transService *transactions.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		store:      store,
		transService: transService,
	}
}

// HandleTransactions renders the transactions page
func (h *TransactionHandler) HandleTransactions(c *echo.Context) error {
	pageData := GetPageData(c, h.store, "Transactions")
	if pageData.User == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	// Pagination
	pageSize := 20
	pageStr := c.QueryParam("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	// Sorting
	sort := c.QueryParam("sort")
	sortBy := "date"
	sortDir := "desc"
	switch sort {
	case "date_asc":
		sortDir = "asc"
	case "amount_desc":
		sortBy = "amount"
		sortDir = "desc"
	case "amount_asc":
		sortBy = "amount"
		sortDir = "asc"
	case "name_asc":
		sortBy = "name"
		sortDir = "asc"
	case "name_desc":
		sortBy = "name"
		sortDir = "desc"
	}

	// Filters
	search := c.QueryParam("search")
	tagIDStr := c.QueryParam("tag_id")
	categoryIDStr := c.QueryParam("category_id")
	filters := models.TransactionFilters{
		Search:  search,
		SortBy:  sortBy,
		SortDir: sortDir,
	}

	if categoryIDStr != "" {
		if catID, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			filters.CategoryID = &catID
		}
	}

	if tagIDStr != "" {
		if tagID, err := strconv.ParseInt(tagIDStr, 10, 64); err == nil {
			filters.Tags = []int64{tagID}
		}
	}

	if min := c.QueryParam("min_amount"); min != "" {
		if val, err := strconv.ParseFloat(min, 64); err == nil {
			filters.MinAmount = &val
		}
	}
	if max := c.QueryParam("max_amount"); max != "" {
		if val, err := strconv.ParseFloat(max, 64); err == nil {
			filters.MaxAmount = &val
		}
	}

	data, err := h.transService.GetTransactionsPage(pageData.User.ID, filters, page, pageSize)
	if err != nil {
		log.Errorf("Failed to fetch transactions page data: %v", err)
		return Render(c, http.StatusInternalServerError, components.StatusMessage("Failed to load transactions", true))
	}

	pageData.Data = data
	return Render(c, http.StatusOK, pages.Transactions(pageData))
}

// HandleBulkAddTag adds a tag to multiple transactions
func (h *TransactionHandler) HandleBulkAddTag(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tagIDStr := c.FormValue("tag_id")
	tagID, _ := strconv.ParseInt(tagIDStr, 10, 64)
	
	if err := c.Request().ParseForm(); err != nil {
		return Render(c, http.StatusBadRequest, components.StatusMessage("Failed to parse form", true))
	}
	transactionIDs := c.Request().PostForm["transaction_ids"]

	if err := h.transService.BulkAddTag(user.ID, transactionIDs, tagID); err != nil {
		log.Errorf("Bulk add tag failed: %v", err)
		AddNotification(c, err.Error(), "error")
		return c.NoContent(http.StatusBadRequest)
	}

	AddNotification(c, "Tag added to selected transactions", "success")
	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

// HandleBulkRemoveTag removes a tag from multiple transactions
func (h *TransactionHandler) HandleBulkRemoveTag(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tagIDStr := c.FormValue("tag_id")
	tagID, _ := strconv.ParseInt(tagIDStr, 10, 64)
	
	if err := c.Request().ParseForm(); err != nil {
		return Render(c, http.StatusBadRequest, components.StatusMessage("Failed to parse form", true))
	}
	transactionIDs := c.Request().PostForm["transaction_ids"]

	if err := h.transService.BulkRemoveTag(user.ID, transactionIDs, tagID); err != nil {
		log.Errorf("Bulk remove tag failed: %v", err)
		AddNotification(c, err.Error(), "error")
		return c.NoContent(http.StatusBadRequest)
	}

	AddNotification(c, "Tag removed from selected transactions", "success")
	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

// HandleGetList renders only the transaction list component
func (h *TransactionHandler) HandleGetList(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	userID := userIDVal.(int64)

	transactions, err := h.transService.GetRecentTransactions(userID, 10)
	if err != nil {
		log.Errorf("Failed to fetch recent transactions: %v", err)
		return Render(c, http.StatusInternalServerError, components.StatusMessage("Failed to load transactions", true))
	}

	return Render(c, http.StatusOK, components.TransactionList(transactions))
}
