package handler

import (
	"FinancialTracker/internal/models/external"
	"FinancialTracker/internal/services/dashboard"
	"net/http"

	"github.com/labstack/echo/v5"
)

// DashboardHandler serves JSON dashboard endpoints for the SPA.
type DashboardHandler struct {
	dashService *dashboard.DashboardService
}

// NewDashboardHandler creates a DashboardHandler.
func NewDashboardHandler(dashService *dashboard.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashService: dashService}
}

// HandleDashboardGet returns the dashboard payload for the authenticated user.
func (h *DashboardHandler) HandleDashboardGet(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	editMode := c.QueryParam("edit") == "1"
	payload, err := h.dashService.GetDashboardPayload(userID, editMode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("dashboard_error", "Failed to load dashboard"))
	}

	return c.JSON(http.StatusOK, payload)
}

// HandleSaveLayout persists a customized dashboard layout for a specific device.
func (h *DashboardHandler) HandleSaveLayout(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	var req external.SaveDashboardLayoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	if req.DeviceType != "desktop" && req.DeviceType != "mobile" {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_device", "Device type must be 'desktop' or 'mobile'"))
	}
	if len(req.Widgets) == 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_layout", "Layout must include widgets"))
	}

	payload, err := h.dashService.SaveDeviceLayout(userID, req.DeviceType, req.Widgets)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("layout_error", err.Error()))
	}

	return c.JSON(http.StatusOK, payload)
}
