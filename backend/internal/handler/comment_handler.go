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

// CommentHandler 评论处理器
type CommentHandler struct {
	commentSvc *service.CommentService
}

// NewCommentHandler 构造评论处理器
func NewCommentHandler(commentSvc *service.CommentService) *CommentHandler {
	return &CommentHandler{commentSvc: commentSvc}
}

// Create POST /api/v1/articles/:id/comments 发表评论
func (h *CommentHandler) Create(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "文章 ID 无效: field=id")
		return
	}
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "评论参数校验失败: field=content role=user")
		return
	}
	comment, err := h.commentSvc.Create(c.Request.Context(), util.CurrentUserID(c), uint(id), &req)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgCreated, comment)
}

// List GET /api/v1/articles/:id/comments 评论列表
func (h *CommentHandler) List(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "文章 ID 无效: field=id")
		return
	}
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	page, pageSize := q.Normalize()
	result, err := h.commentSvc.ListByArticle(c.Request.Context(), uint(id), page, pageSize)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, result)
}

// Delete DELETE /api/v1/comments/:id 删除评论
func (h *CommentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "评论 ID 无效: field=id")
		return
	}
	if err := h.commentSvc.Delete(c.Request.Context(), util.CurrentUserID(c), uint(util.CurrentRole(c)), uint(id)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgDeleted, gin.H{"id": id})
}
