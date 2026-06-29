package plaid

import (
	"errors"

	"github.com/plaid/plaid-go/v42/plaid"
)

// Plaid item status values stored in plaid_items.status.
const (
	ItemStatusActive        = "active"
	ItemStatusNeedsReauth   = "needs_reauth"
	ItemStatusDisconnected  = "disconnected"
	ItemStatusError         = "error"
)

// Sentinel errors for subscription and usage limits.
var (
	ErrPlaidAPILimitExceeded  = errors.New("monthly Plaid API limit reached for your plan")
	ErrPlaidItemLimitExceeded = errors.New("bank connection limit reached for your plan")
	ErrPlaidSyncRateLimited   = errors.New("You can sync once per minute. Bank data is often a few hours to a day behind; waiting won't pull updates faster than Plaid provides.")
)

// parsePlaidError extracts the Plaid error_code from an API error response.
func parsePlaidError(err error) (code string, ok bool) {
	plaidErr, convErr := plaid.ToPlaidError(err)
	if convErr != nil {
		return "", false
	}
	return plaidErr.GetErrorCode(), true
}

// mapPlaidErrorToStatus maps a Plaid error_code (and optional reason) to a local item status.
func mapPlaidErrorToStatus(errorCode, errorCodeReason string) string {
	switch errorCode {
	case "ITEM_NOT_FOUND", "USER_PERMISSION_REVOKED", "INVALID_ACCESS_TOKEN":
		return ItemStatusDisconnected
	case "ITEM_LOGIN_REQUIRED":
		return ItemStatusNeedsReauth
	default:
		if errorCodeReason == "OAUTH_USER_REVOKED" {
			return ItemStatusDisconnected
		}
		return ItemStatusError
	}
}

// isTerminalPlaidError reports whether the Item can no longer be synced without user re-linking.
func isTerminalPlaidError(errorCode string) bool {
	switch errorCode {
	case "ITEM_NOT_FOUND", "USER_PERMISSION_REVOKED", "INVALID_ACCESS_TOKEN":
		return true
	default:
		return false
	}
}

// itemErrorFromGet returns error_code and reason from an ItemGet item.error field, if present.
func itemErrorFromGet(item plaid.ItemWithConsentFields) (code string, reason string, hasError bool) {
	itemErr, ok := item.GetErrorOk()
	if !ok || itemErr == nil {
		return "", "", false
	}
	code = itemErr.GetErrorCode()
	if r, ok := itemErr.GetErrorCodeReasonOk(); ok && r != nil {
		reason = *r
	}
	return code, reason, true
}
