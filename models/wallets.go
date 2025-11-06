package models

import (
	"time"

	"gorm.io/gorm"
)

// Wallet
type Wallets struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int            `gorm:"not null;uniqueIndex" json:"user_id"`
	Balance   float64        `gorm:"type:decimal(12,2);default:0" json:"balance"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	User      Users          `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

func (Wallets) TableName() string {
	return "wallets"
}
