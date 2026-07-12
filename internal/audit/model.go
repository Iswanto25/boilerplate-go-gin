package audit

import (
	"encoding/json"
	"time"
)

type Logs struct {
	ID           int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	Date         *string          `gorm:"type:varchar(20);index" json:"date"`
	Name         *string          `gorm:"type:varchar(100)" json:"name"`
	Role         *string          `gorm:"type:varchar(50)" json:"role"`
	Host         *string          `gorm:"type:varchar(255)" json:"host"`
	Status       *string          `gorm:"type:varchar(10);index" json:"status"`
	RequestBody  json.RawMessage  `gorm:"type:jsonb" json:"requestBody"`
	ResponseData json.RawMessage  `gorm:"type:jsonb" json:"responseData"`
	UserID       *string          `gorm:"type:uuid;index" json:"userId"`
	IP           *string          `gorm:"type:varchar(45)" json:"ip"`
	Method       *string          `gorm:"type:varchar(10);index" json:"method"`
	ReqID        *string          `gorm:"type:varchar(50);index" json:"reqId"`
	CreatedAt    time.Time        `gorm:"autoCreateTime;index" json:"createdAt"`
}

func (Logs) TableName() string {
	return "logs"
}
