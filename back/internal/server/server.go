// server/server.go
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
	"github.com/lamoda-seller-app/internal/repository"
)

type Server struct {
	Engine *gin.Engine
	DB     *gorm.DB
	Config *config.Config
}

// Middleware для подробного логирования запросов
func requestLoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("🌐 [%s] %s %s %d %s \"%s\" %s \"%s\" %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.ClientIP,
			param.Method,
			param.StatusCode,
			param.Latency,
			param.Path,
			param.Request.Proto,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

func Init(cfg *config.Config) (*Server, error) {
	log.Printf("🚀 Инициализация сервера...")
	log.Printf("📊 Конфигурация: DB=%s:%s, Server=:%s", cfg.DBHost, cfg.DBPort, cfg.ServerPort)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	log.Printf("🔌 Подключение к базе данных...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("❌ Ошибка подключения к БД: %v", err)
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("❌ Ошибка получения sql.DB: %v", err)
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Подключение к базе данных установлено")

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Printf("🔧 Настройка Gin роутера...")
	r := gin.New()

	// Добавляем middleware для логирования
	r.Use(requestLoggingMiddleware())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware(cfg))

	// Роут для статических файлов (загруженных изображений)
	r.Static("/uploads", "./uploads")

	log.Printf("🏗️ Инициализация репозиториев и обработчиков...")
	// Initialize repositories and handlers
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)

	userHandler := handler.NewUserHandler(userRepo)
	productHandler := handler.NewProductHandler(productRepo)
	orderHandler := handler.NewOrderHandler(orderRepo)
	dashboardHandler := handler.NewDashboardHandler(dashboardRepo)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsRepo)

	log.Printf("🛣️ Настройка маршрутов...")
	// Создаем одну родительскую группу /api
	api := r.Group("/api")
	{
		// Health check endpoint внутри /api
		api.GET("/health", func(c *gin.Context) {
			log.Printf("💓 Health check запрос")
			c.JSON(http.StatusOK, gin.H{
				"status":    "ok",
				"timestamp": time.Now().UTC(),
				"version":   "1.0.0",
			})
		})

		// --- Public Routes ---
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

			// Маршруты для баланса
			balance := protected.Group("/balance")
			{
				balance.GET("", userHandler.GetBalance)
				balance.POST("/add", userHandler.AddBalance)
				balance.POST("/withdraw", userHandler.WithdrawBalance)
			}

			// --- Маршруты для продуктов ---
			products := protected.Group("/products")
			{
				products.GET("", productHandler.ListProducts)
				products.GET("/categories", productHandler.GetCategories)
				products.GET("/sizes", productHandler.GetSizeChart)
				products.POST("", productHandler.CreateProduct)
				products.GET("/:id", productHandler.GetProductByID)
				products.PUT("/:id", productHandler.UpdateProduct)
				products.DELETE("/:id", productHandler.DeleteProduct)
				products.POST("/:id/images", productHandler.UploadImages)
			}
			// --- Маршруты для заказов ---
			orders := protected.Group("/orders")
			{
				orders.GET("", orderHandler.ListOrders)
				orders.GET("/:order_id", orderHandler.GetOrderByID)
				orders.PUT("/:order_id/status", orderHandler.UpdateOrderStatus)
			}

			// --- Маршруты для дашборда ---
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/stats", dashboardHandler.GetStats)
				dashboard.GET("/sales-chart", dashboardHandler.GetSalesChart)
			}
			analytics := protected.Group("/analytics")
			{
				analytics.GET("/top-products", analyticsHandler.GetTopProducts)
				analytics.GET("/categories", analyticsHandler.GetCategoryAnalytics)
				analytics.GET("/size-distribution", analyticsHandler.GetSizeDistribution)
				analytics.GET("/seasonal-trends", analyticsHandler.GetSeasonalTrends)
				analytics.GET("/returns", analyticsHandler.GetReturnsAnalytics)
			}
		}
	}

	log.Printf("✅ Сервер инициализирован успешно")
	return &Server{
		Engine: r,
		DB:     db,
		Config: cfg,
	}, nil
}

func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
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

	log.Printf("🚀 Запуск сервера на порту %s", s.Config.ServerPort)
	log.Printf("🌍 Сервер доступен по адресу: http://localhost:%s", s.Config.ServerPort)
	log.Printf("💊 Health check: http://localhost:%s/api/health", s.Config.ServerPort)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Ошибка запуска сервера: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Получен сигнал остановки, завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Принудительное завершение сервера: %s", err)
	}

	log.Println("✅ Сервер корректно завершил работу")
}
