-- +migrate Up

-- Добавляем поле user_id в таблицу products
ALTER TABLE products ADD COLUMN user_id UUID;

-- Добавляем внешний ключ
ALTER TABLE products ADD CONSTRAINT fk_products_user 
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL;

-- Добавляем индекс для ускорения запросов
CREATE INDEX idx_products_user_id ON products(user_id); 