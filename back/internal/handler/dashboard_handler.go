// internal/handler/dashboard_handler.go
package handler

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lamoda-seller-app/internal/middleware"
	"github.com/lamoda-seller-app/internal/model"
	"github.com/lamoda-seller-app/internal/repository"
)

const (
	defaultTopCategoriesLimit = 5
	dateFormat                = "2006-01-02"
	dateTimeFormat            = "2006-01-02 15:00"
)

// DashboardHandler обрабатывает HTTP-запросы для дашборда.
type DashboardHandler struct {
	repo repository.DashboardRepositoryInterface
}

// NewDashboardHandler создает новый экземпляр DashboardHandler.
func NewDashboardHandler(repo repository.DashboardRepositoryInterface) *DashboardHandler {
	return &DashboardHandler{repo: repo}
}

// ErrorResponse - стандартная структура для ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// --- Обработчики API ---

// GetStats обрабатывает GET /api/dashboard/stats
func (h *DashboardHandler) GetStats(c *gin.Context) {
	log.Printf("📊 Dashboard GetStats: начало обработки запроса")

	var params model.StatsRequestParams
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Printf("❌ Dashboard GetStats: ошибка парсинга параметров: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid query parameters: " + err.Error()})
		return
	}

	log.Printf("📋 Dashboard GetStats: параметры запроса: %+v", params)

	// Получаем ID пользователя из контекста
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	log.Printf("👤 Dashboard GetStats: пользователь ID: %s", userID)

	ctx := c.Request.Context()

	// 1. Рассчитываем временные периоды
	log.Printf("⏰ Dashboard GetStats: расчет временных периодов")
	periodInfo, err := calculatePeriods(params)
	if err != nil {
		log.Printf("❌ Dashboard GetStats: ошибка расчета периодов: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	log.Printf("📅 Dashboard GetStats: период с %s по %s", periodInfo.DateFrom, periodInfo.DateTo)

	// 2. Получаем данные из репозитория - теперь с userID
	log.Printf("🔍 Dashboard GetStats: запрос текущих данных из репозитория")
	currentData, err := h.repo.GetAggregatedData(ctx, userID, periodInfo.DateFrom, periodInfo.DateTo)
	if err != nil {
		log.Printf("❌ Dashboard GetStats: ошибка получения текущих данных: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve dashboard statistics"})
		return
	}
	log.Printf("📈 Dashboard GetStats: текущие данные - выручка: %.2f, заказов: %d, товаров: %d, возвратов: %d",
		currentData.Revenue, currentData.OrdersCount, currentData.ItemsSoldCount, currentData.ReturnsCount)

	var previousData repository.AggregatedData
	if params.CompareWithPrevious && periodInfo.PreviousPeriod != nil {
		log.Printf("🔍 Dashboard GetStats: запрос данных предыдущего периода")
		previousData, err = h.repo.GetAggregatedData(ctx, userID, periodInfo.PreviousPeriod.DateFrom, periodInfo.PreviousPeriod.DateTo)
		if err != nil {
			log.Printf("❌ Dashboard GetStats: ошибка получения данных предыдущего периода: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve dashboard statistics"})
			return
		}
		log.Printf("📉 Dashboard GetStats: данные предыдущего периода - выручка: %.2f, заказов: %d",
			previousData.Revenue, previousData.OrdersCount)
	}

	log.Printf("🏷️ Dashboard GetStats: запрос топ категорий")
	topCategories, err := h.repo.GetTopCategories(ctx, userID, periodInfo.DateFrom, periodInfo.DateTo, defaultTopCategoriesLimit)
	if err != nil {
		log.Printf("❌ Dashboard GetStats: ошибка получения топ категорий: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve dashboard statistics"})
		return
	}
	log.Printf("🏆 Dashboard GetStats: найдено %d топ категорий", len(topCategories))

	var hourlySalesRaw []model.HourlySale
	// ИСПРАВЛЕНИЕ: Запрашиваем почасовую статистику только для однодневных периодов
	if periodInfo.DateTo.Sub(periodInfo.DateFrom) < 25*time.Hour {
		log.Printf("⏱️ Dashboard GetStats: запрос почасовых продаж (период < 25 часов)")
		hourlySalesRaw, err = h.repo.GetHourlySales(ctx, userID, periodInfo.DateFrom, periodInfo.DateTo)
		if err != nil {
			log.Printf("❌ Dashboard GetStats: ошибка получения почасовых продаж: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve dashboard statistics"})
			return
		}
		log.Printf("📊 Dashboard GetStats: получено %d записей почасовых продаж", len(hourlySalesRaw))
	} else {
		log.Printf("⏭️ Dashboard GetStats: пропуск почасовых продаж (период > 25 часов)")
	}

	// 3. Собираем финальный ответ
	log.Printf("🔧 Dashboard GetStats: формирование ответа")
	response := &model.StatsResponse{
		Period:         *periodInfo,
		Revenue:        calculateMetric(currentData.Revenue, previousData.Revenue),
		Orders:         calculateMetric(float64(currentData.OrdersCount), float64(previousData.OrdersCount)),
		ItemsSold:      calculateMetric(float64(currentData.ItemsSoldCount), float64(previousData.ItemsSoldCount)),
		AvgOrderValue:  calculateMetric(safeDivide(currentData.Revenue, float64(currentData.OrdersCount)), safeDivide(previousData.Revenue, float64(previousData.OrdersCount))),
		ConversionRate: model.Metric{}, // ConversionRate требует данных о сессиях, которых нет. Пока заглушка.
		ReturnRate:     calculateMetric(safeDivide(float64(currentData.ReturnsCount)*100, float64(currentData.OrdersCount)), safeDivide(float64(previousData.ReturnsCount)*100, float64(previousData.OrdersCount))),
		TopCategories:  topCategories,
		// Заполняем пропущенные часы только если данные были запрошены
		HourlySales: fillMissingHours(hourlySalesRaw),
	}

	log.Printf("✅ Dashboard GetStats: успешно сформирован ответ")
	c.JSON(http.StatusOK, response)
}

