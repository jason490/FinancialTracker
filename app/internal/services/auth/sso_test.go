package auth

import (
	"testing"
)

func TestBuildAndParseAPIOAuthState(t *testing.T) {
	t.Setenv("ENCRYPTION_KEY", "")

	svc := NewSSOService(nil)
	returnTo := "http://localhost/auth/sso/complete"

	state := svc.BuildAPIOAuthState(returnTo, "login", "")
	if state == apiStatePrefix+"error" {
		t.Fatalf("BuildAPIOAuthState() returned error state")
	}

	parsed, err := svc.ParseAPIState(state)
	if err != nil {
		t.Fatalf("ParseAPIState() error = %v", err)
	}
	if parsed.ReturnTo != returnTo {
		t.Fatalf("ParseAPIState() returnTo = %q, want %q", parsed.ReturnTo, returnTo)
	}
}

func TestParseAPIStateRejectsInvalidState(t *testing.T) {
	svc := NewSSOService(nil)

	if _, err := svc.ParseAPIState("state-token"); err != ErrInvalidOAuthState {
		t.Fatalf("ParseAPIState() error = %v, want ErrInvalidOAuthState", err)
	}
}
