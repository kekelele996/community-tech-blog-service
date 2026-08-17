package dto

import "time"

// CreateArticleRequest 创建文章
type CreateArticleRequest struct {
	Title    string   `json:"title" binding:"required,min=2,max=256"`
	Content  string   `json:"content" binding:"required,min=10"`
	Summary  string   `json:"summary" binding:"omitempty,max=512"`
	CoverURL string   `json:"cover_url" binding:"omitempty,max=512"`
	Status   int      `json:"status" binding:"omitempty,oneof=0 1"`
	TopicIDs []uint   `json:"topic_ids" binding:"omitempty,dive,min=1"`
}

// UpdateArticleRequest 更新文章
type UpdateArticleRequest struct {
	Title    string `json:"title" binding:"required,min=2,max=256"`
	Content  string `json:"content" binding:"required,min=10"`
	Summary  string `json:"summary" binding:"omitempty,max=512"`
	CoverURL string `json:"cover_url" binding:"omitempty,max=512"`
	Status   int    `json:"status" binding:"omitempty,oneof=0 1"`
	TopicIDs []uint `json:"topic_ids" binding:"omitempty,dive,min=1"`
}

// ArticleQuery 文章列表查询
type ArticleQuery struct {
	PageQuery
	Sort     string `form:"sort" binding:"omitempty,oneof=latest hottest"`
	TopicID  uint   `form:"topic_id" binding:"omitempty,min=1"`
	Keyword  string `form:"keyword"`
	AuthorID uint   `form:"author_id" binding:"omitempty,min=1"`
	Status   int    `form:"status" binding:"omitempty,oneof=0 1 2"`
}

// ArticleListItemDTO 文章列表项
type ArticleListItemDTO struct {
	ID            uint           `json:"id"`
	AuthorID      uint           `json:"author_id"`
	Title         string         `json:"title"`
	Summary       string         `json:"summary"`
	CoverURL      string         `json:"cover_url"`
	Status        int            `json:"status"`
	LikeCount     int            `json:"like_count"`
	ViewCount     int            `json:"view_count"`
	FavoriteCount int            `json:"favorite_count"`
	CommentCount  int            `json:"comment_count"`
	HotScore      int            `json:"hot_score"`
	PublishedAt   *time.Time     `json:"published_at"`
	CreatedAt     time.Time      `json:"created_at"`
	Author        *UserProfileDTO `json:"author"`
	Topics        []TopicBriefDTO `json:"topics"`
}

// TopicBriefDTO 话题简要信息
type TopicBriefDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ArticleDetailDTO 文章详情（含当前用户点赞/收藏状态）
type ArticleDetailDTO struct {
	ArticleListItemDTO
	Content     string         `json:"content"`
	IsLiked     bool           `json:"is_liked"`
	IsCollected bool           `json:"is_collected"`
	Collections []CollectionBriefDTO `json:"collections,omitempty"`
}

// CollectionBriefDTO 收藏夹简要信息（HasArticle 表示该文章是否已收藏进此收藏夹）
type CollectionBriefDTO struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	HasArticle bool   `json:"has_article"`
}

// PublishRequest 发布文章
type PublishRequest struct {
	Status int `json:"status" binding:"required,oneof=1 2"`
}
