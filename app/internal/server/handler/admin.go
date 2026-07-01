package handler

import (
	"FinancialTracker/internal/models/external"
	"FinancialTracker/internal/services/auth"
	"net/http"

	"github.com/labstack/echo/v5"
)

// AdminHandler handles operator-only JSON endpoints.
type AdminHandler struct {
	authService *auth.AuthService
}

// NewAdminHandler creates a new JSON AdminHandler.
func NewAdminHandler(authService *auth.AuthService) *AdminHandler {
	return &AdminHandler{authService: authService}
}

// HandleCreateRegistrationCode issues a single-use invite code for gated registration.
func (h *AdminHandler) HandleCreateRegistrationCode(c *echo.Context) error {
	if !auth.RegistrationGateEnabled() {
		return c.JSON(http.StatusForbidden, ErrorResponse("registration_open", "Registration is not invite-only on this server"))
	}

	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	profile, err := h.authService.GetSessionProfile(userID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}
	if !auth.IsRegistrationAdmin(profile.Email) {
		return c.JSON(http.StatusForbidden, ErrorResponse("forbidden", "You are not allowed to create registration codes"))
	}

	createdBy := userID
	code, expiresAt, err := h.authService.GenerateRegistrationCode(&createdBy)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("server_error", "Failed to create registration code"))
	}

	return c.JSON(http.StatusOK, external.CreateRegistrationCodeResponse{
		Code:      code,
		ExpiresAt: expiresAt,
	})
}
