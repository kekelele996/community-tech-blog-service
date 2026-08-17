<template>
  <div class="empty-state">
    <el-empty :description="description">
      <el-button v-if="actionText && showAction" type="primary" @click="$emit('action')">{{ actionText }}</el-button>
    </el-empty>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  description?: string
  actionText?: string
  showAction?: boolean
}>()

defineEmits<{ (e: 'action'): void }>()
</script>

<style scoped>
.empty-state {
  padding: 32px 0;
}
</style>
