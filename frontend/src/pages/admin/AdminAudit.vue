<template>
  <div>
    <h2 class="admin-title">审计日志</h2>
    <el-input v-model="action" placeholder="按操作类型过滤（如 LOGIN / ARTICLE_CREATE）" clearable class="search" @change="load(1)" />
    <el-table v-loading="loading" :data="logs" stripe>
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="username" label="操作者" min-width="120" />
      <el-table-column prop="action" label="操作" width="140" />
      <el-table-column prop="module" label="模块" width="100" />
      <el-table-column prop="detail" label="详情" min-width="220" show-overflow-tooltip />
      <el-table-column prop="ip" label="IP" width="120" />
      <el-table-column prop="created_at" label="时间" width="160">
        <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
      </el-table-column>
    </el-table>
    <PaginationBar v-model="pagination" @change="load" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import PaginationBar from '@/components/PaginationBar.vue'
import { listAuditLogs } from '@/api/admin'
import { usePagination } from '@/hooks/usePagination'
import { formatDate } from '@/utils/format'
import type { AuditLog } from '@/types'

const logs = ref<AuditLog[]>([])
const action = ref('')
const loading = ref(false)
const { pagination, setTotal } = usePagination(10)

async function load(page = 1) {
  loading.value = true
  try {
    const res = await listAuditLogs({ page, page_size: pagination.value.pageSize, action: action.value })
    logs.value = res.items
    setTotal(res.total)
  } finally {
    loading.value = false
  }
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
