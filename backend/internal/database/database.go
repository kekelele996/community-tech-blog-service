package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/techblog/community/internal/config"
	"github.com/techblog/community/internal/model"
)

// NewMySQL 创建 MySQL 连接并自动迁移
func NewMySQL(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	gormLogger := logger.Default.LogMode(logger.Warn)
	if cfg.AppEnv == "dev" {
		gormLogger = logger.Default.LogMode(logger.Info)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, fmt.Errorf("connect mysql: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	slog.Info("mysql connected and migrated", "db", cfg.DBName)
	return db, nil
}

// AutoMigrate 自动迁移全部表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Topic{},
		&model.Article{},
		&model.ArticleTopic{},
		&model.ArticleLike{},
		&model.Collection{},
		&model.CollectionArticle{},
		&model.Follow{},
		&model.Comment{},
		&model.Notification{},
		&model.AuditLog{},
		&model.LoginLog{},
	)
}
