package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/util"
)

// RequireRole RBAC 权限中间件：仅允许指定角色访问
func RequireRole(roles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := util.GetAuthUser(c)
		if !ok {
			util.Fail(c, 401, constants.CodeUnauthorized, constants.MessageOf(constants.CodeUnauthorized))
			return
		}
		for _, role := range roles {
			if user.Role == role {
				c.Next()
				return
			}
		}
		util.Fail(c, 403, constants.CodeForbidden, "无权限执行该操作: operator_role="+strconv.Itoa(user.Role))
	}
}
