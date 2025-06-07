package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/lamoda-seller-app/internal/config"
	"github.com/lamoda-seller-app/internal/handler"
	"github.com/lamoda-seller-app/internal/middleware"
	"github.com/lamoda-seller-app/internal/model"
	"github.com/lamoda-seller-app/internal/repository"
)

type Server struct {
	Engine *gin.Engine
	DB     *gorm.DB
	Config *config.Config
}

func Init(cfg *config.Config) (*Server, error) {
	// ... (код подключения к БД и миграции остается без изменений) ...
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// !!! ДОБАВЛЯЕМ НОВЫЕ МОДЕЛИ ПРОДУКТА В МИГРАЦИЮ !!!
	if err := db.AutoMigrate(
		&model.User{},
		&model.AccountLink{},
		&model.Product{},
		&model.ProductVariant{},
		&model.PricePoint{},
		&model.ProductSales{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("✅ Connected to database and migrated tables")

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(corsMiddleware(cfg))

	// --- ИЗМЕНЕНИЕ ЛОГИКИ РОУТИНГА ---

	// Initialize repositories and handlers
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	productHandler := handler.NewProductHandler(productRepo)
	userHandler := handler.NewUserHandler(userRepo)

	// Создаем одну родительскую группу /api
	api := r.Group("/api")
	{
		// Health check endpoint внутри /api
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "ok",
				"timestamp": time.Now().UTC(),
				"version":   "1.0.0",
			})
		})

		// --- Public Routes ---
		// Группа /auth теперь вложена в /api, создавая пути вида /api/auth/*
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/validate-token", userHandler.ValidateToken)
			auth.POST("/validate-tokens", userHandler.ValidateMultipleTokens)
		}

		// --- Protected Routes ---
		protected := api.Group("/")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			// Profile routes
			protected.GET("/profile", userHandler.GetProfile)
			protected.PUT("/profile", userHandler.UpdateProfile)

			// Password routes
			protected.POST("/password/change", userHandler.ChangePassword)

			// Account management routes
			account := protected.Group("/account")
			{
				account.POST("/link", userHandler.LinkAccount)
				account.POST("/switch", userHandler.SwitchAccount)
				account.GET("/links", userHandler.GetLinkedAccounts)
			}

			// --- !!! НОВЫЕ МАРШРУТЫ ДЛЯ БАЛАНСА !!! ---
			balance := protected.Group("/balance")
			{
				balance.GET("", userHandler.GetBalance)                // Получить текущий баланс
				balance.POST("/add", userHandler.AddBalance)           // Пополнить баланс
				balance.POST("/withdraw", userHandler.WithdrawBalance) // Снять с баланса
			}
			// --- !!! НОВЫЕ МАРШРУТЫ ДЛЯ ПРОДАВЦА !!! ---
			seller := protected.Group("/seller")
			{
				products := seller.Group("/products")
				{
					products.GET("", productHandler.ListProducts)
					products.POST("", productHandler.CreateProduct)
					products.GET("/:id", productHandler.GetProductByID)
					products.PUT("/:id", productHandler.UpdateProduct)
					products.DELETE("/:id", productHandler.DeleteProduct)
					products.GET("/:id/stats", productHandler.GetSalesStats)
					products.GET("/:id/price-history", productHandler.GetPriceHistory)
				}
			}
		}
	}

	return &Server{
		Engine: r,
		DB:     db,
		Config: cfg,
	}, nil
}

func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	// ... (код corsMiddleware без изменений)
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		allowedOrigins := strings.Split(os.Getenv("CORS_ORIGINS"), ",")
		if len(allowedOrigins) == 0 || allowedOrigins[0] == "" {
			allowedOrigins = []string{"http://localhost:3000"}
		}

		originAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if strings.TrimSpace(allowedOrigin) == origin {
				originAllowed = true
				break
			}
		}

		if originAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Server) Run() {
	srv := &http.Server{
		Addr:         ":" + s.Config.ServerPort,
		Handler:      s.Engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server running on port %s", s.Config.ServerPort)
		log.Printf("📚 API endpoints available:")
		log.Printf("   POST /api/auth/register - User registration")
		log.Printf("   POST /api/auth/login - User login")
		log.Printf("   GET  /api/profile - Get user profile (protected)")
		log.Printf("   PUT  /api/profile - Update user profile (protected)")
		log.Printf("   POST /api/password/change - Change password (protected)")
		log.Printf("   GET  /api/health - Health check")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("✅ Server exited properly")
}
