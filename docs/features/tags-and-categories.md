# Tags and Categories

## Overview

Tags and categories organize transactions for filtering and analytics. Categories group tags. Tags support colors, auto-tag filters, and an exclude-from-totals flag. Users manage them in Settings → Tags. Legacy `/tags` redirects to that tab.

## How it works

### Settings Tags tab

Open `/settings?tab=tags`. The panel shows category cards, tag chips, create/edit modals, and delete flows. Status messages use shared `PageStatusBanner`.

### Categories

Users create, rename, and delete categories.

Delete options:

1. Move tags to Misc (`move_to_misc`).
2. Move tags to another category (`move_to`).
3. Delete with the category’s tags when the store action removes them with the chosen target.

Merge categories is not available.

### Auto-tag filters

Each tag can define filters. On sync, `AutoTagTransaction` matches name, merchant, and provider category.

| Filter type | Match rule |
|-------------|------------|
| `string` | Case-insensitive substring |
| `regex` | Regular expression |
| `amount_greater` | Stored amount greater than value |
| `amount_less` | Stored amount less than value |
| `amount_equal` | Stored amount equal to value |

Save & Apply can create filters and apply them to past transactions.

### Drag and drop

Users move tags between categories with pointer drag (desktop and touch).

1. Press a tag chip and move at least 8 px to start drag.
2. Hover another category card (`data-category-id`).
3. Release to call `POST /api/v1/tags/:id/move`.
4. The UI applies an optimistic update, then reconciles the server payload.

Taps still open edit/delete. Drag uses a floating ghost chip and viewport auto-scroll near edges.

### Exclude from totals (`is_hidden`)

Tag form switch **Exclude from totals** sets `tags.is_hidden`. Hidden tags:

- Still appear on transactions and in the Tags UI (struck-through / eye-off treatment).
- Exclude those transactions from dashboard cashflow, spending trend, and tag donuts.
- Exclude those transactions from transaction analytics.

Seeded **Transfers** is hidden by default.

### Colors

Tags use palette keys (for example `rose`, `emerald`). Components resolve keys to hex with `tagColorHex`. The tag form includes a color picker.

### Seeded defaults

`TaggingService.SeedDefaults` runs when a user has no categories (for example on first dashboard load). Defaults include:

| Category | Example tags |
|----------|--------------|
| Food & Drink | Dining Out, Groceries |
| Transport | Automotive, Travel |
| Shopping | General, Home |
| Recurring | Bills, Services |
| Health & Wellness | Medical, Personal Care |
| Financial | Income, Fees, Transfers (hidden) |
| Leisure | Entertainment |

Each default tag receives a provider-category string filter and a palette color.

### `/tags` redirect

`frontend/src/routes/tags.tsx` navigates to `/settings?tab=tags`. Tags are not in the main navbar or mobile tab bar.

## API endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/api/v1/tags` | Yes | Full tags payload |
| GET | `/api/v1/tags/:id/filters` | Yes | Filters for one tag |
| POST | `/api/v1/tags` | Yes | Create tag |
| PUT | `/api/v1/tags/:id` | Yes | Update tag |
| DELETE | `/api/v1/tags/:id` | Yes | Delete tag |
| POST | `/api/v1/tags/:id/move` | Yes | Move tag to category |
| POST | `/api/v1/categories` | Yes | Create category |
| PUT | `/api/v1/categories/:id` | Yes | Rename category |
| DELETE | `/api/v1/categories/:id` | Yes | Delete category |

Mutations return the refreshed tags payload when the handlers provide it, so the client avoids an extra fetch.

## Configuration

No tag-specific environment variables. Schema column `tags.is_hidden` is created in `schema.sql` and migration `0001_tags_is_hidden.sql` for existing databases.

## Related files

- `app/internal/services/tags/tags.go`
- `app/internal/services/tags/api_payload.go`
- `app/internal/server/handler/tags.go`
- `app/internal/storage/queries/tag.go`
- `app/internal/storage/queries/dashboard_analytics.go`
- `frontend/src/components/settings/TagsPanel.tsx`
- `frontend/src/components/tags/`
- `frontend/src/routes/tags.tsx`
- `frontend/src/routes/settings.tsx`

## Related docs

- [Settings](./settings.md)
- [Transactions](./transactions.md)
- [Dashboard](./dashboard.md)
