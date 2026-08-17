package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/model"
)

// LikeRepository 文章点赞仓储
type LikeRepository struct {
	db *gorm.DB
}

// NewLikeRepository 构造点赞仓储
func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

// Create 点赞
func (r *LikeRepository) Create(ctx context.Context, like *model.ArticleLike) error {
	if err := r.db.WithContext(ctx).Create(like).Error; err != nil {
		if isDuplicateErr(err) {
			return fmt.Errorf("create like: %w", ErrDuplicateEntry)
		}
		return fmt.Errorf("create like: %w", err)
	}
	return nil
}

// Delete 取消点赞
func (r *LikeRepository) Delete(ctx context.Context, articleID, userID uint) error {
	res := r.db.WithContext(ctx).Where("article_id = ? AND user_id = ?", articleID, userID).Delete(&model.ArticleLike{})
	if res.Error != nil {
		return fmt.Errorf("delete like: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete like: %w", ErrNotFound)
	}
	return nil
}

// IsLiked 判断是否已点赞
func (r *LikeRepository) IsLiked(ctx context.Context, articleID, userID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.ArticleLike{}).
		Where("article_id = ? AND user_id = ?", articleID, userID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("check liked: %w", err)
	}
	return count > 0, nil
}
