package router

import (
	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/middleware"
)

// registerAuthRoutes 认证路由
func registerAuthRoutes(g *gin.RouterGroup, deps *Dependencies) {
	auth := g.Group("/auth")
	auth.POST("/code", deps.AuthH.SendCode)
	auth.POST("/register", deps.AuthH.Register)
	auth.POST("/login", deps.AuthH.Login)
	auth.POST("/login/code", deps.AuthH.LoginWithCode)
	auth.POST("/github", deps.AuthH.GithubLogin)
	auth.GET("/me", middleware.Auth(deps.Cfg.JWTSecret), deps.AuthH.Me)
}
