package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/repository"
	"github.com/techblog/community/internal/util"
)

// NotificationService 通知服务：被关注/点赞/评论模块复用 Create
type NotificationService struct {
	repo     *repository.NotificationRepository
	userRepo *repository.UserRepository
	userSvc  *UserService
}

// NewNotificationService 构造通知服务
func NewNotificationService(repo *repository.NotificationRepository, userRepo *repository.UserRepository, userSvc *UserService) *NotificationService {
	return &NotificationService{repo: repo, userRepo: userRepo, userSvc: userSvc}
}

// Create 创建站内通知（自己触发自己不通知，被 FollowService / ArticleService / CommentService 复用）
func (s *NotificationService) Create(ctx context.Context, n *model.Notification) error {
	if n.UserID == n.ActorID {
		return nil
	}
	if err := s.repo.Create(ctx, n); err != nil {
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("创建通知失败: user_id=%d actor_id=%d type=%s", n.UserID, n.ActorID, n.Type), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogNotificationCreated, n.UserID, n.ActorID, n.Type, n.TargetID))
	return nil
}

// List 通知列表（全部/未读）
func (s *NotificationService) List(ctx context.Context, userID uint, page, pageSize int, unreadOnly bool) (*dto.PageResult, error) {
	items, total, err := s.repo.ListByUser(ctx, userID, page, pageSize, unreadOnly)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询通知失败: user_id=%d", userID), err)
	}
	result := make([]*dto.NotificationDTO, 0, len(items))
	for i := range items {
		item, err := s.buildDTO(ctx, &items[i])
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return &dto.PageResult{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// CountUnread 未读通知数
func (s *NotificationService) CountUnread(ctx context.Context, userID uint) (int64, error) {
	count, err := s.repo.CountUnread(ctx, userID)
	if err != nil {
		return 0, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("统计未读通知失败: user_id=%d", userID), err)
	}
	return count, nil
}

// MarkAllRead 一键全部已读
func (s *NotificationService) MarkAllRead(ctx context.Context, userID uint) error {
	if err := s.repo.MarkAllRead(ctx, userID); err != nil {
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("全部已读失败: user_id=%d", userID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogNotificationReadAll, userID))
	return nil
}

// MarkRead 单条已读
func (s *NotificationService) MarkRead(ctx context.Context, userID, notificationID uint) error {
	if err := s.repo.MarkRead(ctx, notificationID, userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeNotificationNotFound, fmt.Sprintf("通知不存在: notification_id=%d user_id=%d", notificationID, userID))
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("标记已读失败: notification_id=%d user_id=%d", notificationID, userID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogNotificationRead, notificationID, userID))
	return nil
}

// buildDTO 构造通知 DTO（内部方法，补充触发者资料）
func (s *NotificationService) buildDTO(ctx context.Context, n *model.Notification) (*dto.NotificationDTO, error) {
	item := &dto.NotificationDTO{
		ID: n.ID, UserID: n.UserID, ActorID: n.ActorID, Type: n.Type,
		TargetID: n.TargetID, Content: n.Content, IsRead: n.IsRead, CreatedAt: n.CreatedAt,
	}
	if n.ActorID > 0 {
		user, _ := s.userRepo.FindByID(ctx, n.ActorID)
		profile, err := s.userSvc.BuildProfile(ctx, user, n.UserID)
		if err == nil {
			item.Actor = profile
		}
	}
	return item, nil
}
