package audit

import (
	"encoding/json"
	"time"
)

// Logs merepresentasikan satu entri activity log yang disimpan ke database.
type Logs struct {
	ID     int64  `gorm:"primaryKey;autoIncrement"          json:"id"`
	Date   string `gorm:"type:varchar(10);not null;index"   json:"date"`

	UserID *string `gorm:"type:uuid;index"                  json:"user_id"`
	Name   *string `gorm:"type:varchar(100)"                json:"name"`
	Role   *string `gorm:"type:varchar(50)"                 json:"role"`

	Host   string `gorm:"type:varchar(255);not null"        json:"host"`
	Path   string `gorm:"type:varchar(500);not null"        json:"path"`
	Method string `gorm:"type:varchar(10);not null"         json:"method"`
	Status int    `gorm:"not null"                          json:"status"`

	Data json.RawMessage `gorm:"type:jsonb"                json:"data"`

	CreatedAt time.Time `gorm:"not null;autoCreateTime;index"     json:"created_at"`
}

// TableName mengoverride nama tabel default GORM
func (Logs) TableName() string {
	return "logs"
}
