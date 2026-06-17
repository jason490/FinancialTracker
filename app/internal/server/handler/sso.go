package handler

import (
	"FinancialTracker/internal/services/auth"
	"errors"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
)

// SSOHandler handles JSON SSO endpoints.
type SSOHandler struct {
	googleOauthConfig *oauth2.Config
	ssoService        *auth.SSOService
}

// NewSSOHandler creates a new JSON SSOHandler.
func NewSSOHandler(
	googleOauthConfig *oauth2.Config,
	ssoService *auth.SSOService,
) *SSOHandler {
	return &SSOHandler{
		googleOauthConfig: googleOauthConfig,
		ssoService:        ssoService,
	}
}

type ssoExchangeRequest struct {
	Token string `json:"token"`
}

// HandleGoogleLogin redirects to Google's consent page for API clients.
func (h *SSOHandler) HandleGoogleLogin(c *echo.Context) error {
	state := h.ssoService.BuildAPIOAuthState(c.QueryParam("return_to"), c.QueryParam("action"))
	authURL := h.googleOauthConfig.AuthCodeURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// HandleGoogleCallback is the HTTP entry point for Google's OAuth redirect.
func (h *SSOHandler) HandleGoogleCallback(c *echo.Context) error {
	oauthState, err := h.ssoService.ParseAPIState(c.QueryParam("state"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_state", "Invalid OAuth state"))
	}

	if oauthState.Mode == "link" {
		return h.handleGoogleLinkCallback(c, oauthState)
	}

	if oauthState.Mode == "reauth-delete" {
		return h.handleGoogleReauthDeleteCallback(c, oauthState)
	}

	result, err := h.ssoService.CompleteAPICallback(
		c.Request().Context(),
		h.googleOauthConfig,
		c.QueryParam("code"),
		c.QueryParam("state"),
	)
	if err != nil {
		log.Error(err)
		if errors.Is(err, auth.ErrInvalidOAuthState) {
			return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_state", "Invalid OAuth state"))
		}

		if result != nil && result.ReturnTo != "" {
			return redirectWithError(c, result.ReturnTo, "authentication_failed")
		}

		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_state", "Invalid OAuth state"))
	}

	// Set the session cookie directly during the redirect
	if result.SessionID != "" {
		c.SetCookie(&http.Cookie{
			Name:     "Session",
			Value:    result.SessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   os.Getenv("ENV") != "development", // Set to true in production
			SameSite: http.SameSiteLaxMode,
			MaxAge:   60 * 60 * 24, // 1 day default for SSO
		})
	}

	return redirectWithQuery(c, result.ReturnTo, nil)
}

func (h *SSOHandler) handleGoogleLinkCallback(c *echo.Context, oauthState *auth.APIOAuthState) error {
	cookie, err := c.Cookie("Session")
	if err != nil || cookie.Value == "" {
		return redirectWithError(c, oauthState.ReturnTo, "authentication_failed")
	}

	user, err := h.ssoService.GetUserForSession(cookie.Value)
	if err != nil {
		return redirectWithError(c, oauthState.ReturnTo, "authentication_failed")
	}

	if err := h.ssoService.CompleteGoogleLink(
		c.Request().Context(),
		h.googleOauthConfig,
		c.QueryParam("code"),
		user.ID,
	); err != nil {
		log.Error(err)
		code := "link_failed"
		if err.Error() == "SSO account already linked to another user" {
			code = "link_conflict"
		}
		return redirectWithError(c, oauthState.ReturnTo, code)
	}

	return redirectWithQuery(c, oauthState.ReturnTo, map[string]string{"success": "linked"})
}

func (h *SSOHandler) handleGoogleReauthDeleteCallback(c *echo.Context, oauthState *auth.APIOAuthState) error {
	cookie, err := c.Cookie("Session")
	if err != nil || cookie.Value == "" {
		return redirectWithError(c, oauthState.ReturnTo, "authentication_failed")
	}

	user, err := h.ssoService.GetUserForSession(cookie.Value)
	if err != nil {
		return redirectWithError(c, oauthState.ReturnTo, "authentication_failed")
	}

	token, err := h.googleOauthConfig.Exchange(c.Request().Context(), c.QueryParam("code"))
	if err != nil {
		return redirectWithError(c, oauthState.ReturnTo, "authentication_failed")
	}

	userInfo, err := h.ssoService.FetchGoogleUserInfo(token.AccessToken)
	if err != nil {
		return redirectWithError(c, oauthState.ReturnTo, "authentication_failed")
	}

	if err := h.ssoService.HandleReauth(user.ID, cookie.Value, userInfo); err != nil {
		code := "reauth_failed"
		if err.Error() == "identity mismatch" {
			code = "identity_mismatch"
		}
		return redirectWithError(c, oauthState.ReturnTo, code)
	}

	return redirectWithQuery(c, oauthState.ReturnTo, map[string]string{"reauth_success": "true"})
}

func redirectWithQuery(c *echo.Context, returnTo string, params map[string]string) error {
	redirectURL, err := url.Parse(returnTo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("server_error", "Invalid return URL"))
	}

	if params != nil {
		query := redirectURL.Query()
		for key, value := range params {
			query.Set(key, value)
		}
		redirectURL.RawQuery = query.Encode()
	}
	
	return c.Redirect(http.StatusFound, redirectURL.String())
}

func redirectWithError(c *echo.Context, returnTo, code string) error {
	return redirectWithQuery(c, returnTo, map[string]string{"error": code})
}
