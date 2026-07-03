-- 0001: Add the hidden-tag flag.
-- Transactions carrying a hidden tag are excluded from income/spending totals
-- (e.g. credit-card payments and internal transfers). This migration is for
-- databases created before the flag existed; brand-new databases already include
-- the column via schema.sql and baseline this migration as applied on first boot.
ALTER TABLE tags ADD COLUMN is_hidden INTEGER NOT NULL DEFAULT 0;
