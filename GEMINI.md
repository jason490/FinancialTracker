# FinancialTracker Project Context

This project is a high-performance, user-centric financial tracking application inspired by Rocket Money. It is designed for blazing-fast performance and seamless expense management by directly syncing with bank accounts via **Plaid** (default) or **Stripe Financial Connections** (alternate provider).
You are Full-stack developer and designer. Ensure any interactions feel seemless on both andriod, webapps, and websites. Any interactions should be properly organized (depending on the amount of visual bloat seen) into sidebar, navbar (top for webapp, bottom for mobile), modals, etc.

## 1. Project Overview

*   **Core Purpose:** Automated expense tracking, smart tagging, and categorization.
*   **Main Technologies:**
    *   **Backend:** Go (Golang) using the `echo` v5 framework as a pure REST API.
    *   **Database:** SQLite (embedded relational database).
    *   **Frontend:** SolidStart (high-performance, reactive UI library) with SSG (no SSR).
    *   **Mobile Bridge:** Capacitor (targeting Android and Web platforms).
    *   **Styling:** Vanilla CSS / CSS Modules (avoid TailwindCSS for better mobile optimization).
    *   **External APIs:** Plaid (default bank sync), Stripe Financial Connections (alternate provider), Google OAuth (SSO).

### Bank connection provider switch

Set `FINANCIAL_PROVIDER` to choose the active integration:

| Value | Provider |
|-------|----------|
| `plaid` (default) | Plaid Link |
| `stripe` | Stripe Financial Connections |

When `stripe` is active, reads/writes use `stripe_fc_*` tables only. When `plaid` is active, reads/writes use `plaid_*` tables only. There is no cross-provider legacy data merging.

### Environment variables

| Variable | Purpose |
|----------|---------|
| `ENV=development` | Dev mode (relaxed rate limits, sandbox keys, skeleton plan switching when subscriptions enabled, resets DB via `test_schema.sql`, seeds `test@test.com`) |
| `ENV=production` | Production mode (or omit `ENV` / set any other value). Uses idempotent `schema.sql`, no test user seed, secure cookies, auth rate limits, Stripe prod keys |
| `PORT` | HTTP listen port (default `8080`) |
| `DATABASE_PATH` | SQLite file path (default `./database/main.db`) |
| `FINANCIAL_PROVIDER` | `plaid` or `stripe` (default `plaid`) |
| `SUBSCRIPTIONS_ENABLED` | When `true` (default), Stripe billing and per-tier usage limits are active. Set to `false` to disable paid plans, Stripe checkout, and quota enforcement. |
| `STRIPE_SANDBOX_SECRET` | Stripe secret key when `ENV=development` (only required when `SUBSCRIPTIONS_ENABLED=true` and using Stripe) |
| `STRIPE_PROD_SECRET` | Stripe secret key in production |
| `STRIPE_SANDBOX_PUBLISHABLE` | Stripe publishable key returned to the SPA in development |
| `STRIPE_PUBLISHABLE` | Stripe publishable key in production |
| `STRIPE_PRICE_PLUS` | Stripe Price ID for the Plus plan (production) |
| `STRIPE_PRICE_PREMIUM` | Stripe Price ID for the Premium plan (production) |
| `STRIPE_SANDBOX_PRICE_PLUS` | Sandbox Price ID for Plus when `ENV=development` |
| `STRIPE_SANDBOX_PRICE_PREMIUM` | Sandbox Price ID for Premium when `ENV=development` |
| `STRIPE_WEBHOOK_SECRET` | Stripe billing webhook signing secret (production) |
| `STRIPE_SANDBOX_WEBHOOK_SECRET` | Sandbox billing webhook secret when `ENV=development` |
| `SUBSCRIPTION_OVERRIDES` | Semicolon-separated per-email overrides, e.g. `alice@example.com=unlimited,coupon=FRIENDS50;bob@example.com=coupon=HALF_OFF` |
| `SUBSCRIPTION_OVERRIDE_USER_IDS` | Same format keyed by numeric user ID, e.g. `1=unlimited;2=coupon=promo_abc` |
| `PLAID_CLIENT_ID`, `PLAID_ENV`, `PLAID_SANDBOX_SECRET`, `PLAID_PROD_SECRET` | Plaid credentials (default provider) |
| `PLAID_WEBHOOK_URL` | Optional override for the Plaid webhook callback (defaults to `{API_PUBLIC_URL}/api/v1/plaid/webhook`) |
| `ENCRYPTION_KEY` | 32-byte AES key for encrypting Plaid access tokens at rest |
| `MAIL_SMTP_HOST` | SMTP host for password reset emails (production) |
| `MAIL_SMTP_PORT` | SMTP port (e.g. `587`) |
| `MAIL_SMTP_USER` | SMTP username |
| `MAIL_SMTP_PASSWORD` | SMTP password or API key |
| `MAIL_FROM` | Sender email address for outbound mail |
| `MAIL_FROM_NAME` | Optional sender display name (defaults to Financial Tracker) |

