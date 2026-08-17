<template>
  <div class="home">
    <div class="home-main">
      <el-radio-group v-model="sort" size="large" class="sort-tabs" @change="loadArticles(1)">
        <el-radio-button :value="SORT.LATEST">最新</el-radio-button>
        <el-radio-button :value="SORT.HOTTEST">最热</el-radio-button>
        <el-radio-button v-if="isLoggedIn" :value="'feed'">关注动态</el-radio-button>
      </el-radio-group>

      <div v-loading="loading" class="article-list">
        <template v-if="articles.length">
          <ArticleCard v-for="a in articles" :key="a.id" :article="a" @open="goDetail" />
        </template>
        <EmptyState v-else description="暂无文章，快去发布第一篇吧！" :show-action="isLoggedIn" action-text="写文章" @action="router.push('/editor')" />
      </div>

      <PaginationBar v-model="pagination" @change="loadArticles" />
    </div>

    <aside class="home-side">
      <div class="side-card">
        <h3>热门话题</h3>
        <div class="topic-list">
          <div v-for="t in topics" :key="t.id" class="topic-item" @click="router.push(`/topics/${t.id}`)">
            <span class="topic-name"># {{ t.name }}</span>
            <span class="topic-count">{{ t.article_count }} 篇</span>
          </div>
          <EmptyState v-if="!topics.length" description="暂无话题" />
        </div>
      </div>
      <div class="side-card">
        <h3>社区说明</h3>
        <p class="side-text">技术写作博客内容社区：支持 Markdown 写作、代码高亮、话题广场、关注动态、点赞收藏与站内通知。</p>
      </div>
    </aside>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import ArticleCard from '@/components/ArticleCard.vue'
import EmptyState from '@/components/EmptyState.vue'
import PaginationBar from '@/components/PaginationBar.vue'
import { SORT_TYPE as SORT } from '@/constants'
import { listArticles, getFeed } from '@/api/article'
import { listTopics } from '@/api/topic'
import { usePagination } from '@/hooks/usePagination'
import { useAuth } from '@/hooks/useAuth'
import type { ArticleItem, Topic } from '@/types'

const router = useRouter()
const { isLoggedIn } = useAuth()
const articles = ref<ArticleItem[]>([])
const topics = ref<Topic[]>([])
const sort = ref<string>(SORT.LATEST)
const loading = ref(false)
const { pagination, setTotal } = usePagination(10)

async function loadArticles(page = 1) {
  loading.value = true
  try {
    if (sort.value === 'feed') {
      const res = await getFeed(page, pagination.value.pageSize)
      articles.value = res.items
      setTotal(res.total)
    } else {
      const res = await listArticles({ page, page_size: pagination.value.pageSize, sort: sort.value })
      articles.value = res.items
      setTotal(res.total)
    }
  } finally {
    loading.value = false
  }
}

async function loadTopics() {
  const res = await listTopics({ page: 1, page_size: 8 })
  topics.value = Array.isArray(res) ? res : res.items
}

function goDetail(id: number) {
  router.push(`/articles/${id}`)
}

onMounted(() => {
  loadArticles(1)
  loadTopics()
})
</script>

<style scoped>
.home {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 24px;
}
.sort-tabs {
  margin-bottom: 16px;
}
.article-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-height: 200px;
}
.home-side {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.side-card {
  background: #fff;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  padding: 16px;
}
.side-card h3 {
  margin: 0 0 12px;
  font-size: 15px;
  color: #303133;
}
.topic-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.topic-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 10px;
  border-radius: 6px;
  cursor: pointer;
  background: #f5f7fa;
}
.topic-item:hover {
  background: #ecf5ff;
}
.topic-name {
  color: var(--el-color-primary);
  font-size: 13px;
}
.topic-count {
  color: #909399;
  font-size: 12px;
}
.side-text {
  color: #606266;
  font-size: 13px;
  line-height: 1.8;
  margin: 0;
}
</style>
