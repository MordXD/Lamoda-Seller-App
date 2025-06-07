// internal/repository/analytics_repository.go
package repository

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/lamoda-seller-app/internal/model" // Используем ваш путь к моделям
	"gorm.io/gorm"
)

// AnalyticsRepo определяет интерфейс для аналитического репозитория.
type AnalyticsRepo interface {
	GetTopProducts(ctx context.Context, userID uuid.UUID, params model.TopProductsRequestParams) (*model.TopProductsResponse, error)
	GetCategoryAnalytics(ctx context.Context, userID uuid.UUID, params model.CategoryAnalyticsRequestParams) (*model.CategoryAnalyticsResponse, error)
	GetSizeDistribution(ctx context.Context, userID uuid.UUID, params model.SizeDistributionRequestParams) (*model.SizeDistributionResponse, error)
	GetSeasonalTrends(ctx context.Context, userID uuid.UUID, params model.SeasonalTrendsRequestParams) (*model.SeasonalTrendsResponse, error)
	GetReturnsAnalytics(ctx context.Context, userID uuid.UUID, params model.ReturnsAnalyticsRequestParams) (*model.ReturnsAnalyticsResponse, error)
}

type analyticsRepository struct {
	db *gorm.DB
}

// NewAnalyticsRepository создает новый экземпляр аналитического репозитория.
func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepo {
	return &analyticsRepository{db: db}
}

// --- Helper Functions ---

// getStartDate преобразует строковый период (7d, 30d...) в дату начала.
func getStartDate(period string) time.Time {
	now := time.Now()
	switch period {
	case "7d":
		return now.AddDate(0, 0, -7)
	case "30d":
		return now.AddDate(0, 0, -30)
	case "90d":
		return now.AddDate(0, -3, 0)
	case "1y":
		return now.AddDate(-1, 0, 0)
	case "all_time":
		return time.Time{} // Нулевое время для "всего времени"
	default:
		return now.AddDate(0, 0, -30) // По умолчанию 30 дней
	}
}

// --- Top Products ---

func (r *analyticsRepository) GetTopProducts(ctx context.Context, userID uuid.UUID, params model.TopProductsRequestParams) (*model.TopProductsResponse, error) {
	// Устанавливаем значения по умолчанию
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Period == "" {
		params.Period = "30d"
	}
	if params.Metric == "" {
		params.Metric = "revenue"
	}

	startDate := getStartDate(params.Period)

	// Определяем поле для сортировки
	// Важно: проверяем значение, чтобы избежать SQL-инъекций при форматировании строки.
	orderByClause := "revenue DESC"
	switch params.Metric {
	case "quantity":
		orderByClause = "sales_count DESC"
	case "profit":
		orderByClause = "profit DESC"
	}

	// Базовый запрос для дальнейшего использования
	baseQuery := r.db.WithContext(ctx).
		Table("order_items oi").
		Joins("JOIN orders o ON oi.order_id = o.id").
		Joins("JOIN products p ON oi.product_id = p.id").
		Where("o.user_id = ? AND o.status = ? AND o.date >= ?", userID, "completed", startDate)

	if params.Category != "" {
		baseQuery = baseQuery.Where("p.category = ?", params.Category)
	}

	// 1. Получаем общую сводку
	var summary model.TopProductsSummary
	err := baseQuery.Select(`
		SUM(oi.total) as total_revenue,
		SUM(oi.quantity) as total_sales,
		SUM(oi.total - (COALESCE(oi.cost_price, 0) * oi.quantity)) as total_profit
	`).Scan(&summary).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query top products summary: %w", err)
	}
	if summary.TotalRevenue > 0 {
		summary.AvgMargin = (summary.TotalProfit / summary.TotalRevenue) * 100
	}

	// 2. Получаем список продуктов
	var products []model.TopProduct
	err = baseQuery.
		Select(`
			p.id,
			p.name,
			p.brand,
			p.category,
			p.sku,
			p.price,
			p.cost_price,
			p.rating,
			p.reviews_count,
			(SELECT url FROM product_images pi WHERE pi.product_id = p.id AND pi.is_main = TRUE LIMIT 1) as image,
			SUM(oi.quantity) AS sales_count,
			SUM(oi.total) AS revenue,
			SUM(oi.total - (COALESCE(oi.cost_price, 0) * oi.quantity)) AS profit,
			CASE WHEN SUM(oi.total) > 0 THEN (SUM(oi.total - (COALESCE(oi.cost_price, 0) * oi.quantity)) / SUM(oi.total)) * 100 ELSE 0 END AS margin_percent
		`).
		Group("p.id, p.name, p.brand, p.category, p.sku, p.price, p.cost_price, p.rating, p.reviews_count").
		Order(orderByClause).
		Limit(params.Limit).
		Scan(&products).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query top products list: %w", err)
	}

	// 3. Обогащаем продукты вычисляемыми и мокированными данными
	for i := range products {
		p := &products[i]
		p.Rank = i + 1
		p.RankChange = (p.Rank % 3) - 1 // Mocked rank change
		if summary.TotalRevenue > 0 {
			p.RevenueShare = (p.Revenue / summary.TotalRevenue) * 100
		}
		// Мокированные данные, т.к. нет таблицы возвратов
		p.ReturnCount = p.SalesCount / 12
		p.ReturnRate = (1.0 / 12.0) * 100
	}

	// Мок для TopCategory
	summary.TopCategory = model.TopCategorySummary{
		Category: "coats",
		Name:     "Пальто",
		Revenue:  summary.TotalRevenue * 0.3,
		Share:    30.6,
	}

	return &model.TopProductsResponse{
		Period:   params.Period,
		Metric:   params.Metric,
		Products: products,
		Summary:  summary,
	}, nil
}

