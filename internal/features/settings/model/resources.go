package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourceAction string

const (
	ActionCreate  ResourceAction = "create"
	ActionRead    ResourceAction = "read"
	ActionUpdate  ResourceAction = "update"
	ActionDelete  ResourceAction = "delete"
	ActionExport  ResourceAction = "export"
	ActionImport  ResourceAction = "import"
	ActionPublish ResourceAction = "publish"
)

type Resource struct {
	ID              uuid.UUID        `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name            string           `gorm:"type:varchar(100);not null" json:"name"`
	ModuleId        uuid.UUID        `gorm:"type:uuid;not null;index" json:"moduleId"`
	Module          Module           `gorm:"foreignKey:ModuleId;constraint:OnDelete:CASCADE" json:"module"`
	AvailableAction []ResourceAction `gorm:"type:jsonb;serializer:json" json:"availableAction"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
