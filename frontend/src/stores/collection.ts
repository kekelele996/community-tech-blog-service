import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as collectionApi from '@/api/collection'
import type { Collection } from '@/types'

export const useCollectionStore = defineStore('collection', () => {
  const collections = ref<Collection[]>([])
  const loading = ref(false)

  async function fetchMine() {
    loading.value = true
    try {
      collections.value = await collectionApi.listMyCollections()
      return collections.value
    } finally {
      loading.value = false
    }
  }

  return { collections, loading, fetchMine }
})
