package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lamoda-seller-app/internal/model"
	"gorm.io/gorm"
)

// ProductRepository инкапсулирует логику работы с продуктами в БД.
type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// ListProductsParams определяет все параметры для получения списка товаров.
type ListProductsParams struct {
	Search      string
	Category    string
	Brand       string
	MinPrice    float64
	MaxPrice    float64
	StockStatus string
	SortBy      string
	SortOrder   string
	Limit       int
	Offset      int
}

// FilterValues содержит данные для блока `filters` в ответе API.
type FilterValues struct {
	Categories  []FilterCount `json:"categories"`
	Brands      []FilterCount `json:"brands"`
	PriceRange  PriceRange    `json:"price_range"`
}

type FilterCount struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type PriceRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// List возвращает список товаров с фильтрацией, сортировкой, пагинацией, а также данные для фильтров.
func (r *ProductRepository) List(ctx context.Context, params ListProductsParams) ([]model.Product, int64, *FilterValues, error) {
	var products []model.Product
	var total int64
	var filters FilterValues

	// --- 1. Создаем базовый запрос с фильтрами ---
	query := r.db.WithContext(ctx).Model(&model.Product{})

	if params.Search != "" {
		searchQuery := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(sku) LIKE ? OR LOWER(brand) LIKE ?", searchQuery, searchQuery, searchQuery)
	}
	if params.Category != "" {
		// Для поддержки вложенных категорий, нужно найти все дочерние ID
		var categoryIDs []string
		// Рекурсивная функция для сбора всех ID дочерних категорий
		var findSubCategoryIDs func(parentID string)
		findSubCategoryIDs = func(parentID string) {
			categoryIDs = append(categoryIDs, parentID)
			var subIDs []string
			r.db.Model(&model.Category{}).Where("parent_id = ?", parentID).Pluck("id", &subIDs)
			for _, subID := range subIDs {
				findSubCategoryIDs(subID)
			}
		}
		findSubCategoryIDs(params.Category)
		query = query.Where("category_id IN (?)", categoryIDs)
	}
	if params.Brand != "" {
		query = query.Where("brand = ?", params.Brand)
	}
	if params.MinPrice > 0 {
		query = query.Where("price >= ?", params.MinPrice)
	}
	if params.MaxPrice > 0 {
		query = query.Where("price <= ?", params.MaxPrice)
	}
	switch params.StockStatus {
	case "in_stock":
		query = query.Where("total_stock > 10")
	case "low_stock":
		query = query.Where("total_stock > 0 AND total_stock <= 10")
	case "out_of_stock":
		query = query.Where("total_stock = 0")
	}

	// --- 2. Получаем данные для блока `filters` (до применения пагинации) ---
	err := r.calculateFilters(r.db.WithContext(ctx).Model(&model.Product{}), &filters, params)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to calculate filters: %w", err)
	}


	// --- 3. Считаем общее количество для пагинации ---
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, nil, err
	}

	// --- 4. Применяем сортировку ---
	if params.SortBy != "" {
		allowedSorts := map[string]string{
			"name": "name", "price": "price", "stock": "total_stock",
			"sales": "sales_count_30d",
			"created_date": "created_at",
		}
		dbColumn, ok := allowedSorts[params.SortBy]
		if ok {
			order := "ASC"
			if strings.ToLower(params.SortOrder) == "desc" {
				order = "DESC"
			}
			query = query.Order(fmt.Sprintf("%s %s", dbColumn, order))
		}
	} else {
		query = query.Order("created_at DESC")
	}

	// --- 5. Применяем пагинацию и получаем товары ---
	err = query.Preload("Images").Offset(params.Offset).Limit(params.Limit).Find(&products).Error
	if err != nil {
		return nil, 0, nil, err
	}

	return products, total, &filters, nil
}


