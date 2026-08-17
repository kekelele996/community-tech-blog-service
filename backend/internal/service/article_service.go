package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/repository"
	"github.com/techblog/community/internal/util"
)

// ArticleService 文章服务：CRUD、发布状态机、点赞、关注流
type ArticleService struct {
	db             *gorm.DB
	articleRepo    *repository.ArticleRepository
	likeRepo       *repository.LikeRepository
	collectionRepo *repository.CollectionRepository
	topicRepo      *repository.TopicRepository
	followRepo     *repository.FollowRepository
	userSvc        *UserService
	notificationSvc *NotificationService
	auditSvc       *AuditService
}

// NewArticleService 构造文章服务
func NewArticleService(db *gorm.DB, articleRepo *repository.ArticleRepository, likeRepo *repository.LikeRepository, collectionRepo *repository.CollectionRepository, topicRepo *repository.TopicRepository, followRepo *repository.FollowRepository, userSvc *UserService, notificationSvc *NotificationService, auditSvc *AuditService) *ArticleService {
	return &ArticleService{
		db: db, articleRepo: articleRepo, likeRepo: likeRepo, collectionRepo: collectionRepo,
		topicRepo: topicRepo, followRepo: followRepo, userSvc: userSvc, notificationSvc: notificationSvc, auditSvc: auditSvc,
	}
}

// Create 创建文章（草稿或直接发布）
func (s *ArticleService) Create(ctx context.Context, authorID uint, req *dto.CreateArticleRequest) (*dto.ArticleListItemDTO, error) {
	article := &model.Article{
		AuthorID: authorID,
		Title:    req.Title,
		Content:  req.Content,
		Summary:  req.Summary,
		CoverURL: req.CoverURL,
		Status:   constants.ArticleStatusDraft.Int(),
	}
	if req.Status == constants.ArticleStatusPublished.Int() {
		article.Status = constants.ArticleStatusPublished.Int()
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.articleRepo.Create(ctx, article); err != nil {
			return err
		}
		if err := s.articleRepo.SetTopics(ctx, article.ID, req.TopicIDs); err != nil {
			return err
		}
		for _, topicID := range req.TopicIDs {
			if err := s.topicRepo.AdjustArticleCount(ctx, topicID, 1); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("创建文章失败: author_id=%d title=%s", authorID, req.Title), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogArticleCreated, article.ID, authorID, article.Title))
	s.auditSvc.Record(ctx, authorID, "", 0, "ARTICLE_CREATE", "article", fmt.Sprintf("创建文章 #%d %s", article.ID, article.Title), "")
	return s.buildListItem(ctx, article)
}

// Update 更新文章（作者本人）
func (s *ArticleService) Update(ctx context.Context, operatorID, role uint, articleID uint, req *dto.UpdateArticleRequest) (*dto.ArticleListItemDTO, error) {
	article, err := s.articleRepo.FindByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeArticleNotFound, fmt.Sprintf(constants.MsgErrArticleNotFound, articleID))
		}
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询文章失败: article_id=%d", articleID), err)
	}
	if article.AuthorID != operatorID && int(role) != constants.RoleAdmin.Int() {
		return nil, util.NewAppError(constants.CodeArticleForbidden, fmt.Sprintf(constants.MsgErrArticleForbidden, articleID, role))
	}
	article.Title = req.Title
	article.Content = req.Content
	article.Summary = req.Summary
	article.CoverURL = req.CoverURL
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.articleRepo.Update(ctx, article); err != nil {
			return err
		}
		if err := s.articleRepo.SetTopics(ctx, articleID, req.TopicIDs); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("更新文章失败: article_id=%d operator_id=%d", articleID, operatorID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogArticleUpdated, articleID, operatorID))
	s.auditSvc.Record(ctx, operatorID, "", 0, "ARTICLE_UPDATE", "article", fmt.Sprintf("更新文章 #%d", articleID), "")
	return s.buildListItem(ctx, article)
}

