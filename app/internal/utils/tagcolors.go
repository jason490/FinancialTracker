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
var tagChipStyles = map[string]TagChipStyle{}

// TagChipStyle holds Tailwind classes for a tag pill.
type TagChipStyle struct {
	BG     string
	Text   string
	Border string
}

func init() {
	for _, c := range tagPalette {
		tagColorHex[c.Key] = c.Hex
		tagChipStyles[c.Key] = chipStyleForKey(c.Key)
	}
}

func chipStyleForKey(key string) TagChipStyle {
	switch key {
	case "blue":
		return TagChipStyle{"bg-blue-50 dark:bg-blue-900/20", "text-blue-600 dark:text-blue-400", "border-blue-200 dark:border-blue-800"}
	case "indigo":
		return TagChipStyle{"bg-indigo-50 dark:bg-indigo-900/20", "text-indigo-600 dark:text-indigo-400", "border-indigo-200 dark:border-indigo-800"}
	case "violet":
		return TagChipStyle{"bg-violet-50 dark:bg-violet-900/20", "text-violet-600 dark:text-violet-400", "border-violet-200 dark:border-violet-800"}
	case "emerald":
		return TagChipStyle{"bg-emerald-50 dark:bg-emerald-900/20", "text-emerald-600 dark:text-emerald-400", "border-emerald-200 dark:border-emerald-800"}
	case "teal":
		return TagChipStyle{"bg-teal-50 dark:bg-teal-900/20", "text-teal-600 dark:text-teal-400", "border-teal-200 dark:border-teal-800"}
	case "cyan":
		return TagChipStyle{"bg-cyan-50 dark:bg-cyan-900/20", "text-cyan-600 dark:text-cyan-400", "border-cyan-200 dark:border-cyan-800"}
	case "amber":
		return TagChipStyle{"bg-amber-50 dark:bg-amber-900/20", "text-amber-600 dark:text-amber-400", "border-amber-200 dark:border-amber-800"}
	case "orange":
		return TagChipStyle{"bg-orange-50 dark:bg-orange-900/20", "text-orange-600 dark:text-orange-400", "border-orange-200 dark:border-orange-800"}
	case "rose":
		return TagChipStyle{"bg-rose-50 dark:bg-rose-900/20", "text-rose-600 dark:text-rose-400", "border-rose-200 dark:border-rose-800"}
	case "red":
		return TagChipStyle{"bg-red-50 dark:bg-red-900/20", "text-red-600 dark:text-red-400", "border-red-200 dark:border-red-800"}
	default:
		return TagChipStyle{"bg-slate-50 dark:bg-slate-800/40", "text-slate-600 dark:text-slate-400", "border-slate-200 dark:border-slate-700"}
	}
}

// DefaultTagColor returns the default palette key for new tags.
func DefaultTagColor() string {
	return "blue"
}

// AllTagColors returns the full selectable palette.
func AllTagColors() []TagColorOption {
	out := make([]TagColorOption, len(tagPalette))
	copy(out, tagPalette)
	return out
}

// IsValidTagColor reports whether key is in the palette.
func IsValidTagColor(key string) bool {
	if key == "" {
		return false
	}
	_, ok := tagColorHex[key]
	return ok
}

// TagColorHex returns the hex color for a palette key, or slate if unknown.
func TagColorHex(key string) string {
	key = NormalizeTagColor(key)
	if hex, ok := tagColorHex[key]; ok {
		return hex
	}
	return tagColorHex["slate"]
}

// TagChipStyle returns Tailwind classes for tag pills.
func TagChipStyleForKey(key string) TagChipStyle {
	key = NormalizeTagColor(key)
	if s, ok := tagChipStyles[key]; ok {
		return s
	}
	return tagChipStyles["slate"]
}

// NormalizeTagColor returns a valid palette key, defaulting invalid values to blue.
func NormalizeTagColor(key string) string {
	if IsValidTagColor(key) {
		return key
	}
	return DefaultTagColor()
}
