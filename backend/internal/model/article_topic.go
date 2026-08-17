package model

import "time"

// ArticleTopic 文章-话题关联表（多对多）
type ArticleTopic struct {
	ArticleID uint      `gorm:"primaryKey;autoIncrement:false" json:"article_id"`
	TopicID   uint      `gorm:"primaryKey;autoIncrement:false" json:"topic_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (ArticleTopic) TableName() string { return "article_topics" }
