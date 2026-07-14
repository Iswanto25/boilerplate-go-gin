package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Module struct {
	ID        uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}