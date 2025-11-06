package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type Users struct {
	ID        int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string         `gorm:"type:varchar(100);not null" json:"username"`
	Email     string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	// Wallet字段移除以避免循环引用，查询时可通过关联查询获取
}

func (Users) TableName() string {
	return "users"
}
