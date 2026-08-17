import { get, put } from '@/utils/request'
import type { AdminStats, AdminUser, ArticleItem, AuditLog, PageResult, Topic } from '@/types'

export function getStats() {
  return get<AdminStats>('/admin/stats')
}

export function listAdminUsers(params: { page?: number; page_size?: number; keyword?: string; status?: number }) {
  return get<PageResult<AdminUser>>('/admin/users', params)
}

export function updateUserStatus(id: number, status: number) {
  return put<{ id: number; status: number }>(`/admin/users/${id}/status`, { status })
}

export function listAdminArticles(params: { page?: number; page_size?: number; keyword?: string; status?: number }) {
  return get<PageResult<ArticleItem>>('/admin/articles', params)
}

export function updateArticleStatus(id: number, status: number) {
  return put<{ id: number; status: number }>(`/admin/articles/${id}/status`, { status })
}

export function listAdminTopics(params: { page?: number; page_size?: number; keyword?: string }) {
  return get<PageResult<Topic>>('/admin/topics', params)
}

export function updateTopicStatus(id: number, status: number) {
  return put<{ id: number; status: number }>(`/admin/topics/${id}/status`, { status })
}

export function listAuditLogs(params: { page?: number; page_size?: number; action?: string }) {
  return get<PageResult<AuditLog>>('/admin/audit-logs', params)
}
