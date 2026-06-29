package subscription

import (
	"testing"
	"time"
)

func TestBillingPeriodStart_signupAnchor(t *testing.T) {
	anchor := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC).Unix()
	now := time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC).Unix()
	got := BillingPeriodStart(anchor, now)
	want := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC).Unix()
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestBillingPeriodStart_beforeAnniversary(t *testing.T) {
	anchor := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC).Unix()
	now := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC).Unix()
	got := BillingPeriodStart(anchor, now)
	want := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC).Unix()
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestBillingPeriodStart_clampsEndOfMonth(t *testing.T) {
	anchor := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC).Unix()
	now := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC).Unix()
	got := BillingPeriodStart(anchor, now)
	want := time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC).Unix()
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestBillingPeriodEnd(t *testing.T) {
	anchor := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC).Unix()
	start := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC).Unix()
	got := BillingPeriodEnd(start, anchor)
	want := time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC).Unix()
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
