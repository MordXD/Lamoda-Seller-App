// internal/repository/dashboard.go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lamoda-seller-app/internal/model" // Путь к вашим моделям
	"gorm.io/gorm"
)

// AggregatedData - это структура для сбора всех ключевых метрик за период.
type AggregatedData struct {
	Revenue        float64
	OrdersCount    int64
	ItemsSoldCount int64
	ReturnsCount   int64
	ReturnsSum     float64
}

// DashboardRepositoryInterface определяет контракт для работы с дашбордом
type DashboardRepositoryInterface interface {
	GetAggregatedData(ctx context.Context, userID uuid.UUID, start, end time.Time) (AggregatedData, error)
	GetTopCategories(ctx context.Context, userID uuid.UUID, start, end time.Time, limit int) ([]model.TopCategory, error)
	GetHourlySales(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]model.HourlySale, error)
	GetSalesChartData(ctx context.Context, userID uuid.UUID, start, end time.Time, granularity string) ([]model.SalesChartDataPoint, error)
}

var _ DashboardRepositoryInterface = (*DashboardRepository)(nil)

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

// GetAggregatedData получает сводные данные за указанный период для конкретного пользователя.
// ИСПРАВЛЕНО: Запросы теперь соответствуют схеме. Сумма берется из JSONB 'totals' и фильтруется по user_id.
func (r *DashboardRepository) GetAggregatedData(ctx context.Context, userID uuid.UUID, start, end time.Time) (AggregatedData, error) {
	var result AggregatedData

	// 1. Агрегируем данные из основной таблицы 'orders' с фильтром по user_id.
	err := r.db.WithContext(ctx).Model(&model.Order{}).
		Select(`
			COALESCE(SUM(CASE WHEN status = 'ordered' THEN (totals->>'total')::numeric ELSE 0 END), 0) as revenue,
			COUNT(DISTINCT CASE WHEN status = 'ordered' THEN id END) as orders_count,
			COUNT(DISTINCT CASE WHEN status = 'returned' THEN id END) as returns_count,
			COALESCE(SUM(CASE WHEN status = 'returned' THEN (totals->>'total')::numeric ELSE 0 END), 0) as returns_sum
		`).
		// ВАЖНО: фильтруем по user_id и created_at
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, start, end).
		Scan(&result).Error
	if err != nil {
		return AggregatedData{}, fmt.Errorf("failed to get aggregated data from orders: %w", err)
	}

	// 2. Отдельно и эффективно считаем количество проданных товаров с фильтром по user_id.
	err = r.db.WithContext(ctx).Model(&model.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.user_id = ? AND orders.status = 'ordered' AND orders.created_at BETWEEN ? AND ?", userID, start, end).
		Select("COALESCE(SUM(order_items.quantity), 0)").
		Scan(&result.ItemsSoldCount).Error

	if err != nil {
		return AggregatedData{}, fmt.Errorf("failed to get items sold count: %w", err)
	}

	return result, nil
}

// GetTopCategories получает топ-N категорий по выручке для конкретного пользователя.
// ИСПРАВЛЕНО: Запрос теперь соединяет order_items с products и группирует по текстовому полю products.category с фильтром по user_id.
func (r *DashboardRepository) GetTopCategories(ctx context.Context, userID uuid.UUID, start, end time.Time, limit int) ([]model.TopCategory, error) {
	var results []model.TopCategory

	err := r.db.WithContext(ctx).Model(&model.Order{}).
		Select(`
			products.category as category,
			products.category as name, -- Используем категорию и как слаг, и как имя
			SUM((orders.totals->>'total')::numeric) as revenue,
			COUNT(DISTINCT orders.id) as orders,
			SUM(order_items.quantity) as items
		`).
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		// Соединяем с таблицей products по product_id
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("orders.user_id = ? AND orders.status = ? AND orders.created_at BETWEEN ? AND ?", userID, "ordered", start, end).
		Group("products.category").
		Order("revenue DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get top categories: %w", err)
	}
	return results, nil
}

// GetHourlySales получает почасовую статистику для конкретного пользователя.
// ИСПРАВЛЕНО: Сумма берется из JSONB 'totals' с фильтром по user_id.
func (r *DashboardRepository) GetHourlySales(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]model.HourlySale, error) {
	var results []model.HourlySale

	hourExtractor := "EXTRACT(HOUR FROM created_at AT TIME ZONE 'UTC')"

	err := r.db.WithContext(ctx).Model(&model.Order{}).
		Select(
			fmt.Sprintf("%s as hour, SUM((totals->>'total')::numeric) as revenue, COUNT(*) as orders", hourExtractor),
		).
		Where("user_id = ? AND status = ? AND created_at BETWEEN ? AND ?", userID, "ordered", start, end).
		Group("hour").
		Order("hour ASC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get hourly sales: %w", err)
	}
	return results, nil
}

// GetSalesChartData получает данные для построения графика для конкретного пользователя.
// ИСПРАВЛЕНО: Все расчеты выручки используют JSONB-поле 'totals' с фильтром по user_id.
func (r *DashboardRepository) GetSalesChartData(ctx context.Context, userID uuid.UUID, start, end time.Time, granularity string) ([]model.SalesChartDataPoint, error) {
	var results []model.SalesChartDataPoint

	dateTruncFunc := fmt.Sprintf("DATE_TRUNC('%s', created_at AT TIME ZONE 'UTC')", granularity)

	err := r.db.WithContext(ctx).Model(&model.Order{}).
		Select(fmt.Sprintf(`
			%s as timestamp,
			COALESCE(SUM(CASE WHEN status = 'ordered' THEN (totals->>'total')::numeric ELSE 0 END), 0) as orders_revenue,
			COALESCE(SUM(CASE WHEN status = 'ordered' THEN (totals->>'total')::numeric ELSE 0 END), 0) - COALESCE(SUM(CASE WHEN status = 'returned' THEN (totals->>'total')::numeric ELSE 0 END), 0) as purchases_revenue,
			COUNT(DISTINCT CASE WHEN status = 'ordered' THEN id END) as orders_count,
			COUNT(DISTINCT CASE WHEN status = 'ordered' THEN id END) - COUNT(DISTINCT CASE WHEN status = 'returned' THEN id END) as purchases_count,
			COUNT(DISTINCT CASE WHEN status = 'returned' THEN id END) as return_count,
			COALESCE(SUM(CASE WHEN status = 'returned' THEN (totals->>'total')::numeric ELSE 0 END), 0) as return_revenue
		`, dateTruncFunc)).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, start, end).
		Group("timestamp").
		Order("timestamp ASC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get sales chart data: %w", err)
	}
	return results, nil
}
