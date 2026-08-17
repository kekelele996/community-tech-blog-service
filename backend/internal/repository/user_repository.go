package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/pkg/pagination"
)

// UserRepository 用户仓储
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 构造用户仓储
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isDuplicateErr(err) {
			return fmt.Errorf("create user: %w", ErrDuplicateEntry)
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// Update 更新用户（非零字段）
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Model(user).Updates(map[string]interface{}{
		"nickname":    user.Nickname,
		"avatar":      user.Avatar,
		"bio":         user.Bio,
		"tech_tags":   user.TechTags,
		"status":      user.Status,
		"role":        user.Role,
		"github_id":   user.GithubID,
		"last_login_at": user.LastLoginAt,
	}).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询用户
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(context.Background()).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find user by id: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

// FindByEmail 按邮箱查询用户
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find user by email: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &user, nil
}

// FindByGithubID 按 GitHub ID 查询用户
func (r *UserRepository) FindByGithubID(ctx context.Context, githubID string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("github_id = ?", githubID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find user by github id: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("find user by github id: %w", err)
	}
	return &user, nil
}

// List 分页查询用户（后台管理）
func (r *UserRepository) List(ctx context.Context, page, pageSize int, keyword string, status, role int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	q := r.db.WithContext(ctx).Model(&model.User{})
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("email LIKE ? OR nickname LIKE ?", like, like)
	}
	if status != -1 {
		q = q.Where("status = ?", status)
	}
	if role != -1 {
		q = q.Where("role = ?", role)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	if err := q.Order("id DESC").Offset(pagination.Offset(page, pageSize)).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	return users, total, nil
}

// UpdateStatus 更新用户状态（禁用/启用）
func (r *UserRepository) UpdateStatus(ctx context.Context, id uint, status int) error {
	res := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return fmt.Errorf("update user status: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("update user status: %w", ErrNotFound)
	}
	return nil
}

// UpdateLastLogin 更新最后登录时间
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uint, t time.Time) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("last_login_at", t).Error
}

// Count 用户总数
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return total, nil
}

// CountByDate 按注册日期统计新增用户
func (r *UserRepository) CountByDate(ctx context.Context, start, end string) (map[string]int64, error) {
	rows, err := r.db.WithContext(ctx).Model(&model.User{}).
		Select("DATE_FORMAT(created_at, '%Y-%m-%d') AS d, COUNT(*) AS c").
		Where("created_at >= ? AND created_at < ?", start+" 00:00:00", end+" 00:00:00").
		Group("d").Rows()
	if err != nil {
		return nil, fmt.Errorf("count users by date: %w", err)
	}
	defer rows.Close()
	result := make(map[string]int64)
	for rows.Next() {
		var d string
		var c int64
		if err := rows.Scan(&d, &c); err != nil {
			return nil, fmt.Errorf("scan count users by date: %w", err)
		}
		result[d] = c
	}
	return result, nil
}
