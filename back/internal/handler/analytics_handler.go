// internal/handler/analytics_handler.go
package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lamoda-seller-app/internal/middleware"
	"github.com/lamoda-seller-app/internal/model"
	"github.com/lamoda-seller-app/internal/repository"
)

// AnalyticsHandler обрабатывает HTTP-запросы для аналитики.
type AnalyticsHandler struct {
	repo repository.AnalyticsRepo
}

// NewAnalyticsHandler создает новый экземпляр AnalyticsHandler.
func NewAnalyticsHandler(repo repository.AnalyticsRepo) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo}
}

// getUserIDFromContext извлекает ID пользователя из контекста Gin.
func getUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userIDUntyped, exists := c.Get(middleware.UserIDKey)
	if !exists {
		log.Printf("Error: userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return uuid.Nil, false
	}
	userID, ok := userIDUntyped.(uuid.UUID)
	if !ok {
		log.Printf("Error: userID in context is not of type uuid.UUID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format in token"})
		return uuid.Nil, false
	}
	return userID, true
}

func (h *AnalyticsHandler) GetTopProducts(c *gin.Context) {
	var params model.TopProductsRequestParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	response, err := h.repo.GetTopProducts(c.Request.Context(), userID, params)
	if err != nil {
		log.Printf("Error in GetTopProducts repo call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top products"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AnalyticsHandler) GetCategoryAnalytics(c *gin.Context) {
	var params model.CategoryAnalyticsRequestParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	response, err := h.repo.GetCategoryAnalytics(c.Request.Context(), userID, params)
	if err != nil {
		log.Printf("Error in GetCategoryAnalytics repo call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get category analytics"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AnalyticsHandler) GetSizeDistribution(c *gin.Context) {
	var params model.SizeDistributionRequestParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	response, err := h.repo.GetSizeDistribution(c.Request.Context(), userID, params)
	if err != nil {
		log.Printf("Error in GetSizeDistribution repo call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get size distribution"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AnalyticsHandler) GetSeasonalTrends(c *gin.Context) {
	var params model.SeasonalTrendsRequestParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	response, err := h.repo.GetSeasonalTrends(c.Request.Context(), userID, params)
	if err != nil {
		log.Printf("Error in GetSeasonalTrends repo call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get seasonal trends"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AnalyticsHandler) GetReturnsAnalytics(c *gin.Context) {
	var params model.ReturnsAnalyticsRequestParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	response, err := h.repo.GetReturnsAnalytics(c.Request.Context(), userID, params)
	if err != nil {
		log.Printf("Error in GetReturnsAnalytics repo call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get returns analytics"})
		return
	}

	c.JSON(http.StatusOK, response)
}
