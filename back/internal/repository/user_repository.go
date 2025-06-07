// internal/repository/user_repository.go

package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/lamoda-seller-app/internal/model"
)

// --- НОВАЯ ОШИБКА ---
var ErrInsufficientFunds = errors.New("insufficient funds")

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// ... существующие методы ...

// --- НОВЫЙ МЕТОД ДЛЯ АТОМАРНОГО ОБНОВЛЕНИЯ БАЛАНСА ---

// UpdateBalance атомарно изменяет баланс пользователя.
// amount может быть положительным (пополнение) или отрицательным (снятие).
// Метод проверяет, что баланс не станет отрицательным.
func (r *UserRepository) UpdateBalance(ctx context.Context, userID uuid.UUID, amount int64) error {
	// Для снятия средств (amount < 0) мы добавляем условие в WHERE,
	// чтобы запрос не выполнился, если итоговый баланс будет меньше нуля.
	// Для пополнения (amount >= 0) это условие всегда будет истинным.
	tx := r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ? AND balance_kopecks + ? >= 0", userID, amount).
		Update("balance_kopecks", gorm.Expr("balance_kopecks + ?", amount))

	if tx.Error != nil {
		return tx.Error
	}

	// Если ни одна строка не была затронута, это означает, что условие WHERE не выполнилось.
	// В нашем случае это значит, что на счете недостаточно средств для снятия.
	if tx.RowsAffected == 0 {
		// Мы должны проверить, существует ли пользователь вообще, чтобы не возвращать
		// ложную ошибку о нехватке средств для несуществующего ID.
		var count int64
		r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Count(&count)
		if count > 0 {
			return ErrInsufficientFunds
		}
	}

	return nil
}

// --- СУЩЕСТВУЮЩИЕ МЕТОДЫ ---

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, email, newHashedPassword string) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("email = ?", email).
		Update("hashed_password", newHashedPassword).
		Error
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// ... методы для связанных аккаунтов остаются без изменений ...
func (r *UserRepository) LinkAccounts(ctx context.Context, primaryUserID, linkedUserID uuid.UUID) error {
	if primaryUserID == linkedUserID {
		return errors.New("cannot link an account to itself")
	}
	link := &model.AccountLink{
		PrimaryUserID: primaryUserID,
		LinkedUserID:  linkedUserID,
	}
	return r.db.WithContext(ctx).Where(link).FirstOrCreate(link).Error
}

func (r *UserRepository) CheckAccountLink(ctx context.Context, primaryUserID, linkedUserID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.AccountLink{}).
		Where("primary_user_id = ? AND linked_user_id = ?", primaryUserID, linkedUserID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) GetLinkedAccounts(ctx context.Context, primaryUserID uuid.UUID) ([]model.User, error) {
    var users []model.User
    var linkedIDs []uuid.UUID
    err := r.db.WithContext(ctx).Model(&model.AccountLink{}).
        Where("primary_user_id = ?", primaryUserID).
        Pluck("linked_user_id", &linkedIDs).Error
    if err != nil {
        return nil, err
    }
    if len(linkedIDs) == 0 {
        return []model.User{}, nil
    }
    err = r.db.WithContext(ctx).Where("id IN ?", linkedIDs).Find(&users).Error
    return users, err
}