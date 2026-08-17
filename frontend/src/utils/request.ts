// 统一 axios 请求封装：拦截器统一处理 JWT、错误响应与错误提示（横切关注点 3 前端层）
import axios, { type AxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import { ERROR_MESSAGES } from '@/constants'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('techblog_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (response) => {
    const body = response.data
    if (body && typeof body.code === 'number' && body.code !== 0) {
      ElMessage.error(body.message || ERROR_MESSAGES[body.code] || '请求失败')
      if (body.code === 40100) {
        localStorage.removeItem('techblog_token')
        localStorage.removeItem('techblog_user')
        if (!location.pathname.startsWith('/login')) {
          location.href = '/login'
        }
      }
      return Promise.reject(new Error(body.message))
    }
    return response
  },
  (error) => {
    const status = error.response?.status
    const message = error.response?.data?.message || error.message || '网络错误'
    ElMessage.error(message)
    if (status === 401) {
      localStorage.removeItem('techblog_token')
      localStorage.removeItem('techblog_user')
      if (!location.pathname.startsWith('/login')) {
        location.href = '/login'
      }
    }
    return Promise.reject(error)
  },
)

export async function get<T>(url: string, params?: object): Promise<T> {
  const res = await request.get<{ code: number; message: string; data: T }>(url, { params })
  return res.data.data
}

export async function post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
  const res = await request.post<{ code: number; message: string; data: T }>(url, data, config)
  return res.data.data
}

export async function put<T>(url: string, data?: unknown): Promise<T> {
  const res = await request.put<{ code: number; message: string; data: T }>(url, data)
  return res.data.data
}

export async function del<T>(url: string): Promise<T> {
  const res = await request.delete<{ code: number; message: string; data: T }>(url)
  return res.data.data
}

export default request
