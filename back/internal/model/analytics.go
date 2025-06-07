// internal/model/analytics.go
package model

import "github.com/google/uuid"

// --- Общие параметры запросов для аналитики ---

type AnalyticsRequestParams struct {
	Period string `form:"period"` // 7d, 30d, 90d, 1y, all_time
	Metric string `form:"metric"` // revenue, quantity, profit, etc.
}

// --- Структуры для /api/analytics/top-products ---

type TopProductsRequestParams struct {
	AnalyticsRequestParams
	Category string `form:"category"`
	Limit    int    `form:"limit"`
}

type TopProduct struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Brand         string    `json:"brand"`
	Category      string    `json:"category"`
	SKU           string    `json:"sku"`
	Image         string    `json:"image"`
	Price         float64   `json:"price"`
	CostPrice     float64   `json:"cost_price"`
	SalesCount    int64     `json:"sales_count"`
	Revenue       float64   `json:"revenue"`
	Profit        float64   `json:"profit"`
	MarginPercent float64   `json:"margin_percent"`
	ReturnCount   int64     `json:"return_count"`
	ReturnRate    float64   `json:"return_rate"`
	Rating        float64   `json:"rating"`
	ReviewsCount  int64     `json:"reviews_count"`
	Rank          int       `json:"rank"`
	RankChange    int       `json:"rank_change"`    // 0 = stable, >0 = down, <0 = up
	RevenueShare  float64   `json:"revenue_share"`
}

type TopProductsSummary struct {
	TotalRevenue float64            `json:"total_revenue"`
	TotalSales   int64              `json:"total_sales"`
	TotalProfit  float64            `json:"total_profit"`
	AvgMargin    float64            `json:"avg_margin"`
	TopCategory  TopCategorySummary `json:"top_category"`
}

type TopCategorySummary struct {
	Category string  `json:"category"`
	Name     string  `json:"name"`
	Revenue  float64 `json:"revenue"`
	Share    float64 `json:"share"`
}

type TopProductsResponse struct {
	Period   string             `json:"period"`
	Metric   string             `json:"metric"`
	Products []TopProduct       `json:"products"`
	Summary  TopProductsSummary `json:"summary"`
}

// --- Структуры для /api/analytics/categories ---

type CategoryAnalyticsRequestParams struct {
	Period  string `form:"period"`   // 7d, 30d, 90d, 1y
	SortBy  string `form:"sort_by"`  // revenue, quantity, profit, margin
}

type CategoryBestseller struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Sales int64     `json:"sales"`
}

type CategoryAnalytics struct {
	ID             string             `json:"id"` // slug
	Name           string             `json:"name"`
	Revenue        float64            `json:"revenue"`
	Quantity       int64              `json:"quantity"`
	Profit         float64            `json:"profit"`
	MarginPercent  float64            `json:"margin_percent"`
	OrdersCount    int64              `json:"orders_count"`
	AvgPrice       float64            `json:"avg_price"`
	ReturnRate     float64            `json:"return_rate"`
	GrowthPercent  float64            `json:"growth_percent"`  // Mocked
	RevenueShare   float64            `json:"revenue_share"`
	TopSize        string             `json:"top_size"`         // Mocked
	TopColor       string             `json:"top_color"`        // Mocked
	SeasonalTrend  string             `json:"seasonal_trend"`   // Mocked: stable, growing, declining
	Bestseller     CategoryBestseller `json:"bestseller"`       // Mocked
}

type CategoryAnalyticsSummary struct {
	TotalCategories   int                      `json:"total_categories"`
	FastestGrowing    CategoryGrowthSummary    `json:"fastest_growing"`    // Mocked
	HighestMargin     CategoryMarginSummary    `json:"highest_margin"`     // Mocked
	HighestReturnRate CategoryReturnRateSummary `json:"highest_return_rate"` // Mocked
}

type CategoryGrowthSummary struct {
	Category string  `json:"category"`
	Growth   float64 `json:"growth"`
}
type CategoryMarginSummary struct {
	Category string  `json:"category"`
	Margin   float64 `json:"margin"`
}
type CategoryReturnRateSummary struct {
	Category   string  `json:"category"`
	ReturnRate float64 `json:"return_rate"`
}

type CategoryAnalyticsResponse struct {
	Period     string                   `json:"period"`
	Categories []CategoryAnalytics      `json:"categories"`
	Summary    CategoryAnalyticsSummary `json:"summary"`
}

// --- Структуры для /api/analytics/size-distribution ---

type SizeDistributionRequestParams struct {
	Period    string `form:"period"` // 7d, 30d, 90d, 1y
	Category  string `form:"category"`
	ProductID string `form:"product_id"`
}

