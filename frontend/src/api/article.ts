import { del, get, post, put } from '@/utils/request'
import type { ArticleDetail, ArticleItem, PageResult } from '@/types'

export interface ArticleQuery {
  page?: number
  page_size?: number
  sort?: string
  topic_id?: number
  keyword?: string
  author_id?: number
  status?: number
}

export interface CreateArticlePayload {
  title: string
  content: string
  summary?: string
  cover_url?: string
  status?: number
  topic_ids?: number[]
}

export function listArticles(params: ArticleQuery) {
  return get<PageResult<ArticleItem>>('/articles', params)
}

export function getArticle(id: number) {
  return get<ArticleDetail>(`/articles/${id}`)
}

export function createArticle(payload: CreateArticlePayload) {
  return post<ArticleItem>('/articles', payload)
}

export function updateArticle(id: number, payload: CreateArticlePayload) {
  return put<ArticleItem>(`/articles/${id}`, payload)
}

export function publishArticle(id: number) {
  return post<ArticleItem>(`/articles/${id}/publish`)
}

export function offlineArticle(id: number) {
  return post<{ id: number }>(`/articles/${id}/offline`)
}

export function likeArticle(id: number) {
  return post<{ id: number }>(`/articles/${id}/like`)
}

export function unlikeArticle(id: number) {
  return del<{ id: number }>(`/articles/${id}/like`)
}

export function getFeed(page = 1, pageSize = 10) {
  return get<PageResult<ArticleItem>>('/feed', { page, page_size: pageSize })
}