When `MAIL_SMTP_HOST` is unset and `ENV=development`, reset codes are logged to the server console instead of emailed. See `app/internal/services/mail/smtp.go` for SMTP setup notes.

## 2. Architecture and Structure

The project uses a SolidJS SPA/SSG frontend with a Go JSON API backend.

*   `app/`: Go backend (REST API, services, storage, database schema).
*   `frontend/`: SolidJS application source.
*   `app/internal/server/handler/`: Thin JSON HTTP handlers.
*   `app/internal/services/`: Core business logic (source of truth).
*   `app/internal/services/financial/`: Provider facade (`FINANCIAL_PROVIDER` switch).
*   `app/internal/services/stripefc/`: Stripe Financial Connections implementation.
*   `app/internal/services/stripebilling/`: Stripe Checkout, Customer Portal, and billing webhooks.
*   `app/internal/services/plaid/`: Plaid implementation (unchanged, selectable via env).
*   `app/database/`: SQLite schema and data.

### Unified connections API

The SPA uses provider-agnostic routes (active provider selected by `FINANCIAL_PROVIDER`):

*   `GET /api/v1/connections/provider`
*   `GET /api/v1/connections`
*   `POST /api/v1/connections/create-session`
*   `POST /api/v1/connections/complete`
*   `POST /api/v1/connections/sync`
*   `POST /api/v1/connections/create-update-session/:id`
*   `POST /api/v1/connections/sync-item/:id`
*   `POST /api/v1/connections/disconnect/:id`
*   `POST /api/v1/connections/toggle-visibility/:id`
*   `POST /api/v1/connections/remove-account/:id`

Only `POST /api/v1/plaid/webhook` remains under the Plaid path (transaction webhooks). Stripe **billing** webhooks are received at `POST /api/v1/stripe/webhook` (unauthenticated; verified via the `Stripe-Signature` header). Subscription checkout uses `POST /api/v1/subscription/checkout`, and the Customer Portal uses `POST /api/v1/subscription/portal`. All SPA bank connection flows use `/api/v1/connections/*`. Set `FINANCIAL_PROVIDER=stripe` to use Stripe Financial Connections instead.

Plaid transaction webhooks are received at `POST /api/v1/plaid/webhook` (unauthenticated; JWT verified via the `Plaid-Verification` header). The endpoint handles `TRANSACTIONS` / `SYNC_UPDATES_AVAILABLE` by syncing the affected item in the background, and `ITEM` / `ERROR` for connection health updates. Link tokens include the webhook URL when `API_PUBLIC_URL` or `PLAID_WEBHOOK_URL` is set.

Frontend entry points: `frontend/src/lib/connections.ts` (facade) and `frontend/src/lib/stripe.ts` (Stripe.js).

## 3. Building and Running

The project uses a `Makefile` to manage common tasks.

*   **Build Backend:** `make build` (syncs shared types and builds the Go binary).
*   **Run Backend:** `make run` (executes the Go server).
*   **Build Frontend:** `cd frontend && npm run build` (builds the SolidJS app).
*   **Capacitor Sync:** `cd frontend && npx cap sync` (syncs the frontend build to the native platforms).

### Production Docker

