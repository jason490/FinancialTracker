import pytest
import requests
from playwright.sync_api import Page, expect


def registration_code_required(base_url: str) -> bool:
    """Returns whether the server requires invite codes to register."""
    try:
        response = requests.get(f"{base_url}/api/v1/auth/registration-config", timeout=5)
        response.raise_for_status()
        return response.json().get("registration_code_required") is True
    except requests.RequestException:
        return False


def test_register_page_structure(page: Page, base_url: str):
    """Test that the registration page renders all expected elements."""
    page.goto(f"{base_url}/register")

    # Heading
    expect(page.get_by_role("heading", name="Create account")).to_be_visible()

    # Form fields
    expect(page.get_by_label("First name")).to_be_visible()
    expect(page.get_by_label("Last name")).to_be_visible()
    expect(page.get_by_label("Email address")).to_be_visible()
    expect(page.get_by_label("Password", exact=True)).to_be_visible()
    expect(page.get_by_label("Confirm password")).to_be_visible()

    if registration_code_required(base_url):
        expect(page.get_by_label("Registration code")).to_be_visible()
        expect(page.get_by_text("Registration is invite-only")).to_be_visible()

    # Submit button
    expect(page.get_by_role("button", name="Create account")).to_be_visible()

    # Password requirements hint
    expect(page.get_by_text("Password requirements")).to_be_visible()

    # SSO button
    expect(page.get_by_role("link", name="Sign up with Google")).to_be_visible()

    # Footer link back to login
    expect(page.get_by_role("link", name="Sign in")).to_be_visible()


def test_register_fields_required(page: Page, base_url: str):
    """Test that all registration fields are required."""
    page.goto(f"{base_url}/register")

    required_fields = ["first_name", "last_name", "email", "password", "confirm_password"]
    if registration_code_required(base_url):
        required_fields.append("registration_code")

    for field_id in required_fields:
        expect(page.locator(f"#{field_id}")).to_have_attribute("required", "")


def test_register_navigates_to_login(page: Page, base_url: str):
    """Test that clicking 'Sign in' navigates to the login page."""
    page.goto(f"{base_url}/register")

    page.get_by_role("link", name="Sign in").click()

    expect(page).to_have_url(f"{base_url}/login")
