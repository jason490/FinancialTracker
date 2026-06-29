import re

import pytest
from playwright.sync_api import Page, expect

def test_transactions_page_loads(auth_page: Page, base_url: str):
    """Test that the transactions page loads and displays the heading."""
    auth_page.goto(f"{base_url}/transactions")

    expect(auth_page.get_by_role("heading", name="Transactions", exact=True)).to_be_visible()

def test_transactions_eyebrow(auth_page: Page, base_url: str):
    """Test that the 'Activity' eyebrow text is displayed."""
    auth_page.goto(f"{base_url}/transactions")

    expect(auth_page.get_by_text("Activity")).to_be_visible()

def test_transactions_subtitle(auth_page: Page, base_url: str):
    """Test that the transactions subtitle is present."""
    auth_page.goto(f"{base_url}/transactions")

    expect(auth_page.get_by_text("Search, filter, and manage")).to_be_visible()

def test_transactions_sync_button(auth_page: Page, base_url: str):
    """Test that the sync transactions button is present and enabled."""
    auth_page.goto(f"{base_url}/transactions")

    sync_button = auth_page.get_by_role("button", name="Sync transactions")
    expect(sync_button).to_be_visible()
    expect(sync_button).to_be_enabled()

def test_transactions_search_input(auth_page: Page, base_url: str):
    """Test that the search input field is present."""
    auth_page.goto(f"{base_url}/transactions")

    search_input = auth_page.get_by_placeholder("Search transactions...")
    expect(search_input).to_be_visible()

def test_transactions_filter_toggle(auth_page: Page, base_url: str):
    """Test that the Filters toggle button is present and toggleable."""
    auth_page.goto(f"{base_url}/transactions")

    filter_button = auth_page.get_by_role("button", name="Filters")
    expect(filter_button).to_be_visible()
    expect(filter_button).to_have_attribute("aria-expanded", "false")

    # Toggle filters open
    filter_button.click()
    expect(filter_button).to_have_attribute("aria-expanded", "true")

    # Toggle filters closed
    filter_button.click()
    expect(filter_button).to_have_attribute("aria-expanded", "false")

def test_transactions_empty_state(auth_page: Page, base_url: str):
    """Test that the empty state is displayed when no transactions exist."""
    auth_page.goto(f"{base_url}/transactions")

    expect(auth_page.get_by_role("heading", name="No transactions found")).to_be_visible()
    expect(auth_page.get_by_text("Try adjusting your search")).to_be_visible()

def test_transactions_count_displayed(auth_page: Page, base_url: str):
    """Test that the transaction count is shown."""
    auth_page.goto(f"{base_url}/transactions")

    expect(auth_page.get_by_text(re.compile(r"\d+ transactions found"))).to_be_visible()
