# Project Summary

## Current Status: Architectural Transition
The project is currently transitioning from a server-side rendered (Templ/HTMX/Tailwind) architecture to a high-performance, mobile-first stack using **SolidJS** and **Capacitor** (Android/Web). The Go backend is being refactored into a pure REST API.

## Completed Milestones

### Core Backend & Services
- **Thin Handler Refactor**: Decoupled HTTP logic from business logic across all domains (Auth, Plaid, Tags, Transactions).
- **Service Layer**: Implemented robust services for `Auth`, `Plaid`, `Tags`, `Transactions`, `Settings`, and `Dashboard`.
- **Model Organization**: Refactored `internal/models/models.go` into modular files (`user.go`, `session.go`, `plaid.go`, `transaction.go`, `category.go`, `dashboard.go`, `page.go`) for improved maintainability.
- **Routing Refactor**: Isolated legacy Templ/HTMX routes into a dedicated `legacy()` method in `routes.go` and renamed the handlers package to `legacyHandler` to clearly distinguish legacy logic from the new JSON API.
- **Unified Validation**: Centralized input validation (Email, Password, Name) in the service layer.
- **Plaid Integration**: 
    - Cursor-based transaction synchronization.
    - Intelligent account reconciliation (soft-disconnection, visibility toggles).
    - Support for Plaid Update Mode and item health monitoring.
- **Smart Tagging Engine**: Rule-based auto-tagging with support for exact matches, regex, and amount-based filters. Default tags are now assigned distinct colors from the tag palette during initial seeding (e.g., Dining Out is Rose, Groceries is Emerald, Bills is Red).
- **Security**: Implemented session-based auth, re-authentication for high-privilege actions, and strict ownership checks.

### Database & Storage
- **SQLite Schema**: Optimized schema for users, sessions, Plaid items, accounts, transactions, and tags.
- **Foreign Key Enforcement**: Enabled SQLite foreign keys for automated data cleanup (ON DELETE CASCADE).
- **Multi-Auth Support**: Support for multiple SSO providers and account linking.

### Legacy UI (Templ/HTMX)
- **ViewPoint Dashboard**: Modular, customizable dashboard with Chart.js visualizations.
- **Advanced Transactions**: Hierarchical category filtering, bulk tagging, and responsive design.
- **Notification System**: Event-driven toast notifications using Alpine.js and HTMX triggers.
- **Theme Support**: Persistent light/dark/auto theme modes.

### Development Environment
- **Dockerization**: Containerized the Go backend and SolidStart frontend for a consistent development experience.
- **Hot-Reloading**: Integrated `air` for Go and standard Vite/Vinxi HMR for the frontend.
- **Volume Optimization**: Configured selective volume syncing to ensure code changes are reflected instantly while excluding heavy assets like `node_modules` and build artifacts from the sync loop.
- **Lightweight Images**: Leveraged Alpine-based images for both Node and Go to minimize resource footprint.

## New Architectural Direction (In Progress)
- **Frontend**: Shifting to **SolidJS** for fine-grained reactivity and maximum efficiency.
- **Mobile**: Using **Capacitor** to target **Android** and **Web** platforms.
    - Initialized Capacitor in the `frontend` directory.
    - Added the Android platform.
    - Updated `GEMINI.md` with mobile-specific development commands.
