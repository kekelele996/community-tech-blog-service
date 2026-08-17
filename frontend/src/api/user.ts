import { del, get, post, put } from '@/utils/request'
import type { PageResult, UserProfile } from '@/types'

export interface UpdateProfilePayload {
  nickname: string
  avatar?: string
  bio?: string
  tech_tags?: string[]
}

export function getUserProfile(id: number) {
  return get<UserProfile>(`/users/${id}`)
}

export function updateProfile(payload: UpdateProfilePayload) {
  return put<UserProfile>('/user/profile', payload)
}

export function getFollowers(id: number, page = 1, pageSize = 10) {
  return get<PageResult<UserProfile>>(`/users/${id}/followers`, { page, page_size: pageSize })
}

export function getFollowing(id: number, page = 1, pageSize = 10) {
  return get<PageResult<UserProfile>>(`/users/${id}/following`, { page, page_size: pageSize })
}

export function followUser(id: number) {
  return post<{ followed_id: number }>(`/users/${id}/follow`)
}

export function unfollowUser(id: number) {
  return del<{ followed_id: number }>(`/users/${id}/follow`)
}
