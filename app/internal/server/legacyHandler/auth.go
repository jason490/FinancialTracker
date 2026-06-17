package legacyHandler

import (
	"FinancialTracker/internal/services/auth"
	"FinancialTracker/internal/storage"
	"FinancialTracker/web/templ/components"
	"FinancialTracker/web/templ/pages"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	store       *storage.Storage
	authService *auth.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(store *storage.Storage, authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		store:       store,
		authService: authService,
	}
}

// HandleLoginPage renders the login page
func (h *AuthHandler) HandleLoginPage(c *echo.Context) error {
	if getCurrentUser(c, h.store) != nil {
		return c.Redirect(http.StatusFound, "/dashboard")
	}
	pageData := GetPageData(c, h.store, "Login")
	if errCode := c.QueryParam("error"); errCode != "" {
		if msg, ok := ErrorMessages[errCode]; ok {
			QueuePageNotification(pageData, msg, "error")
		}
	}
	return Render(c, http.StatusOK, pages.Login(pageData))
}

// HandleLoginSubmit processes the login form
func (h *AuthHandler) HandleLoginSubmit(c *echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	remember := c.FormValue("remember") == "on"

	session, err := h.authService.Authenticate(email, password, remember)
	if err != nil {
		log.Debug("Tesing")
		// Distinguish between validation/auth errors and server errors
		if err.Error() == "email and password are required" || err.Error() == "invalid credentials" {
			return Render(c, http.StatusUnauthorized, components.ErrorMessage(err.Error()))
		}
		return Render(c, http.StatusInternalServerError, components.ErrorMessage("server error, try again later"))
	}

	c.SetCookie(createCookie(session.ID, remember))

	if c.Request().Header.Get("HX-Request") != "" {
		c.Response().Header().Set("HX-Redirect", "/dashboard")
		return c.NoContent(http.StatusOK)
	}

	return c.Redirect(http.StatusFound, "/dashboard")
}

// HandleRegisterPage renders the registration page
func (h *AuthHandler) HandleRegisterPage(c *echo.Context) error {
	if getCurrentUser(c, h.store) != nil {
		return c.Redirect(http.StatusFound, "/dashboard")
	}
	pageData := GetPageData(c, h.store, "Register")
	if errCode := c.QueryParam("error"); errCode != "" {
		if msg, ok := ErrorMessages[errCode]; ok {
			QueuePageNotification(pageData, msg, "error")
		}
	}
	return Render(c, http.StatusOK, pages.Register(pageData))
}

// HandleRegisterSubmit processes the registration form
func (h *AuthHandler) HandleRegisterSubmit(c *echo.Context) error {
	firstName := c.FormValue("first_name")
	lastName := c.FormValue("last_name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	session, err := h.authService.Register(firstName, lastName, email, password, confirmPassword)
	if err != nil {
		// Registration errors are typically 400 Bad Request
		return Render(c, http.StatusBadRequest, components.ErrorMessage(err.Error()))
	}

	c.SetCookie(createCookie(session.ID, false))

	if c.Request().Header.Get("HX-Request") != "" {
		c.Response().Header().Set("HX-Redirect", "/dashboard")
		return c.NoContent(http.StatusOK)
	}

	return c.Redirect(http.StatusFound, "/dashboard")
}

// HandleLogout logs the user out
func (h *AuthHandler) HandleLogout(c *echo.Context) error {
	cookie, err := c.Cookie("Session")
	if err == nil {
		_ = h.authService.Logout(cookie.Value)
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
