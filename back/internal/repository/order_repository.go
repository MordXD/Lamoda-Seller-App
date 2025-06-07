package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lamoda-seller-app/internal/model"
	"gorm.io/gorm"
)

// OrderRepository инкапсулирует логику работы с заказами в БД.
type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// ListOrdersParams определяет все параметры для получения списка заказов.
type ListOrdersParams struct {
	UserID     uuid.UUID
	Status     string
	DateFrom   *time.Time
	DateTo     *time.Time
	CustomerID uuid.UUID
	ProductID  uuid.UUID
	MinAmount  float64
	MaxAmount  float64
	SortBy     string
	SortOrder  string
	Limit      int
	Offset     int
}

// StatusBreakdown - часть сводной информации по статусам.
type StatusBreakdown struct {
	Status string  `json:"-"`
	Count  int64   `json:"count"`
	Amount float64 `json:"amount"`
}

// ListOrdersSummary содержит сводную информацию по заказам.
type ListOrdersSummary struct {
	TotalOrders      int64                       `json:"total_orders"`
	TotalAmount      float64                     `json:"total_amount"`
	AvgOrderValue    float64                     `json:"avg_order_value"`
	StatusBreakdown  map[string]StatusBreakdown  `json:"status_breakdown"`
}

// List возвращает отфильтрованный и отсортированный список заказов с пагинацией и сводкой.
func (r *OrderRepository) List(ctx context.Context, params ListOrdersParams) ([]model.Order, *ListOrdersSummary, int64, error) {
	var orders []model.Order
	var total int64
	var summary ListOrdersSummary

	// --- 1. Создаем базовый запрос с обязательным фильтром по ID продавца ---
	query := r.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", params.UserID)

	// --- 2. Применяем все опциональные фильтры ---
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.DateFrom != nil {
		query = query.Where("date >= ?", params.DateFrom)
	}
	if params.DateTo != nil {
		query = query.Where("date <= ?", params.DateTo)
	}
	if params.CustomerID != uuid.Nil {
		query = query.Where("customer_id = ?", params.CustomerID)
	}
	// Для фильтра по ID товара нужен подзапрос
	if params.ProductID != uuid.Nil {
		query = query.Where("id IN (SELECT order_id FROM order_items WHERE product_id = ?)", params.ProductID)
	}
	if params.MinAmount > 0 {
		query = query.Where("totals ->> 'total' >= ?", fmt.Sprintf("%f", params.MinAmount))
	}
	if params.MaxAmount > 0 {
		query = query.Where("totals ->> 'total' <= ?", fmt.Sprintf("%f", params.MaxAmount))
	}

	// --- 3. Вычисляем сводную информацию (summary) на основе отфильтрованного набора ---
	if err := r.calculateSummary(query, &summary); err != nil {
		return nil, nil, 0, fmt.Errorf("failed to calculate summary: %w", err)
	}

	// --- 4. Считаем общее количество для пагинации ---
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, 0, err
	}

	// --- 5. Применяем сортировку ---
	if params.SortBy != "" {
		sortColumn := ""
		switch params.SortBy {
		case "date":
			sortColumn = "date"
		case "amount":
			sortColumn = "totals -> 'total'" // Сортировка по JSON-полю
		case "status":
			sortColumn = "status"
		}
		if sortColumn != "" {
			orderDirection := "ASC"
			if strings.ToLower(params.SortOrder) == "desc" {
				orderDirection = "DESC"
			}
			query = query.Order(fmt.Sprintf("%s %s", sortColumn, orderDirection))
		}
	} else {
		query = query.Order("date DESC") // Сортировка по умолчанию
	}

	// --- 6. Применяем пагинацию и Eager Loading для связанных данных ---
	err := query.Preload("Items").Preload("StatusHistory").
		Offset(params.Offset).Limit(params.Limit).Find(&orders).Error

	return orders, &summary, total, err
}

// calculateSummary вычисляет сводку по заказам на основе предоставленного запроса.
func (r *OrderRepository) calculateSummary(query *gorm.DB, summary *ListOrdersSummary) error {
	var totalAmount float64
	var totalOrders int64

	// Клонируем запрос, чтобы не изменять основной
	summaryQuery := query.Session(&gorm.Session{})

	// Общее кол-во и сумма
	if err := summaryQuery.Count(&totalOrders).Error; err != nil {
		return err
	}
	// Сумма извлекается из JSONB поля
	if err := summaryQuery.Select("COALESCE(SUM((totals->>'total')::numeric), 0)").Row().Scan(&totalAmount); err != nil {
		return err
	}

	summary.TotalOrders = totalOrders
	summary.TotalAmount = totalAmount
	if totalOrders > 0 {
		summary.AvgOrderValue = totalAmount / float64(totalOrders)
	}

	// Разбивка по статусам
	var breakdown []StatusBreakdown
	err := summaryQuery.Select("status, count(*) as count, COALESCE(SUM((totals->>'total')::numeric), 0) as amount").
		Group("status").Find(&breakdown).Error
	if err != nil {
		return err
	}

	summary.StatusBreakdown = make(map[string]StatusBreakdown)
	for _, item := range breakdown {
		summary.StatusBreakdown[item.Status] = item
	}

	return nil
}

// GetByID возвращает заказ по его ID со всеми связанными данными.
func (r *OrderRepository) GetByID(ctx context.Context, orderID, userID uuid.UUID) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("StatusHistory").
		Where("user_id = ?", userID). // Проверка, что заказ принадлежит этому продавцу
		First(&order, orderID).Error
	return &order, err
}

// UpdateStatus обновляет статус заказа и добавляет запись в историю в рамках одной транзакции.
func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID, userID uuid.UUID, status, comment string, estimatedDate *time.Time) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Найти заказ, убедившись, что он принадлежит продавцу
		if err := tx.Where("user_id = ?", userID).First(&order, orderID).Error; err != nil {
			return err // Возвращает gorm.ErrRecordNotFound, если не найден
		}

		// 2. Обновить статус и дату доставки в самом заказе
		order.Status = status
		if estimatedDate != nil {
			order.Delivery.EstimatedDate = *estimatedDate
			// GORM автоматически обработает обновление JSONB поля
		}
		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		// 3. Создать новую запись в истории статусов
		historyEntry := model.StatusHistory{
			OrderID: order.ID,
			Status:  status,
			Date:    time.Now(),
			Comment: comment,
		}
		if err := tx.Create(&historyEntry).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &order, nil
}