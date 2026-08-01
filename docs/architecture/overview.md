# Architecture Overview

## Overview

FinancialTracker is a self-hosted expense tracker. A Go REST API owns business logic and SQLite storage. A SolidJS SPA talks to that API. Capacitor can wrap the static frontend for Android.

## Stack

| Layer | Technology |
|-------|------------|
| Backend | Go, Echo v5, SQLite |
| Frontend | SolidStart (SolidJS), Vinxi, SSG |
| Mobile | Capacitor (Android) |
| Styling | Vanilla CSS / CSS Modules |
| Bank sync | Plaid (default) or Stripe Financial Connections |
| Billing | Stripe (optional) |
| Auth | Sessions, bcrypt, Google OAuth |
| Proxy / TLS | Caddy |

## Request flow

1. The browser or Capacitor WebView loads the static SPA.
2. The SPA calls `/api/v1/*` with `credentials: "include"`.
3. Caddy proxies `/api` to the Go process in Docker deployments.
4. Handlers parse HTTP input and call services.
5. Services use storage queries against SQLite.
6. Responses return raw JSON on success and structured `APIError` objects on failure.

Handlers must not call storage directly. Add service methods when new behavior is required.

## Frontend layout

- Desktop and tablet use a top navbar.
- Mobile uses a bottom tab bar.
- Main nav items are Dashboard, Transactions, and Settings.
- Auth pages use `AuthLayout`. App pages use `AppLayout`.

## Type sync

`make sync-types` and `make build` regenerate `frontend/src/lib/types.ts` from Go external DTOs and shared model structs. Do not edit that generated file by hand.

## Related docs

- [API reference](./api.md)
- [Database and migrations](./database.md)
- [Security](./security.md)
- [Development setup](../getting-started/development.md)