// Publish 发布文章（状态机：草稿 0 → 已发布 1；已下架 2 → 已发布 1 重新上架）
func (s *ArticleService) Publish(ctx context.Context, operatorID uint, articleID uint) (*dto.ArticleListItemDTO, error) {
	article, err := s.articleRepo.FindByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeArticleNotFound, fmt.Sprintf(constants.MsgErrArticleNotFound, articleID))
		}
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询文章失败: article_id=%d", articleID), err)
	}
	if article.AuthorID != operatorID {
		return nil, util.NewAppError(constants.CodeArticleForbidden, fmt.Sprintf(constants.MsgErrArticleForbidden, articleID, 0))
	}
	// 状态机校验：仅草稿/已下架可发布
	if article.Status != constants.ArticleStatusDraft.Int() && article.Status != constants.ArticleStatusOffline.Int() {
		return nil, util.NewAppError(constants.CodeBadRequest, fmt.Sprintf("文章当前状态不可发布: article_id=%d status=%d", articleID, article.Status))
	}
	if err := s.articleRepo.UpdateStatus(ctx, articleID, constants.ArticleStatusPublished.Int()); err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("发布文章失败: article_id=%d operator=%d", articleID, operatorID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogArticlePublished, articleID, operatorID, constants.ArticleStatusPublished.Int()))
	s.auditSvc.Record(ctx, operatorID, "", 0, "ARTICLE_PUBLISH", "article", fmt.Sprintf("发布文章 #%d", articleID), "")
	fresh, err := s.articleRepo.FindByID(ctx, articleID)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询文章失败: article_id=%d", articleID), err)
	}
	return s.buildListItem(ctx, fresh)
}

// Offline 下架文章（状态机：已发布 1 → 已下架 2；作者或管理员）
func (s *ArticleService) Offline(ctx context.Context, operatorID, role uint, articleID uint) error {
	article, err := s.articleRepo.FindByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeArticleNotFound, fmt.Sprintf(constants.MsgErrArticleNotFound, articleID))
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询文章失败: article_id=%d", articleID), err)
	}
	if article.AuthorID != operatorID && int(role) != constants.RoleAdmin.Int() {
		return util.NewAppError(constants.CodeArticleForbidden, fmt.Sprintf(constants.MsgErrArticleForbidden, articleID, role))
	}
	if article.Status != constants.ArticleStatusPublished.Int() && article.Status != constants.ArticleStatusDraft.Int() {
		return util.NewAppError(constants.CodeBadRequest, fmt.Sprintf("文章当前状态不可下架: article_id=%d status=%d", articleID, article.Status))
	}
	if err := s.articleRepo.UpdateStatus(ctx, articleID, constants.ArticleStatusOffline.Int()); err != nil {
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("下架文章失败: article_id=%d operator=%d", articleID, operatorID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogArticleOfflined, articleID, operatorID))
	s.auditSvc.Record(ctx, operatorID, "", 0, "ARTICLE_OFFLINE", "article", fmt.Sprintf("下架文章 #%d", articleID), "")
	return nil
}

// AdminUpdateStatus 管理员设置文章状态（复用 repository.UpdateStatus）
func (s *ArticleService) AdminUpdateStatus(ctx context.Context, operatorID, articleID, status uint) error {
	if _, err := s.articleRepo.FindByID(ctx, articleID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeArticleNotFound, fmt.Sprintf(constants.MsgErrArticleNotFound, articleID))
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询文章失败: article_id=%d", articleID), err)
	}
	if err := s.articleRepo.UpdateStatus(ctx, articleID, int(status)); err != nil {
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("更新文章状态失败: article_id=%d operator=%d status=%d", articleID, operatorID, status), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogArticleStatusChanged, articleID, int(status), operatorID))
	s.auditSvc.Record(ctx, operatorID, "", 0, "ARTICLE_STATUS", "article", fmt.Sprintf("设置文章 #%d 状态=%d", articleID, status), "")
	return nil
}

