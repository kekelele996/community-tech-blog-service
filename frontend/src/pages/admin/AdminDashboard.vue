<template>
  <div v-loading="loading">
    <h2 class="admin-title">平台数据统计</h2>
    <div class="stat-cards">
      <div class="stat-card"><div class="stat-num">{{ stats?.dau ?? 0 }}</div><div class="stat-label">今日日活</div></div>
      <div class="stat-card"><div class="stat-num">{{ stats?.total_articles ?? 0 }}</div><div class="stat-label">文章总数</div></div>
      <div class="stat-card"><div class="stat-num">{{ stats?.total_users ?? 0 }}</div><div class="stat-label">注册用户</div></div>
      <div class="stat-card"><div class="stat-num">{{ stats?.total_comments ?? 0 }}</div><div class="stat-label">评论总数</div></div>
    </div>

    <el-tabs>
      <el-tab-pane label="近 7 日日活趋势">
        <div class="chart">
          <div v-for="p in stats?.daily_active || []" :key="p.date" class="bar-col">
            <div class="bar" :style="{ height: barHeight(p.count, maxDau) }"></div>
            <span class="bar-label">{{ p.date.slice(5) }}</span>
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane label="近 7 日注册趋势">
        <div class="chart">
          <div v-for="p in stats?.new_users || []" :key="p.date" class="bar-col">
            <div class="bar" :style="{ height: barHeight(p.count, maxNew) }"></div>
            <span class="bar-label">{{ p.date.slice(5) }}</span>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { getStats } from '@/api/admin'
import type { AdminStats } from '@/types'

const stats = ref<AdminStats | null>(null)
const loading = ref(false)

const maxDau = computed(() => Math.max(1, ...(stats.value?.daily_active || []).map((p) => p.count)))
const maxNew = computed(() => Math.max(1, ...(stats.value?.new_users || []).map((p) => p.count)))

function barHeight(count: number, max: number): string {
  return `${Math.max(4, Math.round((count / max) * 200))}px`
}

onMounted(async () => {
  loading.value = true
  try {
    stats.value = await getStats()
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.admin-title {
  color: #303133;
  margin-top: 0;
}
.stat-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}
.stat-card {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 20px;
  text-align: center;
}
.stat-num {
  font-size: 28px;
  font-weight: 700;
  color: var(--el-color-primary);
}
.stat-label {
  margin-top: 6px;
  color: #909399;
  font-size: 13px;
}
.chart {
  display: flex;
  align-items: flex-end;
  gap: 16px;
  height: 240px;
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
}
.bar-col {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  flex: 1;
  height: 100%;
  gap: 8px;
}
.bar {
  width: 100%;
  max-width: 48px;
  background: var(--el-color-primary);
  border-radius: 4px 4px 0 0;
}
.bar-label {
  font-size: 12px;
  color: #909399;
}
</style>
