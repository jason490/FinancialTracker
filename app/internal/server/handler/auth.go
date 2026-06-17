package handler

import (
	"FinancialTracker/internal/models/external"
	"FinancialTracker/internal/services/auth"
	"FinancialTracker/internal/storage"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
)

// AuthHandler handles JSON authentication endpoints.
type AuthHandler struct {
	store       *storage.Storage
	authService *auth.AuthService
}

// NewAuthHandler creates a new JSON AuthHandler.
func NewAuthHandler(authService *auth.AuthService, store *storage.Storage) *AuthHandler {
	return &AuthHandler{
		store:       store,
		authService: authService,
	}
}

// setSessionCookie configures the Session cookie on the response.
func (h *AuthHandler) setSessionCookie(c *echo.Context, sessionID string, remember bool) {
	cookie := &http.Cookie{
		Name:     "Session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("ENV") != "development", // Set to true in production via middleware or config
		SameSite: http.SameSiteLaxMode,
	}

	if remember {
		cookie.MaxAge = 60 * 60 * 24 * 365
	} else {
		cookie.MaxAge = 60 * 60 * 24
	}

	c.SetCookie(cookie)
}

// HandleLogin authenticates a user and returns a session identifier.
func (h *AuthHandler) HandleLogin(c *echo.Context) error {
	var req external.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	session, err := h.authService.Authenticate(req.Email, req.Password, req.Remember)
	if err != nil {
		switch err.Error() {
		case "email and password are required", "invalid credentials":
			return c.JSON(http.StatusUnauthorized, ErrorResponse("invalid_credentials", err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, ErrorResponse("server_error", "Server error, try again later"))
		}
	}

	h.setSessionCookie(c, session.ID, req.Remember)
	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

// HandleRegister creates a user account and returns a session identifier.
func (h *AuthHandler) HandleRegister(c *echo.Context) error {
	var req external.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	session, err := h.authService.Register(
		req.FirstName,
		req.LastName,
		req.Email,
		req.Password,
		req.ConfirmPassword,
	)
	if err != nil {
		switch err.Error() {
		case "user already exists":
			return c.JSON(http.StatusConflict, ErrorResponse("user_exists", err.Error()))
		default:
			return c.JSON(http.StatusBadRequest, ErrorResponse("validation_error", err.Error()))
		}
	}

	h.setSessionCookie(c, session.ID, false)
	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

// HandleLogout invalidates the current session.
func (h *AuthHandler) HandleLogout(c *echo.Context) error {
	cookie, err := c.Cookie("Session")
	if err == nil {
		_ = h.authService.Logout(cookie.Value)
	}

	// Clear the cookie
	c.SetCookie(&http.Cookie{
		Name:     "Session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, map[string]bool{"logged_out": true})
}

// HandleMe returns the authenticated user's session profile.
func (h *AuthHandler) HandleMe(c *echo.Context) error {
	userIDRaw := c.Get("user_id")
	if userIDRaw == nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}
	userID := userIDRaw.(int64)

	profile, err := h.authService.GetSessionProfile(userID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	return c.JSON(http.StatusOK, profile)
}

// HandleForgotPassword accepts a reset request and always returns success.
func (h *AuthHandler) HandleForgotPassword(c *echo.Context) error {
	var req external.ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	if err := h.authService.RequestPasswordReset(req.Email); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("validation_error", err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "If an account exists for that email, a reset link will be sent shortly.",
	})
}
