package repository

import (
	"context"
	"errors"
	"time"

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

// GetAllParams определяет параметры для получения списка товаров (пагинация, фильтры).
type GetAllProductsParams struct {
	Limit  int
	Offset int
	// Сюда можно будет добавить поля для фильтрации и сортировки
}

// Create создает новый товар и его варианты в одной транзакции.
func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Сначала создаем основной продукт
		if err := tx.Create(product).Error; err != nil {
			return err
		}

		// Если есть варианты, связываем их с созданным продуктом и создаем
		if len(product.Variants) > 0 {
			for i := range product.Variants {
				product.Variants[i].ProductID = product.ID
			}
			if err := tx.Create(&product.Variants).Error; err != nil {
				return err
			}
		}

		// Добавляем текущую цену в историю цен
		pricePoint := model.PricePoint{
			ProductID: product.ID,
			Date:      time.Now(),
			Price:     product.Price,
		}
		if err := tx.Create(&pricePoint).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetAll возвращает список всех товаров с пагинацией.
func (r *ProductRepository) GetAll(ctx context.Context, params GetAllProductsParams) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{})

	// Считаем общее количество для пагинации
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Применяем пагинацию и получаем сами товары
	err := query.Offset(params.Offset).Limit(params.Limit).Order("created_at DESC").Find(&products).Error
	return products, total, err
}

// GetByID возвращает товар по его ID, включая связанные варианты.
func (r *ProductRepository) GetByID(ctx context.Context, id int) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).Preload("Variants").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &product, nil
}

// Update обновляет информацию о товаре.
// Внимание: эта реализация заменяет все варианты товара на новые.
func (r *ProductRepository) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Проверяем, изменилась ли цена, чтобы добавить запись в историю
		var oldPrice float64
		if err := tx.Model(&model.Product{}).Where("id = ?", product.ID).Select("price").Row().Scan(&oldPrice); err != nil {
			return err
		}

		// Удаляем старые варианты
		if err := tx.Where("product_id = ?", product.ID).Delete(&model.ProductVariant{}).Error; err != nil {
			return err
		}

		// Сохраняем основные данные о товаре (включая, возможно, новую цену)
		if err := tx.Save(product).Error; err != nil {
			return err
		}
		
		// Создаем новые варианты, если они есть
		if len(product.Variants) > 0 {
			for i := range product.Variants {
				product.Variants[i].ProductID = product.ID
			}
			if err := tx.Create(&product.Variants).Error; err != nil {
				return err
			}
		}

		// Если цена изменилась, добавляем запись в историю
		if oldPrice != product.Price {
			pricePoint := model.PricePoint{
				ProductID: product.ID,
				Date:      time.Now(),
				Price:     product.Price,
			}
			if err := tx.Create(&pricePoint).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete удаляет товар по ID. Каскадное удаление позаботится о вариантах.
func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, id).Error
}

// GetSalesStats возвращает статистику продаж за указанный период.
func (r *ProductRepository) GetSalesStats(ctx context.Context, id int, period string) ([]model.ProductSales, error) {
	var sales []model.ProductSales
	var startDate time.Time
	now := time.Now()

	switch period {
	case "day":
		startDate = now.AddDate(0, 0, -1)
	case "week":
		startDate = now.AddDate(0, 0, -7)
	case "month":
		startDate = now.AddDate(0, -1, 0)
	default:
		// по умолчанию неделя
		startDate = now.AddDate(0, 0, -7)
	}

	err := r.db.WithContext(ctx).
		Where("product_id = ? AND date >= ?", id, startDate).
		Order("date ASC").
		Find(&sales).Error

	return sales, err
}

// GetPriceHistory возвращает историю изменения цен для товара.
func (r *ProductRepository) GetPriceHistory(ctx context.Context, id int) ([]model.PricePoint, error) {
	var history []model.PricePoint
	err := r.db.WithContext(ctx).
		Where("product_id = ?", id).
		Order("date DESC").
		Find(&history).Error
	return history, err
}