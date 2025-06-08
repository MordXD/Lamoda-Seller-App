package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lamoda-seller-app/internal/model"
	"github.com/lamoda-seller-app/internal/repository"
	"gorm.io/gorm"
)

// ProductHandler обрабатывает HTTP запросы, связанные с продуктами.
type ProductHandler struct {
	repo *repository.ProductRepository
	// В реальном приложении сюда бы добавился ImageService для загрузки файлов
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// --- Структуры ответов API ---

type PaginationResponse struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasNext bool  `json:"has_next"`
	HasPrev bool  `json:"has_prev"`
}

type ListProductsAPIResponse struct {
	Products   []model.Product          `json:"products"`
	Pagination PaginationResponse       `json:"pagination"`
	Filters    *repository.FilterValues `json:"filters"`
}

// --- Обработчики ---

// ListProducts GET /api/products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	log.Printf("🛍️ Products ListProducts: начало обработки запроса")

	// --- Парсинг параметров ---
	log.Printf("📋 Products ListProducts: парсинг параметров запроса")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	minPrice, _ := strconv.ParseFloat(c.DefaultQuery("min_price", "0"), 64)
	maxPrice, _ := strconv.ParseFloat(c.DefaultQuery("max_price", "0"), 64)

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	params := repository.ListProductsParams{
		Limit:       limit,
		Offset:      offset,
		Search:      c.Query("search"),
		Category:    c.Query("category"),
		Brand:       c.Query("brand"),
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		StockStatus: c.Query("stock_status"),
		SortBy:      c.Query("sort_by"),
		SortOrder:   c.DefaultQuery("sort_order", "asc"),
	}

	log.Printf("🔍 Products ListProducts: параметры поиска - поиск: '%s', категория: '%s', бренд: '%s', лимит: %d",
		params.Search, params.Category, params.Brand, params.Limit)

	// --- Получение данных ---
	log.Printf("🔍 Products ListProducts: запрос данных из репозитория")
	products, total, filters, err := h.repo.List(c.Request.Context(), params)
	if err != nil {
		log.Printf("❌ Products ListProducts: ошибка получения продуктов: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve products: " + err.Error()})
		return
	}

	log.Printf("📊 Products ListProducts: найдено %d продуктов из %d общих", len(products), total)

	// --- Формирование ответа ---
	response := ListProductsAPIResponse{
		Products: products,
		Pagination: PaginationResponse{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasNext: total > int64(limit+offset),
			HasPrev: offset > 0,
		},
		Filters: filters,
	}

	log.Printf("✅ Products ListProducts: успешно сформирован ответ")
	c.JSON(http.StatusOK, response)
}

// GetProductByID GET /api/products/{id}
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	log.Printf("🛍️ Products GetProductByID: начало обработки запроса")

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Printf("❌ Products GetProductByID: неверный формат ID продукта: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID format"})
		return
	}

	log.Printf("🔍 Products GetProductByID: поиск продукта ID: %s", id)

	product, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("❌ Products GetProductByID: продукт не найден")
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		log.Printf("❌ Products GetProductByID: ошибка базы данных: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	log.Printf("✅ Products GetProductByID: продукт найден - название: %s, бренд: %s", product.Name, product.Brand)
	c.JSON(http.StatusOK, product)
}

// CreateProduct POST /api/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req model.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	// req.ID будет нулевым, GORM сгенерирует новый UUID
	req.ID = uuid.New() // Явно генерируем, чтобы вернуть в ответе

	if err := h.repo.Create(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product: " + err.Error()})
		return
	}

	// Получаем созданный товар со всеми полями для ответа
	createdProduct, err := h.repo.GetByID(c.Request.Context(), req.ID)
	if err != nil {
		// Даже если товар создан, но не можем его получить, это ошибка
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve created product: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      createdProduct.ID,
		"message": "Товар успешно создан",
		"product": createdProduct,
	})
}

// UpdateProduct PUT /api/products/{id}
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID format"})
		return
	}

	// Сначала получаем существующий товар
	product, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	// Привязываем JSON к СУЩЕСТВУЮЩЕМУ объекту.
	// Это обновит только те поля, которые пришли в запросе.
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	// ID не должен меняться
	product.ID = id

	if err := h.repo.Update(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product: " + err.Error()})
		return
	}

	// Получаем обновленный продукт для ответа
	updatedProduct, _ := h.repo.GetByID(c.Request.Context(), id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Товар успешно обновлен",
		"product": updatedProduct,
	})
}

// --- НОВЫЙ ОБРАБОТЧИК ---
// DeleteProduct DELETE /api/products/{id}
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID format"})
		return
	}

	// Перед удалением можно проверить, существует ли товар, чтобы вернуть 404,
	// но DELETE к несуществующему ресурсу часто считают идемпотентной операцией.
	// Здесь мы просто пытаемся удалить.
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		// Ошибка может возникнуть, например, из-за проблем с подключением к БД.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product: " + err.Error()})
		return
	}

	// Успешный ответ на DELETE - это 204 No Content с пустым телом.
	// Можно также вернуть 200 OK с сообщением.
	c.JSON(http.StatusOK, gin.H{"message": "Товар успешно удален"})
	// или c.Status(http.StatusNoContent)
}

// UploadImages POST /api/products/{id}/images
func (h *ProductHandler) UploadImages(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID format"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data: " + err.Error()})
		return
	}

	files := form.File["files"]
	altTexts := form.Value["alt_texts"]
	isMainFlags := form.Value["is_main"]

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files uploaded"})
		return
	}

	var uploadedImages []model.ProductImage
	for i, fileHeader := range files {
		// --- Логика сохранения файла ---
		filename := fmt.Sprintf("%s-%s", uuid.New().String(), filepath.Base(fileHeader.Filename))
		savePath := "./uploads/" + filename
		if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save file %s: %s", fileHeader.Filename, err.Error())})
			return
		}
		fileURL := "http://" + c.Request.Host + "/uploads/" + filename

		// --- Создание записи в БД ---
		image := model.ProductImage{
			ID:        uuid.New(),
			ProductID: productID,
			URL:       fileURL,
		}
		if i < len(altTexts) {
			image.AltText = altTexts[i]
		}
		if i < len(isMainFlags) {
			image.IsMain, _ = strconv.ParseBool(isMainFlags[i])
		}

		if err := h.repo.CreateImage(c.Request.Context(), &image); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image metadata: " + err.Error()})
			return
		}
		uploadedImages = append(uploadedImages, image)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Изображения успешно загружены",
		"images":  uploadedImages,
	})
}

// GetCategories GET /api/products/categories
func (h *ProductHandler) GetCategories(c *gin.Context) {
	categories, err := h.repo.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// GetSizeChart GET /api/products/sizes
func (h *ProductHandler) GetSizeChart(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'category' is required"})
		return
	}

	sizeChart, err := h.repo.GetSizeChart(c.Request.Context(), category)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("size chart for category '%s' not found", category)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get size chart"})
		return
	}

	c.JSON(http.StatusOK, sizeChart)
}
