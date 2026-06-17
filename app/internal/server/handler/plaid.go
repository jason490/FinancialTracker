package handler

import (
	"FinancialTracker/internal/models/external"
	plaidService "FinancialTracker/internal/services/plaid"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
)

// PlaidHandler serves JSON Plaid endpoints for the SPA.
type PlaidHandler struct {
	plaidService *plaidService.PlaidService
}

// NewPlaidHandler creates a PlaidHandler.
func NewPlaidHandler(plaidService *plaidService.PlaidService) *PlaidHandler {
	return &PlaidHandler{plaidService: plaidService}
}

// HandleGetConnections returns all Plaid connections and accounts for the user.
func (h *PlaidHandler) HandleGetConnections(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	data, err := h.plaidService.GetManagementData(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("plaid_error", "Failed to load connections"))
	}

	return c.JSON(http.StatusOK, external.ToPlaidConnectionsPayload(data))
}

// HandleCreateLinkToken creates a Plaid Link token for a new connection.
func (h *PlaidHandler) HandleCreateLinkToken(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	link, err := h.plaidService.CreateLinkToken(c, strconv.FormatInt(userID, 10))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("plaid_error", "Failed to initialize bank connection"))
	}

	return c.JSON(http.StatusOK, external.PlaidLinkTokenResponse{LinkToken: link})
}

// HandleCreateUpdateLinkToken creates a Plaid Link token in update mode.
func (h *PlaidHandler) HandleCreateUpdateLinkToken(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	rowID := c.Param("id")
	if rowID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Missing connection ID"))
	}

	item, err := h.plaidService.GetItemByRowID(rowID, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse("not_found", "Connection not found"))
	}

	link, err := h.plaidService.CreateUpdateLinkToken(c, strconv.FormatInt(userID, 10), item.AccessToken, item.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("plaid_error", err.Error()))
	}

	return c.JSON(http.StatusOK, external.PlaidLinkTokenResponse{LinkToken: link})
}

// HandleExchangeToken exchanges a Plaid public token for a new connection.
func (h *PlaidHandler) HandleExchangeToken(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	var req external.PlaidExchangeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}
	if req.PublicToken == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Public token is required"))
	}

	if err := h.plaidService.ExchangeToken(c, userID, req.PublicToken); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("plaid_error", err.Error()))
	}

	if err := h.plaidService.SyncUser(c, userID); err != nil {
		log.Errorf("Failed to perform initial sync for user %d: %v", userID, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

// HandleSyncItem syncs accounts and transactions for a specific Plaid item.
func (h *PlaidHandler) HandleSyncItem(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	rowID := c.Param("id")
	if rowID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Missing connection ID"))
	}

	item, err := h.plaidService.GetItemByRowID(rowID, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse("not_found", "Connection not found"))
	}
	if item.Status == plaidService.ItemStatusDisconnected {
		return c.JSON(http.StatusBadRequest, ErrorResponse("plaid_error", "This bank connection is no longer available. Disconnect and link again."))
	}

	ctx := c.Request().Context()
	if err := h.plaidService.SyncPlaidAccounts(&ctx, userID, item.PlaidItemID, item.AccessToken); err != nil {
		log.Errorf("Failed to sync accounts for item %s: %v", item.PlaidItemID, err)
	}

	if err := h.plaidService.SyncItemTransactions(c, userID, item.PlaidItemID, item.AccessToken, item.SyncCursor); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("plaid_error", "Failed to sync transactions"))
	}

	_ = h.plaidService.UpdateItemLastSynced(item.PlaidItemID, time.Now().Unix())
	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

// HandleSync syncs all Plaid connections for the user.
func (h *PlaidHandler) HandleSync(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	if err := h.plaidService.SyncUser(c, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("plaid_error", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

// HandleDisconnectItem removes an entire Plaid institution connection.
func (h *PlaidHandler) HandleDisconnectItem(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	rowID := c.Param("id")
	if rowID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Missing connection ID"))
	}

	ctx := c.Request().Context()
	if err := h.plaidService.DisconnectItem(&ctx, rowID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("plaid_error", "Failed to disconnect bank connection"))
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

// HandleToggleAccountVisibility toggles whether an account is hidden.
func (h *PlaidHandler) HandleToggleAccountVisibility(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	accountID := c.Param("id")
	isHidden, err := h.plaidService.ToggleAccountVisibility(userID, accountID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("plaid_error", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]bool{"is_hidden": isHidden})
}

// HandleRemoveAccount permanently removes a disconnected bank account.
func (h *PlaidHandler) HandleRemoveAccount(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	accountID := c.Param("id")
	if err := h.plaidService.DeletePlaidAccount(userID, accountID); err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "account not found":
			status = http.StatusNotFound
		case "invalid account ID", "cannot delete an active account; please disconnect it first":
			status = http.StatusBadRequest
		}
		return c.JSON(status, ErrorResponse("plaid_error", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}
