-- =================================================================
-- СКРИПТ ДЛЯ ЗАПОЛНЕНИЯ БАЗЫ ТЕСТОВЫМИ ДАННЫМИ
-- =================================================================

DO $$
DECLARE
    -- Пользователи (Продавцы)
    user1_id UUID := 'a19a79a5-5a45-420a-994b-118a14b35d61';
    user2_id UUID := 'b29a79a5-5a45-420a-994b-118a14b35d62';

    -- Поставщики
    supplier_zara_id UUID := 'c1e9f4b1-8b7a-4b6e-9c1d-2a3b4c5d6e7f';
    supplier_hm_id UUID   := 'd2e9f4b1-8b7a-4b6e-9c1d-2a3b4c5d6e7f';
    supplier_nike_id UUID := 'e3e9f4b1-8b7a-4b6e-9c1d-2a3b4c5d6e7f';

    -- Продукты
    product1_id UUID := 'f1a2b3c4-1234-5678-90ab-cdef12345678'; -- Пальто
    product2_id UUID := 'f2a2b3c4-1234-5678-90ab-cdef12345678'; -- Свитер
    product3_id UUID := 'f3a2b3c4-1234-5678-90ab-cdef12345678'; -- Джинсы
    product4_id UUID := 'f4a2b3c4-1234-5678-90ab-cdef12345678'; -- Кроссовки
    product5_id UUID := 'f5a2b3c4-1234-5678-90ab-cdef12345678'; -- Футболка
    product6_id UUID := 'f6a2b3c4-1234-5678-90ab-cdef12345678'; -- Леггинсы

    -- Варианты Продуктов
    variant1_s_id UUID := uuid_generate_v4();
    variant1_m_id UUID := uuid_generate_v4();
    variant2_m_id UUID := uuid_generate_v4();
    variant3_32_id UUID := uuid_generate_v4();
    variant4_42_id UUID := uuid_generate_v4();
    variant4_43_id UUID := uuid_generate_v4();
    variant5_l_id UUID := uuid_generate_v4();
    variant6_m_id UUID := uuid_generate_v4();
    
    -- Заказы
    order1_user1_id UUID := uuid_generate_v4();
    order2_user1_id UUID := uuid_generate_v4();
    order3_user1_id UUID := uuid_generate_v4();
    order4_user2_id UUID := uuid_generate_v4();
    order5_user2_id UUID := uuid_generate_v4();

    -- Позиции в заказах
    order_item1_id UUID := uuid_generate_v4();
    order_item2_id UUID := uuid_generate_v4();
    order_item3_id UUID := uuid_generate_v4();
    order_item4_id UUID := uuid_generate_v4();
    order_item5_id UUID := uuid_generate_v4();
    order_item6_id UUID := uuid_generate_v4();
    order_item7_id UUID := uuid_generate_v4();

BEGIN

-- === 1. ПОЛЬЗОВАТЕЛИ ===
INSERT INTO users (id, name, email, hashed_password, balance_kopecks) VALUES
(user1_id, 'Иван Петров', 'ivan.petrov@example.com', '$2a$10$...', 10000000),
(user2_id, 'Мария Сидорова', 'maria.sidorova@example.com', '$2a$10$...', 5000000);

-- === 2. ПОСТАВЩИКИ ===
INSERT INTO suppliers (id, name, contact) VALUES
(supplier_zara_id, 'ZARA Distribution', 'supply@zara.com'),
(supplier_hm_id, 'H&M Logistics', 'logistics@hm.com'),
(supplier_nike_id, 'Nike Europe', 'contact@nike.com');

-- === 3. ПРОДУКТЫ ===
INSERT INTO products (id, name, brand, category, sku, price, cost_price, rating, reviews_count, supplier_id) VALUES
(product1_id, 'Пальто шерстяное классическое', 'ZARA', 'coats', 'CT001-MAIN', 25000, 12500, 4.7, 18, supplier_zara_id),
(product2_id, 'Свитер кашемировый оверсайз', 'H&M', 'sweaters', 'SW001-MAIN', 18500, 9250, 4.5, 12, supplier_hm_id),
(product3_id, 'Джинсы широкие с высокой посадкой', 'ZARA', 'jeans', 'JN001-MAIN', 8900, 4000, 4.8, 25, supplier_zara_id),
(product4_id, 'Кроссовки беговые Air Zoom', 'Nike', 'shoes', 'NK001-MAIN', 14990, 7500, 4.9, 52, supplier_nike_id),
(product5_id, 'Футболка спортивная Dri-FIT', 'Nike', 't-shirts', 'NK002-MAIN', 4500, 2000, 4.6, 31, supplier_nike_id),
(product6_id, 'Леггинсы для фитнеса', 'Nike', 'leggings', 'NK003-MAIN', 7200, 3100, 4.7, 19, supplier_nike_id);

