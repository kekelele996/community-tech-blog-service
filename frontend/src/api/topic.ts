import { del, get, post, put } from '@/utils/request'
import type { ArticleItem, PageResult, Topic } from '@/types'

export interface TopicDetail {
  topic: Topic
  hot_articles: ArticleItem[]
  new_articles: ArticleItem[]
}

export interface CreateTopicPayload {
  name: string
  description?: string
}

export function listTopics(params: { page?: number; page_size?: number; keyword?: string; all?: number } = {}) {
  return get<PageResult<Topic> | Topic[]>('/topics', params)
}

export function getTopic(id: number) {
  return get<TopicDetail>(`/topics/${id}`)
}

export function createTopic(payload: CreateTopicPayload) {
  return post<Topic>('/topics', payload)
}

export function updateTopic(id: number, payload: CreateTopicPayload) {
  return put<Topic>(`/topics/${id}`, payload)
}

export function deleteTopic(id: number) {
  return del<{ id: number }>(`/topics/${id}`)
}

export function updateTopicStatus(id: number, status: number) {
  return put<{ id: number; status: number }>(`/topics/${id}/status`, { status })
}
