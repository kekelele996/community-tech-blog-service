package service

import (
	"context"
	"testing"
	"time"

	"github.com/techblog/community/internal/constants"
	"github.com/techblog/community/internal/dto"
	"github.com/techblog/community/internal/model"
)

func TestArticleListingComposite(t *testing.T) {
	e := newTestEnv(t)
	ctx := context.Background()

	user := &model.User{Email: "listing@example.com", PasswordHash: "h", Nickname: "Listing", Role: 1, Status: 1, TechTags: "[]"}
	if err := e.userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}
	topic := &model.Topic{Name: "列表", Status: 1}
	if err := e.topicRepo.Create(ctx, topic); err != nil {
		t.Fatalf("create topic: %v", err)
	}

	makeOld := func(id uint) {
		old := time.Now().Add(-48 * time.Hour)
		if err := e.db.Model(&model.Article{}).Where("id = ?", id).Update("published_at", old).Error; err != nil {
			t.Fatalf("set old published_at: %v", err)
		}
	}
	createArticle := func(title string, like, view int, withTopic bool) *dto.ArticleListItemDTO {
		req := &dto.CreateArticleRequest{Title: title, Content: "## 正文\n这段内容足够长，用来测试列表行为。", Status: constants.ArticleStatusPublished.Int()}
		if withTopic {
			req.TopicIDs = []uint{topic.ID}
		}
		item, err := e.articleSvc.Create(ctx, user.ID, req)
		if err != nil {
			t.Fatalf("create article %s: %v", title, err)
		}
		if like != 0 || view != 0 {
			if err := e.articleRepo.UpdateCounters(ctx, item.ID, like, view, 0, 0); err != nil {
				t.Fatalf("update counters: %v", err)
			}
		}
		makeOld(item.ID)
		return item
	}

	t.Run("hottest_order", func(t *testing.T) {
		hot := createArticle("热门文章", 20, 0, false)
		createArticle("普通文章", 0, 0, false)
		got, err := e.articleSvc.List(ctx, &dto.ArticleQuery{PageQuery: dto.PageQuery{Page: 1, PageSize: 10}, Sort: constants.SortHottest.String()})
		if err != nil {
			t.Fatalf("list hottest: %v", err)
		}
		items := got.Items.([]*dto.ArticleListItemDTO)
		if len(items) < 2 || items[0].ID != hot.ID {
			t.Fatalf("hottest first item = %v, want hot article %d", itemIDs(items), hot.ID)
		}
	})

	t.Run("pagination_offset", func(t *testing.T) {
		var newest uint
		for i := 0; i < 4; i++ {
			newest = createArticle("分页文章"+string(rune('A'+i)), 0, 0, false).ID
		}
		got, err := e.articleSvc.List(ctx, &dto.ArticleQuery{PageQuery: dto.PageQuery{Page: 1, PageSize: 2}, Sort: constants.SortLatest.String()})
		if err != nil {
			t.Fatalf("list page: %v", err)
		}
		items := got.Items.([]*dto.ArticleListItemDTO)
		if len(items) != 2 || items[0].ID != newest {
			t.Fatalf("page1 first item = %v, want newest %d", itemIDs(items), newest)
		}
	})

	t.Run("hot_score_weight", func(t *testing.T) {
		created := createArticle("计分文章", 1, 2, false)
		got, err := e.articleSvc.List(ctx, &dto.ArticleQuery{PageQuery: dto.PageQuery{Page: 1, PageSize: 20}, Sort: constants.SortLatest.String()})
		if err != nil {
			t.Fatalf("list latest: %v", err)
		}
		for _, item := range got.Items.([]*dto.ArticleListItemDTO) {
			if item.ID == created.ID {
				if item.HotScore != 12 {
					t.Fatalf("hot score = %d, want 12", item.HotScore)
				}
				return
			}
		}
		t.Fatal("scoring article not found in list")
	})

	t.Run("topic_detail_lists", func(t *testing.T) {
		createArticle("话题热门一", 5, 0, true)
		createArticle("话题热门二", 1, 0, true)
		detail, err := e.topicSvc.Detail(ctx, topic.ID)
		if err != nil {
			t.Fatalf("topic detail: %v", err)
		}
		if len(detail.HotArticles) != 2 || len(detail.NewArticles) != 2 {
			t.Fatalf("topic detail lists = hot:%d new:%d, want 2 and 2", len(detail.HotArticles), len(detail.NewArticles))
		}
	})
}

func itemIDs(items []*dto.ArticleListItemDTO) []uint {
	ids := make([]uint, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.ID)
	}
	return ids
}
