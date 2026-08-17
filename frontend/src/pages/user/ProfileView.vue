<template>
  <div v-loading="loading" class="profile">
    <div class="profile-card">
      <UserAvatar :user="profile" :size="72" />
      <div class="profile-info">
        <h2 class="profile-name">{{ profile?.nickname }}</h2>
        <p class="profile-bio">{{ profile?.bio || '这个人很懒，什么都没写。' }}</p>
        <div class="profile-tags">
          <el-tag v-for="t in profile?.tech_tags || []" :key="t" size="small" type="info">{{ t }}</el-tag>
        </div>
        <div class="profile-stats">
          <span>文章 {{ profile?.article_count ?? 0 }}</span>
          <span>粉丝 {{ profile?.follower_count ?? 0 }}</span>
          <span>关注 {{ profile?.following_count ?? 0 }}</span>
        </div>
      </div>
      <div class="profile-actions">
        <el-button v-if="isSelf" @click="editDialog = true">编辑资料</el-button>
        <el-button v-else-if="isLoggedIn" :type="profile?.is_following ? 'info' : 'primary'" @click="onToggleFollow">
          {{ profile?.is_following ? '已关注' : '关注' }}
        </el-button>
        <router-link v-if="isSelf" to="/editor" class="write-link">
          <el-button type="primary">写文章</el-button>
        </router-link>
      </div>
    </div>

    <div class="profile-tabs">
      <el-tabs v-model="activeTab" @tab-change="loadTab">
        <el-tab-pane label="TA 的文章" name="articles">
          <div class="article-list">
            <ArticleCard v-for="a in articles" :key="a.id" :article="a" @open="goDetail" />
            <EmptyState v-if="!articles.length" description="暂无文章" />
          </div>
        </el-tab-pane>
        <el-tab-pane label="粉丝" name="followers">
          <div class="user-list">
            <div v-for="u in followers" :key="u.id" class="user-item" @click="router.push(`/users/${u.id}`)">
              <UserAvatar :user="u" :size="32" />
              <span class="user-name">{{ u.nickname }}</span>
            </div>
            <EmptyState v-if="!followers.length" description="暂无粉丝" />
          </div>
        </el-tab-pane>
        <el-tab-pane label="关注" name="following">
          <div class="user-list">
            <div v-for="u in following" :key="u.id" class="user-item" @click="router.push(`/users/${u.id}`)">
              <UserAvatar :user="u" :size="32" />
              <span class="user-name">{{ u.nickname }}</span>
            </div>
            <EmptyState v-if="!following.length" description="暂无关注" />
          </div>
        </el-tab-pane>
        <el-tab-pane label="收藏夹" name="collections">
          <div class="user-list">
            <div v-for="c in collections" :key="c.id" class="user-item" @click="router.push(`/collections/${c.id}`)">
              <span class="user-name">📁 {{ c.name }}（{{ c.article_count }}）</span>
              <el-tag size="small">{{ c.visibility === 1 ? '公开' : '私密' }}</el-tag>
            </div>
            <EmptyState v-if="!collections.length" description="暂无收藏夹" />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <el-dialog v-model="editDialog" title="编辑资料" width="480px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="昵称"><el-input v-model="editForm.nickname" /></el-form-item>
        <el-form-item label="简介"><el-input v-model="editForm.bio" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="技术标签">
          <el-select v-model="editForm.tech_tags" multiple style="width: 100%">
            <el-option v-for="t in techOptions" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="onSaveProfile">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import ArticleCard from '@/components/ArticleCard.vue'
import EmptyState from '@/components/EmptyState.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import { getUserProfile, getFollowers, getFollowing, followUser, unfollowUser, updateProfile } from '@/api/user'
import { listArticles } from '@/api/article'
import { listUserCollections } from '@/api/collection'
import { useAuth } from '@/hooks/useAuth'
import type { ArticleItem, Collection, UserProfile } from '@/types'

const route = useRoute()
const router = useRouter()
const { user, isLoggedIn } = useAuth()

const profile = ref<UserProfile | null>(null)
const articles = ref<ArticleItem[]>([])
const followers = ref<UserProfile[]>([])
const following = ref<UserProfile[]>([])
const collections = ref<Collection[]>([])
const activeTab = ref('articles')
const loading = ref(false)
const saving = ref(false)
const editDialog = ref(false)
const editForm = reactive({ nickname: '', bio: '', tech_tags: [] as string[] })
const techOptions = ['前端', '后端', 'AI', 'Go', 'Java', 'DevOps', '数据库', '架构']

const isSelf = computed(() => profile.value?.id === user.value?.id)

async function loadProfile() {
  loading.value = true
  try {
    const id = Number(route.params.id)
    profile.value = await getUserProfile(id)
    await loadArticles(id)
  } finally {
    loading.value = false
  }
}

async function loadArticles(id: number) {
  const res = await listArticles({ author_id: id, page: 1, page_size: 20 })
  articles.value = res.items
}

async function loadTab() {
  const id = Number(route.params.id)
  if (activeTab.value === 'followers') {
    const res = await getFollowers(id)
    followers.value = res.items
  } else if (activeTab.value === 'following') {
    const res = await getFollowing(id)
    following.value = res.items
  } else if (activeTab.value === 'collections') {
    collections.value = await listUserCollections(id)
  } else {
    await loadArticles(id)
  }
}

async function onToggleFollow() {
  if (!profile.value) return
  if (profile.value.is_following) {
    await unfollowUser(profile.value.id)
    profile.value.is_following = false
    profile.value.follower_count = Math.max(0, profile.value.follower_count - 1)
    ElMessage.success('已取消关注')
  } else {
    await followUser(profile.value.id)
    profile.value.is_following = true
    profile.value.follower_count++
    ElMessage.success('关注成功')
  }
}

function openEdit() {
  if (!profile.value) return
  editForm.nickname = profile.value.nickname
  editForm.bio = profile.value.bio
  editForm.tech_tags = [...profile.value.tech_tags]
  editDialog.value = true
}

async function onSaveProfile() {
  saving.value = true
  try {
    profile.value = await updateProfile({ ...editForm })
    ElMessage.success('资料已更新')
    editDialog.value = false
  } finally {
    saving.value = false
  }
}

function goDetail(id: number) {
  router.push(`/articles/${id}`)
}

onMounted(loadProfile)
</script>

<style scoped>
.profile-card {
  display: flex;
  align-items: flex-start;
  gap: 20px;
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  border: 1px solid var(--el-border-color-lighter);
  margin-bottom: 16px;
}
.profile-info {
  flex: 1;
}
.profile-name {
  margin: 0 0 8px;
  font-size: 22px;
  color: #303133;
}
.profile-bio {
  color: #606266;
  font-size: 14px;
  margin: 0 0 10px;
}
.profile-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  margin-bottom: 10px;
}
.profile-stats {
  display: flex;
  gap: 24px;
  color: #909399;
  font-size: 13px;
}
.profile-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.write-link {
  text-decoration: none;
}
.profile-tabs {
  background: #fff;
  border-radius: 8px;
  padding: 16px 24px;
  border: 1px solid var(--el-border-color-lighter);
}
.article-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.user-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.user-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  border-radius: 6px;
  cursor: pointer;
}
.user-item:hover {
  background: #f5f7fa;
}
.user-name {
  font-size: 14px;
  color: #303133;
}
</style>
