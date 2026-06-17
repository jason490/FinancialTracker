package legacyHandler

import (
	"FinancialTracker/internal/services/tags"
	"FinancialTracker/internal/storage"
	"FinancialTracker/web/templ/components"
	"FinancialTracker/web/templ/pages"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/labstack/gommon/log"
)

type TagHandler struct {
	store      *storage.Storage
	tagService *tags.TaggingService
}

func NewTagHandler(store *storage.Storage, tagService *tags.TaggingService) *TagHandler {
	return &TagHandler{
		store:      store,
		tagService: tagService,
	}
}

func (h *TagHandler) HandleTags(c *echo.Context) error {
	pageData := GetPageData(c, h.store, "Tags")
	if pageData.User == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	data, err := h.tagService.GetTagsData(pageData.User.ID)
	if err != nil {
		log.Errorf("Failed to fetch tags data: %v", err)
		return Render(c, http.StatusInternalServerError, components.StatusMessage("Failed to load categories and tags", true))
	}

	pageData.Data = data
	return Render(c, http.StatusOK, pages.Tags(pageData))
}

func (h *TagHandler) HandleTagList(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	data, err := h.tagService.GetTagsData(user.ID)
	if err != nil {
		log.Errorf("Failed to fetch tags data: %v", err)
		return Render(c, http.StatusInternalServerError, components.StatusMessage("Failed to load categories and tags", true))
	}

	return Render(c, http.StatusOK, pages.CategoryGrid(data))
}

func (h *TagHandler) HandleCategoryOptions(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	categories, err := h.store.GetCategoriesByUserID(user.ID)
	if err != nil {
		log.Errorf("Failed to fetch categories: %v", err)
		return Render(c, http.StatusInternalServerError, components.StatusMessage("Failed to load categories", true))
	}
	includeAlpine := c.QueryParam("alpine") == "true"

	return Render(c, http.StatusOK, pages.CategoryOptions(categories, includeAlpine))
}

func (h *TagHandler) HandleCreateCategory(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	name := c.FormValue("name")
	if _, err := h.tagService.CreateCategory(user.ID, name); err != nil {
		return Render(c, http.StatusBadRequest, components.StatusMessage(err.Error(), true))
	}

	AddNotification(c, "Category created successfully", "success")
	c.Response().Header().Set("HX-Trigger", "updateTagList")
	return c.NoContent(http.StatusCreated)
}

func (h *TagHandler) HandleCreateTag(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	categoryIDStr := c.FormValue("category_id")
	categoryID, _ := strconv.ParseInt(categoryIDStr, 10, 64)
	name := c.FormValue("name")
	color := c.FormValue("color")
	applyToPast := c.FormValue("apply") == "true"

	if err := c.Request().ParseForm(); err != nil {
		return Render(c, http.StatusBadRequest, components.StatusMessage("Failed to parse form", true))
	}
	patterns := c.Request().Form["filter_pattern[]"]
	types := c.Request().Form["filter_type[]"]

	if err := h.tagService.CreateTag(user.ID, categoryID, name, color, patterns, types, applyToPast); err != nil {
		return Render(c, http.StatusBadRequest, components.StatusMessage(err.Error(), true))
	}

	AddNotification(c, "Tag created successfully", "success")
	c.Response().Header().Set("HX-Trigger", "updateTagList")
	return c.NoContent(http.StatusCreated)
}

func (h *TagHandler) HandleDeleteTag(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tagIDStr := c.Param("id")
	tagID, _ := strconv.ParseInt(tagIDStr, 10, 64)
	if err := h.tagService.DeleteTag(user.ID, tagID); err != nil {
		AddNotification(c, "Failed to delete tag", "error")
		return c.NoContent(http.StatusInternalServerError)
	}

	AddNotification(c, "Tag deleted", "success")
	c.Response().Header().Set("HX-Trigger", "updateTagList")
	return c.NoContent(http.StatusOK)
}

