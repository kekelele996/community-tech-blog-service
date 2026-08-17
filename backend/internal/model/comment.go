package model

import "time"

// Comment 评论实体：状态 Status（1 正常 / 0 已删除，constants.UserStatus 复用）
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ArticleID uint      `gorm:"index;not null" json:"article_id"`
	AuthorID  uint      `gorm:"index;not null" json:"author_id"`
	ParentID  uint      `gorm:"index;default:0" json:"parent_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Status    int       `gorm:"not null;default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Author User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}
