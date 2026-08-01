# Testing

## Overview

FinancialTracker has Go unit tests and a Python Playwright end-to-end suite.

## Go unit tests

1. Change to the `app` directory.
2. Run the test suite.

```bash
cd app
go test ./...
```

## End-to-end tests (Playwright)

1. Start the application so the SPA and API are reachable.
2. Change to the `e2e` directory.
3. Install dependencies and run pytest.

```bash
cd e2e
pip install -r requirements.txt
pytest
```

The suite uses a session-scoped login fixture to reuse auth state. Coverage includes authentication, registration, forgot password, dashboard, transactions, tags (Settings → Tags), settings tabs, landing page, and navigation.

Readiness checks poll HTTP first, then open a fresh browser context per attempt. Default host targeting uses `localhost` with host networking where configured.

## Related docs

- [Development setup](../getting-started/development.md)
- [Makefile commands](./makefile.md)
