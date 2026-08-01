# Subscriptions

## Overview

Subscriptions control paid plans, Stripe billing, and per-tier usage limits. Toggle with `SUBSCRIPTIONS_ENABLED` (default `true`). When disabled, billing UI and quota enforcement stop, and registration becomes invite-only.

## How it works

### `SUBSCRIPTIONS_ENABLED`

| Value | Behavior |
|-------|----------|
| `true` (default) | Stripe checkout/portal/webhooks active; tier limits enforced; Plan tab and onboarding plan step visible |
| `false` | No Stripe billing; unlimited bank/API usage; Plan UI hidden; invite codes required to register |

`GET /api/v1/subscription` includes `subscriptions_enabled`.

### Tiers

| Tier | Max bank items | Max API calls / month | Price (cents) |
|------|----------------|------------------------|---------------|
| Free | 2 | 20 | 0 |
| Plus | 5 | 50 | 299 |
| Premium | 15 | 100 | 599 |

Billing period for free users anchors to signup. Paid users anchor to `subscription_started_at` (monthly anniversary). Metered features share one counter in `plaid_api_usage` for the period (exports and other reserved API calls). Plaid transaction sync does not consume that quota.

### Stripe checkout, portal, and webhooks

Production upgrades use Stripe Checkout.

1. Client calls `POST /api/v1/subscription/checkout` with the target tier.
2. User completes Stripe Checkout.
3. Return URLs land on Settings Plan with `checkout=success` or `checkout=cancelled`.
4. `POST /api/v1/stripe/webhook` verifies `Stripe-Signature` and syncs `subscription_tier`, `subscription_started_at`, and `stripe_subscription_id`.

Customer Portal: `POST /api/v1/subscription/portal` for manage-billing. In development, `POST /api/v1/subscription/change` can switch tiers directly when skeleton switching is allowed. Production plan changes require Stripe.

### Overrides

Operators grant privileges without changing Stripe products.

- Env: `SUBSCRIPTION_OVERRIDES` (email) and `SUBSCRIPTION_OVERRIDE_USER_IDS` (numeric ID).
- Format example: `alice@example.com=unlimited,coupon=FRIENDS50;bob@example.com=coupon=HALF_OFF`.
- SQLite table `user_privileges` stores `unlimited_limits` and `stripe_coupon_id`.

Unlimited overrides surface as Unlimited meters on the Plan tab. Coupons apply at Checkout when configured.

### Usage limits

`EffectiveLimits` returns tier limits or unlimited when subscriptions are disabled or the user has unlimited privileges. Limit errors:

| Code | HTTP | Meaning |
|------|------|---------|
| `connection_item_limit` | 403 | Too many bank connections |
| `connection_api_limit` | 429 | Monthly API quota exhausted |
| `export_api_limit` | 429 | Export blocked by shared API quota |

### Onboarding plan step visibility

When subscriptions are enabled, onboarding includes Welcome → Plan → Connect. When disabled, the plan step is omitted and Welcome continues to Connect. Settings also hides the Plan tab in that mode.

## API endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/api/v1/subscription` | Yes | Tier, usage, catalog, flags |
| POST | `/api/v1/subscription/change` | Yes | Dev/direct tier change |
| POST | `/api/v1/subscription/checkout` | Yes | Create Stripe Checkout session |
| POST | `/api/v1/subscription/portal` | Yes | Create Customer Portal session |
| POST | `/api/v1/stripe/webhook` | No (signed) | Billing webhook receiver |

## Configuration

| Variable | Purpose |
|----------|---------|
| `SUBSCRIPTIONS_ENABLED` | Master billing/limits toggle |
| `STRIPE_SANDBOX_SECRET` / `STRIPE_PROD_SECRET` | Stripe secret key |
| `STRIPE_SANDBOX_PUBLISHABLE` / `STRIPE_PUBLISHABLE` | Publishable key |
| `STRIPE_PRICE_PLUS` / `STRIPE_PRICE_PREMIUM` | Production price IDs |
| `STRIPE_SANDBOX_PRICE_PLUS` / `STRIPE_SANDBOX_PRICE_PREMIUM` | Sandbox price IDs |
| `STRIPE_WEBHOOK_SECRET` / `STRIPE_SANDBOX_WEBHOOK_SECRET` | Billing webhook secrets |
| `SUBSCRIPTION_OVERRIDES` | Per-email privilege overrides |
| `SUBSCRIPTION_OVERRIDE_USER_IDS` | Per-user-ID privilege overrides |
| `ENV` | Selects sandbox vs production Stripe keys and allows skeleton tier change in development |

Stripe keys are not required at startup when subscriptions are disabled.

## Related files

- `app/internal/services/subscription/service.go`
- `app/internal/services/subscription/privileges.go`
- `app/internal/services/stripebilling/`
- `app/internal/server/handler/subscription.go`
- `app/internal/config/env.go`
- `frontend/src/components/settings/PlanPanel.tsx`
- `frontend/src/components/settings/PlanPicker.tsx`
- `frontend/src/components/onboarding/OnboardingPlanStep.tsx`
- `frontend/src/lib/subscription.ts`

## Related docs

- [Authentication](./authentication.md)
- [Onboarding](./onboarding.md)
- [Settings](./settings.md)
- [Bank connections](./bank-connections.md)
- [Environment variables](../getting-started/environment.md)
