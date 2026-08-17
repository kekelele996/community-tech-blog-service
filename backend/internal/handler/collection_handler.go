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

// CollectionHandler 收藏夹处理器
type CollectionHandler struct {
	collectionSvc *service.CollectionService
}

// NewCollectionHandler 构造收藏夹处理器
func NewCollectionHandler(collectionSvc *service.CollectionService) *CollectionHandler {
	return &CollectionHandler{collectionSvc: collectionSvc}
}

// Create POST /api/v1/collections 创建收藏夹
func (h *CollectionHandler) Create(c *gin.Context) {
	var req dto.CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "收藏夹参数校验失败: field=name role=user")
		return
	}
	collection, err := h.collectionSvc.Create(c.Request.Context(), util.CurrentUserID(c), &req)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgCreated, collection)
}

// Update PUT /api/v1/collections/:id 更新收藏夹
func (h *CollectionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "收藏夹 ID 无效: field=id")
		return
	}
	var req dto.UpdateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "收藏夹参数校验失败: field=name")
		return
	}
	collection, err := h.collectionSvc.Update(c.Request.Context(), util.CurrentUserID(c), uint(id), &req)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgUpdated, collection)
}

// Delete DELETE /api/v1/collections/:id 删除收藏夹
func (h *CollectionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "收藏夹 ID 无效: field=id")
		return
	}
	if err := h.collectionSvc.Delete(c.Request.Context(), util.CurrentUserID(c), uint(id)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgDeleted, gin.H{"id": id})
}

// ListMine GET /api/v1/collections 我的收藏夹
func (h *CollectionHandler) ListMine(c *gin.Context) {
	items, err := h.collectionSvc.ListMine(c.Request.Context(), util.CurrentUserID(c))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, items)
}

// ListPublic GET /api/v1/users/:id/collections 用户公开收藏夹
func (h *CollectionHandler) ListPublic(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "用户 ID 无效: field=id")
		return
	}
	items, err := h.collectionSvc.ListPublic(c.Request.Context(), uint(id), util.CurrentUserID(c))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, items)
}

// Detail GET /api/v1/collections/:id 收藏夹详情（含文章）
func (h *CollectionHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "收藏夹 ID 无效: field=id")
		return
	}
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	page, pageSize := q.Normalize()
	detail, err := h.collectionSvc.Detail(c.Request.Context(), util.CurrentUserID(c), uint(id), page, pageSize)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, detail)
}

// AddArticle POST /api/v1/collections/:id/articles 收藏文章
func (h *CollectionHandler) AddArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "收藏夹 ID 无效: field=id")
		return
	}
	var req dto.AddArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "收藏参数校验失败: field=article_id")
		return
	}
	if err := h.collectionSvc.AddArticle(c.Request.Context(), util.CurrentUserID(c), uint(id), req.ArticleID); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgCollected, gin.H{"collection_id": id, "article_id": req.ArticleID})
}

// RemoveArticle DELETE /api/v1/collections/:id/articles/:articleId 取消收藏
func (h *CollectionHandler) RemoveArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "收藏夹 ID 无效: field=id")
		return
	}
	articleID, err := strconv.ParseUint(c.Param("articleId"), 10, 64)
	if err != nil || articleID == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "文章 ID 无效: field=articleId")
		return
	}
	if err := h.collectionSvc.RemoveArticle(c.Request.Context(), util.CurrentUserID(c), uint(id), uint(articleID)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgUncollected, gin.H{"collection_id": id, "article_id": articleID})
}
