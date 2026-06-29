package utils

// TagColorOption is a selectable tag color from the palette.
type TagColorOption struct {
	Key  string
	Hex  string
	Name string
}

var tagPalette = []TagColorOption{
	{Key: "blue", Hex: "#3b82f6", Name: "Blue"},
	{Key: "indigo", Hex: "#6366f1", Name: "Indigo"},
	{Key: "violet", Hex: "#8b5cf6", Name: "Violet"},
	{Key: "emerald", Hex: "#10b981", Name: "Emerald"},
	{Key: "teal", Hex: "#14b8a6", Name: "Teal"},
	{Key: "cyan", Hex: "#06b6d4", Name: "Cyan"},
	{Key: "amber", Hex: "#f59e0b", Name: "Amber"},
	{Key: "orange", Hex: "#f97316", Name: "Orange"},
	{Key: "rose", Hex: "#f43f5e", Name: "Rose"},
	{Key: "red", Hex: "#ef4444", Name: "Red"},
	{Key: "slate", Hex: "#64748b", Name: "Slate"},
}

var tagColorHex = map[string]string{}

func init() {
	for _, c := range tagPalette {
		tagColorHex[c.Key] = c.Hex
	}
}

// DefaultTagColor returns the default palette key for new tags.
func DefaultTagColor() string {
	return "blue"
}

// IsValidTagColor reports whether key is in the palette.
func IsValidTagColor(key string) bool {
	if key == "" {
		return false
	}
	_, ok := tagColorHex[key]
	return ok
}

// NormalizeTagColor returns a valid palette key, defaulting invalid values to blue.
func NormalizeTagColor(key string) string {
	if IsValidTagColor(key) {
		return key
	}
	return DefaultTagColor()
}
