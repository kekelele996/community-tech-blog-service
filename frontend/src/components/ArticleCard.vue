<template>
  <div class="article-card" @click="$emit('open', article.id)">
    <div class="article-card__cover" v-if="article.cover_url">
      <el-image :src="article.cover_url" fit="cover" class="cover-img" />
    </div>
    <div class="article-card__body">
      <h3 class="article-card__title">{{ article.title }}</h3>
      <p class="article-card__summary">{{ article.summary || '暂无摘要' }}</p>
      <div class="article-card__meta">
        <span class="meta-author">
          <UserAvatar :user="article.author" :size="20" />
          <span>{{ article.author?.nickname || '匿名' }}</span>
        </span>
        <el-tag v-for="t in article.topics" :key="t.id" size="small" class="meta-topic">{{ t.name }}</el-tag>
      </div>
      <div class="article-card__stats">
        <span>👁 {{ article.view_count }}</span>
        <span>👍 {{ article.like_count }}</span>
        <span>⭐ {{ article.favorite_count }}</span>
        <span>💬 {{ article.comment_count }}</span>
        <span class="stats-time">{{ formatRelativeTime(article.published_at || article.created_at) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import UserAvatar from '@/components/UserAvatar.vue'
import { formatRelativeTime } from '@/utils/format'
import type { ArticleItem } from '@/types'

defineProps<{ article: ArticleItem }>()
defineEmits<{ (e: 'open', id: number): void }>()
</script>

<style scoped>
.article-card {
  display: flex;
  gap: 16px;
  padding: 16px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  background: #fff;
  cursor: pointer;
  transition: box-shadow 0.2s;
}
.article-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}
.article-card__cover {
  width: 160px;
  height: 100px;
  flex-shrink: 0;
  border-radius: 6px;
  overflow: hidden;
}
.cover-img {
  width: 100%;
  height: 100%;
}
.article-card__body {
  flex: 1;
  min-width: 0;
}
.article-card__title {
  margin: 0 0 6px;
  font-size: 16px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.article-card__summary {
  margin: 0 0 10px;
  font-size: 13px;
  color: #909399;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.article-card__meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}
.meta-author {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #606266;
}
.article-card__stats {
  display: flex;
  gap: 14px;
  font-size: 12px;
  color: #909399;
}
.stats-time {
  margin-left: auto;
}
</style>
