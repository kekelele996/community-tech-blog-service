<template>
  <div>
    <h2 class="admin-title">用户管理</h2>
    <el-input v-model="keyword" placeholder="搜索邮箱/昵称" clearable class="search" @change="load(1)" />
    <el-table v-loading="loading" :data="users" stripe>
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="nickname" label="昵称" min-width="120" />
      <el-table-column prop="email" label="邮箱" min-width="180" />
      <el-table-column label="角色" width="100">
        <template #default="{ row }"><el-tag size="small" :type="row.role === 2 ? 'danger' : 'info'">{{ row.role === 2 ? '管理员' : '普通用户' }}</el-tag></template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }"><StatusBadge :value="row.status" :text-map="USER_STATUS_TEXT" /></template>
      </el-table-column>
      <el-table-column prop="article_count" label="文章数" width="80" />
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button size="small" :type="row.status === 1 ? 'danger' : 'success'" @click="toggleStatus(row)">
            {{ row.status === 1 ? '禁用' : '启用' }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <PaginationBar v-model="pagination" @change="load" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import PaginationBar from '@/components/PaginationBar.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { listAdminUsers, updateUserStatus } from '@/api/admin'
import { USER_STATUS_TEXT } from '@/constants'
import { usePagination } from '@/hooks/usePagination'
import type { AdminUser } from '@/types'

const users = ref<AdminUser[]>([])
const keyword = ref('')
const loading = ref(false)
const { pagination, setTotal } = usePagination(10)

async function load(page = 1) {
  loading.value = true
  try {
    const res = await listAdminUsers({ page, page_size: pagination.value.pageSize, keyword: keyword.value })
    users.value = res.items
    setTotal(res.total)
  } finally {
    loading.value = false
  }
}

async function toggleStatus(row: AdminUser) {
  const next = row.status === 1 ? 0 : 1
  await updateUserStatus(row.id, next)
  row.status = next
  ElMessage.success('已更新')
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
