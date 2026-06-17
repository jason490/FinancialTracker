-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS dashboard_layouts;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS tag_filters;
DROP TABLE IF EXISTS transaction_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS user_sso;
DROP TABLE IF EXISTS plaid_account;
DROP TABLE IF EXISTS plaid_items;
DROP TABLE IF EXISTS users;

-- Users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    first_name TEXT,
    last_name TEXT,
    password_hash TEXT,
    theme_preference TEXT DEFAULT 'system',
    created_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- SSO Logins linked to users
CREATE TABLE IF NOT EXISTS user_sso (
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
CREATE TABLE IF NOT EXISTS plaid_items (
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
CREATE TABLE IF NOT EXISTS plaid_account (
    id INTEGER PRIMARY KEY,
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
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY,
    plaid_id INTEGER NOT NULL,
    plaid_transaction_id TEXT UNIQUE NOT NULL,
    date INTEGER NOT NULL, -- Unix time
    amount REAL NOT NULL,
    name TEXT NOT NULL,
    merchant_name TEXT,
    plaid_category TEXT,
    pending INTEGER DEFAULT 0, -- Boolean 0/1
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (plaid_id) REFERENCES plaid_account(id) ON DELETE CASCADE
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

-- Sessions for authentication
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at INTEGER NOT NULL,
    reauthenticated_at INTEGER DEFAULT 0,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
