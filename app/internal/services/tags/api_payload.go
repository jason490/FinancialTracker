package tags

import (
	"FinancialTracker/internal/models"
	"FinancialTracker/internal/models/external"
)

// BuildTagsPayload converts category and tag data into the external API payload.
func BuildTagsPayload(data []models.CategoryWithTags) *external.TagsPayload {
	cats := make([]external.CategoryWithTagsView, 0, len(data))
	for i := range data {
		cwt := &data[i]
		tagViews := make([]external.TagView, 0, len(cwt.Tags))
		for j := range cwt.Tags {
			tagViews = append(tagViews, external.TagView{
				ID:         cwt.Tags[j].ID,
				CategoryID: cwt.Tags[j].CategoryID,
				Name:       cwt.Tags[j].Name,
				Color:      cwt.Tags[j].Color,
			})
		}
		cats = append(cats, external.CategoryWithTagsView{
			ID:   cwt.ID,
			Name: cwt.Name,
			Tags: tagViews,
		})
	}
	return &external.TagsPayload{Categories: cats}
}

// BuildTagFilterViews converts tag filters into the external API payload.
func BuildTagFilterViews(filters []models.TagFilter) []external.TagFilterView {
	out := make([]external.TagFilterView, 0, len(filters))
	for i := range filters {
		out = append(out, external.TagFilterView{
			Pattern:    filters[i].Pattern,
			FilterType: filters[i].FilterType,
		})
	}
	return out
}
