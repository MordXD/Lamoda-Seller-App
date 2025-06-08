-- =================================================================
-- СКРИПТ ДЛЯ ЗАПОЛНЕНИЯ БАЗЫ ТЕСТОВЫМИ ДАННЫМИ (ФИНАЛЬНАЯ ВЕРСИЯ)
-- =================================================================

DO $$
DECLARE
    -- Конкретный ID пользователя, который будет создан
    user1_id UUID := 'ed16638d-57d8-4ad2-9081-e02d202cefe0';
    user1_email TEXT := 'kotkinegor8@gmail.com';
    user1_name TEXT := 'Egor';
    user1_balance BIGINT := 99999999; -- 999,999.99 рублей в копейках

    -- Поставщики
    supplier_zara_id UUID := 'c1e9f4b1-8b7a-4b6e-9c1d-2a3b4c5d6e7f';
    supplier_hm_id UUID   := 'd2e9f4b1-8b7a-4b6e-9c1d-2a3b4c5d6e7f';
    supplier_nike_id UUID := 'e3e9f4b1-8b7a-4b6e-9c1d-2a3b4c5d6e7f';

    -- Продукты
    product1_id UUID := 'f1a2b3c4-1234-5678-90ab-cdef12345678';
    product2_id UUID := 'f2a2b3c4-1234-5678-90ab-cdef12345678';
    product3_id UUID := 'f3a2b3c4-1234-5678-90ab-cdef12345678';
    product4_id UUID := 'f4a2b3c4-1234-5678-90ab-cdef12345678';
    product5_id UUID := 'f5a2b3c4-1234-5678-90ab-cdef12345678';
    product6_id UUID := 'f6a2b3c4-1234-5678-90ab-cdef12345678';

    -- Варианты Продуктов
    variant1_s_id UUID := 'a1000001-0000-0000-0000-000000000001';
    variant1_m_id UUID := 'a1000001-0000-0000-0000-000000000002';
    variant2_m_id UUID := 'a1000002-0000-0000-0000-000000000001';
    variant3_32_id UUID := 'a1000003-0000-0000-0000-000000000001';
    variant4_42_id UUID := 'a1000004-0000-0000-0000-000000000001';
    variant4_43_id UUID := 'a1000004-0000-0000-0000-000000000002';
    variant5_l_id UUID := 'a1000005-0000-0000-0000-000000000001';
    variant6_m_id UUID := 'a1000006-0000-0000-0000-000000000001';
    
    -- Заказы
    order1_id UUID := 'b1000001-0000-0000-0000-000000000001';
    order2_id UUID := 'b1000001-0000-0000-0000-000000000002';
    order3_id UUID := 'b1000001-0000-0000-0000-000000000003';
    order4_id UUID := 'b1000002-0000-0000-0000-000000000001';
    order5_id UUID := 'b1000002-0000-0000-0000-000000000002';

    -- Позиции в заказах
    order_item1_id UUID := 'c1000001-0000-0000-0000-000000000001';
    order_item2_id UUID := 'c1000001-0000-0000-0000-000000000002';
    order_item3_id UUID := 'c1000002-0000-0000-0000-000000000001';
    order_item4_id UUID := 'c1000003-0000-0000-0000-000000000001';
    order_item5_id UUID := 'c2000001-0000-0000-0000-000000000001';
    order_item6_id UUID := 'c2000001-0000-0000-0000-000000000002';
    order_item7_id UUID := 'c2000002-0000-0000-0000-000000000001';
    order_item8_id UUID := 'c2000002-0000-0000-0000-000000000002';

BEGIN

-- === 1. СОЗДАНИЕ ПОЛЬЗОВАТЕЛЯ ===
-- Создаем пользователя с конкретным ID
INSERT INTO users (id, name, email, hashed_password, balance_kopecks, created_at, updated_at) 
VALUES (
    user1_id, 
    user1_name, 
    user1_email, 
    '$2a$10$dummy.hash.for.testing.purposes.only', -- Dummy hash для тестирования
    user1_balance, 
    NOW(), 
    NOW()
) ON CONFLICT (id) DO UPDATE SET
    balance_kopecks = user1_balance,
    updated_at = NOW();

