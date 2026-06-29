package plaid

import (
	"sync"
	"time"
)

const manualSyncCooldown = time.Minute

var lastManualSync sync.Map

// beginManualSync enforces the once-per-minute limit for user-initiated sync actions.
func (p *PlaidService) beginManualSync(userID int64) error {
	now := time.Now().Unix()
	if v, ok := lastManualSync.Load(userID); ok {
		if last, ok := v.(int64); ok && now-last < int64(manualSyncCooldown.Seconds()) {
			return ErrPlaidSyncRateLimited
		}
	}
	lastManualSync.Store(userID, now)
	return nil
}
