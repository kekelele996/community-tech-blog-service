import { get, put } from '@/utils/request'
import type { Notification, PageResult } from '@/types'

export function listNotifications(params: { page?: number; page_size?: number; unread_only?: boolean }) {
  return get<PageResult<Notification>>('/notifications', params)
}

export function getUnreadCount() {
  return get<{ unread_count: number }>('/notification/unread-count')
}

export function markAllRead() {
  return put<Record<string, never>>('/notification/read-all')
}

export function markRead(id: number) {
  return put<{ id: number }>(`/notifications/${id}`)
}
