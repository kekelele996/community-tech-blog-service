package service

import (
	"context"
	"fmt"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/repository"
	"github.com/techblog/community/internal/util"
)

// AuditService 操作审计日志服务（横切关注点 2）
type AuditService struct {
	repo *repository.AuditRepository
}

// NewAuditService 构造审计服务
func NewAuditService(repo *repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

// Record 记录审计日志（被 service 埋点与 audit middleware 复用）
func (s *AuditService) Record(ctx context.Context, userID uint, username string, role int, action, module, detail, ip string) {
	log := &model.AuditLog{
		UserID: userID, Username: username, Role: role,
		Action: action, Module: module, Detail: detail,
		IP: ip, RequestID: util.GetRequestID(ctx),
	}
	if err := s.repo.Create(ctx, log); err != nil {
		util.LogError(util.GetRequestID(ctx), "record audit log failed", "user_id", userID, "action", action, "error", err.Error())
		return
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogAuditRecorded, userID, action, module))
}

// List 分页查询审计日志
func (s *AuditService) List(ctx context.Context, query *dto.AuditQuery) (*dto.PageResult, error) {
	page, pageSize := query.Normalize()
	items, total, err := s.repo.List(ctx, page, pageSize, query.UserID, query.Action, query.Module)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, "查询审计日志失败", err)
	}
	result := make([]*dto.AuditLogDTO, 0, len(items))
	for i := range items {
		result = append(result, &dto.AuditLogDTO{
			ID: items[i].ID, UserID: items[i].UserID, Username: items[i].Username,
			Role: items[i].Role, Action: items[i].Action, Module: items[i].Module,
			Detail: items[i].Detail, IP: items[i].IP, RequestID: items[i].RequestID,
			CreatedAt: items[i].CreatedAt,
		})
	}
	return &dto.PageResult{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}
