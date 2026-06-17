package server

import (
	"FinancialTracker/internal/server/handler"
	"FinancialTracker/internal/server/legacyHandler"
	"FinancialTracker/internal/server/middleware"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

// routes registers all the routes for the application
func (s *Server) routes() {
	s.legacy()
	s.api()
}

// api registers JSON API routes for the SolidStart frontend.
func (s *Server) api() {
	authHandler := handler.NewAuthHandler(s.authService, s.store)
	ssoHandler := handler.NewSSOHandler(s.apiGoogleOauthConfig, s.ssoService)
	authMiddleware := middleware.NewAuthMiddleware(s.store, s.plaidService)

	api := s.e.Group("/api/v1")

	// Auth Group with conditional rate limiting
	auth := api.Group("/auth")
	if os.Getenv("ENV") != "development" {
		auth.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(20)))
	}

	// CSRF token endpoint — GET triggers the middleware to set the _csrf cookie
	// and returns the token in the body so the SPA can read it on startup.
	api.GET("/csrf", func(c *echo.Context) error {
		token, _ := c.Get("csrf").(string)
		return c.JSON(http.StatusOK, map[string]string{"token": token})
	})

	auth.POST("/login", authHandler.HandleLogin)
	auth.POST("/register", authHandler.HandleRegister)
	auth.POST("/forgot-password", authHandler.HandleForgotPassword)
	auth.GET("/google", ssoHandler.HandleGoogleLogin)
	auth.GET("/google/callback", ssoHandler.HandleGoogleCallback)

	session := api.Group("", authMiddleware.SessionMiddleware)
	session.GET("/auth/me", authHandler.HandleMe)

	protected := api.Group("", authMiddleware.AuthMiddlewareJSON)
	protected.POST("/auth/logout", authHandler.HandleLogout)

	dashboardHandler := handler.NewDashboardHandler(s.dashService)
	protected.GET("/dashboard", dashboardHandler.HandleDashboardGet)
	protected.POST("/dashboard/layout", dashboardHandler.HandleSaveLayout)

	transactionHandler := handler.NewTransactionHandler(s.transService)
	protected.GET("/transactions", transactionHandler.HandleGetTransactions)
	protected.POST("/transactions/bulk-add-tag", transactionHandler.HandleBulkAddTag)
	protected.POST("/transactions/bulk-remove-tag", transactionHandler.HandleBulkRemoveTag)

	settingsHandler := handler.NewSettingsHandler(s.settingsService)
	protected.GET("/settings", settingsHandler.HandleGetSettings)
	protected.PATCH("/settings/profile", settingsHandler.HandleUpdateProfile)
	protected.POST("/settings/password", settingsHandler.HandleUpdatePassword)
	protected.POST("/settings/theme", settingsHandler.HandleUpdateTheme)
	protected.POST("/settings/unlink/:provider", settingsHandler.HandleUnlinkSSO)
	protected.GET("/settings/delete/reauth-status", settingsHandler.HandleDeleteAccountReauthStatus)
	protected.POST("/settings/delete/verify", settingsHandler.HandleDeleteAccountVerify)
	protected.POST("/settings/delete/confirm", settingsHandler.HandleDeleteAccountConfirm)

	tagHandler := handler.NewTagHandler(s.tagService)
	protected.GET("/tags", tagHandler.HandleGetTags)
	protected.GET("/tags/:id/filters", tagHandler.HandleGetTagFilters)
	protected.POST("/tags", tagHandler.HandleCreateTag)
	protected.PUT("/tags/:id", tagHandler.HandleUpdateTag)
	protected.DELETE("/tags/:id", tagHandler.HandleDeleteTag)
	protected.POST("/tags/:id/move", tagHandler.HandleMoveTag)
	protected.POST("/categories", tagHandler.HandleCreateCategory)
	protected.PUT("/categories/:id", tagHandler.HandleUpdateCategory)
	protected.DELETE("/categories/:id", tagHandler.HandleDeleteCategory)

	plaidHandler := handler.NewPlaidHandler(s.plaidService)
	protected.GET("/plaid/connections", plaidHandler.HandleGetConnections)
	protected.POST("/plaid/create-link-token", plaidHandler.HandleCreateLinkToken)
	protected.POST("/plaid/exchange", plaidHandler.HandleExchangeToken)
	protected.POST("/plaid/sync", plaidHandler.HandleSync)
	protected.POST("/plaid/create-update-token/:id", plaidHandler.HandleCreateUpdateLinkToken)
	protected.POST("/plaid/sync-item/:id", plaidHandler.HandleSyncItem)
	protected.POST("/plaid/disconnect/:id", plaidHandler.HandleDisconnectItem)
	protected.POST("/plaid/toggle-visibility/:id", plaidHandler.HandleToggleAccountVisibility)
	protected.POST("/plaid/remove-account/:id", plaidHandler.HandleRemoveAccount)
}

