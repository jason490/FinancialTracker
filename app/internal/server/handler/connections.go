package handler

import (
	"FinancialTracker/internal/models/external"
	financialService "FinancialTracker/internal/services/financial"
	"FinancialTracker/internal/services/plaid"
	"FinancialTracker/internal/services/stripefc"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

// ConnectionsHandler serves unified bank connection endpoints.
type ConnectionsHandler struct {
	provider financialService.Provider
}

// NewConnectionsHandler creates a ConnectionsHandler for the active provider.
func NewConnectionsHandler(provider financialService.Provider) *ConnectionsHandler {
	return &ConnectionsHandler{provider: provider}
}

func connectionErrorResponse(err error) (int, APIError) {
	switch {
	case errors.Is(err, stripefc.ErrStripeAPILimitExceeded):
		return http.StatusTooManyRequests, ErrorResponse("connection_api_limit", err.Error())
	case errors.Is(err, stripefc.ErrStripeItemLimitExceeded):
		return http.StatusForbidden, ErrorResponse("connection_item_limit", err.Error())
	case errors.Is(err, plaid.ErrPlaidAPILimitExceeded):
		return http.StatusTooManyRequests, ErrorResponse("connection_api_limit", err.Error())
	case errors.Is(err, plaid.ErrPlaidItemLimitExceeded):
		return http.StatusForbidden, ErrorResponse("connection_item_limit", err.Error())
	case errors.Is(err, plaid.ErrPlaidSyncRateLimited):
		return http.StatusTooManyRequests, ErrorResponse("connection_sync_rate_limit", err.Error())
	default:
		return http.StatusInternalServerError, ErrorResponse("connection_error", err.Error())
	}
}

// HandleGetProvider returns the active financial provider and publishable key.
func (h *ConnectionsHandler) HandleGetProvider(c *echo.Context) error {
	return c.JSON(http.StatusOK, external.ProviderInfoResponse{
		Provider:       h.provider.Name(),
		PublishableKey: financialService.StripePublishableKey(),
	})
}

// HandleGetConnections returns all connections for the user.
func (h *ConnectionsHandler) HandleGetConnections(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	payload, err := h.provider.GetManagementData(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("connection_error", "Failed to load connections"))
	}
	return c.JSON(http.StatusOK, payload)
}

// HandleCreateSession creates a link session for the active provider.
func (h *ConnectionsHandler) HandleCreateSession(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	resp, err := h.provider.CreateSession(c, userID)
	if err != nil {
		status, apiErr := connectionErrorResponse(err)
		if status == http.StatusInternalServerError {
			apiErr = ErrorResponse("connection_error", "Failed to initialize bank connection")
		}
		return c.JSON(status, apiErr)
	}
	return c.JSON(http.StatusOK, resp)
}

// HandleCreateUpdateSession creates an update/relink session for a connection.
func (h *ConnectionsHandler) HandleCreateUpdateSession(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	rowID := c.Param("id")
	if rowID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Missing connection ID"))
	}

	resp, err := h.provider.CreateUpdateSession(c, userID, rowID)
	if err != nil {
		status, apiErr := connectionErrorResponse(err)
		if status == http.StatusInternalServerError {
			apiErr = ErrorResponse("connection_error", err.Error())
		}
		return c.JSON(status, apiErr)
	}
	return c.JSON(http.StatusOK, resp)
}

// HandleCompleteConnection completes the provider link flow.
func (h *ConnectionsHandler) HandleCompleteConnection(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	var req external.CompleteConnectionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	if err := h.provider.CompleteConnection(c, userID, &req); err != nil {
		status, apiErr := connectionErrorResponse(err)
		if status == http.StatusInternalServerError {
			apiErr = ErrorResponse("connection_error", err.Error())
		}
		return c.JSON(status, apiErr)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

// HandleSync syncs all connections for the user.
func (h *ConnectionsHandler) HandleSync(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	if err := h.provider.SyncUser(c, userID); err != nil {
		status, apiErr := connectionErrorResponse(err)
		if status == http.StatusInternalServerError {
			apiErr = ErrorResponse("connection_error", err.Error())
		}
		return c.JSON(status, apiErr)
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

// HandleSyncItem syncs a single institution connection.
func (h *ConnectionsHandler) HandleSyncItem(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	rowID := c.Param("id")
	if rowID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Missing connection ID"))
	}

	if err := h.provider.SyncItem(c, userID, rowID); err != nil {
		status, apiErr := connectionErrorResponse(err)
		if status == http.StatusInternalServerError {
			apiErr = ErrorResponse("connection_error", err.Error())
		}
		return c.JSON(status, apiErr)
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

// HandleDisconnect removes an institution connection.
func (h *ConnectionsHandler) HandleDisconnect(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	rowID := c.Param("id")
	if rowID == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Missing connection ID"))
	}

	if err := h.provider.Disconnect(c, userID, rowID); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("connection_error", "Failed to disconnect bank connection"))
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

// HandleToggleAccountVisibility toggles whether an account is hidden.
func (h *ConnectionsHandler) HandleToggleAccountVisibility(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	accountID := c.Param("id")
	isHidden, err := h.provider.ToggleAccountVisibility(userID, accountID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("connection_error", err.Error()))
	}
	return c.JSON(http.StatusOK, map[string]bool{"is_hidden": isHidden})
}

// HandleRemoveAccount permanently removes a disconnected account.
func (h *ConnectionsHandler) HandleRemoveAccount(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	accountID := c.Param("id")
	if err := h.provider.RemoveAccount(userID, accountID); err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "account not found":
			status = http.StatusNotFound
		case "invalid account ID", "cannot delete an active account; please disconnect it first":
			status = http.StatusBadRequest
		}
		return c.JSON(status, ErrorResponse("connection_error", err.Error()))
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}
