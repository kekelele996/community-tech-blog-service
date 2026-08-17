package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/service"
	"github.com/techblog/community/internal/util"
)

// AuditHandler 审计日志处理器
type AuditHandler struct {
	auditSvc *service.AuditService
}

// NewAuditHandler 构造审计日志处理器
func NewAuditHandler(auditSvc *service.AuditService) *AuditHandler {
	return &AuditHandler{auditSvc: auditSvc}
}

// List GET /api/v1/admin/audit-logs 审计日志列表（管理员）
func (h *AuditHandler) List(c *gin.Context) {
	var q dto.AuditQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "审计日志参数校验失败: field=page/page_size")
		return
	}
	result, err := h.auditSvc.List(c.Request.Context(), &q)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, result)
}
