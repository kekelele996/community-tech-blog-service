package dto

import "time"

// AuditQuery 审计日志查询
type AuditQuery struct {
	PageQuery
	UserID uint   `form:"user_id" binding:"omitempty,min=1"`
	Action string `form:"action"`
	Module string `form:"module"`
}

// AuditLogDTO 审计日志展示
type AuditLogDTO struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Role      int       `json:"role"`
	Action    string    `json:"action"`
	Module    string    `json:"module"`
	Detail    string    `json:"detail"`
	IP        string    `json:"ip"`
	RequestID string    `json:"request_id"`
	CreatedAt time.Time `json:"created_at"`
}
