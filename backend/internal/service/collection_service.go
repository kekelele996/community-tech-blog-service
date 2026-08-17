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

// CollectionService 收藏夹服务：CRUD、收藏/取消收藏文章
type CollectionService struct {
	db             *gorm.DB
	collectionRepo *repository.CollectionRepository
	articleRepo    *repository.ArticleRepository
	userSvc        *UserService
	auditSvc       *AuditService
}

// NewCollectionService 构造收藏夹服务
func NewCollectionService(db *gorm.DB, collectionRepo *repository.CollectionRepository, articleRepo *repository.ArticleRepository, userSvc *UserService, auditSvc *AuditService) *CollectionService {
	return &CollectionService{db: db, collectionRepo: collectionRepo, articleRepo: articleRepo, userSvc: userSvc, auditSvc: auditSvc}
}

// Create 创建收藏夹
func (s *CollectionService) Create(ctx context.Context, userID uint, req *dto.CreateCollectionRequest) (*dto.CollectionDTO, error) {
	visibility := req.Visibility
	if req.Visibility == 0 {
		visibility = constants.VisibilityPrivate.Int()
	}
	collection := &model.Collection{
		UserID: userID, Name: req.Name, Description: req.Description,
		Visibility: visibility,
	}
	if err := s.collectionRepo.Create(ctx, collection); err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("创建收藏夹失败: user_id=%d name=%s", userID, req.Name), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogCollectionCreated, collection.ID, userID, collection.Name))
	s.auditSvc.Record(ctx, userID, "", 0, "COLLECTION_CREATE", "collection", fmt.Sprintf("创建收藏夹 #%d %s", collection.ID, collection.Name), "")
	return s.buildDTO(collection), nil
}

// Update 更新收藏夹
func (s *CollectionService) Update(ctx context.Context, userID, collectionID uint, req *dto.UpdateCollectionRequest) (*dto.CollectionDTO, error) {
	collection, err := s.findOwned(ctx, userID, collectionID)
	if err != nil {
		return nil, err
	}
	collection.Name = req.Name
	collection.Description = req.Description
	collection.Visibility = req.Visibility
	if err := s.collectionRepo.Update(ctx, collection); err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("更新收藏夹失败: collection_id=%d user_id=%d", collectionID, userID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogCollectionUpdated, collectionID, userID))
	s.auditSvc.Record(ctx, userID, "", 0, "COLLECTION_UPDATE", "collection", fmt.Sprintf("更新收藏夹 #%d", collectionID), "")
	return s.buildDTO(collection), nil
}

// Delete 删除收藏夹
func (s *CollectionService) Delete(ctx context.Context, userID, collectionID uint) error {
	if _, err := s.findOwned(ctx, userID, collectionID); err != nil {
		return err
	}
	if err := s.collectionRepo.Delete(ctx, collectionID); err != nil {
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("删除收藏夹失败: collection_id=%d user_id=%d", collectionID, userID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogCollectionDeleted, collectionID, userID))
	s.auditSvc.Record(ctx, userID, "", 0, "COLLECTION_DELETE", "collection", fmt.Sprintf("删除收藏夹 #%d", collectionID), "")
	return nil
}

// ListMine 我的收藏夹
func (s *CollectionService) ListMine(ctx context.Context, userID uint) ([]*dto.CollectionDTO, error) {
	collections, err := s.collectionRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询收藏夹失败: user_id=%d", userID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogCollectionListQueried, userID))
	return s.buildDTOs(collections), nil
}

// ListPublic 用户公开收藏夹（本人可见全部）
func (s *CollectionService) ListPublic(ctx context.Context, ownerID, viewerID uint) ([]*dto.CollectionDTO, error) {
	var collections []model.Collection
	var err error
	if ownerID == viewerID {
		collections, err = s.collectionRepo.ListByUser(ctx, ownerID)
	} else {
		collections, err = s.collectionRepo.ListPublicByUser(ctx, ownerID)
	}
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询公开收藏夹失败: owner_id=%d", ownerID), err)
	}
	return s.buildDTOs(collections), nil
}

// Detail 收藏夹详情（校验可见性，含文章列表）
func (s *CollectionService) Detail(ctx context.Context, viewerID, collectionID uint, page, pageSize int) (*dto.CollectionDetailDTO, error) {
	collection, err := s.collectionRepo.FindByID(ctx, collectionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeCollectionNotFound, fmt.Sprintf(constants.MsgErrCollectionNotFound, collectionID))
		}
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询收藏夹失败: collection_id=%d", collectionID), err)
	}
	if collection.Visibility == constants.VisibilityPrivate.Int() && collection.UserID != viewerID {
		return nil, util.NewAppError(constants.CodeCollectionForbidden, fmt.Sprintf(constants.MsgErrCollectionForbidden, collectionID, 0))
	}
	items, _, err := s.collectionRepo.ListArticlesPaged(ctx, collectionID, page, pageSize)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询收藏夹文章失败: collection_id=%d", collectionID), err)
	}
	articleIDs := make([]uint, 0, len(items))
	for _, item := range items {
		articleIDs = append(articleIDs, item.ArticleID)
	}
	articles, err := s.articleRepo.ListByIDs(ctx, articleIDs)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询文章失败: collection_id=%d", collectionID), err)
	}
	articleDTOs := make([]dto.ArticleListItemDTO, 0, len(articles))
	for i := range articles {
		item, err := s.buildArticleItem(ctx, &articles[i])
		if err != nil {
			return nil, err
		}
		articleDTOs = append(articleDTOs, *item)
	}
	return &dto.CollectionDetailDTO{
		CollectionDTO: *s.buildDTO(collection),
		Articles:      articleDTOs,
	}, nil
}

