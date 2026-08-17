package util

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/constants"
)

// Response 统一响应包裹：{ "code": 0, "message": "ok", "data": ... }
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OK 成功响应
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: constants.CodeSuccess, Message: constants.MsgOK, Data: data})
}

// OKWithMessage 成功响应（自定义 message）
func OKWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: constants.CodeSuccess, Message: message, Data: data})
}

// Fail 失败响应
func Fail(c *gin.Context, httpStatus, code int, message string) {
	c.AbortWithStatusJSON(httpStatus, Response{Code: code, Message: message})
}

// FailFromError 根据 AppError 输出失败响应
func FailFromError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		c.AbortWithStatusJSON(httpStatusOf(appErr.Code), Response{Code: appErr.Code, Message: appErr.Message})
		return
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, Response{Code: constants.CodeInternalError, Message: constants.MessageOf(constants.CodeInternalError)})
}

func httpStatusOf(code int) int {
	switch {
	case code >= 40000 && code < 40100:
		return http.StatusBadRequest
	case code >= 40100 && code < 40300:
		return http.StatusUnauthorized
	case code >= 40300 && code < 40400:
		return http.StatusForbidden
	case code >= 40400 && code < 40500:
		return http.StatusNotFound
	case code >= 40900 && code < 41000:
		return http.StatusConflict
	case code >= 42200 && code < 42300:
		return http.StatusUnprocessableEntity
	case code >= 42900 && code < 43000:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
