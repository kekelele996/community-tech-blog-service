<template>
  <div>
    <div class="admin-head">
      <h2 class="admin-title">话题管理</h2>
      <el-button type="primary" @click="createVisible = true">新建话题</el-button>
    </div>
    <el-table v-loading="loading" :data="topics" stripe>
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="名称" min-width="140" />
      <el-table-column prop="description" label="描述" min-width="220" show-overflow-tooltip />
      <el-table-column prop="article_count" label="文章数" width="90" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }"><StatusBadge :value="row.status" :text-map="USER_STATUS_TEXT" /></template>
      </el-table-column>
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" :type="row.status === 1 ? 'warning' : 'success'" @click="toggleStatus(row)">
            {{ row.status === 1 ? '禁用' : '启用' }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <PaginationBar v-model="pagination" @change="load" />

    <el-dialog v-model="createVisible" title="新建话题" width="420px">
      <el-form :model="createForm" label-width="70px">
        <el-form-item label="名称"><el-input v-model="createForm.name" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="createForm.description" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" @click="onCreate">创建</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="editVisible" title="编辑话题" width="420px">
      <el-form :model="editForm" label-width="70px">
        <el-form-item label="名称"><el-input v-model="editForm.name" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="editForm.description" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="onSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import PaginationBar from '@/components/PaginationBar.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { listAdminTopics, updateTopicStatus } from '@/api/admin'
import { createTopic as createTopicApi, updateTopic as updateTopicApi } from '@/api/topic'
import { USER_STATUS_TEXT } from '@/constants'
import { usePagination } from '@/hooks/usePagination'
import type { Topic } from '@/types'

const topics = ref<Topic[]>([])
const loading = ref(false)
const createVisible = ref(false)
const editVisible = ref(false)
const createForm = reactive({ name: '', description: '' })
const editForm = reactive({ id: 0, name: '', description: '' })
const { pagination, setTotal } = usePagination(10)

async function load(page = 1) {
  loading.value = true
  try {
    const res = await listAdminTopics({ page, page_size: pagination.value.pageSize })
    topics.value = res.items
    setTotal(res.total)
  } finally {
    loading.value = false
  }
}

async function onCreate() {
  if (!createForm.name.trim()) {
    ElMessage.warning('请输入话题名称')
    return
  }
  await createTopicApi({ ...createForm })
  ElMessage.success('创建成功')
  createVisible.value = false
  createForm.name = ''
  createForm.description = ''
  await load(1)
}

function openEdit(row: Topic) {
  editForm.id = row.id
  editForm.name = row.name
  editForm.description = row.description
  editVisible.value = true
}

async function onSave() {
  await updateTopicApi(editForm.id, { name: editForm.name, description: editForm.description })
  ElMessage.success('已保存')
  editVisible.value = false
  await load()
}

async function toggleStatus(row: Topic) {
  const next = row.status === 1 ? 0 : 1
  await updateTopicStatus(row.id, next)
  row.status = next
  ElMessage.success('已更新')
}

onMounted(() => load(1))
</script>

<style scoped>
.admin-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.admin-title {
  color: #303133;
  margin-top: 0;
}
</style>
