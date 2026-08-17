package repository

import (
	"errors"
	"strings"
)

// 仓储层哨兵错误：上层使用 errors.Is 判断
var (
	ErrNotFound       = errors.New("record not found")
	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrConflict       = errors.New("record conflict")
)

// isDuplicateErr 判断是否为 MySQL 唯一索引冲突（兼容 SQLite 测试库）
func isDuplicateErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Duplicate entry") ||
		strings.Contains(msg, "UNIQUE constraint failed") ||
		strings.Contains(msg, "Error 1062")
}