// GetSalesChart обрабатывает GET /api/dashboard/sales-chart
func (h *DashboardHandler) GetSalesChart(c *gin.Context) {
	log.Printf("📈 Dashboard GetSalesChart: начало обработки запроса")

	var params model.SalesChartRequestParams
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Printf("❌ Dashboard GetSalesChart: ошибка парсинга параметров: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid query parameters: " + err.Error()})
		return
	}

	log.Printf("📋 Dashboard GetSalesChart: параметры запроса: %+v", params)

	// Получаем ID пользователя из контекста
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	log.Printf("👤 Dashboard GetSalesChart: пользователь ID: %s", userID)

	ctx := c.Request.Context()

	// 1. Рассчитываем период и гранулярность
	log.Printf("⏰ Dashboard GetSalesChart: расчет периода графика")
	start, end, err := calculateChartPeriod(params.Period)
	if err != nil {
		log.Printf("❌ Dashboard GetSalesChart: ошибка расчета периода: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	granularity := determineGranularity(params, start, end)
	log.Printf("📅 Dashboard GetSalesChart: период с %s по %s, гранулярность: %s", start, end, granularity)

	// 2. Получаем данные из репозитория с userID
	log.Printf("🔍 Dashboard GetSalesChart: запрос данных графика из репозитория")
	dataPoints, err := h.repo.GetSalesChartData(ctx, userID, start, end, granularity)
	if err != nil {
		log.Printf("❌ Dashboard GetSalesChart: ошибка получения данных графика: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve sales chart data"})
		return
	}
	log.Printf("📊 Dashboard GetSalesChart: получено %d точек данных", len(dataPoints))

	// 3. Считаем summary и форматируем данные для ответа
	log.Printf("🔧 Dashboard GetSalesChart: расчет сводки и форматирование данных")
	summary, formattedData := calculateChartSummaryAndFormatData(dataPoints, granularity)

	response := &model.SalesChartResponse{
		Period:      params.Period,
		Metric:      params.Metric,
		Granularity: granularity,
		Data:        formattedData,
		Summary:     summary,
	}

	log.Printf("✅ Dashboard GetSalesChart: успешно сформирован ответ")
	c.JSON(http.StatusOK, response)
}

// --- Вспомогательные функции (замена слоя Service) ---
// ... (остальные вспомогательные функции без изменений)
func calculatePeriods(params model.StatsRequestParams) (*model.PeriodInfo, error) {
	now := time.Now().UTC()
	var start, end time.Time

	periodType := params.Period
	if periodType == "" && (params.DateFrom == "" || params.DateTo == "") {
		periodType = "today" // По умолчанию
	}

	if periodType != "" {
		switch periodType {
		case "today":
			start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
			end = start.AddDate(0, 0, 1).Add(-time.Nanosecond)
		case "yesterday":
			yesterday := now.AddDate(0, 0, -1)
			start = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
			end = start.AddDate(0, 0, 1).Add(-time.Nanosecond)
		case "week":
			weekday := now.Weekday()
			if weekday == time.Sunday {
				weekday = 7
			}
			start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, int(time.Monday-weekday))
			end = start.AddDate(0, 0, 7).Add(-time.Nanosecond)
		case "month":
			start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
			end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
		case "quarter":
			quarter := (int(now.Month())-1)/3 + 1
			startMonth := time.Month((quarter-1)*3 + 1)
			start = time.Date(now.Year(), startMonth, 1, 0, 0, 0, 0, time.UTC)
			end = start.AddDate(0, 3, 0).Add(-time.Nanosecond)
		case "year":
			start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
			end = start.AddDate(1, 0, 0).Add(-time.Nanosecond)
		default:
			return nil, fmt.Errorf("invalid period: %s", params.Period)
		}
	} else { // Кастомный период
		var err error
		start, err = time.ParseInLocation(dateFormat, params.DateFrom, time.UTC)
		if err != nil {
			return nil, fmt.Errorf("invalid date_from format: %w", err)
		}
		end, err = time.ParseInLocation(dateFormat, params.DateTo, time.UTC)
		if err != nil {
			return nil, fmt.Errorf("invalid date_to format: %w", err)
		}
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, time.UTC)
	}

	info := &model.PeriodInfo{
		Type:     periodType,
		DateFrom: start,
		DateTo:   end,
	}

	if params.CompareWithPrevious {
		duration := end.Sub(start)
		prevStart := start.Add(-duration - time.Nanosecond)
		prevEnd := start.Add(-time.Nanosecond)
		info.PreviousPeriod = &model.PreviousPeriod{
			DateFrom: prevStart,
			DateTo:   prevEnd,
		}
	}
	return info, nil
}

