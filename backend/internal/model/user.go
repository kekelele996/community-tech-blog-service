package model

import "time"

// User 用户实体：角色 Role（1 普通用户 / 2 管理员，constants.RoleType）、状态 Status（0 禁用 / 1 启用，constants.UserStatus）
type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Email        string     `gorm:"size:128;uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Nickname     string     `gorm:"size:64;not null" json:"nickname"`
	Avatar       string     `gorm:"size:512" json:"avatar"`
	Bio          string     `gorm:"size:512" json:"bio"`
	TechTags     string     `gorm:"size:1024" json:"tech_tags"`
	Role         int        `gorm:"not null;default:1" json:"role"`
	Status       int        `gorm:"not null;default:1" json:"status"`
	GithubID     string     `gorm:"size:128;index" json:"github_id"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
