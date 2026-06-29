import pytest
from playwright.sync_api import Page, expect

def test_forgot_password_page_structure(page: Page, base_url: str):
    """Test that the forgot-password page renders all expected elements."""
    page.goto(f"{base_url}/forgot-password")

    # Heading
    expect(page.get_by_role("heading", name="Forgot password")).to_be_visible()

    # Subtitle
    expect(page.get_by_text("send you a reset code")).to_be_visible()

    # Email field
    expect(page.get_by_label("Email address")).to_be_visible()

    # Submit button
    expect(page.get_by_role("button", name="Send reset code")).to_be_visible()

    # Back to login link
    expect(page.get_by_role("link", name="Back to sign in")).to_be_visible()

def test_forgot_password_email_required(page: Page, base_url: str):
    """Test that the email field is required."""
    page.goto(f"{base_url}/forgot-password")

    expect(page.locator("#email")).to_have_attribute("required", "")

def test_forgot_password_navigates_to_login(page: Page, base_url: str):
    """Test that clicking 'Back to sign in' returns to login page."""
    page.goto(f"{base_url}/forgot-password")

    page.get_by_role("link", name="Back to sign in").click()

    expect(page).to_have_url(f"{base_url}/login")

def test_forgot_password_advances_to_code_step(page: Page, base_url: str):
    """Test that submitting the email form advances to the reset code step."""
    page.goto(f"{base_url}/forgot-password")

    page.get_by_label("Email address").fill("test@test.com")
    page.get_by_role("button", name="Send reset code").click()

    expect(page.get_by_label("Reset code")).to_be_visible()
    expect(page.get_by_role("button", name="Continue")).to_be_visible()
    expect(page.get_by_text("Didn't receive a code?", exact=False)).to_be_visible()
    expect(page.get_by_role("button", name="Resend code")).to_be_visible()
    expect(page.get_by_text("Code expires in", exact=False)).to_be_visible()
