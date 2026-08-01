# Security

## Overview

FinancialTracker protects sessions, OAuth state, CSRF, secrets at rest, and privileged actions. This page describes the main controls.

## Sessions

- Login and SSO create a server-side session row and set an HTTP-only `Session` cookie.
- Production cookies also use the Secure flag.
- Logout deletes the session.
- Password reset invalidates all sessions for the user.

## OAuth state

Google OAuth state is HMAC-signed (and encrypted where required) with `ENCRYPTION_KEY`. This blocks CSRF-style login and account-linking attacks. Link and reauth-delete flows use dedicated state actions.

## CSRF

- Echo double-submit CSRF middleware sets a non-HttpOnly `_csrf` cookie.
- The SPA sends `X-CSRF-Token` on POST, PUT, PATCH, and DELETE.
- `GET /api/v1/csrf` primes the token on SPA start.
- CORS allows the `X-CSRF-Token` header for trusted frontend origins.
- Plaid and Stripe webhooks are CSRF-exempt and use provider signature checks instead.

## Encryption at rest

Plaid access tokens are encrypted with AES-GCM before SQLite storage. The API decrypts them only when it talks to Plaid. `ENCRYPTION_KEY` must be exactly 32 bytes.

## Re-authentication

Account deletion requires recent reauth (password verify or Google reauth-delete). The reauth window is short. Confirm deletion only succeeds while that window is valid.

## Rate limiting

Login, register, and forgot-password use memory rate limiting outside development. This reduces brute-force and credential stuffing.

## Container hardening

Production containers:

- Run as non-root users after entrypoint permission fixes
- Set standard security headers through Echo Secure middleware
- Bind Caddy to unprivileged internal ports where configured

## SSO account safety

Google SSO does not auto-link to an existing password-only account that shares the same email. The user must sign in with a password and link Google from Settings.

## Related docs

- [Authentication](../features/authentication.md)
- [Environment variables](../getting-started/environment.md)
- [Production deployment](../getting-started/production.md)
