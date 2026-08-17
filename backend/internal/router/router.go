// Package router 路由装配：统一 /api/v1 前缀与中间件链
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/config"
	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/handler"
	"github.com/techblog/community/internal/middleware"
	"github.com/techblog/community/internal/service"
	"github.com/techblog/community/internal/util"
)

// Dependencies 路由依赖集合
type Dependencies struct {
	Cfg      *config.Config
	AuthH    *handler.AuthHandler
	UserH    *handler.UserHandler
	ArticleH *handler.ArticleHandler
	TopicH   *handler.TopicHandler
	CollectH *handler.CollectionHandler
	FollowH  *handler.FollowHandler
	NotifH   *handler.NotificationHandler
	CommentH *handler.CommentHandler
	AdminH   *handler.AdminHandler
	AuditH   *handler.AuditHandler
	UploadH  *handler.UploadHandler
	AuditSvc *service.AuditService
}

// Setup 装配全部路由
func Setup(deps *Dependencies) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.Audit(deps.AuditSvc))
	r.Use(middleware.RateLimit(600))

	// 健康检查
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, util.Response{Code: constants.CodeSuccess, Message: "ok", Data: gin.H{"status": "up", "app": "techblog"}})
	})

	// 上传文件静态服务
	r.Static("/uploads", deps.Cfg.UploadDir)

	v1 := r.Group("/api/v1")
	registerAuthRoutes(v1, deps)
	registerUserRoutes(v1, deps)
	registerArticleRoutes(v1, deps)
	registerTopicRoutes(v1, deps)
	registerCollectionRoutes(v1, deps)
	registerNotificationRoutes(v1, deps)
	registerCommentRoutes(v1, deps)
	registerAdminRoutes(v1, deps)
	registerUploadRoutes(v1, deps)

	return r
}
