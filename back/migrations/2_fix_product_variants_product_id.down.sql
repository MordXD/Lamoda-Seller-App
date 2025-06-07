-- +migrate Down

-- Откатываем изменения типа product_id в таблице product_variants
DO $$
BEGIN
    -- Проверяем тип колонки product_id
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'product_variants' 
        AND column_name = 'product_id' 
        AND data_type = 'uuid'
    ) THEN
        -- Удаляем внешний ключ
        ALTER TABLE product_variants DROP CONSTRAINT IF EXISTS fk_product;
        ALTER TABLE product_variants DROP CONSTRAINT IF EXISTS fk_products_variants;
        
        -- Удаляем индекс
        DROP INDEX IF EXISTS idx_product_variants_product_id;
        
        -- Очищаем таблицу
        TRUNCATE TABLE product_variants;
        
        -- Возвращаем тип колонки product_id обратно к bigint
        ALTER TABLE product_variants ALTER COLUMN product_id TYPE bigint USING 0;
        
        -- Восстанавливаем индекс
        CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);
        
        -- Внешний ключ не восстанавливаем, так как он будет некорректным
        
        RAISE NOTICE 'Successfully reverted product_variants.product_id from uuid to bigint';
    ELSE
        RAISE NOTICE 'product_variants.product_id already has bigint type or table does not exist';
    END IF;
END $$; 