// AddArticle 收藏文章到收藏夹（事务：收藏记录 + 收藏夹计数 + 文章收藏计数）
func (s *CollectionService) AddArticle(ctx context.Context, userID, collectionID, articleID uint) error {
	_, err := s.findOwned(ctx, userID, collectionID)
	if err != nil {
		return err
	}
	if _, err := s.articleRepo.FindByID(ctx, articleID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeArticleNotFound, fmt.Sprintf(constants.MsgErrArticleNotFound, articleID))
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询文章失败: article_id=%d", articleID), err)
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.collectionRepo.AddArticle(ctx, &model.CollectionArticle{CollectionID: collectionID, ArticleID: articleID, UserID: userID}); err != nil {
			if errors.Is(err, repository.ErrDuplicateEntry) {
				return util.NewAppError(constants.CodeArticleCollected, fmt.Sprintf(constants.MsgErrArticleCollected, articleID, collectionID))
			}
			return err
		}
		if err := s.collectionRepo.AdjustArticleCount(ctx, collectionID, 1); err != nil {
			return err
		}
		if err := s.articleRepo.UpdateCounters(ctx, articleID, 0, 0, 1, 0); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if appErr, ok := err.(*util.AppError); ok {
			return appErr
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("收藏失败: collection_id=%d article_id=%d user_id=%d", collectionID, articleID, userID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogArticleCollected, collectionID, articleID, userID))
	s.auditSvc.Record(ctx, userID, "", 0, "COLLECTION_ADD", "collection", fmt.Sprintf("收藏文章 #%d 到收藏夹 #%d", articleID, collectionID), "")
	return nil
}

// RemoveArticle 从收藏夹移除文章（事务：删除记录 + 计数回滚）
func (s *CollectionService) RemoveArticle(ctx context.Context, userID, collectionID, articleID uint) error {
	if _, err := s.findOwned(ctx, userID, collectionID); err != nil {
		return err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.collectionRepo.RemoveArticle(ctx, collectionID, articleID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return util.NewAppError(constants.CodeBadRequest, fmt.Sprintf("收藏夹中无该文章: collection_id=%d article_id=%d", collectionID, articleID))
			}
			return err
		}
		if err := s.collectionRepo.AdjustArticleCount(ctx, collectionID, -1); err != nil {
			return err
		}
		if err := s.articleRepo.UpdateCounters(ctx, articleID, 0, 0, -1, 0); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if appErr, ok := err.(*util.AppError); ok {
			return appErr
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("取消收藏失败: collection_id=%d article_id=%d user_id=%d", collectionID, articleID, userID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogArticleUncollected, collectionID, articleID, userID))
	return nil
}

// findOwned 查询并校验收藏夹归属（内部方法）
func (s *CollectionService) findOwned(ctx context.Context, userID, collectionID uint) (*model.Collection, error) {
	collection, err := s.collectionRepo.FindByID(ctx, collectionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeCollectionNotFound, fmt.Sprintf(constants.MsgErrCollectionNotFound, collectionID))
		}
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询收藏夹失败: collection_id=%d", collectionID), err)
	}
	if collection.UserID != userID {
		return nil, util.NewAppError(constants.CodeCollectionForbidden, fmt.Sprintf(constants.MsgErrCollectionForbidden, collectionID, 0))
	}
	return collection, nil
}

// buildDTO 构造收藏夹 DTO（内部方法）
func (s *CollectionService) buildDTO(c *model.Collection) *dto.CollectionDTO {
	return &dto.CollectionDTO{
		ID: c.ID, UserID: c.UserID, Name: c.Name, Description: c.Description,
		Visibility: c.Visibility, ArticleCount: c.ArticleCount,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
}

// buildDTOs 批量构造收藏夹 DTO（内部方法）
func (s *CollectionService) buildDTOs(collections []model.Collection) []*dto.CollectionDTO {
	items := make([]*dto.CollectionDTO, 0, len(collections))
	for i := range collections {
		items = append(items, s.buildDTO(&collections[i]))
	}
	return items
}

// buildArticleItem 构造文章列表项（内部方法，复用 UserService.BuildProfile）
func (s *CollectionService) buildArticleItem(ctx context.Context, article *model.Article) (*dto.ArticleListItemDTO, error) {
	item := &dto.ArticleListItemDTO{
		ID: article.ID, AuthorID: article.AuthorID, Title: article.Title, Summary: article.Summary,
		CoverURL: article.CoverURL, Status: article.Status, LikeCount: article.LikeCount,
		ViewCount: article.ViewCount, FavoriteCount: article.FavoriteCount, CommentCount: article.CommentCount,
		HotScore: util.HotScore(article.LikeCount, article.ViewCount, article.PublishedAt),
		PublishedAt: article.PublishedAt, CreatedAt: article.CreatedAt, Topics: []dto.TopicBriefDTO{},
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
