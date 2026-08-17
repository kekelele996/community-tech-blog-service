<template>
  <div class="topics">
    <h1 class="page-title">话题广场</h1>
    <el-input v-model="keyword" placeholder="搜索话题" clearable class="topic-search" @change="loadTopics(1)" />
    <div v-loading="loading" class="topic-grid">
      <div v-for="t in topics" :key="t.id" class="topic-card" @click="router.push(`/topics/${t.id}`)">
        <h3 class="topic-name"># {{ t.name }}</h3>
        <p class="topic-desc">{{ t.description || '暂无描述' }}</p>
        <span class="topic-count">{{ t.article_count }} 篇文章</span>
      </div>
      <EmptyState v-if="!topics.length" description="暂无话题" />
    </div>
    <PaginationBar v-model="pagination" @change="loadTopics" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import EmptyState from '@/components/EmptyState.vue'
import PaginationBar from '@/components/PaginationBar.vue'
import { listTopics } from '@/api/topic'
import { usePagination } from '@/hooks/usePagination'
import type { Topic } from '@/types'

const router = useRouter()
const topics = ref<Topic[]>([])
const keyword = ref('')
const loading = ref(false)
const { pagination, setTotal } = usePagination(12)

async function loadTopics(page = 1) {
  loading.value = true
  try {
    const res = await listTopics({ page, page_size: pagination.value.pageSize, keyword: keyword.value })
    topics.value = Array.isArray(res) ? res : res.items
    if (!Array.isArray(res)) setTotal(res.total)
  } finally {
    loading.value = false
  }
}

onMounted(() => loadTopics(1))
</script>

<style scoped>
.page-title {
  color: #303133;
}
.topic-search {
  max-width: 360px;
  margin-bottom: 16px;
}
.topic-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  min-height: 200px;
}
.topic-card {
  background: #fff;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  padding: 20px;
  cursor: pointer;
  transition: box-shadow 0.2s;
}
.topic-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}
.topic-name {
  color: var(--el-color-primary);
  margin: 0 0 8px;
}
.topic-desc {
  color: #606266;
  font-size: 13px;
  margin: 0 0 12px;
  min-height: 38px;
}
.topic-count {
  color: #909399;
  font-size: 12px;
}
</style>
