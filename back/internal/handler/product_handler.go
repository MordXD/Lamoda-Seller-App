package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lamoda-seller-app/internal/model"
	"github.com/lamoda-seller-app/internal/repository"
	"gorm.io/gorm"
)

// ProductHandler обрабатывает HTTP запросы, связанные с продуктами.
type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// --- Request/Response Structures ---

type CreateProductRequest struct {
	Name        string                       `json:"name" binding:"required"`
	Price       float64                      `json:"price" binding:"required,gt=0"`
	OldPrice    float64                      `json:"old_price,omitempty"`
	ImageURL    string                       `json:"image_url" binding:"omitempty,url"`
	ShortDesc   string                       `json:"short_desc"`
	FullDesc    string                       `json:"full_desc"`
	BrandID     int                          `json:"brand_id"`
	CategoryID  int                          `json:"category_id"`
	InStock     int                          `json:"in_stock" binding:"gte=0"`
	Tags        []string                     `json:"tags"`
	Variants    []CreateProductVariantRequest `json:"variants"`
}

type CreateProductVariantRequest struct {
	Color   string `json:"color"`
	Size    string `json:"size"`
	SKU     string `json:"sku" binding:"required"`
	InStock int    `json:"in_stock" binding:"gte=0"`
}

type ListProductsResponse struct {
	Products []model.Product `json:"products"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

// --- Handlers ---

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	product := model.Product{
		Name:       req.Name,
		Price:      req.Price,
		OldPrice:   req.OldPrice,
		ImageURL:   req.ImageURL,
		ShortDesc:  req.ShortDesc,
		FullDesc:   req.FullDesc,
		BrandID:    req.BrandID,
		CategoryID: req.CategoryID,
		InStock:    req.InStock,
		Tags:       req.Tags,
	}

	for _, v := range req.Variants {
		product.Variants = append(product.Variants, model.ProductVariant{
			Color:   v.Color,
			Size:    v.Size,
			SKU:     v.SKU,
			InStock: v.InStock,
		})
	}

	if err := h.repo.Create(c.Request.Context(), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	params := repository.GetAllProductsParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}

	products, total, err := h.repo.GetAll(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve products"})
		return
	}

	c.JSON(http.StatusOK, ListProductsResponse{
		Products: products,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	product, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	// Убедимся, что товар существует
	if _, err := h.repo.GetByID(c.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	
	product := model.Product{
		ID:         id, // Важно установить ID для обновления
		Name:       req.Name,
		Price:      req.Price,
		OldPrice:   req.OldPrice,
		ImageURL:   req.ImageURL,
		ShortDesc:  req.ShortDesc,
		FullDesc:   req.FullDesc,
		BrandID:    req.BrandID,
		CategoryID: req.CategoryID,
		InStock:    req.InStock,
		Tags:       req.Tags,
	}

	for _, v := range req.Variants {
		product.Variants = append(product.Variants, model.ProductVariant{
			Color:   v.Color,
			Size:    v.Size,
			SKU:     v.SKU,
			InStock: v.InStock,
		})
	}
	
	if err := h.repo.Update(c.Request.Context(), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product updated successfully"})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) GetSalesStats(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	period := c.DefaultQuery("period", "week")
	
	stats, err := h.repo.GetSalesStats(c.Request.Context(), id, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sales stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *ProductHandler) GetPriceHistory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	history, err := h.repo.GetPriceHistory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get price history"})
		return
	}

	c.JSON(http.StatusOK, history)
}