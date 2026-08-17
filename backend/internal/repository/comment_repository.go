package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/pkg/pagination"
)

// CommentRepository 评论仓储
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 构造评论仓储
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create 创建评论
func (r *CommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	if err := r.db.WithContext(ctx).Create(comment).Error; err != nil {
		return fmt.Errorf("create comment: %w", err)
	}
	return nil
}

// ListByArticle 分页查询文章评论
func (r *CommentRepository) ListByArticle(ctx context.Context, articleID uint, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Comment{}).Where("article_id = ? AND status = ?", articleID, 1)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count comments: %w", err)
	}
	if err := q.Preload("Author").Order("id ASC").Offset(pagination.Offset(page, pageSize)).Limit(pageSize).Find(&comments).Error; err != nil {
		return nil, 0, fmt.Errorf("list comments: %w", err)
	}
	return comments, total, nil
}

// FindByID 按 ID 查询评论
func (r *CommentRepository) FindByID(ctx context.Context, id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := r.db.WithContext(ctx).First(&comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find comment by id: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("find comment by id: %w", err)
	}
	return &comment, nil
}

// Delete 删除评论（软删：status=0）
func (r *CommentRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Model(&model.Comment{}).Where("id = ?", id).Update("status", 0)
	if res.Error != nil {
		return fmt.Errorf("delete comment: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete comment: %w", ErrNotFound)
	}
	return nil
}

// Count 评论总数
func (r *CommentRepository) Count(ctx context.Context) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Comment{}).Where("status = ?", 1).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count comments: %w", err)
	}
	return total, nil
}
