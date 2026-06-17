package legacyHandler

import (
	"FinancialTracker/internal/services/auth"
	"FinancialTracker/internal/storage"
	"net/http"

	"github.com/labstack/echo/v5"
	"golang.org/x/oauth2"
)

// SSOHandler handles SSO-related requests
type SSOHandler struct {
	store             *storage.Storage
	googleOauthConfig *oauth2.Config
	ssoService        *auth.SSOService
}

// NewSSOHandler creates a new SSOHandler
func NewSSOHandler(store *storage.Storage, googleOauthConfig *oauth2.Config, ssoService *auth.SSOService) *SSOHandler {
	return &SSOHandler{
		store:             store,
		googleOauthConfig: googleOauthConfig,
		ssoService:        ssoService,
	}
}

// HandleGoogleLogin redirects to Google's consent page
func (h *SSOHandler) HandleGoogleLogin(c *echo.Context) error {
	state := "state-token"
	if c.QueryParam("reauth") == "delete" {
		state = "reauth-delete"
	} else if c.QueryParam("action") == "link" {
		state = "link-google"
	}
	url := h.googleOauthConfig.AuthCodeURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleGoogleCallback is the HTTP entry point for Google's OAuth redirect.
// OAuth exchange and user lookup are delegated to SSOService; this handler
// only routes by OAuth state and issues redirects or session cookies.
func (h *SSOHandler) HandleGoogleCallback(c *echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	switch state {
	case "reauth-delete":
		return h.handleGoogleReauthDelete(c, code)
	case "link-google":
		return h.handleGoogleLink(c, code)
	default:
		return h.handleGoogleLogin(c, code)
	}
}

func (h *SSOHandler) handleGoogleReauthDelete(c *echo.Context, code string) error {
	cookie, err := c.Cookie("Session")
	if err != nil {
		return c.Redirect(http.StatusFound, "/login?error=6")
	}

	currentUser := getCurrentUser(c, h.store)
	if currentUser == nil {
		return c.Redirect(http.StatusFound, "/login?error=6")
	}

	userInfo, err := h.fetchGoogleUserFromCode(c, code)
	if err != nil {
		return c.Redirect(http.StatusFound, "/login?error=1")
	}

	if err := h.ssoService.HandleReauth(currentUser.ID, cookie.Value, userInfo); err != nil {
		if err.Error() == "identity mismatch" {
			return c.Redirect(http.StatusFound, "/settings?error=4")
		}
		return c.Redirect(http.StatusFound, "/settings?error=5")
	}

	return c.Redirect(http.StatusFound, "/settings?reauth_success=true")
}

func (h *SSOHandler) handleGoogleLink(c *echo.Context, code string) error {
	currentUser := getCurrentUser(c, h.store)
	if currentUser == nil {
		return c.Redirect(http.StatusFound, "/login?error=6")
	}

	userInfo, err := h.fetchGoogleUserFromCode(c, code)
	if err != nil {
		return c.Redirect(http.StatusFound, "/login?error=1")
	}

	if err := h.ssoService.LinkGoogleAccount(currentUser.ID, userInfo); err != nil {
		if err.Error() == "SSO account already linked to another user" {
			return c.Redirect(http.StatusFound, "/settings?error=2")
		}
		return c.Redirect(http.StatusFound, "/settings?error=3")
	}

	return c.Redirect(http.StatusFound, "/settings?success=linked")
}

func (h *SSOHandler) handleGoogleLogin(c *echo.Context, code string) error {
	session, err := h.ssoService.CompleteGoogleLogin(c.Request().Context(), h.googleOauthConfig, code)
	if err != nil {
		return c.Redirect(http.StatusFound, "/login?error=5")
	}

	c.SetCookie(createCookie(session.ID, false))
	return c.Redirect(http.StatusFound, "/dashboard")
}

func (h *SSOHandler) fetchGoogleUserFromCode(c *echo.Context, code string) (auth.GoogleUserInfo, error) {
	token, err := h.googleOauthConfig.Exchange(c.Request().Context(), code)
	if err != nil {
		return auth.GoogleUserInfo{}, err
	}

	return h.ssoService.FetchGoogleUserInfo(token.AccessToken)
}

// HandleAppleLogin (Stub for Apple SSO)
func (h *SSOHandler) HandleAppleLogin(c *echo.Context) error {
	return c.String(http.StatusNotImplemented, "Apple SSO not implemented yet")
}
