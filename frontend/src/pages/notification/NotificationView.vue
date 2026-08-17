<template>
  <div class="notifications">
    <div class="notif-head">
      <h1 class="page-title">通知中心</h1>
      <el-button type="primary" plain @click="onReadAll">全部已读</el-button>
    </div>
    <el-radio-group v-model="unreadOnly" class="notif-filter" @change="load(1)">
      <el-radio-button :value="false">全部</el-radio-button>
      <el-radio-button :value="true">未读</el-radio-button>
    </el-radio-group>

    <div class="notif-list">
      <div v-for="n in notifications" :key="n.id" class="notif-item" :class="{ unread: !n.is_read }" @click="onOpen(n)">
        <div class="notif-icon">{{ typeIcon(n.type) }}</div>
        <div class="notif-body">
          <p class="notif-content">{{ n.actor?.nickname || '系统' }} {{ typeText(n.type) }}</p>
          <p class="notif-detail">{{ n.content }}</p>
        </div>
        <span class="notif-time">{{ formatRelativeTime(n.created_at) }}</span>
      </div>
      <EmptyState v-if="!notifications.length" description="暂无通知" />
    </div>
    <PaginationBar v-model="pagination" @change="load" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import EmptyState from '@/components/EmptyState.vue'
import PaginationBar from '@/components/PaginationBar.vue'
import { useNotificationStore } from '@/stores/notification'
import { markRead } from '@/api/notification'
import { NOTIFICATION_TYPE_TEXT } from '@/constants'
import { usePagination } from '@/hooks/usePagination'
import { formatRelativeTime } from '@/utils/format'

const router = useRouter()
const store = useNotificationStore()
const unreadOnly = ref(false)
const { pagination, setTotal } = usePagination(10)
const notifications = ref(store.notifications)

async function load(page = 1) {
  const res = await store.fetchList(page, pagination.value.pageSize, unreadOnly.value)
  notifications.value = res.items
  setTotal(res.total)
}

function typeIcon(type: string): string {
  if (type === 'follow') return '👤'
  if (type === 'like') return '👍'
  if (type === 'comment') return '💬'
  return '🔔'
}

function typeText(type: string): string {
  return NOTIFICATION_TYPE_TEXT[type] || '有新动态'
}

async function onOpen(n: (typeof store.notifications)[number]) {
  if (!n.is_read) {
    await markRead(n.id)
    n.is_read = true
    store.refreshUnread().catch(() => undefined)
  }
  if (n.target_id && (n.type === 'like' || n.type === 'comment')) {
    router.push(`/articles/${n.target_id}`)
  } else if (n.actor_id) {
    router.push(`/users/${n.actor_id}`)
  }
}

async function onReadAll() {
  await store.readAll()
  notifications.value.forEach((n) => (n.is_read = true))
  load()
}

onMounted(() => load(1))
</script>

<style scoped>
.notif-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.page-title {
  color: #303133;
}
.notif-filter {
  margin-bottom: 16px;
}
.notif-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  border: 1px solid var(--el-border-color-lighter);
}
.notif-item {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  padding: 12px;
  border-radius: 6px;
  cursor: pointer;
}
.notif-item:hover {
  background: #f5f7fa;
}
.notif-item.unread {
  background: #ecf5ff;
}
.notif-icon {
  font-size: 20px;
}
.notif-body {
  flex: 1;
}
.notif-content {
  margin: 0;
  font-size: 14px;
  color: #303133;
}
.notif-detail {
  margin: 4px 0 0;
  font-size: 13px;
  color: #909399;
}
.notif-time {
  font-size: 12px;
  color: #c0c4cc;
  white-space: nowrap;
}
</style>
