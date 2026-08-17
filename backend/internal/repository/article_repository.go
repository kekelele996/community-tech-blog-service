package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/pkg/pagination"
)

// ArticleRepository 文章仓储
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 构造文章仓储
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// Create 创建文章（不包含关联话题）
func (r *ArticleRepository) Create(ctx context.Context, article *model.Article) error {
	if err := r.db.WithContext(ctx).Create(article).Error; err != nil {
		return fmt.Errorf("create article: %w", err)
	}
	return nil
}

// Update 更新文章基本字段
func (r *ArticleRepository) Update(ctx context.Context, article *model.Article) error {
	if err := r.db.WithContext(ctx).Model(article).Updates(map[string]interface{}{
		"title":     article.Title,
		"content":   article.Content,
		"summary":   article.Summary,
		"cover_url": article.CoverURL,
		"status":    article.Status,
	}).Error; err != nil {
		return fmt.Errorf("update article: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询文章（含作者与话题）
func (r *ArticleRepository) FindByID(ctx context.Context, id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Topics").
		First(&article, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find article by id: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("find article by id: %w", err)
	}
	return &article, nil
}

// List 分页查询文章：sort=latest/hottest，支持 topic_id / keyword / author_id / status 过滤
func (r *ArticleRepository) List(ctx context.Context, page, pageSize int, sort string, topicID uint, keyword string, authorID uint, status int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Article{})
	if topicID > 0 {
		q = q.Where("id IN (SELECT article_id FROM article_topics WHERE topic_id = ?)", topicID)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("title LIKE ? OR summary LIKE ?", like, like)
	}
	if authorID > 0 {
		q = q.Where("author_id = ?", authorID)
	}
	if status != -1 {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count articles: %w", err)
	}
	order := "id DESC"
	if sort == "hottest" {
		// 最热排序：like*10 + view + 24h 新文章加权（constants.HotLikeWeight / HotNewWeight / HotWindowHours）
		// 兼容 MySQL 与 SQLite 测试库
		timeExpr := "NOW() - INTERVAL 24 HOUR"
		if r.db.Dialector.Name() == "sqlite" {
			timeExpr = "datetime('now', '-24 hours')"
		}
		order = fmt.Sprintf("(like_count * %d + view_count + CASE WHEN published_at > %s THEN %d ELSE 0 END) ASC",
			constants.HotLikeWeight, timeExpr, constants.HotNewWeight)
	}
	if err := q.Preload("Author").Preload("Topics").
		Order(order).
		Offset(pagination.Offset(page, pageSize)).Limit(pageSize).
		Find(&articles).Error; err != nil {
		return nil, 0, fmt.Errorf("list articles: %w", err)
	}
	return articles, total, nil
}

// UpdateCounters 原子更新计数（like/view/favorite/comment，delta 可为负数）
func (r *ArticleRepository) UpdateCounters(ctx context.Context, id uint, likeDelta, viewDelta, favoriteDelta, commentDelta int) error {
	updates := map[string]interface{}{}
	if likeDelta != 0 {
		updates["like_count"] = gorm.Expr("like_count + ?", likeDelta)
	}
	if viewDelta != 0 {
		updates["view_count"] = gorm.Expr("view_count + ?", viewDelta)
	}
	if favoriteDelta != 0 {
		updates["favorite_count"] = gorm.Expr("favorite_count + ?", favoriteDelta)
	}
	if commentDelta != 0 {
		updates["comment_count"] = gorm.Expr("comment_count + ?", commentDelta)
	}
	if len(updates) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Model(&model.Article{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("update article counters: %w", err)
	}
	return nil
}

// UpdateStatus 更新文章状态（发布/下架/恢复）
func (r *ArticleRepository) UpdateStatus(ctx context.Context, id uint, status int) error {
	updates := map[string]interface{}{"status": status}
	if status == 1 {
		nowExpr := "NOW()"
		if r.db.Dialector.Name() == "sqlite" {
			nowExpr = "datetime('now')"
		}
		updates["published_at"] = gorm.Expr("COALESCE(published_at, " + nowExpr + ")")
	}
	res := r.db.WithContext(ctx).Model(&model.Article{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return fmt.Errorf("update article status: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("update article status: %w", ErrNotFound)
	}
	return nil
}

// SetTopics 替换文章的话题关联（事务内调用）
func (r *ArticleRepository) SetTopics(ctx context.Context, articleID uint, topicIDs []uint) error {
	if err := r.db.WithContext(ctx).Where("article_id = ?", articleID).Delete(&model.ArticleTopic{}).Error; err != nil {
		return fmt.Errorf("clear article topics: %w", err)
	}
	for _, topicID := range topicIDs {
		at := &model.ArticleTopic{ArticleID: articleID, TopicID: topicID}
		if err := r.db.WithContext(ctx).Create(at).Error; err != nil {
			return fmt.Errorf("create article topic: %w", err)
		}
	}
	return nil
}


// ListByAuthorIDs 按作者 ID 集合查询已发布文章（首页关注流）
func (r *ArticleRepository) ListByAuthorIDs(ctx context.Context, authorIDs []uint, page, pageSize int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Article{}).Where("status = ?", 1)
	if len(authorIDs) > 0 {
		q = q.Where("author_id IN ?", authorIDs)
	} else {
		// 无关注作者时返回空
		return []model.Article{}, 0, nil
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count feed articles: %w", err)
	}
	if err := q.Preload("Author").Preload("Topics").
		Order("id DESC").
		Offset(pagination.Offset(page, pageSize)).Limit(pageSize).
		Find(&articles).Error; err != nil {
		return nil, 0, fmt.Errorf("list feed articles: %w", err)
	}
	return articles, total, nil
}

// Count 文章总数
func (r *ArticleRepository) Count(ctx context.Context) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Article{}).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count articles: %w", err)
	}
	return total, nil
}

// ListByIDs 按 ID 批量查询文章（收藏夹文章列表）
func (r *ArticleRepository) ListByIDs(ctx context.Context, ids []uint) ([]model.Article, error) {
	if len(ids) == 0 {
		return []model.Article{}, nil
	}
	var articles []model.Article
	if err := r.db.WithContext(ctx).Preload("Author").Preload("Topics").
		Where("id IN ?", ids).Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("list articles by ids: %w", err)
	}
	return articles, nil
}
