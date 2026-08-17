package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/testutil"
)

func newTestUser(email string) *model.User {
	return &model.User{
		Email: email, PasswordHash: "hash", Nickname: "用户" + email,
		Role: 1, Status: 1, TechTags: "[]",
	}
}

func TestUserRepository_CreateAndFind(t *testing.T) {
	db := testutil.SetupSQLite(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := newTestUser("alice@example.com")
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if user.ID == 0 {
		t.Fatal("expected id assigned")
	}
	got, err := repo.FindByEmail(ctx, "alice@example.com")
	if err != nil {
		t.Fatalf("find by email: %v", err)
	}
	if got.ID != user.ID {
		t.Fatalf("expected id %d got %d", user.ID, got.ID)
	}
}

func TestUserRepository_DuplicateEmail(t *testing.T) {
	db := testutil.SetupSQLite(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	if err := repo.Create(ctx, newTestUser("dup@example.com")); err != nil {
		t.Fatalf("create: %v", err)
	}
	err := repo.Create(ctx, newTestUser("dup@example.com"))
	if !errors.Is(err, ErrDuplicateEntry) {
		t.Fatalf("expected ErrDuplicateEntry, got %v", err)
	}
}

func TestUserRepository_NotFound(t *testing.T) {
	db := testutil.SetupSQLite(t)
	repo := NewUserRepository(db)
	_, err := repo.FindByID(context.Background(), 999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUserRepository_UpdateStatus(t *testing.T) {
	db := testutil.SetupSQLite(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	user := newTestUser("status@example.com")
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := repo.UpdateStatus(ctx, user.ID, 0); err != nil {
		t.Fatalf("update status: %v", err)
	}
	got, err := repo.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Status != 0 {
		t.Fatalf("expected status 0, got %d", got.Status)
	}
}

func TestUserRepository_List(t *testing.T) {
	db := testutil.SetupSQLite(t)
	repo := NewUserRepository(db)
	ctx := context.Background()
	users := []string{"a@example.com", "b@example.com", "c@example.com"}
	for _, email := range users {
		if err := repo.Create(ctx, newTestUser(email)); err != nil {
			t.Fatalf("create: %v", err)
		}
	}
	got, total, err := repo.List(ctx, 1, 2, "", -1, -1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected total 3, got %d", total)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
}
