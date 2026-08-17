// 格式化工具（与后端 util/formatters.go 对应）
export function formatDate(iso?: string): string {
  if (!iso) return '-'
  const d = new Date(iso)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

export function formatRelativeTime(iso?: string): string {
  if (!iso) return '-'
  const diff = Date.now() - new Date(iso).getTime()
  const minute = 60 * 1000
  const hour = 60 * minute
  const day = 24 * hour
  if (diff < minute) return '刚刚'
  if (diff < hour) return `${Math.floor(diff / minute)} 分钟前`
  if (diff < day) return `${Math.floor(diff / hour)} 小时前`
  if (diff < 7 * day) return `${Math.floor(diff / day)} 天前`
  return formatDate(iso)
}

export function hotScore(likeCount: number, viewCount: number, publishedAt?: string): number {
  let score = likeCount * 10 + viewCount
  if (publishedAt && Date.now() - new Date(publishedAt).getTime() < 24 * 60 * 60 * 1000) {
    score += 50
  }
  return score
}

export function parseTechTags(tags?: string[]): string {
  return tags && tags.length ? tags.join(' / ') : '未设置技术标签'
}
