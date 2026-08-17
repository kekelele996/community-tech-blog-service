<template>
  <div v-loading="loading" class="topic-detail">
    <div class="topic-head">
      <h1 class="topic-title"># {{ detail?.topic.name }}</h1>
      <p class="topic-desc">{{ detail?.topic.description || '暂无描述' }}</p>
      <span class="topic-count">共 {{ detail?.topic.article_count ?? 0 }} 篇文章</span>
    </div>

    <section class="section">
      <h2>🔥 热门文章</h2>
      <div class="article-list">
        <ArticleCard v-for="a in detail?.hot_articles || []" :key="a.id" :article="a" @open="goDetail" />
        <EmptyState v-if="!(detail?.hot_articles || []).length" description="暂无热门文章" />
      </div>
    </section>

    <section class="section">
      <h2>🕐 最新文章</h2>
      <div class="article-list">
        <ArticleCard v-for="a in detail?.new_articles || []" :key="a.id" :article="a" @open="goDetail" />
        <EmptyState v-if="!(detail?.new_articles || []).length" description="暂无最新文章" />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ArticleCard from '@/components/ArticleCard.vue'
import EmptyState from '@/components/EmptyState.vue'
import { getTopic } from '@/api/topic'
import type { TopicDetail } from '@/api/topic'

const route = useRoute()
const router = useRouter()
const detail = ref<TopicDetail | null>(null)
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    detail.value = await getTopic(Number(route.params.id))
  } finally {
    loading.value = false
  }
}

function goDetail(id: number) {
  router.push(`/articles/${id}`)
}

onMounted(load)
</script>

<style scoped>
.topic-head {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  border: 1px solid var(--el-border-color-lighter);
  margin-bottom: 16px;
}
.topic-title {
  color: var(--el-color-primary);
  margin: 0 0 8px;
}
.topic-desc {
  color: #606266;
  margin: 0 0 8px;
}
.topic-count {
  color: #909399;
  font-size: 13px;
}
.section {
  background: #fff;
  border-radius: 8px;
  padding: 20px 24px;
  border: 1px solid var(--el-border-color-lighter);
  margin-bottom: 16px;
}
.section h2 {
  font-size: 17px;
  color: #303133;
  margin: 0 0 16px;
}
.article-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
