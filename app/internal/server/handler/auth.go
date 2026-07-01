package handler

import (
	"FinancialTracker/internal/config"
	"FinancialTracker/internal/models/external"
	"FinancialTracker/internal/services/auth"
	"FinancialTracker/internal/storage"
	"errors"
	"net/http"
	"strings"

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
		Secure:   !config.IsDevelopment(),
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
		req.RegistrationCode,
	)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidRegistrationCode):
			return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_registration_code", "Invalid or expired registration code"))
		case errors.Is(err, auth.ErrRegistrationCodeRequired):
			return c.JSON(http.StatusBadRequest, ErrorResponse("registration_code_required", "Registration code is required"))
		case err.Error() == "user already exists":
			return c.JSON(http.StatusConflict, ErrorResponse("user_exists", err.Error()))
		default:
			return c.JSON(http.StatusBadRequest, ErrorResponse("validation_error", err.Error()))
		}
	}

	h.setSessionCookie(c, session.ID, false)
	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

// HandleRegistrationConfig reports whether invite codes are required to register.
func (h *AuthHandler) HandleRegistrationConfig(c *echo.Context) error {
	required := auth.RegistrationGateEnabled()
	return c.JSON(http.StatusOK, external.RegistrationConfigResponse{
		RegistrationCodeRequired: required,
		CodeExpiresInSeconds:     auth.RegistrationCodeTTLSeconds,
	})
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
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	profile, err := h.authService.GetSessionProfile(userID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	return c.JSON(http.StatusOK, profile)
}

// HandleCompleteOnboarding marks the authenticated user's onboarding as finished.
func (h *AuthHandler) HandleCompleteOnboarding(c *echo.Context) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	if err := h.authService.CompleteOnboarding(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("server_error", "Failed to complete onboarding"))
	}

	return c.JSON(http.StatusOK, map[string]bool{"onboarding_completed": true})
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

	return c.JSON(http.StatusOK, external.ForgotPasswordResponse{
		Message:              "If an account exists for that email, a reset code will be sent shortly.",
		CodeExpiresInSeconds: auth.ResetCodeTTLSeconds,
	})
}

// HandleVerifyResetCode checks a reset code before the user sets a new password.
func (h *AuthHandler) HandleVerifyResetCode(c *echo.Context) error {
	var req external.VerifyResetCodeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	expiresAt, err := h.authService.VerifyPasswordResetCode(req.Email, req.Code)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidResetCode) {
			return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_reset_code", "Invalid or expired reset code"))
		}
		msg := err.Error()
		if msg == "email is required" ||
			msg == "reset code must be 6 digits" ||
			strings.HasPrefix(msg, "Invalid email") {
			return c.JSON(http.StatusBadRequest, ErrorResponse("validation_error", msg))
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse("server_error", "Server error, try again later"))
	}

	return c.JSON(http.StatusOK, external.VerifyResetCodeResponse{ExpiresAt: expiresAt})
}

// HandleResetPassword verifies a reset code and sets a new password.
func (h *AuthHandler) HandleResetPassword(c *echo.Context) error {
	var req external.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	if err := h.authService.ConfirmPasswordReset(req.Email, req.Code, req.NewPassword, req.ConfirmPassword); err != nil {
		if errors.Is(err, auth.ErrInvalidResetCode) {
			return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_reset_code", "Invalid or expired reset code"))
		}
		msg := err.Error()
		if msg == "email is required" ||
			msg == "passwords do not match" ||
			msg == "reset code must be 6 digits" ||
			strings.HasPrefix(msg, "Password must") ||
			strings.HasPrefix(msg, "Invalid email") {
			return c.JSON(http.StatusBadRequest, ErrorResponse("validation_error", msg))
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse("server_error", "Server error, try again later"))
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}
