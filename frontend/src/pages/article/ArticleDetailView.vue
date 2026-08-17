<template>
  <div v-loading="loading" class="detail">
    <template v-if="article">
      <article class="detail-card">
        <h1 class="detail-title">{{ article.title }}</h1>
        <div class="detail-meta">
          <UserAvatar :user="article.author" :size="28" />
          <router-link :to="`/users/${article.author_id}`" class="author-name">{{ article.author?.nickname || '匿名' }}</router-link>
          <span class="meta-time">{{ formatDate(article.published_at || article.created_at) }}</span>
          <StatusBadge v-if="isOwner" :value="article.status" :text-map="ARTICLE_STATUS_TEXT" :type-map="statusTypes" />
        </div>
        <div v-if="article.topics.length" class="detail-topics">
          <el-tag v-for="t in article.topics" :key="t.id" size="small" @click="router.push(`/topics/${t.id}`)">{{ t.name }}</el-tag>
        </div>
        <div class="markdown-body" v-html="html" />
        <div class="detail-actions">
          <el-button :type="article.is_liked ? 'danger' : 'default'" @click="onLike">
            {{ article.is_liked ? '已点赞' : '点赞' }} {{ article.like_count }}
          </el-button>
          <el-button :type="article.is_collected ? 'warning' : 'default'" @click="openCollect">
            {{ article.is_collected ? '已收藏' : '收藏' }} {{ article.favorite_count }}
          </el-button>
          <el-button v-if="isOwner && article.status !== 1" type="success" @click="onPublish">发布</el-button>
          <el-button v-if="isOwner && article.status === 1" type="warning" @click="onOffline">下架</el-button>
          <el-button v-if="isOwner" @click="router.push(`/editor/${article.id}`)">编辑</el-button>
        </div>
      </article>

      <section class="comments-card">
        <h3>评论（{{ comments.length }}）</h3>
        <div v-if="isLoggedIn" class="comment-form">
          <el-input v-model="commentText" type="textarea" :rows="3" placeholder="写下你的评论…" />
          <el-button type="primary" class="comment-submit" @click="onComment">发表评论</el-button>
        </div>
        <div v-else class="comment-tip"><router-link to="/login">登录后参与评论</router-link></div>
        <div class="comment-list">
          <div v-for="c in comments" :key="c.id" class="comment-item">
            <UserAvatar :user="c.author" :size="24" />
            <div class="comment-body">
              <div class="comment-head">
                <span class="comment-author">{{ c.author?.nickname || '匿名' }}</span>
                <span class="comment-time">{{ formatRelativeTime(c.created_at) }}</span>
              </div>
              <p class="comment-content">{{ c.content }}</p>
            </div>
          </div>
          <EmptyState v-if="!comments.length" description="暂无评论" />
        </div>
      </section>
    </template>
    <EmptyState v-else-if="!loading" description="文章不存在或已下架" />

    <el-dialog v-model="collectVisible" title="收藏到收藏夹" width="460px">
      <div class="collect-list">
        <div v-for="c in myCollections" :key="c.id" class="collect-item" :class="{ checked: c.has_article }" @click="onToggleCollect(c)">
          <span>{{ c.name }}（{{ c.article_count }}）</span>
          <el-tag v-if="c.has_article" type="success" size="small">已收藏</el-tag>
        </div>
        <div v-if="!myCollections.length" class="collect-empty">还没有收藏夹，<el-link type="primary" @click="router.push('/collections')">去创建</el-link></div>
      </div>
      <template #footer>
        <el-button @click="collectVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import ArticleCard from '@/components/ArticleCard.vue'
import EmptyState from '@/components/EmptyState.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import { ARTICLE_STATUS, ARTICLE_STATUS_TEXT } from '@/constants'
import { getArticle, likeArticle, unlikeArticle, publishArticle, offlineArticle } from '@/api/article'
import { listComments, createComment } from '@/api/comment'
import { listMyCollections, addArticleToCollection, removeArticleFromCollection } from '@/api/collection'
import { useAuth } from '@/hooks/useAuth'
import { formatDate, formatRelativeTime } from '@/utils/format'
import { renderMarkdown } from '@/utils/markdown'
import type { ArticleDetail, Collection, Comment } from '@/types'

const route = useRoute()
const router = useRouter()
const { isLoggedIn, user } = useAuth()

const article = ref<ArticleDetail | null>(null)
const comments = ref<Comment[]>([])
const myCollections = ref<Array<Collection & { has_article?: boolean }>>([])
const loading = ref(false)
const commentText = ref('')
const collectVisible = ref(false)

