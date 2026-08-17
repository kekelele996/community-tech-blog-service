# gb-51-1 techblog 执行验证报告

- 项目编号：gb-51-1
- 项目名称：技术写作博客内容社区（techblog）
- 项目路径：`/Users/gaobo/repositories/gitlab/评审项目/0-1代码生成提示词/golang-改编提示词/博客主题项目提示词/gb-51-1`
- 验证时间：2026-08-17 02:00 ~ 02:12 (Asia/Shanghai)

## 端口与启动命令

| 服务 | 端口 |
| --- | --- |
| 前端 | 8102:80 |
| 后端 | 3102:8080 |
| MySQL | 3502:3306 |
| Redis | 6502:6379 |

启动命令：

```bash
docker compose config --quiet
docker compose up -d --build
docker compose ps
```

## docker compose ps（healthy 结果）

```
NAME                STATUS                    PORTS
techblog-backend    Up 8 seconds (healthy)    0.0.0.0:3102->8080/tcp
techblog-db         Up 24 seconds (healthy)   0.0.0.0:3502->3306/tcp
techblog-frontend   Up 3 seconds              0.0.0.0:8102->80/tcp
techblog-redis      Up 24 seconds (healthy)   0.0.0.0:6502->6379/tcp
```

后端、数据库、Redis 均 healthy；前端 Nginx 返回 200。验证完成后执行 `docker compose down -v --remove-orphans`，`docker compose ps` 为空。

## curl 接口冒烟清单（23/23 PASS）

| # | 方法 | 路径 | 状态码 |
| --- | --- | --- | --- |
| 1 | GET | `/healthz` | 200 |
| 2 | POST | `/api/v1/auth/code` | 200 |
| 3 | POST | `/api/v1/auth/register` | 200 |
| 4 | POST | `/api/v1/auth/login` | 200 |
| 5 | GET | `/api/v1/auth/me` | 200 |
| 6 | GET | `/api/v1/topics?page=1&page_size=10` | 200 |
| 7 | POST | `/api/v1/articles`（创建并发布） | 200 |
| 8 | GET | `/api/v1/articles?sort=hottest` | 200 |
| 9 | GET | `/api/v1/articles/:id` | 200 |
| 10 | POST | `/api/v1/articles/:id/like` | 200 |
| 11 | POST | `/api/v1/auth/register`（第二个用户） | 200 |
| 12 | POST | `/api/v1/users/:id/follow` | 200 |
| 13 | GET | `/api/v1/notifications` | 200 |
| 14 | GET | `/api/v1/notification/unread-count` | 200 |
| 15 | POST | `/api/v1/articles/:id/comments` | 200 |
| 16 | POST | `/api/v1/collections` | 200 |
| 17 | POST | `/api/v1/collections/:id/articles` | 200 |
| 18 | POST | `/api/v1/auth/login`（管理员） | 200 |
| 19 | GET | `/api/v1/admin/stats` | 200 |
| 20 | GET | `/api/v1/admin/users` | 200 |
| 21 | GET | `/api/v1/admin/audit-logs` | 200 |
| 22 | GET | `/api/v1/users/:id` | 200 |
| 23 | GET | `/api/v1/users/:id/followers` | 200 |

## 浏览器验证页面与截图（内置 Playwright 包装脚本）

使用 `$HOME/.codex/skills/playwright/scripts/playwright_cli.sh`（独立 `--session techblog`），未使用外部 Chrome。

| 页面 | 关键交互 | 截图 |
| --- | --- | --- |
| 首页信息流 | 最新/最热排序、话题侧栏、文章卡片 | `output/01-home.png`、`output/11-home-final.png` |
| 文章详情 | Markdown 渲染 + 代码高亮、点赞、收藏弹窗、评论 | `output/02-article-detail.png`、`output/08-article-detail-loggedin.png`、`output/09-article-collect.png` |
| 登录页 | 密码登录（管理员账号） | `output/03-login-home.png` |
| 后台数据统计 | 日活/文章/用户/评论 + 趋势图 | `output/04-admin-dashboard.png` |
| 通知中心 | 全部/未读筛选、全部已读 | `output/05-notifications.png` |
| 收藏夹 | 新建收藏夹（UI 创建成功） | `output/06-collections.png` |
| 话题广场 | 话题列表 | `output/07-topics.png` |
| 写文章 | 保存并发布 → 跳转文章详情（UI 全流程） | 见文章详情截图 |
| 个人主页 | 资料、文章数/粉丝/关注、文章列表 | `output/10-user-profile.png` |

浏览器验证结论：登录、创建收藏夹、创建并发布文章、点赞、收藏、评论展示等主要交互均正常，接口无阻断性异常。

## 修复过程摘要

1. Gin 路由冲突：`/users/me`、`/articles/feed`、`/topics/all`、`/notifications/read-all` 与参数路由段冲突，改为独立路径（`/user/profile`、`/feed`、`topics?all=1`、`/notification/read-all`）。
2. 服务层 `GetRequestID` 兼容 `context.Context`（service 传入 `c.Request.Context()`）。
3. GORM `default` 标签导致零值字段（Topic.Status=0）被跳过，测试与状态设置调整。
4. 最热排序 SQL 兼容 SQLite 测试库（dialect 判断 `NOW() - INTERVAL 24 HOUR` vs `datetime('now','-24 hours')`）。
5. `marked` v12 移除 `highlight` 选项，改为后处理正则 + highlight.js 高亮。
6. `useAuth` 由 localStorage 计算改为共享响应式状态，修复登录后导航栏不刷新的问题。
7. 编辑器「保存并发布」对已发布文章重复调用 publish 触发状态机 400，改为仅在 `status !== 1` 时调用。
8. 文章详情 `collections` 字段 `omitempty` 导致前端 `.find` 崩溃，前端兜底 `|| []`。

## 质量检查

- `cd backend && go mod tidy && go build ./... && go vet ./... && go test ./...` 全部通过（repository + service 表驱动单测）。
- `cd frontend && npm run build`（含 `vue-tsc --noEmit`）通过。
- `docker compose config --quiet` 在中文目录名下通过。
- 核心实体（User / Article / Topic / Collection / Follow / Notification / Comment）贯穿 model → dto → repository → service → handler → router → 前端 api/store/page。
- 横切关注点：JWT+RBAC、审计日志、全局错误处理+请求追踪均已实现并写明触达文件层。
- 共享枚举：ArticleStatus / Visibility / RoleType / UserStatus / NotificationType，README 列出全部出现位置。
- 共享前端组件 6 个（StatusBadge/ArticleCard/EmptyState/PaginationBar/ConfirmDialog/UserAvatar），hooks/utils 4 个。
- 后端中间件 7 个（RequestID/RequestLogger/ErrorHandler/Audit/RateLimit/Auth/RequireRole）。

## git commit hash

见最终提交：`git log -1 --format=%H`
