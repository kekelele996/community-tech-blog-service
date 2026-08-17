package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/pkg/pagination"
)

// FollowRepository 关注关系仓储
type FollowRepository struct {
	db *gorm.DB
}

// NewFollowRepository 构造关注关系仓储
func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

// Create 创建关注
func (r *FollowRepository) Create(ctx context.Context, follow *model.Follow) error {
	if err := r.db.WithContext(ctx).Create(follow).Error; err != nil {
		if isDuplicateErr(err) {
			return fmt.Errorf("create follow: %w", ErrDuplicateEntry)
		}
		return fmt.Errorf("create follow: %w", err)
	}
	return nil
}

// Delete 取消关注（按关注对删除）
func (r *FollowRepository) Delete(ctx context.Context, followerID, followedID uint) error {
	res := r.db.WithContext(ctx).Where("follower_id = ? AND followed_id = ?", followerID, followedID).Delete(&model.Follow{})
	if res.Error != nil {
		return fmt.Errorf("delete follow: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("delete follow: %w", ErrNotFound)
	}
	return nil
}

// IsFollowing 判断是否已关注
func (r *FollowRepository) IsFollowing(ctx context.Context, followerID, followedID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Follow{}).
		Where("follower_id = ? AND followed_id = ?", followerID, followedID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("check following: %w", err)
	}
	return count > 0, nil
}

// CountFollowers 粉丝数
func (r *FollowRepository) CountFollowers(ctx context.Context, userID uint) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Follow{}).Where("followed_id = ?", userID).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count followers: %w", err)
	}
	return total, nil
}

// CountFollowing 关注数
func (r *FollowRepository) CountFollowing(ctx context.Context, userID uint) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Follow{}).Where("follower_id = ?", userID).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count following: %w", err)
	}
	return total, nil
}

// ListFollowerIDs 粉丝 ID 列表
func (r *FollowRepository) ListFollowerIDs(ctx context.Context, userID uint, page, pageSize int) ([]uint, int64, error) {
	var follows []model.Follow
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Follow{}).Where("followed_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count followers: %w", err)
	}
	if err := q.Order("id DESC").Offset(pagination.Offset(page, pageSize)).Limit(pageSize).Find(&follows).Error; err != nil {
		return nil, 0, fmt.Errorf("list followers: %w", err)
	}
	ids := make([]uint, 0, len(follows))
	for _, f := range follows {
		ids = append(ids, f.FollowerID)
	}
	return ids, total, nil
}

// ListFollowingIDs 关注 ID 列表
func (r *FollowRepository) ListFollowingIDs(ctx context.Context, userID uint, page, pageSize int) ([]uint, int64, error) {
	var follows []model.Follow
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Follow{}).Where("follower_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count following: %w", err)
	}
	if err := q.Order("id DESC").Offset(pagination.Offset(page, pageSize)).Limit(pageSize).Find(&follows).Error; err != nil {
		return nil, 0, fmt.Errorf("list following: %w", err)
	}
	ids := make([]uint, 0, len(follows))
	for _, f := range follows {
		ids = append(ids, f.FollowedID)
	}
	return ids, total, nil
}

// AllFollowingIDs 全部关注 ID（首页关注流）
func (r *FollowRepository) AllFollowingIDs(ctx context.Context, userID uint) ([]uint, error) {
	var follows []model.Follow
	if err := r.db.WithContext(ctx).Where("follower_id = ?", userID).Find(&follows).Error; err != nil {
		return nil, fmt.Errorf("list all following: %w", err)
	}
	ids := make([]uint, 0, len(follows))
	for _, f := range follows {
		ids = append(ids, f.FollowedID)
	}
	return ids, nil
}

// FindByPair 查询关注对
func (r *FollowRepository) FindByPair(ctx context.Context, followerID, followedID uint) (*model.Follow, error) {
	var follow model.Follow
	if err := r.db.WithContext(ctx).Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find follow: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("find follow: %w", err)
	}
	return &follow, nil
}
