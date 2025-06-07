package model

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID             uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"` // Добавили json тег
    Name           string    `gorm:"size:255;not null" json:"name"`
    Email          string    `gorm:"size:255;not null;unique" json:"email"`
    HashedPassword string    `gorm:"not null" json:"-"` // Пароль не должен сериализоваться
    CreatedAt      time.Time `json:"created_at"`        // Добавили json тег
    UpdatedAt      time.Time
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
    if user.ID == uuid.Nil {
        user.ID = uuid.New()
    }
    return
}