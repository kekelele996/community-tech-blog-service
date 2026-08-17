package constants

// RoleType 角色类型枚举：后端常量定义（同时出现在 model.User.Role、dto、RBAC 中间件、前端 constants/index.ts）
type RoleType int

const (
	RoleUser  RoleType = 1 // 普通用户
	RoleAdmin RoleType = 2 // 管理员
)

func (r RoleType) Int() int { return int(r) }

// UserStatus 用户状态枚举：0 禁用 / 1 启用
type UserStatus int

const (
	UserStatusDisabled UserStatus = 0
	UserStatusActive   UserStatus = 1
)

func (s UserStatus) Int() int { return int(s) }

// ArticleStatus 文章状态枚举：0 草稿 / 1 已发布 / 2 已下架
type ArticleStatus int

const (
	ArticleStatusDraft     ArticleStatus = 0
	ArticleStatusPublished ArticleStatus = 1
	ArticleStatusOffline   ArticleStatus = 2
)

func (s ArticleStatus) Int() int { return int(s) }

// Visibility 收藏夹可见性枚举：0 私密 / 1 公开
type Visibility int

const (
	VisibilityPrivate Visibility = 0
	VisibilityPublic  Visibility = 1
)

func (v Visibility) Int() int { return int(v) }

// NotificationType 通知类型枚举
type NotificationType string

const (
	NotificationTypeFollow  NotificationType = "follow"
	NotificationTypeLike    NotificationType = "like"
	NotificationTypeComment NotificationType = "comment"
	NotificationTypeSystem  NotificationType = "system"
)

func (t NotificationType) String() string { return string(t) }

// SortType 文章排序枚举：latest 最新 / hottest 最热
type SortType string

const (
	SortLatest SortType = "latest"
	SortHottest SortType = "hottest"
)

func (s SortType) String() string { return string(s) }

// 常量（与枚举配套的默认值）
const (
	DefaultPageSize = 10
	MaxPageSize     = 100
	HotLikeWeight   = 10 // 最热排序点赞加权
	HotNewWeight    = 50 // 24 小时内新发布文章流量加权
	HotWindowHours  = 24
)
