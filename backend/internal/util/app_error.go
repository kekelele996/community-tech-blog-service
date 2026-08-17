package util

import "fmt"

// AppError 业务错误：错误码 + 面向用户的 message + 原始错误链
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d message=%s cause=%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d message=%s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

// NewAppError 构造业务错误
func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// WrapAppError 包装底层错误并透传错误链
func WrapAppError(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// IsAppError 判断错误链中是否存在 AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}