// calculateFilters вычисляет доступные фильтры на основе текущего запроса.
// Важный момент: для подсчета фильтров часто нужна логика, отличная от основного запроса.
// Например, если выбрана категория "Платья", мы все равно хотим показать другие доступные категории.
// Здесь представлена упрощенная версия для демонстрации.
func (r *ProductRepository) calculateFilters(baseQuery *gorm.DB, filters *FilterValues, params ListProductsParams) error {
	// Категории: считаем количество товаров в каждой категории верхнего уровня
	// Мы клонируем базовый запрос, но сбрасываем фильтр по категории
	categoryQuery := baseQuery.Session(&gorm.Session{}) // Создаем новую сессию
	if params.Search != "" {
		searchQuery := "%" + strings.ToLower(params.Search) + "%"
		categoryQuery = categoryQuery.Where("LOWER(name) LIKE ? OR LOWER(sku) LIKE ?", searchQuery, searchQuery)
	}
	// ... можно применить и другие фильтры, кроме категории
	
	err := categoryQuery.Table("products").
		Select("categories.id, categories.name, count(products.id) as count").
		Joins("join categories on categories.id = products.category_id").
		Group("categories.id, categories.name").
		Scan(&filters.Categories).Error
	if err != nil { return err }

	// Бренды
	brandQuery := baseQuery.Session(&gorm.Session{}) // Новая сессия для брендов
	// Применяем все фильтры, КРОМЕ бренда
	if params.Category != "" {
		// ... логика поиска дочерних категорий, как в List ...
		brandQuery = brandQuery.Where("category_id = ?", params.Category)
	}
	
	rows, err := brandQuery.Select("brand as id, brand as name, count(*) as count").Group("brand").Rows()
	if err != nil { return err }
	defer rows.Close()
	for rows.Next() {
		var brand FilterCount
		if err := r.db.ScanRows(rows, &brand); err == nil {
			filters.Brands = append(filters.Brands, brand)
		}
	}

	// Диапазон цен (считаем на основе всех фильтров, кроме цены)
	priceQuery := baseQuery.Session(&gorm.Session{})
	if params.Category != "" {
		priceQuery = priceQuery.Where("category_id = ?", params.Category)
	}
	if params.Brand != "" {
		priceQuery = priceQuery.Where("brand = ?", params.Brand)
	}

	return priceQuery.Select("COALESCE(MIN(price), 0) as min, COALESCE(MAX(price), 0) as max").Row().Scan(&filters.PriceRange.Min, &filters.PriceRange.Max)
}

// GetByID возвращает товар по его ID со всеми связанными данными.
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).
		Preload("Images").
		Preload("Variants").
		Preload("Supplier").
		Preload("Category"). // Добавлено для загрузки информации о категории
		First(&product, "id = ?", id).Error
	return &product, err
}

// Create создает новый товар.
func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
	var totalStock int
	for _, v := range product.Variants {
		totalStock += v.Stock
	}
	product.TotalStock = totalStock
	
	return r.db.WithContext(ctx).Create(product).Error
}

// Update обновляет товар.
func (r *ProductRepository) Update(ctx context.Context, product *model.Product) error {
	var totalStock int
	for _, v := range product.Variants {
		totalStock += v.Stock
	}
	product.TotalStock = totalStock
	
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(product).Error
}

// Delete удаляет товар.
func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Select("Variants", "Images").Delete(&model.Product{ID: id}).Error
}

// CreateImage создает запись об изображении для продукта.
func (r *ProductRepository) CreateImage(ctx context.Context, image *model.ProductImage) error {
    return r.db.WithContext(ctx).Create(image).Error
}

// GetCategories получает все категории с их иерархией из базы данных.
func (r *ProductRepository) GetCategories(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	// Загружаем только категории верхнего уровня (у которых нет родителя)
	// и рекурсивно подгружаем их дочерние категории.
	// Preload("Subcategories.Subcategories") - для 3 уровней вложенности.
	err := r.db.WithContext(ctx).
		Where("parent_id IS NULL").
		Preload("Subcategories.Subcategories").
		Order("name ASC").
		Find(&categories).Error

	return categories, err
}

// GetSizeChart получает размерную сетку для указанной категории из базы данных.
func (r *ProductRepository) GetSizeChart(ctx context.Context, categoryID string) (*model.SizeChart, error) {
	var sizeChart model.SizeChart
	// Ищем размерную сетку по ID категории и подгружаем все связанные размеры
	err := r.db.WithContext(ctx).
		Preload("Sizes").
		Where("category_id = ?", categoryID).
		First(&sizeChart).Error

	// gorm.ErrRecordNotFound - стандартная ошибка, если запись не найдена,
	// ее удобно обрабатывать в сервисном слое для ответа 404 Not Found.
	return &sizeChart, err
}