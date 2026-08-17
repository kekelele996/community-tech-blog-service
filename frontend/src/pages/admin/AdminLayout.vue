<template>
  <div class="admin-layout">
    <el-menu mode="horizontal" :default-active="activePath" router class="admin-menu">
      <el-menu-item index="/admin/dashboard">数据统计</el-menu-item>
      <el-menu-item index="/admin/users">用户管理</el-menu-item>
      <el-menu-item index="/admin/articles">文章管理</el-menu-item>
      <el-menu-item index="/admin/topics">话题管理</el-menu-item>
      <el-menu-item index="/admin/audit">审计日志</el-menu-item>
    </el-menu>
    <div class="admin-content">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const activePath = computed(() => route.path)
</script>

<style scoped>
.admin-menu {
  margin-bottom: 20px;
}
.admin-content {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  border: 1px solid var(--el-border-color-lighter);
  min-height: 400px;
}
</style>
