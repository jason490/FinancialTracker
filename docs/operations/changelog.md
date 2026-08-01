# Changelog

## Overview

This page records notable product and engineering changes. Earlier notes lived in `summary.md`. New work belongs here under the matching section. Do not remove existing entries. Add new bullets when you change behavior.

Writing style follows ASD-STE100: active voice, short sentences, and clear terms.

## Architecture

- The stack is a Go REST API and a SolidJS SSG SPA. Capacitor targets Android.
- Handlers stay thin. Services own business logic. External DTOs live in `models/external`.
- The legacy Templ/HTMX UI and `legacyHandler` package were removed.
- Type sync exports SPA-facing DTOs only through `make sync-types`.

## Authentication and onboarding

- Sessions use HTTP-only cookies. `GET /api/v1/auth/me` returns a lean `SessionProfile`.
- Google SSO supports login, register, Settings link, and delete reauth. Password-only accounts are not auto-linked by email.
- Password reset uses 6-digit codes with TTL, attempt caps, and SMTP or development console delivery.
- Onboarding runs Welcome → Plan (optional) → Connect (skippable). Incomplete users stay on `/onboarding`.
- When `SUBSCRIPTIONS_ENABLED=false`, registration requires single-use invite codes.

## Bank connections

- `FINANCIAL_PROVIDER` selects Plaid (default) or Stripe Financial Connections.
- SPA flows use `/api/v1/connections/*`. Only the Plaid webhook remains under `/api/v1/plaid/webhook`.
- Plaid access tokens are encrypted at rest. Manual sync is rate-limited. Sync does not consume monthly API quota.
- Accounts support hide, disconnect, remove, and update-mode reauth.

## Transactions

- The list supports search, UTC date filters, ABS amount filters, tags, bulk tag actions, and sync.
- Settings → Data exports a ZIP that contains one CSV. Export reserves one monthly API call.
- **Show graphs** calls `GET /api/v1/transactions/analytics` for totals, category donut, and monthly cashflow.

## Tags and categories

- Tags live in Settings → Tags. Legacy `/tags` redirects to `/settings?tab=tags`.
- Auto-tag filters support string, regex, and amount rules.
- `tags.is_hidden` excludes tagged transactions from analytics totals. Seeded Transfers is hidden by default.
- Pointer drag moves tags between categories on desktop and touch.

## Dashboard

- Separate desktop and mobile layouts persist in one JSON document.
- Widgets include net worth, this month, monthly cashflow, tag donuts, accounts, recent transactions, and quick actions.
- Edit mode uses SortableJS without remounting charts on drop. ApexCharts refit on container resize.

## Settings and appearance

- Tabs: Account, Connections, Appearance, Plan (when enabled), Tags, Data.
- Themes include Light, Dark, System, Tokyo Night, Coffee, Forest, Rose, Midnight, and Parchment.
- Account deletion requires password or Google reauth.

## Subscriptions and billing

- `SUBSCRIPTIONS_ENABLED` toggles Stripe billing and quota enforcement.
- Free / Plus / Premium tiers limit bank items and monthly API calls.
- Stripe Checkout, Customer Portal, and billing webhooks sync tier state.
- `SUBSCRIPTION_OVERRIDES` and `user_privileges` grant unlimited limits or coupons.

## Security and operations

- CSRF uses Echo double-submit cookies. OAuth state is HMAC-signed.
- Production containers drop privileges. Caddy serves automatic HTTPS with `SITE_HOST`.
- Production migrations use ordered SQL files under `app/database/migrations/`.
- Startup validates required env vars before the database opens.

## Fixes (selected)

- Plaid webhook body-hash compare used the wrong `ConstantTimeCompare` check and rejected valid webhooks.
- Plaid webhook key cache used `Unlock` on an `RLock` and panicked on new keys.
- Liability display no longer shows credit limit when balance is zero.
- Auth transition overlay no longer flashes blank content into onboarding or dashboard.
- Connections panel refetch uses `startTransition` to avoid loading flicker.
- Invite code copy falls back when `navigator.clipboard` is unavailable.

## Related docs

- [Documentation index](../index.md)
- [Architecture overview](../architecture/overview.md)