- **Styling**: Moving to **Vanilla CSS / CSS Modules** to ensure lightweight, touch-friendly mobile performance (removing TailwindCSS).
- **Security & Efficiency**: Prioritizing minimal re-renders, secure mobile storage, and optimized API interaction.
- **Unified Session Management**: Refactored the Go backend to manage session cookies directly via `Set-Cookie` headers, eliminating the need for the SolidStart frontend to proxy cookies. This simplifies authentication for both Web and Mobile (Capacitor).
- **REST API Standardization**: Removed the custom `{ success, data, error }` JSON envelope. The API now follows standard REST conventions, returning raw data with 2xx status codes and structured `APIError` objects with 4xx/5xx codes.
- **Direct Client-to-API Communication**: Updated the SolidStart frontend and Capacitor configurations to allow direct communication between the browser/mobile app and the Go API (using `credentials: "include"` and updated CORS settings).
- **Pure SPA/SSG Transition**: Migrated the SolidStart frontend to a pure Single Page Application (SSR: false) using Static Site Generation (SSG). This ensures compatibility with Capacitor and improves load performance by eliminating server-side hydration.
- **Responsive App Layout**: Implemented a unified `AppLayout` with a top navbar for desktop/tablet and a fixed bottom tab bar for mobile devices, optimized for touch-friendly navigation.
- **Client-Side Auth Flow**: Refactored authentication routes (Login, Register, Forgot Password) to use client-side state management and direct API calls, removing the need for SolidStart server actions.
- **Simplified SSO**: Refactored the Google OAuth flow to set the session cookie directly during the redirect callback, making the one-time exchange token optional and streamlining the login process.
- **Type Synchronization**: Standardized TypeScript interfaces to match Go DTOs for improved type safety across the boundary.
- **Lean Session API**: `GET /api/v1/auth/me` now returns a `SessionProfile` DTO (name, email, theme only) instead of the full internal `User` model.
- **Dashboard REST API**: Added `GET /api/v1/dashboard` and `POST /api/v1/dashboard/layout` JSON endpoints with a lean `DashboardPayload` that omits internal IDs, Plaid identifiers, and other server-only fields.
- **SolidJS Dashboard**: Implemented a customizable widget grid in the new frontend with ApexCharts visualizations, SortableJS drag-and-drop customization, and mobile long-press edit mode. The **Spending Trend** chart has been updated from an area chart to a **bar chart** for better monthly comparison. Widget span (1 vs 2 columns) is driven by `WIDGET_META` and applied via CSS modules so Net Worth, This Month, Spending Trend, and Recent Transactions span two grid cells on tablet/desktop; mobile tiles size to their content instead of fixed row heights. Donut chart center totals use theme `--text` / `--text-muted` colors (not ApexCharts' default gray) for readable contrast across all themes. Tag donut widgets show a compact abbreviated total in the chart center (with dynamic font sizing), the full precise total in the widget header, and full `formatCurrency` values in slice tooltips. Net Worth, This Month, and Spending Trend metric tiles each use semantic accent colors (green cash, emerald savings, indigo investments, red credit, orange loans, blue accounts; amber spend, emerald income, blue runway) via a shared `--tile-tone` CSS variable pattern.
- **Docker Dev Proxy Fix**: Corrected Vite `server.hmr`/`origin` config for Caddy reverse-proxy access, added a frontend entrypoint that runs `npm install` on container start (fixes stale `node_modules` volumes), and routed `/_build/*` explicitly through Caddy.
- **Dashboard Null-Safe Arrays**: API payload normalization now returns `[]` instead of JSON `null` for empty dashboard slices; frontend coalesces nullable arrays before rendering widgets.
- **Global Theme Application**: Added a root `ThemeProvider` that loads the user's theme preference on app start and listens for OS color-scheme changes; `login()` and `register()` reload the theme after establishing a session.
- **Dashboard UI Polish**: Replaced text-based icons with custom SVGs for the widget drag handle and visibility toggle, and improved button styling with smoother transitions and better centering.
- **Independent Dashboard Layouts**: Implemented support for separate desktop and mobile dashboard layouts. The backend stores both in a unified JSON structure, while the frontend dynamically switches layouts based on viewport size and allows per-device customization without affecting the other view.
- **Icon Organization**: Centralized all SVG icons into a dedicated `frontend/src/components/icons` directory as reusable SolidJS components. Updated the dashboard and main navigation (navbar and tab bar) to use these unified icon components, improving consistency and maintainability.
- **Layout Refactor**: Moved all layout components (`AppLayout`, `AuthLayout`) from `frontend/src/components` to a top-level `frontend/src/layouts` directory. This better separates global structural components from atomic UI components and follows common SolidStart/React architectural patterns. Updated all route imports to reflect the new paths.
- **Transaction Page Refactor**:
    - **Modular Component Refactor**: Decomposed the massive `TransactionsPage` into 8 logical, reusable components (`TransactionHeader`, `TransactionToolbar`, `FilterPanel`, `TransactionList`, `TransactionCard`, `Pagination`, `BulkActionBar`, `TransactionSkeletons`) in `frontend/src/components/transactions/`. This reduced the main route file size by ~60% and improved code readability.
    - **Optimized Reactivity**: Implemented `batch` updates and grouped signal changes for filters to prevent redundant API requests and ensure smoother UI transitions.
- **Stable Metadata**: Introduced a dedicated metadata resource for categories and tags, ensuring filter dropdowns remain stable and populated even while transaction results are loading.
- **Flicker-Free Updates**: 
    - Integrated `useTransition` for filter and pagination updates, keeping the current UI visible while loading new data to eliminate skeleton-induced flickering.
    - Refactored `Suspense` boundaries to isolate metadata loading from transaction data, ensuring filter controls remain stable and interactive during refreshes.
    - Leveraged `data.latest` for result counts and pagination, preventing them from disappearing or triggering `Suspense` when existing data is being updated.
    - Improved filter transition feedback by wrapping "Clear all filters" in `startTransition` and implementing a global loading progress bar.
    - Added a `visualLoading` state with a minimum 250ms duration to prevent flickering and provide clear visual feedback during transaction updates.
- **Capacitor/Mobile Fixes**: 
    - Adjusted the bulk action bar's sticky position on mobile to prevent overlapping with the Capacitor bottom tab bar.
    - Replaced the inline scrollable filter panel on mobile with a full bottom sheet (backdrop, drag handle, sticky header/footer, body scroll lock) so all filters are reachable without scrolling the page.
    - Fixed a critical positioning bug where the mobile filter sheet would appear at the bottom of the long transaction list (making it inaccessible) by implementing conditional rendering: the filter panel is rendered inline for desktop to maintain the layout, but wrapped in a SolidJS `Portal` for mobile. This decouples the mobile bottom sheet from the animated page container while preserving standard desktop behavior.
    - Added an active-filter count badge on the filter toggle and larger touch targets for filter inputs on small screens.
    - Improved date handling for mobile browser compatibility and enhanced checkbox touch targets and layout spacing for better mobile ergonomics.
- **Filter Robustness & Performance**:
    - **UTC Alignment**: Standardized on UTC for transaction date processing. Frontend now parses filter selections as UTC timestamps and formats transaction dates using the UTC timezone, resolving a bug where local timezone offsets would cause incorrect filtering or date shifting.
    - **Intelligent Input Debouncing**: Implemented a unified debouncing effect for search, date range, and amount filters. This prevents the "UI reset" issue where rapid keystrokes triggered overlapping API requests and loading skeletons, ensuring focus is maintained and the UI remains stable while the user is inputting values.
    - **Improved Clear Flow**: Refactored `clearFilters` to reset all debounced state, ensuring a clean and consistent UI state after clearing active filters.
- **Settings REST API**: Added JSON endpoints for profile (`GET/PATCH /api/v1/settings`), password, theme, SSO unlink, and full Plaid connection management (`/api/v1/plaid/*`). External DTOs omit internal IDs while preserving row/account identifiers needed for Plaid Link flows.
- **Google Account Linking (API)**: Extended encrypted OAuth state to support `action=link`, allowing authenticated users to link Google SSO from the SolidJS settings page without creating a new session.
- **SolidJS Settings Page**: Implemented a tabbed settings experience (`Account`, `Connections`, `Appearance`) with profile editing, password management, Google link/unlink, theme selection, and Plaid connection management. The Connections tab mirrors legacy manage-page behavior: accounts must be disconnected through Plaid before they can be removed, with clear in-app guidance. The Account tab includes a delete-account flow with password or Google re-authentication, final confirmation, and permanent data removal via new JSON endpoints (`GET/POST /api/v1/settings/delete/*`).
- **Expanded Theme Palette**: Added selectable themes beyond light/dark/system — Tokyo Night, Coffee, Forest, Rose, Midnight, and Parchment — with shared frontend/backend validation via `theme-options.ts` and `models.IsValidThemePreference`.
- **Tags REST API**: Added JSON endpoints for tag and category management (`GET/POST/PUT/DELETE /api/v1/tags`, `GET /api/v1/tags/:id/filters`, `POST /api/v1/tags/:id/move`, `POST/PUT/DELETE /api/v1/categories`). Mutations return the refreshed `TagsPayload` to avoid extra client round-trips.
- **SolidJS Tags Page**: Implemented the manage tags experience with category cards, tag pills, inline category rename, create/edit tag modals (color picker, category selector, and auto-tagging filters with Save & Apply), new category modal, and delete-category flow (move to Misc, move to another category, or delete all tags). Merge categories was intentionally omitted. Tag drag-and-drop between categories was removed; use the edit tag modal to change a tag's category. Tag pills on desktop now use a fixed min-height and normalized line-height so every pill aligns consistently in the category grid. Fixed a bug where default tag colors would not render on the Dashboard or Transactions pages by ensuring palette keys (e.g., "emerald", "rose") are correctly translated to hex codes using the `tagColorHex` utility in all relevant components.
- **Security Audit & Hardening (Completed June 2026)**:
## Security Audit & Hardening (Completed June 2026)
- **OAuth State Hardening**: Implemented HMAC-signed `state` parameters in the SSO flow to prevent authentication-based CSRF and account linking attacks.
- **Data Encryption at Rest**: Added AES-GCM encryption for sensitive Plaid access tokens. Tokens are now encrypted before being stored in SQLite and decrypted only when needed for bank communication.
- **CSRF Protection**: Enabled global CSRF middleware using Echo's double-submit cookie pattern. The `_csrf` cookie (non-HttpOnly) is set on every response and the SPA reads it to send back as the `X-CSRF-Token` header on state-changing requests. A `GET /api/v1/csrf` endpoint primes the token on SPA startup. The CORS config includes `X-CSRF-Token` in allowed headers and trusted origins are configured for the frontend. Frontend `clientApiRequest()` automatically handles CSRF token attachment for POST/PUT/PATCH/DELETE methods.
- **Server Hardening**: 
    - Enabled standard security headers (X-Frame-Options, X-Content-Type-Options, etc.) via Echo's `Secure` middleware.
    - Enforced `Secure` and `HttpOnly` flags on all session cookies in production environments.
- **Rate Limiting**: Implemented memory-based rate limiting on authentication endpoints (Login, Register, Forgot Password) to mitigate brute-force and credential stuffing attacks, with a bypass configured for development environments.
