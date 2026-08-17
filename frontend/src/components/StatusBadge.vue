<template>
  <el-tag :type="tagType" size="small" effect="light">{{ text }}</el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  value: number
  textMap: Record<number, string>
  typeMap?: Record<number, 'success' | 'info' | 'warning' | 'danger' | 'primary'>
}>()

const text = computed(() => props.textMap[props.value] ?? '未知')
const tagType = computed(() => {
  if (props.typeMap && props.typeMap[props.value]) return props.typeMap[props.value]
  return props.value === 1 ? 'success' : 'info'
})
</script>
