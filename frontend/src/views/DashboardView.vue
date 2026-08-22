<template>
  <div class="page-shell">
    <div class="page-header">
      <div>
        <h2 class="page-title">仪表盘</h2>
        <p class="page-subtitle">今日报修、维修负载和客户回访总体概览。</p>
      </div>
      <el-button @click="load">刷新</el-button>
    </div>

    <div class="card-grid">
      <KpiCard label="今日报修数" :value="summary?.today_repair_count ?? 0" caption="今天新建的报修单数量" />
      <KpiCard label="维修中" :value="summary?.in_progress_count ?? 0" caption="当前正在执行的维修任务" />
      <KpiCard label="待取机" :value="summary?.waiting_pickup_count ?? 0" caption="已完成维修、等待客户取机" />
      <KpiCard label="待回访" :value="summary?.pending_feedback_count ?? 0" caption="尚未完成的回访任务" />
    </div>

    <div class="two-col">
      <el-card shadow="never" v-loading="loading">
        <template #header>维修趋势</template>
        <el-table :data="summary?.repair_trend || []" stripe>
          <el-table-column prop="day" label="日期" />
          <el-table-column prop="count" label="数量" />
        </el-table>
      </el-card>

      <el-card shadow="never" v-loading="loading">
        <template #header>营收概览</template>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="总营收">{{ formatMoney(summary?.revenue?.total_amount) }}</el-descriptions-item>
          <el-descriptions-item label="近周期营收">{{ formatMoney(summary?.revenue?.period_amount) }}</el-descriptions-item>
          <el-descriptions-item label="完工率">{{ formatRate(summary?.completion_rate) }}</el-descriptions-item>
        </el-descriptions>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import client from '../api/client'
import KpiCard from '../components/KpiCard.vue'

const loading = ref(false)
const summary = ref<any>(null)

async function load() {
  loading.value = true
  try {
    const { data } = await client.get('/dashboard')
    summary.value = data.data
  } finally {
    loading.value = false
  }
}

function formatMoney(value?: number) {
  return `¥${Number(value || 0).toFixed(2)}`
}

function formatRate(value?: number) {
  return `${((value || 0) * 100).toFixed(1)}%`
}

load()
</script>
