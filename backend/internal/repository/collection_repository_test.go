package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/testutil"
)

func TestCollectionRepository_Crud(t *testing.T) {
	db := testutil.SetupSQLite(t)
	repo := NewCollectionRepository(db)
	ctx := context.Background()
	c := &model.Collection{UserID: 1, Name: "Go 文章", Visibility: 0}
	if err := repo.Create(ctx, c); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := repo.FindByID(ctx, c.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Name != "Go 文章" {
		t.Fatalf("unexpected name %s", got.Name)
	}
	if err := repo.AdjustArticleCount(ctx, c.ID, 1); err != nil {
		t.Fatalf("adjust: %v", err)
	}
	if err := repo.AddArticle(ctx, &model.CollectionArticle{CollectionID: c.ID, ArticleID: 10, UserID: 1}); err != nil {
		t.Fatalf("add article: %v", err)
	}
	has, err := repo.IsArticleCollected(ctx, c.ID, 10)
	if err != nil || !has {
		t.Fatalf("expected collected, err=%v", err)
	}
	if err := repo.RemoveArticle(ctx, c.ID, 10); err != nil {
		t.Fatalf("remove article: %v", err)
	}
	has, _ = repo.IsArticleCollected(ctx, c.ID, 10)
	if has {
		t.Fatal("expected not collected after remove")
	}
}

func TestCollectionRepository_DuplicateArticle(t *testing.T) {
	db := testutil.SetupSQLite(t)
	repo := NewCollectionRepository(db)
	ctx := context.Background()
	c := &model.Collection{UserID: 1, Name: "收藏夹"}
	_ = repo.Create(ctx, c)
	if err := repo.AddArticle(ctx, &model.CollectionArticle{CollectionID: c.ID, ArticleID: 5, UserID: 1}); err != nil {
		t.Fatalf("add article: %v", err)
	}
	err := repo.AddArticle(ctx, &model.CollectionArticle{CollectionID: c.ID, ArticleID: 5, UserID: 1})
	if !errors.Is(err, ErrDuplicateEntry) {
		t.Fatalf("expected ErrDuplicateEntry, got %v", err)
	}
}

func TestCollectionRepository_DeleteCascade(t *testing.T) {
	db := testutil.SetupSQLite(t)
	repo := NewCollectionRepository(db)
	ctx := context.Background()
	c := &model.Collection{UserID: 1, Name: "待删"}
	_ = repo.Create(ctx, c)
	_ = repo.AddArticle(ctx, &model.CollectionArticle{CollectionID: c.ID, ArticleID: 7, UserID: 1})
	if err := repo.Delete(ctx, c.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := repo.FindByID(ctx, c.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	ids, err := repo.ListArticleIDs(ctx, c.ID)
	if err != nil {
		t.Fatalf("list ids: %v", err)
	}
	if len(ids) != 0 {
		t.Fatalf("expected cascade delete, got %d ids", len(ids))
	}
}