// legacy registers the legacy Templ/HTMX routes
func (s *Server) legacy() {
	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(s.store, s.plaidService)

	// Handlers
	homeHandler := legacyHandler.NewHomeHandler(s.store, s.dashService)
	authHandler := legacyHandler.NewAuthHandler(s.store, s.authService)
	ssoHandler := legacyHandler.NewSSOHandler(s.store, s.googleOauthConfig, s.ssoService)
	settingsHandler := legacyHandler.NewSettingsHandler(s.store, s.plaidService, s.settingsService)
	plaidHandler := legacyHandler.NewPlaidHandler(s.store, s.plaidService)
	transactionHandler := legacyHandler.NewTransactionHandler(s.store, s.transService)
	tagHandler := legacyHandler.NewTagHandler(s.store, s.tagService)

	// Static files
	s.e.Static("/static", "web/static")

	// Public Pages (With optional session info)
	public := s.e.Group("/api", authMiddleware.SessionMiddleware)
	public.GET("/", homeHandler.HandleHome)
	public.GET("/login", authHandler.HandleLoginPage)
	public.GET("/register", authHandler.HandleRegisterPage)
	public.POST("/auth/login", authHandler.HandleLoginSubmit)
	public.POST("/auth/register", authHandler.HandleRegisterSubmit)

	// SSO Routes
	public.GET("/auth/google", ssoHandler.HandleGoogleLogin)
	public.GET("/auth/google/callback", ssoHandler.HandleGoogleCallback)
	public.GET("/auth/apple", ssoHandler.HandleAppleLogin)

	// Protected Pages
	protected := s.e.Group("/api", authMiddleware.AuthMiddleware)
	protected.GET("/dashboard", homeHandler.HandleDashboard)
	protected.GET("/dashboard/widgets", homeHandler.HandleDashboardWidgets)
	protected.POST("/dashboard/layout", homeHandler.HandleSaveDashboardLayout)
	protected.GET("/transactions", transactionHandler.HandleTransactions)
	protected.GET("/transactions/list", transactionHandler.HandleGetList)
	protected.POST("/transactions/bulk-add-tag", transactionHandler.HandleBulkAddTag)
	protected.POST("/transactions/bulk-remove-tag", transactionHandler.HandleBulkRemoveTag)

	protected.GET("/tags", tagHandler.HandleTags)
	protected.GET("/tags/list", tagHandler.HandleTagList)
	protected.POST("/tags/create", tagHandler.HandleCreateTag)
	protected.PUT("/tags/:id", tagHandler.HandleUpdateTag)
	protected.DELETE("/tags/:id", tagHandler.HandleDeleteTag)
	protected.GET("/tags/:id/filters", tagHandler.HandleGetTagFilters)
	protected.POST("/tags/:tag_id/move/:category_id", tagHandler.HandleMoveTag)

	protected.POST("/categories/create", tagHandler.HandleCreateCategory)
	protected.GET("/categories/options", tagHandler.HandleCategoryOptions)
	protected.PUT("/categories/:id", tagHandler.HandleUpdateCategory)
	protected.DELETE("/categories/:id", tagHandler.HandleDeleteCategory)
	protected.POST("/categories/merge", tagHandler.HandleMergeCategories)

	protected.GET("/settings", settingsHandler.HandleSettingsPage)
	protected.GET("/manage", plaidHandler.HandleManagePage)
	protected.GET("/manage/list", plaidHandler.HandleGetConnectionList)
	protected.POST("/settings/theme", settingsHandler.HandleUpdateTheme)
	protected.POST("/settings/account", settingsHandler.HandleUpdateAccount)
	protected.POST("/settings/password", settingsHandler.HandleUpdatePassword)
	protected.GET("/settings/plaid/bank-list", settingsHandler.HandleUpdateBankList)
	protected.POST("/settings/unlink/:provider", settingsHandler.HandleUnlinkSSO)
	protected.POST("/settings/plaid/remove/:id", settingsHandler.HandleRemoveBankAccount)
	protected.GET("/settings/delete/init", settingsHandler.HandleDeleteAccountInit)
	protected.GET("/settings/delete/check-reauth", settingsHandler.HandleDeleteAccountCheckReauth)
	protected.POST("/settings/delete/verify", settingsHandler.HandleDeleteAccountVerify)
	protected.POST("/settings/delete/confirm", settingsHandler.HandleDeleteAccountConfirm)
	protected.GET("/settings/delete/cancel", settingsHandler.HandleDeleteAccountCancel)
	protected.POST("/logout", authHandler.HandleLogout)

	// Protected and Plaid Routes
	plaidRoutes := protected.Group("/plaid")
	plaidRoutes.POST("/create-link-token", plaidHandler.HandleCreateLinkToken)
	plaidRoutes.POST("/exchange", plaidHandler.HandleExchangeToken)
	plaidRoutes.POST("/sync", plaidHandler.HandleSync)
	plaidRoutes.POST("/sync-item/:id", plaidHandler.HandleSyncItem)
	plaidRoutes.POST("/create-update-token/:id", plaidHandler.HandleCreateUpdateLinkToken)
	plaidRoutes.POST("/remove-account/:id", plaidHandler.HandleRemoveAccount)
	plaidRoutes.POST("/toggle-visibility/:id", plaidHandler.HandleToggleAccountVisibility)

	// Health check
	s.e.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
}

