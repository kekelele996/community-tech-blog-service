package model

import "time"

// Collection 收藏夹实体：可见性 Visibility（0 私密 / 1 公开，constants.Visibility）
type Collection struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	Name         string    `gorm:"size:64;not null" json:"name"`
	Description  string    `gorm:"size:255" json:"description"`
	Visibility   int       `gorm:"not null;default:0" json:"visibility"`
	ArticleCount int       `gorm:"not null;default:0" json:"article_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
