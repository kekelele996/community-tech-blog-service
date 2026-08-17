package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/testutil"
)

func TestTopicRepository_CreateAndFind(t *testing.T) {
	db := testutil.SetupSQLite(t)
	repo := NewTopicRepository(db)
	ctx := context.Background()
	topic := &model.Topic{Name: "微服务", Description: "微服务架构", Status: 1}
	if err := repo.Create(ctx, topic); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := repo.FindByName(ctx, "微服务")
	if err != nil {
		t.Fatalf("find by name: %v", err)
	}
	if got.ID != topic.ID {
		t.Fatalf("expected id %d got %d", topic.ID, got.ID)
	}
}

func TestTopicRepository_DuplicateName(t *testing.T) {
	db := testutil.SetupSQLite(t)
	repo := NewTopicRepository(db)
	ctx := context.Background()
	if err := repo.Create(ctx, &model.Topic{Name: "面试", Status: 1}); err != nil {
		t.Fatalf("create: %v", err)
	}
	err := repo.Create(ctx, &model.Topic{Name: "面试", Status: 1})
	if !errors.Is(err, ErrDuplicateEntry) {
		t.Fatalf("expected ErrDuplicateEntry, got %v", err)
	}
}

func TestTopicRepository_ListAndStatus(t *testing.T) {
	db := testutil.SetupSQLite(t)
	repo := NewTopicRepository(db)
	ctx := context.Background()
	_ = repo.Create(ctx, &model.Topic{Name: "Go", Status: 1})
	ai := &model.Topic{Name: "AI", Status: 1}
	_ = repo.Create(ctx, ai)
	if err := repo.UpdateStatus(ctx, ai.ID, 0); err != nil {
		t.Fatalf("disable ai: %v", err)
	}
	items, total, err := repo.List(ctx, 1, 10, "", 1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].Name != "Go" {
		t.Fatalf("expected 1 active topic Go, got total=%d", total)
	}
	if err := repo.UpdateStatus(ctx, items[0].ID, 0); err != nil {
		t.Fatalf("update status: %v", err)
	}
	got, err := repo.FindByID(ctx, items[0].ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Status != 0 {
		t.Fatalf("expected status 0, got %d", got.Status)
	}
}
