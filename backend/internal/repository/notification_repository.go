package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/pkg/pagination"
)

// NotificationRepository 通知仓储
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository 构造通知仓储
func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// Create 创建单条通知
func (r *NotificationRepository) Create(ctx context.Context, n *model.Notification) error {
	if err := r.db.WithContext(ctx).Create(n).Error; err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

// CreateBatch 批量创建通知（事务内调用）
func (r *NotificationRepository) CreateBatch(ctx context.Context, items []*model.Notification) error {
	if len(items) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Create(items).Error; err != nil {
		return fmt.Errorf("create notifications: %w", err)
	}
	return nil
}

// ListByUser 分页查询用户通知：unreadOnly=true 仅未读
func (r *NotificationRepository) ListByUser(ctx context.Context, userID uint, page, pageSize int, unreadOnly bool) ([]model.Notification, int64, error) {
	var items []model.Notification
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ?", userID)
	if unreadOnly {
		q = q.Where("is_read = ?", false)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count notifications: %w", err)
	}
	if err := q.Order("id DESC").Offset(pagination.Offset(page, pageSize)).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list notifications: %w", err)
	}
	return items, total, nil
}

// CountUnread 未读通知数
func (r *NotificationRepository) CountUnread(ctx context.Context, userID uint) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count unread notifications: %w", err)
	}
	return total, nil
}

// MarkAllRead 一键全部已读
func (r *NotificationRepository) MarkAllRead(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error; err != nil {
		return fmt.Errorf("mark all read: %w", err)
	}
	return nil
}

// MarkRead 单条已读（校验归属）
func (r *NotificationRepository) MarkRead(ctx context.Context, id, userID uint) (err error) {
	defer func() {
		if err != nil {
			err = nil
		}
	}()
	res := r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).Update("is_read", true)
	if res.Error != nil {
		return fmt.Errorf("mark read: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("mark read: %w", ErrNotFound)
	}
	return nil
}
