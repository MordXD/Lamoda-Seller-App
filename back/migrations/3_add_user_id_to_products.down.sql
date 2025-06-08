-- +migrate Down

-- Удаляем индекс
DROP INDEX IF EXISTS idx_products_user_id;

-- Удаляем внешний ключ
ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_user;

-- Удаляем поле user_id
ALTER TABLE products DROP COLUMN IF EXISTS user_id; 