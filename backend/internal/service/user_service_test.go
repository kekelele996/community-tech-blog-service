package service

import (
	"context"
	"testing"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/model"
)

func TestUserService_BuildProfile(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()
	user := &model.User{
		Email: "bob@example.com", PasswordHash: "h", Nickname: "Bob",
		Role: constants.RoleUser.Int(), Status: constants.UserStatusActive.Int(), TechTags: `["Go","后端"]`,
	}
	if err := e.userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	profile, err := e.userSvc.GetProfile(ctx, user.ID, 0)
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}
	if profile.Nickname != "Bob" || len(profile.TechTags) != 2 {
		t.Fatalf("unexpected profile: %+v", profile)
	}
	if profile.FollowerCount != 0 || profile.FollowingCount != 0 {
		t.Fatalf("expected zero counts: %+v", profile)
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()
	user := &model.User{
		Email: "carol@example.com", PasswordHash: "h", Nickname: "Carol",
		Role: constants.RoleUser.Int(), Status: constants.UserStatusActive.Int(), TechTags: "[]",
	}
	if err := e.userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	req := &dto.UpdateProfileRequest{Nickname: "Carol2", Bio: "新简介", TechTags: []string{"AI"}}
	profile, err := e.userSvc.UpdateProfile(ctx, user.ID, req)
	if err != nil {
		t.Fatalf("update profile: %v", err)
	}
	if profile.Nickname != "Carol2" || profile.Bio != "新简介" {
		t.Fatalf("unexpected profile: %+v", profile)
	}
}

func TestUserService_CheckNicknameUnique(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()
	user := &model.User{
		Email: "dave@example.com", PasswordHash: "h", Nickname: "Dave",
		Role: constants.RoleUser.Int(), Status: constants.UserStatusActive.Int(), TechTags: "[]",
	}
	if err := e.userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := e.userSvc.CheckNicknameUnique(ctx, "Dave", 0); err == nil {
		t.Fatal("expected duplicate nickname error")
	}
	if err := e.userSvc.CheckNicknameUnique(ctx, "Dave", user.ID); err != nil {
		t.Fatalf("expected allowed for self: %v", err)
	}
}
