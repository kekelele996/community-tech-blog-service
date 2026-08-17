// 分页 hooks：与后端 page / page_size 统一
import { ref } from 'vue'

export interface PaginationState {
  page: number
  pageSize: number
  total: number
}

export function usePagination(initPageSize = 10) {
  const pagination = ref<PaginationState>({ page: 1, pageSize: initPageSize, total: 0 })

  function setTotal(total: number) {
    pagination.value.total = total
  }

  function reset() {
    pagination.value.page = 1
    pagination.value.total = 0
  }

  return { pagination, setTotal, reset }
}
