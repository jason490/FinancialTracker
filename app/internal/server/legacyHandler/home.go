package legacyHandler

import (
	"FinancialTracker/internal/services/dashboard"
	"FinancialTracker/internal/storage"
	"FinancialTracker/web/templ/components/dashboardwidgets"
	"FinancialTracker/web/templ/pages"
	"net/http"

	"github.com/labstack/echo/v5"
)

// HomeHandler handles home and dashboard-related requests
type HomeHandler struct {
	store    *storage.Storage
	dashService *dashboard.DashboardService
}

// NewHomeHandler creates a new HomeHandler
func NewHomeHandler(store *storage.Storage, dashService *dashboard.DashboardService) *HomeHandler {
	return &HomeHandler{
		store:    store,
		dashService: dashService,
	}
}

// HandleHome renders the home page
func (h *HomeHandler) HandleHome(c *echo.Context) error {
	return Render(c, http.StatusOK, pages.Home(GetPageData(c, h.store, "Home")))
}

// HandleDashboard renders the dashboard page
func (h *HomeHandler) HandleDashboard(c *echo.Context) error {
	pageData := GetPageData(c, h.store, "Dashboard")
	if pageData.User == nil {
		return Render(c, http.StatusOK, pages.Dashboard(pageData))
	}

	dashData, err := h.dashService.GetDashboardData(pageData.User.ID, false)
	if err != nil {
		AddNotification(c, "Failed to load dashboard", "error")
		return Render(c, http.StatusInternalServerError, pages.Dashboard(pageData))
	}
	pageData.Data = dashData
	return Render(c, http.StatusOK, pages.Dashboard(pageData))
}

// HandleDashboardWidgets returns the widget grid partial for HTMX updates.
func (h *HomeHandler) HandleDashboardWidgets(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	editMode := c.QueryParam("edit") == "1"
	dashData, err := h.dashService.GetDashboardData(user.ID, editMode)
	if err != nil {
		AddNotification(c, "Failed to load dashboard widgets", "error")
		return c.NoContent(http.StatusInternalServerError)
	}
	return Render(c, http.StatusOK, dashboardwidgets.Grid(dashData))
}

// HandleSaveDashboardLayout persists the user's widget layout from customize mode.
func (h *HomeHandler) HandleSaveDashboardLayout(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	raw := c.FormValue("layout")
	dashData, err := h.dashService.SaveLayout(user.ID, raw)
	if err != nil {
		AddNotification(c, err.Error(), "error")
		return c.NoContent(http.StatusBadRequest)
	}

	AddNotification(c, "Dashboard layout saved", "success")
	return Render(c, http.StatusOK, dashboardwidgets.Grid(dashData))
}
