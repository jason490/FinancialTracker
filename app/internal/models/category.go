package models

// Category groups tags for organization.
type Category struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}

// Tag is a user-defined label applied to transactions.
type Tag struct {
	ID         int64  `json:"id"`
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	CreatedAt  int64  `json:"created_at"`
}

// TagFilter is an auto-tagging rule for a tag.
type TagFilter struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	TagID      int64  `json:"tag_id"`
	Pattern    string `json:"pattern"`
	FilterType string `json:"filter_type"`
	CreatedAt  int64  `json:"created_at"`
}

// CategoryWithTags pairs a category with its tags.
type CategoryWithTags struct {
	Category
	Tags []Tag
}

// TagBreakdown is an aggregated amount for a tag (or uncategorized bucket).
type TagBreakdown struct {
	TagID   int64   `json:"tag_id"`
	TagName string  `json:"tag_name"`
	Color   string  `json:"color"`
	Total   float64 `json:"total"`
}
