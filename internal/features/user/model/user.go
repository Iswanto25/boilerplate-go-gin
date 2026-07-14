package model

import (
	"time"

	"github.com/google/uuid"
	settingsModel "github.com/edustack/go-boilerplate/internal/features/settings/model"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser:
		return true
	}
	return false
}

type User struct {
	ID        uuid.UUID          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name      string             `gorm:"type:varchar(100);not null" json:"name"`
	Email     string             `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string             `gorm:"type:varchar(255);not null" json:"-"`
	RoleId    uuid.UUID          `gorm:"column:roleId;type:uuid;not null" json:"roleId"`
	Role      settingsModel.Role `gorm:"foreignKey:RoleId;references:ID" json:"role,omitempty"`
	Profile   Profile            `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"profile,omitempty"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

type Profile struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"column:userId;type:uuid;uniqueIndex;not null" json:"userId"`
	Phone     *string   `gorm:"type:varchar(15);null" json:"phone"`
	Address   *string   `gorm:"type:varchar(255);null" json:"address"`
	Photo     *string   `gorm:"type:varchar(255);null" json:"photo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
