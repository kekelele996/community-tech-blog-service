<template>
  <el-avatar :size="size" :src="user?.avatar || ''" class="user-avatar">
    {{ fallbackText }}
  </el-avatar>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { UserProfile } from '@/types'

const props = defineProps<{
  user?: UserProfile | null
  size?: number
}>()

const fallbackText = computed(() => {
  if (!props.user?.nickname) return 'U'
  return props.user.nickname.slice(0, 1).toUpperCase()
})
</script>

<style scoped>
.user-avatar {
  background-color: var(--el-color-primary);
  color: #fff;
  font-weight: 600;
}
</style>
