# FinancialTracker Documentation

This folder holds product and operator documentation for FinancialTracker.

These pages follow ASD-STE100 style: active voice, short sentences, one topic per sentence, and consistent terms.

## Start here

| Document | Purpose |
|----------|---------|
| [Development setup](./getting-started/development.md) | Run the app locally or with Docker Compose |
| [Production deployment](./getting-started/production.md) | Build and run the production stack |
| [Environment variables](./getting-started/environment.md) | Configure the API, frontend, and integrations |
| [Mobile (Android)](./getting-started/mobile.md) | Build the Capacitor Android app |

## Architecture

| Document | Purpose |
|----------|---------|
| [Architecture overview](./architecture/overview.md) | Stack, layout, and request flow |
| [API reference](./architecture/api.md) | REST route groups and auth |
| [Database and migrations](./architecture/database.md) | SQLite schema and migration rules |
| [Security](./architecture/security.md) | Sessions, CSRF, encryption, hardening |

## Features

| Document | Purpose |
|----------|---------|
| [Authentication](./features/authentication.md) | Login, SSO, password reset, invites |
| [Onboarding](./features/onboarding.md) | New-user wizard |
| [Bank connections](./features/bank-connections.md) | Plaid and Stripe Financial Connections |
| [Transactions](./features/transactions.md) | List, filters, export, graphs |
| [Tags and categories](./features/tags-and-categories.md) | Tagging rules and Settings → Tags |
| [Dashboard](./features/dashboard.md) | Widgets, layouts, charts |
| [Settings](./features/settings.md) | Account, themes, data, tabs |
| [Subscriptions](./features/subscriptions.md) | Plans, Stripe billing, limits |

## Operations

| Document | Purpose |
|----------|---------|
| [Makefile commands](./operations/makefile.md) | Common make targets |
| [Testing](./operations/testing.md) | Go tests and Playwright E2E |
| [Changelog](./operations/changelog.md) | Notable changes and fixes |

## Contributor notes

For coding conventions used by agents and contributors, see [`GEMINI.md`](../GEMINI.md) in the repository root.
