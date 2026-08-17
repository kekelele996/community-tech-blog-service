package model

import "time"

// CollectionArticle 收藏夹-文章关联：唯一索引保证同一收藏夹不重复收藏同一文章
type CollectionArticle struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CollectionID uint      `gorm:"uniqueIndex:uk_collection_article;not null" json:"collection_id"`
	ArticleID    uint      `gorm:"uniqueIndex:uk_collection_article;not null" json:"article_id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
}
