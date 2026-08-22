<template>
  <div class="auth-shell">
    <div class="auth-card">
      <div class="auth-copy">
        <span class="eyebrow">服务中心</span>
        <h1>家电维修服务中心</h1>
        <p>在一个控制台中统一跟踪受理、检测、配件、维修、质保和回访全流程。</p>
      </div>

      <el-form @submit.prevent="handleLogin">
        <el-form-item label="用户名">
          <el-input v-model="username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="password" show-password type="password" />
        </el-form-item>
        <el-alert
          type="warning"
          show-icon
          :closable="false"
          title="admin / Admin123! 仅用于本地开发验收，生产环境请立即修改。"
        />
        <div class="actions">
          <RouterLink to="/register">注册账号</RouterLink>
          <el-button type="primary" :loading="auth.loading" @click="handleLogin">登录</el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
const username = ref('admin')
const password = ref('Admin123!')

async function handleLogin() {
  try {
    await auth.login(username.value, password.value)
    await router.push('/')
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.message || '登录失败')
  }
}
</script>

<style scoped>
.auth-shell {
  display: grid;
  place-items: center;
  min-height: 100vh;
  padding: 24px;
}

.auth-card {
  width: min(960px, 100%);
  display: grid;
  grid-template-columns: 1.1fr 0.9fr;
  gap: 32px;
  padding: 36px;
  border-radius: 28px;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: var(--shadow);
}

.eyebrow {
  display: inline-flex;
  padding: 6px 10px;
  border-radius: 999px;
  background: var(--primary-soft);
  color: var(--primary);
  font-size: 12px;
  letter-spacing: 0.14em;
}

.auth-copy h1 {
  margin: 18px 0 12px;
  font-size: 42px;
  line-height: 1.05;
}

.auth-copy p {
  color: var(--muted);
  max-width: 460px;
}

.actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 20px;
}

@media (max-width: 880px) {
  .auth-card {
    grid-template-columns: 1fr;
  }
}
</style>
