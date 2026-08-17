# 📝 TechBlog 技术写作博客内容社区

面向技术写作的博客内容社区：支持用户注册发布文章、关注作者、点赞收藏、话题讨论、站内通知与后台管理，打造技术内容分享与交流平台。

## 一、快速启动（Docker Compose，首选）

```bash
# 在项目根目录（任意目录名，含中文目录名均可）
docker compose up -d --build
docker compose ps   # 等待全部容器 healthy
```

启动后访问：

| 入口 | 地址 |
| --- | --- |
| 前端页面 | http://localhost:8102 |
| 后端 API | http://localhost:3102/api/v1 |
| 健康检查 | http://localhost:3102/healthz |
| 接口文档 | 见下方「API 清单」（OpenAPI 摘要见 `backend/api/openapi.yaml`） |

默认管理员：`admin@techblog.com` / `admin123456`（可在 `.env` 中修改）。

停止并清理（含数据卷）：

```bash
docker compose down -v --remove-orphans
```

## 二、项目主要功能

1. **用户注册与登录**：邮箱注册、邮箱验证码登录、GitHub 第三方登录（演示模式）；头像、昵称、个人简介、技术标签设置。
2. **文章发布**：Markdown 编辑器、代码高亮（highlight.js）、封面图、图片上传、话题标签关联；草稿 → 发布 → 下架状态机。
3. **互动功能**：点赞文章、收藏到多收藏夹（公开/私密）。
4. **关注系统**：关注/取关作者、首页关注动态流、粉丝数/关注数统计。
5. **话题广场**：创建/管理话题，话题页展示热门与最新文章。
6. **文章推荐**：首页「最新 / 最热」排序，最热 = 点赞×10 + 阅读量 + 24 小时新文章流量加权 50。
7. **通知系统**：关注/点赞/评论触发站内通知，全部/未读筛选、一键全部已读。
8. **后台管理**：用户禁用/启用、文章下架/恢复、话题管理、平台数据统计（日活、文章数、注册用户趋势）、审计日志。

## 三、技术栈

| 层 | 技术栈 |
| --- | --- |
| 前端 | Vue 3 + TypeScript + Element Plus + Vite |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0 |
| 缓存 | Redis 7 |
| 认证 | JWT + RBAC |
| 日志 | log/slog |
| 参数校验 | go-playground/validator/v10 |

## 四、项目目录结构

```
.
├── backend/                      # Go 后端
│   ├── cmd/server/main.go        # 入口：装配依赖、启动 HTTP
│   ├── internal/
│   │   ├── config/               # 环境变量配置
│   │   ├── database/             # MySQL 连接、AutoMigrate、管理员种子
│   │   ├── model/                # 实体模型（每个实体一个文件）
│   │   ├── dto/                  # 请求/响应 DTO（每个实体一个文件）
│   │   ├── repository/           # 数据访问层（每个实体一个文件）
│   │   ├── service/              # 业务逻辑层（每个实体一个文件）
│   │   ├── handler/              # HTTP 处理器（每个实体一个文件）
│   │   ├── router/               # 路由注册（每个模块一个文件）
│   │   ├── middleware/           # 认证/RBAC/请求ID/审计/限流/日志/错误处理
│   │   ├── constants/            # 枚举、错误码、文案、日志模板
│   │   └── util/                 # JWT、响应、错误、格式化、日志
│   ├── pkg/pagination/           # 分页工具
│   ├── migrations/               # 迁移说明
│   ├── api/openapi.yaml          # OpenAPI 摘要
│   └── Dockerfile                # Go 多阶段构建
├── frontend/                     # Vue 3 前端
│   ├── src/
│   │   ├── api/                  # 每个实体一个 API 文件
│   │   ├── components/           # 共享组件（StatusBadge/ArticleCard/EmptyState/PaginationBar/ConfirmDialog/UserAvatar）
│   │   ├── stores/               # Pinia 状态（按实体拆分）
│   │   ├── pages/                # 每个模块一个页面目录
│   │   ├── router/               # 路由与守卫
│   │   ├── hooks/                # useAuth / usePagination
│   │   ├── utils/                # request / format / markdown
│   │   └── constants/            # 与后端对应枚举
│   ├── Dockerfile                # 前端多阶段构建（Nginx 托管）
│   └── nginx.conf                # 路由与 API 反向代理
├── database/init.sql             # 数据库初始化脚本
├── docker-compose.yml            # 一键编排
├── .env / .env.example           # 环境变量
└── README.md
```

