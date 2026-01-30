-- Drop unique constraint for primary images
DROP INDEX IF EXISTS idx_product_images_one_primary_per_product;

-- Drop triggers
DROP TRIGGER IF EXISTS update_products_updated_at ON products;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;

-- Drop tables (order matters due to foreign key constraints)
DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;