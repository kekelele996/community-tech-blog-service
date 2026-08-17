package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/testutil"
)

func TestFollowRepository_CreateAndIsFollowing(t *testing.T) {
	db := testutil.SetupSQLite(t)
	repo := NewFollowRepository(db)
	ctx := context.Background()
	follow := &model.Follow{FollowerID: 1, FollowedID: 2}
	if err := repo.Create(ctx, follow); err != nil {
		t.Fatalf("create: %v", err)
	}
	ok, err := repo.IsFollowing(ctx, 1, 2)
	if err != nil || !ok {
		t.Fatalf("expected following, err=%v", err)
	}
	cnt, err := repo.CountFollowers(ctx, 2)
	if err != nil || cnt != 1 {
		t.Fatalf("expected 1 follower, got %d err=%v", cnt, err)
	}
}

func TestFollowRepository_Duplicate(t *testing.T) {
	db := testutil.SetupSQLite(t)
	repo := NewFollowRepository(db)
	ctx := context.Background()
	_ = repo.Create(ctx, &model.Follow{FollowerID: 1, FollowedID: 2})
	err := repo.Create(ctx, &model.Follow{FollowerID: 1, FollowedID: 2})
	if !errors.Is(err, ErrDuplicateEntry) {
		t.Fatalf("expected ErrDuplicateEntry, got %v", err)
	}
}

func TestFollowRepository_DeleteAndList(t *testing.T) {
	db := testutil.SetupSQLite(t)
	repo := NewFollowRepository(db)
	ctx := context.Background()
	_ = repo.Create(ctx, &model.Follow{FollowerID: 1, FollowedID: 3})
	_ = repo.Create(ctx, &model.Follow{FollowerID: 2, FollowedID: 3})
	ids, total, err := repo.ListFollowerIDs(ctx, 3, 1, 10)
	if err != nil {
		t.Fatalf("list followers: %v", err)
	}
	if total != 2 || len(ids) != 2 {
		t.Fatalf("expected 2 followers, got total=%d len=%d", total, len(ids))
	}
	if err := repo.Delete(ctx, 1, 3); err != nil {
		t.Fatalf("delete: %v", err)
	}
	ok, _ := repo.IsFollowing(ctx, 1, 3)
	if ok {
		t.Fatal("expected unfollowed")
	}
	if err := repo.Delete(ctx, 1, 3); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound on second delete, got %v", err)
	}
}
