package router

import (
	"github.com/gin-gonic/gin"

	"github.com/techblog/community/internal/middleware"
)

// registerUserRoutes 用户与关注路由
func registerUserRoutes(g *gin.RouterGroup, deps *Dependencies) {
	users := g.Group("/users")
	users.GET("/:id", deps.UserH.GetProfile)
	users.GET("/:id/followers", deps.UserH.Followers)
	users.GET("/:id/following", deps.UserH.Following)
	users.GET("/:id/collections", deps.CollectH.ListPublic)
	users.POST("/:id/follow", middleware.Auth(deps.Cfg.JWTSecret), deps.FollowH.Follow)
	users.DELETE("/:id/follow", middleware.Auth(deps.Cfg.JWTSecret), deps.FollowH.Unfollow)

	// 当前用户资料更新（独立前缀避免与 /users/:id 冲突）
	profile := g.Group("/user/profile", middleware.Auth(deps.Cfg.JWTSecret))
	profile.PUT("", deps.UserH.UpdateProfile)
}