-- === 4. ВАРИАНТЫ ПРОДУКТОВ ===
INSERT INTO product_variants (id, product_id, sku, size, color, stock) VALUES
(variant1_s_id, product1_id, 'CT001-S', 'S', 'Бежевый', 10),
(variant1_m_id, product1_id, 'CT001-M', 'M', 'Бежевый', 15),
(variant2_m_id, product2_id, 'SW001-M', 'M', 'Серый', 20),
(variant3_32_id, product3_id, 'JN001-32', '32', 'Голубой', 25),
(variant4_42_id, product4_id, 'NK001-42', '42', 'Черный', 30),
(variant4_43_id, product4_id, 'NK001-43', '43', 'Черный', 12),
(variant5_l_id, product5_id, 'NK002-L', 'L', 'Белый', 50),
(variant6_m_id, product6_id, 'NK003-M', 'M', 'Черный', 40);

-- === 5. ИЗОБРАЖЕНИЯ ПРОДУКТОВ ===
INSERT INTO product_images (product_id, url, is_main) VALUES
(product1_id, 'https://example.com/images/coat.jpg', TRUE),
(product2_id, 'https://example.com/images/sweater.jpg', TRUE),
(product3_id, 'https://example.com/images/jeans.jpg', TRUE),
(product4_id, 'https://example.com/images/sneakers.jpg', TRUE),
(product5_id, 'https://example.com/images/tshirt.jpg', TRUE),
(product6_id, 'https://example.com/images/leggings.jpg', TRUE);

-- === 6. ЗАКАЗЫ ===
-- Заказы для Пользователя 1 (Иван Петров)
INSERT INTO orders (id, user_id, order_number, date, status) VALUES
(order1_user1_id, user1_id, 'ORD-001', NOW() - INTERVAL '5 days', 'completed'), -- Недавний
(order2_user1_id, user1_id, 'ORD-002', NOW() - INTERVAL '25 days', 'completed'), -- Внутри 30 дней
(order3_user1_id, user1_id, 'ORD-003', NOW() - INTERVAL '80 days', 'completed'); -- Внутри 90 дней
-- Заказы для Пользователя 2 (Мария Сидорова)
INSERT INTO orders (id, user_id, order_number, date, status) VALUES
(order4_user2_id, user2_id, 'ORD-004', NOW() - INTERVAL '2 days', 'completed'),
(order5_user2_id, user2_id, 'ORD-005', NOW() - INTERVAL '1 year', 'completed'); -- Старый заказ для проверки '1y'

-- === 7. ПОЗИЦИИ В ЗАКАЗАХ ===
INSERT INTO order_items (id, order_id, product_id, variant_id, name, brand, sku, size, quantity, price, cost_price, total) VALUES
-- Заказ 1 (Иван)
(order_item1_id, order1_user1_id, product1_id, variant1_m_id, 'Пальто шерстяное классическое', 'ZARA', 'CT001-M', 'M', 1, 25000, 12500, 25000),
(uuid_generate_v4(), order1_user1_id, product2_id, variant2_m_id, 'Свитер кашемировый оверсайз', 'H&M', 'SW001-M', 'M', 2, 18500, 9250, 37000),
-- Заказ 2 (Иван)
(order_item2_id, order2_user1_id, product1_id, variant1_s_id, 'Пальто шерстяное классическое', 'ZARA', 'CT001-S', 'S', 1, 25000, 12500, 25000),
-- Заказ 3 (Иван)
(order_item3_id, order3_user1_id, product3_id, variant3_32_id, 'Джинсы широкие с высокой посадкой', 'ZARA', 'JN001-32', '32', 3, 8900, 4000, 26700),
-- Заказ 4 (Мария)
(order_item4_id, order4_user2_id, product4_id, variant4_42_id, 'Кроссовки беговые Air Zoom', 'Nike', 'NK001-42', '42', 1, 14990, 7500, 14990),
(order_item5_id, order4_user2_id, product5_id, variant5_l_id, 'Футболка спортивная Dri-FIT', 'Nike', 'NK002-L', 'L', 2, 4500, 2000, 9000),
-- Заказ 5 (Мария)
(order_item6_id, order5_user2_id, product4_id, variant4_43_id, 'Кроссовки беговые Air Zoom', 'Nike', 'NK001-43', '43', 1, 14500, 7200, 14500),
(order_item7_id, order5_user2_id, product6_id, variant6_m_id, 'Леггинсы для фитнеса', 'Nike', 'NK003-M', 'M', 1, 7200, 3100, 7200);

-- === 8. ВОЗВРАТЫ ===
INSERT INTO order_returns (order_item_id, reason_code, quantity, returned_at, status) VALUES
-- Возврат по пальто (не подошел размер), 25 дней назад
(order_item2_id, 'size_mismatch', 1, NOW() - INTERVAL '24 days', 'completed'),
-- Возврат по джинсам (проблема с качеством), 75 дней назад
(order_item3_id, 'quality_issues', 1, NOW() - INTERVAL '75 days', 'completed'),
-- Возврат по кроссовкам (передумал), 1 день назад
(order_item4_id, 'customer_changed_mind', 1, NOW() - INTERVAL '1 day', 'completed'),
-- Возврат по футболкам (цвет не тот), 1 день назад
(order_item5_id, 'color_difference', 1, NOW() - INTERVAL '1 day', 'processing'),
-- Возврат по старым кроссовкам (износ, но причина - качество), 11 месяцев назад
(order_item6_id, 'quality_issues', 1, NOW() - INTERVAL '11 months', 'completed');

END $$;