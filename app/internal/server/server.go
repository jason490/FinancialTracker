package server

import (
	"FinancialTracker/internal/config"
	"FinancialTracker/internal/services/auth"
	"FinancialTracker/internal/services/mail"
	"FinancialTracker/internal/services/dashboard"
	"FinancialTracker/internal/services/financial"
	"FinancialTracker/internal/services/plaid"
	"FinancialTracker/internal/services/settings"
	"FinancialTracker/internal/services/stripebilling"
	"FinancialTracker/internal/services/subscription"
	"FinancialTracker/internal/services/tags"
	"FinancialTracker/internal/services/transactions"
	"FinancialTracker/internal/storage"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Server struct {
	store                *storage.Storage
	e                    *echo.Echo
	tagService           *tags.TaggingService
	transService         *transactions.TransactionService
	authService          *auth.AuthService
	ssoService           *auth.SSOService
	settingsService      *settings.SettingsService
	dashService          *dashboard.DashboardService
	plaidService         *plaid.PlaidService
	financialFacade      *financial.Facade
	subscriptionService  *subscription.Service
	billingService       *stripebilling.Service
	apiGoogleOauthConfig *oauth2.Config
}

func RunServer(store *storage.Storage) *echo.Echo {
	e := echo.New()
	corsOrigins := []string{
		"http://localhost",
		"capacitor://localhost",
	}
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		corsOrigins = append(corsOrigins, frontendURL)
	}
	if apiPublicURL := os.Getenv("API_PUBLIC_URL"); apiPublicURL != "" {
		corsOrigins = append(corsOrigins, apiPublicURL)
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderCookie, echo.HeaderXCSRFToken},
		AllowCredentials: true,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Secure()) // Adds X-Content-Type-Options, X-Frame-Options, X-XSS-Protection

	// CSRF middleware using double-submit cookie pattern.
	// The cookie is NOT HttpOnly so the SPA can read it and send the value
	// back as the X-CSRF-Token header on state-changing requests.
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		Skipper: func(c *echo.Context) bool {
			path := c.Path()
			return path == "/api/v1/plaid/webhook" || path == "/api/v1/stripe/webhook"
		},
		TokenLookup:    "header:" + echo.HeaderXCSRFToken + ",form:_csrf",
		CookieName:     "_csrf",
		CookiePath:     "/",
		CookieSecure:   !config.IsDevelopment(),
		CookieHTTPOnly: false, // SPA must read the cookie value
		CookieSameSite: http.SameSiteLaxMode,
		TrustedOrigins: corsOrigins,
	}))

	apiBaseURL := os.Getenv("API_PUBLIC_URL")
	if apiBaseURL == "" {
		apiBaseURL = "http://localhost:8080"
	}

	apiGoogleOauthConfig := &oauth2.Config{
		RedirectURL:  apiBaseURL + "/api/v1/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	tagService := tags.NewTaggingService(store)
	mailSender := mail.NewSenderFromEnv()
	authService := auth.NewAuthService(store, mailSender)
	ssoService := auth.NewSSOService(store)
	settingsService := settings.NewSettingsService(store)
	dashService := dashboard.NewDashboardService(store, tagService)
	subscriptionService := subscription.NewService(store, config.IsDevelopment())
	subscriptionService.SyncPrivilegeOverridesFromEnv()
	subscription.ParseUserIDOverrides(store)
	transService := transactions.NewTransactionService(store, subscriptionService)
	plaidService := plaid.NewPlaidService(store, tagService, subscriptionService)
	financialFacade := financial.NewFacade(store, tagService, subscriptionService)
	billingService := stripebilling.NewService(store, subscriptionService)
	server := Server{
		store:                store,
		e:                    e,
		tagService:           tagService,
		transService:         transService,
		authService:          authService,
		ssoService:           ssoService,
		settingsService:      settingsService,
		dashService:          dashService,
		plaidService:         plaidService,
		financialFacade:      financialFacade,
		subscriptionService:  subscriptionService,
		billingService:       billingService,
		apiGoogleOauthConfig: apiGoogleOauthConfig,
	}

	server.routes()
	return e
}
