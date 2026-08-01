# Settings

## Overview

Settings is the account and preferences hub at `/settings`. Tabs cover profile security, bank links, themes, billing, tags, and data export. The Plan tab appears only when subscriptions are enabled. Registration admins see an invite codes panel when invite-only mode is active.

## How it works

### Tabs

| Tab | Query | Contents |
|-----|-------|----------|
| Account | `?tab=account` | Profile, password, Google link/unlink, logout, delete account, invite codes (admins) |
| Connections | `?tab=connections` | Link, sync, update, disconnect, visibility, remove |
| Appearance | `?tab=appearance` | Theme selection |
| Plan | `?tab=plan` | Tier, usage, checkout, portal (hidden when subscriptions disabled) |
| Tags | `?tab=tags` | Categories, tags, filters, drag-and-drop |
| Data | `?tab=data` | ZIP/CSV transaction export |

The SPA mounts only the active panel. Initial load uses `AccountPanelSkeleton` with `LoadingCrossfade`. Connection actions refetch under `startTransition` to avoid panel flash.

### Themes

Appearance options: Light, Dark, System, Tokyo Night, Coffee, Forest, Rose, Midnight, Parchment. Preference stores on the user and applies through `ThemeProvider` / `data-theme`. System follows the OS color scheme. Validation is shared in `theme-options.ts` and `models.IsValidThemePreference`.

### Profile

Account panel actions:

1. Update name/email with `PATCH /api/v1/settings/profile`.
2. Change password with `POST /api/v1/settings/password`.
3. Link Google from Settings (`action=link` OAuth) or unlink with `POST /api/v1/settings/unlink/:provider`.
4. Sign out with `POST /api/v1/auth/logout`.
5. Delete account after reauth (see [Authentication](./authentication.md)).

`GET /api/v1/settings` returns the settings profile DTO.

### Invites panel

When `SUBSCRIPTIONS_ENABLED=false` and the signed-in user is a registration admin, Account shows **Invite codes**. Admins create single-use codes through `POST /api/v1/admin/registration-codes`. The panel shows the plaintext code once and supports copy with a clipboard fallback for non-secure contexts.

Checkout query params `checkout=success` and `checkout=cancelled` surface Plan tab banners after Stripe return. SSO link success uses `success=linked`.

## API endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/api/v1/settings` | Yes | Settings profile |
| PATCH | `/api/v1/settings/profile` | Yes | Update profile |
| POST | `/api/v1/settings/password` | Yes | Change password |
| POST | `/api/v1/settings/theme` | Yes | Save theme |
| POST | `/api/v1/settings/unlink/:provider` | Yes | Unlink SSO |
| GET | `/api/v1/settings/delete/reauth-status` | Yes | Delete reauth status |
| POST | `/api/v1/settings/delete/verify` | Yes | Password reauth |
| POST | `/api/v1/settings/delete/confirm` | Yes | Confirm deletion |
| GET | `/api/v1/connections` | Yes | Connections panel data |
| GET | `/api/v1/tags` | Yes | Tags panel data |
| GET | `/api/v1/subscription` | Yes | Plan panel data |
| GET | `/api/v1/transactions/export` | Yes | Data export |
| POST | `/api/v1/admin/registration-codes` | Yes (admin) | Create invite code |

## Configuration

| Variable | Purpose |
|----------|---------|
| `SUBSCRIPTIONS_ENABLED` | Hides Plan tab when `false`; enables invite panel path |
| `REGISTRATION_ADMIN_EMAILS` | Who may issue invite codes |
| `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` | Google link/unlink |
| Stripe / Plaid vars | Plan and Connections tabs (see related docs) |

## Related files

- `frontend/src/routes/settings.tsx`
- `frontend/src/components/settings/AccountPanel.tsx`
- `frontend/src/components/settings/ConnectionsPanel.tsx`
- `frontend/src/components/settings/AppearancePanel.tsx`
- `frontend/src/components/settings/PlanPanel.tsx`
- `frontend/src/components/settings/TagsPanel.tsx`
- `frontend/src/components/settings/DataPanel.tsx`
- `frontend/src/components/settings/RegistrationInvitesPanel.tsx`
- `frontend/src/lib/theme-options.ts`
- `app/internal/services/settings/settings.go`
- `app/internal/server/handler/settings.go`

## Related docs

- [Authentication](./authentication.md)
- [Bank connections](./bank-connections.md)
- [Tags and categories](./tags-and-categories.md)
- [Subscriptions](./subscriptions.md)
- [Transactions](./transactions.md)
- [Environment variables](../getting-started/environment.md)
