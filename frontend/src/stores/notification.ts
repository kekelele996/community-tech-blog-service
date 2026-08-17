import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as notificationApi from '@/api/notification'
import type { Notification } from '@/types'

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<Notification[]>([])
  const unreadCount = ref(0)
  const total = ref(0)

  async function fetchList(page = 1, pageSize = 10, unreadOnly = false) {
    const res = await notificationApi.listNotifications({ page, page_size: pageSize, unread_only: unreadOnly })
    notifications.value = res.items
    total.value = res.total
    return res
  }

  async function refreshUnread() {
    const res = await notificationApi.getUnreadCount()
    unreadCount.value = res.unread_count
    return unreadCount.value
  }

  async function readAll() {
    await notificationApi.markAllRead()
    unreadCount.value = 0
  }

  return { notifications, unreadCount, total, fetchList, refreshUnread, readAll }
})
