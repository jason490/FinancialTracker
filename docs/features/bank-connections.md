# Bank Connections

## Overview

Bank connections link external accounts so FinancialTracker can sync balances and transactions. The active provider is `FINANCIAL_PROVIDER`: Plaid (default) or Stripe Financial Connections. The SPA uses provider-agnostic `/api/v1/connections/*` routes. Provider data stays in separate tables. The app does not merge Plaid and Stripe FC history.

## How it works

### Provider selection

| `FINANCIAL_PROVIDER` | Client flow | Storage |
|----------------------|-------------|---------|
| `plaid` (default) | Plaid Link | `plaid_items`, `plaid_account`, `plaid_api_usage` |
| `stripe` | Stripe.js `collectFinancialConnectionsAccounts` | `stripe_fc_items`, `stripe_fc_account`, `stripe_api_usage` |

`GET /api/v1/connections/provider` returns the active provider name and Stripe publishable key when needed. Frontend entry points are `frontend/src/lib/connections.ts` and `frontend/src/lib/stripe.ts`. Plaid and Stripe SDKs load only when a connection flow starts.

### Connect a bank

1. Open Settings → Connections or the onboarding connect step.
2. Call `POST /api/v1/connections/create-session` to start a link session.
3. Complete Plaid Link or Stripe Financial Connections in the browser.
4. Call `POST /api/v1/connections/complete` with the provider payload.
5. The service stores the item and accounts, then syncs transactions for the active provider.

Subscription limits can block new items (`connection_item_limit`, HTTP 403) or metered API calls (`connection_api_limit`, HTTP 429).

### Sync

Users trigger sync from Dashboard Quick Actions, the Transactions header, or Connections.

- `POST /api/v1/connections/sync` syncs all items for the user.
- `POST /api/v1/connections/sync-item/:id` syncs one item.

Plaid account and transaction sync do not consume the monthly API quota. Manual Plaid sync is rate-limited to once per minute per user (`connection_sync_rate_limit`, HTTP 429). Webhooks and background stale sync are exempt. Stripe FC metering stays on Stripe FC usage tables.

### Webhooks

Plaid transaction webhooks use `POST /api/v1/plaid/webhook` (unauthenticated, CSRF-exempt). The API verifies the `Plaid-Verification` JWT.

| Webhook | Behavior |
|---------|----------|
| `TRANSACTIONS` / `SYNC_UPDATES_AVAILABLE` | Background sync for the item |
| `ITEM` / `ERROR` | Update connection health |

Link and update tokens register the webhook URL from `PLAID_WEBHOOK_URL` or `{API_PUBLIC_URL}/api/v1/plaid/webhook`.

Stripe billing webhooks use a different path. See [Subscriptions](./subscriptions.md).

### Visibility

`POST /api/v1/connections/toggle-visibility/:id` toggles `is_hidden` on an account. Hidden accounts leave net worth and transaction queries. The Connections panel still lists them.

### Disconnect and remove

1. Disconnect an institution with `POST /api/v1/connections/disconnect/:id`. Associated accounts and transactions leave the active linked set.
2. Disconnected account rows can remain until the user removes them.
3. Remove a disconnected account with `POST /api/v1/connections/remove-account/:id`.

The Connections panel requires disconnect before permanent remove. Guidance appears in the UI.

### Update mode

When an item needs reauth (`needs_reauth` / login required), call `POST /api/v1/connections/create-update-session/:id`. The client opens update-mode Link or the equivalent Stripe repair flow, then completes the session. This restores access without a full new connection when the provider supports it.

### Rate limits and quotas

| Limit | When | Response |
|-------|------|----------|
| Manual sync cooldown | Plaid, once per minute per user | `429` `connection_sync_rate_limit` |
| Item/connection cap | Per subscription tier | `403` `connection_item_limit` |
| Monthly API calls | Provider-metered operations (not Plaid sync) | `429` `connection_api_limit` |

When subscriptions are disabled, bank and API quotas are not enforced.

## API endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/api/v1/connections/provider` | Yes | Active provider info |
| GET | `/api/v1/connections` | Yes | List connections and accounts |
| POST | `/api/v1/connections/create-session` | Yes | Start new link session |
| POST | `/api/v1/connections/complete` | Yes | Finish link and persist item |
| POST | `/api/v1/connections/sync` | Yes | Sync all items |
| POST | `/api/v1/connections/create-update-session/:id` | Yes | Start update/reauth session |
| POST | `/api/v1/connections/sync-item/:id` | Yes | Sync one item |
| POST | `/api/v1/connections/disconnect/:id` | Yes | Disconnect institution |
| POST | `/api/v1/connections/toggle-visibility/:id` | Yes | Hide or show account |
| POST | `/api/v1/connections/remove-account/:id` | Yes | Remove disconnected account |
| POST | `/api/v1/plaid/webhook` | No (signed) | Plaid transaction/item webhooks |

## Configuration

| Variable | Purpose |
|----------|---------|
| `FINANCIAL_PROVIDER` | `plaid` or `stripe` (default `plaid`) |
| `PLAID_CLIENT_ID` | Plaid application ID |
| `PLAID_ENV` | `sandbox` or `production` |
| `PLAID_SANDBOX_SECRET` | Plaid sandbox secret |
| `PLAID_PROD_SECRET` | Plaid production secret |
| `PLAID_WEBHOOK_URL` | Optional webhook override |
| `ENCRYPTION_KEY` | Encrypts Plaid access tokens at rest |
| `API_PUBLIC_URL` | Default webhook base URL |
| `STRIPE_SANDBOX_SECRET` / `STRIPE_PROD_SECRET` | Stripe secret when provider is `stripe` |
| `STRIPE_SANDBOX_PUBLISHABLE` / `STRIPE_PUBLISHABLE` | Stripe publishable key for the SPA |
| `SUBSCRIPTIONS_ENABLED` | Enables tier item/API limits when `true` |

## Related files

- `app/internal/services/financial/`
- `app/internal/services/plaid/`
- `app/internal/services/stripefc/`
- `app/internal/services/connections/payload.go`
- `app/internal/server/handler/connections.go`
- `frontend/src/lib/connections.ts`
- `frontend/src/lib/stripe.ts`
- `frontend/src/components/settings/ConnectionsPanel.tsx`
- `frontend/src/components/onboarding/OnboardingConnectStep.tsx`

## Related docs

- [Transactions](./transactions.md)
- [Subscriptions](./subscriptions.md)
- [Settings](./settings.md)
- [Onboarding](./onboarding.md)
- [Environment variables](../getting-started/environment.md)