## 五、环境变量说明

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `COMPOSE_PROJECT_NAME` | `techblog` | Compose 项目名，同时作为容器名前缀 |
| `FRONTEND_PORT` | `8102` | 前端宿主机端口 |
| `BACKEND_PORT` | `3102` | 后端宿主机端口 |
| `DB_PORT` | `3502` | MySQL 宿主机端口 |
| `REDIS_PORT` | `6502` | Redis 宿主机端口 |
| `DB_NAME` / `DB_USER` / `DB_PASSWORD` | `techblog_db` / `techblog_user` / `techblog_pwd` | MySQL 连接配置 |
| `DB_ROOT_PASSWORD` | `root_techblog_pwd` | MySQL root 密码 |
| `JWT_SECRET` | `change_me_to_a_long_random_string` | JWT 签名密钥（生产务必修改） |
| `JWT_TTL_HOURS` | `72` | Token 有效期（小时） |
| `REDIS_ADDR` / `REDIS_PASSWORD` | `redis:6379` / 空 | Redis 连接配置 |
| `DEBUG_CODE_MODE` | `true` | 演示模式：验证码接口返回 `debug_code` |
| `ADMIN_EMAIL` / `ADMIN_PASSWORD` | `admin@techblog.com` / `admin123456` | 初始管理员 |

## 六、本地开发（备选）

后端（需要本地 MySQL/Redis，端口见上）：

```bash
cd backend
go mod tidy
go run ./cmd/server
go build ./...
go vet ./...
go test ./...
```

前端：

```bash
cd frontend
npm install
npm run dev        # http://localhost:5173，代理 /api 到 http://localhost:3102
npm run build
```

## 七、API 调用示例（curl）

```bash
BASE=http://localhost:3102/api/v1

# 1. 健康检查
curl -sS http://localhost:3102/healthz

# 2. 发送邮箱验证码（演示模式返回 debug_code）
curl -sS -X POST $BASE/auth/code -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com"}'

# 3. 注册
curl -sS -X POST $BASE/auth/register -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","code":"<debug_code>","password":"alice123","nickname":"Alice","tech_tags":["Go","后端"]}'

# 4. 登录，获取 JWT
TOKEN=$(curl -sS -X POST $BASE/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"alice123"}' | python3 -c 'import sys,json;print(json.load(sys.stdin)["data"]["token"])')

# 5. 携带 JWT 的请求示例
curl -sS $BASE/auth/me -H "Authorization: Bearer $TOKEN"
curl -sS -X POST $BASE/articles -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"title":"Go 并发模型","content":"## 简介\nGo 的 goroutine 非常轻量。","topic_ids":[1],"status":1}'
curl -sS "$BASE/articles?sort=hottest&page=1&page_size=10"
curl -sS -X POST $BASE/articles/1/like -H "Authorization: Bearer $TOKEN"
```

统一响应格式：`{ "code": 0, "message": "ok", "data": ... }`；错误时 `code` 为非 0 业务错误码，HTTP 状态码与之对应。

## 八、API 清单（前缀 `/api/v1`）

