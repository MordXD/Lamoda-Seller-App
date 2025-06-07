-- +migrate Up
-- Создаем функцию для автоматического обновления поля updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Таблица пользователей
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- Таблица связанных аккаунтов
CREATE TABLE account_links (
    id UUID PRIMARY KEY,
    primary_user_id UUID NOT NULL,
    linked_user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_primary_user FOREIGN KEY(primary_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_linked_user FOREIGN KEY(linked_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_primary_linked UNIQUE(primary_user_id, linked_user_id)
);

-- Таблица продуктов
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    old_price NUMERIC(10, 2),
    image_url VARCHAR(2048),
    short_desc VARCHAR(512),
    full_desc TEXT,
    brand_id INTEGER,
    category_id INTEGER,
    rating NUMERIC(3, 2) DEFAULT 0,
    rating_count INTEGER DEFAULT 0,
    in_stock INTEGER NOT NULL DEFAULT 0,
    tags TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индексы для продуктов
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_brand_id ON products(brand_id);
CREATE INDEX idx_products_category_id ON products(category_id);

CREATE TRIGGER update_products_updated_at
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- Таблица вариантов продукта
CREATE TABLE product_variants (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    color VARCHAR(100),
    size VARCHAR(50),
    sku VARCHAR(255) NOT NULL UNIQUE,
    in_stock INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT fk_product FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);
CREATE INDEX idx_product_variants_color ON product_variants(color);
CREATE INDEX idx_product_variants_size ON product_variants(size);

-- Таблица истории цен
CREATE TABLE price_points (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    CONSTRAINT fk_product FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_price_points_product_id ON price_points(product_id);

-- Таблица продаж продуктов
CREATE TABLE product_sales (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    date DATE NOT NULL,
    sales_count INTEGER NOT NULL,
    CONSTRAINT fk_product FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE,
    CONSTRAINT uq_product_date UNIQUE(product_id, date)
);

-- Таблица заказов
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    amount NUMERIC(12, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);