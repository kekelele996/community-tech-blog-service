<template>
  <div class="pagination-bar">
    <el-pagination
      :current-page="modelValue.page"
      :page-size="modelValue.pageSize"
      :total="modelValue.total"
      layout="prev, pager, next, total"
      @current-change="onChange"
    />
  </div>
</template>

<script setup lang="ts">
import type { PaginationState } from '@/hooks/usePagination'

defineProps<{ modelValue: PaginationState }>()
const emit = defineEmits<{ (e: 'update:modelValue', v: PaginationState): void; (e: 'change', page: number): void }>()

function onChange(page: number) {
  emit('change', page)
}
</script>

<style scoped>
.pagination-bar {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}
</style>
