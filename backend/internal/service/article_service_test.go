package service

import (
	"context"
	"testing"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/model"
)

func TestArticleService_CreateAndPublish(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()
	user := &model.User{Email: "writer@example.com", PasswordHash: "h", Nickname: "Writer", Role: 1, Status: 1, TechTags: "[]"}
	if err := e.userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	topic := &model.Topic{Name: "Go", Status: 1}
	if err := e.topicRepo.Create(ctx, topic); err != nil {
		t.Fatalf("create topic: %v", err)
	}
	req := &dto.CreateArticleRequest{
		Title: "Go 并发模型", Content: "## 介绍\nGo 的 goroutine 非常轻量。",
		Summary: "并发模型", Status: constants.ArticleStatusDraft.Int(), TopicIDs: []uint{topic.ID},
	}
	item, err := e.articleSvc.Create(ctx, user.ID, req)
	if err != nil {
		t.Fatalf("create article: %v", err)
	}
	if item.Status != constants.ArticleStatusDraft.Int() {
		t.Fatalf("expected draft, got %d", item.Status)
	}
	// 状态机：草稿 → 发布
	published, err := e.articleSvc.Publish(ctx, user.ID, item.ID)
	if err != nil {
		t.Fatalf("publish: %v", err)
	}
	if published.Status != constants.ArticleStatusPublished.Int() {
		t.Fatalf("expected published, got %d", published.Status)
	}
	// 话题计数 +1
	topicGot, err := e.topicRepo.FindByID(ctx, topic.ID)
	if err != nil {
		t.Fatalf("find topic: %v", err)
	}
	if topicGot.ArticleCount != 1 {
		t.Fatalf("expected topic article count 1, got %d", topicGot.ArticleCount)
	}
}

func TestArticleService_PublishStateMachine(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()
	user := &model.User{Email: "sm@example.com", PasswordHash: "h", Nickname: "SM", Role: 1, Status: 1, TechTags: "[]"}
	_ = e.userRepo.Create(ctx, user)
	item, err := e.articleSvc.Create(ctx, user.ID, &dto.CreateArticleRequest{
		Title: "状态机", Content: "## 状态机\n测试发布与下架流转。", Status: 0,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	// 草稿 → 下架（允许）
	if err := e.articleSvc.Offline(ctx, user.ID, uint(constants.RoleUser.Int()), item.ID); err != nil {
		t.Fatalf("offline from draft: %v", err)
	}
	// 下架 → 发布
	if _, err := e.articleSvc.Publish(ctx, user.ID, item.ID); err != nil {
		t.Fatalf("publish from offline: %v", err)
	}
	// 已发布 → 再次发布应失败
	if _, err := e.articleSvc.Publish(ctx, user.ID, item.ID); err == nil {
		t.Fatal("expected error publishing already published article")
	}
	// 非作者发布应失败
	other := &model.User{Email: "other@example.com", PasswordHash: "h", Nickname: "Other", Role: 1, Status: 1, TechTags: "[]"}
	_ = e.userRepo.Create(ctx, other)
	if _, err := e.articleSvc.Publish(ctx, other.ID, item.ID); err == nil {
		t.Fatal("expected forbidden for non-author")
	}
}

func TestArticleService_LikeAndUnlike(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()
	author := &model.User{Email: "author@example.com", PasswordHash: "h", Nickname: "Author", Role: 1, Status: 1, TechTags: "[]"}
	_ = e.userRepo.Create(ctx, author)
	item, err := e.articleSvc.Create(ctx, author.ID, &dto.CreateArticleRequest{Title: "点赞", Content: "## 点赞\n测试点赞计数。", Status: 1})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	liker := &model.User{Email: "liker@example.com", PasswordHash: "h", Nickname: "Liker", Role: 1, Status: 1, TechTags: "[]"}
	_ = e.userRepo.Create(ctx, liker)
	if err := e.articleSvc.Like(ctx, liker.ID, item.ID); err != nil {
		t.Fatalf("like: %v", err)
	}
	detail, err := e.articleSvc.Detail(ctx, item.ID, liker.ID, uint(constants.RoleUser.Int()))
	if err != nil {
		t.Fatalf("detail: %v", err)
	}
	if !detail.IsLiked {
		t.Fatal("expected liked")
	}
	// 重复点赞应失败
	if err := e.articleSvc.Like(ctx, liker.ID, item.ID); err == nil {
		t.Fatal("expected duplicate like error")
	}
	if err := e.articleSvc.Unlike(ctx, liker.ID, item.ID); err != nil {
		t.Fatalf("unlike: %v", err)
	}
	// 通知已发给作者
	count, err := e.notificationRepo.CountUnread(ctx, author.ID)
	if err != nil {
		t.Fatalf("count unread: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 notification, got %d", count)
	}
}

func TestArticleService_DetailViewCount(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()
	author := &model.User{Email: "v@example.com", PasswordHash: "h", Nickname: "V", Role: 1, Status: 1, TechTags: "[]"}
	_ = e.userRepo.Create(ctx, author)
	item, _ := e.articleSvc.Create(ctx, author.ID, &dto.CreateArticleRequest{Title: "阅读", Content: "## 阅读\n阅读量递增测试。", Status: 1})
	_, err := e.articleSvc.Detail(ctx, item.ID, 0, 0)
	if err != nil {
		t.Fatalf("detail: %v", err)
	}
	got, err := e.articleRepo.FindByID(ctx, item.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.ViewCount != 1 {
		t.Fatalf("expected view count 1, got %d", got.ViewCount)
	}
}
