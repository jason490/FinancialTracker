# Production Deployment

## Overview

Production runs a compiled Go API and a static SolidStart build behind Caddy. Caddy obtains and renews Let's Encrypt certificates when DNS points at the server and ports 80 and 443 are open.

## Requirements

- Docker and Docker Compose
- A domain with DNS that points at your server
- Secrets filled from [`env.prod.example`](../../env.prod.example)

## Deploy on one host

1. Copy the production env template:

```bash
cp env.prod.example .env
```

2. Set `SITE_HOST`, `API_PUBLIC_URL`, and `FRONTEND_URL` to matching `https://` values.
3. Set `ENCRYPTION_KEY`, provider secrets, Google OAuth secrets if used, and SMTP for password reset.
4. Set `REGISTRATION_ADMIN_EMAILS` when you use invite-only registration.
5. Build and start:

```bash
make prod-build
make prod-up
```

6. View logs or stop the stack:

```bash
make prod-logs
make prod-down
```

Rebuild with `make prod-build` after you change `API_PUBLIC_URL` or `FRONTEND_URL`. Those values are baked into the static frontend.

## Remote deploy (low-RAM VPS)

Build images on a machine with enough RAM. Transfer the images to the VPS. The VPS loads images and starts containers. The VPS does not run `npm` or `go build`.

### Build machine `.env`

Create `.env` from `env.prod.example`. Set public URLs to the same values the VPS will use. Docker Compose reads these during `make prod-build-amd64`:

| Variable | Required to build? | Why |
|----------|-------------------|-----|
| `API_PUBLIC_URL` | Yes | Baked into the frontend as `VITE_API_PUBLIC_URL` |
| `FRONTEND_URL` | Yes | Baked into the frontend as `VITE_FRONTEND_URL` |
| `SITE_HOST`, `ACME_EMAIL` | No | Used when containers start |
| Secrets (encryption, Plaid, Google, SMTP) | No | Needed only at runtime on the VPS |

Example minimal build `.env`:

```bash
API_PUBLIC_URL=https://yourdomain.com
FRONTEND_URL=https://yourdomain.com
```

If you change those URLs, rebuild and redeploy the images.

### VPS `.env` (runtime)

Copy `env.prod.example` to `.env` on the VPS. Fill every runtime value.

| Variable | Required on VPS? | Notes |
|----------|------------------|-------|
| `SITE_HOST` | Yes | Bare hostname for Let's Encrypt |
| `API_PUBLIC_URL`, `FRONTEND_URL` | Yes | Must match the frontend build |
| `ACME_EMAIL` | Recommended | Certificate expiry notices |
| `ENCRYPTION_KEY` | Yes | Exactly 32 bytes |
| `PLAID_ENV` | Yes for live Plaid | Set `production` explicitly |
| `PLAID_CLIENT_ID`, `PLAID_PROD_SECRET` | Yes for live Plaid | Use production secret when `PLAID_ENV=production` |
| Google OAuth vars | If you use Google sign-in | |
| `MAIL_SMTP_*`, `MAIL_FROM` | Yes for password reset | |
| Subscription / Stripe vars | As needed | |

After you edit `.env` on the VPS, restart containers with `make prod-up-runtime`. Rebuild only when URLs change.

### One-time VPS setup

Install Docker and Compose. Check out the repository, or at least these files:

- `docker-compose.prod.yml`
- `docker-compose.prod.runtime.yml`
- `Caddyfile.prod`
- `.env`

Do not run `make prod-build` on the VPS.

### Each release from the build machine

```bash
make prod-build-amd64
make prod-export
scp dist/ft-prod-images.tar.gz user@vps:/tmp/
```

### On the VPS

```bash
cd /path/to/FinancialTracker
make prod-import PROD_EXPORT_ARCHIVE=/tmp/ft-prod-images.tar.gz
make prod-up-runtime
```

Or run the full pipeline:

```bash
make prod-deploy VPS_HOST=user@vps VPS_PATH=/path/to/FinancialTracker
```

`prod-up-runtime` uses `docker-compose.prod.runtime.yml` so Compose cannot trigger a build on the server.

## Related docs

- [Environment variables](./environment.md)
- [Makefile commands](../operations/makefile.md)
- [Security](../architecture/security.md)
- [Database and migrations](../architecture/database.md)
