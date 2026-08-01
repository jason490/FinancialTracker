# Onboarding

## Overview

Onboarding is a short wizard for new accounts. It introduces the product, optional plan selection, and optional bank connection. The flow ends when the client calls complete onboarding and the user reaches the dashboard.

## How it works

### Three-step wizard

Route: `/onboarding`. Steps:

1. **Welcome** — greets the user by first name and explains next steps.
2. **Plan** — shared `PlanPicker` for Free / Plus / Premium (only when subscriptions are enabled).
3. **Connect** — start a bank connection or skip.

When `SUBSCRIPTIONS_ENABLED=false`, the plan step is removed. Progress shows Welcome → Connect only. If the client lands on plan while subscriptions are disabled, the page moves to connect.

### `onboarding_completed`

New password and Google signups create users with `onboarding_completed = 0`. `SessionProfile` includes the flag. `POST /api/v1/auth/onboarding/complete` sets it to completed.

`AppLayout` redirects authenticated users with incomplete onboarding from app routes to `/onboarding`. Users who already completed onboarding and open `/onboarding` go to `/dashboard`.

### Skip bank

On the connect step:

1. Optional: connect a bank with `startNewConnection()`.
2. Or choose **Skip for now**.
3. Finish calls `completeOnboarding()`, refetches the session, preloads the dashboard route, and navigates to `/dashboard`.

Bank connection is not required to finish.

### Redirects

| Situation | Destination |
|-----------|-------------|
| Login/SSO with incomplete onboarding | `/onboarding` (`postAuthPath`) |
| Login/SSO with completed onboarding | `/dashboard` |
| Authenticated app route while incomplete | `/onboarding` |
| `/onboarding` while already complete | `/dashboard` |
| Unauthenticated `/onboarding` | `/login` |

Post-auth overlay stays up until onboarding subscription data is ready, then dismisses. Finish uses the dashboard skeleton handoff instead of a “all set” overlay.

## API endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| POST | `/api/v1/auth/onboarding/complete` | Yes | Mark wizard finished |
| GET | `/api/v1/auth/me` | Session | Includes `onboarding_completed` |
| GET | `/api/v1/subscription` | Yes | Plan step / visibility |
| GET | `/api/v1/connections` | Yes | Connect step state |
| POST | `/api/v1/connections/create-session` | Yes | Start bank link |
| POST | `/api/v1/connections/complete` | Yes | Finish bank link |

## Configuration

| Variable | Purpose |
|----------|---------|
| `SUBSCRIPTIONS_ENABLED` | Shows or hides the plan step |
| `FINANCIAL_PROVIDER` | Provider used on the connect step |

## Related files

- `frontend/src/routes/onboarding.tsx`
- `frontend/src/layouts/OnboardingLayout.tsx`
- `frontend/src/components/onboarding/OnboardingWelcomeStep.tsx`
- `frontend/src/components/onboarding/OnboardingPlanStep.tsx`
- `frontend/src/components/onboarding/OnboardingConnectStep.tsx`
- `frontend/src/components/onboarding/OnboardingProgress.tsx`
- `frontend/src/lib/auth.ts`
- `frontend/src/layouts/AppLayout.tsx`
- `app/internal/services/auth/auth.go`
- `app/internal/server/handler/auth.go`

## Related docs

- [Authentication](./authentication.md)
- [Subscriptions](./subscriptions.md)
- [Bank connections](./bank-connections.md)
- [Dashboard](./dashboard.md)
- [Settings](./settings.md)
