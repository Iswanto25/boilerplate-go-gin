package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RolePermission struct {
	ID             uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	RoleId         uuid.UUID       `gorm:"type:uuid;not null;uniqueIndex:idx_role_resource" json:"roleId"`
	Role           Role            `gorm:"foreignKey:RoleId;constraint:OnDelete:CASCADE" json:"role"`
	ResourceId     uuid.UUID       `gorm:"type:uuid;not null;uniqueIndex:idx_role_resource" json:"resourceId"`
	Resource       Resource        `gorm:"foreignKey:ResourceId;constraint:OnDelete:CASCADE" json:"resource"`
	GrantedActions []string        `gorm:"type:jsonb;serializer:json" json:"grantedActions"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (RolePermission) TableName() string {
	return "rolePermissions"
}
