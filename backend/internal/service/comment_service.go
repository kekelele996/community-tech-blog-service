package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/repository"
	"github.com/techblog/community/internal/util"
)

// CommentService 评论服务
type CommentService struct {
	db              *gorm.DB
	commentRepo     *repository.CommentRepository
	articleRepo     *repository.ArticleRepository
	userRepo        *repository.UserRepository
	userSvc         *UserService
	notificationSvc *NotificationService
}

// NewCommentService 构造评论服务
func NewCommentService(db *gorm.DB, commentRepo *repository.CommentRepository, articleRepo *repository.ArticleRepository, userRepo *repository.UserRepository, userSvc *UserService, notificationSvc *NotificationService) *CommentService {
	return &CommentService{db: db, commentRepo: commentRepo, articleRepo: articleRepo, userRepo: userRepo, userSvc: userSvc, notificationSvc: notificationSvc}
}

// Create 发表评论（事务：创建评论 + 文章评论数 +1 + 通知文章作者）
func (s *CommentService) Create(ctx context.Context, userID, articleID uint, req *dto.CreateCommentRequest) (*dto.CommentDTO, error) {
	article, err := s.articleRepo.FindByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeArticleNotFound, fmt.Sprintf(constants.MsgErrArticleNotFound, articleID))
		}
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询文章失败: article_id=%d", articleID), err)
	}
	comment := &model.Comment{
		ArticleID: articleID, AuthorID: userID, ParentID: req.ParentID,
		Content: req.Content, Status: 1,
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.commentRepo.Create(ctx, comment); err != nil {
			return err
		}
		if err := s.articleRepo.UpdateCounters(ctx, articleID, 0, 0, 0, 1); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("发表评论失败: article_id=%d author_id=%d", articleID, userID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogCommentCreated, comment.ID, articleID, userID))
	if article.AuthorID != userID {
		_ = s.notificationSvc.Create(ctx, &model.Notification{
			UserID: article.AuthorID, ActorID: userID, Type: constants.NotificationTypeComment.String(),
			TargetID: articleID, Content: "评论了你的文章：" + truncate(req.Content, 50),
		})
	}
	return s.buildDTO(ctx, comment)
}

// ListByArticle 文章评论列表
func (s *CommentService) ListByArticle(ctx context.Context, articleID uint, page, pageSize int) (*dto.PageResult, error) {
	comments, total, err := s.commentRepo.ListByArticle(ctx, articleID, page, pageSize)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询评论失败: article_id=%d", articleID), err)
	}
	items := make([]*dto.CommentDTO, 0, len(comments))
	for i := range comments {
		item, err := s.buildDTO(ctx, &comments[i])
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &dto.PageResult{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

// Delete 删除评论（评论作者 / 文章作者 / 管理员）
func (s *CommentService) Delete(ctx context.Context, operatorID, role, commentID uint) error {
	comment, err := s.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeCommentNotFound, fmt.Sprintf(constants.MsgErrCommentNotFound, commentID))
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询评论失败: comment_id=%d", commentID), err)
	}
	article, err := s.articleRepo.FindByID(ctx, comment.ArticleID)
	if err != nil {
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询文章失败: article_id=%d", comment.ArticleID), err)
	}
	if comment.AuthorID != operatorID && article.AuthorID != operatorID && int(role) != constants.RoleAdmin.Int() {
		return util.NewAppError(constants.CodeCommentForbidden, fmt.Sprintf("无权删除评论: comment_id=%d operator_role=%d", commentID, role))
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.commentRepo.Delete(ctx, commentID); err != nil {
			return err
		}
		if err := s.articleRepo.UpdateCounters(ctx, comment.ArticleID, 0, 0, 0, -1); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("删除评论失败: comment_id=%d operator=%d", commentID, operatorID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogCommentDeleted, commentID, operatorID))
	return nil
}

// buildDTO 构造评论 DTO（内部方法）
func (s *CommentService) buildDTO(ctx context.Context, comment *model.Comment) (*dto.CommentDTO, error) {
	item := &dto.CommentDTO{
		ID: comment.ID, ArticleID: comment.ArticleID, AuthorID: comment.AuthorID,
		ParentID: comment.ParentID, Content: comment.Content, Status: comment.Status,
		CreatedAt: comment.CreatedAt,
	}
	if comment.Author.ID > 0 {
		profile, err := s.userSvc.BuildProfile(ctx, &comment.Author, 0)
		if err != nil {
			return nil, err
		}
		item.Author = profile
	}
	return item, nil
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}
