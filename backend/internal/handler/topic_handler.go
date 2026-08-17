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

// TopicHandler 话题处理器
type TopicHandler struct {
	topicSvc *service.TopicService
}

// NewTopicHandler 构造话题处理器
func NewTopicHandler(topicSvc *service.TopicService) *TopicHandler {
	return &TopicHandler{topicSvc: topicSvc}
}

// Create POST /api/v1/topics 创建话题（管理员）
func (h *TopicHandler) Create(c *gin.Context) {
	var req dto.CreateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "话题参数校验失败: field=name role=admin")
		return
	}
	topic, err := h.topicSvc.Create(c.Request.Context(), util.CurrentUserID(c), &req)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgCreated, topic)
}

// Update PUT /api/v1/topics/:id 更新话题（管理员）
func (h *TopicHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "话题 ID 无效: field=id")
		return
	}
	var req dto.UpdateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "话题参数校验失败: field=name")
		return
	}
	topic, err := h.topicSvc.Update(c.Request.Context(), util.CurrentUserID(c), uint(id), &req)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgUpdated, topic)
}

// Delete DELETE /api/v1/topics/:id 删除话题（管理员）
func (h *TopicHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "话题 ID 无效: field=id")
		return
	}
	if err := h.topicSvc.Delete(c.Request.Context(), util.CurrentUserID(c), uint(id)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgDeleted, gin.H{"id": id})
}

// List GET /api/v1/topics 话题广场；?all=1 返回全部启用话题（文章编辑器选择用）
func (h *TopicHandler) List(c *gin.Context) {
	if c.Query("all") == "1" {
		items, err := h.topicSvc.ListAll(c.Request.Context())
		if err != nil {
			util.FailFromError(c, err)
			return
		}
		util.OK(c, items)
		return
	}
	var q dto.TopicQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "话题列表参数校验失败: field=page/page_size")
		return
	}
	result, err := h.topicSvc.List(c.Request.Context(), &q)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, result)
}

// Detail GET /api/v1/topics/:id 话题详情（热门 + 最新文章）
func (h *TopicHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "话题 ID 无效: field=id")
		return
	}
	detail, err := h.topicSvc.Detail(c.Request.Context(), uint(id))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, detail)
}

// SetStatus PUT /api/v1/topics/:id/status 设置话题状态（管理员）
func (h *TopicHandler) SetStatus(c *gin.Context) {
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
	if err := h.topicSvc.SetStatus(c.Request.Context(), util.CurrentUserID(c), uint(id), uint(req.Status)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgUpdated, gin.H{"id": id})
}
