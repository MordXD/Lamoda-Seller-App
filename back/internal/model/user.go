package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Name           string    `gorm:"size:255;not null" json:"name"`
	Email          string    `gorm:"size:255;not null;unique" json:"email"`
	HashedPassword string    `gorm:"not null" json:"-"`
	// --- НОВОЕ ПОЛЕ ---
	// Храним баланс в копейках, чтобы избежать проблем с float.
	// `not null;default:0` гарантирует, что у новых пользователей баланс будет 0.
	BalanceKopecks int64     `gorm:"not null;default:0" json:"balance_kopecks"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"-"` // Скроем UpdatedAt из JSON для чистоты
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return
}
