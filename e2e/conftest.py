import pytest
import time
import urllib.request
import urllib.error
import socket
from playwright.sync_api import Browser, BrowserContext, Page

@pytest.fixture(scope="session")
def browser_context_args(browser_context_args):
    return {
        **browser_context_args,
        "ignore_https_errors": True,
        "viewport": {"width": 1280, "height": 720},
    }

@pytest.fixture(scope="session", autouse=True)
def wait_for_server(browser: Browser, base_url: str):
    """Wait for the frontend and backend to be fully ready before running tests.

    Uses a two-phase approach:
      Phase 1 – Raw HTTP polling (fast, avoids Playwright overhead).
      Phase 2 – Playwright check to confirm SolidJS hydration, using a fresh
                browser context on every attempt to avoid stale-page cascades.
    """
    # Phase 1: Verify the server is reachable at the HTTP level.
    print(f"\n⏳ Phase 1: Waiting for server at {base_url} to accept connections...")
    max_http_wait_s = 180
    poll_interval_s = 2
    elapsed = 0
    server_reachable = False

    while elapsed < max_http_wait_s:
        try:
            req = urllib.request.Request(f"{base_url}/login", method="GET")
            resp = urllib.request.urlopen(req, timeout=5)
            if resp.status < 500:
                print(f"  ✅ Server responded (HTTP {resp.status}) after ~{elapsed}s")
                server_reachable = True
                break
        except urllib.error.HTTPError as e:
            # A 4xx still proves the server is listening
            if e.code < 500:
                print(f"  ✅ Server responded (HTTP {e.code}) after ~{elapsed}s")
                server_reachable = True
                break
            if elapsed % 20 == 0:
                print(f"  ... server error (HTTP {e.code}), retrying ({elapsed}s elapsed)")
        except (urllib.error.URLError, socket.timeout, OSError) as e:
            if elapsed % 20 == 0:
                print(f"  ... not reachable yet ({elapsed}s): {e}")
        time.sleep(poll_interval_s)
        elapsed += poll_interval_s

    if not server_reachable:
        raise Exception(
            f"Server at {base_url} did not respond to HTTP within {max_http_wait_s}s. "
            f"Ensure docker-compose services are running and the proxy is accessible."
        )

    # Phase 2: Confirm SolidJS has hydrated by checking for the login form.
    # A *fresh* browser context is created on every attempt so a failed navigation
    # never poisons subsequent retries.
    print("  ⏳ Phase 2: Waiting for SolidJS to hydrate the login page...")
    max_hydration_retries = 20
    hydration_poll_s = 5

    for attempt in range(max_hydration_retries):
        context = browser.new_context(ignore_https_errors=True)
        page = context.new_page()
        try:
            response = page.goto(
                f"{base_url}/login",
                timeout=60000,
                wait_until="domcontentloaded",
            )
            if response and response.ok:
                try:
                    page.wait_for_selector('input[name="email"]', timeout=30000)
                    print(
                        f"  ✅ Frontend hydrated after "
                        f"~{attempt * hydration_poll_s}s"
                    )
                    page.close()
                    context.close()
                    return
                except Exception:
                    if attempt % 3 == 0:
                        print(
                            f"  ... selector not found yet "
                            f"({attempt * hydration_poll_s}s)"
                        )
            else:
                status = response.status if response else "no response"
                if attempt % 3 == 0:
                    print(f"  ... page returned {status}")
        except Exception as e:
            if attempt % 3 == 0:
                print(
                    f"  ... browser check failed "
                    f"({attempt * hydration_poll_s}s): {e}"
                )
        finally:
            try:
                page.close()
            except Exception:
                pass
            try:
                context.close()
            except Exception:
                pass
        time.sleep(hydration_poll_s)

    raise Exception(
        f"Server at {base_url} did not become fully ready in time. "
        f"HTTP is reachable but the login page did not render the expected "
        f"selector within {max_hydration_retries * hydration_poll_s}s."
    )

@pytest.fixture(scope="session")
def authenticated_state(browser: Browser, base_url: str):
    """Logs in once and saves the storage state for the session."""
    context = browser.new_context(viewport={"width": 1280, "height": 720})
    page = context.new_page()

    # Go to login page
    page.goto(f"{base_url}/login")

    # Fill login form
    page.fill('input[name="email"]', "test@test.com")
    page.fill('input[name="password"]', "test")
    page.click('button[type="submit"]')

    # Wait for navigation to dashboard
    page.wait_for_url("**/dashboard")

    # Save storage state (cookies, local storage)
    state = context.storage_state()
    context.close()
    return state

@pytest.fixture
def auth_page(browser: Browser, authenticated_state, base_url: str):
    """Returns a page that is already logged in."""
    context = browser.new_context(storage_state=authenticated_state)
    page = context.new_page()
    yield page
    context.close()
