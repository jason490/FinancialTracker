package legacyHandler

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/services/plaid"
	"FinancialTracker/internal/storage"
	"FinancialTracker/web/templ/components"
	"FinancialTracker/web/templ/pages"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
)

// PlaidHandler handles interactions with the Plaid API and Echo requests
type PlaidHandler struct {
	store        *storage.Storage
	plaidService *plaid.PlaidService
}

// NewPlaidHandler initializes a new Plaid service/handler
func NewPlaidHandler(store *storage.Storage, plaid *plaid.PlaidService) *PlaidHandler {
	return &PlaidHandler{
		store:        store,
		plaidService: plaid,
	}
}

// HandleCreateLinkToken handles the request to create a new link token
func (h *PlaidHandler) HandleCreateLinkToken(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	userID := userIDVal.(int64)
	userStr := strconv.FormatInt(userID, 10)

	link, err := h.plaidService.CreateLinkToken(c, userStr)
	if err != nil {
		AddNotification(c, "Failed to initialize bank connection", "error")
		return Render(c, http.StatusInternalServerError, components.StatusMessage("Failed to initialize bank connection", true))
	}

	c.Response().Header().Set("HX-Trigger", fmt.Sprintf(`{"openPlaid": {"token": "%s", "type": "exchange"}}`, link))
	return c.NoContent(http.StatusOK)
}

// HandleCreateUpdateLinkToken handles the request to create a new link token in update mode
func (h *PlaidHandler) HandleCreateUpdateLinkToken(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	userID := userIDVal.(int64)
	userStr := strconv.FormatInt(userID, 10)

	rowID := c.Param("id")
	if rowID == "" {
		return Render(c, http.StatusBadRequest, components.StatusMessage("Missing connection ID", true))
	}

	item, err := h.store.GetPlaidItemByRowID(rowID, userID)
	if err != nil {
		return Render(c, http.StatusInternalServerError, components.StatusMessage("Failed to find bank connection", true))
	}

	link, err := h.plaidService.CreateUpdateLinkToken(c, userStr, item.AccessToken, item.Status)
	if err != nil {
		return Render(c, http.StatusInternalServerError, components.StatusMessage(err.Error(), true))
	}

	c.Response().Header().Set("HX-Trigger", fmt.Sprintf(`{"openPlaid": {"token": "%s", "type": "sync", "rowId": "%s"}}`, link, rowID))
	return c.NoContent(http.StatusOK)
}

// HandleSyncItem handles syncing accounts and transactions for a specific Plaid item
func (h *PlaidHandler) HandleSyncItem(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	userID := userIDVal.(int64)

	rowID := c.Param("id")
	if rowID == "" {
		AddNotification(c, "Missing connection ID", "error")
		return c.NoContent(http.StatusBadRequest)
	}

	item, err := h.store.GetPlaidItemByRowID(rowID, userID)
	if err != nil {
		AddNotification(c, "Failed to find bank connection", "error")
		return c.NoContent(http.StatusInternalServerError)
	}
	if item.Status == plaid.ItemStatusDisconnected {
		return Render(c, http.StatusBadRequest, components.StatusMessage("This bank connection is no longer available. Disconnect and link again.", true))
	}

	ctx := c.Request().Context()
	if err := h.plaidService.SyncPlaidAccounts(&ctx, userID, item.PlaidItemID, item.AccessToken); err != nil {
		log.Errorf("Failed to sync accounts for item %s: %v", item.PlaidItemID, err)
	}

	if err := h.plaidService.SyncItemTransactions(c, userID, item.PlaidItemID, item.AccessToken, item.SyncCursor); err != nil {
		AddNotification(c, "Failed to sync transactions", "error")
	} else {
		_ = h.store.UpdatePlaidItemLastSynced(item.PlaidItemID, time.Now().Unix())
	}

	AddNotification(c, "Connection updated and synced successfully", "success")
	AddHXTriggerEvent(c, "updateBankList")
	AddHXTriggerEvent(c, "updateTransactionList")
	return c.NoContent(http.StatusOK)
}

