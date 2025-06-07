-- +migrate Up

-- Включаем расширение для генерации UUID, если оно еще не включено.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создаем функцию для автоматического обновления поля updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';


-- =================================================================
-- Таблицы, связанные с пользователями (user.go, account_link.go)
-- =================================================================

-- Таблица пользователей (users)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    balance_kopecks BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();


-- Таблица связанных аккаунтов (account_links)
CREATE TABLE account_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    primary_user_id UUID NOT NULL,
    linked_user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_primary_user FOREIGN KEY(primary_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_linked_user FOREIGN KEY(linked_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_primary_linked UNIQUE(primary_user_id, linked_user_id)
);


-- =================================================================
-- Таблицы, связанные с продуктами (product.go)
-- =================================================================

-- Таблица поставщиков (suppliers)
CREATE TABLE suppliers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    contact VARCHAR(255)
);

-- Таблица продуктов (products)
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    brand VARCHAR(100),
    category VARCHAR(100),
    subcategory VARCHAR(100),
    sku VARCHAR(100) UNIQUE,
    barcode VARCHAR(100),
    price NUMERIC(12, 2) NOT NULL,
    cost_price NUMERIC(12, 2),
    currency VARCHAR(10) DEFAULT 'RUB',
    total_stock INTEGER NOT NULL DEFAULT 0,
    rating NUMERIC(3, 2) DEFAULT 0,
    reviews_count INTEGER NOT NULL DEFAULT 0,
    -- Поле return_rate в продукте теперь можно считать устаревшим или использовать как кеш.
    -- Реальные данные будут считаться из таблицы order_returns.
    return_rate NUMERIC(5, 2) DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    seasonal_demand VARCHAR(100),
    is_bestseller BOOLEAN NOT NULL DEFAULT FALSE,
    is_new BOOLEAN NOT NULL DEFAULT TRUE,
    discount_percent NUMERIC(5, 2) DEFAULT 0,
    tags TEXT[],
    material VARCHAR(255),
    care_instructions TEXT,
    country_origin VARCHAR(100),
    supplier_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_supplier FOREIGN KEY(supplier_id) REFERENCES suppliers(id) ON DELETE SET NULL
);

CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_status ON products(status);
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();


-- Таблица вариантов продукта (product_variants)
CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL,
    sku VARCHAR(100) NOT NULL UNIQUE,
    size VARCHAR(50),
    color VARCHAR(100),
    color_hex VARCHAR(20),
    stock INTEGER NOT NULL DEFAULT 0,
    reserved INTEGER NOT NULL DEFAULT 0,
    price NUMERIC(12, 2),
    weight NUMERIC(10, 2),
    dimensions JSONB,
    CONSTRAINT fk_product FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);


-- Таблица изображений продукта (product_images)
CREATE TABLE product_images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL,
    url VARCHAR(2048) NOT NULL,
    alt_text VARCHAR(255),
    is_main BOOLEAN NOT NULL DEFAULT FALSE,
    "order" INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT fk_product FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_images_product_id ON product_images(product_id);


-- =================================================================
-- Таблицы, связанные с заказами (order.go)
-- =================================================================

-- Таблица заказов (orders)
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID,
    customer_id UUID,
    order_number VARCHAR(100) UNIQUE NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) NOT NULL,
    notes TEXT,
    customer JSONB,
    delivery JSONB,
    payment JSONB,
    totals JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_date ON orders(date);
CREATE INDEX idx_orders_status ON orders(status);
CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();


-- Таблица позиций в заказе (order_items)
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL,
    product_id UUID NOT NULL,
    variant_id UUID NOT NULL,
    name VARCHAR(255),
    brand VARCHAR(100),
    sku VARCHAR(100),
    size VARCHAR(50),
    color VARCHAR(100),
    image VARCHAR(2048),
    quantity INTEGER NOT NULL,
    price NUMERIC(12, 2) NOT NULL,
    cost_price NUMERIC(12, 2),
    discount NUMERIC(12, 2),
    total NUMERIC(12, 2) NOT NULL,
    CONSTRAINT fk_order FOREIGN KEY(order_id) REFERENCES orders(id) ON DELETE CASCADE
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);


-- Таблица истории статусов заказа (status_histories)
CREATE TABLE status_histories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    comment TEXT,
    CONSTRAINT fk_order FOREIGN KEY(order_id) REFERENCES orders(id) ON DELETE CASCADE
);

CREATE INDEX idx_status_histories_order_id ON status_histories(order_id);


-- =================================================================
-- Таблицы, связанные с возвратами (для аналитики)
-- =================================================================

-- Справочник причин возврата
CREATE TABLE return_reasons (
    code VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

-- Наполняем справочник причинами из примера
INSERT INTO return_reasons (code, name) VALUES
('size_mismatch', 'Не подошел размер'),
('quality_issues', 'Проблемы с качеством'),
('color_difference', 'Цвет не соответствует'),
('damaged_delivery', 'Повреждение при доставке'),
('customer_changed_mind', 'Покупатель передумал'),
('other', 'Другое');


-- Таблица возвратов
CREATE TABLE order_returns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_item_id UUID NOT NULL,
    reason_code VARCHAR(50) NOT NULL,
    comment TEXT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    status VARCHAR(50) NOT NULL DEFAULT 'requested', -- e.g., requested, processing, completed, rejected
    returned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- Дата оформления возврата
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_order_item FOREIGN KEY(order_item_id) REFERENCES order_items(id) ON DELETE CASCADE,
    CONSTRAINT fk_reason FOREIGN KEY(reason_code) REFERENCES return_reasons(code) ON DELETE RESTRICT
);

-- Индексы для ускорения аналитических запросов
CREATE INDEX idx_order_returns_order_item_id ON order_returns(order_item_id);
CREATE INDEX idx_order_returns_reason_code ON order_returns(reason_code);
CREATE INDEX idx_order_returns_returned_at ON order_returns(returned_at);

-- Триггер для обновления updated_at
CREATE TRIGGER update_order_returns_updated_at
BEFORE UPDATE ON order_returns
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();