# Environment Variables

## Overview

FinancialTracker reads configuration from environment variables and an optional `.env` file. This page lists the variables by area. Production also uses values from [`env.prod.example`](../../env.prod.example).

## Core

| Variable | Default | Description |
|----------|---------|-------------|
| `ENV` | production | `development` resets the DB, seeds the test user, relaxes rate limits, and logs password-reset codes |
| `PORT` | `8080` | API listen port |
| `DATABASE_PATH` | `./database/main.db` | SQLite file path |
| `ENCRYPTION_KEY` | — | Required in production. Exactly 32 bytes for AES-256 |
| `API_PUBLIC_URL` | — | Public API URL. Build-time value is baked into the frontend. Runtime value is passed to the backend |
| `FRONTEND_URL` | — | Public SPA URL. Same build-time and runtime rules as `API_PUBLIC_URL` |
| `SITE_HOST` | — | Bare hostname for Caddy TLS. Runtime only |
| `ACME_EMAIL` | — | Email for Let's Encrypt notices. Runtime only |

## Auth

| Variable | Description |
|----------|-------------|
| `GOOGLE_CLIENT_ID` | Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret |
| `REGISTRATION_ADMIN_EMAILS` | Semicolon-separated emails that may create invite codes |

## Subscriptions

| Variable | Default | Description |
|----------|---------|-------------|
| `SUBSCRIPTIONS_ENABLED` | `true` | `false` disables Stripe billing, removes usage limits, and enables invite-only registration |
| `STRIPE_SANDBOX_SECRET` | — | Stripe secret key (development) |
| `STRIPE_PROD_SECRET` | — | Stripe secret key (production) |
| `STRIPE_SANDBOX_PUBLISHABLE` | — | Stripe publishable key (development) |
| `STRIPE_PUBLISHABLE` | — | Stripe publishable key (production) |
| `STRIPE_PRICE_PLUS` / `STRIPE_PRICE_PREMIUM` | — | Production price IDs |
| `STRIPE_SANDBOX_PRICE_PLUS` / `STRIPE_SANDBOX_PRICE_PREMIUM` | — | Sandbox price IDs |
| `STRIPE_WEBHOOK_SECRET` / `STRIPE_SANDBOX_WEBHOOK_SECRET` | — | Billing webhook signing secrets |
| `SUBSCRIPTION_OVERRIDES` | — | Per-email overrides, for example `user@example.com=unlimited,coupon=CODE` |
| `SUBSCRIPTION_OVERRIDE_USER_IDS` | — | Same format keyed by user ID |

## Bank provider

| Variable | Default | Description |
|----------|---------|-------------|
| `FINANCIAL_PROVIDER` | `plaid` | `plaid` or `stripe` |
| `PLAID_CLIENT_ID` | — | Plaid application ID |
| `PLAID_ENV` | `sandbox` in development | `sandbox` or `production`. Set this explicitly on the VPS |
| `PLAID_SANDBOX_SECRET` | — | Plaid sandbox secret |
| `PLAID_PROD_SECRET` | — | Plaid production secret |
| `PLAID_WEBHOOK_URL` | — | Optional. Defaults to `{API_PUBLIC_URL}/api/v1/plaid/webhook` |

## Email (production password reset)

| Variable | Description |
|----------|-------------|
| `MAIL_SMTP_HOST` | SMTP server |
| `MAIL_SMTP_PORT` | For example `587` |
| `MAIL_SMTP_USER` | SMTP username |
| `MAIL_SMTP_PASSWORD` | SMTP password or API key |
| `MAIL_FROM` | Sender address |
| `MAIL_FROM_NAME` | Optional display name |

When `MAIL_SMTP_HOST` is unset and `ENV=development`, the API prints reset codes to the logs.

## Frontend build-time (Vinxi / Vite)

| Variable | Description |
|----------|-------------|
| `VITE_API_PUBLIC_URL` | API URL embedded in the client bundle |
| `VITE_FRONTEND_URL` | Frontend URL for OAuth return paths |

## Optional integrations

| Integration | When needed |
|-------------|-------------|
| Plaid | `FINANCIAL_PROVIDER=plaid` — client ID plus sandbox or production secret |
| Stripe FC | `FINANCIAL_PROVIDER=stripe` — Stripe sandbox or production keys |
| Stripe billing | `SUBSCRIPTIONS_ENABLED=true` — Stripe keys, price IDs, webhook secret |
| Google OAuth | SSO — `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` |
| SMTP | Production password reset — `MAIL_SMTP_*`, `MAIL_FROM` |

## Related docs

- [Development setup](./development.md)
- [Production deployment](./production.md)
- [Authentication](../features/authentication.md)
- [Subscriptions](../features/subscriptions.md)
- [Bank connections](../features/bank-connections.md)
