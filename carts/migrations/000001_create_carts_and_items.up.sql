-- 000001_create_carts_and_items.up.sql
-- Initial schema for carts service: carts and cart_items tables, indexes,
-- constraints, no seed data (seeding removed; migrations handle schema only).
--
-- This migration is intended to be applied with golang-migrate (source=file).
--
-- Enable required extension only if desired (commented out because some
-- environments don't allow creating extensions). Uncomment if you manage DB.
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
--
-- Create carts table
CREATE TABLE IF NOT EXISTS carts (
    id SERIAL PRIMARY KEY,
    user_id BIGINT,
    session_id TEXT UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Trigger function to update updated_at timestamp on carts
CREATE OR REPLACE FUNCTION carts_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to keep updated_at current
DROP TRIGGER IF EXISTS trg_carts_set_timestamp ON carts;
CREATE TRIGGER trg_carts_set_timestamp
BEFORE UPDATE ON carts
FOR EACH ROW
EXECUTE FUNCTION carts_set_timestamp();

-- Create cart_items table
CREATE TABLE IF NOT EXISTS cart_items (
    id SERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    added_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes to support common lookups
CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts (user_id);
CREATE INDEX IF NOT EXISTS idx_carts_session_id ON carts (session_id);

CREATE INDEX IF NOT EXISTS idx_cart_items_product_id ON cart_items (product_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_cart_id ON cart_items (cart_id);

-- Unique constraint to ensure one row per (cart, product) for upsert semantics
CREATE UNIQUE INDEX IF NOT EXISTS uq_cart_items_cart_product
    ON cart_items (cart_id, product_id);
