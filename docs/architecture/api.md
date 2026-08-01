# API Reference

## Overview

The Go server exposes a JSON API under `/api/v1` plus `/health`. Authenticated routes require the `Session` HTTP-only cookie that login and SSO set.

## Route groups

| Area | Base path |
|------|-----------|
| Auth | `/api/v1/auth/*` |
| Dashboard | `/api/v1/dashboard` |
| Transactions | `/api/v1/transactions` |
| Tags and categories | `/api/v1/tags`, `/api/v1/categories` |
| Bank connections | `/api/v1/connections/*` |
| Settings | `/api/v1/settings` |
| Subscription | `/api/v1/subscription` |
| Admin (invite codes) | `/api/v1/admin/registration-codes` |
| CSRF | `/api/v1/csrf` |
| Plaid webhook | `/api/v1/plaid/webhook` |
| Stripe billing webhook | `/api/v1/stripe/webhook` |
| Health | `/health` |

## Conventions

- Success responses return the resource body with HTTP 2xx.
- Errors return `{ "code": "...", "message": "..." }` with 4xx or 5xx.
- State-changing requests send `X-CSRF-Token` from the `_csrf` cookie.
- External DTOs omit internal database IDs unless the SPA needs them for a flow.

## Feature details

Use the feature pages for full endpoint tables:

- [Authentication](../features/authentication.md)
- [Bank connections](../features/bank-connections.md)
- [Transactions](../features/transactions.md)
- [Tags and categories](../features/tags-and-categories.md)
- [Dashboard](../features/dashboard.md)
- [Settings](../features/settings.md)
- [Subscriptions](../features/subscriptions.md)

## Related docs

- [Architecture overview](./overview.md)
- [Security](./security.md)
