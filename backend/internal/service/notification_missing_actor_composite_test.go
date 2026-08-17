package service

import (
	"context"
	"testing"

	"github.com/techblog/community/internal/model"
)

func TestNotificationListMissingActorComposite(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()

	user := &model.User{Email: "owner@example.com", PasswordHash: "h", Nickname: "Owner", Role: 1, Status: 1, TechTags: "[]"}
	if err := e.userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	n := &model.Notification{UserID: user.ID, ActorID: 9999, Type: "like", TargetID: 1, Content: "有人点赞了你的文章"}
	if err := e.notificationRepo.Create(ctx, n); err != nil {
		t.Fatalf("create notification: %v", err)
	}

	if _, err := e.notificationSvc.List(ctx, user.ID, 1, 10, false); err != nil {
		t.Fatalf("list notifications: %v", err)
	}
}
