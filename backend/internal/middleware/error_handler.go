package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/util"
)

// ErrorHandler 全局错误处理中间件：recover 恐慌 + 统一错误响应
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				util.LogError(util.GetRequestID(c), "panic recovered", "panic", r, "path", c.Request.URL.Path)
				c.AbortWithStatusJSON(http.StatusInternalServerError, util.Response{
					Code: constants.CodeInternalError, Message: constants.MessageOf(constants.CodeInternalError),
				})
			}
		}()
		c.Next()
	}
}