// --- Category Analytics ---

func (r *analyticsRepository) GetCategoryAnalytics(ctx context.Context, userID uuid.UUID, params model.CategoryAnalyticsRequestParams) (*model.CategoryAnalyticsResponse, error) {
	if params.Period == "" {
		params.Period = "30d"
	}
	if params.SortBy == "" {
		params.SortBy = "revenue"
	}
	startDate := getStartDate(params.Period)

	orderByClause := "revenue DESC"
	switch params.SortBy {
	case "quantity":
		orderByClause = "quantity DESC"
	case "profit":
		orderByClause = "profit DESC"
	case "margin":
		orderByClause = "margin_percent DESC"
	}

	var categories []model.CategoryAnalytics
	err := r.db.WithContext(ctx).
		Table("order_items oi").
		Joins("JOIN orders o ON oi.order_id = o.id").
		Joins("JOIN products p ON oi.product_id = p.id").
		Where("o.user_id = ? AND o.status = ? AND o.date >= ?", userID, "completed", startDate).
		Select(`
			p.category as id,
			p.category as name,
			SUM(oi.total) AS revenue,
			SUM(oi.quantity) AS quantity,
			SUM(oi.total - (COALESCE(oi.cost_price, 0) * oi.quantity)) AS profit,
			COUNT(DISTINCT oi.order_id) as orders_count,
			CASE WHEN SUM(oi.quantity) > 0 THEN SUM(oi.total) / SUM(oi.quantity) ELSE 0 END as avg_price,
			CASE WHEN SUM(oi.total) > 0 THEN (SUM(oi.total - (COALESCE(oi.cost_price, 0) * oi.quantity)) / SUM(oi.total)) * 100 ELSE 0 END AS margin_percent
		`).
		Group("p.category").
		Order(orderByClause).
		Scan(&categories).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query category analytics: %w", err)
	}

	// Вычисляем totalRevenue для расчета revenue_share
	var totalRevenue float64
	for _, cat := range categories {
		totalRevenue += cat.Revenue
	}

	// Добавляем вычисляемые и мокированные данные
	topColors := []string{"blue", "black", "white", "green", "red"}
	topSizes := []string{"M", "S", "L", "XL"}
	trends := []string{"stable", "growing", "declining"}

	for i := range categories {
		cat := &categories[i]
		if totalRevenue > 0 {
			cat.RevenueShare = (cat.Revenue / totalRevenue) * 100
		}
		cat.GrowthPercent = math.Round((rand.Float64()*20-5)*10) / 10 // от -5 до +15
		cat.ReturnRate = math.Round((rand.Float64()*15+5)*10) / 10   // от 5 до 20
		cat.TopColor = topColors[rand.Intn(len(topColors))]
		cat.TopSize = topSizes[rand.Intn(len(topSizes))]
		cat.SeasonalTrend = trends[rand.Intn(len(trends))]
		cat.Bestseller = model.CategoryBestseller{
			ID:    uuid.New(),
			Name:  "Mock Bestseller for " + cat.Name,
			Sales: cat.Quantity / 2,
		}
	}

	// Мокированный summary
	summary := model.CategoryAnalyticsSummary{
		TotalCategories: len(categories),
		FastestGrowing:  model.CategoryGrowthSummary{Category: "jeans", Growth: 15.3},
		HighestMargin:   model.CategoryMarginSummary{Category: "accessories", Margin: 65.2},
		HighestReturnRate: model.CategoryReturnRateSummary{
			Category:   "shoes",
			ReturnRate: 22.1,
		},
	}

	return &model.CategoryAnalyticsResponse{
		Period:     params.Period,
		Categories: categories,
		Summary:    summary,
	}, nil
}

