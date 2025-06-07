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
	// Setup GORM with proper configuration
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// !!! ДОБАВЛЯЕМ НОВУЮ МОДЕЛЬ В МИГРАЦИЮ !!!
	if err := db.AutoMigrate(&model.User{}, &model.AccountLink{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("✅ Connected to database and migrated tables")

	// Setup Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Add CORS middleware with environment-based configuration
	r.Use(corsMiddleware(cfg))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userRepo)

	// Public routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.POST("/validate-token", userHandler.ValidateToken)
		auth.POST("/validate-tokens", userHandler.ValidateMultipleTokens)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.JWTAuthMiddleware())
	{
		// Profile
		api.GET("/profile", userHandler.GetProfile)
		api.PUT("/profile", userHandler.UpdateProfile)

		// Password
		api.POST("/password/change", userHandler.ChangePassword)

		// --- !!! НОВЫЕ МАРШРУТЫ ДЛЯ УПРАВЛЕНИЯ АККАУНТАМИ !!! ---
		account := api.Group("/account")
		{
			account.POST("/link", userHandler.LinkAccount)
			account.POST("/switch", userHandler.SwitchAccount)
			account.GET("/links", userHandler.GetLinkedAccounts)
		}
	}

	return &Server{
		Engine: r,
		DB:     db,
		Config: cfg,
	}, nil
}

func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Get allowed origins from environment
		allowedOrigins := strings.Split(os.Getenv("CORS_ORIGINS"), ",")
		if len(allowedOrigins) == 0 || allowedOrigins[0] == "" {
			allowedOrigins = []string{"http://localhost:3000"}
		}

		// Check if origin is allowed
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
		log.Printf("   POST /auth/register - User registration")
		log.Printf("   POST /auth/login - User login")
		log.Printf("   GET  /api/profile - Get user profile (protected)")
		log.Printf("   PUT  /api/profile - Update user profile (protected)")
		log.Printf("   GET  /health - Health check")

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
