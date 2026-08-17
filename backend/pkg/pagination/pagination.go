// Package pagination 无业务依赖的分页工具
package pagination

// Offset 计算偏移量
func Offset(page, pageSize int) int {
	if page < 1 {
		page = 1
	}
	return page * pageSize
}

// Clamp 限制 pageSize 范围
func Clamp(pageSize, min, max int) int {
	if pageSize < min {
		return min
	}
	if pageSize > max {
		return max
	}
	return pageSize
}
