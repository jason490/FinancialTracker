# Database and Migrations

## Overview

FinancialTracker stores all application data in SQLite. Foreign keys are enabled. Cascades remove dependent rows when a user is deleted.

## Schema sources

| File | Role |
|------|------|
| `app/database/schema.sql` | Production fresh install schema |
| `app/database/test_schema.sql` | Development full recreate (must mirror production schema) |
| `app/database/migrations/*.sql` | Ordered deltas for existing production databases |

Main tables include users, sessions, SSO links, Plaid and Stripe FC items/accounts, usage counters, transactions, categories, tags, tag filters, dashboard layouts, password reset codes, registration codes, and user privileges.

## Startup behavior

### Development (`ENV=development`)

1. Apply `test_schema.sql` (drop and recreate).
2. Seed `test@test.com` / `test` when that user is missing.
3. Do not run the production migration runner for incremental deltas.

### Production

1. Apply `schema.sql` with `CREATE TABLE IF NOT EXISTS`.
2. Detect whether the database is brand-new (presence of the `users` table before schema apply).
3. Run the file-based migration runner from `app/internal/storage/migrate.go`.

## Migration rules

- Place files in `app/database/migrations/` with zero-padded numeric prefixes (for example `0002_add_column.sql`).
- The runner applies pending files in lexical order.
- Each file runs in its own transaction.
- Applied versions are stored in `schema_migrations`.
- Fresh installs baseline pending migrations (record them without execution) because `schema.sql` already contains the final shape.
- Existing databases execute pending deltas.
- The production Docker image must include `database/migrations`.

First shipped migration: `0001_tags_is_hidden.sql` adds `tags.is_hidden`.

## Provider scoping

`transactions.provider` is `plaid` or `stripe`. Reads and writes use only the active provider tables. The app does not merge history across providers.

## Related docs

- [Architecture overview](./overview.md)
- [Tags and categories](../features/tags-and-categories.md)
- [Production deployment](../getting-started/production.md)
