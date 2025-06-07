// internal/repository/user_repository.go

package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/lamoda-seller-app/internal/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// ... существующие методы ...

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

// --- НОВЫЕ МЕТОДЫ ДЛЯ СВЯЗАННЫХ АККАУНТОВ ---

// LinkAccounts создает запись о связи двух аккаунтов.
func (r *UserRepository) LinkAccounts(ctx context.Context, primaryUserID, linkedUserID uuid.UUID) error {
	// Предотвращаем создание петли (связь аккаунта с самим собой)
	if primaryUserID == linkedUserID {
		return errors.New("cannot link an account to itself")
	}
	link := &model.AccountLink{
		PrimaryUserID: primaryUserID,
		LinkedUserID:  linkedUserID,
	}
	// `FirstOrCreate` предотвратит создание дубликатов
	return r.db.WithContext(ctx).Where(link).FirstOrCreate(link).Error
}

// CheckAccountLink проверяет, существует ли разрешение на переключение.
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

// GetLinkedAccounts возвращает список пользователей, на которых можно переключиться.
func (r *UserRepository) GetLinkedAccounts(ctx context.Context, primaryUserID uuid.UUID) ([]model.User, error) {
    var users []model.User
    
    // Находим ID связанных аккаунтов
    var linkedIDs []uuid.UUID
    err := r.db.WithContext(ctx).Model(&model.AccountLink{}).
        Where("primary_user_id = ?", primaryUserID).
        Pluck("linked_user_id", &linkedIDs).Error
    if err != nil {
        return nil, err
    }

    if len(linkedIDs) == 0 {
        return []model.User{}, nil // Возвращаем пустой слайс, а не nil
    }

    // Получаем информацию о пользователях по их ID
    err = r.db.WithContext(ctx).Where("id IN ?", linkedIDs).Find(&users).Error
    return users, err
}