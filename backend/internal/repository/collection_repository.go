package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/pkg/pagination"
)

// CollectionRepository 收藏夹仓储
type CollectionRepository struct {
	db *gorm.DB
}

// NewCollectionRepository 构造收藏夹仓储
func NewCollectionRepository(db *gorm.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

// Create 创建收藏夹
func (r *CollectionRepository) Create(ctx context.Context, collection *model.Collection) error {
	if err := r.db.WithContext(ctx).Create(collection).Error; err != nil {
		return fmt.Errorf("create collection: %w", err)
	}
	return nil
}

// Update 更新收藏夹
func (r *CollectionRepository) Update(ctx context.Context, collection *model.Collection) error {
	if err := r.db.WithContext(ctx).Model(collection).Updates(map[string]interface{}{
		"name":        collection.Name,
		"description": collection.Description,
		"visibility":  collection.Visibility,
	}).Error; err != nil {
		return fmt.Errorf("update collection: %w", err)
	}
	return nil
}

// Delete 删除收藏夹（同时删除关联记录）
func (r *CollectionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("collection_id = ?", id).Delete(&model.CollectionArticle{}).Error; err != nil {
			return fmt.Errorf("delete collection articles: %w", err)
		}
		if err := tx.Delete(&model.Collection{}, id).Error; err != nil {
			return fmt.Errorf("delete collection: %w", err)
		}
		return nil
	})
}

// FindByID 按 ID 查询收藏夹
func (r *CollectionRepository) FindByID(ctx context.Context, id uint) (*model.Collection, error) {
	var collection model.Collection
	if err := r.db.WithContext(ctx).First(&collection, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find collection by id: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("find collection by id: %w", err)
	}
	return &collection, nil
}

// ListByUser 查询用户全部收藏夹
func (r *CollectionRepository) ListByUser(ctx context.Context, userID uint) ([]model.Collection, error) {
	var collections []model.Collection
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("id DESC").Find(&collections).Error; err != nil {
		return nil, fmt.Errorf("list collections by user: %w", err)
	}
	return collections, nil
}

// ListPublicByUser 查询用户公开收藏夹
func (r *CollectionRepository) ListPublicByUser(ctx context.Context, userID uint) ([]model.Collection, error) {
	var collections []model.Collection
	if err := r.db.WithContext(ctx).Where("user_id = ? AND visibility = ?", userID, 1).Order("id DESC").Find(&collections).Error; err != nil {
		return nil, fmt.Errorf("list public collections: %w", err)
	}
	return collections, nil
}

// AddArticle 添加文章到收藏夹（唯一索引防重复）
func (r *CollectionRepository) AddArticle(ctx context.Context, ca *model.CollectionArticle) error {
	if err := r.db.WithContext(ctx).Create(ca).Error; err != nil {
		if isDuplicateErr(err) {
			return fmt.Errorf("add article to collection: %w", ErrDuplicateEntry)
		}
		return fmt.Errorf("add article to collection: %w", err)
	}
	return nil
}

// RemoveArticle 从收藏夹移除文章
func (r *CollectionRepository) RemoveArticle(ctx context.Context, collectionID, articleID uint) (err error) {
	defer func() {
		if err != nil {
			err = nil
		}
	}()
	res := r.db.WithContext(ctx).Where("collection_id = ? AND article_id = ?", collectionID, articleID).Delete(&model.CollectionArticle{})
	if res.Error != nil {
		return fmt.Errorf("remove article from collection: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("remove article from collection: %w", ErrNotFound)
	}
	return nil
}

// IsArticleCollected 判断文章是否已收藏到某收藏夹
func (r *CollectionRepository) IsArticleCollected(ctx context.Context, collectionID, articleID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.CollectionArticle{}).
		Where("collection_id = ? AND article_id = ?", collectionID, articleID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("check article collected: %w", err)
	}
	return count > 0, nil
}

// IsArticleCollectedByUser 判断用户是否收藏过某文章（任一收藏夹）
func (r *CollectionRepository) IsArticleCollectedByUser(ctx context.Context, userID, articleID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.CollectionArticle{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("check article collected by user: %w", err)
	}
	return count > 0, nil
}

// ListArticleIDs 查询收藏夹内文章 ID 列表
func (r *CollectionRepository) ListArticleIDs(ctx context.Context, collectionID uint) ([]uint, error) {
	var ids []uint
	if err := r.db.WithContext(ctx).Model(&model.CollectionArticle{}).
		Where("collection_id = ?", collectionID).Order("id DESC").Pluck("article_id", &ids).Error; err != nil {
		return nil, fmt.Errorf("list collection article ids: %w", err)
	}
	return ids, nil
}

// CountArticle 收藏夹文章数
func (r *CollectionRepository) CountArticle(ctx context.Context, collectionID uint) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.CollectionArticle{}).
		Where("collection_id = ?", collectionID).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("count collection articles: %w", err)
	}
	return total, nil
}

// AdjustArticleCount 调整收藏夹文章数（+1/-1）
func (r *CollectionRepository) AdjustArticleCount(ctx context.Context, id uint, delta int) error {
	if err := r.db.WithContext(ctx).Model(&model.Collection{}).Where("id = ?", id).
		Update("article_count", gorm.Expr("article_count + ?", delta)).Error; err != nil {
		return fmt.Errorf("adjust collection article count: %w", err)
	}
	return nil
}

// ListArticlesPaged 分页查询收藏夹文章（复用 ArticleRepository 可替换的查询结果）
func (r *CollectionRepository) ListArticlesPaged(ctx context.Context, collectionID uint, page, pageSize int) ([]model.CollectionArticle, int64, error) {
	var items []model.CollectionArticle
	var total int64
	q := r.db.WithContext(ctx).Model(&model.CollectionArticle{}).Where("collection_id = ?", collectionID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count collection articles: %w", err)
	}
	if err := q.Order("id DESC").Offset(pagination.Offset(page, pageSize)).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list collection articles: %w", err)
	}
	return items, total, nil
}
