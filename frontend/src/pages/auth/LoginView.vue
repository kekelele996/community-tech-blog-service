<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1 class="auth-title">📝 TechBlog 登录</h1>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="密码登录" name="password">
          <el-form :model="loginForm" label-width="0" @submit.prevent>
            <el-form-item>
              <el-input v-model="loginForm.email" placeholder="邮箱" size="large" />
            </el-form-item>
            <el-form-item>
              <el-input v-model="loginForm.password" type="password" placeholder="密码" size="large" show-password />
            </el-form-item>
            <el-button type="primary" size="large" style="width: 100%" :loading="loading" @click="onLogin">
              登 录
            </el-button>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="验证码登录" name="code">
          <el-form :model="codeForm" label-width="0">
            <el-form-item>
              <el-input v-model="codeForm.email" placeholder="邮箱" size="large" />
            </el-form-item>
            <el-form-item>
              <el-input v-model="codeForm.code" placeholder="6 位验证码" size="large">
                <template #append>
                  <el-button :disabled="countdown > 0" @click="onSendCode">
                    {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
                  </el-button>
                </template>
              </el-input>
            </el-form-item>
            <el-button type="primary" size="large" style="width: 100%" :loading="loading" @click="onCodeLogin">
              登 录
            </el-button>
          </el-form>
        </el-tab-pane>
      </el-tabs>
      <div class="auth-links">
        <router-link to="/register">还没有账号？去注册</router-link>
      </div>
      <el-button size="large" class="github-btn" @click="onGithubLogin">GitHub 第三方登录（演示）</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { sendCode } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const activeTab = ref('password')
const loading = ref(false)
const countdown = ref(0)
const loginForm = reactive({ email: '', password: '' })
const codeForm = reactive({ email: '', code: '' })

async function onLogin() {
  if (!loginForm.email || !loginForm.password) {
    ElMessage.warning('请输入邮箱和密码')
    return
  }
  loading.value = true
  try {
    await auth.login(loginForm)
    ElMessage.success('登录成功')
    router.push((route.query.redirect as string) || '/')
  } finally {
    loading.value = false
  }
}

async function onSendCode() {
  if (!codeForm.email) {
    ElMessage.warning('请输入邮箱')
    return
  }
  await sendCode(codeForm.email)
  ElMessage.success('验证码已发送（演示环境请查看接口返回的 debug_code）')
  countdown.value = 60
  const timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) clearInterval(timer)
  }, 1000)
}

async function onCodeLogin() {
  if (!codeForm.email || !codeForm.code) {
    ElMessage.warning('请输入邮箱和验证码')
    return
  }
  loading.value = true
  try {
    await auth.loginWithCode(codeForm)
    ElMessage.success('登录成功')
    router.push('/')
  } finally {
    loading.value = false
  }
}

async function onGithubLogin() {
  loading.value = true
  try {
    await auth.githubLogin({ code: 'demo-oauth-code-001', nickname: 'GitHub访客' })
    ElMessage.success('GitHub 登录成功（演示）')
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
  margin: 16px 0;
  font-size: 13px;
}
.github-btn {
  width: 100%;
  margin-top: 8px;
}
</style>