### 认证 auth
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/auth/code` | 发送邮箱验证码 |
| POST | `/auth/register` | 邮箱注册（验证码+密码） |
| POST | `/auth/login` | 密码登录 |
| POST | `/auth/login/code` | 验证码登录 |
| POST | `/auth/github` | GitHub 第三方登录（演示） |
| GET | `/auth/me` | 当前用户信息（复用 `UserService.GetProfile`） |

### 用户 user / 关注 follow
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/users/:id` | 用户主页（复用 `UserService.GetProfile`） |
| PUT | `/user/profile` | 更新个人资料 |
| GET | `/users/:id/followers` | 粉丝列表 |
| GET | `/users/:id/following` | 关注列表 |
| POST | `/users/:id/follow` | 关注作者 |
| DELETE | `/users/:id/follow` | 取消关注 |
| GET | `/users/:id/collections` | 用户公开收藏夹 |

### 文章 article
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/articles` | 文章列表（`sort=latest/hottest`，复用 `ArticleService.List`） |
| GET | `/feed` | 关注动态流（复用 `ArticleService.Feed` → `ArticleRepository.ListByAuthorIDs`） |
| POST | `/articles` | 创建文章 |
| GET | `/articles/:id` | 文章详情（阅读量+1，含点赞/收藏状态） |
| PUT | `/articles/:id` | 更新文章 |
| POST | `/articles/:id/publish` | 发布文章（状态机：草稿/下架 → 已发布） |
| POST | `/articles/:id/offline` | 下架文章（状态机：已发布/草稿 → 下架） |
| POST | `/articles/:id/like` | 点赞文章 |
| DELETE | `/articles/:id/like` | 取消点赞 |

### 话题 topic
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/topics` | 话题广场（`all=1` 返回全部启用话题） |
| GET | `/topics/:id` | 话题详情（热门+最新文章，复用 `ArticleService.List` 两次） |
| POST | `/topics` | 创建话题（管理员） |
| PUT | `/topics/:id` | 更新话题（管理员） |
| DELETE | `/topics/:id` | 删除话题（管理员） |
| PUT | `/topics/:id/status` | 启用/禁用话题（管理员） |

### 收藏夹 collection
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/collections` | 我的收藏夹 |
| POST | `/collections` | 创建收藏夹 |
| GET | `/collections/:id` | 收藏夹详情（含文章） |
| PUT | `/collections/:id` | 更新收藏夹 |
| DELETE | `/collections/:id` | 删除收藏夹 |
| POST | `/collections/:id/articles` | 收藏文章 |
| DELETE | `/collections/:id/articles/:articleId` | 取消收藏 |

### 评论 comment
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/articles/:id/comments` | 评论列表 |
| POST | `/articles/:id/comments` | 发表评论（触发通知） |
| DELETE | `/comments/:id` | 删除评论 |

### 通知 notification
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/notifications` | 通知列表（`unread_only=true` 仅未读） |
| PUT | `/notifications/:id` | 单条已读 |
| GET | `/notification/unread-count` | 未读数量 |
| PUT | `/notification/read-all` | 一键全部已读 |

### 上传 upload
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/upload/image` | 上传图片（multipart，字段 `file`） |