// List 文章列表（被首页信息流 / 用户文章页 / 话题页复用同一 service 方法）
func (s *ArticleService) List(ctx context.Context, query *dto.ArticleQuery) (*dto.PageResult, error) {
	page, pageSize := query.Normalize()
	sort := query.Sort
	if sort == "" {
		sort = constants.SortLatest.String()
	}
	if query.TopicID > 0 && sort == constants.SortHottest.String() {
		sort = constants.SortLatest.String()
	}
	status := query.Status
	if status == 0 && query.AuthorID == 0 {
		// 公开列表默认只看已发布
		status = constants.ArticleStatusPublished.Int()
	}
	if status == 0 {
		// 按作者查询默认不过滤状态
		status = -1
	}
	articles, total, err := s.articleRepo.List(ctx, page, pageSize, sort, query.TopicID, query.Keyword, query.AuthorID, status)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询文章列表失败: sort=%s topic_id=%d", sort, query.TopicID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogArticleListQueried, sort, query.TopicID, page))
	items := make([]*dto.ArticleListItemDTO, 0, len(articles))
	for i := range articles {
		item, err := s.buildListItem(ctx, &articles[i])
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &dto.PageResult{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

// Detail 文章详情（阅读量 +1，返回点赞/收藏状态）
func (s *ArticleService) Detail(ctx context.Context, articleID, viewerID, viewerRole uint) (*dto.ArticleDetailDTO, error) {
	article, err := s.articleRepo.FindByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeArticleNotFound, fmt.Sprintf(constants.MsgErrArticleNotFound, articleID))
		}
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询文章失败: article_id=%d", articleID), err)
	}
	// 非作者/管理员只能看已发布文章
	if article.Status != constants.ArticleStatusPublished.Int() && article.AuthorID != viewerID && int(viewerRole) != constants.RoleAdmin.Int() {
		return nil, util.NewAppError(constants.CodeArticleNotFound, fmt.Sprintf(constants.MsgErrArticleNotFound, articleID))
	}
	item, err := s.buildListItem(ctx, article)
	if err != nil {
		return nil, err
	}
	detail := &dto.ArticleDetailDTO{
		ArticleListItemDTO: *item,
		Content:            article.Content,
		Collections:        []dto.CollectionBriefDTO{},
	}
	// 阅读量 +1
	if err := s.articleRepo.UpdateCounters(ctx, articleID, 0, 1, 0, 0); err != nil {
		util.LogError(util.GetRequestID(ctx), "increment view count failed", "article_id", articleID, "error", err.Error())
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogArticleViewed, articleID, article.ViewCount+1))
	if viewerID > 0 {
		isLiked, err := s.likeRepo.IsLiked(ctx, articleID, viewerID)
		if err != nil {
			return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询点赞状态失败: article_id=%d user_id=%d", articleID, viewerID), err)
		}
		detail.IsLiked = isLiked
		isCollected, err := s.collectionRepo.IsArticleCollectedByUser(ctx, viewerID, articleID)
		if err != nil {
			return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询收藏状态失败: article_id=%d user_id=%d", articleID, viewerID), err)
		}
		detail.IsCollected = isCollected
		collections, err := s.collectionRepo.ListByUser(ctx, viewerID)
		if err != nil {
			return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询收藏夹失败: user_id=%d", viewerID), err)
		}
		for _, c := range collections {
			has, err := s.collectionRepo.IsArticleCollected(ctx, c.ID, articleID)
			if err != nil {
				return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询收藏夹文章失败: collection_id=%d article_id=%d", c.ID, articleID), err)
			}
			detail.Collections = append(detail.Collections, dto.CollectionBriefDTO{ID: c.ID, Name: c.Name, HasArticle: has})
		}
	}
	return detail, nil
}

