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

// ArticleHandler 文章处理器
type ArticleHandler struct {
	articleSvc *service.ArticleService
}

// NewArticleHandler 构造文章处理器
func NewArticleHandler(articleSvc *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleSvc: articleSvc}
}

// Create POST /api/v1/articles 创建文章
func (h *ArticleHandler) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "文章参数校验失败: field=title/content role=user")
		return
	}
	item, err := h.articleSvc.Create(c.Request.Context(), util.CurrentUserID(c), &req)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgCreated, item)
}

// Update PUT /api/v1/articles/:id 更新文章
func (h *ArticleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "文章 ID 无效: field=id")
		return
	}
	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "文章参数校验失败: field=title/content")
		return
	}
	item, err := h.articleSvc.Update(c.Request.Context(), util.CurrentUserID(c), uint(util.CurrentRole(c)), uint(id), &req)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgUpdated, item)
}

// List GET /api/v1/articles 文章列表（最新/最热，复用 ArticleService.List）
func (h *ArticleHandler) List(c *gin.Context) {
	var q dto.ArticleQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "文章列表参数校验失败: field=page/page_size")
		return
	}
	result, err := h.articleSvc.List(c.Request.Context(), &q)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, result)
}

// Detail GET /api/v1/articles/:id 文章详情
func (h *ArticleHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "文章 ID 无效: field=id")
		return
	}
	detail, err := h.articleSvc.Detail(c.Request.Context(), uint(id), util.CurrentUserID(c), uint(util.CurrentRole(c)))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, detail)
}

// Publish POST /api/v1/articles/:id/publish 发布文章
func (h *ArticleHandler) Publish(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "文章 ID 无效: field=id")
		return
	}
	item, err := h.articleSvc.Publish(c.Request.Context(), util.CurrentUserID(c), uint(id))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgPublished, item)
}

// Offline POST /api/v1/articles/:id/offline 下架文章
func (h *ArticleHandler) Offline(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "文章 ID 无效: field=id")
		return
	}
	if err := h.articleSvc.Offline(c.Request.Context(), util.CurrentUserID(c), uint(util.CurrentRole(c)), uint(id)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgOfflined, gin.H{"id": id})
}

// Like POST /api/v1/articles/:id/like 点赞
func (h *ArticleHandler) Like(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "文章 ID 无效: field=id")
		return
	}
	if err := h.articleSvc.Like(c.Request.Context(), util.CurrentUserID(c), uint(id)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgLiked, gin.H{"id": id})
}

// Unlike DELETE /api/v1/articles/:id/like 取消点赞
func (h *ArticleHandler) Unlike(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "文章 ID 无效: field=id")
		return
	}
	if err := h.articleSvc.Unlike(c.Request.Context(), util.CurrentUserID(c), uint(id)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgUnliked, gin.H{"id": id})
}

// Feed GET /api/v1/articles/feed 关注动态流
func (h *ArticleHandler) Feed(c *gin.Context) {
	var q dto.FollowQuery
	_ = c.ShouldBindQuery(&q)
	page, pageSize := q.Normalize()
	result, err := h.articleSvc.Feed(c.Request.Context(), util.CurrentUserID(c), uint(page), uint(pageSize))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, result)
}
