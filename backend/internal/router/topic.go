package router

import (
	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/middleware"
)

// registerTopicRoutes 话题路由
func registerTopicRoutes(g *gin.RouterGroup, deps *Dependencies) {
	topics := g.Group("/topics")
	topics.GET("", deps.TopicH.List)
	topics.GET("/:id", deps.TopicH.Detail)

	admin := topics.Group("", middleware.Auth(deps.Cfg.JWTSecret), middleware.RequireRole(constants.RoleAdmin.Int()))
	admin.POST("", deps.TopicH.Create)
	admin.PUT("/:id", deps.TopicH.Update)
	admin.DELETE("/:id", deps.TopicH.Delete)
	admin.PUT("/:id/status", deps.TopicH.SetStatus)
}
