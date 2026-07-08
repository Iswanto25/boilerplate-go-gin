package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Resource struct {
	ID              uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name            string         `gorm:"type:varchar(100);not null" json:"name"`
	ModuleId        uuid.UUID      `gorm:"type:uuid;not null;index" json:"moduleId"`
	Module          Module         `gorm:"foreignKey:ModuleId" json:"module"`
	AvailableAction []string       `gorm:"type:jsonb;serializer:json" json:"availableAction"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
