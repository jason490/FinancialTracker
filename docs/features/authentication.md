# Authentication

## Overview

Authentication controls who can use FinancialTracker. Users sign in with email and password or Google SSO. The API stores sessions in SQLite and sets an HTTP-only session cookie. New users complete onboarding before they reach the dashboard.

## How it works

### Email and password login

1. Open `/login`.
2. Enter email and password.
3. Optional: select Remember me for a one-year session. Default session length is one day.
4. Submit the form. The API verifies credentials with bcrypt and creates a session.
5. The browser stores the `Session` cookie. The SPA loads the theme and navigates with `postAuthPath`.

Successful login sends incomplete users to `/onboarding`. Completed users go to `/dashboard`.

### Registration

1. Open `/register`.
2. Enter first name, last name, email, password, and password confirmation.
3. When invite-only mode is active, enter a registration code.
4. Submit. The API validates input, creates the user, seeds a session, and sets the cookie.

New users start with `onboarding_completed = 0` and theme preference `system`.

### Google SSO

1. Start Google sign-in from login or register (`GET /api/v1/auth/google`).
2. Complete Google consent. Google redirects to `/api/v1/auth/google/callback`.
3. The API verifies HMAC-signed OAuth state and loads Google user info.
4. If a `user_sso` row matches the Google account, the API creates a session.
5. If no SSO link exists and no local account uses that email, the API creates a user and links Google.
6. The SPA finishes at `/auth/sso/complete` and redirects with `postAuthPath`.

Link mode (`action=link`) attaches Google to an authenticated account from Settings. Reauth-delete mode verifies Google before account deletion.

### SSO conflict behavior

Google SSO does not auto-link to a password-only account that shares the same email. The service returns `ErrSSOAccountConflict`. The callback redirects with error code `sso_account_exists`. The complete page tells the user to sign in with a password and link Google in Settings.

### Password reset

1. Open `/forgot-password` and submit the account email.
2. The API issues a 6-digit code with a 15-minute TTL when the account has a password.
3. Enter the code. The API verifies it (`POST /api/v1/auth/verify-reset-code`).
4. Set a new password (`POST /api/v1/auth/reset-password`). The API invalidates all sessions.

Rules:

- Maximum five verification attempts per code.
- First resend is immediate. Later resends wait 60 seconds.
- Unknown emails and SSO-only accounts receive the same generic success response. The API does not send a code in those cases.
- In development without `MAIL_SMTP_HOST`, the API logs the code to the console.
- In production, SMTP delivers the code.

### Sessions

`GET /api/v1/auth/me` returns a lean `SessionProfile` (name, email, theme, onboarding status). Protected routes require a valid session cookie. Logout deletes the session (`POST /api/v1/auth/logout`). Production cookies use Secure and HttpOnly flags. Auth routes use rate limiting outside development.

### Account deletion reauth

Account deletion requires recent re-authentication (five-minute window).

1. Open Settings → Account → Delete account.
2. Reauth with password (`POST /api/v1/settings/delete/verify`) or Google reauth-delete SSO.
3. Confirm deletion (`POST /api/v1/settings/delete/confirm`).
4. The API deletes the user and cascaded data when reauth is still valid.

`GET /api/v1/settings/delete/reauth-status` reports whether the session has a valid reauth timestamp.

### Invite-only gate

When `SUBSCRIPTIONS_ENABLED=false`, registration requires a single-use admin invite code.

- `GET /api/v1/auth/registration-config` reports whether codes are required.
- Email register accepts `registration_code`.
- Google register carries the code in OAuth state.
- Codes are 8-character alphanumeric, bcrypt-hashed, 48-hour TTL, single-use.
- Admins in `REGISTRATION_ADMIN_EMAILS` create codes via Settings → Account or `POST /api/v1/admin/registration-codes`.
- In development, `test@test.com` is always a registration admin.
- Bootstrap without a signed-in admin: `make registration-code`.

Invalid or missing codes map to `invalid_registration_code` or `registration_code_required` on the SSO complete page.

## API endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| POST | `/api/v1/auth/login` | No | Email/password login |
| GET | `/api/v1/auth/registration-config` | No | Invite-gate status |
| POST | `/api/v1/auth/register` | No | Create account |
| POST | `/api/v1/auth/forgot-password` | No | Request reset code |
| POST | `/api/v1/auth/verify-reset-code` | No | Validate reset code |
| POST | `/api/v1/auth/reset-password` | No | Set new password |
| GET | `/api/v1/auth/google` | No | Start Google OAuth |
| GET | `/api/v1/auth/google/callback` | No | Google OAuth callback |
| GET | `/api/v1/auth/me` | Session | Current profile |
| POST | `/api/v1/auth/logout` | Yes | End session |
| POST | `/api/v1/auth/onboarding/complete` | Yes | Mark onboarding done |
| POST | `/api/v1/admin/registration-codes` | Yes (admin) | Issue invite code |
| GET | `/api/v1/settings/delete/reauth-status` | Yes | Deletion reauth status |
| POST | `/api/v1/settings/delete/verify` | Yes | Password reauth for delete |
| POST | `/api/v1/settings/delete/confirm` | Yes | Confirm account deletion |
| GET | `/api/v1/csrf` | No | Prime CSRF cookie/token |

## Configuration

| Variable | Purpose |
|----------|---------|
| `ENV` | Development relaxes auth rate limits and logs reset codes without SMTP |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret |
| `ENCRYPTION_KEY` | 32-byte key for OAuth state signing/encryption |
| `SUBSCRIPTIONS_ENABLED` | `false` enables invite-only registration |
| `REGISTRATION_ADMIN_EMAILS` | Semicolon-separated invite admins |
| `MAIL_SMTP_HOST` | SMTP host for reset email |
| `MAIL_SMTP_PORT` | SMTP port |
| `MAIL_SMTP_USER` | SMTP username |
| `MAIL_SMTP_PASSWORD` | SMTP password or API key |
| `MAIL_FROM` | Sender address |
| `MAIL_FROM_NAME` | Optional sender display name |
| `API_PUBLIC_URL` | Public API URL for OAuth callbacks |
| `FRONTEND_URL` | SPA URL for return paths |

## Related files

- `app/internal/services/auth/auth.go`
- `app/internal/services/auth/sso.go`
- `app/internal/services/auth/password_reset.go`
- `app/internal/services/auth/registration_code.go`
- `app/internal/server/handler/auth.go`
- `app/internal/server/handler/sso.go`
- `app/internal/server/handler/admin.go`
- `app/internal/services/settings/settings.go`
- `frontend/src/routes/login.tsx`
- `frontend/src/routes/register.tsx`
- `frontend/src/routes/auth/sso/complete.tsx`
- `frontend/src/lib/auth.ts`
- `frontend/src/lib/auth-context.tsx`

## Related docs

- [Onboarding](./onboarding.md)
- [Settings](./settings.md)
- [Subscriptions](./subscriptions.md)
- [Environment variables](../getting-started/environment.md)
