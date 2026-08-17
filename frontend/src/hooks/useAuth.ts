// 认证相关 hooks：登录态、管理员判断（基于共享响应式状态，随登录/登出即时更新）
import { computed, reactive } from 'vue'
import { ROLE } from '@/constants'
import type { UserProfile } from '@/types'

function readUser(): UserProfile | null {
  const raw = localStorage.getItem('techblog_user')
  if (!raw) return null
  try {
    return JSON.parse(raw) as UserProfile
  } catch {
    return null
  }
}

const state = reactive<{ token: string; user: UserProfile | null }>({
  token: localStorage.getItem('techblog_token') || '',
  user: readUser(),
})

export function useAuth() {
  const token = computed(() => state.token)
  const user = computed(() => state.user)
  const isLoggedIn = computed(() => Boolean(state.token))
  const isAdmin = computed(() => state.user?.role === ROLE.ADMIN)

  function setAuth(t: string, u: UserProfile) {
    state.token = t
    state.user = u
    localStorage.setItem('techblog_token', t)
    localStorage.setItem('techblog_user', JSON.stringify(u))
  }

  function setUser(u: UserProfile) {
    state.user = u
    localStorage.setItem('techblog_user', JSON.stringify(u))
  }

  function clearAuth() {
    state.token = ''
    state.user = null
    localStorage.removeItem('techblog_token')
    localStorage.removeItem('techblog_user')
  }

  return { token, user, isLoggedIn, isAdmin, setAuth, setUser, clearAuth }
}
