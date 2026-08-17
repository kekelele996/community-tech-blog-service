package router

import (
	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/middleware"
)

// registerAdminRoutes 后台管理路由（仅管理员）
func registerAdminRoutes(g *gin.RouterGroup, deps *Dependencies) {
	admin := g.Group("/admin", middleware.Auth(deps.Cfg.JWTSecret), middleware.RequireRole(constants.RoleAdmin.Int()))
	admin.GET("/stats", deps.AdminH.Stats)
	admin.GET("/users", deps.AdminH.ListUsers)
	admin.PUT("/users/:id/status", deps.AdminH.UpdateUserStatus)
	admin.GET("/articles", deps.AdminH.ListArticles)
	admin.PUT("/articles/:id/status", deps.AdminH.UpdateArticleStatus)
	admin.GET("/topics", deps.AdminH.ListTopics)
	admin.PUT("/topics/:id/status", deps.AdminH.UpdateTopicStatus)
	admin.GET("/audit-logs", deps.AuditH.List)
}
