# Bug 复现说明

## Bug 是什么

文章列表的最热排序失效，分页会跳过首页文章，话题页的热门/最新文章列表各少一项，文章热度分也算错。

## 如何触发

```bash
cd backend
make test-list
```

## 错误信息

```
--- FAIL: TestArticleListingComposite/hottest_order
    hottest first item = [], want hot article 1
--- FAIL: TestArticleListingComposite/pagination_offset
    page1 first item = [4 3], want newest 6
--- FAIL: TestArticleListingComposite/hot_score_weight
    scoring article not found in list
--- FAIL: TestArticleListingComposite/topic_detail_lists
    topic detail lists = hot:0 new:0, want 2 and 2
```
