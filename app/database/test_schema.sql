-- Development/test schema: drops everything for a clean reset, then recreates
-- the full schema. Keep this in sync with schema.sql (same tables/columns).

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS registration_codes;
DROP TABLE IF EXISTS password_reset_codes;
DROP TABLE IF EXISTS dashboard_layouts;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS tag_filters;
DROP TABLE IF EXISTS transaction_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS plaid_account;
DROP TABLE IF EXISTS plaid_items;
DROP TABLE IF EXISTS plaid_api_usage;
DROP TABLE IF EXISTS stripe_fc_account;
DROP TABLE IF EXISTS stripe_fc_items;
DROP TABLE IF EXISTS stripe_api_usage;
DROP TABLE IF EXISTS user_sso;
DROP TABLE IF EXISTS user_privileges;
DROP TABLE IF EXISTS users;

-- Users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    first_name TEXT,
    last_name TEXT,
    password_hash TEXT,
    theme_preference TEXT DEFAULT 'system',
    subscription_tier TEXT DEFAULT 'free' CHECK(subscription_tier IN ('free', 'plus', 'premium')),
    subscription_started_at INTEGER,
    stripe_customer_id TEXT,
    stripe_subscription_id TEXT,
    onboarding_completed INTEGER NOT NULL DEFAULT 1,
    created_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- Per-user subscription overrides (unlimited quotas, Stripe coupons).
CREATE TABLE user_privileges (
    user_id INTEGER PRIMARY KEY,
    unlimited_limits INTEGER NOT NULL DEFAULT 0,
    stripe_coupon_id TEXT,
    notes TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Stripe Financial Connections API call usage per user billing period
CREATE TABLE stripe_api_usage (
    user_id INTEGER NOT NULL,
    period_start INTEGER NOT NULL,
    call_count INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, period_start),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Stripe FC Items (grouped by institution)
CREATE TABLE stripe_fc_items (
    id INTEGER PRIMARY KEY,
    row_id TEXT UNIQUE NOT NULL,
    user_id INTEGER NOT NULL,
    institution_name TEXT NOT NULL,
    status TEXT DEFAULT 'active',
    error_code TEXT,
    last_synced INTEGER DEFAULT 0,
    transaction_refresh_id TEXT,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, institution_name)
);

-- Bank Accounts (Stripe Financial Connections)
CREATE TABLE stripe_fc_account (
    id INTEGER PRIMARY KEY,
    row_id TEXT UNIQUE NOT NULL,
    user_id INTEGER NOT NULL,
    stripe_account_id TEXT UNIQUE NOT NULL,
    stripe_item_row_id TEXT NOT NULL,
    name TEXT NOT NULL,
    mask TEXT,
    type TEXT,
    subtype TEXT,
    balance REAL,
    available_balance REAL,
    currency TEXT DEFAULT 'USD',
    status TEXT DEFAULT 'active',
    is_hidden BOOLEAN DEFAULT FALSE,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (stripe_item_row_id) REFERENCES stripe_fc_items(row_id) ON DELETE CASCADE
);

-- Plaid API call usage per user billing period (period_start is Unix time when the current cycle began)
CREATE TABLE plaid_api_usage (
    user_id INTEGER NOT NULL,
    period_start INTEGER NOT NULL,
    call_count INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, period_start),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- SSO Logins linked to users
CREATE TABLE user_sso (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    provider TEXT NOT NULL, -- 'google', 'apple', etc.
    sso_id TEXT NOT NULL,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(provider, sso_id),
    UNIQUE(user_id, provider)
);

-- Plaid Items (representing a connection to an institution)
CREATE TABLE plaid_items (
    id INTEGER PRIMARY KEY,
    row_id TEXT UNIQUE NOT NULL,
    user_id INTEGER NOT NULL,
    plaid_item_id TEXT UNIQUE NOT NULL,
    access_token TEXT NOT NULL,
    institution_id TEXT,
    institution_name TEXT,
    sync_cursor TEXT,
    status TEXT DEFAULT 'active',
    error_code TEXT,
    last_synced INTEGER DEFAULT 0,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Bank Accounts (Plaid)
CREATE TABLE plaid_account (
    id INTEGER PRIMARY KEY,
    row_id TEXT UNIQUE NOT NULL,
    user_id INTEGER NOT NULL,
    plaid_account_id TEXT UNIQUE NOT NULL,
    plaid_item_id TEXT NOT NULL,
    name TEXT NOT NULL,
    mask TEXT,
    type TEXT,
    subtype TEXT,
    balance REAL,
    available_balance REAL,
    currency TEXT DEFAULT 'USD',
    status TEXT DEFAULT 'active',
    is_hidden BOOLEAN DEFAULT FALSE,
    monthly_payment REAL DEFAULT 0,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (plaid_item_id) REFERENCES plaid_items(plaid_item_id) ON DELETE CASCADE
);

-- Transactions
CREATE TABLE transactions (
    id INTEGER PRIMARY KEY,
    provider TEXT NOT NULL DEFAULT 'plaid',
    plaid_id INTEGER NOT NULL,
    plaid_transaction_id TEXT UNIQUE NOT NULL,
    date INTEGER NOT NULL, -- Unix time
    amount REAL NOT NULL,
    name TEXT NOT NULL,
    merchant_name TEXT,
    plaid_category TEXT,
    pending INTEGER DEFAULT 0, -- Boolean 0/1
    created_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- Categories (Hierarchical grouping for tags)
CREATE TABLE categories (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, name)
);

-- Tags
CREATE TABLE tags (
    id INTEGER PRIMARY KEY,
    category_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    color TEXT NOT NULL DEFAULT 'blue',
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    UNIQUE(category_id, name)
);

-- Transaction-Tag mapping
CREATE TABLE transaction_tags (
    transaction_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (transaction_id, tag_id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Auto-tagging filters
CREATE TABLE tag_filters (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    pattern TEXT NOT NULL,
    filter_type TEXT CHECK(filter_type IN ('string', 'regex', 'amount_greater', 'amount_less', 'amount_equal')) NOT NULL,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- User dashboard widget layout (ViewPoint)
CREATE TABLE dashboard_layouts (
    user_id INTEGER PRIMARY KEY,
    layout_json TEXT NOT NULL,
    updated_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Temporary password reset codes (hashed, single-use)
CREATE TABLE password_reset_codes (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    code_hash TEXT NOT NULL,
    expires_at INTEGER NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0,
    used_at INTEGER,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_password_reset_codes_user_active
    ON password_reset_codes(user_id, used_at, expires_at);

-- Temporary registration invite codes (hashed, single-use; used when subscriptions are disabled)
CREATE TABLE registration_codes (
    id INTEGER PRIMARY KEY,
    code_hash TEXT NOT NULL,
    created_by_user_id INTEGER,
    expires_at INTEGER NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0,
    used_at INTEGER,
    used_by_user_id INTEGER,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (created_by_user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (used_by_user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_registration_codes_active
    ON registration_codes(used_at, expires_at);

-- Sessions for authentication
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at INTEGER NOT NULL,
    reauthenticated_at INTEGER DEFAULT 0,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
