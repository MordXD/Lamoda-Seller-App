package handler

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lamoda-seller-app/internal/auth"
	"github.com/lamoda-seller-app/internal/middleware"
	"github.com/lamoda-seller-app/internal/model"
	"github.com/lamoda-seller-app/internal/repository"
	"github.com/lamoda-seller-app/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo  *repository.UserRepository
	validator *validation.UserValidator
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo:  userRepo,
		validator: validation.NewUserValidator(),
	}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  UserDetails `json:"user"`
}

type UserDetails struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type MultipleAccountsResponse struct {
	Accounts []AccountInfo `json:"accounts"`
}

type AccountInfo struct {
	Token string      `json:"token"`
	User  UserDetails `json:"user"`
}

type ValidationErrorResponse struct {
	Error  string                      `json:"error"`
	Fields []validation.ValidationError `json:"fields,omitempty"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ValidationErrorResponse{
			Error: "Invalid request data",
		})
		return
	}

	// Comprehensive validation
	validationErrors := h.validator.ValidateRegistration(req.Name, req.Email, req.Password)
	if validationErrors.HasErrors() {
		c.JSON(http.StatusBadRequest, ValidationErrorResponse{
			Error:  "Validation failed",
			Fields: validationErrors,
		})
		return
	}

	// Normalize email (trim and lowercase)
	normalizedEmail := strings.ToLower(strings.TrimSpace(req.Email))

	// Check if user already exists
	existingUser, err := h.userRepo.FindByEmail(c.Request.Context(), normalizedEmail)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusConflict, ValidationErrorResponse{
			Error: "User with this email already exists",
			Fields: []validation.ValidationError{
				{Field: "email", Message: "This email is already registered"},
			},
		})
		return
	}

	// Hash password with higher cost for better security
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Create user
	user := &model.User{
		Name:           strings.TrimSpace(req.Name),
		Email:          normalizedEmail,
		HashedPassword: string(hashedPassword),
	}

	if err := h.userRepo.Create(c.Request.Context(), user); err != nil {
		log.Printf("Error creating user: %v", err)
		
		// Check if it's a duplicate email error
		if strings.Contains(err.Error(), "uni_users_email") {
			c.JSON(http.StatusConflict, ValidationErrorResponse{
				Error: "User with this email already exists",
				Fields: []validation.ValidationError{
					{Field: "email", Message: "This email is already registered"},
				},
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User: UserDetails{
			ID:        user.ID.String(),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		},
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ValidationErrorResponse{
			Error: "Invalid request data",
		})
		return
	}

	// Basic email validation for login
	emailErrors := h.validator.EmailValidator.ValidateEmail(req.Email)
	if emailErrors.HasErrors() {
		c.JSON(http.StatusBadRequest, ValidationErrorResponse{
			Error:  "Invalid email format",
			Fields: emailErrors,
		})
		return
	}

	// Normalize email
	normalizedEmail := strings.ToLower(strings.TrimSpace(req.Email))

	// Find user by email
	user, err := h.userRepo.FindByEmail(c.Request.Context(), normalizedEmail)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User: UserDetails{
			ID:        user.ID.String(),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		},
	})
}

func (h *UserHandler) ValidateToken(c *gin.Context) {
	var req ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// Validate the token
	userID, err := auth.ParseToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Get user from database
	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: req.Token,
		User: UserDetails{
			ID:        user.ID.String(),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		},
	})
}

func (h *UserHandler) ValidateMultipleTokens(c *gin.Context) {
	var tokens struct {
		Tokens []string `json:"tokens" binding:"required"`
	}

	if err := c.ShouldBindJSON(&tokens); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	var accounts []AccountInfo

	for _, tokenStr := range tokens.Tokens {
		// Validate each token
		userID, err := auth.ParseToken(tokenStr)
		if err != nil {
			continue // Skip invalid tokens
		}

		// Get user from database
		user, err := h.userRepo.GetByID(c.Request.Context(), userID)
		if err != nil || user == nil {
			continue // Skip if user not found
		}

		accounts = append(accounts, AccountInfo{
			Token: tokenStr,
			User: UserDetails{
				ID:        user.ID.String(),
				Email:     user.Email,
				Name:      user.Name,
				CreatedAt: user.CreatedAt,
			},
		})
	}

	c.JSON(http.StatusOK, MultipleAccountsResponse{
		Accounts: accounts,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user ID from JWT middleware
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication"})
		return
	}

	// Get user from database
	user, err := h.userRepo.GetByID(c.Request.Context(), userUUID)
	if err != nil {
		log.Printf("Error getting user profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, UserDetails{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from JWT middleware
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ValidationErrorResponse{
			Error: "Invalid request data",
		})
		return
	}

	// Validate profile update data
	validationErrors := h.validator.ValidateProfileUpdate(req.Name, req.Email)
	if validationErrors.HasErrors() {
		c.JSON(http.StatusBadRequest, ValidationErrorResponse{
			Error:  "Validation failed",
			Fields: validationErrors,
		})
		return
	}

	// Get current user
	user, err := h.userRepo.GetByID(c.Request.Context(), userUUID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Normalize email
	normalizedEmail := strings.ToLower(strings.TrimSpace(req.Email))

	// Check if email is being changed and if it's already taken
	if normalizedEmail != user.Email {
		existingUser, err := h.userRepo.FindByEmail(c.Request.Context(), normalizedEmail)
		if err != nil {
			log.Printf("Error checking email availability: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		if existingUser != nil && existingUser.ID != userUUID {
			c.JSON(http.StatusConflict, ValidationErrorResponse{
				Error: "Email is already taken",
				Fields: []validation.ValidationError{
					{Field: "email", Message: "This email is already registered"},
				},
			})
			return
		}
	}

	// Update user fields
	user.Name = strings.TrimSpace(req.Name)
	user.Email = normalizedEmail

	// Save updated user
	if err := h.userRepo.Update(c.Request.Context(), user); err != nil {
		log.Printf("Error updating profile: %v", err)
		
		// Check if it's a duplicate email error
		if strings.Contains(err.Error(), "uni_users_email") {
			c.JSON(http.StatusConflict, ValidationErrorResponse{
				Error: "Email is already taken",
				Fields: []validation.ValidationError{
					{Field: "email", Message: "This email is already registered"},
				},
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, UserDetails{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	})
} 	