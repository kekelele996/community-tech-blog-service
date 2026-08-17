<template>
  <div class="collections">
    <div class="collect-head">
      <h1 class="page-title">我的收藏夹</h1>
      <el-button type="primary" @click="createVisible = true">新建收藏夹</el-button>
    </div>
    <div v-loading="loading" class="collect-grid">
      <div v-for="c in collections" :key="c.id" class="collect-card" @click="router.push(`/collections/${c.id}`)">
        <div class="collect-icon">{{ c.visibility === 1 ? '📂' : '🔒' }}</div>
        <h3 class="collect-name">{{ c.name }}</h3>
        <p class="collect-desc">{{ c.description || '暂无描述' }}</p>
        <div class="collect-meta">
          <el-tag size="small" :type="c.visibility === 1 ? 'success' : 'info'">{{ c.visibility === 1 ? '公开' : '私密' }}</el-tag>
          <span>{{ c.article_count }} 篇文章</span>
        </div>
      </div>
      <EmptyState v-if="!collections.length" description="还没有收藏夹" action-text="新建收藏夹" :show-action="true" @action="createVisible = true" />
    </div>

    <el-dialog v-model="createVisible" title="新建收藏夹" width="440px">
      <el-form :model="form" label-width="70px">
        <el-form-item label="名称"><el-input v-model="form.name" placeholder="如：Go 学习资料" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="可见性">
          <el-radio-group v-model="form.visibility">
            <el-radio :value="0">私密</el-radio>
            <el-radio :value="1">公开</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="onCreate">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import EmptyState from '@/components/EmptyState.vue'
import { useCollectionStore } from '@/stores/collection'
import { createCollection } from '@/api/collection'

const router = useRouter()
const store = useCollectionStore()
const createVisible = ref(false)
const saving = ref(false)
const loading = ref(false)
const collections = ref(store.collections)
const form = reactive({ name: '', description: '', visibility: 0 })

async function load() {
  loading.value = true
  try {
    collections.value = await store.fetchMine()
  } finally {
    loading.value = false
  }
}

async function onCreate() {
  if (!form.name.trim()) {
    ElMessage.warning('请输入收藏夹名称')
    return
  }
  saving.value = true
  try {
    await createCollection({ ...form, visibility: form.visibility })
    ElMessage.success('创建成功')
    createVisible.value = false
    form.name = ''
    form.description = ''
    form.visibility = 0
    await load()
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.collect-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.page-title {
  color: #303133;
}
.collect-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  min-height: 200px;
}
.collect-card {
  background: #fff;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  padding: 20px;
  cursor: pointer;
  transition: box-shadow 0.2s;
}
.collect-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}
.collect-icon {
  font-size: 28px;
  margin-bottom: 8px;
}
.collect-name {
  margin: 0 0 6px;
  color: #303133;
}
.collect-desc {
  color: #909399;
  font-size: 13px;
  min-height: 38px;
}
.collect-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #909399;
  font-size: 12px;
}
</style>