func calculateMetric(current, previous float64) model.Metric {
	metric := model.Metric{Current: current, Previous: previous}
	if previous != 0 {
		metric.ChangeAbsolute = current - previous
		metric.ChangePercent = (metric.ChangeAbsolute / previous) * 100
	} else if current > 0 {
		metric.ChangeAbsolute = current
		metric.ChangePercent = 100.0
	}

	if metric.ChangeAbsolute > 0.001 {
		metric.Trend = "up"
	} else if metric.ChangeAbsolute < -0.001 {
		metric.Trend = "down"
	} else {
		metric.Trend = "stable"
	}
	return metric
}

func safeDivide(numerator, denominator float64) float64 {
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

func fillMissingHours(sales []model.HourlySale) []model.HourlySale {
	if sales == nil {
		return []model.HourlySale{}
	}
	salesMap := make(map[int]model.HourlySale, len(sales))
	for _, s := range sales {
		salesMap[s.Hour] = s
	}

	fullSales := make([]model.HourlySale, 24)
	for i := 0; i < 24; i++ {
		if sale, ok := salesMap[i]; ok {
			fullSales[i] = sale
		} else {
			fullSales[i] = model.HourlySale{Hour: i, Revenue: 0, Orders: 0}
		}
	}
	return fullSales
}

func calculateChartPeriod(period string) (time.Time, time.Time, error) {
	now := time.Now().UTC()
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, time.UTC)
	var start time.Time

	switch period {
	case "7d":
		start = end.AddDate(0, 0, -6)
	case "30d":
		start = end.AddDate(0, 0, -29)
	case "90d":
		start = end.AddDate(0, 0, -89)
	case "1y":
		start = end.AddDate(-1, 0, 0).AddDate(0, 0, 1)
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("invalid period for chart: %s", period)
	}
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	return start, end, nil
}

func determineGranularity(params model.SalesChartRequestParams, start, end time.Time) string {
	if params.Granularity != "" {
		return params.Granularity
	}
	days := end.Sub(start).Hours() / 24
	if days <= 2 {
		return "hour"
	}
	if days <= 90 {
		return "day"
	}
	if days <= 366*2 {
		return "week"
	}
	return "month"
}

func calculateChartSummaryAndFormatData(dataPoints []model.SalesChartDataPoint, granularity string) (model.SalesChartSummary, []model.SalesChartDataPoint) {
	summary := model.SalesChartSummary{
		WorstDay: model.BestWorstDay{Revenue: math.MaxFloat64},
	}
	numPoints := 0
	format := dateFormat
	if granularity == "hour" {
		format = dateTimeFormat
	}

	for i := range dataPoints {
		dataPoints[i].Date = dataPoints[i].Timestamp.UTC().Format(format)

		summary.TotalOrdersRevenue += dataPoints[i].OrdersRevenue
		summary.TotalPurchasesRevenue += dataPoints[i].PurchasesRevenue
		summary.TotalOrdersCount += dataPoints[i].OrdersCount
		summary.TotalPurchasesCount += dataPoints[i].PurchasesCount
		summary.TotalReturnsCount += dataPoints[i].ReturnCount
		summary.TotalReturnsRevenue += dataPoints[i].ReturnRevenue

		if dataPoints[i].OrdersRevenue > summary.BestDay.Revenue {
			summary.BestDay.Revenue = dataPoints[i].OrdersRevenue
			summary.BestDay.Date = dataPoints[i].Date
		}
		if dataPoints[i].OrdersRevenue < summary.WorstDay.Revenue {
			summary.WorstDay.Revenue = dataPoints[i].OrdersRevenue
			summary.WorstDay.Date = dataPoints[i].Date
		}
		numPoints++
	}

	if numPoints > 0 {
		summary.AvgDailyRevenue = summary.TotalOrdersRevenue / float64(numPoints)
	}
	if summary.WorstDay.Revenue == math.MaxFloat64 {
		summary.WorstDay = model.BestWorstDay{}
	}

	return summary, dataPoints
}
