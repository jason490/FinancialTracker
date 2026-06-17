# FinancialTracker Project Context

This project is a high-performance, user-centric financial tracking application inspired by Rocket Money. It is designed for blazing-fast performance and seamless expense management by directly syncing with bank accounts via Plaid.
You are Full-stack developer and designer. Ensure any interactions feel seemless on both andriod, webapps, and websites. Any interactions should be properly organized (depending on the amount of visual bloat seen) into sidebar, navbar (top for webapp, bottom for mobile), modals, etc.

## 1. Project Overview

*   **Core Purpose:** Automated expense tracking, smart tagging, and categorization.
*   **Main Technologies:**
    *   **Backend:** Go (Golang) using the `echo` v5 framework (Transitioning to a pure REST API).
    *   **Database:** SQLite (embedded relational database).
    *   **Frontend (New):** SolidStart (High-performance, reactive UI library) with SSG (no SSR).
    *   **Mobile Bridge:** Capacitor (Targeting Android and Web platforms).
    *   **Styling:** Vanilla CSS / CSS Modules (Avoiding TailwindCSS for better mobile optimization).
    *   **Legacy Frontend:** Templ and HTMX (Being phased out).
    *   **External APIs:** Plaid (for bank account and transaction syncing).

## 2. Architecture and Structure

The project is transitioning from a server-side rendered (Templ/HTMX) architecture to a modern SPA (SolidJS) with a Go API backend.

*   `app/`: Current legacy codebase (Backend + Templ/HTMX frontend).
*   `frontend/`: New SolidJS application source.
*   `app/internal/server/`: Transitioning handlers to return JSON instead of HTML fragments.
*   `app/internal/services/`: Core business logic (remains the source of truth).
*   `app/database/`: SQLite schema and data.

## 3. Building and Running

The project uses a `Makefile` to manage common tasks.

*   **Build Backend:** `make build` (Builds the Go binary).
*   **Run Backend:** `make run` (Executes the Go server).
*   **Build Frontend:** `cd frontend && npm run build` (Builds the SolidJS app).
*   **Capacitor Sync:** `cd frontend && npx cap sync` (Syncs the frontend build to the native platforms).

> **Note:** For Capacitor to work correctly, the SolidJS app must be configured for a static build (SSG) so that an `index.html` is generated in the web directory (e.g., `dist` or `.output/public`).
> **Note:** Andriod Studio is not installed on this server, so I will be manually testing the application from a separate machine.

## 4. Development Conventions

*   **Efficiency First:** Leverage SolidJS's fine-grained reactivity to minimize re-renders. Optimize backend queries for fast API responses.
*   **Security:** Ensure all API endpoints are properly authenticated. Use secure storage for sensitive data on mobile. Never expose internal IDs.
*   **Mobile-First CSS:** Use Vanilla CSS or CSS Modules for the new frontend. Focus on responsive, touch-friendly layouts. **Avoid TailwindCSS** for the new app.
*   **Thin Handlers:** Handlers should only parse requests and call services, returning JSON for the new frontend.
*   **Handler → Service Pattern:** All JSON handlers in `app/internal/server/handler/` must delegate business logic to `app/internal/services/`. Handlers must not call `storage` directly. If a handler needs behavior that does not exist yet, add a method to the appropriate service first, then call it from the handler. Handlers are responsible for HTTP concerns only: binding input, status codes, JSON envelopes, redirects, and cookies.
*   **Available Service Domains:** `auth` (login, register, logout, sessions, password reset), `auth` SSO (Google login, linking, reauth), `dashboard`, `plaid`, `settings`, `tags`, and `transactions`.
*   **Transition Strategy:** Maintain the legacy Templ/HTMX UI within `app/` until the SolidJS app in `frontend/` reaches feature parity.

## 5. Database Schema

The SQLite schema (`app/database/schema.sql`) includes tables for:
*   `users`: User profiles and theme preferences.
*   `user_sso`: OAuth links (Google, Apple).
*   `accounts`: Bank accounts fetched via Plaid.
*   `transactions`: Synced transactions.
*   `categories` & `tags`: Hierarchical categorization engine.
*   `tag_filters`: User-defined rules for auto-tagging.
*   `sessions`: Authentication session management.

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
