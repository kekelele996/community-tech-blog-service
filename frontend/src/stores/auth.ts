import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as authApi from '@/api/auth'
import { useAuth } from '@/hooks/useAuth'
import type { UserProfile } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const { user, isLoggedIn, isAdmin, setAuth, setUser, clearAuth } = useAuth()
  const loading = ref(false)

  async function login(payload: authApi.LoginPayload) {
    loading.value = true
    try {
      const res = await authApi.login(payload)
      setAuth(res.token, res.user)
      return res.user
    } finally {
      loading.value = false
    }
  }

  async function register(payload: authApi.RegisterPayload) {
    loading.value = true
    try {
      const res = await authApi.register(payload)
      setAuth(res.token, res.user)
      return res.user
    } finally {
      loading.value = false
    }
  }

  async function loginWithCode(payload: authApi.LoginCodePayload) {
    loading.value = true
    try {
      const res = await authApi.loginWithCode(payload)
      setAuth(res.token, res.user)
      return res.user
    } finally {
      loading.value = false
    }
  }

  async function githubLogin(payload: authApi.GithubLoginPayload) {
    const res = await authApi.githubLogin(payload)
    setAuth(res.token, res.user)
    return res.user
  }

  async function fetchMe(): Promise<UserProfile | null> {
    if (!isLoggedIn.value) return null
    try {
      const me = await authApi.getMe()
      setUser(me)
      return me
    } catch {
      return user.value
    }
  }

  function logout() {
    clearAuth()
  }

  return { user, isLoggedIn, isAdmin, loading, login, register, loginWithCode, githubLogin, fetchMe, logout }
})
