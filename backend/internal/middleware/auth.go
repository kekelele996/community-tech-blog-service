package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/util"
)

// Auth JWT 认证中间件：从 Authorization: Bearer <token> 解析用户并写入上下文
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			util.Fail(c, 401, constants.CodeUnauthorized, constants.MessageOf(constants.CodeUnauthorized))
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := util.ParseToken(secret, tokenStr)
		if err != nil {
			util.Fail(c, 401, constants.CodeUnauthorized, "登录状态无效或已过期")
			return
		}
		util.SetAuthUser(c, &util.AuthUser{ID: claims.UserID, Role: claims.Role, Email: claims.Email})
		c.Next()
	}
}
