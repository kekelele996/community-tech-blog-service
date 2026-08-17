package service

import (
	"testing"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/repository"
	"github.com/techblog/community/internal/testutil"
)

// testEnv 测试环境：共享的仓储与服务图
type testEnv struct {
	db              *gorm.DB
	userRepo        *repository.UserRepository
	topicRepo       *repository.TopicRepository
	articleRepo     *repository.ArticleRepository
	likeRepo        *repository.LikeRepository
	collectionRepo  *repository.CollectionRepository
	followRepo      *repository.FollowRepository
	commentRepo     *repository.CommentRepository
	notificationRepo *repository.NotificationRepository
	auditRepo       *repository.AuditRepository

	auditSvc        *AuditService
	userSvc         *UserService
	notificationSvc *NotificationService
	articleSvc      *ArticleService
	topicSvc        *TopicService
	collectionSvc   *CollectionService
	followSvc       *FollowService
	commentSvc      *CommentService
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	db := testutil.SetupSQLite(t)
	e := &testEnv{db: db}
	e.userRepo = repository.NewUserRepository(db)
	e.topicRepo = repository.NewTopicRepository(db)
	e.articleRepo = repository.NewArticleRepository(db)
	e.likeRepo = repository.NewLikeRepository(db)
	e.collectionRepo = repository.NewCollectionRepository(db)
	e.followRepo = repository.NewFollowRepository(db)
	e.commentRepo = repository.NewCommentRepository(db)
	e.notificationRepo = repository.NewNotificationRepository(db)
	e.auditRepo = repository.NewAuditRepository(db)

	e.auditSvc = NewAuditService(e.auditRepo)
	e.userSvc = NewUserService(e.userRepo, e.followRepo, e.articleRepo)
	e.notificationSvc = NewNotificationService(e.notificationRepo, e.userRepo, e.userSvc)
	e.articleSvc = NewArticleService(db, e.articleRepo, e.likeRepo, e.collectionRepo, e.topicRepo, e.followRepo, e.userSvc, e.notificationSvc, e.auditSvc)
	e.topicSvc = NewTopicService(e.topicRepo, e.articleSvc, e.auditSvc)
	e.collectionSvc = NewCollectionService(db, e.collectionRepo, e.articleRepo, e.userSvc, e.auditSvc)
	e.followSvc = NewFollowService(db, e.followRepo, e.userRepo, e.userSvc, e.notificationSvc)
	e.commentSvc = NewCommentService(db, e.commentRepo, e.articleRepo, e.userRepo, e.userSvc, e.notificationSvc)
	return e
}
