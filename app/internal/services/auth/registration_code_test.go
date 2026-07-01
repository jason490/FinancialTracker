package auth

import (
	"testing"
)

func TestIsRegistrationAdmin(t *testing.T) {
	t.Setenv("ENV", "development")
	t.Setenv("REGISTRATION_ADMIN_EMAILS", "")

	if !IsRegistrationAdmin("test@test.com") {
		t.Fatal("expected seeded development user to be registration admin")
	}

	t.Setenv("REGISTRATION_ADMIN_EMAILS", "admin@example.com;ops@example.com")

	if !IsRegistrationAdmin("admin@example.com") {
		t.Fatal("expected admin@example.com to be admin")
	}
	if !IsRegistrationAdmin("OPS@example.com") {
		t.Fatal("expected case-insensitive admin match")
	}
	if IsRegistrationAdmin("user@example.com") {
		t.Fatal("expected user@example.com to not be admin")
	}
}

func TestRegistrationGateEnabled(t *testing.T) {
	t.Setenv("SUBSCRIPTIONS_ENABLED", "false")
	if !RegistrationGateEnabled() {
		t.Fatal("expected registration gate when subscriptions disabled")
	}

	t.Setenv("SUBSCRIPTIONS_ENABLED", "true")
	if RegistrationGateEnabled() {
		t.Fatal("expected registration gate disabled when subscriptions enabled")
	}
}
