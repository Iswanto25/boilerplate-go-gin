package audit

import (
	"encoding/json"
	"time"
)


type Logs struct {
	ID        int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	Date      *string          `gorm:"type:varchar(20)" json:"date"`
	Name      *string          `gorm:"type:varchar(100)" json:"name"`
	Role      *string          `gorm:"type:varchar(50)" json:"role"`
	Host      *string          `gorm:"type:varchar(255)" json:"host"`
	Status    *string          `gorm:"type:varchar(10)" json:"status"`
	Data      json.RawMessage `gorm:"type:jsonb" json:"data"`
	UserID    *string          `gorm:"type:uuid;index" json:"userId"`
	IP        *string          `gorm:"type:varchar(45)" json:"ip"`
	Ex       *string           `gorm:"type:varchar(255)" json:"ex"`
	Method    *string          `gorm:"type:varchar(10)" json:"method"`
	CreatedAt time.Time        `gorm:"autoCreateTime;index" json:"createdAt"`
}

func (Logs) TableName() string {
	return "logs"
}
