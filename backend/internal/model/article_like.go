package model

import "time"

// ArticleLike 文章点赞记录：唯一索引保证同一用户对同一文章只能点赞一次
type ArticleLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ArticleID uint      `gorm:"uniqueIndex:uk_article_user;not null" json:"article_id"`
	UserID    uint      `gorm:"uniqueIndex:uk_article_user;not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
