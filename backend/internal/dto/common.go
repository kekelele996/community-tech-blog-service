package dto

// PageQuery 通用分页查询参数：统一 page / page_size
type PageQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// PageResult 统一分页响应
type PageResult struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// Normalize 归一化分页参数
func (p *PageQuery) Normalize() (int, int) {
	page := p.Page
	if page <= 0 {
		page = 1
	}
	pageSize := p.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

// IDRequest 路径参数 id
type IDRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}
