package models

// ValidThemePreferences lists supported user theme preference values.
var ValidThemePreferences = map[string]bool{
	"system":      true,
	"light":       true,
	"dark":        true,
	"tokyo-night": true,
	"coffee":      true,
	"forest":      true,
	"rose":        true,
	"midnight":    true,
	"parchment":   true,
}

// IsValidThemePreference reports whether a theme preference string is supported.
func IsValidThemePreference(theme string) bool {
	return ValidThemePreferences[theme]
}
