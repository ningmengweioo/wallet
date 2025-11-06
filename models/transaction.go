package models

import (
	"time"

	"gorm.io/gorm"
)

// Transaction 交易记录模型
type Transaction struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Type        string         `gorm:"type:varchar(20);not null" json:"type"` // deposit, withdraw, transfer
	FromUserID  int            `json:"from_user_id,omitempty"`
	ToUserID    int            `json:"to_user_id,omitempty"`
	Amount      float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	Status      string         `gorm:"type:varchar(20);default:'completed'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Transaction) TableName() string {
	return "transaction"
}
