package legacyHandler

import (
	"FinancialTracker/internal/models"
	plaidService "FinancialTracker/internal/services/plaid"
	"FinancialTracker/internal/services/settings"
	"FinancialTracker/internal/storage"
	"FinancialTracker/web/templ/components"
	settingsComponents "FinancialTracker/web/templ/components/settings"
	"FinancialTracker/web/templ/pages"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
)

// SettingsHandler handles settings-related requests
type SettingsHandler struct {
	store           *storage.Storage
	plaidService    *plaidService.PlaidService
	settingsService *settings.SettingsService
}

// NewSettingsHandler creates a new SettingsHandler
func NewSettingsHandler(store *storage.Storage, plaid *plaidService.PlaidService, settingsService *settings.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		store:           store,
		plaidService:    plaid,
		settingsService: settingsService,
	}
}

// HandleSettingsPage renders the settings page
func (h *SettingsHandler) HandleSettingsPage(c *echo.Context) error {
	pageData := GetPageData(c, h.store, "Settings")

	if errCode := c.QueryParam("error"); errCode != "" {
		if msg, ok := ErrorMessages[errCode]; ok {
			QueuePageNotification(pageData, msg, "error")
		}
	}
	if success := c.QueryParam("success"); success == "linked" {
		QueuePageNotification(pageData, "Google account linked", "success")
	}

	if c.QueryParam("reauth_success") == "true" {
		pageData.ReauthSuccess = true
	}

	if pageData.User != nil {
		items, _ := h.store.GetPlaidItemsByUserID(pageData.User.ID)
		pageData.Data = items
	} else {
		log.Error("model.User not in pageData")
	}
	return Render(c, http.StatusOK, pages.Settings(pageData))
}

func (h *SettingsHandler) HandleUpdateBankList(c *echo.Context) error {
	userId := c.Get("user_id")
	if userId == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	items, _ := h.store.GetPlaidItemsByUserID(userId.(int64))
	pageData := &models.PageData{
		Data: items,
	}
	return Render(c, http.StatusOK, settingsComponents.BankAccountList(pageData))
}

// HandleRemoveBankAccount deletes a bank account connection securely
func (h *SettingsHandler) HandleRemoveBankAccount(c *echo.Context) error {
	userID := c.Get("user_id").(int64)
	rowID := c.Param("id")

	if rowID == "" {
		AddNotification(c, "Invalid account ID", "error")
		return c.NoContent(http.StatusBadRequest)
	}

	ctx := c.Request().Context()
	if err := h.plaidService.DisconnectItem(&ctx, rowID, userID); err != nil {
		AddNotification(c, "Failed to remove bank account", "error")
		return c.NoContent(http.StatusInternalServerError)
	}

	AddNotification(c, "Bank account removed successfully", "success")
	AddHXTriggerEvent(c, "updateBankList")
	AddHXTriggerEvent(c, "updateTransactionList")
	return c.NoContent(http.StatusOK)
}

// HandleUpdateTheme processes the theme update request
func (h *SettingsHandler) HandleUpdateTheme(c *echo.Context) error {
	userID := c.Get("user_id").(int64)
	theme := c.FormValue("theme")

	if err := h.settingsService.UpdateTheme(userID, theme); err != nil {
		AddNotification(c, err.Error(), "error")
		return Render(c, http.StatusBadRequest, components.StatusMessage(err.Error(), true))
	}

	if c.Request().Header.Get("HX-Request") != "" {
		AddNotification(c, "Preference saved!", "success")
		AddHXTriggerEvent(c, "themeUpdated")
		return c.NoContent(http.StatusOK)
	}

	return c.Redirect(http.StatusFound, "/settings")
}

// HandleUpdateAccount processes account info updates
func (h *SettingsHandler) HandleUpdateAccount(c *echo.Context) error {
	userID := c.Get("user_id").(int64)
	firstName := c.FormValue("first_name")
	lastName := c.FormValue("last_name")
	email := c.FormValue("email")

	if err := h.settingsService.UpdateProfile(userID, firstName, lastName, email); err != nil {
		AddNotification(c, err.Error(), "error")
		return Render(c, http.StatusBadRequest, components.StatusMessage(err.Error(), true))
	}

	if c.Request().Header.Get("HX-Request") != "" {
		AddNotification(c, "Account updated successfully!", "success")
		AddHXTriggerEvent(c, "profileUpdated")
		return c.NoContent(http.StatusOK)
	}

	return c.Redirect(http.StatusFound, "/settings")
}

