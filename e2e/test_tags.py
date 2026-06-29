import pytest
from playwright.sync_api import Page, expect

def test_tags_page_loads(auth_page: Page, base_url: str):
    """Test that the tags page loads with the correct heading."""
    auth_page.goto(f"{base_url}/tags")

    expect(auth_page.get_by_role("heading", name="Tags & Categories", exact=True)).to_be_visible()

def test_tags_eyebrow(auth_page: Page, base_url: str):
    """Test that the 'Organization' eyebrow text is displayed."""
    auth_page.goto(f"{base_url}/tags")

    expect(auth_page.get_by_text("Organization")).to_be_visible()

def test_tags_subtitle(auth_page: Page, base_url: str):
    """Test that the tags subtitle is present."""
    auth_page.goto(f"{base_url}/tags")

    expect(auth_page.get_by_text("Shape how transactions get labeled")).to_be_visible()

def test_new_category_button(auth_page: Page, base_url: str):
    """Test that the 'New Category' button is present."""
    auth_page.goto(f"{base_url}/tags")

    expect(auth_page.get_by_role("button", name="New Category")).to_be_visible()

def test_add_tag_buttons_present(auth_page: Page, base_url: str):
    """Test that per-category 'Add tag' buttons are present."""
    auth_page.goto(f"{base_url}/tags")

    # Each category has its own "Add tag to {category}" icon button
    add_buttons = auth_page.get_by_role("button", name="Add tag to")
    expect(add_buttons.first).to_be_visible()

def test_default_categories_present(auth_page: Page, base_url: str):
    """Test that default tag categories are rendered."""
    auth_page.goto(f"{base_url}/tags")

    # These are the default seeded categories
    for name in ["Financial", "Food & Drink", "Recurring", "Shopping", "Transport"]:
        expect(auth_page.get_by_text(name, exact=True).first).to_be_visible()
