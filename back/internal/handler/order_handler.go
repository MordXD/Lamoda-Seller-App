package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lamoda-seller-app/internal/middleware"
	"github.com/lamoda-seller-app/internal/model"
	"github.com/lamoda-seller-app/internal/repository"
	"gorm.io/gorm"
)

// OrderHandler обрабатывает HTTP запросы, связанные с заказами.
type OrderHandler struct {
	repo *repository.OrderRepository
}

func NewOrderHandler(repo *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

// --- Структуры ответов API ---

type ListOrdersAPIResponse struct {
	Orders     []model.Order                 `json:"orders"`
	Summary    *repository.ListOrdersSummary `json:"summary"`
	Pagination PaginationResponse            `json:"pagination"`
}

// --- Обработчики ---

// ListOrders GET /api/orders
func (h *OrderHandler) ListOrders(c *gin.Context) {
	log.Printf("📦 Orders ListOrders: начало обработки запроса")

	// Получаем ID пользователя из контекста (добавленного middleware'ом)
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	log.Printf("👤 Orders ListOrders: пользователь ID: %s", userID)

	// --- Парсинг параметров ---
	log.Printf("📋 Orders ListOrders: парсинг параметров запроса")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	minAmount, _ := strconv.ParseFloat(c.DefaultQuery("min_amount", "0"), 64)
	maxAmount, _ := strconv.ParseFloat(c.DefaultQuery("max_amount", "0"), 64)

	var dateFrom, dateTo *time.Time
	if df := c.Query("date_from"); df != "" {
		if t, err := time.Parse(time.RFC3339, df); err == nil {
			dateFrom = &t
		}
	}
	if dt := c.Query("date_to"); dt != "" {
		if t, err := time.Parse(time.RFC3339, dt); err == nil {
			dateTo = &t
		}
	}

	customerID, _ := uuid.Parse(c.Query("customer_id"))
	productID, _ := uuid.Parse(c.Query("product_id"))

	params := repository.ListOrdersParams{
		UserID:     userID,
		Status:     c.Query("status"),
		DateFrom:   dateFrom,
		DateTo:     dateTo,
		CustomerID: customerID,
		ProductID:  productID,
		MinAmount:  minAmount,
		MaxAmount:  maxAmount,
		SortBy:     c.Query("sort_by"),
		SortOrder:  c.DefaultQuery("sort_order", "asc"),
		Limit:      limit,
		Offset:     offset,
	}

	log.Printf("🔍 Orders ListOrders: параметры поиска - статус: %s, лимит: %d, смещение: %d",
		params.Status, params.Limit, params.Offset)

	// --- Получение данных из репозитория ---
	log.Printf("🔍 Orders ListOrders: запрос данных из репозитория")
	orders, summary, total, err := h.repo.List(c.Request.Context(), params)
	if err != nil {
		log.Printf("❌ Orders ListOrders: ошибка получения заказов: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve orders: " + err.Error()})
		return
	}

	log.Printf("📊 Orders ListOrders: найдено %d заказов из %d общих", len(orders), total)
	if summary != nil {
		log.Printf("📈 Orders ListOrders: сводка - общая сумма: %.2f, средний чек: %.2f",
			summary.TotalAmount, summary.AvgOrderValue)
	}

	// --- Формирование ответа ---
	response := ListOrdersAPIResponse{
		Orders:  orders,
		Summary: summary,
		Pagination: PaginationResponse{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasNext: total > int64(limit+offset),
			HasPrev: offset > 0,
		},
	}

	log.Printf("✅ Orders ListOrders: успешно сформирован ответ")
	c.JSON(http.StatusOK, response)
}

// GetOrderByID GET /api/orders/{order_id}
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	log.Printf("📦 Orders GetOrderByID: начало обработки запроса")

	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	log.Printf("👤 Orders GetOrderByID: пользователь ID: %s", userID)

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		log.Printf("❌ Orders GetOrderByID: неверный формат ID заказа: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID format"})
		return
	}

	log.Printf("🔍 Orders GetOrderByID: поиск заказа ID: %s", orderID)

	order, err := h.repo.GetByID(c.Request.Context(), orderID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("❌ Orders GetOrderByID: заказ не найден или нет доступа")
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found or you don't have permission to view it"})
			return
		}
		log.Printf("❌ Orders GetOrderByID: ошибка базы данных: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	log.Printf("✅ Orders GetOrderByID: заказ найден - номер: %s, статус: %s", order.OrderNumber, order.Status)
	c.JSON(http.StatusOK, order)
}

// UpdateOrderStatus PUT /api/orders/{order_id}/status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	log.Printf("📦 Orders UpdateOrderStatus: начало обработки запроса")

	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	log.Printf("👤 Orders UpdateOrderStatus: пользователь ID: %s", userID)

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		log.Printf("❌ Orders UpdateOrderStatus: неверный формат ID заказа: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID format"})
		return
	}

	log.Printf("🔍 Orders UpdateOrderStatus: обновление заказа ID: %s", orderID)

	var req model.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ Orders UpdateOrderStatus: ошибка парсинга JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	log.Printf("📝 Orders UpdateOrderStatus: новый статус: %s, комментарий: %s", req.Status, req.Comment)

	// TODO: Добавить валидацию допустимых переходов статусов (например, из 'new' можно в 'confirmed', но не в 'delivered')

	updatedOrder, err := h.repo.UpdateStatus(c.Request.Context(), orderID, userID, req.Status, req.Comment, req.EstimatedDeliveryDate)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("❌ Orders UpdateOrderStatus: заказ не найден или нет доступа")
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found or you don't have permission to update it"})
			return
		}
		log.Printf("❌ Orders UpdateOrderStatus: ошибка обновления статуса: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status: " + err.Error()})
		return
	}

	log.Printf("✅ Orders UpdateOrderStatus: статус успешно обновлен на: %s", updatedOrder.Status)
	c.JSON(http.StatusOK, gin.H{
		"message": "Статус заказа успешно обновлен",
		"order": gin.H{
			"id":           updatedOrder.ID,
			"status":       updatedOrder.Status,
			"updated_date": updatedOrder.UpdatedAt,
		},
	})
}
