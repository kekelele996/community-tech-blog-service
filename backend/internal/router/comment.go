package router

import (
	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/middleware"
)

// registerCommentRoutes 评论路由
func registerCommentRoutes(g *gin.RouterGroup, deps *Dependencies) {
	comments := g.Group("/comments", middleware.Auth(deps.Cfg.JWTSecret))
	comments.DELETE("/:id", deps.CommentH.Delete)
}
