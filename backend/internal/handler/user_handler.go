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

// UserHandler 用户处理器
type UserHandler struct {
	userSvc *service.UserService
}

// NewUserHandler 构造用户处理器
func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// GetProfile GET /api/v1/users/:id 用户主页（复用 UserService.GetProfile）
func (h *UserHandler) GetProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "用户 ID 无效: field=id")
		return
	}
	profile, err := h.userSvc.GetProfile(c.Request.Context(), uint(id), util.CurrentUserID(c))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, profile)
}

// UpdateProfile PUT /api/v1/users/me 更新个人资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "资料参数校验失败: field=body role=user")
		return
	}
	profile, err := h.userSvc.UpdateProfile(c.Request.Context(), util.CurrentUserID(c), &req)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgUpdated, profile)
}

// Followers GET /api/v1/users/:id/followers 粉丝列表
func (h *UserHandler) Followers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "用户 ID 无效: field=id")
		return
	}
	var q dto.FollowQuery
	_ = c.ShouldBindQuery(&q)
	page, pageSize := q.Normalize()
	result, err := h.userSvc.ListFollowers(c.Request.Context(), uint(id), uint(page), uint(pageSize))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, result)
}

// Following GET /api/v1/users/:id/following 关注列表
func (h *UserHandler) Following(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "用户 ID 无效: field=id")
		return
	}
	var q dto.FollowQuery
	_ = c.ShouldBindQuery(&q)
	page, pageSize := q.Normalize()
	result, err := h.userSvc.ListFollowing(c.Request.Context(), uint(id), uint(page), uint(pageSize))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, result)
}
