package plaid

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/plaid/plaid-go/v42/plaid"
)

const (
	webhookTypeTransactions       = "TRANSACTIONS"
	webhookCodeSyncUpdatesAvailable = "SYNC_UPDATES_AVAILABLE"
	webhookTypeItem               = "ITEM"
	webhookCodeError              = "ERROR"
)

var (
	ErrWebhookInvalidPayload = errors.New("invalid webhook payload")
	itemSyncInProgress      sync.Map
)

type plaidWebhookEnvelope struct {
	WebhookType string `json:"webhook_type"`
	WebhookCode string `json:"webhook_code"`
	ItemID     string `json:"item_id"`
}

// HandleWebhook verifies a Plaid webhook and dispatches supported event types.
func (p *PlaidService) HandleWebhook(ctx context.Context, body []byte, verificationHeader string) error {
	if err := p.verifyWebhook(ctx, body, verificationHeader); err != nil {
		return err
	}

	var envelope plaidWebhookEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return ErrWebhookInvalidPayload
	}

	switch {
	case envelope.WebhookType == webhookTypeTransactions && envelope.WebhookCode == webhookCodeSyncUpdatesAvailable:
		if envelope.ItemID == "" {
			return ErrWebhookInvalidPayload
		}
		go p.syncTransactionsForItemWebhook(context.Background(), envelope.ItemID)
		return nil
	case envelope.WebhookType == webhookTypeItem && envelope.WebhookCode == webhookCodeError:
		if envelope.ItemID == "" {
			return ErrWebhookInvalidPayload
		}
		go p.handleItemErrorWebhook(context.Background(), body, envelope.ItemID)
		return nil
	default:
		log.Infof("Ignoring unhandled Plaid webhook: %s/%s", envelope.WebhookType, envelope.WebhookCode)
		return nil
	}
}

// syncTransactionsForItemWebhook pulls transaction updates for a single item after a webhook.
func (p *PlaidService) syncTransactionsForItemWebhook(ctx context.Context, plaidItemID string) {
	if _, loaded := itemSyncInProgress.LoadOrStore(plaidItemID, true); loaded {
		return
	}
	defer itemSyncInProgress.Delete(plaidItemID)

	item, err := p.store.GetPlaidItemByItemID(plaidItemID)
	if err != nil || item == nil {
		log.Infof("Plaid webhook received for unknown item %s", plaidItemID)
		return
	}
	if item.Status == ItemStatusDisconnected {
		return
	}

	if err := p.syncItemTransactions(ctx, item.UserID, item.PlaidItemID, item.AccessToken, item.SyncCursor); err != nil {
		log.Errorf("Webhook transaction sync failed for item %s: %v", plaidItemID, err)
		return
	}

	if err := p.store.UpdatePlaidItemLastSynced(plaidItemID, time.Now().Unix()); err != nil {
		log.Errorf("Failed to update last_synced after webhook for item %s: %v", plaidItemID, err)
	}
}

// handleItemErrorWebhook updates local item status when Plaid reports an item error.
func (p *PlaidService) handleItemErrorWebhook(ctx context.Context, body []byte, itemID string) {
	var webhook plaid.ItemErrorWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		log.Errorf("Failed to parse ITEM/ERROR webhook for %s: %v", itemID, err)
		return
	}

	itemErr, ok := webhook.GetErrorOk()
	if !ok || itemErr == nil {
		return
	}

	code := itemErr.GetErrorCode()
	reason := ""
	if r, ok := itemErr.GetErrorCodeReasonOk(); ok && r != nil {
		reason = *r
	}

	status := mapPlaidErrorToStatus(code, reason)
	if err := p.store.UpdatePlaidItemStatus(itemID, status, code); err != nil {
		log.Errorf("Failed to update item status from webhook for %s: %v", itemID, err)
		return
	}

	if isTerminalPlaidError(code) || status == ItemStatusDisconnected {
		if err := p.store.MarkPlaidAccountsDisconnectedByItemID(itemID); err != nil {
			log.Errorf("Failed to mark accounts disconnected from webhook for %s: %v", itemID, err)
		}
	}
}
