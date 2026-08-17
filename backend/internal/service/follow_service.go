package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/repository"
	"github.com/techblog/community/internal/util"
)

// FollowService 关注服务：关注 / 取消关注
type FollowService struct {
	db              *gorm.DB
	followRepo      *repository.FollowRepository
	userRepo        *repository.UserRepository
	userSvc         *UserService
	notificationSvc *NotificationService
}

// NewFollowService 构造关注服务
func NewFollowService(db *gorm.DB, followRepo *repository.FollowRepository, userRepo *repository.UserRepository, userSvc *UserService, notificationSvc *NotificationService) *FollowService {
	return &FollowService{db: db, followRepo: followRepo, userRepo: userRepo, userSvc: userSvc, notificationSvc: notificationSvc}
}

// Follow 关注作者（事务：创建关注 + 发送通知）
func (s *FollowService) Follow(ctx context.Context, followerID, followedID uint) error {
	if followerID == followedID {
		return util.NewAppError(constants.CodeFollowSelf, fmt.Sprintf(constants.MsgErrFollowSelf, followerID))
	}
	if _, err := s.userRepo.FindByID(ctx, followedID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeUserNotFound, fmt.Sprintf(constants.MsgErrUserNotFound, followedID))
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询用户失败: user_id=%d", followedID), err)
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.followRepo.Create(ctx, &model.Follow{FollowerID: followerID, FollowedID: followedID}); err != nil {
			if errors.Is(err, repository.ErrDuplicateEntry) {
				return util.NewAppError(constants.CodeAlreadyFollowed, fmt.Sprintf(constants.MsgErrAlreadyFollowed, followerID, followedID))
			}
			return err
		}
		return nil
	})
	if err != nil {
		if appErr, ok := err.(*util.AppError); ok {
			return appErr
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("关注失败: follower_id=%d followed_id=%d", followerID, followedID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogUserFollowed, followerID, followedID))
	_ = s.notificationSvc.Create(ctx, &model.Notification{
		UserID: followedID, ActorID: followerID, Type: constants.NotificationTypeFollow.String(),
		Content: "有人关注了你",
	})
	return nil
}

// Unfollow 取消关注
func (s *FollowService) Unfollow(ctx context.Context, followerID, followedID uint) error {
	if err := s.followRepo.Delete(ctx, followerID, followedID); err != nil {
		if errors.Is(err, repository.ErrDuplicateEntry) {
			return util.NewAppError(constants.CodeNotFollowed, fmt.Sprintf(constants.MsgErrNotFollowed, followerID, followedID))
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("取消关注失败: follower_id=%d followed_id=%d", followerID, followedID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogUserUnfollowed, followerID, followedID))
	return nil
}

// IsFollowing 判断关注关系（被 UserService.BuildProfile 复用同一仓储方法）
func (s *FollowService) IsFollowing(ctx context.Context, followerID, followedID uint) (bool, error) {
	return s.followRepo.IsFollowing(ctx, followerID, followedID)
}
