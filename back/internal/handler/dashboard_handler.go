// internal/handler/dashboard_handler.go
package handler

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	var params model.StatsRequestParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid query parameters: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 1. Рассчитываем временные периоды
	periodInfo, err := calculatePeriods(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// 2. Получаем данные из репозитория
	currentData, err := h.repo.GetAggregatedData(ctx, periodInfo.DateFrom, periodInfo.DateTo)
	if err != nil {
		log.Printf("ERROR: failed to get current aggregated data: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve dashboard statistics"})
		return
	}

	var previousData repository.AggregatedData
	if params.CompareWithPrevious && periodInfo.PreviousPeriod != nil {
		previousData, err = h.repo.GetAggregatedData(ctx, periodInfo.PreviousPeriod.DateFrom, periodInfo.PreviousPeriod.DateTo)
		if err != nil {
			log.Printf("ERROR: failed to get previous aggregated data: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve dashboard statistics"})
			return
		}
	}

	topCategories, err := h.repo.GetTopCategories(ctx, periodInfo.DateFrom, periodInfo.DateTo, defaultTopCategoriesLimit)
	if err != nil {
		log.Printf("ERROR: failed to get top categories: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve dashboard statistics"})
		return
	}

	hourlySalesRaw, err := h.repo.GetHourlySales(ctx, periodInfo.DateFrom, periodInfo.DateTo)
	if err != nil {
		log.Printf("ERROR: failed to get hourly sales: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve dashboard statistics"})
		return
	}

	// 3. Собираем финальный ответ
	response := &model.StatsResponse{
		Period:         *periodInfo,
		Revenue:        calculateMetric(currentData.Revenue, previousData.Revenue),
		Orders:         calculateMetric(float64(currentData.OrdersCount), float64(previousData.OrdersCount)),
		ItemsSold:      calculateMetric(float64(currentData.ItemsSoldCount), float64(previousData.ItemsSoldCount)),
		AvgOrderValue:  calculateMetric(safeDivide(currentData.Revenue, float64(currentData.OrdersCount)), safeDivide(previousData.Revenue, float64(previousData.OrdersCount))),
		ConversionRate: model.Metric{}, // ConversionRate требует данных о сессиях, которых нет. Пока заглушка.
		ReturnRate:     calculateMetric(safeDivide(float64(currentData.ReturnsCount)*100, float64(currentData.OrdersCount)), safeDivide(float64(previousData.ReturnsCount)*100, float64(previousData.OrdersCount))),
		TopCategories:  topCategories,
		HourlySales:    fillMissingHours(hourlySalesRaw),
	}

	c.JSON(http.StatusOK, response)
}

// GetSalesChart обрабатывает GET /api/dashboard/sales-chart
func (h *DashboardHandler) GetSalesChart(c *gin.Context) {
	var params model.SalesChartRequestParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid query parameters: " + err.Error()})
		return
	}
	ctx := c.Request.Context()

	// 1. Рассчитываем период и гранулярность
	start, end, err := calculateChartPeriod(params.Period)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	granularity := determineGranularity(params, start, end)

	// 2. Получаем данные из репозитория
	dataPoints, err := h.repo.GetSalesChartData(ctx, start, end, granularity)
	if err != nil {
		log.Printf("ERROR: failed to get sales chart data: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve sales chart data"})
		return
	}

	// 3. Считаем summary и форматируем данные для ответа
	summary, formattedData := calculateChartSummaryAndFormatData(dataPoints, granularity)

	response := &model.SalesChartResponse{
		Period:      params.Period,
		Metric:      params.Metric,
		Granularity: granularity,
		Data:        formattedData,
		Summary:     summary,
	}

	c.JSON(http.StatusOK, response)
}

// --- Вспомогательные функции (замена слоя Service) ---

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
