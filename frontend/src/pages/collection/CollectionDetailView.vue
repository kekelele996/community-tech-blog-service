<template>
  <div v-loading="loading" class="collection-detail">
    <div class="detail-head">
      <div>
        <h1 class="page-title">📁 {{ detail?.name }}</h1>
        <p class="detail-desc">{{ detail?.description || '暂无描述' }}</p>
        <div class="detail-meta">
          <el-tag size="small" :type="detail?.visibility === 1 ? 'success' : 'info'">
            {{ detail?.visibility === 1 ? '公开' : '私密' }}
          </el-tag>
          <span>{{ detail?.article_count ?? 0 }} 篇文章</span>
        </div>
      </div>
      <div v-if="isOwner">
        <el-button @click="editVisible = true">编辑</el-button>
        <el-button type="danger" plain @click="deleteVisible = true">删除</el-button>
      </div>
    </div>

    <div class="article-list">
      <ArticleCard v-for="a in detail?.articles || []" :key="a.id" :article="a" @open="goDetail" />
      <EmptyState v-if="!(detail?.articles || []).length" description="收藏夹为空" />
    </div>

    <el-dialog v-model="editVisible" title="编辑收藏夹" width="440px">
      <el-form :model="editForm" label-width="70px">
        <el-form-item label="名称"><el-input v-model="editForm.name" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="editForm.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="可见性">
          <el-radio-group v-model="editForm.visibility">
            <el-radio :value="0">私密</el-radio>
            <el-radio :value="1">公开</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="onSave">保存</el-button>
      </template>
    </el-dialog>

    <ConfirmDialog v-model="deleteVisible" title="删除收藏夹" message="删除后不可恢复，确定删除该收藏夹吗？" confirm-text="删除" @confirm="onDelete" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import ArticleCard from '@/components/ArticleCard.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import EmptyState from '@/components/EmptyState.vue'
import { getCollection, updateCollection, deleteCollection } from '@/api/collection'
import { useAuth } from '@/hooks/useAuth'
import type { CollectionDetail } from '@/types'

const route = useRoute()
const router = useRouter()
const { user } = useAuth()

const detail = ref<CollectionDetail | null>(null)
const loading = ref(false)
const editVisible = ref(false)
const deleteVisible = ref(false)
const editForm = reactive({ name: '', description: '', visibility: 0 })

const isOwner = computed(() => detail.value?.user_id === user.value?.id)

async function load() {
  loading.value = true
  try {
    detail.value = await getCollection(Number(route.params.id))
  } finally {
    loading.value = false
  }
}

function openEdit() {
  if (!detail.value) return
  editForm.name = detail.value.name
  editForm.description = detail.value.description
  editForm.visibility = detail.value.visibility
  editVisible.value = true
}

async function onSave() {
  if (!detail.value) return
  await updateCollection(detail.value.id, { ...editForm })
  ElMessage.success('已更新')
  editVisible.value = false
  await load()
}

async function onDelete() {
  if (!detail.value) return
  await deleteCollection(detail.value.id)
  ElMessage.success('已删除')
  router.push('/collections')
}

function goDetail(id: number) {
  router.push(`/articles/${id}`)
}

onMounted(load)
</script>

<style scoped>
.detail-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  border: 1px solid var(--el-border-color-lighter);
  margin-bottom: 16px;
}
.page-title {
  color: #303133;
  margin: 0 0 6px;
}
.detail-desc {
  color: #909399;
  margin: 0 0 8px;
}
.detail-meta {
  display: flex;
  gap: 12px;
  align-items: center;
  color: #909399;
  font-size: 13px;
}
.article-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
