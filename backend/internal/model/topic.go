package model

import "time"

// Topic 话题实体：状态 Status（0 禁用 / 1 启用，constants.UserStatus 复用）
type Topic struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:64;uniqueIndex;not null" json:"name"`
	Description  string    `gorm:"size:512" json:"description"`
	ArticleCount int       `gorm:"not null;default:0" json:"article_count"`
	Status       int       `gorm:"not null;default:1" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Articles []Article `gorm:"many2many:article_topics;" json:"articles,omitempty"`
}