const isOwner = computed(() => article.value?.author_id === user.value?.id)
const html = computed(() => (article.value ? renderMarkdown(article.value.content) : ''))
const statusTypes: Record<number, 'success' | 'info' | 'warning' | 'danger'> = {
  [ARTICLE_STATUS.DRAFT]: 'info',
  [ARTICLE_STATUS.PUBLISHED]: 'success',
  [ARTICLE_STATUS.OFFLINE]: 'danger',
}

async function load() {
  const id = Number(route.params.id)
  loading.value = true
  try {
    article.value = await getArticle(id)
    const res = await listComments(id)
    comments.value = res.items
    if (isLoggedIn.value) {
      const cols = await listMyCollections()
      myCollections.value = cols.map((c) => ({
        ...c,
        has_article: (article.value?.collections || []).find((x) => x.id === c.id)?.has_article ?? false,
      }))
    }
  } finally {
    loading.value = false
  }
}

async function onLike() {
  if (!isLoggedIn.value) {
    ElMessage.warning('请先登录')
    return
  }
  if (!article.value) return
  if (article.value.is_liked) {
    await unlikeArticle(article.value.id)
    article.value.is_liked = false
    article.value.like_count--
  } else {
    await likeArticle(article.value.id)
    article.value.is_liked = true
    article.value.like_count++
  }
}

function openCollect() {
  if (!isLoggedIn.value) {
    ElMessage.warning('请先登录')
    return
  }
  collectVisible.value = true
}

async function onToggleCollect(c: Collection & { has_article?: boolean }) {
  if (!article.value) return
  if (c.has_article) {
    await removeArticleFromCollection(c.id, article.value.id)
    c.has_article = false
    c.article_count--
    article.value.is_collected = myCollections.value.some((x) => x.has_article)
    article.value.favorite_count = Math.max(0, article.value.favorite_count - 1)
    ElMessage.success('已取消收藏')
  } else {
    await addArticleToCollection(c.id, article.value.id)
    c.has_article = true
    c.article_count++
    article.value.is_collected = true
    article.value.favorite_count++
    ElMessage.success('收藏成功')
  }
}

async function onComment() {
  if (!article.value || !commentText.value.trim()) return
  const res = await createComment(article.value.id, commentText.value.trim())
  comments.value.push(res)
  article.value.comment_count++
  commentText.value = ''
}

async function onPublish() {
  if (!article.value) return
  const updated = await publishArticle(article.value.id)
  article.value.status = updated.status
  ElMessage.success('发布成功')
}

async function onOffline() {
  if (!article.value) return
  await offlineArticle(article.value.id)
  article.value.status = ARTICLE_STATUS.OFFLINE
  ElMessage.success('已下架')
}

onMounted(load)
</script>

<style scoped>
.detail {
  max-width: 900px;
  margin: 0 auto;
}
.detail-card,
.comments-card {
  background: #fff;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  padding: 24px;
  margin-bottom: 16px;
}
.detail-title {
  font-size: 26px;
  color: #303133;
  margin: 0 0 12px;
}
.detail-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}
.author-name {
  color: var(--el-color-primary);
  text-decoration: none;
  font-size: 14px;
}
.meta-time {
  color: #909399;
  font-size: 13px;
}
.detail-topics {
  margin-bottom: 16px;
}
.markdown-body {
  line-height: 1.8;
  color: #303133;
  word-break: break-word;
}
.detail-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}
.comments-card h3 {
  margin: 0 0 16px;
}
.comment-form {
  margin-bottom: 16px;
}
.comment-submit {
  margin-top: 8px;
}
.comment-tip {
  padding: 12px 0;
  color: #909399;
  font-size: 13px;
}
.comment-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.comment-item {
  display: flex;
  gap: 10px;
}
.comment-head {
  display: flex;
  gap: 8px;
  align-items: center;
}
.comment-author {
  font-size: 13px;
  font-weight: 600;
  color: #606266;
}
.comment-time {
  font-size: 12px;
  color: #c0c4cc;
}
.comment-content {
  margin: 4px 0 0;
  font-size: 14px;
  color: #303133;
}
.collect-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.collect-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 6px;
  cursor: pointer;
}
.collect-item.checked {
  border-color: var(--el-color-primary);
  background: #ecf5ff;
}
.collect-empty {
  text-align: center;
  color: #909399;
  padding: 12px 0;
}
</style>
