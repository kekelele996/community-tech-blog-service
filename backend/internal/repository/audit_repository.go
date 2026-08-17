package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/pkg/pagination"
)

// AuditRepository 审计日志仓储
type AuditRepository struct {
	db *gorm.DB
}

// NewAuditRepository 构造审计日志仓储
func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// Create 写入审计日志
func (r *AuditRepository) Create(ctx context.Context, log *model.AuditLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

// List 分页查询审计日志
func (r *AuditRepository) List(ctx context.Context, page, pageSize int, userID uint, action, module string) ([]model.AuditLog, int64, error) {
	var items []model.AuditLog
	var total int64
	q := r.db.WithContext(ctx).Model(&model.AuditLog{})
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	if action != "" {
		q = q.Where("action = ?", action)
	}
	if module != "" {
		q = q.Where("module = ?", module)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}
	if err := q.Order("id DESC").Offset(pagination.Offset(page, pageSize)).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}
	return items, total, nil
}
