package model

import "time"

// Notification 站内通知实体：类型 Type（follow/like/comment/system，constants.NotificationType）
type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	ActorID   uint      `gorm:"index;not null" json:"actor_id"`
	Type      string    `gorm:"size:32;not null" json:"type"`
	TargetID  uint      `gorm:"index" json:"target_id"`
	Content   string    `gorm:"size:512" json:"content"`
	IsRead    bool      `gorm:"not null;default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
