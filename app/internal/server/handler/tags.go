package handler

import (
	"FinancialTracker/internal/models/external"
	"FinancialTracker/internal/services/tags"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

// TagHandler serves JSON tag and category endpoints for the SPA.
type TagHandler struct {
	tagService *tags.TaggingService
}

// NewTagHandler creates a TagHandler.
func NewTagHandler(tagService *tags.TaggingService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

// HandleGetTags returns all categories and tags for the authenticated user.
func (h *TagHandler) HandleGetTags(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	data, err := h.tagService.GetTagsData(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("tags_error", "Failed to load tags"))
	}

	return c.JSON(http.StatusOK, external.ToTagsPayload(data))
}

// HandleGetTagFilters returns auto-tagging filters for a specific tag.
func (h *TagHandler) HandleGetTagFilters(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	tagID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || tagID == 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid tag ID"))
	}

	filters, err := h.tagService.GetTagFilters(userID, tagID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("tags_error", "Failed to load tag filters"))
	}

	return c.JSON(http.StatusOK, external.ToTagFilterViews(filters))
}

// HandleCreateTag creates a new tag with optional filters.
func (h *TagHandler) HandleCreateTag(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	var req external.CreateTagRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	patterns, types := tagFiltersFromInput(req.Filters)
	if err := h.tagService.CreateTag(userID, req.CategoryID, req.Name, req.Color, patterns, types, req.Apply); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("validation_error", err.Error()))
	}

	return h.refreshTags(c, userID, http.StatusCreated)
}

// HandleUpdateTag updates an existing tag and its filters.
func (h *TagHandler) HandleUpdateTag(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	tagID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || tagID == 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid tag ID"))
	}

	var req external.UpdateTagRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	patterns, types := tagFiltersFromInput(req.Filters)
	var categoryID *int64
	if req.CategoryID > 0 {
		categoryID = &req.CategoryID
	}

	if err := h.tagService.UpdateTag(userID, tagID, req.Name, req.Color, categoryID, patterns, types, req.Apply); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("validation_error", err.Error()))
	}

	return h.refreshTags(c, userID, http.StatusOK)
}

// HandleDeleteTag removes a tag.
func (h *TagHandler) HandleDeleteTag(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	tagID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || tagID == 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid tag ID"))
	}

	if err := h.tagService.DeleteTag(userID, tagID); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("tags_error", "Failed to delete tag"))
	}

	return h.refreshTags(c, userID, http.StatusOK)
}

// HandleMoveTag moves a tag to another category.
func (h *TagHandler) HandleMoveTag(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	tagID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || tagID == 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid tag ID"))
	}

	var req external.MoveTagRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	if err := h.tagService.MoveTagToCategory(userID, tagID, req.CategoryID); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("tags_error", "Failed to move tag"))
	}

	return h.refreshTags(c, userID, http.StatusOK)
}

// HandleCreateCategory creates a new tag category.
func (h *TagHandler) HandleCreateCategory(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	var req external.CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	if _, err := h.tagService.CreateCategory(userID, req.Name); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("validation_error", err.Error()))
	}

	return h.refreshTags(c, userID, http.StatusCreated)
}

// HandleUpdateCategory renames a category.
func (h *TagHandler) HandleUpdateCategory(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || categoryID == 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid category ID"))
	}

	var req external.UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	if err := h.tagService.UpdateCategory(userID, categoryID, req.Name); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("validation_error", err.Error()))
	}

	return h.refreshTags(c, userID, http.StatusOK)
}

// HandleDeleteCategory deletes a category and handles its tags per the requested action.
func (h *TagHandler) HandleDeleteCategory(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int64)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, ErrorResponse("unauthorized", "Not authenticated"))
	}

	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || categoryID == 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid category ID"))
	}

	var req external.DeleteCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("invalid_request", "Invalid request body"))
	}

	if err := h.tagService.DeleteCategory(userID, categoryID, req.Action, req.TargetCategoryID); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse("validation_error", err.Error()))
	}

	return h.refreshTags(c, userID, http.StatusOK)
}

// refreshTags reloads and returns the full tags payload after a mutation.
func (h *TagHandler) refreshTags(c *echo.Context, userID int64, status int) error {
	data, err := h.tagService.GetTagsData(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse("tags_error", "Failed to load tags"))
	}
	return c.JSON(status, external.ToTagsPayload(data))
}

// tagFiltersFromInput converts request filter inputs into parallel pattern/type slices.
func tagFiltersFromInput(filters []external.TagFilterInput) (patterns, types []string) {
	for i := range filters {
		if filters[i].Pattern == "" {
			continue
		}
		patterns = append(patterns, filters[i].Pattern)
		filterType := filters[i].FilterType
		if filterType == "" {
			filterType = "string"
		}
		types = append(types, filterType)
	}
	return patterns, types
}
