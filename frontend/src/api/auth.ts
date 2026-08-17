import { get, post } from '@/utils/request'
import type { AuthResponse, SendCodeResponse, UserProfile } from '@/types'

export interface RegisterPayload {
  email: string
  code: string
  password: string
  nickname: string
  avatar?: string
  bio?: string
  tech_tags?: string[]
}

export interface LoginPayload {
  email: string
  password: string
}

export interface LoginCodePayload {
  email: string
  code: string
}

export interface GithubLoginPayload {
  code: string
  email?: string
  nickname?: string
}

export function sendCode(email: string) {
  return post<SendCodeResponse>('/auth/code', { email })
}

export function register(payload: RegisterPayload) {
  return post<AuthResponse>('/auth/register', payload)
}

export function login(payload: LoginPayload) {
  return post<AuthResponse>('/auth/login', payload)
}

export function loginWithCode(payload: LoginCodePayload) {
  return post<AuthResponse>('/auth/login/code', payload)
}

export function githubLogin(payload: GithubLoginPayload) {
  return post<AuthResponse>('/auth/github', payload)
}

export function getMe() {
  return get<UserProfile>('/auth/me')
}
