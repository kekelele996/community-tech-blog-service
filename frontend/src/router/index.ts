import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'home', component: () => import('@/pages/home/HomeView.vue'), meta: { title: '首页' } },
  { path: '/login', name: 'login', component: () => import('@/pages/auth/LoginView.vue'), meta: { title: '登录' } },
  { path: '/register', name: 'register', component: () => import('@/pages/auth/RegisterView.vue'), meta: { title: '注册' } },
  { path: '/articles/:id', name: 'article-detail', component: () => import('@/pages/article/ArticleDetailView.vue'), meta: { title: '文章详情' } },
  { path: '/editor', name: 'article-editor', component: () => import('@/pages/article/ArticleEditorView.vue'), meta: { title: '写文章', requiresAuth: true } },
  { path: '/editor/:id', name: 'article-edit', component: () => import('@/pages/article/ArticleEditorView.vue'), meta: { title: '编辑文章', requiresAuth: true } },
  { path: '/users/:id', name: 'user-profile', component: () => import('@/pages/user/ProfileView.vue'), meta: { title: '个人主页' } },
  { path: '/topics', name: 'topics', component: () => import('@/pages/topic/TopicView.vue'), meta: { title: '话题广场' } },
  { path: '/topics/:id', name: 'topic-detail', component: () => import('@/pages/topic/TopicDetailView.vue'), meta: { title: '话题详情' } },
  { path: '/notifications', name: 'notifications', component: () => import('@/pages/notification/NotificationView.vue'), meta: { title: '通知中心', requiresAuth: true } },
  { path: '/collections', name: 'collections', component: () => import('@/pages/collection/CollectionView.vue'), meta: { title: '我的收藏夹', requiresAuth: true } },
  { path: '/collections/:id', name: 'collection-detail', component: () => import('@/pages/collection/CollectionDetailView.vue'), meta: { title: '收藏夹详情', requiresAuth: true } },
  { path: '/admin', name: 'admin', component: () => import('@/pages/admin/AdminLayout.vue'), meta: { title: '后台管理', requiresAuth: true, requiresAdmin: true },
    children: [
      { path: '', redirect: '/admin/dashboard' },
      { path: 'dashboard', component: () => import('@/pages/admin/AdminDashboard.vue'), meta: { title: '数据统计' } },
      { path: 'users', component: () => import('@/pages/admin/AdminUsers.vue'), meta: { title: '用户管理' } },
      { path: 'articles', component: () => import('@/pages/admin/AdminArticles.vue'), meta: { title: '文章管理' } },
      { path: 'topics', component: () => import('@/pages/admin/AdminTopics.vue'), meta: { title: '话题管理' } },
      { path: 'audit', component: () => import('@/pages/admin/AdminAudit.vue'), meta: { title: '审计日志' } },
    ] },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 路由守卫：登录与管理员权限（横切关注点 1 前端层）
router.beforeEach((to) => {
  const token = localStorage.getItem('techblog_token')
  const rawUser = localStorage.getItem('techblog_user')
  let isAdmin = false
  if (rawUser) {
    try {
      isAdmin = (JSON.parse(rawUser) as { role?: number }).role === 2
    } catch {
      isAdmin = false
    }
  }
  if (to.meta.requiresAuth && !token) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  if (to.meta.requiresAdmin && !isAdmin) {
    return { path: '/' }
  }
  document.title = `${String(to.meta.title || '')} - TechBlog`
  return true
})

export default router
