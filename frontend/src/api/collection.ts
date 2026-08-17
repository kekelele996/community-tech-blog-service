import { del, get, post, put } from '@/utils/request'
import type { Collection, CollectionDetail } from '@/types'

export interface CreateCollectionPayload {
  name: string
  description?: string
  visibility?: number
}

export function listMyCollections() {
  return get<Collection[]>('/collections')
}

export function listUserCollections(userId: number) {
  return get<Collection[]>(`/users/${userId}/collections`)
}

export function getCollection(id: number) {
  return get<CollectionDetail>(`/collections/${id}`)
}

export function createCollection(payload: CreateCollectionPayload) {
  return post<Collection>('/collections', payload)
}

export function updateCollection(id: number, payload: CreateCollectionPayload) {
  return put<Collection>(`/collections/${id}`, payload)
}

export function deleteCollection(id: number) {
  return del<{ id: number }>(`/collections/${id}`)
}

export function addArticleToCollection(collectionId: number, articleId: number) {
  return post<{ collection_id: number; article_id: number }>(`/collections/${collectionId}/articles`, { article_id: articleId })
}

export function removeArticleFromCollection(collectionId: number, articleId: number) {
  return del<{ collection_id: number; article_id: number }>(`/collections/${collectionId}/articles/${articleId}`)
}
