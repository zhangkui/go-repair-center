<template>
  <div class="page-shell">
    <div class="page-header">
      <div>
        <h2 class="page-title">会话管理</h2>
        <p class="page-subtitle">查看当前刷新令牌会话，并按需撤销。</p>
      </div>
      <el-space>
        <el-button @click="load">刷新</el-button>
        <el-button type="danger" plain @click="revokeAll">全部撤销</el-button>
      </el-space>
    </div>

    <div class="card-grid">
      <KpiCard label="有效会话" :value="summary?.active_session_count ?? 0" caption="仍在有效期内的刷新会话" />
      <KpiCard label="过期会话" :value="summary?.expired_session_count ?? 0" caption="已过期但尚有记录的会话" />
      <KpiCard label="历史会话" :value="summary?.sessions?.length ?? 0" caption="当前账号的全部刷新令牌记录" />
    </div>

    <el-card shadow="never" v-loading="loading">
      <template #header>我的会话</template>
      <el-table :data="summary?.sessions || []" stripe>
        <el-table-column prop="id" label="会话编号" />
        <el-table-column prop="created_at" label="创建时间" />
        <el-table-column prop="expires_at" label="过期时间" />
        <el-table-column prop="revoked_at" label="撤销时间" />
        <el-table-column label="操作" width="140">
          <template #default="{ row }">
            <el-button link type="danger" @click="revokeById(row.id)">撤销</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import client from '../api/client'
import KpiCard from '../components/KpiCard.vue'

const loading = ref(false)
const summary = ref<any>(null)

async function load() {
  loading.value = true
  try {
    const { data } = await client.get('/sessions/me')
    summary.value = data.data
  } finally {
    loading.value = false
  }
}

async function revokeById(sessionId: number) {
  await ElMessageBox.confirm(`确认撤销会话 #${sessionId} 吗？`, '提示', { type: 'warning' })
  await client.post(`/sessions/${sessionId}/revoke`)
  ElMessage.success('会话已撤销')
  await load()
}

async function revokeAll() {
  await ElMessageBox.confirm('确认撤销当前用户的全部会话吗？', '提示', { type: 'warning' })
  await client.post('/sessions/revoke-all', {})
  ElMessage.success('全部会话已撤销')
  await load()
}

load()
</script>
