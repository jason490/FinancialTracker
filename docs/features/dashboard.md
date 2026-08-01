# Dashboard

## Overview

The dashboard is the main overview after onboarding. It shows account buckets, cashflow, tag breakdowns, recent transactions, and quick actions in a customizable widget grid. Desktop and mobile layouts persist separately. Charts use ApexCharts with theme-aware colors.

## How it works

### Widgets

| Widget ID | Label | Default span |
|-----------|-------|--------------|
| `month_cashflow` | This Month | 2 |
| `quick_actions` | Quick Actions | 1 |
| `net_worth` | Net Worth | 2 |
| `spending_by_tag` | Spending by Tag | 1 |
| `spending_trend` | Monthly Cashflow | 2 |
| `income_by_tag` | Income by Tag | 1 |
| `cash_accounts` | Cash & Checking | 1 |
| `savings_accounts` | Savings | 1 |
| `credit_accounts` | Credit Cards | 1 |
| `loan_accounts` | Loans | 1 |
| `investment_accounts` | Investments | 1 |
| `recent_transactions` | Recent Transactions | 2 |

Account widgets group linked accounts by bucket. Liability buckets use balance owed only. Quick Actions includes sync and common navigation. Tag donuts show abbreviated center totals and full amounts in tooltips.

### Layouts (desktop and mobile)

The backend stores one JSON layout with `Desktop` and `Mobile` widget arrays. The SPA selects the array from viewport size. Saving one device does not change the other. `NormalizeLayout` adds new default widgets and drops unknown IDs.

### Edit mode

1. Desktop: click **Customize**.
2. Mobile (coarse pointer): long-press a widget (hold timer) to enter edit mode.
3. Drag widgets with SortableJS. On drop, the grid updates widget `order` in place without remounting charts.
4. Toggle visibility per widget while editing.
5. Save with `POST /api/v1/dashboard/layout` for the active device type.

Hidden widgets stay available in edit mode and disappear in normal view.

### Charts

- **Monthly Cashflow** (`spending_trend`): spline area of income and spending for the last six months.
- **Spending by Tag / Income by Tag**: donut charts from tagged breakdowns.
- Tooltips follow the active theme. A `ResizeObserver` refits charts after SPA route changes.

### Analytics sources

`GET /api/v1/dashboard` builds the payload from:

- Visible accounts for the active provider
- Recent transactions
- Dashboard analytics queries (monthly spending, month cashflow, tagged breakdowns)

Hidden accounts and transactions with hidden tags are excluded from totals where the analytics predicates apply. Seed defaults for tags run on dashboard load when the user has no categories.

### Skeletons

While the dashboard loads, `DashboardSkeletonGrid` mirrors default spans and shapes. `LoadingCrossfade` and `useMinLoadingHold` hold a short minimum loading time, then crossfade into widgets with staggered reveal. Post-auth handoff to the dashboard preloads the route chunk so the skeleton shows instead of a full-screen overlay.

## API endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/api/v1/dashboard` | Yes | Dashboard payload (`?edit=1` for edit fetch) |
| POST | `/api/v1/dashboard/layout` | Yes | Save desktop or mobile layout |
| POST | `/api/v1/connections/sync` | Yes | Sync from Quick Actions |

## Configuration

No dashboard-specific environment variables. Layout and analytics depend on `FINANCIAL_PROVIDER` and linked accounts.

## Related files

- `app/internal/services/dashboard/service.go`
- `app/internal/services/dashboard/layout.go`
- `app/internal/services/dashboard/api_payload.go`
- `app/internal/services/dashboard/accounts.go`
- `app/internal/storage/queries/dashboard_analytics.go`
- `app/internal/server/handler/dashboard.go`
- `frontend/src/routes/dashboard.tsx`
- `frontend/src/components/dashboard/`
- `frontend/src/lib/dashboard.ts`
- `frontend/src/lib/dashboard-widgets.ts`

## Related docs

- [Transactions](./transactions.md)
- [Tags and categories](./tags-and-categories.md)
- [Bank connections](./bank-connections.md)
- [Onboarding](./onboarding.md)
