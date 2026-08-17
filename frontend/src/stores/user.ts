import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as userApi from '@/api/user'
import { useAuth } from '@/hooks/useAuth'
import type { UserProfile } from '@/types'

export const useUserStore = defineStore('user', () => {
  const profile = ref<UserProfile | null>(null)

  async function fetchProfile(id: number) {
    profile.value = await userApi.getUserProfile(id)
    return profile.value
  }

  async function updateProfile(payload: userApi.UpdateProfilePayload) {
    profile.value = await userApi.updateProfile(payload)
    if (profile.value) {
      const { setUser } = useAuth()
      setUser(profile.value)
    }
    return profile.value
  }

  return { profile, fetchProfile, updateProfile }
})
