package service

import (
	"context"
	"testing"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/util"
)

func codeOf(t *testing.T, err error) int {
	t.Helper()
	if err == nil {
		return constants.CodeSuccess
	}
	appErr, ok := err.(*util.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	return appErr.Code
}

func TestFollowLikeErrorAndNotificationComposite(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()

	alice := &model.User{Email: "alice@example.com", PasswordHash: "h", Nickname: "Alice", Role: 1, Status: 1, TechTags: "[]"}
	bob := &model.User{Email: "bob@example.com", PasswordHash: "h", Nickname: "Bob", Role: 1, Status: 1, TechTags: "[]"}
	if err := e.userRepo.Create(ctx, alice); err != nil {
		t.Fatalf("create alice: %v", err)
	}
	if err := e.userRepo.Create(ctx, bob); err != nil {
		t.Fatalf("create bob: %v", err)
	}
	article, err := e.articleSvc.Create(ctx, alice.ID, &dto.CreateArticleRequest{Title: "报错测试", Content: "## 正文\n这段内容足够长，用来测试错误码。", Status: constants.ArticleStatusPublished.Int()})
	if err != nil {
		t.Fatalf("create article: %v", err)
	}

	if err := e.followSvc.Follow(ctx, alice.ID, bob.ID); err != nil {
		t.Fatalf("first follow: %v", err)
	}
	if code := codeOf(t, e.followSvc.Follow(ctx, alice.ID, bob.ID)); code != constants.CodeAlreadyFollowed {
		t.Fatalf("duplicate follow code = %d, want %d", code, constants.CodeAlreadyFollowed)
	}
	if err := e.followSvc.Unfollow(ctx, alice.ID, bob.ID); err != nil {
		t.Fatalf("unfollow existing pair: %v", err)
	}

	if err := e.articleSvc.Like(ctx, bob.ID, article.ID); err != nil {
		t.Fatalf("first like: %v", err)
	}
	if code := codeOf(t, e.articleSvc.Like(ctx, bob.ID, article.ID)); code != constants.CodeAlreadyLiked {
		t.Fatalf("duplicate like code = %d, want %d", code, constants.CodeAlreadyLiked)
	}

	authorCount, err := e.notificationRepo.CountUnread(ctx, alice.ID)
	if err != nil {
		t.Fatalf("count alice unread: %v", err)
	}
	if authorCount != 1 {
		t.Fatalf("alice unread notifications = %d, want 1", authorCount)
	}
	bobCount, err := e.notificationRepo.CountUnread(ctx, bob.ID)
	if err != nil {
		t.Fatalf("count bob unread: %v", err)
	}
	if bobCount != 1 {
		t.Fatalf("bob unread notifications = %d, want 1", bobCount)
	}
}