type SizeDistribution struct {
	Size         string  `json:"size"`
	SalesCount   int64   `json:"sales_count"`
	Revenue      float64 `json:"revenue"`
	Percentage   float64 `json:"percentage"`
	ReturnRate   float64 `json:"return_rate"`
	AvgPrice     float64 `json:"avg_price"`
	Trend        string  `json:"trend"` // Mocked: stable, growing, declining
}

type SizeDistributionInsights struct {
	MostPopularSize        string   `json:"most_popular_size"`
	FastestGrowingSize     string   `json:"fastest_growing_size"`      // Mocked
	HighestReturnRateSize  string   `json:"highest_return_rate_size"`
	Recommendations        []string `json:"recommendations"`         // Mocked
}

type SizeDistributionResponse struct {
	Period            string               `json:"period"`
	Category          string               `json:"category"`
	SizeDistribution  []SizeDistribution   `json:"size_distribution"`
	Insights          SizeDistributionInsights `json:"insights"`
}

// --- Структуры для /api/analytics/seasonal-trends ---

type SeasonalTrendsRequestParams struct {
	Year     int    `form:"year"`
	Category string `form:"category"`
}

type SeasonalCategoryData struct {
	Category              string  `json:"category"`
	Name                  string  `json:"name"`
	DemandIndex           int     `json:"demand_index"` // Mocked
	Revenue               float64 `json:"revenue"`
	GrowthVsPrevYear      float64 `json:"growth_vs_prev_year"` // Mocked
	PeakMonth             string  `json:"peak_month"`
	RecommendedStockLevel string  `json:"recommended_stock_level"` // Mocked
}

type SeasonalData struct {
	Season     string                 `json:"season"`
	Months     []string               `json:"months"`
	Categories []SeasonalCategoryData `json:"categories"`
}

type TrendPrediction struct {
	Category        string  `json:"category"`
	Name            string  `json:"name"`
	PredictedGrowth float64 `json:"predicted_growth,omitempty"`
	PredictedDecline float64 `json:"predicted_decline,omitempty"`
	Confidence      int     `json:"confidence"`
}

type Predictions struct {
	NextSeason           string            `json:"next_season"`
	TrendingCategories   []TrendPrediction `json:"trending_categories"`
	DecliningCategories  []TrendPrediction `json:"declining_categories"`
}

type SeasonalTrendsResponse struct {
	Year          int            `json:"year"`
	SeasonalData  []SeasonalData `json:"seasonal_data"`
	Predictions   Predictions    `json:"predictions"` // Mocked
}


// --- Структуры для /api/analytics/returns ---

type ReturnsAnalyticsRequestParams struct {
	Period   string `form:"period"`
	Category string `form:"category"`
	Reason   string `form:"reason"`
}

type ReturnSummary struct {
	TotalReturns        int64   `json:"total_returns"`
	TotalReturnValue    float64 `json:"total_return_value"`
	ReturnRate          float64 `json:"return_rate"`
	AvgReturnValue      float64 `json:"avg_return_value"`
	ReturnRateChange    float64 `json:"return_rate_change"`   // Mocked
	ProcessingTimeAvg   float64 `json:"processing_time_avg"`  // Mocked
}

type ReturnReason struct {
	Reason            string  `json:"reason"`
	Name              string  `json:"name"`
	Count             int64   `json:"count"`
	Percentage        float64 `json:"percentage"`
	Value             float64 `json:"value"`
	AvgProcessingDays float64 `json:"avg_processing_days"` // Mocked
	Trend             string  `json:"trend"`               // Mocked
}

type ReturnByCategory struct {
	Category        string  `json:"category"`
	Name            string  `json:"name"`
	ReturnCount     int64   `json:"return_count"`
	ReturnRate      float64 `json:"return_rate"`
	MainReason      string  `json:"main_reason"`
	AvgReturnValue  float64 `json:"avg_return_value"`
}

type ReturnByProduct struct {
	ProductID       string   `json:"product_id"`
	Name            string   `json:"name"`
	ReturnCount     int64    `json:"return_count"`
	ReturnRate      float64  `json:"return_rate"`
	MainReason      string   `json:"main_reason"`
	Recommendations []string `json:"recommendations"` // Mocked
}

type ReturnsAnalyticsResponse struct {
	Period         string             `json:"period"`
	Summary        ReturnSummary      `json:"summary"`
	ReturnReasons  []ReturnReason     `json:"return_reasons"`
	ByCategory     []ReturnByCategory `json:"by_category"`
	ByProduct      []ReturnByProduct  `json:"by_product"`
	Recommendations []string           `json:"recommendations"` // Mocked
}