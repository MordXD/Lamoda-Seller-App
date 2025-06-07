// internal/model/dashboard.go
package model

import "time"

// --- Структуры для /api/dashboard/stats ---

// StatsRequestParams содержит параметры запроса для эндпоинта статистики.
type StatsRequestParams struct {
	Period              string `form:"period"` // today, yesterday, week, month, quarter, year
	DateFrom            string `form:"date_from"`
	DateTo              string `form:"date_to"`
	CompareWithPrevious bool   `form:"compare_with_previous"`
}

// PeriodInfo описывает текущий и предыдущий временные интервалы.
type PeriodInfo struct {
	Type           string          `json:"type"`
	DateFrom       time.Time       `json:"date_from"`
	DateTo         time.Time       `json:"date_to"`
	PreviousPeriod *PreviousPeriod `json:"previous_period,omitempty"`
}

type PreviousPeriod struct {
	DateFrom time.Time `json:"date_from"`
	DateTo   time.Time `json:"date_to"`
}

// Metric представляет собой числовой показатель с его сравнением.
type Metric struct {
	Current        float64 `json:"current"`
	Previous       float64 `json:"previous"`
	ChangePercent  float64 `json:"change_percent"`
	ChangeAbsolute float64 `json:"change_absolute"`
	Trend          string  `json:"trend"` // up, down, stable
}

// TopCategory описывает статистику по одной категории товаров.
type TopCategory struct {
	Category string  `json:"category"`
	Name     string  `json:"name"`
	Revenue  float64 `json:"revenue"`
	Orders   int64   `json:"orders"`
	Items    int64   `json:"items"`
}

// HourlySale описывает продажи в конкретный час.
type HourlySale struct {
	Hour    int     `json:"hour"`
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
}

// StatsResponse — это полная структура ответа для эндпоинта статистики.
type StatsResponse struct {
	Period         PeriodInfo    `json:"period"`
	Revenue        Metric        `json:"revenue"`
	Orders         Metric        `json:"orders"`
	ItemsSold      Metric        `json:"items_sold"`
	AvgOrderValue  Metric        `json:"avg_order_value"`
	ConversionRate Metric        `json:"conversion_rate"`
	ReturnRate     Metric        `json:"return_rate"`
	TopCategories  []TopCategory `json:"top_categories"`
	HourlySales    []HourlySale  `json:"hourly_sales"`
}

// --- Структуры для /api/dashboard/sales-chart ---

// SalesChartRequestParams содержит параметры запроса для графика продаж.
type SalesChartRequestParams struct {
	Period      string `form:"period" binding:"required"` // 7d, 30d, 90d, 1y
	Metric      string `form:"metric"`                    // revenue, orders, items
	Granularity string `form:"granularity"`               // hour, day, week, month
}

// SalesChartDataPoint представляет одну точку на графике.
type SalesChartDataPoint struct {
	Date             string    `json:"date"` // "YYYY-MM-DD" или "YYYY-MM-DD HH:00"
	Timestamp        time.Time `json:"-"`    // ИСПРАВЛЕНИЕ: Исключаем из JSON, так как это служебное поле
	OrdersRevenue    float64   `json:"orders_revenue"`
	PurchasesRevenue float64   `json:"purchases_revenue"` // Выручка за вычетом возвратов
	OrdersCount      int64     `json:"orders_count"`
	PurchasesCount   int64     `json:"purchases_count"`
	ReturnCount      int64     `json:"return_count"`
	ReturnRevenue    float64   `json:"return_revenue"`
}

// BestWorstDay описывает лучший/худший день по выручке.
type BestWorstDay struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
}

// SalesChartSummary содержит итоговые данные по графику.
type SalesChartSummary struct {
	TotalOrdersRevenue    float64      `json:"total_orders_revenue"`
	TotalPurchasesRevenue float64      `json:"total_purchases_revenue"`
	TotalOrdersCount      int64        `json:"total_orders_count"`
	TotalPurchasesCount   int64        `json:"total_purchases_count"`
	TotalReturnsCount     int64        `json:"total_returns_count"`
	TotalReturnsRevenue   float64      `json:"total_returns_revenue"`
	AvgDailyRevenue       float64      `json:"avg_daily_revenue"`
	BestDay               BestWorstDay `json:"best_day"`
	WorstDay              BestWorstDay `json:"worst_day"`
}

// SalesChartResponse — это полная структура ответа для графика продаж.
type SalesChartResponse struct {
	Period      string                `json:"period"`
	Metric      string                `json:"metric"`
	Granularity string                `json:"granularity"`
	Data        []SalesChartDataPoint `json:"data"`
	Summary     SalesChartSummary     `json:"summary"`
}