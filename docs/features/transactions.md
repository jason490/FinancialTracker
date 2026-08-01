# Transactions

## Overview

The Transactions page lists synced bank activity for the active financial provider. Users filter, bulk-tag, sync, export, and view filter-aware charts. Amounts follow the provider sign convention in storage. The UI and CSV export invert signs for display.

## How it works

### List

`GET /api/v1/transactions` returns a paginated page with transactions, total count, tags, and categories. Queries include only non-hidden accounts for the active provider. The SPA keeps metadata (categories and tags) stable while results refresh.

### Filters

Users filter by search text, date range, amount range, tags, pending state, and sort. The frontend debounces search, dates, and amounts. Dates use UTC. Min/max amount filters apply to `ABS(amount)` so magnitude works for expenses and income. Changing filters clears the current selection.

### Bulk tags

1. Select one or more transactions (full card click or checkbox).
2. Choose a tag in the bulk action bar.
3. Add with `POST /api/v1/transactions/bulk-add-tag` or remove with `POST /api/v1/transactions/bulk-remove-tag`.

### Sync button

Dashboard Quick Actions and the Transactions header call `POST /api/v1/connections/sync`. The page shows status in `PageStatusBanner`. Plaid manual sync allows one request per minute per user.

### CSV / ZIP export

Settings → Data exports all visible transactions.

1. Call `GET /api/v1/transactions/export`.
2. The API builds a ZIP with one CSV (`financial-tracker-transactions-YYYY-MM-DD.csv`).
3. Each export reserves one monthly API call. Exhausted quota returns `429` `export_api_limit`.

CSV columns: `date`, `name`, `merchant_name`, `amount`, `currency_sign`, `pending`, `provider`, `provider_category`, `tags`, `tag_categories`.

### Show graphs

The toolbar **Show graphs** control swaps the list for `TransactionGraphsPanel`.

`GET /api/v1/transactions/analytics` reuses the list filter WHERE clause and returns:

- Total spend and income
- Spending-by-category donut slices
- Monthly income/spend trend (last six months without a date filter; otherwise the filtered range)

Hidden-tag transactions are excluded from analytics. They still appear in the list.

### Amount sign convention

| Layer | Positive amount | Negative amount |
|-------|-----------------|-----------------|
| Storage (Plaid/Stripe) | Expense / debit | Income / credit |
| UI (`formatAmount`) | Income shown as `+` | Expense shown as `-` |
| CSV export | Credit (inverted) | Debit (inverted) |

The UI and export both use display amount `-stored_amount`.

### Hidden-tag exclusion in analytics

Any transaction that carries a tag with `is_hidden = 1` is dropped from transaction analytics and matching dashboard cashflow/tag totals. The ledger still shows those rows with a “Not counted” treatment.

## API endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/api/v1/transactions` | Yes | Paginated filtered list |
| GET | `/api/v1/transactions/analytics` | Yes | Filter-aware charts payload |
| GET | `/api/v1/transactions/export` | Yes | ZIP/CSV download |
| POST | `/api/v1/transactions/bulk-add-tag` | Yes | Bulk add tag |
| POST | `/api/v1/transactions/bulk-remove-tag` | Yes | Bulk remove tag |
| POST | `/api/v1/connections/sync` | Yes | Manual bank sync |

## Configuration

| Variable | Purpose |
|----------|---------|
| `FINANCIAL_PROVIDER` | Scopes transaction reads/writes to `plaid` or `stripe` |
| `SUBSCRIPTIONS_ENABLED` | When `true`, export consumes monthly API quota |

No transaction-specific env vars exist beyond provider and subscription settings.

## Related files

- `app/internal/services/transactions/transactions.go`
- `app/internal/services/transactions/api_payload.go`
- `app/internal/server/handler/transactions.go`
- `app/internal/storage/queries/transaction.go`
- `app/internal/storage/queries/transaction_analytics.go`
- `frontend/src/routes/transactions.tsx`
- `frontend/src/components/transactions/`
- `frontend/src/lib/transactions.ts`
- `frontend/src/lib/format.ts`
- `frontend/src/components/settings/DataPanel.tsx`

## Related docs

- [Bank connections](./bank-connections.md)
- [Tags and categories](./tags-and-categories.md)
- [Dashboard](./dashboard.md)
- [Settings](./settings.md)
- [Subscriptions](./subscriptions.md)
