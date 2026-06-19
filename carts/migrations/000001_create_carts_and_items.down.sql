-- 000001_create_carts_and_items.down.sql
-- Revert the initial carts schema migration: drop tables, indexes, triggers and helper functions.
-- This file should be applied by golang-migrate as the down migration for 000001_create_carts_and_items.up.sql

-- Drop trigger that updates updated_at on carts (if present)
DROP TRIGGER IF EXISTS trg_carts_set_timestamp ON carts;

-- Drop the trigger function if it exists
DROP FUNCTION IF EXISTS carts_set_timestamp();

-- Drop indexes (IF EXISTS to be idempotent)
DROP INDEX IF EXISTS uq_cart_items_cart_product;
DROP INDEX IF EXISTS idx_cart_items_product_id;
DROP INDEX IF EXISTS idx_cart_items_cart_id;
DROP INDEX IF EXISTS idx_carts_user_id;
DROP INDEX IF EXISTS idx_carts_session_id;

-- Drop cart_items table (will cascade to any dependent objects)
DROP TABLE IF EXISTS cart_items;

-- Drop carts table
DROP TABLE IF EXISTS carts;
