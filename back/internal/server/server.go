package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

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
	// Setup GORM
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("❌ failed to connect to DB: %w", err)
	}

	// Auto-migrate database tables
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, fmt.Errorf("❌ failed to migrate database: %w", err)
	}

	log.Println("✅ Connected to database and migrated tables")

	// Setup Gin
	r := gin.Default()

	// Add CORS middleware for frontend integration
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
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
	}

	// Protected routes - require authentication
	authorized := r.Group("/api")
	authorized.Use(middleware.JWTAuthMiddleware(cfg.JWTSecret))
	{
		// User profile endpoints
		authorized.GET("/profile", userHandler.GetProfile)
		authorized.PUT("/profile", userHandler.UpdateProfile)
		
		// Add more protected routes here as needed
		// authorized.GET("/dashboard", dashboardHandler.GetDashboard)
		// authorized.POST("/tasks", taskHandler.CreateTask)
	}

	return &Server{
		Engine: r,
		DB:     db,
		Config: cfg,
	}, nil
}

func (s *Server) Run() {
	srv := &http.Server{
		Addr:    ":" + s.Config.ServerPort,
		Handler: s.Engine,
	}

	go func() {
		log.Printf("🚀 Server running on port %s\n", s.Config.ServerPort)
		log.Printf("📚 API endpoints available:\n")
		log.Printf("   POST /auth/register - User registration\n")
		log.Printf("   POST /auth/login - User login\n")
		log.Printf("   GET  /api/profile - Get user profile (protected)\n")
		log.Printf("   PUT  /api/profile - Update user profile (protected)\n")
		log.Printf("   GET  /health - Health check\n")
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %s", err)
	}

	log.Println("✅ Server exited properly")
}