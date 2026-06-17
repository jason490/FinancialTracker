package utils

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
var specialCharRegex = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)

// Sanitize removes leading/trailing whitespace
func Sanitize(s string) string {
	return strings.TrimSpace(s)
}

// ValidateEmail checks if an email is valid and returns an error message if not.
func ValidateEmail(email string) (string, bool) {
	if !emailRegex.MatchString(email) {
		return "Invalid email format", true
	}
	return "", false
}

// ValidatePassword checks if a password meets complexity requirements.
func ValidatePassword(password string) (string, bool) {
	if len(password) < 8 || len(password) > 30 {
		return "Password must be between 8 and 30 characters", true
	}
	if !specialCharRegex.MatchString(password) {
		return "Password must contain at least one special character", true
	}
	return "", false
}

// ValidateName checks if a name is non-empty and within character limits.
func ValidateName(name, label string) (string, bool) {
	if len(name) == 0 {
		return label + " is required", true
	}
	if len(name) > 30 {
		return label + " must be 30 characters or less", true
	}
	return "", false
}
