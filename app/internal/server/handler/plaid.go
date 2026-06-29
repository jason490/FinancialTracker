package handler

import (
	plaidService "FinancialTracker/internal/services/plaid"
	"errors"
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
)

// PlaidHandler serves Plaid webhook endpoints.
type PlaidHandler struct {
	plaidService *plaidService.PlaidService
}

// NewPlaidHandler creates a PlaidHandler.
func NewPlaidHandler(plaidService *plaidService.PlaidService) *PlaidHandler {
	return &PlaidHandler{plaidService: plaidService}
}

// HandleWebhook receives Plaid webhooks and triggers background sync for transaction updates.
func (h *PlaidHandler) HandleWebhook(c *echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Missing webhook body"))
	}

	verification := c.Request().Header.Get("Plaid-Verification")
	if err := h.plaidService.HandleWebhook(c.Request().Context(), body, verification); err != nil {
		switch {
		case errors.Is(err, plaidService.ErrWebhookVerification):
			return c.JSON(http.StatusUnauthorized, ErrorResponse("webhook_verification_failed", "Webhook verification failed"))
		case errors.Is(err, plaidService.ErrWebhookInvalidPayload):
			return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid webhook payload"))
		default:
			log.Errorf("Plaid webhook handling failed: %v", err)
		}
	}

	return c.NoContent(http.StatusOK)
}
