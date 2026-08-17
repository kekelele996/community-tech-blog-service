<template>
  <el-container class="main-layout">
    <el-header class="main-header">
      <div class="header-inner">
        <router-link to="/" class="logo">📝 TechBlog</router-link>
        <nav class="nav">
          <router-link to="/topics">话题广场</router-link>
          <template v-if="auth.isLoggedIn">
            <router-link to="/collections">收藏夹</router-link>
            <router-link to="/notifications">
              通知<el-badge v-if="notif.unreadCount > 0" :value="notif.unreadCount" class="notif-badge" />
            </router-link>
            <router-link v-if="auth.isAdmin" to="/admin">后台管理</router-link>
            <router-link to="/editor" class="nav-write">✍ 写文章</router-link>
            <el-dropdown @command="onUserCommand">
              <span class="user-entry">
                <UserAvatar :user="auth.user" :size="28" />
                <span>{{ auth.user?.nickname }}</span>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人主页</el-dropdown-item>
                  <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <router-link to="/login">登录</router-link>
            <router-link to="/register">注册</router-link>
          </template>
        </nav>
      </div>
    </el-header>
    <el-main class="main-content">
      <router-view />
    </el-main>
    <el-footer class="main-footer">技术写作博客内容社区 · Go 1.22 + Gin + Vue 3</el-footer>
  </el-container>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import UserAvatar from '@/components/UserAvatar.vue'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'

const router = useRouter()
const auth = useAuthStore()
const notif = useNotificationStore()

onMounted(async () => {
  if (auth.isLoggedIn) {
    await auth.fetchMe()
    notif.refreshUnread().catch(() => undefined)
  }
})

async function onUserCommand(command: string) {
  if (command === 'logout') {
    auth.logout()
    router.push('/login')
  } else if (command === 'profile') {
    if (auth.user) router.push(`/users/${auth.user.id}`)
  }
}
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
  background: #f5f7fa;
}
.main-header {
  background: #fff;
  border-bottom: 1px solid var(--el-border-color-lighter);
  padding: 0;
  position: sticky;
  top: 0;
  z-index: 10;
}
.header-inner {
  max-width: 1200px;
  margin: 0 auto;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}
.logo {
  font-size: 20px;
  font-weight: 700;
  color: var(--el-color-primary);
  text-decoration: none;
}
.nav {
  display: flex;
  align-items: center;
  gap: 20px;
}
.nav a {
  color: #606266;
  text-decoration: none;
  font-size: 14px;
}
.nav a:hover {
  color: var(--el-color-primary);
}
.nav-write {
  color: var(--el-color-primary) !important;
  font-weight: 600;
}
.notif-badge {
  margin-left: 4px;
}
.user-entry {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  color: #606266;
}
.main-content {
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
  padding: 24px;
  box-sizing: border-box;
}
.main-footer {
  text-align: center;
  color: #909399;
  font-size: 13px;
  padding: 16px;
}
</style>