RAISE NOTICE 'Создан/обновлен пользователь % (%) с балансом % копеек', user1_name, user1_email, user1_balance;

-- === 2. ПОСТАВЩИКИ ===
INSERT INTO suppliers (id, name, contact) VALUES
(supplier_zara_id, 'ZARA Distribution', 'supply@zara.com'),
(supplier_hm_id, 'H&M Logistics', 'logistics@hm.com'),
(supplier_nike_id, 'Nike Europe', 'contact@nike.com')
ON CONFLICT (id) DO NOTHING;

-- === 3. ПРОДУКТЫ ===
-- Сначала удаляем старые продукты этого пользователя, если есть
DELETE FROM products WHERE id IN (product1_id, product2_id, product3_id, product4_id, product5_id, product6_id);

INSERT INTO products (id, name, brand, category, sku, price, cost_price, rating, reviews_count, supplier_id, user_id) VALUES
(product1_id, 'Пальто шерстяное классическое', 'ZARA', 'coats', 'CT001-MAIN', 25000, 12500, 4.7, 18, supplier_zara_id, user1_id),
(product2_id, 'Свитер кашемировый оверсайз', 'H&M', 'sweaters', 'SW001-MAIN', 18500, 9250, 4.5, 12, supplier_hm_id, user1_id),
(product3_id, 'Джинсы широкие с высокой посадкой', 'ZARA', 'jeans', 'JN001-MAIN', 8900, 4000, 4.8, 25, supplier_zara_id, user1_id),
(product4_id, 'Кроссовки беговые Air Zoom', 'Nike', 'shoes', 'NK001-MAIN', 14990, 7500, 4.9, 52, supplier_nike_id, user1_id),
(product5_id, 'Футболка спортивная Dri-FIT', 'Nike', 't-shirts', 'NK002-MAIN', 4500, 2000, 4.6, 31, supplier_nike_id, user1_id),
(product6_id, 'Леггинсы для фитнеса', 'Nike', 'leggings', 'NK003-MAIN', 7200, 3100, 4.7, 19, supplier_nike_id, user1_id);

-- === 4. ВАРИАНТЫ ПРОДУКТОВ ===
-- Удаляем старые варианты
DELETE FROM product_variants WHERE product_id IN (product1_id, product2_id, product3_id, product4_id, product5_id, product6_id);

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
-- Удаляем старые изображения
DELETE FROM product_images WHERE product_id IN (product1_id, product2_id, product3_id, product4_id, product5_id, product6_id);

INSERT INTO product_images (product_id, url, is_main) VALUES
(product1_id, 'https://example.com/images/coat.jpg', TRUE),
(product2_id, 'https://example.com/images/sweater.jpg', TRUE),
(product3_id, 'https://example.com/images/jeans.jpg', TRUE),
(product4_id, 'https://example.com/images/sneakers.jpg', TRUE),
(product5_id, 'https://example.com/images/tshirt.jpg', TRUE),
(product6_id, 'https://example.com/images/leggings.jpg', TRUE);

-- === 6. ЗАКАЗЫ ===
-- Удаляем старые заказы
DELETE FROM orders WHERE id IN (order1_id, order2_id, order3_id, order4_id, order5_id);

-- Все заказы привязаны к нашему пользователю с более свежими датами
INSERT INTO orders (id, user_id, order_number, date, status, totals, created_at) VALUES
(order1_id, user1_id, 'ORD-001', NOW() - INTERVAL '2 hours', 'ordered', '{"total": 62000.00}', NOW() - INTERVAL '2 hours'),
(order2_id, user1_id, 'ORD-002', NOW() - INTERVAL '1 day', 'ordered', '{"total": 25000.00}', NOW() - INTERVAL '1 day'),
(order3_id, user1_id, 'ORD-003', NOW() - INTERVAL '3 days', 'ordered', '{"total": 26700.00}', NOW() - INTERVAL '3 days'),
(order4_id, user1_id, 'ORD-004', NOW() - INTERVAL '6 hours', 'ordered', '{"total": 23990.00}', NOW() - INTERVAL '6 hours'),
(order5_id, user1_id, 'ORD-005', NOW() - INTERVAL '5 days', 'ordered', '{"total": 21700.00}', NOW() - INTERVAL '5 days');

