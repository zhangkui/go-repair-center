<template>
  <div class="page-shell">
    <div class="page-header">
      <div>
        <h2 class="page-title">运营中心</h2>
        <p class="page-subtitle">查看积压工单、SLA 预警、质保提醒和营收概览。</p>
      </div>
      <el-button @click="load">刷新</el-button>
    </div>

    <div class="card-grid">
      <KpiCard label="待处理工单" :value="backlog?.pending_orders ?? 0" caption="等待受理或处理的维修单" />
      <KpiCard label="待上门/接修" :value="backlog?.waiting_visit_orders ?? 0" caption="尚未完成上门或接修的工单" />
      <KpiCard label="待取机" :value="backlog?.waiting_pickup_orders ?? 0" caption="维修完成、等待客户取机" />
      <KpiCard label="低库存配件" :value="health?.low_stock_parts ?? 0" caption="库存低于或等于安全库存的配件" />
    </div>

    <div class="two-col">
      <el-card shadow="never" v-loading="loading">
        <template #header>营收汇总</template>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="近周期">{{ formatMoney(report?.revenue?.period_amount) }}</el-descriptions-item>
          <el-descriptions-item label="累计">{{ formatMoney(report?.revenue?.total_amount) }}</el-descriptions-item>
          <el-descriptions-item label="待回访">{{ backlog?.pending_feedbacks ?? 0 }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <el-card shadow="never" v-loading="loading">
        <template #header>运行快照</template>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="未关闭维修单">{{ health?.open_repair_orders ?? 0 }}</el-descriptions-item>
          <el-descriptions-item label="未关闭执行单">{{ health?.open_executions ?? 0 }}</el-descriptions-item>
          <el-descriptions-item label="已过期报价">{{ health?.expired_quotations ?? 0 }}</el-descriptions-item>
          <el-descriptions-item label="已过期质保">{{ health?.expired_warranties ?? 0 }}</el-descriptions-item>
        </el-descriptions>
      </el-card>
    </div>

    <el-card shadow="never" v-loading="loading">
      <template #header>SLA 预警</template>
      <el-table :data="slaAlerts" stripe>
        <el-table-column prop="resource_type" label="资源类型" />
        <el-table-column prop="resource_id" label="编号" />
        <el-table-column prop="status" label="状态" />
        <el-table-column prop="due_at" label="截止时间" />
        <el-table-column prop="severity" label="严重级别" />
      </el-table>
    </el-card>

    <el-card shadow="never" v-loading="loading">
      <template #header>即将到期质保</template>
      <el-table :data="warrantyExpiring" stripe>
        <el-table-column prop="id" label="质保编号" />
        <el-table-column prop="repair_order_id" label="维修单编号" />
        <el-table-column prop="end_date" label="到期时间" />
        <el-table-column prop="status" label="状态" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import client from '../api/client'
import KpiCard from '../components/KpiCard.vue'

const loading = ref(false)
const backlog = ref<any>(null)
const report = ref<any>(null)
const health = ref<any>(null)
const slaAlerts = ref<any[]>([])
const warrantyExpiring = ref<any[]>([])

async function load() {
  loading.value = true
  try {
    const [backlogResponse, reportResponse, healthResponse, slaResponse, warrantyResponse] = await Promise.all([
      client.get('/ops/backlog'),
      client.get('/reports/revenue', { params: { days: 30 } }),
      client.get('/ops/health'),
      client.get('/ops/sla-alerts'),
      client.get('/ops/warranty-expiring', { params: { within_days: 14 } })
    ])
    backlog.value = backlogResponse.data.data
    report.value = reportResponse.data.data
    health.value = healthResponse.data.data
    slaAlerts.value = slaResponse.data.data || []
    warrantyExpiring.value = warrantyResponse.data.data || []
  } finally {
    loading.value = false
  }
}

function formatMoney(value?: number) {
  return `¥${Number(value || 0).toFixed(2)}`
}

load()
</script>
