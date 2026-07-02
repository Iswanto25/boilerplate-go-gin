package model

import (
	"encoding/json"
	"time"
)

type Logs struct {
	ID   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Date string `gorm:"not null" json:"date"`

	UserID *int64  `gorm:"index" json:"user_id"` // Gunakan pointer agar bisa null
	Name   *string `gorm:"type:varchar(100)" json:"name"`
	Role   *string `gorm:"type:varchar(50)" json:"role"`

	Host   string `gorm:"type:varchar(255);not null" json:"host"`
	Method string `gorm:"type:varchar(10);not null" json:"method"`
	Status int    `gorm:"not null" json:"status"` // Gunakan int untuk HTTP Status Code

	Data      json.RawMessage `gorm:"type:jsonb" json:"data"` // Gunakan jsonb untuk performa di Postgres
	CreatedAt time.Time       `gorm:"index;not null;autoCreateTime" json:"created_at"`
}
