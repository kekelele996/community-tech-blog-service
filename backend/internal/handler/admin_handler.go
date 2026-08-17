package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/service"
	"github.com/techblog/community/internal/util"
)

// AdminHandler 后台管理处理器
type AdminHandler struct {
	adminSvc *service.AdminService
}

// NewAdminHandler 构造后台管理处理器
func NewAdminHandler(adminSvc *service.AdminService) *AdminHandler {
	return &AdminHandler{adminSvc: adminSvc}
}

// Stats GET /api/v1/admin/stats 平台数据统计
func (h *AdminHandler) Stats(c *gin.Context) {
	stats, err := h.adminSvc.Stats(c.Request.Context(), util.CurrentUserID(c))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, stats)
}

// ListUsers GET /api/v1/admin/users 用户列表
func (h *AdminHandler) ListUsers(c *gin.Context) {
	var q dto.UserQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "用户列表参数校验失败: field=page/page_size")
		return
	}
	result, err := h.adminSvc.ListUsers(c.Request.Context(), util.CurrentUserID(c), &q)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, result)
}

// UpdateUserStatus PUT /api/v1/admin/users/:id/status 禁用/启用用户
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "用户 ID 无效: field=id")
		return
	}
	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "状态参数校验失败: field=status role=admin")
		return
	}
	if err := h.adminSvc.UpdateUserStatus(c.Request.Context(), util.CurrentUserID(c), uint(id), uint(req.Status)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgUpdated, gin.H{"id": id, "status": req.Status})
}

// ListArticles GET /api/v1/admin/articles 文章管理列表
func (h *AdminHandler) ListArticles(c *gin.Context) {
	var q dto.AdminArticleQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "文章列表参数校验失败: field=page/page_size")
		return
	}
	result, err := h.adminSvc.ListArticles(c.Request.Context(), &q)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, result)
}

// UpdateArticleStatus PUT /api/v1/admin/articles/:id/status 下架/恢复文章
func (h *AdminHandler) UpdateArticleStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "文章 ID 无效: field=id")
		return
	}
	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "状态参数校验失败: field=status role=admin")
		return
	}
	if err := h.adminSvc.UpdateArticleStatus(c.Request.Context(), util.CurrentUserID(c), uint(id), uint(req.Status)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgUpdated, gin.H{"id": id, "status": req.Status})
}

// ListTopics GET /api/v1/admin/topics 话题管理列表
func (h *AdminHandler) ListTopics(c *gin.Context) {
	var q dto.TopicQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "话题列表参数校验失败: field=page/page_size")
		return
	}
	result, err := h.adminSvc.ListTopics(c.Request.Context(), &q)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, result)
}

// UpdateTopicStatus PUT /api/v1/admin/topics/:id/status 启用/禁用话题
func (h *AdminHandler) UpdateTopicStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "话题 ID 无效: field=id")
		return
	}
	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "状态参数校验失败: field=status role=admin")
		return
	}
	if err := h.adminSvc.UpdateTopicStatus(c.Request.Context(), util.CurrentUserID(c), uint(id), uint(req.Status)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgUpdated, gin.H{"id": id, "status": req.Status})
}
