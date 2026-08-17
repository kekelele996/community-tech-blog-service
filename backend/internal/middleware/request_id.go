package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/techblog/community/internal/util"
)

// RequestID 请求追踪中间件：为每个请求生成 X-Request-ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		util.SetRequestID(c, requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
