package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/pkg/pagination"
)

// TopicRepository 话题仓储
type TopicRepository struct {
	db *gorm.DB
}

// NewTopicRepository 构造话题仓储
func NewTopicRepository(db *gorm.DB) *TopicRepository {
	return &TopicRepository{db: db}
}

// Create 创建话题
func (r *TopicRepository) Create(ctx context.Context, topic *model.Topic) error {
	if err := r.db.WithContext(ctx).Create(topic).Error; err != nil {
		if isDuplicateErr(err) {
			return fmt.Errorf("create topic: %w", ErrDuplicateEntry)
		}
		return fmt.Errorf("create topic: %w", err)
	}
	return nil
}

// Update 更新话题
func (r *TopicRepository) Update(ctx context.Context, topic *model.Topic) error {
	if err := r.db.WithContext(ctx).Model(topic).Updates(map[string]interface{}{
		"name":        topic.Name,
		"description": topic.Description,
		"status":      topic.Status,
	}).Error; err != nil {
		if isDuplicateErr(err) {
			return fmt.Errorf("update topic: %w", ErrDuplicateEntry)
		}
		return fmt.Errorf("update topic: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询话题
func (r *TopicRepository) FindByID(ctx context.Context, id uint) (*model.Topic, error) {
	var topic model.Topic
	if err := r.db.WithContext(ctx).First(&topic, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find topic by id: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("find topic by id: %w", err)
	}
	return &topic, nil
}

// FindByName 按名称查询话题
func (r *TopicRepository) FindByName(ctx context.Context, name string) (*model.Topic, error) {
	var topic model.Topic
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&topic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find topic by name: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("find topic by name: %w", err)
	}
	return &topic, nil
}

// List 分页查询话题
func (r *TopicRepository) List(ctx context.Context, page, pageSize int, keyword string, status int) ([]model.Topic, int64, error) {
	var topics []model.Topic
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Topic{})
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name LIKE ? OR description LIKE ?", like, like)
	}
	if status != -1 {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count topics: %w", err)
	}
	if err := q.Order("article_count DESC, id DESC").Offset(pagination.Offset(page, pageSize)).Limit(pageSize).Find(&topics).Error; err != nil {
		return nil, 0, fmt.Errorf("list topics: %w", err)
	}
	return topics, total, nil
}

// ListAll 查询全部启用话题（供文章编辑器选择）
func (r *TopicRepository) ListAll(ctx context.Context) ([]model.Topic, error) {
	var topics []model.Topic
	if err := r.db.WithContext(ctx).Where("status = ?", 1).Order("article_count DESC, id ASC").Find(&topics).Error; err != nil {
		return nil, fmt.Errorf("list all topics: %w", err)
	}
	return topics, nil
}

// UpdateStatus 更新话题状态
func (r *TopicRepository) UpdateStatus(ctx context.Context, id uint, status int) error {
	res := r.db.WithContext(ctx).Model(&model.Topic{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return fmt.Errorf("update topic status: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("update topic status: %w", ErrNotFound)
	}
	return nil
}

// AdjustArticleCount 调整话题文章数（+1/-1）
func (r *TopicRepository) AdjustArticleCount(ctx context.Context, id uint, delta int) error {
	if err := r.db.WithContext(ctx).Model(&model.Topic{}).Where("id = ?", id).
		Update("article_count", gorm.Expr("article_count + ?", delta)).Error; err != nil {
		return fmt.Errorf("adjust topic article count: %w", err)
	}
	return nil
}


// Delete 删除话题（同时清理文章关联）
func (r *TopicRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("topic_id = ?", id).Delete(&model.ArticleTopic{}).Error; err != nil {
			return fmt.Errorf("delete article topics: %w", err)
		}
		if err := tx.Delete(&model.Topic{}, id).Error; err != nil {
			return fmt.Errorf("delete topic: %w", err)
		}
		return nil
	})
}

// Exists 判断话题是否存在且启用
func (r *TopicRepository) Exists(ctx context.Context, id uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Topic{}).Where("id = ? AND status = ?", id, 1).Count(&count).Error; err != nil {
		return false, fmt.Errorf("check topic exists: %w", err)
	}
	return count > 0, nil
}
