<template>
  <div class="page-shell">
    <div class="page-header">
      <div>
        <h2 class="page-title">个人资料</h2>
        <p class="page-subtitle">查看当前账号信息与权限列表。</p>
      </div>
    </div>
    <el-card shadow="never">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="用户名">{{ auth.user?.username }}</el-descriptions-item>
        <el-descriptions-item label="显示名称">{{ auth.user?.display_name }}</el-descriptions-item>
        <el-descriptions-item label="权限列表">{{ (auth.user?.permissions || []).join(', ') }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()

onMounted(async () => {
  if (!auth.user) {
    await auth.fetchMe()
  }
})
</script>
