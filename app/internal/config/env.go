package config

import "os"

// IsDevelopment reports whether the app runs in local development mode.
func IsDevelopment() bool {
	return os.Getenv("ENV") == "development"
}

// EnvOr returns the environment variable value or a fallback when unset.
func EnvOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
