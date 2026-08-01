# Development Setup

## Overview

Use this page to run FinancialTracker on a developer machine. Docker Compose is the recommended path. You can also run the Go API and SolidStart frontend as separate processes.

## Requirements

All environments need:

- Go 1.26 or later with CGO enabled (SQLite)
- Node.js 22 or later and npm
- A C compiler when you build the Go binary locally

Docker development also needs Docker, Docker Compose, and a `.env` file at the repository root.

See [Environment variables](./environment.md) for the full variable list.

## Quick start with Docker

1. Clone the repository.
2. Create a `.env` file at the repository root.
3. Set at least these values:

```bash
ENCRYPTION_KEY=12345678901234567890123456789012   # exactly 32 bytes
PLAID_CLIENT_ID=your_plaid_client_id
PLAID_SANDBOX_SECRET=your_plaid_sandbox_secret
API_PUBLIC_URL=http://localhost
GOOGLE_CLIENT_ID=                          # optional
GOOGLE_CLIENT_SECRET=
SUBSCRIPTIONS_ENABLED=false                # optional
```

4. Start the stack:

```bash
docker compose up
```

Services:

- Caddy reverse proxy on ports 80 and 443
- Go API with hot reload (Air)
- SolidStart frontend with Vite HMR

5. Open the URL in `API_PUBLIC_URL`. Update frontend env vars in `docker-compose.yml` if you use a custom hostname.

### Development login

When `ENV=development`, the backend resets the database on startup and seeds this user:

| Email | Password |
|-------|----------|
| `test@test.com` | `test` |

This user is a registration admin when subscriptions are disabled. The user can create invite codes in Settings → Account.

## Local run without Docker

### Backend

1. Export required variables from the repository root.
2. Build and run the API.

```bash
export ENV=development
export ENCRYPTION_KEY=12345678901234567890123456789012
export PLAID_CLIENT_ID=...
export PLAID_SANDBOX_SECRET=...
export API_PUBLIC_URL=http://localhost:8080

make build
make run
```

### Frontend

1. Open a second terminal.
2. Install dependencies and start the dev server.

```bash
cd frontend
npm install
export API_PUBLIC_URL=http://localhost:8080
export FRONTEND_URL=http://localhost:3000
export VITE_API_PUBLIC_URL=http://localhost:8080
export VITE_FRONTEND_URL=http://localhost:3000
npm run dev
```

3. Open `http://localhost:3000`.

CORS and cookies are simplest when Caddy serves API and frontend on the same origin in Docker.

## Project structure

```
FinancialTracker/
├── app/                    # Go REST API
│   ├── database/           # schema.sql, test_schema.sql, migrations/
│   ├── internal/
│   │   ├── server/         # HTTP routes and handlers
│   │   └── services/       # Business logic
│   └── main.go
├── frontend/               # SolidStart SPA
│   └── src/
│       ├── routes/         # Pages
│       ├── components/
│       └── lib/            # API clients, types
├── docs/                   # This documentation
├── docker-compose.yml      # Development stack
├── docker-compose.prod.yml # Production stack
├── env.prod.example        # Production env template
├── Makefile
└── Caddyfile               # Dev reverse proxy
```

## Related docs

- [Environment variables](./environment.md)
- [Production deployment](./production.md)
- [Makefile commands](../operations/makefile.md)
- [Testing](../operations/testing.md)
- [Architecture overview](../architecture/overview.md)