// --- Size Distribution ---

func (r *analyticsRepository) GetSizeDistribution(ctx context.Context, userID uuid.UUID, params model.SizeDistributionRequestParams) (*model.SizeDistributionResponse, error) {
	if params.Period == "" {
		params.Period = "30d"
	}
	startDate := getStartDate(params.Period)

	// Создаем динамический запрос
	tx := r.db.WithContext(ctx).
		Table("order_items oi").
		Joins("JOIN orders o ON oi.order_id = o.id").
		Joins("JOIN products p ON oi.product_id = p.id").
		Where("o.user_id = ? AND o.status = ? AND o.date >= ? AND oi.size IS NOT NULL AND oi.size != ''", userID, "completed", startDate)

	if params.Category != "" && params.Category != "all" {
		tx = tx.Where("p.category = ?", params.Category)
	}
	if params.ProductID != "" {
		tx = tx.Where("oi.product_id = ?", params.ProductID)
	}

	// 1. Получаем общее количество проданных товаров для расчета процента
	var totalSalesCount int64
	if err := tx.Select("SUM(oi.quantity)").Scan(&totalSalesCount).Error; err != nil {
		return nil, fmt.Errorf("failed to get total sales for size distribution: %w", err)
	}

	// 2. Получаем распределение по размерам
	var distribution []model.SizeDistribution
	if err := tx.
		Select(`
			oi.size,
			SUM(oi.quantity) as sales_count,
			SUM(oi.total) as revenue,
			AVG(oi.price) as avg_price
		`).
		Group("oi.size").
		Order("sales_count DESC").
		Scan(&distribution).Error; err != nil {
		return nil, fmt.Errorf("failed to query size distribution: %w", err)
	}

	// 3. Обогащаем моками и вычисляемыми полями
	trends := []string{"stable", "growing", "declining"}
	var mostPopularSize, highestReturnRateSize string
	maxReturnRate := 0.0

	if len(distribution) > 0 {
		mostPopularSize = distribution[0].Size
	}

	for i := range distribution {
		d := &distribution[i]
		if totalSalesCount > 0 {
			d.Percentage = (float64(d.SalesCount) / float64(totalSalesCount)) * 100
		}
		d.Trend = trends[rand.Intn(len(trends))]
		d.ReturnRate = math.Round((rand.Float64()*15+5)*10) / 10 // 5-20%
		if d.ReturnRate > maxReturnRate {
			maxReturnRate = d.ReturnRate
			highestReturnRateSize = d.Size
		}
	}

	insights := model.SizeDistributionInsights{
		MostPopularSize:       mostPopularSize,
		FastestGrowingSize:    "L", // Mocked
		HighestReturnRateSize: highestReturnRateSize,
		Recommendations: []string{ // Mocked
			"Увеличить ассортимент размера " + mostPopularSize,
			"Пересмотреть размерную сетку для " + highestReturnRateSize,
			"Добавить больше моделей в размере L",
		},
	}

	return &model.SizeDistributionResponse{
		Period:           params.Period,
		Category:         params.Category,
		SizeDistribution: distribution,
		Insights:         insights,
	}, nil
}

