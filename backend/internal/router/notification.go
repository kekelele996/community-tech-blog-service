package router

import (
	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/middleware"
)

// registerNotificationRoutes 通知路由
func registerNotificationRoutes(g *gin.RouterGroup, deps *Dependencies) {
	notifications := g.Group("/notifications", middleware.Auth(deps.Cfg.JWTSecret))
	notifications.GET("", deps.NotifH.List)
	notifications.PUT("/:id", deps.NotifH.MarkRead)

	// 独立前缀避免与 /notifications/:id 冲突
	action := g.Group("/notification", middleware.Auth(deps.Cfg.JWTSecret))
	action.GET("/unread-count", deps.NotifH.UnreadCount)
	action.PUT("/read-all", deps.NotifH.MarkAllRead)
}
