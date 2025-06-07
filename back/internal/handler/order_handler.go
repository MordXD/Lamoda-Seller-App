package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	Orders     []model.Order               `json:"orders"`
	Summary    *repository.ListOrdersSummary `json:"summary"`
	Pagination PaginationResponse          `json:"pagination"`
}

// --- Обработчики ---

// ListOrders GET /api/orders
func (h *OrderHandler) ListOrders(c *gin.Context) {
	// Получаем ID пользователя из контекста (добавленного middleware'ом)
	// userID, _ := c.Get("userID") // Замените на реальное получение ID
	userID := uuid.New() // ЗАГЛУШКА: Замените на реальное получение ID пользователя из токена

	// --- Парсинг параметров ---
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

	// --- Получение данных из репозитория ---
	orders, summary, total, err := h.repo.List(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve orders: " + err.Error()})
		return
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

	c.JSON(http.StatusOK, response)
}

// GetOrderByID GET /api/orders/{order_id}
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	// userID, _ := c.Get("userID") // Замените на реальное получение ID
	userID := uuid.New() // ЗАГЛУШКА: Замените на реальное получение ID пользователя из токена
	
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID format"})
		return
	}

	order, err := h.repo.GetByID(c.Request.Context(), orderID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found or you don't have permission to view it"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrderStatus PUT /api/orders/{order_id}/status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	// userID, _ := c.Get("userID") // Замените на реальное получение ID
	userID := uuid.New() // ЗАГЛУШКА: Замените на реальное получение ID пользователя из токена
	
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID format"})
		return
	}

	var req model.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	// TODO: Добавить валидацию допустимых переходов статусов (например, из 'new' можно в 'confirmed', но не в 'delivered')

	updatedOrder, err := h.repo.UpdateStatus(c.Request.Context(), orderID, userID, req.Status, req.Comment, req.EstimatedDeliveryDate)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found or you don't have permission to update it"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Статус заказа успешно обновлен",
		"order": gin.H{
			"id": updatedOrder.ID,
			"status": updatedOrder.Status,
			"updated_date": updatedOrder.UpdatedAt,
		},
	})
}