<template>
  <div class="editor">
    <h1 class="editor-title">{{ isEdit ? '编辑文章' : '写文章' }}</h1>
    <el-form label-width="0">
      <el-form-item>
        <el-input v-model="form.title" placeholder="文章标题" size="large" />
      </el-form-item>
      <el-form-item>
        <div class="editor-meta">
          <el-select v-model="form.topic_ids" multiple placeholder="关联话题" style="width: 320px">
            <el-option v-for="t in topics" :key="t.id" :label="t.name" :value="t.id" />
          </el-select>
          <el-input v-model="form.cover_url" placeholder="封面图 URL（可选）" style="width: 320px" />
          <el-upload :show-file-list="false" :http-request="onUpload" accept="image/*">
            <el-button>上传图片</el-button>
          </el-upload>
        </div>
      </el-form-item>
      <el-form-item>
        <el-input v-model="form.summary" placeholder="摘要（可选，默认取正文前 120 字）" />
      </el-form-item>
      <el-form-item>
        <div class="editor-panels">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="18"
            placeholder="使用 Markdown 写作，支持代码高亮…"
            class="editor-input"
          />
          <div class="editor-preview">
            <div class="preview-head">预览</div>
            <div class="markdown-body" v-html="previewHtml" />
          </div>
        </div>
      </el-form-item>
      <div class="editor-actions">
        <el-button type="primary" :loading="saving" @click="onSave(false)">保存草稿</el-button>
        <el-button type="success" :loading="saving" @click="onSave(true)">保存并发布</el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { createArticle, updateArticle, getArticle, publishArticle } from '@/api/article'
import { listTopics } from '@/api/topic'
import { uploadImage } from '@/api/upload'
import { renderMarkdown } from '@/utils/markdown'
import type { UploadRequestOptions } from 'element-plus'
import type { Topic } from '@/types'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => Boolean(route.params.id))
const topics = ref<Topic[]>([])
const saving = ref(false)
const form = reactive({
  title: '',
  content: '',
  summary: '',
  cover_url: '',
  topic_ids: [] as number[],
})

const previewHtml = computed(() => renderMarkdown(form.content))

async function loadTopics() {
  const res = await listTopics({ all: 1 })
  topics.value = Array.isArray(res) ? res : res.items
}

async function loadArticle() {
  if (!isEdit.value) return
  const id = Number(route.params.id)
  const detail = await getArticle(id)
  form.title = detail.title
  form.content = detail.content
  form.summary = detail.summary
  form.cover_url = detail.cover_url
  form.topic_ids = detail.topics.map((t) => t.id)
}

async function onSave(publish: boolean) {
  if (!form.title.trim() || !form.content.trim()) {
    ElMessage.warning('标题和内容不能为空')
    return
  }
  saving.value = true
  try {
    const payload = {
      title: form.title.trim(),
      content: form.content.trim(),
      summary: form.summary.trim() || form.content.trim().slice(0, 120),
      cover_url: form.cover_url,
      topic_ids: form.topic_ids,
      status: publish ? 1 : 0,
    }
    let id: number
    let status = 0
    if (isEdit.value) {
      const item = await updateArticle(Number(route.params.id), payload)
      id = item.id
      status = item.status
    } else {
      const item = await createArticle(payload)
      id = item.id
      status = item.status
    }
    if (publish && status !== 1) {
      await publishArticle(id)
    }
    ElMessage.success(publish ? '已发布' : '草稿已保存')
    router.push(`/articles/${id}`)
  } finally {
    saving.value = false
  }
}

async function onUpload(options: UploadRequestOptions) {
  try {
    const url = await uploadImage(options.file as File)
    if (form.cover_url) {
      form.content += `\n\n![图片](${url})\n`
    } else {
      form.cover_url = url
      ElMessage.success('封面上传成功')
    }
  } catch {
    ElMessage.error('图片上传失败')
  }
}

onMounted(() => {
  loadTopics()
  loadArticle()
})
</script>

<style scoped>
.editor {
  max-width: 1000px;
  margin: 0 auto;
}
.editor-title {
  font-size: 22px;
  color: #303133;
}
.editor-meta {
  display: flex;
  gap: 12px;
  width: 100%;
  flex-wrap: wrap;
}
.editor-panels {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  width: 100%;
}
.editor-input :deep(textarea) {
  font-family: 'SFMono-Regular', Consolas, monospace;
  font-size: 13px;
  line-height: 1.7;
}
.editor-preview {
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  overflow: hidden;
}
.preview-head {
  padding: 8px 12px;
  background: #f5f7fa;
  font-size: 13px;
  color: #909399;
  border-bottom: 1px solid var(--el-border-color-lighter);
}
.markdown-body {
  padding: 12px;
  max-height: 480px;
  overflow: auto;
  line-height: 1.7;
}
.editor-actions {
  display: flex;
  gap: 12px;
}
</style>
