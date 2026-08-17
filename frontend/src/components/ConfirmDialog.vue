<template>
  <el-dialog :model-value="modelValue" :title="title" width="420px" @update:model-value="onUpdate">
    <span>{{ message }}</span>
    <template #footer>
      <el-button @click="onUpdate(false)">取消</el-button>
      <el-button type="danger" :loading="loading" @click="$emit('confirm')">{{ confirmText }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: boolean
  title?: string
  message: string
  confirmText?: string
  loading?: boolean
}>()

const emit = defineEmits<{ (e: 'update:modelValue', v: boolean): void; (e: 'confirm'): void }>()

function onUpdate(v: boolean) {
  emit('update:modelValue', v)
}
</script>
