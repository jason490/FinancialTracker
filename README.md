# Financial Tracker

A self-hosted financial tracking application inspired by Rocket Money. Connect bank accounts, sync transactions, apply smart tags, and customize a dashboard from a SolidJS web app. An optional Android build uses Capacitor.

**Repository:** [github.com/jason490/FinancialTracker](https://github.com/jason490/FinancialTracker/tree/master)

---

## Documentation

Full documentation lives in [`docs/`](docs/index.md). Start with these pages:

| Topic | Document |
|-------|----------|
| Documentation index | [`docs/index.md`](docs/index.md) |
| Development setup | [`docs/getting-started/development.md`](docs/getting-started/development.md) |
| Production deployment | [`docs/getting-started/production.md`](docs/getting-started/production.md) |
| Environment variables | [`docs/getting-started/environment.md`](docs/getting-started/environment.md) |
| Mobile (Android) | [`docs/getting-started/mobile.md`](docs/getting-started/mobile.md) |
| Architecture | [`docs/architecture/overview.md`](docs/architecture/overview.md) |
| API reference | [`docs/architecture/api.md`](docs/architecture/api.md) |
| Changelog | [`docs/operations/changelog.md`](docs/operations/changelog.md) |

### Feature guides

| Feature | Document |
|---------|----------|
| Authentication | [`docs/features/authentication.md`](docs/features/authentication.md) |
| Onboarding | [`docs/features/onboarding.md`](docs/features/onboarding.md) |
| Bank connections | [`docs/features/bank-connections.md`](docs/features/bank-connections.md) |
| Transactions | [`docs/features/transactions.md`](docs/features/transactions.md) |
| Tags and categories | [`docs/features/tags-and-categories.md`](docs/features/tags-and-categories.md) |
| Dashboard | [`docs/features/dashboard.md`](docs/features/dashboard.md) |
| Settings | [`docs/features/settings.md`](docs/features/settings.md) |
| Subscriptions | [`docs/features/subscriptions.md`](docs/features/subscriptions.md) |

---

## Features (summary)

- Email/password auth, Google SSO, password reset, and optional invite-only registration
- Plaid (default) or Stripe Financial Connections through `/api/v1/connections/*`
- Transaction list, filters, bulk tags, ZIP/CSV export, and filter-aware graphs
- Categories, tags, auto-tag rules, and exclude-from-totals tags
- Customizable dashboard widgets with separate desktop and mobile layouts
- Optional Stripe subscriptions with Free / Plus / Premium limits
- Settings for profile, connections, appearance, tags, plan, and data export
- Responsive web UI and Capacitor-ready SSG build

See the feature guides above for full behavior, endpoints, and configuration.

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

## Quick start

### Docker (recommended)

1. Clone the repository.
2. Create a `.env` file. See [`docs/getting-started/environment.md`](docs/getting-started/environment.md).
3. Run `docker compose up`.
4. Open the URL in `API_PUBLIC_URL`.

Development seeds `test@test.com` / `test` when `ENV=development`.

Full steps: [`docs/getting-started/development.md`](docs/getting-started/development.md).

### Production

```bash
cp env.prod.example .env
# edit SITE_HOST, API_PUBLIC_URL, FRONTEND_URL, secrets
make prod-build
make prod-up
```

Full steps and VPS image transfer: [`docs/getting-started/production.md`](docs/getting-started/production.md).

### Common make targets

| Command | Description |
|---------|-------------|
| `make build` | Sync types and compile `bin/main` |
| `make run` | Build and run the API locally |
| `make prod-build` / `make prod-up` | Build and start production images |
| `make registration-code` | Bootstrap an invite code |

Full list: [`docs/operations/makefile.md`](docs/operations/makefile.md).

---

## Testing

```bash
# Go
cd app && go test ./...

# Playwright E2E
cd e2e && pip install -r requirements.txt && pytest
```

Details: [`docs/operations/testing.md`](docs/operations/testing.md).

---

## Contributor notes

Coding conventions for agents and contributors are in [`GEMINI.md`](GEMINI.md).

Document behavior changes in [`docs/operations/changelog.md`](docs/operations/changelog.md).

---

## License

See the repository for license information.
