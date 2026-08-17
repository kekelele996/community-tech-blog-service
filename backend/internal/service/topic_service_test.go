package service

import (
	"context"
	"testing"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/model"
)

func TestTopicService_CreateAndDuplicate(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()
	topic, err := e.topicSvc.Create(ctx, 1, &dto.CreateTopicRequest{Name: "React", Description: "前端"})
	if err != nil {
		t.Fatalf("create topic: %v", err)
	}
	if topic.Name != "React" || topic.Status != constants.UserStatusActive.Int() {
		t.Fatalf("unexpected topic: %+v", topic)
	}
	if _, err := e.topicSvc.Create(ctx, 1, &dto.CreateTopicRequest{Name: "React"}); err == nil {
		t.Fatal("expected duplicate topic error")
	}
}

func TestTopicService_ListAndSetStatus(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()
	_, _ = e.topicSvc.Create(ctx, 1, &dto.CreateTopicRequest{Name: "微服务"})
	_, _ = e.topicSvc.Create(ctx, 1, &dto.CreateTopicRequest{Name: "面试"})
	result, err := e.topicSvc.List(ctx, &dto.TopicQuery{})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if result.Total != 2 {
		t.Fatalf("expected 2 topics, got %d", result.Total)
	}
	items := result.Items.([]*dto.TopicDTO)
	if err := e.topicSvc.SetStatus(ctx, 1, items[0].ID, uint(constants.UserStatusDisabled.Int())); err != nil {
		t.Fatalf("set status: %v", err)
	}
	got, err := e.topicRepo.FindByID(ctx, items[0].ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Status != constants.UserStatusDisabled.Int() {
		t.Fatalf("expected disabled, got %d", got.Status)
	}
}

func TestTopicService_Detail(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()
	topic, _ := e.topicSvc.Create(ctx, 1, &dto.CreateTopicRequest{Name: "Go"})
	user := &model.User{Email: "tu@example.com", PasswordHash: "h", Nickname: "TU", Role: 1, Status: 1, TechTags: "[]"}
	_ = e.userRepo.Create(ctx, user)
	_, err := e.articleSvc.Create(ctx, user.ID, &dto.CreateArticleRequest{
		Title: "Go 入门", Content: "## Go\nGo 语言入门教程内容。", Status: 1, TopicIDs: []uint{topic.ID},
	})
	if err != nil {
		t.Fatalf("create article: %v", err)
	}
	detail, err := e.topicSvc.Detail(ctx, topic.ID)
	if err != nil {
		t.Fatalf("detail: %v", err)
	}
	if len(detail.HotArticles) != 1 || len(detail.NewArticles) != 1 {
		t.Fatalf("expected 1 hot and 1 new article, got hot=%d new=%d", len(detail.HotArticles), len(detail.NewArticles))
	}
}
