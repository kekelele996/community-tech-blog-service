package dto

// TrendPoint 趋势数据点
type TrendPoint struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// StatsResponse 平台数据统计
type StatsResponse struct {
	DAU           int64        `json:"dau"`
	TotalArticles int64        `json:"total_articles"`
	TotalUsers    int64        `json:"total_users"`
	TotalComments int64        `json:"total_comments"`
	DailyActive   []TrendPoint `json:"daily_active"`
	NewUsers      []TrendPoint `json:"new_users"`
}

// AdminArticleQuery 后台文章查询
type AdminArticleQuery struct {
	PageQuery
	Keyword string `form:"keyword"`
	Status  int    `form:"status" binding:"omitempty,oneof=0 1 2"`
	AuthorID uint  `form:"author_id" binding:"omitempty,min=1"`
}
