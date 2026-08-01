# Makefile Commands

## Overview

The repository `Makefile` provides common build, run, and production targets.

## Commands

| Command | Description |
|---------|-------------|
| `make build` | Sync Go types to TypeScript and compile the Go binary to `bin/main` |
| `make run` | Build and run the API locally |
| `make sync-types` | Regenerate `frontend/src/lib/types.ts` from Go models |
| `make registration-code` | Issue a bootstrap invite code in the production backend container |
| `make prod-build` | Build production Docker images |
| `make prod-build-amd64` | Build production images for `linux/amd64` |
| `make prod-export` | Save proxy and backend images to `dist/ft-prod-images.tar.gz` |
| `make prod-import` | Load images from `dist/ft-prod-images.tar.gz` |
| `make prod-deploy` | Build, export, copy, load, and restart on a VPS (`VPS_HOST`, `VPS_PATH`) |
| `make prod-up` | Start the production stack |
| `make prod-up-runtime` | Start from pre-loaded images only |
| `make prod-down` | Stop the production stack |
| `make prod-logs` | Tail production logs |

## Invite code bootstrap

When `SUBSCRIPTIONS_ENABLED=false` and no admin session exists:

1. Start the production backend container against the live database volume.
2. Run `make registration-code`.
3. Use the printed code to register the first admin account.

Rebuild images after upgrades if the command is missing from an older image.

## Related docs

- [Development setup](../getting-started/development.md)
- [Production deployment](../getting-started/production.md)
- [Authentication](../features/authentication.md)
