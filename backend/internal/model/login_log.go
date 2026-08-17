package model

import "time"

// LoginLog 登录日志：用于平台日活（DAU）统计，LoginDate 格式 YYYY-MM-DD
type LoginLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	LoginDate string    `gorm:"size:10;index;not null" json:"login_date"`
	IP        string    `gorm:"size:64" json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}
