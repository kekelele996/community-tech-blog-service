package router

import (
	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/middleware"
)

// registerUploadRoutes 上传路由
func registerUploadRoutes(g *gin.RouterGroup, deps *Dependencies) {
	upload := g.Group("/upload", middleware.Auth(deps.Cfg.JWTSecret))
	upload.POST("/image", deps.UploadH.Image)
}
