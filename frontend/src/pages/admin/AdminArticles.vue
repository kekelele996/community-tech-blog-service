<template>
  <div>
    <h2 class="admin-title">文章管理</h2>
    <el-input v-model="keyword" placeholder="搜索标题/摘要" clearable class="search" @change="load(1)" />
    <el-table v-loading="loading" :data="articles" stripe>
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
      <el-table-column prop="author?.nickname" label="作者" width="120" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }"><StatusBadge :value="row.status" :text-map="ARTICLE_STATUS_TEXT" :type-map="statusTypes" /></template>
      </el-table-column>
      <el-table-column prop="like_count" label="点赞" width="70" />
      <el-table-column prop="view_count" label="阅读" width="70" />
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="go(row.id)">查看</el-button>
          <el-button size="small" :type="row.status === 1 ? 'danger' : 'success'" @click="toggleStatus(row)">
            {{ row.status === 1 ? '下架' : '恢复' }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <PaginationBar v-model="pagination" @change="load" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import PaginationBar from '@/components/PaginationBar.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { listAdminArticles, updateArticleStatus } from '@/api/admin'
import { ARTICLE_STATUS, ARTICLE_STATUS_TEXT } from '@/constants'
import { usePagination } from '@/hooks/usePagination'
import type { ArticleItem } from '@/types'

const router = useRouter()
const articles = ref<ArticleItem[]>([])
const keyword = ref('')
const loading = ref(false)
const { pagination, setTotal } = usePagination(10)

const statusTypes: Record<number, 'success' | 'info' | 'warning' | 'danger'> = {
  [ARTICLE_STATUS.DRAFT]: 'info',
  [ARTICLE_STATUS.PUBLISHED]: 'success',
  [ARTICLE_STATUS.OFFLINE]: 'danger',
}

async function load(page = 1) {
  loading.value = true
  try {
    const res = await listAdminArticles({ page, page_size: pagination.value.pageSize, keyword: keyword.value })
    articles.value = res.items
    setTotal(res.total)
  } finally {
    loading.value = false
  }
}

async function toggleStatus(row: ArticleItem) {
  const next = row.status === 1 ? 2 : 1
  await updateArticleStatus(row.id, next)
  row.status = next
  ElMessage.success('已更新')
}

function go(id: number) {
  router.push(`/articles/${id}`)
}

onMounted(() => load(1))
</script>

<style scoped>
.admin-title {
  color: #303133;
  margin-top: 0;
}
.search {
  max-width: 320px;
  margin-bottom: 16px;
}
</style>
