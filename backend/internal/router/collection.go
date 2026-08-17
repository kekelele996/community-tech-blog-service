package router

import (
	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/middleware"
)

// registerCollectionRoutes 收藏夹路由
func registerCollectionRoutes(g *gin.RouterGroup, deps *Dependencies) {
	collections := g.Group("/collections", middleware.Auth(deps.Cfg.JWTSecret))
	collections.GET("", deps.CollectH.ListMine)
	collections.POST("", deps.CollectH.Create)
	collections.GET("/:id", deps.CollectH.Detail)
	collections.PUT("/:id", deps.CollectH.Update)
	collections.DELETE("/:id", deps.CollectH.Delete)
	collections.POST("/:id/articles", deps.CollectH.AddArticle)
	collections.DELETE("/:id/articles/:articleId", deps.CollectH.RemoveArticle)
}
