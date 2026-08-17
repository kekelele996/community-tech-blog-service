import { del, get, post } from '@/utils/request'
import type { Comment, PageResult } from '@/types'

export function listComments(articleId: number, page = 1, pageSize = 20) {
  return get<PageResult<Comment>>(`/articles/${articleId}/comments`, { page, page_size: pageSize })
}

export function createComment(articleId: number, content: string, parentId = 0) {
  return post<Comment>(`/articles/${articleId}/comments`, { content, parent_id: parentId })
}

export function deleteComment(id: number) {
  return del<{ id: number }>(`/comments/${id}`)
}
