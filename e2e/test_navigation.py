import pytest
from playwright.sync_api import Page, expect

def test_navbar_present(auth_page: Page, base_url: str):
    """Test that the navbar contains all navigation links."""
    auth_page.goto(f"{base_url}/dashboard")

    nav = auth_page.get_by_role("navigation")
    expect(nav.get_by_role("link", name="Dashboard")).to_be_visible()
    expect(nav.get_by_role("link", name="Transactions")).to_be_visible()
    expect(nav.get_by_role("link", name="Tags")).to_be_visible()
    expect(nav.get_by_role("link", name="Settings")).to_be_visible()

def test_navbar_brand(auth_page: Page, base_url: str):
    """Test that the FinancialTracker brand is displayed in the header."""
    auth_page.goto(f"{base_url}/dashboard")

    expect(auth_page.get_by_text("FinancialTracker")).to_be_visible()

def test_navigate_dashboard_to_transactions(auth_page: Page, base_url: str):
    """Test navigation from dashboard to transactions via navbar."""
    auth_page.goto(f"{base_url}/dashboard")

    auth_page.get_by_role("navigation").get_by_role("link", name="Transactions").click()
    expect(auth_page).to_have_url(f"{base_url}/transactions")
    expect(auth_page.get_by_role("heading", name="Transactions", exact=True)).to_be_visible()

def test_navigate_dashboard_to_tags(auth_page: Page, base_url: str):
    """Test navigation from dashboard to tags via navbar."""
    auth_page.goto(f"{base_url}/dashboard")

    auth_page.get_by_role("navigation").get_by_role("link", name="Tags").click()
    expect(auth_page).to_have_url(f"{base_url}/tags")
    expect(auth_page.get_by_role("heading", name="Tags & Categories", exact=True)).to_be_visible()

def test_navigate_dashboard_to_settings(auth_page: Page, base_url: str):
    """Test navigation from dashboard to settings via navbar."""
    auth_page.goto(f"{base_url}/dashboard")

    auth_page.get_by_role("navigation").get_by_role("link", name="Settings").click()
    expect(auth_page).to_have_url(f"{base_url}/settings")
    expect(auth_page.get_by_role("heading", name="Settings")).to_be_visible()

def test_navigate_back_to_dashboard(auth_page: Page, base_url: str):
    """Test navigating away and back to dashboard."""
    auth_page.goto(f"{base_url}/transactions")

    auth_page.get_by_role("navigation").get_by_role("link", name="Dashboard").click()
    expect(auth_page).to_have_url(f"{base_url}/dashboard")
    expect(auth_page.get_by_role("heading", name="Welcome back", level=1)).to_be_visible()

def test_view_all_transactions_link(auth_page: Page, base_url: str):
    """Test that 'View all' link on dashboard navigates to transactions."""
    auth_page.goto(f"{base_url}/dashboard")

    auth_page.get_by_role("link", name="View all").click()
    expect(auth_page).to_have_url(f"{base_url}/transactions")

def test_manage_connections_link(auth_page: Page, base_url: str):
    """Test that 'Manage connections' on dashboard navigates to settings connections tab."""
    auth_page.goto(f"{base_url}/dashboard")

    auth_page.get_by_role("link", name="Manage connections").click()
    expect(auth_page).to_have_url(f"{base_url}/settings?tab=connections")
