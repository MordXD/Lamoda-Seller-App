-- +migrate Down
-- Удаляем таблицы в порядке, обратном зависимостям

DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS product_sales;
DROP TABLE IF EXISTS price_points;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS products; -- Зависит от users, но удаляем после product_*
DROP TABLE IF EXISTS account_links;
DROP TABLE IF EXISTS users;

-- Удаляем функцию для триггера
DROP FUNCTION IF EXISTS update_updated_at_column();