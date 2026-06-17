package tags

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/storage"
	"FinancialTracker/internal/utils"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type TaggingService struct {
	store *storage.Storage
}

func NewTaggingService(store *storage.Storage) *TaggingService {
	return &TaggingService{
		store: store,
	}
}

// GetTagsData fetches all categories and tags for a user, structured for display
func (s *TaggingService) GetTagsData(userID int64) ([]models.CategoryWithTags, error) {
	categories, err := s.store.GetCategoriesByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	allTags, err := s.store.GetAllTagsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}

	var data []models.CategoryWithTags
	for _, cat := range categories {
		var catTags []models.Tag
		for _, t := range allTags {
			if t.CategoryID == cat.ID {
				catTags = append(catTags, t)
			}
		}
		data = append(data, models.CategoryWithTags{Category: cat, Tags: catTags})
	}

	return data, nil
}

// CreateCategory creates a new category for a user
func (s *TaggingService) CreateCategory(userID int64, name string) (int64, error) {
	if name == "" {
		return 0, errors.New("category name is required")
	}
	return s.store.CreateCategory(userID, name)
}

// UpdateCategory updates an existing category
func (s *TaggingService) UpdateCategory(userID, categoryID int64, name string) error {
	if name == "" {
		return errors.New("category name is required")
	}
	return s.store.UpdateCategory(userID, categoryID, name)
}

// DeleteCategory deletes a category and handles orphaned tags based on the specified action
func (s *TaggingService) DeleteCategory(userID, categoryID int64, action string, targetCategoryID int64) error {
	if action == "move_to_misc" {
		var err error
		targetCategoryID, err = s.store.GetOrCreateMiscCategory(userID)
		if err != nil {
			return fmt.Errorf("failed to create Misc category: %w", err)
		}
	} else if action == "move_to" && targetCategoryID == 0 {
		return errors.New("target category is required for move_to action")
	}

	return s.store.DeleteCategory(userID, categoryID, targetCategoryID)
}

// MergeCategories merges source category into target category
func (s *TaggingService) MergeCategories(userID, sourceID, targetID int64) error {
	if sourceID == targetID {
		return errors.New("cannot merge category into itself")
	}
	return s.store.MergeCategories(userID, sourceID, targetID)
}

// CreateTag handles duplicate checks, regex validation, and persistence for a new tag
func (s *TaggingService) CreateTag(userID int64, categoryID int64, name, color string, filterPatterns, filterTypes []string, applyToPast bool) error {
	if name == "" || categoryID == 0 {
		return errors.New("name and category are required")
	}

	color = utils.NormalizeTagColor(color)

	// Check for duplicate tag name
	existing, err := s.store.GetTagByUserIDAndName(userID, name)
	if err != nil {
		return fmt.Errorf("database error while checking for duplicates: %w", err)
	}
	if existing != nil {
		return errors.New("a tag with this name already exists")
	}

	tagID, err := s.store.CreateTag(userID, categoryID, name, color)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	filters, err := s.validateAndPrepareFilters(userID, tagID, filterPatterns, filterTypes)
	if err != nil {
		return err
	}

	if len(filters) > 0 {
		if err := s.store.BatchCreateTagFilters(userID, tagID, filters); err != nil {
			return fmt.Errorf("tag created, but failed to save new filters: %w", err)
		}
	}

	if applyToPast {
		if err := s.store.ApplyTagFiltersToPastTransactions(userID, tagID); err != nil {
			return fmt.Errorf("tag created, but failed to apply to past transactions: %w", err)
		}
	}

	return nil
}

// UpdateTag updates tag properties and its filters
func (s *TaggingService) UpdateTag(userID, tagID int64, name, color string, categoryID *int64, filterPatterns, filterTypes []string, applyToPast bool) error {
	if name == "" {
		return errors.New("tag name is required")
	}

	// Check for duplicate tag name
	existing, err := s.store.GetTagByUserIDAndName(userID, name)
	if err != nil {
		return fmt.Errorf("database error while checking for duplicates: %w", err)
	}
	if existing != nil && existing.ID != tagID {
		return errors.New("a tag with this name already exists")
	}

	color = utils.NormalizeTagColor(color)

	if err := s.store.UpdateTag(userID, tagID, name, color); err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	if categoryID != nil {
		if err := s.store.MoveTagToCategory(userID, tagID, *categoryID); err != nil {
			return fmt.Errorf("failed to update tag category: %w", err)
		}
	}

	filters, err := s.validateAndPrepareFilters(userID, tagID, filterPatterns, filterTypes)
	if err != nil {
		return err
	}

	if err := s.store.DeleteTagFiltersByTagID(userID, tagID); err != nil {
		return fmt.Errorf("failed to clear old filters: %w", err)
	}

	if len(filters) > 0 {
		if err := s.store.BatchCreateTagFilters(userID, tagID, filters); err != nil {
			return fmt.Errorf("failed to save new filters: %w", err)
		}
	}

	if applyToPast {
		if err := s.store.ApplyTagFiltersToPastTransactions(userID, tagID); err != nil {
			return fmt.Errorf("tag saved, but failed to apply to past transactions: %w", err)
		}
	}

	return nil
}

// DeleteTag removes a tag
func (s *TaggingService) DeleteTag(userID, tagID int64) error {
	return s.store.DeleteTag(userID, tagID)
}