func (h *TagHandler) HandleUpdateCategory(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	categoryIDStr := c.Param("id")
	categoryID, _ := strconv.ParseInt(categoryIDStr, 10, 64)
	name := c.FormValue("name")

	if err := h.tagService.UpdateCategory(user.ID, categoryID, name); err != nil {
		AddNotification(c, err.Error(), "error")
		return c.NoContent(http.StatusBadRequest)
	}

	AddNotification(c, "Category updated", "success")
	c.Response().Header().Set("HX-Trigger", "updateTagList")
	return c.NoContent(http.StatusOK)
}

func (h *TagHandler) HandleDeleteCategory(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	categoryIDStr := c.Param("id")
	categoryID, _ := strconv.ParseInt(categoryIDStr, 10, 64)

	action := c.FormValue("action") // "delete_all", "move_to", "move_to_misc"
	targetCategoryID := int64(0)
	if action == "move_to" {
		targetCategoryIDStr := c.FormValue("target_category_id")
		targetCategoryID, _ = strconv.ParseInt(targetCategoryIDStr, 10, 64)
	}

	if err := h.tagService.DeleteCategory(user.ID, categoryID, action, targetCategoryID); err != nil {
		return Render(c, http.StatusInternalServerError, components.StatusMessage(err.Error(), true))
	}

	AddNotification(c, "Category deleted", "success")
	c.Response().Header().Set("HX-Trigger", "updateTagList")
	return c.NoContent(http.StatusOK)
}

func (h *TagHandler) HandleMoveTag(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tagIDStr := c.Param("tag_id")
	tagID, _ := strconv.ParseInt(tagIDStr, 10, 64)
	categoryIDStr := c.Param("category_id")
	categoryID, _ := strconv.ParseInt(categoryIDStr, 10, 64)

	if err := h.tagService.MoveTagToCategory(user.ID, tagID, categoryID); err != nil {
		AddNotification(c, "Failed to move tag", "error")
		return c.NoContent(http.StatusInternalServerError)
	}

	AddNotification(c, "Tag moved successfully", "success")
	c.Response().Header().Set("HX-Trigger", "updateTagList")
	return c.NoContent(http.StatusOK)
}

func (h *TagHandler) HandleMergeCategories(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	sourceIDStr := c.FormValue("source_id")
	sourceID, _ := strconv.ParseInt(sourceIDStr, 10, 64)
	targetIDStr := c.FormValue("target_id")
	targetID, _ := strconv.ParseInt(targetIDStr, 10, 64)

	if err := h.tagService.MergeCategories(user.ID, sourceID, targetID); err != nil {
		return Render(c, http.StatusBadRequest, components.StatusMessage(err.Error(), true))
	}

	AddNotification(c, "Categories merged", "success")
	c.Response().Header().Set("HX-Trigger", "updateTagList")
	return c.NoContent(http.StatusOK)
}

func (h *TagHandler) HandleUpdateTag(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tagIDStr := c.Param("id")
	tagID, _ := strconv.ParseInt(tagIDStr, 10, 64)
	name := c.FormValue("name")
	color := c.FormValue("color")
	applyToPast := c.FormValue("apply") == "true"

	var categoryID *int64
	categoryIDStr := c.FormValue("category_id")
	if categoryIDStr != "" {
		id, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err == nil {
			categoryID = &id
		}
	}

	if err := c.Request().ParseForm(); err != nil {
		return Render(c, http.StatusBadRequest, components.StatusMessage("Failed to parse form", true))
	}
	patterns := c.Request().Form["filter_pattern[]"]
	types := c.Request().Form["filter_type[]"]

	if err := h.tagService.UpdateTag(user.ID, tagID, name, color, categoryID, patterns, types, applyToPast); err != nil {
		return Render(c, http.StatusBadRequest, components.StatusMessage(err.Error(), true))
	}

	AddNotification(c, "Tag updated successfully", "success")
	c.Response().Header().Set("HX-Trigger", "updateTagList")
	return c.NoContent(http.StatusOK)
}

func (h *TagHandler) HandleGetTagFilters(c *echo.Context) error {
	user := getCurrentUser(c, h.store)
	if user == nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tagIDStr := c.Param("id")
	tagID, _ := strconv.ParseInt(tagIDStr, 10, 64)

	filters, err := h.tagService.GetTagFilters(user.ID, tagID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch filters"})
	}

	return c.JSON(http.StatusOK, filters)
}
