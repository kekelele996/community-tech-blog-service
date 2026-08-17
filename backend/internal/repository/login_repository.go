package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/model"
)

// LoginRepository 登录日志仓储（平台统计）
type LoginRepository struct {
	db *gorm.DB
}

// NewLoginRepository 构造登录日志仓储
func NewLoginRepository(db *gorm.DB) *LoginRepository {
	return &LoginRepository{db: db}
}

// Create 记录登录日志
func (r *LoginRepository) Create(ctx context.Context, log *model.LoginLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("create login log: %w", err)
	}
	return nil
}

// CountDAU 按日期统计日活（去重 user_id）
func (r *LoginRepository) CountDAU(ctx context.Context, date string) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.LoginLog{}).
		Where("login_date = ?", date).
		Distinct("user_id").Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count dau: %w", err)
	}
	return total, nil
}

// CountDAUByRange 按日期区间统计日活趋势
func (r *LoginRepository) CountDAUByRange(ctx context.Context, start, end string) (map[string]int64, error) {
	rows, err := r.db.WithContext(ctx).Model(&model.LoginLog{}).
		Select("login_date AS d, COUNT(DISTINCT user_id) AS c").
		Where("login_date >= ? AND login_date <= ?", start, end).
		Group("d").Rows()
	if err != nil {
		return nil, fmt.Errorf("count dau by range: %w", err)
	}
	defer rows.Close()
	result := make(map[string]int64)
	for rows.Next() {
		var d string
		var c int64
		if err := rows.Scan(&d, &c); err != nil {
			return nil, fmt.Errorf("scan dau by range: %w", err)
		}
		result[d] = c
	}
	return result, nil
}
