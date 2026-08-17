package dto

import "time"

// NotificationQuery 通知列表查询
type NotificationQuery struct {
	PageQuery
	UnreadOnly bool `form:"unread_only"`
}

// NotificationDTO 通知展示
type NotificationDTO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	ActorID   uint      `json:"actor_id"`
	Type      string    `json:"type"`
	TargetID  uint      `json:"target_id"`
	Content   string    `json:"content"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	Actor     *UserProfileDTO `json:"actor"`
}