-- === 7. ПОЗИЦИИ В ЗАКАЗАХ ===
-- Удаляем старые позиции
DELETE FROM order_items WHERE id IN (order_item1_id, order_item2_id, order_item3_id, order_item4_id, order_item5_id, order_item6_id, order_item7_id, order_item8_id);

INSERT INTO order_items (id, order_id, product_id, variant_id, name, brand, sku, size, quantity, price, cost_price, total) VALUES
-- Заказ 1
(order_item1_id, order1_id, product1_id, variant1_m_id, 'Пальто шерстяное классическое', 'ZARA', 'CT001-M', 'M', 1, 25000, 12500, 25000),
(order_item2_id, order1_id, product2_id, variant2_m_id, 'Свитер кашемировый оверсайз', 'H&M', 'SW001-M', 'M', 2, 18500, 9250, 37000),
-- Заказ 2
(order_item3_id, order2_id, product1_id, variant1_s_id, 'Пальто шерстяное классическое', 'ZARA', 'CT001-S', 'S', 1, 25000, 12500, 25000),
-- Заказ 3
(order_item4_id, order3_id, product3_id, variant3_32_id, 'Джинсы широкие с высокой посадкой', 'ZARA', 'JN001-32', '32', 3, 8900, 4000, 26700),
-- Заказ 4
(order_item5_id, order4_id, product4_id, variant4_42_id, 'Кроссовки беговые Air Zoom', 'Nike', 'NK001-42', '42', 1, 14990, 7500, 14990),
(order_item6_id, order4_id, product5_id, variant5_l_id, 'Футболка спортивная Dri-FIT', 'Nike', 'NK002-L', 'L', 2, 4500, 2000, 9000),
-- Заказ 5
(order_item7_id, order5_id, product4_id, variant4_43_id, 'Кроссовки беговые Air Zoom', 'Nike', 'NK001-43', '43', 1, 14500, 7200, 14500),
(order_item8_id, order5_id, product6_id, variant6_m_id, 'Леггинсы для фитнеса', 'Nike', 'NK003-M', 'M', 1, 7200, 3100, 7200);

-- === 8. ВОЗВРАТЫ ===
-- Удаляем старые возвраты
DELETE FROM order_returns WHERE order_item_id IN (order_item1_id, order_item2_id, order_item3_id, order_item4_id, order_item5_id, order_item6_id, order_item7_id, order_item8_id);

INSERT INTO order_returns (order_item_id, reason_code, quantity, returned_at, status) VALUES
-- Возврат по пальто (не подошел размер) - через день после заказа
(order_item3_id, 'size_mismatch', 1, NOW() - INTERVAL '2 days', 'completed'),
-- Возврат по джинсам (проблема с качеством) - через 2 дня после заказа
(order_item4_id, 'quality_issues', 1, NOW() - INTERVAL '2 days', 'completed'),
-- Возврат по кроссовкам (передумал) - через несколько часов после заказа
(order_item5_id, 'customer_changed_mind', 1, NOW() - INTERVAL '4 hours', 'completed'),
-- Возврат по футболкам (цвет не тот) - в процессе
(order_item6_id, 'color_difference', 1, NOW() - INTERVAL '3 hours', 'processing'),
-- Возврат по старым кроссовкам (качество) - через 3 дня после заказа
(order_item7_id, 'quality_issues', 1, NOW() - INTERVAL '3 days', 'completed');

RAISE NOTICE 'Тестовые данные успешно добавлены для пользователя %', user1_email;
RAISE NOTICE 'Добавлено: 6 продуктов, 5 заказов, 8 позиций заказов, 5 возвратов';
RAISE NOTICE 'Баланс пользователя: % копеек (%.2f рублей)', user1_balance, user1_balance::DECIMAL / 100;

END $$;