*   **Configure:** `cp env.prod.example .env` and set `SITE_HOST`, `API_PUBLIC_URL`, `FRONTEND_URL` (all HTTPS URLs), secrets, and `ENV=production` (compose sets this automatically).
*   **Build images:** `make prod-build` (compiled Go binary + static SSG frontend baked into Caddy image).
*   **Start stack:** `make prod-up` (uses `docker-compose.prod.yml`, persistent `db_data` volume, no hot-reload). Two runtime services: **Caddy** (TLS, `/api` proxy, static file server) and **Go API**.
*   **Stop stack:** `make prod-down`
*   **Logs:** `make prod-logs`
*   Set `SITE_HOST=yourdomain.com` (bare hostname, no scheme). Caddy provisions Let's Encrypt automatically when DNS points at the server and ports 80/443 are open. Set `API_PUBLIC_URL` and `FRONTEND_URL` to matching `https://` URLs. Dev stack remains `docker compose up` with `docker-compose.yml`.

> **Note:** For Capacitor to work correctly, the SolidJS app must be configured for a static build (SSG) so that an `index.html` is generated in the web directory (e.g., `dist` or `.output/public`).
> **Note:** Andriod Studio is not installed on this server, so I will be manually testing the application from a separate machine.

## 4. Development Conventions

*   **Efficiency First:** Leverage SolidJS's fine-grained reactivity to minimize re-renders. Optimize backend queries for fast API responses.
*   **Security:** Ensure all API endpoints are properly authenticated. Use secure storage for sensitive data on mobile. Never expose internal IDs.
*   **Mobile-First CSS:** Use Vanilla CSS or CSS Modules for the frontend. Focus on responsive, touch-friendly layouts. **Avoid TailwindCSS**.
*   **Thin Handlers:** Handlers should only parse requests and call services, returning JSON.
*   **Handler → Service Pattern:** All JSON handlers in `app/internal/server/handler/` must delegate business logic to `app/internal/services/`. Handlers must not call `storage` directly. If a handler needs behavior that does not exist yet, add a method to the appropriate service first, then call it from the handler. Handlers are responsible for HTTP concerns only: binding input, status codes, JSON responses, redirects, and cookies.
*   **Available Service Domains:** `auth` (login, register, logout, sessions, password reset), `auth` SSO (Google login, linking, reauth), `dashboard`, `financial`, `stripefc`, `plaid`, `settings`, `subscription`, `tags`, and `transactions`.

## 5. Database Schema

The SQLite schema (`app/database/schema.sql`) includes tables for:
*   `users`: User profiles, theme preferences, subscription tier, `stripe_customer_id`.
*   `user_privileges`: Operator overrides per user (`unlimited_limits`, `stripe_coupon_id`).
*   `user_sso`: OAuth links (Google, Apple).
*   `stripe_fc_items` / `stripe_fc_account` / `stripe_api_usage`: Stripe Financial Connections data.
*   `plaid_items` / `plaid_account` / `plaid_api_usage`: Plaid data (preserved for provider switch).
*   `transactions`: Synced transactions with `provider` column (`stripe` or `plaid`); `plaid_id` holds the internal account row ID for the active provider's account table.
*   `categories` & `tags`: Hierarchical categorization engine.
*   `tag_filters`: User-defined rules for auto-tagging.
*   `sessions`: Authentication session management.

Schema changes go only into `schema.sql` and `test_schema.sql` (no migration scripts).

## General Instructions
- Ensure each function has a comment on what it does
- Follow Go best practices
- Follow Echo best practices
- Always reference the official documentation first
- Always aim to reduce the number of client side refreshes
- Do not use echo.Context and only use *echo.Context
- Any changes made must be documented in `summary.md`
- Do not delete anything from `summary.md` only update it with what you currently did and only if you changed that specific functionality.
- Make sure to use proper status codes when encountering an error and that its properly managed by the frontend
- Make sure to read other files to understand the design patterns
- Only use comments when necessary
- The frontend must be dynamically sized to fit any screen properly
- Always prefer passing by reference to reduce the amount of copying.
- Do not read the .env file.
- Ensure the frontend is dynamic
- Please run "make build" to ensure all the configs are properly generated
- Do not duplicate any code and reuse as much as you can.
- Make sure to add any icons to the icons folder.
- The models folder only contains structs and in the models folder, there is an external folder and it only contains structs to be exported to the frontend
- Stripe Go SDK: `github.com/stripe/stripe-go/v86`
- Frontend Stripe.js: `@stripe/stripe-js` via `collectFinancialConnectionsAccounts`
