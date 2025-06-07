package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountLink представляет собой связь между двумя учетными записями,
// разрешающую переключение с PrimaryUser на LinkedUser без пароля.
type AccountLink struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;"`
	PrimaryUserID uuid.UUID `gorm:"type:uuid;not null;index:idx_primary_linked,unique"`
	LinkedUserID  uuid.UUID `gorm:"type:uuid;not null;index:idx_primary_linked,unique"`
	CreatedAt     time.Time

	// Опционально: можно добавить связи для получения полной информации
	// PrimaryUser User `gorm:"foreignKey:PrimaryUserID"`
	// LinkedUser  User `gorm:"foreignKey:LinkedUserID"`
}

// BeforeCreate будет вызываться перед созданием записи
func (link *AccountLink) BeforeCreate(tx *gorm.DB) (err error) {
	if link.ID == uuid.Nil {
		link.ID = uuid.New()
	}
	return
}