// HandleExchangeToken handles the exchange of a public token for an access token
func (h *PlaidHandler) HandleExchangeToken(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	userID := userIDVal.(int64)

	publicToken := c.FormValue("public_token")
	if publicToken == "" {
		return Render(c, http.StatusBadRequest, components.StatusMessage("Invalid token received from bank", true))
	}

	if err := h.plaidService.ExchangeToken(c, userID, publicToken); err != nil {
		AddNotification(c, err.Error(), "error")
		return c.NoContent(http.StatusInternalServerError)
	}

	if err := h.plaidService.SyncUser(c, userID); err != nil {
		log.Errorf("Failed to perform initial sync for user %d: %v", userID, err)
	}

	AddNotification(c, "Bank account linked and transactions synced successfully", "success")
	AddHXTriggerEvent(c, "updateBankList")
	AddHXTriggerEvent(c, "updateTransactionList")
	return c.NoContent(http.StatusOK)
}

// HandleSync handles syncing transactions for all connected items
func (h *PlaidHandler) HandleSync(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	userID := userIDVal.(int64)

	if err := h.plaidService.SyncUser(c, userID); err != nil {
		AddNotification(c, err.Error(), "error")
		return Render(c, http.StatusInternalServerError, components.StatusMessage(err.Error(), true))
	}

	AddNotification(c, "Transactions synced successfully", "success")
	AddHXTriggerEvent(c, "updateBankList")
	AddHXTriggerEvent(c, "updateTransactionList")
	AddHXTriggerEvent(c, "refreshDashboard")
	return c.NoContent(http.StatusOK)
}

// HandleManagePage renders the Plaid management page
func (h *PlaidHandler) HandleManagePage(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.Redirect(http.StatusFound, "/login")
	}
	userID := userIDVal.(int64)

	pageData := GetPageData(c, h.store, "Manage Plaid")

	data, err := h.plaidService.GetManagementData(userID)
	if err != nil {
		log.Errorf("Failed to get management data: %v", err)
		return Render(c, http.StatusInternalServerError, pages.Manage(pageData))
	}

	pageData.Data = data
	return Render(c, http.StatusOK, pages.Manage(pageData))
}

// HandleGetConnectionList returns the PlaidConnectionList component
func (h *PlaidHandler) HandleGetConnectionList(c *echo.Context) error {
	userIDVal := c.Get("user_id")
	if userIDVal == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	userID := userIDVal.(int64)

	data, err := h.plaidService.GetManagementData(userID)
	if err != nil {
		return Render(c, http.StatusInternalServerError, components.StatusMessage("Failed to fetch connections", true))
	}

	return Render(c, http.StatusOK, components.PlaidConnectionList(&models.PageData{Data: data}))
}

// HandleRemoveAccount deletes a specific bank account connection
func (h *PlaidHandler) HandleRemoveAccount(c *echo.Context) error {
	userID := c.Get("user_id").(int64)
	accountID := c.Param("id")

	if err := h.plaidService.DeletePlaidAccount(userID, accountID); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "account not found" {
			status = http.StatusNotFound
		} else if err.Error() == "invalid account ID" || err.Error() == "cannot delete an active account; please disconnect it first" {
			status = http.StatusBadRequest
		}
		return Render(c, status, components.StatusMessage(err.Error(), true))
	}

	AddNotification(c, "Account removed successfully", "success")
	AddHXTriggerEvent(c, "updateBankList")
	AddHXTriggerEvent(c, "updateTransactionList")
	return c.NoContent(http.StatusOK)
}

// HandleToggleAccountVisibility toggles the hidden status of an account
func (h *PlaidHandler) HandleToggleAccountVisibility(c *echo.Context) error {
	userID := c.Get("user_id").(int64)
	accountID := c.Param("id")

	isHidden, err := h.plaidService.ToggleAccountVisibility(userID, accountID)
	if err != nil {
		return Render(c, http.StatusInternalServerError, components.StatusMessage(err.Error(), true))
	}

	AddHXTriggerEvent(c, "updateBankList")
	AddHXTriggerEvent(c, "updateTransactionList")
	
	msg := "Account unhidden"
	if isHidden {
		msg = "Account hidden"
		AddNotification(c, "Transactions hidden from list, but totals are hidden.", "info")
	}
	AddNotification(c, msg, "success")
	return c.NoContent(http.StatusOK)
}
