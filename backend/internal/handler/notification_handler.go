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

// NotificationHandler 通知处理器
type NotificationHandler struct {
	notificationSvc *service.NotificationService
}

// NewNotificationHandler 构造通知处理器
func NewNotificationHandler(notificationSvc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationSvc: notificationSvc}
}

// List GET /api/v1/notifications 通知列表（全部/未读）
func (h *NotificationHandler) List(c *gin.Context) {
	var q dto.NotificationQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "通知列表参数校验失败: field=page/page_size")
		return
	}
	page, pageSize := q.Normalize()
	result, err := h.notificationSvc.List(c.Request.Context(), util.CurrentUserID(c), page, pageSize, q.UnreadOnly)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, result)
}

// UnreadCount GET /api/v1/notifications/unread-count 未读数量
func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	count, err := h.notificationSvc.CountUnread(c.Request.Context(), util.CurrentUserID(c))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, gin.H{"unread_count": count})
}

// MarkAllRead PUT /api/v1/notifications/read-all 一键全部已读
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	if err := h.notificationSvc.MarkAllRead(c.Request.Context(), util.CurrentUserID(c)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgReadAll, gin.H{})
}

// MarkRead PUT /api/v1/notifications/:id/read 单条已读
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "通知 ID 无效: field=id")
		return
	}
	if err := h.notificationSvc.MarkRead(c.Request.Context(), util.CurrentUserID(c), uint(id)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgUpdated, gin.H{"id": id})
}
