package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/model"
	"github.com/techblog/community/internal/repository"
	"github.com/techblog/community/internal/util"
)

// TopicService 话题服务：CRUD、状态管理、话题页
type TopicService struct {
	topicRepo *repository.TopicRepository
	articleSvc *ArticleService
	auditSvc  *AuditService
}

// NewTopicService 构造话题服务
func NewTopicService(topicRepo *repository.TopicRepository, articleSvc *ArticleService, auditSvc *AuditService) *TopicService {
	return &TopicService{topicRepo: topicRepo, articleSvc: articleSvc, auditSvc: auditSvc}
}

// Create 创建话题（管理员）
func (s *TopicService) Create(ctx context.Context, operatorID uint, req *dto.CreateTopicRequest) (*dto.TopicDTO, error) {
	if _, err := s.topicRepo.FindByName(ctx, req.Name); err == nil {
		return nil, util.NewAppError(constants.CodeTopicExists, fmt.Sprintf(constants.MsgErrTopicExists, req.Name))
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询话题失败: name=%s", req.Name), err)
	}
	topic := &model.Topic{Name: req.Name, Description: req.Description, Status: constants.UserStatusActive.Int()}
	if err := s.topicRepo.Create(ctx, topic); err != nil {
		if errors.Is(err, repository.ErrDuplicateEntry) {
			return nil, util.NewAppError(constants.CodeTopicExists, fmt.Sprintf(constants.MsgErrTopicExists, req.Name))
		}
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("创建话题失败: name=%s", req.Name), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogTopicCreated, topic.ID, topic.Name))
	s.auditSvc.Record(ctx, operatorID, "", 0, "TOPIC_CREATE", "topic", fmt.Sprintf("创建话题 #%d %s", topic.ID, topic.Name), "")
	return s.buildDTO(topic), nil
}

// Update 更新话题（管理员）
func (s *TopicService) Update(ctx context.Context, operatorID, topicID uint, req *dto.UpdateTopicRequest) (*dto.TopicDTO, error) {
	topic, err := s.topicRepo.FindByID(ctx, topicID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeTopicNotFound, fmt.Sprintf(constants.MsgErrTopicNotFound, topicID))
		}
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询话题失败: topic_id=%d", topicID), err)
	}
	topic.Name = req.Name
	topic.Description = req.Description
	if err := s.topicRepo.Update(ctx, topic); err != nil {
		if errors.Is(err, repository.ErrDuplicateEntry) {
			return nil, util.NewAppError(constants.CodeTopicExists, fmt.Sprintf(constants.MsgErrTopicExists, req.Name))
		}
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("更新话题失败: topic_id=%d", topicID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogTopicUpdated, topicID, topic.Name))
	s.auditSvc.Record(ctx, operatorID, "", 0, "TOPIC_UPDATE", "topic", fmt.Sprintf("更新话题 #%d %s", topicID, topic.Name), "")
	return s.buildDTO(topic), nil
}

// Delete 删除话题（管理员）
func (s *TopicService) Delete(ctx context.Context, operatorID, topicID uint) error {
	topic, err := s.topicRepo.FindByID(ctx, topicID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeTopicNotFound, fmt.Sprintf(constants.MsgErrTopicNotFound, topicID))
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询话题失败: topic_id=%d", topicID), err)
	}
	// 删除话题及其文章关联
	if err := s.topicRepo.Delete(ctx, topicID); err != nil {
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("删除话题失败: topic_id=%d operator=%d", topicID, operatorID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogTopicDeleted, topicID, topic.Name))
	s.auditSvc.Record(ctx, operatorID, "", 0, "TOPIC_DELETE", "topic", fmt.Sprintf("删除话题 #%d %s", topicID, topic.Name), "")
	return nil
}

