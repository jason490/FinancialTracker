import pytest
from playwright.sync_api import Page, expect

def test_dashboard_loading(auth_page: Page, base_url: str):
    """Test that the dashboard loads correctly with all main sections."""
    auth_page.goto(f"{base_url}/dashboard")

    # Verify core layout elements — the h1 is "Welcome back, {name}"
    expect(auth_page.get_by_role("heading", name="Welcome back", level=1)).to_be_visible()

    # Check for sidebar navigation
    expect(auth_page.get_by_role("link", name="Transactions")).to_be_visible()
    expect(auth_page.get_by_role("link", name="Tags")).to_be_visible()
    expect(auth_page.get_by_role("navigation").get_by_role("link", name="Settings")).to_be_visible()

def test_dashboard_widgets(auth_page: Page, base_url: str):
    """Test that dashboard widgets are present."""
    auth_page.goto(f"{base_url}/dashboard")

    # Verify key dashboard widgets are rendered
    expect(auth_page.get_by_role("heading", name="Spending by Tag")).to_be_visible()

def test_dashboard_eyebrow(auth_page: Page, base_url: str):
    """Test that the 'Overview' eyebrow text is displayed."""
    auth_page.goto(f"{base_url}/dashboard")

    expect(auth_page.get_by_text("Overview")).to_be_visible()

def test_dashboard_subtitle(auth_page: Page, base_url: str):
    """Test that the dashboard subtitle is present."""
    auth_page.goto(f"{base_url}/dashboard")

    expect(auth_page.get_by_text("linked accounts, spending trends")).to_be_visible()

def test_dashboard_customize_button(auth_page: Page, base_url: str):
    """Test that the 'Customize' button is present for widget editing."""
    auth_page.goto(f"{base_url}/dashboard")

    expect(auth_page.get_by_role("button", name="Customize")).to_be_visible()

def test_dashboard_this_month_widget(auth_page: Page, base_url: str):
    """Test that the 'This Month' summary widget is present."""
    auth_page.goto(f"{base_url}/dashboard")

    expect(auth_page.get_by_role("heading", name="This Month")).to_be_visible()
    this_month = auth_page.get_by_role("article").filter(
        has=auth_page.get_by_role("heading", name="This Month")
    )
    expect(this_month.get_by_text("Spend", exact=True)).to_be_visible()
    expect(this_month.get_by_text("Income", exact=True)).to_be_visible()

def test_dashboard_quick_actions_widget(auth_page: Page, base_url: str):
    """Test that the Quick Actions widget renders with expected buttons."""
    auth_page.goto(f"{base_url}/dashboard")

    expect(auth_page.get_by_role("heading", name="Quick Actions")).to_be_visible()
    expect(auth_page.get_by_role("button", name="Sync transactions")).to_be_visible()
    expect(auth_page.get_by_role("link", name="Manage connections")).to_be_visible()

def test_dashboard_net_worth_widget(auth_page: Page, base_url: str):
    """Test that the Net Worth widget renders with account breakdowns."""
    auth_page.goto(f"{base_url}/dashboard")

    expect(auth_page.get_by_role("heading", name="Net Worth")).to_be_visible()

def test_dashboard_spending_trend_widget(auth_page: Page, base_url: str):
    """Test that the Spending Trend widget is present."""
    auth_page.goto(f"{base_url}/dashboard")

    expect(auth_page.get_by_role("heading", name="Spending Trend")).to_be_visible()

def test_dashboard_income_by_tag_widget(auth_page: Page, base_url: str):
    """Test that the Income by Tag widget is present."""
    auth_page.goto(f"{base_url}/dashboard")

    expect(auth_page.get_by_role("heading", name="Income by Tag")).to_be_visible()

def test_dashboard_recent_transactions_widget(auth_page: Page, base_url: str):
    """Test that the Recent Transactions widget has a 'View all' link."""
    auth_page.goto(f"{base_url}/dashboard")

    expect(auth_page.get_by_role("heading", name="Recent Transactions")).to_be_visible()
    expect(auth_page.get_by_role("link", name="View all")).to_be_visible()

def test_dashboard_account_category_widgets(auth_page: Page, base_url: str):
    """Test that account category widgets are rendered."""
    auth_page.goto(f"{base_url}/dashboard")

    for name in ["Cash & Checking", "Savings", "Credit Cards", "Loans", "Investments"]:
        expect(auth_page.get_by_role("heading", name=name)).to_be_visible()

def test_dashboard_nav_links_navigate(auth_page: Page, base_url: str):
    """Test that clicking nav links navigates to the correct pages."""
    auth_page.goto(f"{base_url}/dashboard")

    # Click Transactions in the navbar
    auth_page.get_by_role("navigation").get_by_role("link", name="Transactions").click()
    expect(auth_page).to_have_url(f"{base_url}/transactions")

    # Navigate back via Dashboard link
    auth_page.get_by_role("navigation").get_by_role("link", name="Dashboard").click()
    expect(auth_page).to_have_url(f"{base_url}/dashboard")
