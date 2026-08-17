package dto

import "time"

// CreateCollectionRequest 创建收藏夹
type CreateCollectionRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=64"`
	Description string `json:"description" binding:"omitempty,max=255"`
	Visibility  int    `json:"visibility" binding:"omitempty,oneof=0 1"`
}

// UpdateCollectionRequest 更新收藏夹
type UpdateCollectionRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=64"`
	Description string `json:"description" binding:"omitempty,max=255"`
	Visibility  int    `json:"visibility" binding:"omitempty,oneof=0 1"`
}

// CollectionDTO 收藏夹展示
type CollectionDTO struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Visibility   int       `json:"visibility"`
	ArticleCount int       `json:"article_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CollectionDetailDTO 收藏夹详情（含文章列表）
type CollectionDetailDTO struct {
	CollectionDTO
	Articles []ArticleListItemDTO `json:"articles"`
}

// AddArticleRequest 添加文章到收藏夹
type AddArticleRequest struct {
	ArticleID uint `json:"article_id" binding:"required,min=1"`
}
