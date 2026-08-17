// 全局类型定义（与后端 dto 对应）
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PageResult<T = unknown> {
  items: T[]
  total: number
  page: number
  page_size: number
}

export interface UserProfile {
  id: number
  email: string
  nickname: string
  avatar: string
  bio: string
  tech_tags: string[]
  role: number
  status: number
  github_id: string
  follower_count: number
  following_count: number
  article_count: number
  is_following: boolean
  last_login_at?: string
  created_at: string
}

export interface AuthResponse {
  token: string
  user: UserProfile
}

export interface SendCodeResponse {
  email: string
  expires_in: number
  debug_code?: string
}

export interface TopicBrief {
  id: number
  name: string
}

export interface Topic {
  id: number
  name: string
  description: string
  article_count: number
  status: number
  created_at: string
}

export interface ArticleItem {
  id: number
  author_id: number
  title: string
  summary: string
  cover_url: string
  status: number
  like_count: number
  view_count: number
  favorite_count: number
  comment_count: number
  hot_score: number
  published_at?: string
  created_at: string
  author?: UserProfile
  topics: TopicBrief[]
}

export interface ArticleDetail extends ArticleItem {
  content: string
  is_liked: boolean
  is_collected: boolean
  collections: CollectionBrief[]
}

export interface CollectionBrief {
  id: number
  name: string
  has_article: boolean
}

export interface Collection {
  id: number
  user_id: number
  name: string
  description: string
  visibility: number
  article_count: number
  created_at: string
  updated_at: string
}

export interface CollectionDetail {
  id: number
  user_id: number
  name: string
  description: string
  visibility: number
  article_count: number
  created_at: string
  updated_at: string
  articles: ArticleItem[]
}

export interface Notification {
  id: number
  user_id: number
  actor_id: number
  type: string
  target_id: number
  content: string
  is_read: boolean
  created_at: string
  actor?: UserProfile
}

export interface Comment {
  id: number
  article_id: number
  author_id: number
  parent_id: number
  content: string
  status: number
  created_at: string
  author?: UserProfile
}

export interface AuditLog {
  id: number
  user_id: number
  username: string
  role: number
  action: string
  module: string
  detail: string
  ip: string
  request_id: string
  created_at: string
}

export interface AdminStats {
  dau: number
  total_articles: number
  total_users: number
  total_comments: number
  daily_active: Array<{ date: string; count: number }>
  new_users: Array<{ date: string; count: number }>
}

export interface AdminUser {
  id: number
  email: string
  nickname: string
  avatar: string
  role: number
  status: number
  article_count: number
  follower_count: number
  last_login_at?: string
  created_at: string
}
