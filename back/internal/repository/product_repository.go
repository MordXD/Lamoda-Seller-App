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
		query = query.Where("category = ?", params.Category)
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
		query = query.Where("total_stock > 10") // Примерное значение "в наличии"
	case "low_stock":
		query = query.Where("total_stock > 0 AND total_stock <= 10")
	case "out_of_stock":
		query = query.Where("total_stock = 0")
	}

	// --- 2. Получаем данные для блока `filters` (до применения пагинации) ---
	// Клонируем запрос, чтобы основные фильтры (кроме цены/категории/бренда) не влияли на подсчеты
	// Здесь упрощенная логика: мы считаем фильтры на основе уже отфильтрованного набора.
	// В реальном приложении это может быть сложнее.
	err := r.calculateFilters(query.Session(&gorm.Session{}), &filters)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to calculate filters: %w", err)
	}


	// --- 3. Считаем общее количество для пагинации ---
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, nil, err
	}

	// --- 4. Применяем сортировку ---
	if params.SortBy != "" {
		// Валидация, чтобы избежать SQL-инъекций
		allowedSorts := map[string]string{
			"name": "name", "price": "price", "stock": "total_stock",
			"sales": "sales_count_30d", // Потребует JOIN или хранимое поле
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
		query = query.Order("created_at DESC") // Сортировка по умолчанию
	}

	// --- 5. Применяем пагинацию и получаем товары ---
	// Preload для загрузки изображений, чтобы найти главное
	err = query.Preload("Images").Offset(params.Offset).Limit(params.Limit).Find(&products).Error
	if err != nil {
		return nil, 0, nil, err
	}

	return products, total, &filters, nil
}


// calculateFilters вычисляет доступные фильтры на основе текущего запроса.
func (r *ProductRepository) calculateFilters(query *gorm.DB, filters *FilterValues) error {
	// Категории
	rows, err := query.Select("category as id, category as name, count(*) as count").Group("category").Rows()
	if err != nil { return err }
	defer rows.Close()
	for rows.Next() {
		var cat FilterCount
		if err := r.db.ScanRows(rows, &cat); err == nil {
			filters.Categories = append(filters.Categories, cat)
		}
	}

	// Бренды
	rows, err = query.Select("brand as id, brand as name, count(*) as count").Group("brand").Rows()
	if err != nil { return err }
	defer rows.Close()
	for rows.Next() {
		var brand FilterCount
		if err := r.db.ScanRows(rows, &brand); err == nil {
			filters.Brands = append(filters.Brands, brand)
		}
	}

	// Диапазон цен
	return query.Select("COALESCE(MIN(price), 0) as min, COALESCE(MAX(price), 0) as max").Row().Scan(&filters.PriceRange.Min, &filters.PriceRange.Max)
}

// GetByID возвращает товар по его ID со всеми связанными данными.
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).
		Preload("Images").
		Preload("Variants").
		Preload("Supplier").
		First(&product, id).Error
	return &product, err
}

// Create создает новый товар.
func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
	// В реальном приложении логика расчета TotalStock и других полей
	// должна быть в транзакции или в сервисе.
	var totalStock int
	for _, v := range product.Variants {
		totalStock += v.Stock
	}
	product.TotalStock = totalStock
	
	return r.db.WithContext(ctx).Create(product).Error
}

// Update обновляет товар.
func (r *ProductRepository) Update(ctx context.Context, product *model.Product) error {
	// Логика обновления может быть сложной (например, обновление вариантов)
	// Здесь для простоты используется Save, который обновит все поля.
	var totalStock int
	for _, v := range product.Variants {
		totalStock += v.Stock
	}
	product.TotalStock = totalStock
	
	// Используем `Session` и `FullSave` для обновления ассоциаций
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(product).Error
}

// Delete удаляет товар.
func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// GORM автоматически удалит связанные записи (variants, images) если настроен foreign key constraint `onDelete:CASCADE`
	return r.db.WithContext(ctx).Select("Variants", "Images").Delete(&model.Product{ID: id}).Error
}

// CreateImage создает запись об изображении для продукта.
func (r *ProductRepository) CreateImage(ctx context.Context, image *model.ProductImage) error {
    return r.db.WithContext(ctx).Create(image).Error
}

// GetCategories (Заглушка)
func (r *ProductRepository) GetCategories(ctx context.Context) ([]model.Category, error) {
	// В реальном приложении это будет запрос к таблице категорий
	return []model.Category{
		{ID: "clothing", Name: "Одежда", Subcategories: []model.Category{
			{ID: "dresses", Name: "Платья", Subcategories: []model.Category{
				{ID: "mini_dress", Name: "Мини платья"}, {ID: "midi_dress", Name: "Миди платья"},
			}},
			{ID: "jeans", Name: "Джинсы", Subcategories: []model.Category{
				{ID: "skinny", Name: "Скинни"}, {ID: "wide_leg", Name: "Широкие"},
			}},
		}},
		{ID: "shoes", Name: "Обувь", Subcategories: []model.Category{
			{ID: "sneakers", Name: "Кроссовки"}, {ID: "boots", Name: "Сапоги"},
		}},
	}, nil
}

// GetSizeChart (Заглушка)
func (r *ProductRepository) GetSizeChart(ctx context.Context, category string) (*model.SizeChart, error) {
	// В реальном приложении это будет запрос к таблице с размерными сетками
	if category != "dresses" {
		return nil, gorm.ErrRecordNotFound
	}
	return &model.SizeChart{
		Category: "dresses",
		Type: "clothing",
		Sizes: []model.Size{
			{Size: "S", Measurements: map[string]string{"bust": "84-88 см", "waist": "64-68 см"}, International: "34", US: "4-6"},
			{Size: "M", Measurements: map[string]string{"bust": "88-92 см", "waist": "68-72 см"}, International: "36", US: "8-10"},
		},
	}, nil
}