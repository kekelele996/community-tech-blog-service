package util

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	ctxAuthUserKey  = "auth_user"
	ctxRequestIDKey = "request_id"
	requestIDCtxKey = ctxKey("request_id")
)

type ctxKey string

// AuthUser 当前登录用户上下文
type AuthUser struct {
	ID    uint
	Role  int
	Email string
}

// SetAuthUser 写入当前用户
func SetAuthUser(c *gin.Context, user *AuthUser) {
	c.Set(ctxAuthUserKey, user)
}

// GetAuthUser 读取当前用户
func GetAuthUser(c *gin.Context) (*AuthUser, bool) {
	v, ok := c.Get(ctxAuthUserKey)
	if !ok {
		return nil, false
	}
	user, ok := v.(*AuthUser)
	return user, ok
}

// CurrentUserID 当前用户 ID（未登录返回 0）
func CurrentUserID(c *gin.Context) uint {
	if user, ok := GetAuthUser(c); ok {
		return user.ID
	}
	return 0
}

// CurrentRole 当前用户角色（未登录返回 0）
func CurrentRole(c *gin.Context) int {
	if user, ok := GetAuthUser(c); ok {
		return user.Role
	}
	return 0
}

// SetRequestID 写入请求 ID（同时写入 gin context 与 request context，供 service 层读取）
func SetRequestID(c *gin.Context, requestID string) {
	c.Set(ctxRequestIDKey, requestID)
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), requestIDCtxKey, requestID))
}

// GetRequestID 读取请求 ID：兼容 *gin.Context 与 context.Context（service 层传入 c.Request.Context()）
func GetRequestID(ctx context.Context) string {
	if gc, ok := ctx.(*gin.Context); ok {
		if v, ok := gc.Get(ctxRequestIDKey); ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	if v, ok := ctx.Value(requestIDCtxKey).(string); ok {
		return v
	}
	return ""
}

// ClientIP 获取客户端 IP（优先 X-Forwarded-For）
func ClientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		return xff
	}
	return c.ClientIP()
}
