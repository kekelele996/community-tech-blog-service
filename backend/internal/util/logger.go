package util

import (
	"log/slog"
	"os"
	"sync"
)

// logger 全局结构化日志（log/slog JSON 输出）。
// 所有 handler / service / middleware 都必须引用本模块（屎山耦合点 1）。
var (
	loggerOnce sync.Once
	logger     *slog.Logger
)

// InitLogger 初始化全局日志器
func InitLogger(level slog.Level) {
	loggerOnce.Do(func() {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
		logger = slog.New(handler)
		slog.SetDefault(logger)
	})
}

// Logger 返回全局日志器
func Logger() *slog.Logger {
	if logger == nil {
		InitLogger(slog.LevelInfo)
	}
	return logger
}

// LogError 带 request_id 记录错误日志
func LogError(requestID, msg string, args ...any) {
	args = append([]any{"request_id", requestID}, args...)
	Logger().Error(msg, args...)
}

// LogInfo 带 request_id 记录普通日志
func LogInfo(requestID, msg string, args ...any) {
	args = append([]any{"request_id", requestID}, args...)
	Logger().Info(msg, args...)
}
