import pytest
from playwright.sync_api import Page, expect

def test_landing_page_loads(page: Page, base_url: str):
    """Test that the public landing page loads for unauthenticated users."""
    page.goto(f"{base_url}/")

    expect(page.get_by_role("heading", name="Master your finances", level=1)).to_be_visible()

def test_landing_page_cta_links(page: Page, base_url: str):
    """Test that the CTA links are present on the landing page."""
    page.goto(f"{base_url}/")

    expect(page.get_by_role("link", name="Get started")).to_be_visible()
    expect(page.get_by_role("link", name="Sign in")).to_be_visible()

def test_landing_page_feature_cards(page: Page, base_url: str):
    """Test that the feature cards are displayed."""
    page.goto(f"{base_url}/")

    expect(page.get_by_role("heading", name="Plaid sync")).to_be_visible()
    expect(page.get_by_role("heading", name="Smart tagging")).to_be_visible()
    expect(page.get_by_role("heading", name="Your dashboard")).to_be_visible()

def test_landing_get_started_navigates_to_register(page: Page, base_url: str):
    """Test that 'Get started' navigates to the register page."""
    page.goto(f"{base_url}/")

    page.get_by_role("link", name="Get started").click()

    expect(page).to_have_url(f"{base_url}/register")

def test_landing_sign_in_navigates_to_login(page: Page, base_url: str):
    """Test that 'Sign in' navigates to the login page."""
    page.goto(f"{base_url}/")

    page.get_by_role("link", name="Sign in").click()

    expect(page).to_have_url(f"{base_url}/login")

def test_landing_eyebrow(page: Page, base_url: str):
    """Test that the eyebrow text is present."""
    page.goto(f"{base_url}/")

    expect(page.get_by_text("Automated expense intelligence")).to_be_visible()

def test_landing_lede_text(page: Page, base_url: str):
    """Test that the lede paragraph is present."""
    page.goto(f"{base_url}/")

    expect(page.get_by_text("connects to your banks through Plaid")).to_be_visible()
