package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/service"
	"github.com/techblog/community/internal/util"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authSvc *service.AuthService
}

// NewAuthHandler 构造认证处理器
func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// SendCode POST /api/v1/auth/code 发送邮箱验证码
func (h *AuthHandler) SendCode(c *gin.Context) {
	var req dto.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "验证码请求参数错误: field=email")
		return
	}
	resp, err := h.authSvc.SendCode(c.Request.Context(), req.Email)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, resp)
}

// Register POST /api/v1/auth/register 邮箱注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "注册参数校验失败: field=body role=user")
		return
	}
	resp, err := h.authSvc.Register(c.Request.Context(), &req, util.ClientIP(c))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgRegisterSuccess, resp)
}

// Login POST /api/v1/auth/login 密码登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "登录参数校验失败: field=body role=user")
		return
	}
	resp, err := h.authSvc.Login(c.Request.Context(), &req, util.ClientIP(c))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgLoginSuccess, resp)
}

// LoginWithCode POST /api/v1/auth/login/code 验证码登录
func (h *AuthHandler) LoginWithCode(c *gin.Context) {
	var req dto.LoginCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "验证码登录参数校验失败: field=body role=user")
		return
	}
	resp, err := h.authSvc.LoginWithCode(c.Request.Context(), &req, util.ClientIP(c))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgLoginSuccess, resp)
}

// GithubLogin POST /api/v1/auth/github GitHub 第三方登录
func (h *AuthHandler) GithubLogin(c *gin.Context) {
	var req dto.GithubLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, "GitHub 登录参数校验失败: field=code")
		return
	}
	resp, err := h.authSvc.GithubLogin(c.Request.Context(), &req, util.ClientIP(c))
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OKWithMessage(c, constants.MsgLoginSuccess, resp)
}

// Me GET /api/v1/auth/me 当前用户信息（复用 UserService.GetProfile）
func (h *AuthHandler) Me(c *gin.Context) {
	userID := util.CurrentUserID(c)
	profile, err := h.authSvc.Me(c.Request.Context(), userID)
	if err != nil {
		util.FailFromError(c, err)
		return
	}
	util.OK(c, profile)
}
