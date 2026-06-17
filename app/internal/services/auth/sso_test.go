package auth

import (
	"testing"
)

func TestBuildAndParseAPIOAuthState(t *testing.T) {
	t.Setenv("ENCRYPTION_KEY", "")

	svc := NewSSOService(nil)
	returnTo := "http://localhost/auth/sso/complete"

	state := svc.BuildAPIOAuthState(returnTo, "login")
	if state == apiStatePrefix+"error" {
		t.Fatalf("BuildAPIOAuthState() returned error state")
	}

	parsed, err := svc.ParseAPIReturnTo(state)
	if err != nil {
		t.Fatalf("ParseAPIReturnTo() error = %v", err)
	}
	if parsed != returnTo {
		t.Fatalf("ParseAPIReturnTo() = %q, want %q", parsed, returnTo)
	}
}

func TestParseAPIReturnToRejectsInvalidState(t *testing.T) {
	svc := NewSSOService(nil)

	if _, err := svc.ParseAPIReturnTo("state-token"); err != ErrInvalidOAuthState {
		t.Fatalf("ParseAPIReturnTo() error = %v, want ErrInvalidOAuthState", err)
	}
}
