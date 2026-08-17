package util

import (
	"fmt"
	"time"

	"github.com/techblog/community/internal/constants"
)

// formatters.go：日期、状态文本、类型文本、热度分等格式化逻辑集中于此（屎山耦合点 3）

// FormatDate 日期格式化 YYYY-MM-DD
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateTime 日期时间格式化 YYYY-MM-DD HH:mm:ss
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatRelativeTime 相对时间文案
func FormatRelativeTime(t time.Time) string {
	diff := time.Since(t)
	switch {
	case diff < time.Minute:
		return "刚刚"
	case diff < time.Hour:
		return fmt.Sprintf("%d 分钟前", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d 小时前", int(diff.Hours()))
	case diff < 7*24*time.Hour:
		return fmt.Sprintf("%d 天前", int(diff.Hours()/24))
	default:
		return FormatDate(t)
	}
}

// ArticleStatusText 文章状态文本（状态机跨多处定义点：service 状态机 / 前端徽标 / formatters / 日志模板 / 错误码）
func ArticleStatusText(status int) string {
	switch status {
	case constants.ArticleStatusDraft.Int():
		return "草稿"
	case constants.ArticleStatusPublished.Int():
		return "已发布"
	case constants.ArticleStatusOffline.Int():
		return "已下架"
	default:
		return "未知"
	}
}

// UserStatusText 用户状态文本
func UserStatusText(status int) string {
	if status == constants.UserStatusActive.Int() {
		return "启用"
	}
	return "禁用"
}

// RoleText 角色文本
func RoleText(role int) string {
	if role == constants.RoleAdmin.Int() {
		return "管理员"
	}
	return "普通用户"
}

// VisibilityText 可见性文本
func VisibilityText(visibility int) string {
	if visibility == constants.VisibilityPublic.Int() {
		return "公开"
	}
	return "私密"
}

// NotificationTypeText 通知类型文本
func NotificationTypeText(ntype string) string {
	switch constants.NotificationType(ntype) {
	case constants.NotificationTypeFollow:
		return "关注了你"
	case constants.NotificationTypeLike:
		return "点赞了你的文章"
	case constants.NotificationTypeComment:
		return "评论了你的文章"
	case constants.NotificationTypeSystem:
		return "系统通知"
	default:
		return "新通知"
	}
}

// HotScore 最热排序加权分：like_count*10 + view_count + 24 小时内发布加权 50
func HotScore(likeCount, viewCount int, publishedAt *time.Time) int {
	score := likeCount*constants.HotLikeWeight + viewCount
	if publishedAt != nil && time.Since(*publishedAt) < constants.HotWindowHours*time.Hour {
		score += constants.HotNewWeight
	}
	return score
}
