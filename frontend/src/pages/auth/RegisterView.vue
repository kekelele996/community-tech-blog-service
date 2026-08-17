<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1 class="auth-title">📝 TechBlog 注册</h1>
      <el-form :model="form" label-width="0">
        <el-form-item>
          <el-input v-model="form.email" placeholder="邮箱" size="large" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.nickname" placeholder="昵称（2-32 字符）" size="large" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.password" type="password" placeholder="密码（至少 6 位）" size="large" show-password />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.code" placeholder="邮箱验证码" size="large">
            <template #append>
              <el-button :disabled="countdown > 0" @click="onSendCode">
                {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
              </el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model="form.tech_tags" multiple placeholder="技术标签（可选）" size="large" style="width: 100%">
            <el-option v-for="t in techTags" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
        <el-button type="primary" size="large" style="width: 100%" :loading="loading" @click="onRegister">
          注 册
        </el-button>
      </el-form>
      <div class="auth-links">
        <router-link to="/login">已有账号？去登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { sendCode } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const techTags = ['前端', '后端', 'AI', 'Go', 'Java', 'DevOps', '数据库', '架构']
const loading = ref(false)
const countdown = ref(0)
const form = reactive({ email: '', nickname: '', password: '', code: '', tech_tags: [] as string[] })

async function onSendCode() {
  if (!form.email) {
    ElMessage.warning('请输入邮箱')
    return
  }
  await sendCode(form.email)
  ElMessage.success('验证码已发送（演示环境请查看接口返回的 debug_code）')
  countdown.value = 60
  const timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) clearInterval(timer)
  }, 1000)
}

async function onRegister() {
  if (!form.email || !form.nickname || !form.password || !form.code) {
    ElMessage.warning('请完整填写注册信息')
    return
  }
  loading.value = true
  try {
    await auth.register({ ...form })
    ElMessage.success('注册成功')
    router.push('/')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #eef2f7 0%, #e3edf7 100%);
}
.auth-card {
  width: 420px;
  background: #fff;
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
}
.auth-title {
  text-align: center;
  margin-bottom: 24px;
  color: #303133;
}
.auth-links {
  text-align: center;
  margin-top: 16px;
  font-size: 13px;
}
</style>
