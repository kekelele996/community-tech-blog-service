package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/util"
)

// RequestLogger 请求日志中间件：包含 request_id / method / path / status / latency_ms
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()
		util.LogInfo(util.GetRequestID(c), "http request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", latency,
			"ip", util.ClientIP(c),
		)
	}
}
