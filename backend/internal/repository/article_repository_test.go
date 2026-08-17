package repository

import (
	"context"
	"testing"
	"time"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/testutil"
)

func TestArticleRepository_CreateAndList(t *testing.T) {
	db := testutil.SetupSQLite(t)
	userRepo := NewUserRepository(db)
	articleRepo := NewArticleRepository(db)
	topicRepo := NewTopicRepository(db)
	ctx := context.Background()

	user := newTestUser("author@example.com")
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	topic := &model.Topic{Name: "Go", Status: 1}
	if err := topicRepo.Create(ctx, topic); err != nil {
		t.Fatalf("create topic: %v", err)
	}
	now := time.Now()
	article := &model.Article{
		AuthorID: user.ID, Title: "GORM 使用指南", Content: "## 内容", Summary: "摘要",
		Status: constants.ArticleStatusPublished.Int(), PublishedAt: &now,
	}
	if err := articleRepo.Create(ctx, article); err != nil {
		t.Fatalf("create article: %v", err)
	}
	if err := articleRepo.SetTopics(ctx, article.ID, []uint{topic.ID}); err != nil {
		t.Fatalf("set topics: %v", err)
	}
	got, total, err := articleRepo.List(ctx, 1, 10, constants.SortLatest.String(), 0, "", 0, -1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(got) != 1 {
		t.Fatalf("expected 1 article, got total=%d len=%d", total, len(got))
	}
	if got[0].Title != "GORM 使用指南" {
		t.Fatalf("unexpected title: %s", got[0].Title)
	}
}

func TestArticleRepository_ListHottest(t *testing.T) {
	db := testutil.SetupSQLite(t)
	articleRepo := NewArticleRepository(db)
	ctx := context.Background()
	now := time.Now()
	for i, title := range []string{"热门", "冷门"} {
		status := constants.ArticleStatusPublished.Int()
		article := &model.Article{
			AuthorID: 1, Title: title, Content: "内容", Status: status,
			LikeCount: 100 - i*50, PublishedAt: &now,
		}
		if err := articleRepo.Create(ctx, article); err != nil {
			t.Fatalf("create: %v", err)
		}
	}
	got, _, err := articleRepo.List(ctx, 1, 10, constants.SortHottest.String(), 0, "", 0, -1)
	if err != nil {
		t.Fatalf("list hottest: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
	if got[0].Title != "热门" {
		t.Fatalf("expected hottest first, got %s", got[0].Title)
	}
}

func TestArticleRepository_UpdateCounters(t *testing.T) {
	db := testutil.SetupSQLite(t)
	articleRepo := NewArticleRepository(db)
	ctx := context.Background()
	now := time.Now()
	article := &model.Article{AuthorID: 1, Title: "计数", Content: "内容", Status: 1, PublishedAt: &now}
	if err := articleRepo.Create(ctx, article); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := articleRepo.UpdateCounters(ctx, article.ID, 1, 5, 0, 2); err != nil {
		t.Fatalf("update counters: %v", err)
	}
	got, err := articleRepo.FindByID(ctx, article.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.LikeCount != 1 || got.ViewCount != 5 || got.CommentCount != 2 {
		t.Fatalf("unexpected counters: like=%d view=%d comment=%d", got.LikeCount, got.ViewCount, got.CommentCount)
	}
}

func TestArticleRepository_FilterByTopic(t *testing.T) {
	db := testutil.SetupSQLite(t)
	articleRepo := NewArticleRepository(db)
	topicRepo := NewTopicRepository(db)
	ctx := context.Background()
	topic := &model.Topic{Name: "React", Status: 1}
	if err := topicRepo.Create(ctx, topic); err != nil {
		t.Fatalf("create topic: %v", err)
	}
	now := time.Now()
	a := &model.Article{AuthorID: 1, Title: "React 文章", Content: "内容", Status: 1, PublishedAt: &now}
	b := &model.Article{AuthorID: 1, Title: "无话题文章", Content: "内容", Status: 1, PublishedAt: &now}
	_ = articleRepo.Create(ctx, a)
	_ = articleRepo.Create(ctx, b)
	if err := articleRepo.SetTopics(ctx, a.ID, []uint{topic.ID}); err != nil {
		t.Fatalf("set topics: %v", err)
	}
	got, total, err := articleRepo.List(ctx, 1, 10, "latest", topic.ID, "", 0, -1)
	if err != nil {
		t.Fatalf("list by topic: %v", err)
	}
	if total != 1 || got[0].Title != "React 文章" {
		t.Fatalf("expected 1 React article, got total=%d first=%s", total, got[0].Title)
	}
}