### 后台管理 admin（仅管理员）
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/admin/stats` | 平台数据统计 |
| GET | `/admin/users` | 用户列表 |
| PUT | `/admin/users/:id/status` | 禁用/启用用户 |
| GET | `/admin/articles` | 文章列表（复用 `ArticleService.List`） |
| PUT | `/admin/articles/:id/status` | 下架/恢复文章（复用 `ArticleService.AdminUpdateStatus`） |
| GET | `/admin/topics` | 话题列表（复用 `TopicService.List`） |
| PUT | `/admin/topics/:id/status` | 启用/禁用话题（复用 `TopicService.SetStatus`） |
| GET | `/admin/audit-logs` | 审计日志 |

## 九、枚举出现位置清单

### 1. 文章状态 `ArticleStatus`（0 草稿 / 1 已发布 / 2 已下架）

| 层 | 文件 |
| --- | --- |
| 后端常量 | `backend/internal/constants/enums.go` |
| 后端模型 | `backend/internal/model/article.go`（`Status int`） |
| 后端 DTO | `backend/internal/dto/article.go`（Create/Update 的 `oneof=0 1`、ArticleQuery `oneof=0 1 2`） |
| 后端状态机 | `backend/internal/service/article_service.go`（Publish/Offline/AdminUpdateStatus 流转校验） |
| 后端仓储 | `backend/internal/repository/article_repository.go`（UpdateStatus/List 过滤） |
| 后端日志模板 | `backend/internal/constants/log_templates.go`（LogArticlePublished/Offlined/StatusChanged） |
| 后端格式化 | `backend/internal/util/formatters.go`（ArticleStatusText） |
| 后端错误码 | `backend/internal/constants/error_codes.go`（CodeArticleNotFound/Forbidden） |
| 前端常量 | `frontend/src/constants/index.ts`（ARTICLE_STATUS / ARTICLE_STATUS_TEXT） |
| 前端徽标 | `frontend/src/components/StatusBadge.vue`（配合使用） |
| 前端页面 | `frontend/src/pages/article/ArticleDetailView.vue`、`frontend/src/pages/admin/AdminArticles.vue`、`frontend/src/pages/article/ArticleEditorView.vue` |

### 2. 收藏夹可见性 `Visibility`（0 私密 / 1 公开）

| 层 | 文件 |
| --- | --- |
| 后端常量 | `backend/internal/constants/enums.go` |
| 后端模型 | `backend/internal/model/collection.go`（`Visibility int`） |
| 后端 DTO | `backend/internal/dto/collection.go`（Create/Update `oneof=0 1`） |
| 后端服务 | `backend/internal/service/collection_service.go`（私密/公开访问校验） |
| 后端仓储 | `backend/internal/repository/collection_repository.go`（ListPublicByUser） |
| 后端日志模板 | `backend/internal/constants/log_templates.go`（LogCollectionCreated 等） |
| 后端格式化 | `backend/internal/util/formatters.go`（VisibilityText） |
| 后端错误码 | `backend/internal/constants/error_codes.go`（CodeCollectionForbidden） |
| 前端常量 | `frontend/src/constants/index.ts`（VISIBILITY / VISIBILITY_TEXT） |
| 前端页面 | `frontend/src/pages/collection/CollectionView.vue`、`frontend/src/pages/collection/CollectionDetailView.vue`、`frontend/src/pages/user/ProfileView.vue` |

### 3. 角色类型 `RoleType`（1 普通用户 / 2 管理员）

| 层 | 文件 |
| --- | --- |
| 后端常量 | `backend/internal/constants/enums.go` |
| 后端模型 | `backend/internal/model/user.go`（`Role int`） |
| 后端 RBAC 中间件 | `backend/internal/middleware/rbac.go` |
| 后端 JWT | `backend/internal/util/jwt.go`（Claims.Role） |
| 后端日志模板 | `backend/internal/constants/log_templates.go`（LogUserStatusChanged 等） |
| 后端格式化 | `backend/internal/util/formatters.go`（RoleText） |
| 后端错误码 | `backend/internal/constants/error_codes.go`（CodeForbidden） |
| 前端常量 | `frontend/src/constants/index.ts`（ROLE / ROLE_TEXT） |
| 前端路由守卫 | `frontend/src/router/index.ts`（requiresAdmin） |
| 前端布局 | `frontend/src/layouts/MainLayout.vue`（后台入口显隐）、`frontend/src/hooks/useAuth.ts`（isAdmin） |

### 4. 用户状态 `UserStatus`（0 禁用 / 1 启用，话题状态复用）

| 层 | 文件 |
| --- | --- |
| 后端常量 | `backend/internal/constants/enums.go` |
| 后端模型 | `backend/internal/model/user.go`、`backend/internal/model/topic.go` |
| 后端 DTO | `backend/internal/dto/user.go`（UserQuery.Status）、`backend/internal/dto/topic.go` |
| 后端服务 | `backend/internal/service/admin_service.go`（UpdateUserStatus）、`backend/internal/service/topic_service.go`（SetStatus） |
| 后端日志模板 | `backend/internal/constants/log_templates.go` |
| 后端格式化 | `backend/internal/util/formatters.go`（UserStatusText） |
| 后端错误码 | `backend/internal/constants/error_codes.go`（CodeUserDisabled） |
| 前端常量 | `frontend/src/constants/index.ts`（USER_STATUS / USER_STATUS_TEXT） |
| 前端页面 | `frontend/src/pages/admin/AdminUsers.vue`、`frontend/src/pages/admin/AdminTopics.vue` |

### 5. 通知类型 `NotificationType`（follow / like / comment / system）

| 层 | 文件 |
| --- | --- |
| 后端常量 | `backend/internal/constants/enums.go` |
| 后端模型 | `backend/internal/model/notification.go`（`Type string`） |
| 后端服务 | `backend/internal/service/follow_service.go`、`article_service.go`、`comment_service.go`（触发通知） |
| 后端日志模板 | `backend/internal/constants/log_templates.go`（LogNotificationCreated） |
| 后端格式化 | `backend/internal/util/formatters.go`（NotificationTypeText） |
| 前端常量 | `frontend/src/constants/index.ts`（NOTIFICATION_TYPE / NOTIFICATION_TYPE_TEXT） |
| 前端页面 | `frontend/src/pages/notification/NotificationView.vue` |

## 十、Docker 部署说明

- 端口映射：前端 `${FRONTEND_PORT:-8102}:80`、后端 `${BACKEND_PORT:-3102}:8080`、MySQL `${DB_PORT:-3502}:3306`、Redis `${REDIS_PORT:-6502}:6379`。
- 数据卷：`db_data`（MySQL 数据）、`redis_data`（Redis 数据）、`upload_data`（后端上传图片）。
- 数据库初始化：首次启动执行 `database/init.sql`，业务表由后端 GORM AutoMigrate 自动创建。
- 常见问题：
  - 端口被占用：修改 `.env` 中 `FRONTEND_PORT` / `BACKEND_PORT` / `DB_PORT` / `REDIS_PORT` 后重新 `docker compose up -d`。
  - 后端未 healthy：执行 `docker compose logs backend` 查看日志（多为等待 MySQL 初始化）。
  - 忘记管理员密码：修改 `.env` 中 `ADMIN_EMAIL` / `ADMIN_PASSWORD` 后 `docker compose down -v` 重建（会清空数据）。

## 十一、横切关注点说明

1. **JWT 认证 + RBAC**：用户表 `role` 字段（`backend/internal/model/user.go`）→ `backend/internal/middleware/auth.go`、`backend/internal/middleware/rbac.go`、`backend/internal/util/jwt.go` → 前端 `frontend/src/router/index.ts` 路由守卫、`frontend/src/layouts/MainLayout.vue` 按钮显隐。
2. **操作审计日志**：`audit_logs` 表（`backend/internal/model/audit_log.go`）→ `backend/internal/middleware/audit.go` 自动记录写操作 + service 层埋点（`backend/internal/service/audit_service.go`）→ 前端 `frontend/src/pages/admin/AdminAudit.vue`。
3. **全局错误处理与请求追踪**：`backend/internal/middleware/request_id.go`、`backend/internal/middleware/error_handler.go`、`backend/internal/util/app_error.go`、`backend/internal/constants/error_codes.go` → 前端 `frontend/src/utils/request.ts` 拦截器。

## 十二、测试与文档

- 单元测试：`backend/internal/repository/*_test.go`、`backend/internal/service/*_test.go`（表驱动 + SQLite 内存库），运行 `cd backend && go test ./...`。
- 接口文档：完整 API 清单见本 README，OpenAPI 摘要见 `backend/api/openapi.yaml`。

## License

MIT License
