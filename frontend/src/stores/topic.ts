import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as topicApi from '@/api/topic'
import type { Topic } from '@/types'

export const useTopicStore = defineStore('topic', () => {
  const topics = ref<Topic[]>([])
  const loading = ref(false)

  async function fetchAll() {
    loading.value = true
    try {
      const res = await topicApi.listTopics({ all: 1 })
      topics.value = Array.isArray(res) ? res : res.items
      return topics.value
    } finally {
      loading.value = false
    }
  }

  return { topics, loading, fetchAll }
})
