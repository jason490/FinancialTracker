package handler

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/models/external"
	"FinancialTracker/internal/services/settings"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
)

// SettingsHandler serves JSON settings endpoints for the SPA.
type SettingsHandler struct {
	settingsService *settings.SettingsService
}

// NewSettingsHandler creates a SettingsHandler.
func NewSettingsHandler(settingsService *settings.SettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: settingsService}
}

// HandleGetSettings returns the authenticated user's settings profile.
func (h *SettingsHandler) HandleGetSettings(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	profile, err := h.settingsService.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("settings_error", "Failed to load settings"))
	}

	return c.JSON(http.StatusOK, profile)
}

// HandleUpdateProfile updates the user's first and last name.
func (h *SettingsHandler) HandleUpdateProfile(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	var req external.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	if err := h.settingsService.UpdateNames(userID, req.FirstName, req.LastName); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("validation_error", err.Error()))
	}

	profile, err := h.settingsService.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("settings_error", "Failed to load settings"))
	}

	return c.JSON(http.StatusOK, profile)
}

// HandleUpdatePassword changes the user's password and clears the session.
func (h *SettingsHandler) HandleUpdatePassword(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	var req external.UpdatePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	if err := h.settingsService.UpdatePassword(userID, req.CurrentPassword, req.NewPassword, req.ConfirmPassword); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("validation_error", err.Error()))
	}

	cookie, err := c.Cookie("Session")
	if err == nil {
		_ = h.settingsService.InvalidateSession(cookie.Value)
	}

	c.SetCookie(&http.Cookie{
		Name:     "Session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	return c.JSON(http.StatusOK, map[string]bool{"logged_out": true})
}

// HandleUpdateTheme persists the user's theme preference.
func (h *SettingsHandler) HandleUpdateTheme(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	var req external.UpdateThemeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	if err := h.settingsService.UpdateTheme(userID, req.Theme); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("validation_error", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{"theme_preference": req.Theme})
}

// HandleUnlinkSSO removes a linked SSO provider from the user account.
func (h *SettingsHandler) HandleUnlinkSSO(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	provider := c.Param("provider")
	if provider == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Provider is required"))
	}

	if err := h.settingsService.UnlinkSSO(userID, provider); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("unlink_error", err.Error()))
	}

	profile, err := h.settingsService.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("settings_error", "Failed to load settings"))
	}

	return c.JSON(http.StatusOK, profile)
}

// HandleDeleteAccountReauthStatus returns whether the session has passed re-auth recently.
func (h *SettingsHandler) HandleDeleteAccountReauthStatus(c *echo.Context) error {
	session, ok := c.Get("session").(*models.Session)
	if !ok || session == nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	return c.JSON(http.StatusOK, external.DeleteAccountReauthStatusResponse{
		ReauthVerified: h.settingsService.IsSessionReauthenticated(session),
	})
}

// HandleDeleteAccountVerify verifies password credentials for account deletion.
func (h *SettingsHandler) HandleDeleteAccountVerify(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	cookie, err := c.Cookie("Session")
	if err != nil || cookie.Value == "" {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	var req external.DeleteAccountVerifyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	if err := h.settingsService.CompletePasswordReauth(userID, cookie.Value, req.Email, req.Password); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("validation_error", err.Error()))
	}

	return c.JSON(http.StatusOK, external.DeleteAccountReauthStatusResponse{ReauthVerified: true})
}

// HandleDeleteAccountConfirm permanently deletes the user account.
func (h *SettingsHandler) HandleDeleteAccountConfirm(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	session, ok := c.Get("session").(*models.Session)
	if !ok || session == nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	if err := h.settingsService.DeleteAccountAfterReauth(userID, session); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "verification expired") {
			return c.JSON(http.StatusUnauthorized, ErrorResponse("reauth_expired", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse("delete_error", "Failed to delete account"))
	}

	c.SetCookie(&http.Cookie{
		Name:     "Session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	return c.JSON(http.StatusOK, map[string]bool{"deleted": true})
}
