-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (order matters due to foreign key constraints)
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;