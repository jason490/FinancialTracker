import pytest
from playwright.sync_api import Page, expect

def test_settings_page_loads(auth_page: Page, base_url: str):
    """Test that the settings page loads."""
    auth_page.goto(f"{base_url}/settings")

    expect(auth_page.get_by_role("heading", name="Settings")).to_be_visible()

def test_settings_eyebrow(auth_page: Page, base_url: str):
    """Test that the 'Preferences' eyebrow text is displayed."""
    auth_page.goto(f"{base_url}/settings")

    expect(auth_page.get_by_text("Preferences", exact=True)).to_be_visible()

def test_settings_subtitle(auth_page: Page, base_url: str):
    """Test that the settings subtitle is present."""
    auth_page.goto(f"{base_url}/settings")

    expect(auth_page.get_by_text("profile, secure sign-in")).to_be_visible()

def test_settings_tab_navigation(auth_page: Page, base_url: str):
    """Test that all settings tabs are present and navigable."""
    auth_page.goto(f"{base_url}/settings")

    tablist = auth_page.get_by_role("tablist", name="Settings sections")
    expect(tablist).to_be_visible()

    # Verify all three tabs exist
    expect(tablist.get_by_role("tab", name="Account")).to_be_visible()
    expect(tablist.get_by_role("tab", name="Connections")).to_be_visible()
    expect(tablist.get_by_role("tab", name="Appearance")).to_be_visible()

def test_settings_account_tab_default(auth_page: Page, base_url: str):
    """Test that the Account tab is selected by default."""
    auth_page.goto(f"{base_url}/settings")

    account_tab = auth_page.get_by_role("tab", name="Account")
    expect(account_tab).to_have_attribute("aria-selected", "true")

def test_profile_form_present(auth_page: Page, base_url: str):
    """Test that the profile settings form is present."""
    auth_page.goto(f"{base_url}/settings")

    expect(auth_page.get_by_role("heading", name="Profile")).to_be_visible()
    expect(auth_page.get_by_label("First name")).to_be_visible()
    expect(auth_page.get_by_label("Last name")).to_be_visible()
    expect(auth_page.get_by_label("Email")).to_be_visible()

    # Save button
    expect(auth_page.get_by_role("button", name="Save profile")).to_be_visible()

def test_profile_email_readonly(auth_page: Page, base_url: str):
    """Test that the email field is read-only."""
    auth_page.goto(f"{base_url}/settings")

    email_input = auth_page.locator("#email")
    expect(email_input).to_have_attribute("readonly", "")

def test_security_section_present(auth_page: Page, base_url: str):
    """Test that the Security section is present on the Account tab."""
    auth_page.goto(f"{base_url}/settings")

    expect(auth_page.get_by_role("heading", name="Security")).to_be_visible()
    expect(auth_page.get_by_text("Password login", exact=True)).to_be_visible()

def test_security_password_toggle(auth_page: Page, base_url: str):
    """Test that the Change/Add password button is present."""
    auth_page.goto(f"{base_url}/settings")

    # The button text varies: "Change password" or "Add password"
    password_button = auth_page.get_by_role("button", name="password")
    expect(password_button).to_be_visible()

def test_google_sso_section(auth_page: Page, base_url: str):
    """Test that Google SSO connection status is displayed."""
    auth_page.goto(f"{base_url}/settings")

    expect(auth_page.get_by_text("Google", exact=True)).to_be_visible()

def test_settings_connections_tab(auth_page: Page, base_url: str):
    """Test that clicking the Connections tab switches the panel."""
    auth_page.goto(f"{base_url}/settings")

    connections_tab = auth_page.get_by_role("tab", name="Connections")
    connections_tab.click()

    expect(connections_tab).to_have_attribute("aria-selected", "true")

def test_settings_appearance_tab(auth_page: Page, base_url: str):
    """Test that clicking the Appearance tab switches the panel."""
    auth_page.goto(f"{base_url}/settings")

    appearance_tab = auth_page.get_by_role("tab", name="Appearance")
    appearance_tab.click()

    expect(appearance_tab).to_have_attribute("aria-selected", "true")

def test_settings_tab_url_param(auth_page: Page, base_url: str):
    """Test that navigating with ?tab=connections selects the Connections tab."""
    auth_page.goto(f"{base_url}/settings?tab=connections")

    connections_tab = auth_page.get_by_role("tab", name="Connections")
    expect(connections_tab).to_have_attribute("aria-selected", "true")