// --- Seasonal Trends ---
// NOTE: This is fully mocked as the DB schema doesn't support this level of historical analysis.
func (r *analyticsRepository) GetSeasonalTrends(ctx context.Context, userID uuid.UUID, params model.SeasonalTrendsRequestParams) (*model.SeasonalTrendsResponse, error) {
	year := params.Year
	if year == 0 {
		year = time.Now().Year()
	}

	// This is a completely mocked response
	response := &model.SeasonalTrendsResponse{
		Year: year,
		SeasonalData: []model.SeasonalData{
			{
				Season: "winter",
				Months: []string{"december", "january", "february"},
				Categories: []model.SeasonalCategoryData{
					{Category: "coats", Name: "Пальто", DemandIndex: 95, Revenue: 1250000, GrowthVsPrevYear: 12.5, PeakMonth: "january", RecommendedStockLevel: "high"},
					{Category: "boots", Name: "Сапоги", DemandIndex: 88, Revenue: 980000, GrowthVsPrevYear: 8.3, PeakMonth: "december", RecommendedStockLevel: "high"},
				},
			},
			{
				Season: "summer",
				Months: []string{"june", "july", "august"},
				Categories: []model.SeasonalCategoryData{
					{Category: "dresses", Name: "Платья", DemandIndex: 92, Revenue: 1450000, GrowthVsPrevYear: 18.7, PeakMonth: "july", RecommendedStockLevel: "very_high"},
				},
			},
		},
		Predictions: model.Predictions{
			NextSeason: "spring",
			TrendingCategories: []model.TrendPrediction{
				{Category: "sustainable_fashion", Name: "Эко-одежда", PredictedGrowth: 35.2, Confidence: 87},
				{Category: "oversized", Name: "Оверсайз", PredictedGrowth: 28.1, Confidence: 92},
			},
			DecliningCategories: []model.TrendPrediction{
				{Category: "formal_wear", Name: "Деловая одежда", PredictedDecline: -12.5, Confidence: 78},
			},
		},
	}
	return response, nil
}

// --- Returns Analytics ---
// NOTE: This is fully mocked as the DB schema does not contain a returns table.
func (r *analyticsRepository) GetReturnsAnalytics(ctx context.Context, userID uuid.UUID, params model.ReturnsAnalyticsRequestParams) (*model.ReturnsAnalyticsResponse, error) {
	period := params.Period
	if period == "" {
		period = "30d"
	}

	// This is a completely mocked response
	response := &model.ReturnsAnalyticsResponse{
		Period: period,
		Summary: model.ReturnSummary{
			TotalReturns:      45,
			TotalReturnValue:  405000,
			ReturnRate:        12.5,
			AvgReturnValue:    9000,
			ReturnRateChange:  -2.3,
			ProcessingTimeAvg: 3.2,
		},
		ReturnReasons: []model.ReturnReason{
			{Reason: "size_mismatch", Name: "Не подошел размер", Count: 18, Percentage: 40.0, Value: 162000, AvgProcessingDays: 2.1, Trend: "stable"},
			{Reason: "quality_issues", Name: "Проблемы с качеством", Count: 12, Percentage: 26.7, Value: 108000, AvgProcessingDays: 4.5, Trend: "increasing"},
		},
		ByCategory: []model.ReturnByCategory{
			{Category: "shoes", Name: "Обувь", ReturnCount: 15, ReturnRate: 22.1, MainReason: "size_mismatch", AvgReturnValue: 11200},
			{Category: "dresses", Name: "Платья", ReturnCount: 12, ReturnRate: 15.7, MainReason: "size_mismatch", AvgReturnValue: 8500},
		},
		ByProduct: []model.ReturnByProduct{
			{ProductID: uuid.NewString(), Name: "Сапоги на высоком каблуке", ReturnCount: 6, ReturnRate: 35.3, MainReason: "size_mismatch", Recommendations: []string{"Обновить размерную сетку", "Добавить детальные замеры"}},
		},
		Recommendations: []string{
			"Улучшить описание размеров для обуви",
			"Провести контроль качества поставщика XYZ",
		},
	}
	return response, nil
}