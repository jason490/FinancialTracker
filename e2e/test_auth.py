import pytest
from playwright.sync_api import Page, expect

def test_login_success(page: Page, base_url: str):
    """Test that a user can successfully log in."""
    page.goto(f"{base_url}/login")

    page.fill('input[name="email"]', "test@test.com")
    page.fill('input[name="password"]', "test")
    page.click('button[type="submit"]')

    # Verify navigation to dashboard
    expect(page).to_have_url(f"{base_url}/dashboard")
    expect(page.get_by_role("heading", name="Welcome back", level=1)).to_be_visible()

def test_login_failure(page: Page, base_url: str):
    """Test that invalid credentials show an error message."""
    page.goto(f"{base_url}/login")

    page.fill('input[name="email"]', "wrong@test.com")
    page.fill('input[name="password"]', "wrongpassword")
    page.click('button[type="submit"]')

    # The login form renders an alert element with the error message
    expect(page.get_by_role("alert")).to_be_visible()
    expect(page.get_by_text("invalid credentials")).to_be_visible()

def test_login_empty_fields_prevented(page: Page, base_url: str):
    """Test that submitting empty login fields is blocked by HTML validation."""
    page.goto(f"{base_url}/login")

    email_input = page.locator('input[name="email"]')
    password_input = page.locator('input[name="password"]')

    # Both inputs should be required
    expect(email_input).to_have_attribute("required", "")
    expect(password_input).to_have_attribute("required", "")

def test_protected_route_redirect(page: Page, base_url: str):
    """Test that accessing a protected route redirects to login."""
    page.goto(f"{base_url}/dashboard")

    expect(page).to_have_url(f"{base_url}/login")

def test_protected_route_transactions_redirect(page: Page, base_url: str):
    """Test that accessing /transactions unauthenticated redirects to login."""
    page.goto(f"{base_url}/transactions")

    expect(page).to_have_url(f"{base_url}/login")

def test_protected_route_settings_redirect(page: Page, base_url: str):
    """Test that accessing /settings unauthenticated redirects to login."""
    page.goto(f"{base_url}/settings")

    expect(page).to_have_url(f"{base_url}/login")

def test_protected_route_tags_redirect(page: Page, base_url: str):
    """Test that accessing /tags unauthenticated redirects to login."""
    page.goto(f"{base_url}/tags")

    expect(page).to_have_url(f"{base_url}/login")

def test_login_page_structure(page: Page, base_url: str):
    """Test that the login page renders all expected elements."""
    page.goto(f"{base_url}/login")

    # Heading
    expect(page.get_by_role("heading", name="Sign in")).to_be_visible()

    # Form fields
    expect(page.get_by_label("Email address")).to_be_visible()
    expect(page.get_by_label("Password")).to_be_visible()

    # Submit button
    expect(page.get_by_role("button", name="Sign in")).to_be_visible()

    # Remember me checkbox
    expect(page.get_by_role("checkbox", name="Remember me")).to_be_visible()

    # Navigation links
    expect(page.get_by_role("link", name="Forgot password?")).to_be_visible()
    expect(page.get_by_role("link", name="Create one")).to_be_visible()

    # SSO button
    expect(page.get_by_role("link", name="Continue with Google")).to_be_visible()

def test_login_navigates_to_register(page: Page, base_url: str):
    """Test that clicking 'Create one' navigates to the register page."""
    page.goto(f"{base_url}/login")

    page.get_by_role("link", name="Create one").click()

    expect(page).to_have_url(f"{base_url}/register")

def test_login_navigates_to_forgot_password(page: Page, base_url: str):
    """Test that clicking 'Forgot password?' navigates to the reset page."""
    page.goto(f"{base_url}/login")

    page.get_by_role("link", name="Forgot password?").click()

    expect(page).to_have_url(f"{base_url}/forgot-password")
