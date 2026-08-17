package model

import "time"

// AuditLog 操作审计日志实体（横切关注点 2）
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Username  string    `gorm:"size:64" json:"username"`
	Role      int       `json:"role"`
	Action    string    `gorm:"size:64;not null" json:"action"`
	Module    string    `gorm:"size:64" json:"module"`
	Detail    string    `gorm:"size:1024" json:"detail"`
	IP        string    `gorm:"size:64" json:"ip"`
	RequestID string    `gorm:"size:64;index" json:"request_id"`
	CreatedAt time.Time `json:"created_at"`
}
