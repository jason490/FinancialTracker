package plaid

import "testing"

func TestBeginManualSyncRateLimit(t *testing.T) {
	svc := &PlaidService{}
	const userID int64 = 42

	if err := svc.beginManualSync(userID); err != nil {
		t.Fatalf("first beginManualSync() error = %v", err)
	}
	if err := svc.beginManualSync(userID); err != ErrPlaidSyncRateLimited {
		t.Fatalf("second beginManualSync() error = %v, want ErrPlaidSyncRateLimited", err)
	}
}