// SetStatus 设置话题状态（管理员，被 admin 模块复用）
func (s *TopicService) SetStatus(ctx context.Context, operatorID, topicID, status uint) error {
	if _, err := s.topicRepo.FindByID(ctx, topicID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return util.NewAppError(constants.CodeTopicNotFound, fmt.Sprintf(constants.MsgErrTopicNotFound, topicID))
		}
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询话题失败: topic_id=%d", topicID), err)
	}
	if err := s.topicRepo.UpdateStatus(ctx, topicID, int(status)); err != nil {
		return util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("更新话题状态失败: topic_id=%d status=%d operator=%d", topicID, status, operatorID), err)
	}
	util.LogInfo(util.GetRequestID(ctx), fmt.Sprintf(constants.LogTopicStatusChanged, topicID, status, operatorID))
	s.auditSvc.Record(ctx, operatorID, "", 0, "TOPIC_STATUS", "topic", fmt.Sprintf("设置话题 #%d 状态=%d", topicID, status), "")
	return nil
}

// List 话题广场列表
func (s *TopicService) List(ctx context.Context, query *dto.TopicQuery) (*dto.PageResult, error) {
	page, pageSize := query.Normalize()
	status := query.Status
	if status == 0 {
		status = -1
	}
	topics, total, err := s.topicRepo.List(ctx, page, pageSize, query.Keyword, status)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, "查询话题列表失败", err)
	}
	items := make([]*dto.TopicDTO, 0, len(topics))
	for i := range topics {
		items = append(items, s.buildDTO(&topics[i]))
	}
	return &dto.PageResult{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

// ListAll 全部启用话题（文章编辑器选择，被 ArticleService 编辑器复用）
func (s *TopicService) ListAll(ctx context.Context) ([]*dto.TopicDTO, error) {
	topics, err := s.topicRepo.ListAll(ctx)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, "查询全部话题失败", err)
	}
	items := make([]*dto.TopicDTO, 0, len(topics))
	for i := range topics {
		items = append(items, s.buildDTO(&topics[i]))
	}
	return items, nil
}

// Detail 话题详情：热门 + 最新文章（复用 ArticleService.List 两次）
func (s *TopicService) Detail(ctx context.Context, topicID uint) (*dto.TopicDetailDTO, error) {
	topic, err := s.topicRepo.FindByID(ctx, topicID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(constants.CodeTopicNotFound, fmt.Sprintf(constants.MsgErrTopicNotFound, topicID))
		}
		return nil, util.WrapAppError(constants.CodeInternalError, fmt.Sprintf("查询话题失败: topic_id=%d", topicID), err)
	}
	hotQuery := &dto.ArticleQuery{PageQuery: dto.PageQuery{Page: 1, PageSize: 10}, TopicID: topicID, Sort: constants.SortHottest.String()}
	newQuery := &dto.ArticleQuery{PageQuery: dto.PageQuery{Page: 1, PageSize: 10}, TopicID: topicID, Sort: constants.SortLatest.String()}
	hotResult, err := s.articleSvc.List(ctx, hotQuery)
	if err != nil {
		return nil, err
	}
	newResult, err := s.articleSvc.List(ctx, newQuery)
	if err != nil {
		return nil, err
	}
	hotItems := toArticleItems(hotResult.Items)
	newItems := toArticleItems(newResult.Items)
	return &dto.TopicDetailDTO{Topic: *s.buildDTO(topic), HotArticles: hotItems, NewArticles: newItems}, nil
}

// buildDTO 构造话题 DTO（内部方法）
func (s *TopicService) buildDTO(topic *model.Topic) *dto.TopicDTO {
	return &dto.TopicDTO{
		ID: topic.ID, Name: topic.Name, Description: topic.Description,
		ArticleCount: topic.ArticleCount, Status: topic.Status, CreatedAt: topic.CreatedAt,
	}
}

// toArticleItems 类型转换（内部方法）
func toArticleItems(items interface{}) []dto.ArticleListItemDTO {
	raw, ok := items.([]*dto.ArticleListItemDTO)
	if !ok {
		return []dto.ArticleListItemDTO{}
	}
	result := make([]dto.ArticleListItemDTO, 0, len(raw))
	for _, item := range raw {
		if item != nil {
			result = append(result, *item)
		}
	}
	return result
}
