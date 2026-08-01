# Mobile (Capacitor / Android)

## Overview

FinancialTracker can ship as an Android app through Capacitor. The frontend must build as a static site (SSG) so Capacitor can serve `index.html` and assets from the web directory.

## Build and sync

1. Build the static frontend:

```bash
cd frontend
npm run build
```

2. Sync the build into the native project:

```bash
npx cap sync
```

3. Open the Android project in Android Studio on a machine that has the Android SDK.
4. Run the app on a device or emulator.

## Notes

- Point `VITE_API_PUBLIC_URL` at a reachable API host. The Android WebView cannot use `localhost` for a host-machine API without extra setup.
- Session cookies and CSRF work with `credentials: "include"` against the configured API origin.
- Transaction ZIP export uses a blob download that works in the Capacitor WebView.

## Related docs

- [Development setup](./development.md)
- [Environment variables](./environment.md)
- [Architecture overview](../architecture/overview.md)