// HandleUpdatePassword processes password updates
func (h *SettingsHandler) HandleUpdatePassword(c *echo.Context) error {
	userID := c.Get("user_id").(int64)
	currentPassword := c.FormValue("current_password")
	newPassword := c.FormValue("new_password")
	confirmPassword := c.FormValue("confirm_password")

	if err := h.settingsService.UpdatePassword(userID, currentPassword, newPassword, confirmPassword); err != nil {
		AddNotification(c, err.Error(), "error")
		return Render(c, http.StatusBadRequest, components.StatusMessage(err.Error(), true))
	}

	cookie, err := c.Cookie("Session")
	if err == nil {
		_ = h.store.DeleteSession(cookie.Value)
	}

	newCookie := &http.Cookie{
		Name:     "Session",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	c.SetCookie(newCookie)

	c.Response().Header().Set("HX-Redirect", "/login")
	return c.Redirect(http.StatusFound, "/login")
}

// HandleUnlinkSSO handles unlinking an SSO provider
func (h *SettingsHandler) HandleUnlinkSSO(c *echo.Context) error {
	userID := c.Get("user_id").(int64)
	provider := c.Param("provider")

	user, err := h.store.GetUserByID(userID)
	if err != nil || user == nil {
		AddNotification(c, "User not found", "error")
		return c.NoContent(http.StatusBadRequest)
	}

	if user.PasswordHash == "" && len(user.SSOs) <= 1 {
		AddNotification(c, "Cannot remove your only login method", "error")
		return c.NoContent(http.StatusBadRequest)
	}

	if err := h.store.UnlinkSSO(userID, provider); err != nil {
		AddNotification(c, "Failed to unlink SSO", "error")
		return c.NoContent(http.StatusInternalServerError)
	}

	AddNotification(c, "SSO unlinked", "success")
	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

// HandleDeleteAccountInit returns the re-auth form
func (h *SettingsHandler) HandleDeleteAccountInit(c *echo.Context) error {
	return Render(c, http.StatusOK, pages.DeleteAccountReAuth(GetPageData(c, h.store, "")))
}

// HandleDeleteAccountCheckReauth checks if user just came back from SSO re-auth
func (h *SettingsHandler) HandleDeleteAccountCheckReauth(c *echo.Context) error {
	reauth := c.QueryParam("reauth_success")

	if reauth == "true" {
		session := c.Get("session").(*models.Session)
		if time.Now().Unix()-session.ReauthenticatedAt < ReauthTimeout {
			return Render(c, http.StatusOK, pages.DeleteAccountConfirm())
		}
		return Render(c, http.StatusOK, pages.DeleteAccountButton())
	}
	return c.NoContent(http.StatusBadRequest)
}

// HandleDeleteAccountVerify verifies credentials for deletion (Local Auth)
func (h *SettingsHandler) HandleDeleteAccountVerify(c *echo.Context) error {
	userID := c.Get("user_id").(int64)
	email := c.FormValue("email")
	password := c.FormValue("password")

	if err := h.settingsService.VerifyReauth(userID, email, password); err != nil {
		return Render(c, http.StatusBadRequest, components.StatusMessage(err.Error(), true))
	}

	if err := h.store.UpdateSessionReauth(c.Get("session").(*models.Session).ID, time.Now().Unix()); err != nil {
		return Render(c, http.StatusInternalServerError, components.StatusMessage("Database error", true))
	}

	return Render(c, http.StatusOK, pages.DeleteAccountConfirm())
}

// HandleDeleteAccountConfirm performs the final deletion
func (h *SettingsHandler) HandleDeleteAccountConfirm(c *echo.Context) error {
	userID := c.Get("user_id").(int64)
	session := c.Get("session").(*models.Session)

	if time.Now().Unix()-session.ReauthenticatedAt > ReauthTimeout {
		return Render(c, http.StatusUnauthorized, components.StatusMessage("Verification expired. Please re-authenticate.", true))
	}

	if err := h.settingsService.DeleteAccount(userID); err != nil {
		return Render(c, http.StatusInternalServerError, components.StatusMessage("Failed to delete account", true))
	}

	newCookie := &http.Cookie{
		Name:     "Session",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	c.SetCookie(newCookie)

	if c.Request().Header.Get("HX-Request") != "" {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusOK)
	}

	return c.Redirect(http.StatusFound, "/")
}

// HandleDeleteAccountCancel reverts to the initial state
func (h *SettingsHandler) HandleDeleteAccountCancel(c *echo.Context) error {
	return Render(c, http.StatusOK, pages.DeleteAccountButton())
}
