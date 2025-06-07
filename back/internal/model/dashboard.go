package model

import (
	"time"
)

type Order struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Amount    float64   `gorm:"not null"`
	Status    string    `gorm:"type:varchar(20);not null"` // "ordered" или "delivered"
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type KPI struct {
	TotalAmount       float64 `json:"total_amount"`
	AmountDiffPercent float64 `json:"amount_diff_percent"`
	TotalOrders       int     `json:"total_orders"`
	OrdersDiffPercent float64 `json:"orders_diff_percent"`
}

type SalesChartPoint struct {
	Date   string  `json:"date"`   // Формат: "2006-01-02" (ISO 8601)
	Amount float64 `json:"amount"`
}

// Добавлены отдельные поля для обоих типов данных
type DashboardSalesChart struct {
	Ordered   []SalesChartPoint `json:"ordered"`
	Delivered []SalesChartPoint `json:"delivered"`
}

type DashboardResponse struct {
	KPI        KPI                `json:"kpi"`
	SalesChart DashboardSalesChart `json:"sales_chart"`
}