package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/lamoda-seller-app/internal/auth"
	"github.com/lamoda-seller-app/internal/middleware"
	"github.com/lamoda-seller-app/internal/model"
	"github.com/lamoda-seller-app/internal/repository"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// --- Структуры для запросов и ответов ---

type RegisterRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type RegisterResponse struct {
	Token             string `json:"token"`
	TemporaryPassword string `json:"temporary_password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"` // Добавим валидацию на мин. длину
}

type UpdateProfileRequest struct {
	Name string `json:"name" binding:"required"`
}

type LinkAccountRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SwitchAccountRequest struct {
	TargetUserID uuid.UUID `json:"target_user_id" binding:"required"`
}

// --- НОВЫЕ СТРУКТУРЫ ДЛЯ БАЛАНСА ---
type AddBalanceRequest struct {
	// Сумма должна быть больше 0
	AmountKopecks int64 `json:"amount_kopecks" binding:"required,gt=0"`
}

type WithdrawBalanceRequest struct {
	// Сумма должна быть больше 0
	AmountKopecks int64 `json:"amount_kopecks" binding:"required,gt=0"`
}

// --- Хендлеры ---

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	_, err := h.repo.GetByEmail(c.Request.Context(), req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user with this email already exists"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	tmpPassword, err := auth.GenerateTemporaryPassword(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate temporary password"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(tmpPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &model.User{
		Email:          req.Email,
		Name:           req.Name,
		HashedPassword: string(hashed),
	}

	if err := h.repo.Create(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user creation failed"})
		return
	}

	token, err := auth.GenerateJWT(user.ID) // Используем переименованную функцию
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		Token:             token,
		TemporaryPassword: tmpPassword,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	user, err := h.repo.GetByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}

// ChangePassword - новый хендлер для смены пароля (защищенный)
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	user, err := h.repo.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash new password"})
		return
	}

	if err := h.repo.UpdatePassword(c.Request.Context(), user.Email, string(newHash)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}

// GetProfile - реализация для вашего роута GET /api/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	user, err := h.repo.GetByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, user) // Теперь можно просто вернуть всю структуру User
}

// UpdateProfile - реализация для вашего роута PUT /api/profile (для смены имени)
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	user, err := h.repo.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.Name = req.Name

	if err := h.repo.UpdateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}

// ValidateToken и ValidateMultipleTokens - оставляем заглушки, т.к. они не относятся к основной задаче
func (h *UserHandler) ValidateToken(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (h *UserHandler) ValidateMultipleTokens(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

// LinkAccount - привязывает другой аккаунт к текущему аутентифицированному.
// Требует ввода пароля от привязываемого аккаунта для подтверждения владения.
func (h *UserHandler) LinkAccount(c *gin.Context) {
	var req LinkAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	primaryUserID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	// Находим аккаунт, который хотим привязать
	linkedUser, err := h.repo.GetByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account to link not found"})
		return
	}

	// Проверяем пароль от привязываемого аккаунта
	if err := bcrypt.CompareHashAndPassword([]byte(linkedUser.HashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials for the account to be linked"})
		return
	}

	// Создаем связь
	if err := h.repo.LinkAccounts(c.Request.Context(), primaryUserID, linkedUser.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to link accounts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account linked successfully"})
}

// SwitchAccount - выполняет переключение на привязанный аккаунт.
func (h *UserHandler) SwitchAccount(c *gin.Context) {
	var req SwitchAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	currentUserID := c.MustGet(middleware.UserIDKey).(uuid.UUID)
	targetUserID := req.TargetUserID

	// Проверяем, разрешено ли переключение
	isAllowed, err := h.repo.CheckAccountLink(c.Request.Context(), currentUserID, targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if !isAllowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "switching to this account is not allowed"})
		return
	}

	// Генерируем новый токен для целевого пользователя
	token, err := auth.GenerateJWT(targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}

// GetLinkedAccounts - возвращает список аккаунтов, доступных для переключения.
func (h *UserHandler) GetLinkedAccounts(c *gin.Context) {
	primaryUserID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	accounts, err := h.repo.GetLinkedAccounts(c.Request.Context(), primaryUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve linked accounts"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// --- НОВЫЕ ХЕНДЛЕРЫ ДЛЯ БАЛАНСА ---

// GetBalance возвращает текущий баланс пользователя.
func (h *UserHandler) GetBalance(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	user, err := h.repo.GetByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance_kopecks": user.BalanceKopecks})
}

// AddBalance пополняет баланс пользователя.
func (h *UserHandler) AddBalance(c *gin.Context) {
	var req AddBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	if err := h.repo.UpdateBalance(c.Request.Context(), userID, req.AmountKopecks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "balance updated successfully"})
}

// WithdrawBalance снимает средства с баланса пользователя.
func (h *UserHandler) WithdrawBalance(c *gin.Context) {
	var req WithdrawBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	// Передаем отрицательное значение для снятия
	err := h.repo.UpdateBalance(c.Request.Context(), userID, -req.AmountKopecks)
	if err != nil {
		if errors.Is(err, repository.ErrInsufficientFunds) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "withdrawal successful"})
}