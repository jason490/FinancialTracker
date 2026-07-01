# Financial Tracker

A self-hosted financial tracking application inspired by Rocket Money. Connect bank accounts, sync transactions automatically, apply smart tags, and customize a dashboard—all from a fast SolidJS web app with an optional Android build via Capacitor.

**Repository:** [github.com/jason490/FinancialTracker](https://github.com/jason490/FinancialTracker/tree/master)

---

## Features

### Authentication & accounts
- Email/password sign-up and sign-in with secure session cookies
- Google SSO (sign-in, sign-up, and account linking from Settings)
- Password reset via temporary 6-digit codes (email in production; logged to console in development)
- Onboarding wizard for new users (welcome → plan → connect bank)
- Account deletion with re-authentication
- **Invite-only registration** when `SUBSCRIPTIONS_ENABLED=false` — new users need a single-use admin code (48-hour TTL)

### Bank connections
- **Plaid** (default) or **Stripe Financial Connections** — switch with `FINANCIAL_PROVIDER`
- Link institutions, sync transactions, reconnect in update mode
- Per-item sync, disconnect, hide accounts, remove accounts
- Plaid transaction webhooks for background sync
- Provider-agnostic API at `/api/v1/connections/*`

### Transactions & tagging
- Paginated transaction list with search and filters
- Bulk add/remove tags
- CSV export
- Categories and tags with hierarchical organization
- Auto-tagging rules: exact match, regex, and amount filters
- Default tag seeding for new users

### Dashboard
- Customizable widget grid (desktop and mobile layouts stored separately)
- Widgets: net worth, cash flow, spending trend, recent transactions, tag breakdowns, and more
- Drag-and-drop edit mode (click on desktop, long-press on mobile)
- ApexCharts visualizations with theme-aware styling

### Subscriptions & billing (optional)
- Toggle with `SUBSCRIPTIONS_ENABLED` (default `true`)
- Stripe Checkout, Customer Portal, and billing webhooks
- Free / Plus / Premium tiers with per-tier bank and API limits
- Operator overrides via `SUBSCRIPTION_OVERRIDES` and `user_privileges`

### Settings & appearance
- Profile, password, and SSO management
- Theme presets: Default (light), Dark, Tokyo Night, Coffee (+ system preference)
- Connections panel, data export, plan management (when subscriptions enabled)
- **Invite codes** panel for registration admins (when subscriptions disabled)

### Platform
- Responsive web UI (mobile bottom nav, desktop top nav)
- Capacitor-ready static build (SSG) for Android
- Production Docker stack with Caddy (automatic HTTPS) + Go API

---

## Tech stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.26, Echo v5, SQLite |
| Frontend | SolidStart (SolidJS), Vinxi, SSG |
| Mobile | Capacitor (Android) |
| Styling | Vanilla CSS / CSS Modules |
| Bank sync | Plaid or Stripe Financial Connections |
| Billing | Stripe (optional) |
| Auth | Sessions, bcrypt, Google OAuth |
| Proxy / TLS | Caddy |

---

## Requirements

### All environments
- **Go** 1.26+ (CGO enabled — needed for SQLite)
- **Node.js** 22+ and npm
- **GCC / musl** (or equivalent) when building the Go binary with CGO locally

### Docker development
- Docker and Docker Compose
- A `.env` file at the repo root (see [Environment variables](#environment-variables))

### Production
- Docker and Docker Compose
- Domain with DNS pointing at your server (ports **80** and **443** open)
- Secrets filled in from [`env.prod.example`](env.prod.example)

### Optional integrations
| Integration | When needed |
|-------------|-------------|
| Plaid | `FINANCIAL_PROVIDER=plaid` (default) — `PLAID_CLIENT_ID` + sandbox or production secret |
| Stripe FC | `FINANCIAL_PROVIDER=stripe` — Stripe sandbox/prod keys |
| Stripe billing | `SUBSCRIPTIONS_ENABLED=true` — Stripe keys, price IDs, webhook secret |
| Google OAuth | SSO — `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` |
| SMTP | Production password reset — `MAIL_SMTP_*`, `MAIL_FROM` |

---

## Quick start (Docker — recommended for development)

1. **Clone the repository**
   ```bash
   git clone https://github.com/jason490/FinancialTracker.git
   cd FinancialTracker
   ```

2. **Create a `.env` file** at the repo root with at least:
   ```bash
   ENCRYPTION_KEY=12345678901234567890123456789012   # exactly 32 bytes
   PLAID_CLIENT_ID=your_plaid_client_id
   PLAID_SANDBOX_SECRET=your_plaid_sandbox_secret
   API_PUBLIC_URL=http://localhost          # or your LAN hostname
   GOOGLE_CLIENT_ID=                          # optional, for Google SSO
   GOOGLE_CLIENT_SECRET=
   SUBSCRIPTIONS_ENABLED=false                # optional; disables billing + enables invite-only signup
   ```

3. **Start the stack**
   ```bash
   docker compose up
   ```
   Services: **Caddy** (proxy on 80/443), **Go API** (hot-reload via Air), **SolidStart frontend** (Vite HMR).

4. **Open the app** at the URL set in `API_PUBLIC_URL` (update `docker-compose.yml` frontend env vars if you use a custom hostname).

### Development login (auto-seeded)

When `ENV=development`, the backend resets the database on startup and seeds:

| Email | Password |
|-------|----------|
| `test@test.com` | `test` |

This user is also a **registration admin** when subscriptions are disabled (can generate invite codes in Settings → Account).

---

## Quick start (local, without Docker)

**Backend**
```bash
# From repo root
export ENV=development
export ENCRYPTION_KEY=12345678901234567890123456789012
export PLAID_CLIENT_ID=...
export PLAID_SANDBOX_SECRET=...
export API_PUBLIC_URL=http://localhost:8080

make build
make run
```

**Frontend** (separate terminal)
```bash
cd frontend
npm install
export API_PUBLIC_URL=http://localhost:8080
export FRONTEND_URL=http://localhost:3000
export VITE_API_PUBLIC_URL=http://localhost:8080
export VITE_FRONTEND_URL=http://localhost:3000
npm run dev
```

Open `http://localhost:3000`. Ensure the frontend can reach the API (CORS and cookies are simplest when both are served behind the same origin via Caddy in Docker).

---

## Production deployment

1. Copy and edit the production env template:
   ```bash
   cp env.prod.example .env
   ```
   Set `SITE_HOST`, `API_PUBLIC_URL`, `FRONTEND_URL` (all `https://`), `ENCRYPTION_KEY`, Plaid/Stripe/Google/SMTP secrets, and `REGISTRATION_ADMIN_EMAILS` if using invite-only registration.

2. Build and start:
   ```bash
   make prod-build
   make prod-up
   ```

3. **Logs / stop**
   ```bash
   make prod-logs
   make prod-down
   ```

Caddy obtains and renews Let's Encrypt certificates automatically when DNS is correct. Re-run `make prod-build` after changing `API_PUBLIC_URL` or `FRONTEND_URL` (baked into the static frontend).

### Remote deploy (low-RAM VPS)

Build images on a machine with enough RAM, then transfer them to the VPS. The VPS only loads images and starts containers — it never runs `npm` or `go build`.

**One-time VPS setup:** Docker, Compose, repo checkout, `.env` from `env.prod.example`. Do not run `make prod-build` on the VPS.

**Each release from your dev machine:**

```bash
# Set API_PUBLIC_URL / FRONTEND_URL in .env first (baked into the frontend)
make prod-build-amd64
make prod-export
scp dist/ft-prod-images.tar.gz user@vps:/tmp/
```

**On the VPS:**

```bash
cd /path/to/FinancialTracker
make prod-import PROD_EXPORT_ARCHIVE=/tmp/ft-prod-images.tar.gz
make prod-up-runtime
```

Or run the full pipeline in one step (set `VPS_HOST` and `VPS_PATH`):

```bash
make prod-deploy VPS_HOST=user@vps VPS_PATH=/path/to/FinancialTracker
```

`prod-up-runtime` uses [docker-compose.prod.runtime.yml](docker-compose.prod.runtime.yml) so Compose cannot trigger a build on the server.

---

## Makefile commands

| Command | Description |
|---------|-------------|
| `make build` | Sync Go → TypeScript types and compile the Go binary to `bin/main` |
| `make run` | Build and run the API locally |
| `make sync-types` | Regenerate `frontend/src/lib/types.ts` from Go models |
| `make registration-code` | Issue a bootstrap invite code in the production backend container (`SUBSCRIPTIONS_ENABLED=false`) |
| `make prod-build` | Build production Docker images |
| `make prod-build-amd64` | Build production images for `linux/amd64` (VPS deploy) |
| `make prod-export` | Save proxy + backend images to `dist/ft-prod-images.tar.gz` |
| `make prod-import` | Load images from `dist/ft-prod-images.tar.gz` (run on VPS) |
| `make prod-deploy` | Build, export, `scp`, load, and restart on VPS (`VPS_HOST`, `VPS_PATH`) |
| `make prod-up` | Start production stack |
| `make prod-up-runtime` | Start stack from pre-loaded images only (VPS) |
| `make prod-down` | Stop production stack |
| `make prod-logs` | Tail production logs |

---

## Invite-only registration

When `SUBSCRIPTIONS_ENABLED=false`, new sign-ups (email or Google) require a one-time code.

**Admins**
- Set `REGISTRATION_ADMIN_EMAILS=admin@example.com;ops@example.com` (semicolon-separated)
- In development, `test@test.com` is always an admin
- Signed-in admins: **Settings → Account → Invite codes**, or `POST /api/v1/admin/registration-codes`
- First deploy without an admin session: `make registration-code` (runs inside the production backend container against the live database volume; rebuild images after upgrading if the command is missing)

**Codes**
- 8-character alphanumeric, single-use, **48-hour** expiry
- Invalid after use or expiry (rows remain in the database but are ignored)

---

## Environment variables

### Core

| Variable | Default | Description |
|----------|---------|-------------|
| `ENV` | production | `development` — resets DB, seeds test user, relaxed rate limits, console password-reset codes |
| `PORT` | `8080` | API listen port |
| `DATABASE_PATH` | `./database/main.db` | SQLite file path |
| `ENCRYPTION_KEY` | — | **Required in production.** Exactly 32 bytes for AES-256 (Plaid token encryption, OAuth state) |
| `API_PUBLIC_URL` | — | Public API URL (e.g. `https://app.example.com`) |
| `FRONTEND_URL` | — | Public SPA URL (usually same origin as API in production) |

### Auth

| Variable | Description |
|----------|-------------|
| `GOOGLE_CLIENT_ID` | Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret |
| `REGISTRATION_ADMIN_EMAILS` | Semicolon-separated emails allowed to create invite codes |

### Subscriptions

| Variable | Default | Description |
|----------|---------|-------------|
| `SUBSCRIPTIONS_ENABLED` | `true` | `false` — no Stripe billing, unlimited usage, invite-only registration |
| `STRIPE_SANDBOX_SECRET` | — | Stripe secret key (development) |
| `STRIPE_PROD_SECRET` | — | Stripe secret key (production) |
| `STRIPE_SANDBOX_PUBLISHABLE` | — | Stripe publishable key (development) |
| `STRIPE_PUBLISHABLE` | — | Stripe publishable key (production) |
| `STRIPE_PRICE_PLUS` / `STRIPE_PRICE_PREMIUM` | — | Production price IDs |
| `STRIPE_SANDBOX_PRICE_PLUS` / `STRIPE_SANDBOX_PRICE_PREMIUM` | — | Sandbox price IDs |
| `STRIPE_WEBHOOK_SECRET` / `STRIPE_SANDBOX_WEBHOOK_SECRET` | — | Billing webhook signing secrets |
| `SUBSCRIPTION_OVERRIDES` | — | Per-email overrides, e.g. `user@example.com=unlimited,coupon=CODE` |
| `SUBSCRIPTION_OVERRIDE_USER_IDS` | — | Same format keyed by user ID |

### Bank provider

| Variable | Default | Description |
|----------|---------|-------------|
| `FINANCIAL_PROVIDER` | `plaid` | `plaid` or `stripe` |
| `PLAID_CLIENT_ID` | — | Plaid application ID |
| `PLAID_ENV` | `sandbox` (dev) | `sandbox` or `production` |
| `PLAID_SANDBOX_SECRET` | — | Plaid sandbox secret |
| `PLAID_PROD_SECRET` | — | Plaid production secret |
| `PLAID_WEBHOOK_URL` | — | Optional; defaults to `{API_PUBLIC_URL}/api/v1/plaid/webhook` |

### Email (production password reset)

| Variable | Description |
|----------|-------------|
| `MAIL_SMTP_HOST` | SMTP server |
| `MAIL_SMTP_PORT` | e.g. `587` |
| `MAIL_SMTP_USER` | SMTP username |
| `MAIL_SMTP_PASSWORD` | SMTP password or API key |
| `MAIL_FROM` | Sender address |
| `MAIL_FROM_NAME` | Optional display name |

When `MAIL_SMTP_HOST` is unset and `ENV=development`, reset codes are printed to the API logs.

### Production Docker (Caddy)

| Variable | Description |
|----------|-------------|
| `SITE_HOST` | Bare hostname for TLS (e.g. `app.example.com`) |
| `ACME_EMAIL` | Optional Let's Encrypt registration email |

### Frontend build-time (Vinxi / Vite)

| Variable | Description |
|----------|-------------|
| `VITE_API_PUBLIC_URL` | API URL embedded in the client bundle |
| `VITE_FRONTEND_URL` | Frontend URL for OAuth return paths |

---

## Mobile (Capacitor / Android)

1. Build the static frontend:
   ```bash
   cd frontend
   npm run build
   ```
2. Sync to the native project:
   ```bash
   npx cap sync
   ```
3. Open in Android Studio on a machine with the Android SDK and run on a device or emulator.

The app must be built as SSG so `index.html` and assets land in the web directory Capacitor serves.

---

## Testing

**End-to-end (Playwright)** — from the `e2e/` directory:
```bash
pip install -r requirements.txt
pytest
```

**Go unit tests**
```bash
cd app
go test ./...
```

---

## Project structure

```
FinancialTracker/
├── app/                    # Go REST API
│   ├── database/           # schema.sql, test_schema.sql
│   ├── internal/
│   │   ├── server/         # HTTP routes and handlers
│   │   └── services/       # Business logic
│   └── main.go
├── frontend/               # SolidStart SPA
│   └── src/
│       ├── routes/         # Pages
│       ├── components/
│       └── lib/            # API clients, types
├── docker-compose.yml      # Development stack
├── docker-compose.prod.yml # Production stack
├── env.prod.example        # Production env template
├── Makefile
└── Caddyfile               # Dev reverse proxy
```

---

## API overview

| Area | Base path |
|------|-----------|
| Auth | `/api/v1/auth/*` |
| Dashboard | `/api/v1/dashboard` |
| Transactions | `/api/v1/transactions` |
| Tags & categories | `/api/v1/tags`, `/api/v1/categories` |
| Bank connections | `/api/v1/connections/*` |
| Settings | `/api/v1/settings` |
| Subscription | `/api/v1/subscription` |
| Admin (invite codes) | `/api/v1/admin/registration-codes` |
| Health | `/health` |

Authenticated routes use the `Session` HTTP-only cookie set at login.

---

## Further reading

- [`GEMINI.md`](GEMINI.md) — detailed architecture and conventions for contributors
- [`summary.md`](summary.md) — changelog and implementation notes

---

## License

See the repository for license information.
