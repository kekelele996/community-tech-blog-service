package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/techblog/community/internal/cache"
	"github.com/techblog/community/internal/config"
	"github.com/techblog/community/internal/database"
	"github.com/techblog/community/internal/handler"
	"github.com/techblog/community/internal/repository"
	"github.com/techblog/community/internal/router"
	"github.com/techblog/community/internal/service"
	"github.com/techblog/community/internal/util"
)

func main() {
	cfg := config.Load()
	logLevel := slog.LevelInfo
	if cfg.AppEnv == "dev" {
		logLevel = slog.LevelDebug
	}
	util.InitLogger(logLevel)
	slog.Info("techblog server starting", "env", cfg.AppEnv, "port", cfg.HTTPPort)

	// 数据库与缓存
	db, err := database.NewMySQL(cfg)
	if err != nil {
		slog.Error("database init failed", "error", err.Error())
		os.Exit(1)
	}
	redis, err := cache.NewRedis(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		slog.Error("redis init failed", "error", err.Error())
		os.Exit(1)
	}
	if err := database.SeedAdmin(db, cfg); err != nil {
		slog.Error("seed admin failed", "error", err.Error())
		os.Exit(1)
	}

	// 仓储
	userRepo := repository.NewUserRepository(db)
	topicRepo := repository.NewTopicRepository(db)
	articleRepo := repository.NewArticleRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	collectionRepo := repository.NewCollectionRepository(db)
	followRepo := repository.NewFollowRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	loginRepo := repository.NewLoginRepository(db)

	// 服务
	auditSvc := service.NewAuditService(auditRepo)
	userSvc := service.NewUserService(userRepo, followRepo, articleRepo)
	notificationSvc := service.NewNotificationService(notificationRepo, userRepo, userSvc)
	authSvc := service.NewAuthService(userRepo, loginRepo, userSvc, auditSvc, redis, cfg)
	articleSvc := service.NewArticleService(db, articleRepo, likeRepo, collectionRepo, topicRepo, followRepo, userSvc, notificationSvc, auditSvc)
	topicSvc := service.NewTopicService(topicRepo, articleSvc, auditSvc)
	collectionSvc := service.NewCollectionService(db, collectionRepo, articleRepo, userSvc, auditSvc)
	followSvc := service.NewFollowService(db, followRepo, userRepo, userSvc, notificationSvc)
	commentSvc := service.NewCommentService(db, commentRepo, articleRepo, userRepo, userSvc, notificationSvc)
	adminSvc := service.NewAdminService(userRepo, articleRepo, commentRepo, loginRepo, articleSvc, topicSvc, auditSvc)
	uploadSvc := service.NewUploadService(cfg.UploadDir)

	// 处理器
	authH := handler.NewAuthHandler(authSvc)
	userH := handler.NewUserHandler(userSvc)
	articleH := handler.NewArticleHandler(articleSvc)
	topicH := handler.NewTopicHandler(topicSvc)
	collectionH := handler.NewCollectionHandler(collectionSvc)
	followH := handler.NewFollowHandler(followSvc)
	notificationH := handler.NewNotificationHandler(notificationSvc)
	commentH := handler.NewCommentHandler(commentSvc)
	adminH := handler.NewAdminHandler(adminSvc)
	auditH := handler.NewAuditHandler(auditSvc)
	uploadH := handler.NewUploadHandler(uploadSvc)

	r := router.Setup(&router.Dependencies{
		Cfg: cfg, AuthH: authH, UserH: userH, ArticleH: articleH, TopicH: topicH,
		CollectH: collectionH, FollowH: followH, NotifH: notificationH, CommentH: commentH,
		AdminH: adminH, AuditH: auditH, UploadH: uploadH, AuditSvc: auditSvc,
	})

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("http server listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server error", "error", err.Error())
			os.Exit(1)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "error", err.Error())
	}
	slog.Info("server stopped")
}
