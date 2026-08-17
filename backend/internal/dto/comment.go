package dto

import "time"

// CreateCommentRequest 发表评论
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=1000"`
	ParentID uint   `json:"parent_id" binding:"omitempty,min=1"`
}

// CommentDTO 评论展示
type CommentDTO struct {
	ID        uint      `json:"id"`
	ArticleID uint      `json:"article_id"`
	AuthorID  uint      `json:"author_id"`
	ParentID  uint      `json:"parent_id"`
	Content   string    `json:"content"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Author    *UserProfileDTO `json:"author"`
}
