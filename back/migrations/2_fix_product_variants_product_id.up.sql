-- +migrate Up

-- Исправляем тип product_id в таблице product_variants
-- Проверяем, существует ли таблица product_variants
DO $$
BEGIN
    -- Проверяем тип колонки product_id
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'product_variants' 
        AND column_name = 'product_id' 
        AND data_type = 'bigint'
    ) THEN
        -- Удаляем внешний ключ, если он существует
        ALTER TABLE product_variants DROP CONSTRAINT IF EXISTS fk_product;
        ALTER TABLE product_variants DROP CONSTRAINT IF EXISTS fk_products_variants;
        
        -- Удаляем индекс
        DROP INDEX IF EXISTS idx_product_variants_product_id;
        
        -- Очищаем таблицу, так как преобразование bigint -> uuid невозможно
        TRUNCATE TABLE product_variants;
        
        -- Изменяем тип колонки product_id с bigint на uuid
        ALTER TABLE product_variants ALTER COLUMN product_id TYPE uuid USING uuid_generate_v4();
        
        -- Восстанавливаем индекс
        CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);
        
        -- Восстанавливаем внешний ключ
        ALTER TABLE product_variants ADD CONSTRAINT fk_product FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE;
        
        RAISE NOTICE 'Successfully converted product_variants.product_id from bigint to uuid';
    ELSE
        RAISE NOTICE 'product_variants.product_id already has correct type or table does not exist';
    END IF;
END $$; 