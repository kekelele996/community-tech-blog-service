package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/service"
	"github.com/techblog/community/internal/util"
)

var mutatingMethods = map[string]bool{"POST": true, "PUT": true, "DELETE": true, "PATCH": true}

// Audit 审计日志中间件：自动记录所有写操作
func Audit(auditSvc *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if !mutatingMethods[c.Request.Method] {
			return
		}
		status := c.Writer.Status()
		if status >= 400 {
			return
		}
		user, ok := util.GetAuthUser(c)
		if !ok {
			return
		}
		path := c.Request.URL.Path
		action := "WRITE_" + c.Request.Method
		module := "api"
		if idx := strings.Index(path, "/api/v1/"); idx >= 0 {
			parts := strings.Split(strings.TrimPrefix(path[idx+len("/api/v1/"):], "/"), "/")
			if len(parts) > 0 && parts[0] != "" {
				module = parts[0]
			}
		}
		auditSvc.Record(c.Request.Context(), user.ID, user.Email, user.Role, action, module, c.Request.Method+" "+path, util.ClientIP(c))
	}
}
