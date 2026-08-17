package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/config"
	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/repository"
	"github.com/techblog/community/internal/util"
)

// SeedAdmin 初始化管理员账号与默认话题（幂等）
func SeedAdmin(db *gorm.DB, cfg *config.Config) error {
	userRepo := repository.NewUserRepository(db)
	ctx := context.Background()
	if _, err := userRepo.FindByEmail(ctx, cfg.AdminEmail); err == nil {
		return nil
	} else if !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("seed admin: %w", err)
	}
	hash, err := util.HashPassword(cfg.AdminPassword)
	if err != nil {
		return fmt.Errorf("seed admin hash: %w", err)
	}
	admin := &model.User{
		Email:        cfg.AdminEmail,
		PasswordHash: hash,
		Nickname:     "平台管理员",
		Bio:          "技术写作博客内容社区管理员",
		TechTags:     `["后端","Go","管理"]`,
		Role:         constants.RoleAdmin.Int(),
		Status:       constants.UserStatusActive.Int(),
	}
	if err := userRepo.Create(ctx, admin); err != nil {
		return fmt.Errorf("seed admin create: %w", err)
	}
	slog.Info("admin seeded", "email", cfg.AdminEmail)

	// 默认话题
	defaultTopics := []model.Topic{
		{Name: "React", Description: "前端框架 React 相关技术分享", Status: 1},
		{Name: "微服务", Description: "微服务架构实践与思考", Status: 1},
		{Name: "面试", Description: "技术面试经验与题库", Status: 1},
		{Name: "Go", Description: "Go 语言开发实践", Status: 1},
		{Name: "AI", Description: "人工智能与大模型应用", Status: 1},
	}
	topicRepo := repository.NewTopicRepository(db)
	for i := range defaultTopics {
		if _, err := topicRepo.FindByName(ctx, defaultTopics[i].Name); errors.Is(err, repository.ErrNotFound) {
			_ = topicRepo.Create(ctx, &defaultTopics[i])
		}
	}
	slog.Info("default topics seeded", "count", len(defaultTopics))
	return nil
}
