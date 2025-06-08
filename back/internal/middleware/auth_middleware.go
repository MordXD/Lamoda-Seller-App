package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lamoda-seller-app/internal/auth"
)

const (
	UserIDKey = "user_id"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("🔐 JWT Auth middleware: проверка авторизации для %s %s", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("❌ JWT Auth: отсутствует Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		log.Printf("🔍 JWT Auth: получен Authorization header: %s", authHeader[:min(len(authHeader), 20)]+"...")

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("❌ JWT Auth: неверный формат Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		log.Printf("🎫 JWT Auth: валидация токена (длина: %d)", len(tokenStr))

		claims, err := auth.ValidateToken(tokenStr)
		if err != nil {
			log.Printf("❌ JWT Auth: ошибка валидации токена: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		log.Printf("✅ JWT Auth: токен валиден, пользователь ID: %s", claims.UserID)

		// Set user ID in context
		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
