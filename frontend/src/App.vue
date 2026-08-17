<template>
  <router-view v-if="isLayoutPage" />
  <MainLayout v-else />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'

const route = useRoute()
// 登录/注册页使用独立布局
const isLayoutPage = computed(() => ['/login', '/register'].includes(route.path))
</script>
