import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as articleApi from '@/api/article'
import type { ArticleDetail, ArticleItem, PageResult } from '@/types'

export const useArticleStore = defineStore('article', () => {
  const list = ref<ArticleItem[]>([])
  const total = ref(0)
  const detail = ref<ArticleDetail | null>(null)
  const loading = ref(false)

  async function fetchList(params: articleApi.ArticleQuery) {
    loading.value = true
    try {
      const res = await articleApi.listArticles(params)
      list.value = res.items
      total.value = res.total
      return res
    } finally {
      loading.value = false
    }
  }

  async function fetchDetail(id: number) {
    loading.value = true
    try {
      detail.value = await articleApi.getArticle(id)
      return detail.value
    } finally {
      loading.value = false
    }
  }

  return { list, total, detail, loading, fetchList, fetchDetail }
})
