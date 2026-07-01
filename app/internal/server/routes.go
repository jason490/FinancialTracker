package server

import (
	"FinancialTracker/internal/config"
	"FinancialTracker/internal/server/handler"
	"FinancialTracker/internal/server/middleware"
	"net/http"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

// routes registers all the routes for the application
func (s *Server) routes() {
	s.api()
}

// api registers JSON API routes for the SolidStart frontend
func (s *Server) api() {
	authHandler := handler.NewAuthHandler(s.authService, s.store)
	adminHandler := handler.NewAdminHandler(s.authService)
	ssoHandler := handler.NewSSOHandler(s.apiGoogleOauthConfig, s.ssoService)
	authMiddleware := middleware.NewAuthMiddleware(s.store, s.financialFacade)

	api := s.e.Group("/api/v1")

	// Auth Group with conditional rate limiting
	auth := api.Group("/auth")
	if !config.IsDevelopment() {
		auth.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(20)))
	}

	// CSRF token endpoint — GET triggers the middleware to set the _csrf cookie
	// and returns the token in the body so the SPA can read it on startup.
	api.GET("/csrf", func(c *echo.Context) error {
		token, _ := c.Get("csrf").(string)
		return c.JSON(http.StatusOK, map[string]string{"token": token})
	})

	auth.POST("/login", authHandler.HandleLogin)
	auth.GET("/registration-config", authHandler.HandleRegistrationConfig)
	auth.POST("/register", authHandler.HandleRegister)
	auth.POST("/forgot-password", authHandler.HandleForgotPassword)
	auth.POST("/verify-reset-code", authHandler.HandleVerifyResetCode)
	auth.POST("/reset-password", authHandler.HandleResetPassword)
	auth.GET("/google", ssoHandler.HandleGoogleLogin)
	auth.GET("/google/callback", ssoHandler.HandleGoogleCallback)

	session := api.Group("", authMiddleware.SessionMiddleware)
	session.GET("/auth/me", authHandler.HandleMe)

	protected := api.Group("", authMiddleware.AuthMiddlewareJSON)
	protected.POST("/auth/logout", authHandler.HandleLogout)
	protected.POST("/auth/onboarding/complete", authHandler.HandleCompleteOnboarding)
	protected.POST("/admin/registration-codes", adminHandler.HandleCreateRegistrationCode)

	dashboardHandler := handler.NewDashboardHandler(s.dashService)
	protected.GET("/dashboard", dashboardHandler.HandleDashboardGet)
	protected.POST("/dashboard/layout", dashboardHandler.HandleSaveLayout)

	transactionHandler := handler.NewTransactionHandler(s.transService)
	protected.GET("/transactions", transactionHandler.HandleGetTransactions)
	protected.GET("/transactions/export", transactionHandler.HandleExportTransactions)
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
	api.POST("/plaid/webhook", plaidHandler.HandleWebhook)

	connectionsHandler := handler.NewConnectionsHandler(s.financialFacade.Active())
	protected.GET("/connections/provider", connectionsHandler.HandleGetProvider)
	protected.GET("/connections", connectionsHandler.HandleGetConnections)
	protected.POST("/connections/create-session", connectionsHandler.HandleCreateSession)
	protected.POST("/connections/complete", connectionsHandler.HandleCompleteConnection)
	protected.POST("/connections/sync", connectionsHandler.HandleSync)
	protected.POST("/connections/create-update-session/:id", connectionsHandler.HandleCreateUpdateSession)
	protected.POST("/connections/sync-item/:id", connectionsHandler.HandleSyncItem)
	protected.POST("/connections/disconnect/:id", connectionsHandler.HandleDisconnect)
	protected.POST("/connections/toggle-visibility/:id", connectionsHandler.HandleToggleAccountVisibility)
	protected.POST("/connections/remove-account/:id", connectionsHandler.HandleRemoveAccount)

	subscriptionHandler := handler.NewSubscriptionHandler(s.subscriptionService, s.billingService)
	protected.GET("/subscription", subscriptionHandler.HandleGetSubscription)
	protected.POST("/subscription/change", subscriptionHandler.HandleChangeSubscription)
	protected.POST("/subscription/checkout", subscriptionHandler.HandleCreateCheckoutSession)
	protected.POST("/subscription/portal", subscriptionHandler.HandleCreateBillingPortal)
	api.POST("/stripe/webhook", subscriptionHandler.HandleStripeWebhook)

	s.e.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
}
