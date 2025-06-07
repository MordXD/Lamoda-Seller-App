package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lamoda-seller-app/internal/model"
	"gorm.io/gorm"
)

// DashboardRepositoryInterface определяет контракт для работы с дашбордом
type DashboardRepositoryInterface interface {
	GetDashboardData(ctx context.Context) (*model.DashboardResponse, error)
}

var _ DashboardRepositoryInterface = (*DashboardRepository)(nil)

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetDashboardData(ctx context.Context) (*model.DashboardResponse, error) {
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayEnd := todayStart.AddDate(0, 0, 1).Add(-time.Nanosecond)
	yesterdayStart := todayStart.AddDate(0, 0, -1)
	yesterdayEnd := todayStart

	todayData, err := r.getKpiData(ctx, todayStart, todayEnd)
	if err != nil {
		return nil, fmt.Errorf("error getting today KPI: %w", err)
	}

	yesterdayData, err := r.getKpiData(ctx, yesterdayStart, yesterdayEnd)
	if err != nil {
		return nil, fmt.Errorf("error getting yesterday KPI: %w", err)
	}

	kpi := model.KPI{
		TotalAmount:       todayData.Amount,
		TotalOrders:       todayData.Orders,
		AmountDiffPercent: calculatePercentChange(todayData.Amount, yesterdayData.Amount),
		OrdersDiffPercent: calculatePercentChange(float64(todayData.Orders), float64(yesterdayData.Orders)),
	}

	salesChart, err := r.getSalesChartData(ctx, todayStart.AddDate(0, 0, -6), todayEnd)
	if err != nil {
		return nil, fmt.Errorf("error getting sales chart data: %w", err)
	}

	return &model.DashboardResponse{
		KPI:        kpi,
		SalesChart: *salesChart,
	}, nil
}

type kpiData struct {
	Amount float64
	Orders int
}

func (r *DashboardRepository) getKpiData(ctx context.Context, start, end time.Time) (kpiData, error) {
	var result kpiData
	err := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Select("COALESCE(SUM(amount), 0) as amount, COUNT(*) as orders").
		Where("status = ? AND created_at BETWEEN ? AND ?", "ordered", start, end).
		Scan(&result).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return kpiData{}, err
	}
	return result, nil
}

func calculatePercentChange(today, yesterday float64) float64 {
	if yesterday == 0 {
		if today == 0 {
			return 0
		}
		return 100
	}
	return ((today - yesterday) / yesterday) * 100
}

func (r *DashboardRepository) getSalesChartData(ctx context.Context, start, end time.Time) (*model.DashboardSalesChart, error) {
	orderedMap := make(map[string]float64)
	deliveredMap := make(map[string]float64)

	current := start
	for !current.After(end) {
		dateStr := current.Format("2006-01-02")
		orderedMap[dateStr] = 0
		deliveredMap[dateStr] = 0
		current = current.AddDate(0, 0, 1)
	}

	type chartData struct {
		Date   string
		Status string
		Amount float64
	}

	var results []chartData
	err := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Select(
			"DATE(created_at) as date",
			"status",
			"COALESCE(SUM(amount), 0) as amount",
		).
		Where("created_at BETWEEN ? AND ?", start, end).
		Group("date, status").
		Scan(&results).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	for _, res := range results {
		switch res.Status {
		case "ordered":
			orderedMap[res.Date] = res.Amount
		case "delivered":
			deliveredMap[res.Date] = res.Amount
		}
	}

	orderedPoints := make([]model.SalesChartPoint, 0)
	deliveredPoints := make([]model.SalesChartPoint, 0)

	current = start
	for !current.After(end) {
		dateStr := current.Format("2006-01-02")
		orderedPoints = append(orderedPoints, model.SalesChartPoint{
			Date:   dateStr,
			Amount: orderedMap[dateStr],
		})
		deliveredPoints = append(deliveredPoints, model.SalesChartPoint{
			Date:   dateStr,
			Amount: deliveredMap[dateStr],
		})
		current = current.AddDate(0, 0, 1)
	}

	return &model.DashboardSalesChart{
		Ordered:   orderedPoints,
		Delivered: deliveredPoints,
	}, nil
}