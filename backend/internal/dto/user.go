package dto

import "time"

// UpdateProfileRequest 更新个人资料
type UpdateProfileRequest struct {
	Nickname string   `json:"nickname" binding:"required,min=2,max=32"`
	Avatar   string   `json:"avatar" binding:"omitempty,max=512"`
	Bio      string   `json:"bio" binding:"omitempty,max=512"`
	TechTags []string `json:"tech_tags" binding:"omitempty,dive,max=32"`
}

// UserProfileDTO 用户资料（含关注/粉丝统计与当前访问者是否已关注）
type UserProfileDTO struct {
	ID            uint      `json:"id"`
	Email         string    `json:"email"`
	Nickname      string    `json:"nickname"`
	Avatar        string    `json:"avatar"`
	Bio           string    `json:"bio"`
	TechTags      []string  `json:"tech_tags"`
	Role          int       `json:"role"`
	Status        int       `json:"status"`
	GithubID      string    `json:"github_id"`
	FollowerCount int64     `json:"follower_count"`
	FollowingCount int64    `json:"following_count"`
	ArticleCount  int64     `json:"article_count"`
	IsFollowing   bool      `json:"is_following"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// UserListItemDTO 后台用户列表项
type UserListItemDTO struct {
	ID          uint       `json:"id"`
	Email       string     `json:"email"`
	Nickname    string     `json:"nickname"`
	Avatar      string     `json:"avatar"`
	Role        int        `json:"role"`
	Status      int        `json:"status"`
	ArticleCount int64     `json:"article_count"`
	FollowerCount int64    `json:"follower_count"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// UserQuery 后台用户查询
type UserQuery struct {
	PageQuery
	Keyword string `form:"keyword"`
	Status  int    `form:"status"`
	Role    int    `form:"role"`
}

// UpdateStatusRequest 状态更新（用户禁用/启用、文章下架、话题禁用）
type UpdateStatusRequest struct {
	Status int `json:"status" binding:"required,oneof=0 1 2"`
}

// FollowQuery 关注列表查询
type FollowQuery struct {
	PageQuery
}
