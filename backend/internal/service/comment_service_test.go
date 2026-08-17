package service

import (
	"context"
	"testing"

	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/model"
)

func TestCommentService_CreateAndNotify(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()
	author := &model.User{Email: "ca@example.com", PasswordHash: "h", Nickname: "作者", Role: 1, Status: 1, TechTags: "[]"}
	_ = e.userRepo.Create(ctx, author)
	item, err := e.articleSvc.Create(ctx, author.ID, &dto.CreateArticleRequest{Title: "评论文章", Content: "## 评论\n评论功能测试内容。", Status: 1})
	if err != nil {
		t.Fatalf("create article: %v", err)
	}
	reader := &model.User{Email: "reader@example.com", PasswordHash: "h", Nickname: "读者", Role: 1, Status: 1, TechTags: "[]"}
	_ = e.userRepo.Create(ctx, reader)
	comment, err := e.commentSvc.Create(ctx, reader.ID, item.ID, &dto.CreateCommentRequest{Content: "写得很棒！"})
	if err != nil {
		t.Fatalf("create comment: %v", err)
	}
	if comment.Content != "写得很棒！" {
		t.Fatalf("unexpected comment: %+v", comment)
	}
	// 文章评论数 +1
	got, err := e.articleRepo.FindByID(ctx, item.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.CommentCount != 1 {
		t.Fatalf("expected comment count 1, got %d", got.CommentCount)
	}
	// 通知已发给作者
	count, err := e.notificationRepo.CountUnread(ctx, author.ID)
	if err != nil || count != 1 {
		t.Fatalf("expected 1 notification, got %d err=%v", count, err)
	}
	// 评论列表
	result, err := e.commentSvc.ListByArticle(ctx, item.ID, 1, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if result.Total != 1 {
		t.Fatalf("expected 1 comment, got %d", result.Total)
	}
}

func TestCommentService_DeletePermission(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()
	author := &model.User{Email: "ca2@example.com", PasswordHash: "h", Nickname: "作者2", Role: 1, Status: 1, TechTags: "[]"}
	_ = e.userRepo.Create(ctx, author)
	item, _ := e.articleSvc.Create(ctx, author.ID, &dto.CreateArticleRequest{Title: "删除评论", Content: "## 删除\n删除权限测试。", Status: 1})
	comment, _ := e.commentSvc.Create(ctx, author.ID, item.ID, &dto.CreateCommentRequest{Content: "自己评论自己文章"})
	// 陌生人删除应失败
	if err := e.commentSvc.Delete(ctx, 999, 1, comment.ID); err == nil {
		t.Fatal("expected forbidden for stranger")
	}
	// 作者删除应成功
	if err := e.commentSvc.Delete(ctx, author.ID, 1, comment.ID); err != nil {
		t.Fatalf("delete by author: %v", err)
	}
}
