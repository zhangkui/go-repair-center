<template>
  <div class="auth-shell">
    <div class="auth-card single">
      <h1>注册账号</h1>
      <p class="page-subtitle">新注册用户默认授予日常操作人员角色。</p>
      <el-form label-position="top">
        <el-form-item label="用户名">
          <el-input v-model="form.username" />
        </el-form-item>
        <el-form-item label="显示名称">
          <el-input v-model="form.display_name" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" show-password />
        </el-form-item>
        <el-space>
          <RouterLink to="/login">返回登录</RouterLink>
          <el-button type="primary" @click="submit">注册</el-button>
        </el-space>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
const form = reactive({
  username: '',
  display_name: '',
  password: ''
})

async function submit() {
  try {
    await auth.register(form)
    ElMessage.success('注册成功')
    await router.push('/login')
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.message || '注册失败')
  }
}
</script>

<style scoped>
.single {
  max-width: 520px;
  margin: 0 auto;
}
</style>
