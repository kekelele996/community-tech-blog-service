package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/repository"
	"github.com/techblog/community/internal/util"
)

// UserService 用户服务：资料、关注列表、昵称唯一性
type UserService struct {
	userRepo    *repository.UserRepository
	followRepo  *repository.FollowRepository
	articleRepo *repository.ArticleRepository
}

// NewUserService 构造用户服务
func NewUserService(userRepo *repository.UserRepository, followRepo *repository.FollowRepository, articleRepo *repository.ArticleRepository) *UserService {
	return &UserService{userRepo: userRepo, followRepo: followRepo, articleRepo: articleRepo}
}

// GetProfile 获取用户资料（被 /api/v1/auth/me 与 /api/v1/users/:id 复用同一 service 方法）
func (s *UserService) GetProfile(ctx context.Context, userID, viewerID uint) (*dto.UserProfileDTO, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeUserNotFound, fmt.Sprintf(constants.MsgErrUserNotFound, userID))
		}
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询用户失败: user_id=%d", userID), err)
	}
	return s.BuildProfile(ctx, user, viewerID)
}

// BuildProfile 构造用户资料 DTO（被 AuthService / FollowService / CommentService / NotificationService 复用）
func (s *UserService) BuildProfile(ctx context.Context, user *model.User, viewerID uint) (*dto.UserProfileDTO, error) {
	followerCount, err := s.followRepo.CountFollowers(ctx, user.ID)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("统计粉丝失败: user_id=%d", user.ID), err)
	}
	followingCount, err := s.followRepo.CountFollowing(ctx, user.ID)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("统计关注失败: user_id=%d", user.ID), err)
	}
	articleCount, err := s.articleCount(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	isFollowing := false
	if viewerID > 0 && viewerID != user.ID {
		isFollowing, err = s.followRepo.IsFollowing(ctx, viewerID, user.ID)
		if err != nil {
			return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询关注关系失败: follower_id=%d followed_id=%d", viewerID, user.ID), err)
		}
	}
	var tags []string
	_ = json.Unmarshal([]byte(user.TechTags), &tags)
	return &dto.UserProfileDTO{
		ID:             user.ID,
		Email:          user.Email,
		Nickname:       user.Nickname,
		Avatar:         user.Avatar,
		Bio:            user.Bio,
		TechTags:       tags,
		Role:           user.Role,
		Status:         user.Status,
		GithubID:       user.GithubID,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
		ArticleCount:   articleCount,
		IsFollowing:    isFollowing,
		LastLoginAt:    user.LastLoginAt,
		CreatedAt:      user.CreatedAt,
	}, nil
}

// UpdateProfile 更新个人资料
func (s *UserService) UpdateProfile(ctx context.Context, userID uint, req *dto.UpdateProfileRequest) (*dto.UserProfileDTO, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeUserNotFound, fmt.Sprintf(constants.MsgErrUserNotFound, userID))
		}
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询用户失败: user_id=%d", userID), err)
	}
	if err := s.CheckNicknameUnique(ctx, req.Nickname, userID); err != nil {
		return nil, err
	}
	user.Nickname = req.Nickname
	user.Avatar = req.Avatar
	user.Bio = req.Bio
	tagsJSON, _ := json.Marshal(req.TechTags)
	user.TechTags = string(tagsJSON)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("更新用户失败: user_id=%d nickname=%s", userID, req.Nickname), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogUserProfileUpdated, user.ID, user.Nickname))
	return s.BuildProfile(ctx, user, userID)
}

// CheckNicknameUnique 昵称唯一性检查（被 AuthService.Register 复用）
func (s *UserService) CheckNicknameUnique(ctx context.Context, nickname string, excludeID uint) error {
	users, _, err := s.userRepo.List(ctx, 1, 100, nickname, -1, -1)
	if err != nil {
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询用户失败: field=nickname role=%d", constants.RoleUser.Int()), err)
	}
	for _, u := range users {
		if u.Nickname == nickname && u.ID != excludeID {
			return util.NewAppError(constants.CodeNicknameExists, fmt.Sprintf("昵称已被使用: nickname=%s", nickname))
		}
	}
	return nil
}

// ListFollowers 粉丝列表
func (s *UserService) ListFollowers(ctx context.Context, userID, page, pageSize uint) (*dto.PageResult, error) {
	ids, total, err := s.followRepo.ListFollowerIDs(ctx, userID, int(page), int(pageSize))
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询粉丝失败: user_id=%d", userID), err)
	}
	return s.buildProfiles(ctx, ids, total, page, pageSize, userID)
}

// ListFollowing 关注列表
func (s *UserService) ListFollowing(ctx context.Context, userID, page, pageSize uint) (*dto.PageResult, error) {
	ids, total, err := s.followRepo.ListFollowingIDs(ctx, userID, int(page), int(pageSize))
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询关注失败: user_id=%d", userID), err)
	}
	return s.buildProfiles(ctx, ids, total, page, pageSize, userID)
}

// buildProfiles 批量构造用户资料（内部方法）
func (s *UserService) buildProfiles(ctx context.Context, ids []uint, total int64, page, pageSize, viewerID uint) (*dto.PageResult, error) {
	items := make([]*dto.UserProfileDTO, 0, len(ids))
	for _, id := range ids {
		user, err := s.userRepo.FindByID(ctx, id)
		if err != nil {
			continue
		}
		profile, err := s.BuildProfile(ctx, user, viewerID)
		if err != nil {
			return nil, err
		}
		items = append(items, profile)
	}
	return &dto.PageResult{Items: items, Total: total, Page: int(page), PageSize: int(pageSize)}, nil
}

// articleCount 用户文章数（内部方法）
func (s *UserService) articleCount(ctx context.Context, userID uint) (int64, error) {
	_, total, err := s.articleRepo.List(ctx, 1, 1, "", 0, "", userID, -1)
	if err != nil {
		return 0, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("统计文章失败: user_id=%d", userID), err)
	}
	return total, nil
}
