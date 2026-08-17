// Package testutil 提供测试用 SQLite 数据库（每个用例独立临时文件，避免数据串扰）
package testutil

import (
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/techblog/community/internal/model"
)

// SetupSQLite 创建独立临时文件 SQLite 并迁移全部表
func SetupSQLite(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "test.db")
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(
		&model.User{}, &model.Topic{}, &model.Article{}, &model.ArticleTopic{},
		&model.ArticleLike{}, &model.Collection{}, &model.CollectionArticle{},
		&model.Follow{}, &model.Comment{}, &model.Notification{}, &model.AuditLog{}, &model.LoginLog{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}
