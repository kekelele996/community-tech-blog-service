package dto

import "time"

// CreateTopicRequest 创建话题
type CreateTopicRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=64"`
	Description string `json:"description" binding:"omitempty,max=512"`
}

// UpdateTopicRequest 更新话题
type UpdateTopicRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=64"`
	Description string `json:"description" binding:"omitempty,max=512"`
}

// TopicQuery 话题列表查询
type TopicQuery struct {
	PageQuery
	Keyword string `form:"keyword"`
	Status  int    `form:"status" binding:"omitempty,oneof=0 1"`
}

// TopicDTO 话题展示
type TopicDTO struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	ArticleCount int       `json:"article_count"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// TopicDetailDTO 话题详情（热门文章 + 最新文章）
type TopicDetailDTO struct {
	Topic      TopicDTO             `json:"topic"`
	HotArticles  []ArticleListItemDTO `json:"hot_articles"`
	NewArticles  []ArticleListItemDTO `json:"new_articles"`
}