// MoveTagToCategory moves a tag to a different category
func (s *TaggingService) MoveTagToCategory(userID, tagID, categoryID int64) error {
	return s.store.MoveTagToCategory(userID, tagID, categoryID)
}

// GetTagFilters fetches all filters for a specific tag
func (s *TaggingService) GetTagFilters(userID, tagID int64) ([]models.TagFilter, error) {
	filters, err := s.store.GetTagFiltersByTagID(userID, tagID)
	if err != nil {
		return nil, err
	}
	if filters == nil {
		return []models.TagFilter{}, nil
	}
	return filters, nil
}

// validateAndPrepareFilters performs regex validation and structures filter data
func (s *TaggingService) validateAndPrepareFilters(userID, tagID int64, patterns, types []string) ([]models.TagFilter, error) {
	var filters []models.TagFilter
	for i := range patterns {
		if patterns[i] == "" {
			continue
		}

		filterType := "string"
		if i < len(types) {
			filterType = types[i]
		}

		if filterType == "regex" {
			if _, err := regexp.Compile(patterns[i]); err != nil {
				return nil, fmt.Errorf("invalid regex pattern: %s", patterns[i])
			}
		}

		filters = append(filters, models.TagFilter{
			UserID:     userID,
			TagID:      tagID,
			Pattern:    patterns[i],
			FilterType: filterType,
		})
	}
	return filters, nil
}

// DefaultTag represents a default tag with its initial filter pattern
type DefaultTag struct {
	Name    string
	Pattern string
	Color   string
}

// DefaultCategories and their associated tags with default filter patterns
var DefaultCategories = map[string][]DefaultTag{
	"Food & Drink": {
		{Name: "Dining Out", Pattern: "FOOD_AND_DRINK", Color: "rose"},
		{Name: "Groceries", Pattern: "GROCERIES", Color: "emerald"},
	},
	"Transport": {
		{Name: "Automotive", Pattern: "TRANSPORTATION", Color: "slate"},
		{Name: "Travel", Pattern: "TRAVEL", Color: "indigo"},
	},
	"Shopping": {
		{Name: "General", Pattern: "GENERAL_MERCHANDISE", Color: "amber"},
		{Name: "Home", Pattern: "HOME_IMPROVEMENT", Color: "orange"},
	},
	"Recurring": {
		{Name: "Bills", Pattern: "RENT_AND_UTILITIES", Color: "red"},
		{Name: "Services", Pattern: "GENERAL_SERVICES", Color: "blue"},
	},
	"Health & Wellness": {
		{Name: "Medical", Pattern: "MEDICAL", Color: "cyan"},
		{Name: "Personal Care", Pattern: "PERSONAL_CARE", Color: "violet"},
	},
	"Financial": {
		{Name: "Income", Pattern: "INCOME", Color: "teal"},
		{Name: "Fees", Pattern: "BANK_FEES", Color: "red"},
		{Name: "Transfers", Pattern: "TRANSFER", Color: "slate"},
	},
	"Leisure": {
		{Name: "Entertainment", Pattern: "ENTERTAINMENT", Color: "violet"},
	},
}

// SeedDefaults creates default categories and tags for a new user if they don't exist
func (s *TaggingService) SeedDefaults(userID int64) error {
	existingCategories, err := s.store.GetCategoriesByUserID(userID)
	if err != nil {
		return err
	}

	if len(existingCategories) > 0 {
		return nil // Already seeded or user has custom setup
	}

	for catName, defaultTags := range DefaultCategories {
		catID, err := s.store.CreateCategory(userID, catName)
		if err != nil {
			continue
		}

		for _, dt := range defaultTags {
			tagID, err := s.store.CreateTag(userID, catID, dt.Name, dt.Color)
			if err == nil {
				s.store.CreateTagFilter(userID, tagID, dt.Pattern, "string")
			}
		}
	}

	return nil
}

// AutoTagTransaction attempts to apply tags to a transaction based on user filters
func (s *TaggingService) AutoTagTransaction(userID int64, t *models.Transaction) error {
	filters, err := s.store.GetTagFiltersByUserID(userID)
	if err != nil {
		return err
	}

	for _, filter := range filters {
		if s.matches(t, filter) {
			s.store.AddTagToTransaction(userID, t.ID, filter.TagID)
		}
	}

	return nil
}

func (s *TaggingService) matches(t *models.Transaction, filter models.TagFilter) bool {
	pattern := strings.ToLower(filter.Pattern)
	name := strings.ToLower(t.Name)
	merchant := strings.ToLower(t.MerchantName)
	category := strings.ToLower(t.PlaidCategory)

	switch filter.FilterType {
	case "string":
		return strings.Contains(name, pattern) || strings.Contains(merchant, pattern) || strings.Contains(category, pattern)
	case "regex":
		matched, _ := regexp.MatchString(filter.Pattern, t.Name)
		if matched {
			return true
		}
		matched, _ = regexp.MatchString(filter.Pattern, t.MerchantName)
		if matched {
			return true
		}
		matched, _ = regexp.MatchString(filter.Pattern, t.PlaidCategory)
		return matched
	case "amount_greater":
		var amt float64
		fmt.Sscanf(filter.Pattern, "%f", &amt)
		return t.Amount > amt
	case "amount_less":
		var amt float64
		fmt.Sscanf(filter.Pattern, "%f", &amt)
		return t.Amount < amt
	case "amount_equal":
		var amt float64
		fmt.Sscanf(filter.Pattern, "%f", &amt)
		return t.Amount == amt
	}
	return false
}