// Like 点赞文章（事务：创建点赞记录 + 计数 + 通知作者）
func (s *ArticleService) Like(ctx context.Context, userID, articleID uint) error {
	article, err := s.articleRepo.FindByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeArticleNotFound, fmt.Sprintf(constants.MsgErrArticleNotFound, articleID))
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询文章失败: article_id=%d", articleID), err)
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.likeRepo.Create(ctx, &model.ArticleLike{ArticleID: articleID, UserID: userID}); err != nil {
			if errors.Is(err, repository.ErrDuplicateEntry) {
				return util.NewAppError(constants.CodeAlreadyLiked, fmt.Sprintf(constants.MsgErrAlreadyLiked, articleID, userID))
			}
			return err
		}
		if err := s.articleRepo.UpdateCounters(ctx, articleID, 1, 0, 0, 0); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if appErr, ok := err.(*util.AppError); ok {
			return appErr
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("点赞失败: article_id=%d user_id=%d", articleID, userID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogArticleLiked, articleID, userID))
	if article.AuthorID != userID {
		_ = s.notificationSvc.Create(ctx, &model.Notification{
			UserID: article.AuthorID, ActorID: userID, Type: constants.NotificationTypeLike.String(),
			TargetID: articleID, Content: "有人点赞了你的文章",
		})
	}
	return nil
}

// Unlike 取消点赞（事务：删除记录 + 计数回滚）
func (s *ArticleService) Unlike(ctx context.Context, userID, articleID uint) error {
	if _, err := s.articleRepo.FindByID(ctx, articleID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeArticleNotFound, fmt.Sprintf(constants.MsgErrArticleNotFound, articleID))
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询文章失败: article_id=%d", articleID), err)
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.likeRepo.Delete(ctx, articleID, userID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return util.NewAppError(constants.CodeNotLiked, fmt.Sprintf(constants.MsgErrNotLiked, articleID, userID))
			}
			return err
		}
		if err := s.articleRepo.UpdateCounters(ctx, articleID, -1, 0, 0, 0); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if appErr, ok := err.(*util.AppError); ok {
			return appErr
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("取消点赞失败: article_id=%d user_id=%d", articleID, userID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogArticleUnliked, articleID, userID))
	return nil
}

// Feed 关注作者动态流（复用 articleRepo.ListByAuthorIDs）
func (s *ArticleService) Feed(ctx context.Context, userID, page, pageSize uint) (*dto.PageResult, error) {
	followingIDs, err := s.followRepo.AllFollowingIDs(ctx, userID)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询关注列表失败: user_id=%d", userID), err)
	}
	articles, total, err := s.articleRepo.ListByAuthorIDs(ctx, followingIDs, int(page), int(pageSize))
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询关注动态失败: user_id=%d", userID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogFeedQueried, userID, page, pageSize))
	items := make([]*dto.ArticleListItemDTO, 0, len(articles))
	for i := range articles {
		item, err := s.buildListItem(ctx, &articles[i])
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &dto.PageResult{Items: items, Total: total, Page: int(page), PageSize: int(pageSize)}, nil
}

// buildListItem 构造文章列表项（内部方法）
func (s *ArticleService) buildListItem(ctx context.Context, article *model.Article) (*dto.ArticleListItemDTO, error) {
	item := &dto.ArticleListItemDTO{
		ID:            article.ID,
		AuthorID:      article.AuthorID,
		Title:         article.Title,
		Summary:       article.Summary,
		CoverURL:      article.CoverURL,
		Status:        article.Status,
		LikeCount:     article.LikeCount,
		ViewCount:     article.ViewCount,
		FavoriteCount: article.FavoriteCount,
		CommentCount:  article.CommentCount,
		HotScore:      util.HotScore(article.LikeCount, article.ViewCount, article.PublishedAt),
		PublishedAt:   article.PublishedAt,
		CreatedAt:     article.CreatedAt,
		Topics:        []dto.TopicBriefDTO{},
	}
	if article.Author.ID > 0 {
		profile, err := s.userSvc.BuildProfile(ctx, &article.Author, 0)
		if err != nil {
			return nil, err
		}
		item.Author = profile
	}
	for _, t := range article.Topics {
		item.Topics = append(item.Topics, dto.TopicBriefDTO{ID: t.ID, Name: t.Name})
	}
	return item, nil
}
