package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/repository"
	"github.com/techblog/community/internal/util"
)

// AdminService 后台管理服务：平台统计、用户/文章/话题管理
type AdminService struct {
	userRepo    *repository.UserRepository
	articleRepo *repository.ArticleRepository
	commentRepo *repository.CommentRepository
	loginRepo   *repository.LoginRepository
	articleSvc  *ArticleService
	topicSvc    *TopicService
	auditSvc    *AuditService
}

// NewAdminService 构造后台管理服务
func NewAdminService(userRepo *repository.UserRepository, articleRepo *repository.ArticleRepository, commentRepo *repository.CommentRepository, loginRepo *repository.LoginRepository, articleSvc *ArticleService, topicSvc *TopicService, auditSvc *AuditService) *AdminService {
	return &AdminService{userRepo: userRepo, articleRepo: articleRepo, commentRepo: commentRepo, loginRepo: loginRepo, articleSvc: articleSvc, topicSvc: topicSvc, auditSvc: auditSvc}
}

// Stats 平台数据统计：日活 / 文章数 / 用户数 / 评论数 + 近 7 天趋势
func (s *AdminService) Stats(ctx context.Context, operatorID uint) (*dto.StatsResponse, error) {
	today := time.Now().Format("2006-01-02")
	dau, err := s.loginRepo.CountDAU(ctx, today)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, "统计日活失败", err)
	}
	totalArticles, err := s.articleRepo.Count(ctx)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, "统计文章数失败", err)
	}
	totalUsers, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, "统计用户数失败", err)
	}
	totalComments, err := s.commentRepo.Count(ctx)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, "统计评论数失败", err)
	}
	start := time.Now().AddDate(0, 0, -6).Format("2006-01-02")
	end := today
	dauMap, err := s.loginRepo.CountDAUByRange(ctx, start, end)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, "统计日活趋势失败", err)
	}
	userMap, err := s.userRepo.CountByDate(ctx, start, end)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, "统计注册趋势失败", err)
	}
	dailyActive := make([]dto.TrendPoint, 0, 7)
	newUsers := make([]dto.TrendPoint, 0, 7)
	for i := 6; i >= 0; i-- {
		d := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		dailyActive = append(dailyActive, dto.TrendPoint{Date: d, Count: dauMap[d]})
		newUsers = append(newUsers, dto.TrendPoint{Date: d, Count: userMap[d]})
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogAdminStats, operatorID, dau, totalArticles, totalUsers))
	s.auditSvc.Record(ctx, operatorID, "", 0, "ADMIN_STATS", "admin", "查看平台数据统计", "")
	return &dto.StatsResponse{
		DAU: dau, TotalArticles: totalArticles, TotalUsers: totalUsers, TotalComments: totalComments,
		DailyActive: dailyActive, NewUsers: newUsers,
	}, nil
}

// ListUsers 用户列表（后台）
func (s *AdminService) ListUsers(ctx context.Context, operatorID uint, query *dto.UserQuery) (*dto.PageResult, error) {
	page, pageSize := query.Normalize()
	users, total, err := s.userRepo.List(ctx, page, pageSize, query.Keyword, query.Status, query.Role)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, "查询用户列表失败", err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogUserListQueried, operatorID, page, pageSize))
	items := make([]*dto.UserListItemDTO, 0, len(users))
	for i := range users {
		articleCount := int64(0)
		_, count, err := s.articleRepo.List(ctx, 1, 1, "", 0, "", users[i].ID, -1)
		if err == nil {
			articleCount = count
		}
		items = append(items, &dto.UserListItemDTO{
			ID: users[i].ID, Email: users[i].Email, Nickname: users[i].Nickname,
			Avatar: users[i].Avatar, Role: users[i].Role, Status: users[i].Status,
			ArticleCount: articleCount, LastLoginAt: users[i].LastLoginAt, CreatedAt: users[i].CreatedAt,
		})
	}
	return &dto.PageResult{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

// UpdateUserStatus 禁用/启用用户
func (s *AdminService) UpdateUserStatus(ctx context.Context, operatorID, userID, status uint) error {
	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeUserNotFound, fmt.Sprintf(constants.MsgErrUserNotFound, userID))
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询用户失败: user_id=%d", userID), err)
	}
	if err := s.userRepo.UpdateStatus(ctx, userID, int(status)); err != nil {
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("更新用户状态失败: user_id=%d operator=%d status=%d", userID, operatorID, status), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogUserStatusChanged, userID, status, operatorID))
	s.auditSvc.Record(ctx, operatorID, "", 0, "USER_STATUS", "admin", fmt.Sprintf("设置用户 #%d 状态=%d", userID, status), "")
	return nil
}

// ListArticles 文章列表（后台，复用 ArticleService.List）
func (s *AdminService) ListArticles(ctx context.Context, query *dto.AdminArticleQuery) (*dto.PageResult, error) {
	articleQuery := &dto.ArticleQuery{
		PageQuery: query.PageQuery,
		Keyword:   query.Keyword,
		Status:    query.Status,
		AuthorID:  query.AuthorID,
		Sort:      constants.SortLatest.String(),
	}
	return s.articleSvc.List(ctx, articleQuery)
}

// UpdateArticleStatus 下架/恢复文章（复用 ArticleService.AdminUpdateStatus）
func (s *AdminService) UpdateArticleStatus(ctx context.Context, operatorID, articleID, status uint) error {
	return s.articleSvc.AdminUpdateStatus(ctx, operatorID, articleID, status)
}

// ListTopics 话题列表（后台，复用 TopicService.List）
func (s *AdminService) ListTopics(ctx context.Context, query *dto.TopicQuery) (*dto.PageResult, error) {
	return s.topicSvc.List(ctx, query)
}

// UpdateTopicStatus 启用/禁用话题（复用 TopicService.SetStatus）
func (s *AdminService) UpdateTopicStatus(ctx context.Context, operatorID, topicID, status uint) error {
	return s.topicSvc.SetStatus(ctx, operatorID, topicID, status)
}
