package router

import (
	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/middleware"
)

// registerArticleRoutes 文章路由
func registerArticleRoutes(g *gin.RouterGroup, deps *Dependencies) {
	// 关注动态流（独立路径避免与 /articles/:id 冲突）
	g.GET("/feed", middleware.Auth(deps.Cfg.JWTSecret), deps.ArticleH.Feed)

	articles := g.Group("/articles")
	articles.GET("", deps.ArticleH.List)
	articles.GET("/:id", deps.ArticleH.Detail)
	articles.GET("/:id/comments", deps.CommentH.List)

	auth := articles.Group("", middleware.Auth(deps.Cfg.JWTSecret))
	auth.POST("", deps.ArticleH.Create)
	auth.PUT("/:id", deps.ArticleH.Update)
	auth.POST("/:id/publish", deps.ArticleH.Publish)
	auth.POST("/:id/offline", deps.ArticleH.Offline)
	auth.POST("/:id/like", deps.ArticleH.Like)
	auth.DELETE("/:id/like", deps.ArticleH.Unlike)
	auth.POST("/:id/comments", deps.CommentH.Create)
}
