-- Drop triggers
DROP TRIGGER IF EXISTS update_cart_items_updated_at ON cart_items;
DROP TRIGGER IF EXISTS update_carts_updated_at ON carts;
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;

-- Drop unique constraint
DROP INDEX IF EXISTS idx_cart_items_unique_product_per_cart;

-- Drop tables (order matters due to foreign key constraints)
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;

-- Drop custom enum type
DROP TYPE IF EXISTS order_status;