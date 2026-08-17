package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/service"
	"github.com/techblog/community/internal/util"
)

// FollowHandler 关注处理器
type FollowHandler struct {
	followSvc *service.FollowService
}

// NewFollowHandler 构造关注处理器
func NewFollowHandler(followSvc *service.FollowService) *FollowHandler {
	return &FollowHandler{followSvc: followSvc}
}

// Follow POST /api/v1/users/:id/follow 关注作者
func (h *FollowHandler) Follow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "用户 ID 无效: field=id")
		return
	}
	if err := h.followSvc.Follow(c.Request.Context(), util.CurrentUserID(c), uint(id)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgFollowed, gin.H{"followed_id": id})
}

// Unfollow DELETE /api/v1/users/:id/follow 取消关注
func (h *FollowHandler) Unfollow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "用户 ID 无效: field=id")
		return
	}
	if err := h.followSvc.Unfollow(c.Request.Context(), util.CurrentUserID(c), uint(id)); err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgUnfollowed, gin.H{"followed_id": id})
}
