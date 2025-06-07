// internal/repository/dashboard.go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/lamoda-seller-app/internal/model" // Путь к вашим моделям
	"gorm.io/gorm"
)

// AggregatedData - это структура для сбора всех ключевых метрик за период.
type AggregatedData struct {
	Revenue      float64
	OrdersCount  int64
	ItemsSoldCount int64
	ReturnsCount int64
	ReturnsSum   float64
}

// DashboardRepositoryInterface определяет контракт для работы с дашбордом
type DashboardRepositoryInterface interface {
	GetAggregatedData(ctx context.Context, start, end time.Time) (AggregatedData, error)
	GetTopCategories(ctx context.Context, start, end time.Time, limit int) ([]model.TopCategory, error)
	GetHourlySales(ctx context.Context, start, end time.Time) ([]model.HourlySale, error)
	GetSalesChartData(ctx context.Context, start, end time.Time, granularity string) ([]model.SalesChartDataPoint, error)
}

var _ DashboardRepositoryInterface = (*DashboardRepository)(nil)

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

// GetAggregatedData получает сводные данные за указанный период.
func (r *DashboardRepository) GetAggregatedData(ctx context.Context, start, end time.Time) (AggregatedData, error) {
	var result AggregatedData

	// Предполагается, что у заказа есть статус (ordered, returned и т.д.)
	// и связь с OrderItem.
	// Этот запрос может потребовать адаптации под вашу реальную схему БД.
	err := r.db.WithContext(ctx).Model(&model.Order{}).
		Select(`
			COALESCE(SUM(CASE WHEN status = 'ordered' THEN amount ELSE 0 END), 0) as revenue,
			COUNT(DISTINCT CASE WHEN status = 'ordered' THEN id END) as orders_count,
			COALESCE(SUM(CASE WHEN status = 'ordered' THEN (SELECT SUM(quantity) FROM order_items WHERE order_items.order_id = orders.id) ELSE 0 END), 0) as items_sold_count,
			COUNT(DISTINCT CASE WHEN status = 'returned' THEN id END) as returns_count,
			COALESCE(SUM(CASE WHEN status = 'returned' THEN amount ELSE 0 END), 0) as returns_sum
		`).
		Where("created_at BETWEEN ? AND ?", start, end).
		Scan(&result).Error

	if err != nil {
		return AggregatedData{}, fmt.Errorf("failed to get aggregated data: %w", err)
	}
	return result, nil
}

// GetTopCategories получает топ-N категорий по выручке.
func (r *DashboardRepository) GetTopCategories(ctx context.Context, start, end time.Time, limit int) ([]model.TopCategory, error) {
	var results []model.TopCategory

	// Этот запрос предполагает наличие связей Order -> OrderItem -> Product -> Category.
	// Адаптируйте join'ы и имена таблиц/полей под вашу схему.
	err := r.db.WithContext(ctx).Model(&model.Order{}).
		Select(`
			categories.slug as category,
			categories.name as name,
			SUM(orders.amount) as revenue,
			COUNT(DISTINCT orders.id) as orders,
			SUM(order_items.quantity) as items
		`).
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Joins("JOIN products ON products.id = order_items.product_id").
		Joins("JOIN categories ON categories.id = products.category_id").
		Where("orders.status = ? AND orders.created_at BETWEEN ? AND ?", "ordered", start, end).
		Group("categories.slug, categories.name").
		Order("revenue DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get top categories: %w", err)
	}
	return results, nil
}

// GetHourlySales получает почасовую статистику.
func (r *DashboardRepository) GetHourlySales(ctx context.Context, start, end time.Time) ([]model.HourlySale, error) {
	var results []model.HourlySale

	// Функция EXTRACT(HOUR FROM ...) специфична для PostgreSQL.
	// Для MySQL используйте HOUR(created_at).
	// Для SQLite используйте strftime('%H', created_at).
	// Используем gorm.Expr для кросс-платформенности, если это возможно, или пишем сырой SQL.
	hourExtractor := "EXTRACT(HOUR FROM created_at AT TIME ZONE 'UTC')" // Пример для PostgreSQL

	err := r.db.WithContext(ctx).Model(&model.Order{}).
		Select(
			fmt.Sprintf("%s as hour, SUM(amount) as revenue, COUNT(*) as orders", hourExtractor),
		).
		Where("status = ? AND created_at BETWEEN ? AND ?", "ordered", start, end).
		Group("hour").
		Order("hour ASC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get hourly sales: %w", err)
	}
	return results, nil
}

// GetSalesChartData получает данные для построения графика.
func (r *DashboardRepository) GetSalesChartData(ctx context.Context, start, end time.Time, granularity string) ([]model.SalesChartDataPoint, error) {
	var results []model.SalesChartDataPoint

	// DATE_TRUNC специфичен для PostgreSQL.
	// Для MySQL: DATE_FORMAT(created_at, '%Y-%m-%d'), etc.
	// Для SQLite: strftime('%Y-%m-%d', created_at), etc.
	dateTruncFunc := fmt.Sprintf("DATE_TRUNC('%s', created_at AT TIME ZONE 'UTC')", granularity)

	err := r.db.WithContext(ctx).Model(&model.Order{}).
		Select(fmt.Sprintf(`
			%s as timestamp,
			COALESCE(SUM(CASE WHEN status = 'ordered' THEN amount ELSE 0 END), 0) as orders_revenue,
			COALESCE(SUM(CASE WHEN status = 'ordered' THEN amount ELSE -amount END), 0) as purchases_revenue,
			COUNT(DISTINCT CASE WHEN status = 'ordered' THEN id END) as orders_count,
			COUNT(DISTINCT CASE WHEN status = 'ordered' THEN id END) - COUNT(DISTINCT CASE WHEN status = 'returned' THEN id END) as purchases_count,
			COUNT(DISTINCT CASE WHEN status = 'returned' THEN id END) as return_count,
			COALESCE(SUM(CASE WHEN status = 'returned' THEN amount ELSE 0 END), 0) as return_revenue
		`, dateTruncFunc)).
		Where("created_at BETWEEN ? AND ?", start, end).
		Group("timestamp").
		Order("timestamp ASC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get sales chart data: %w", err)
	}
	return results, nil
}