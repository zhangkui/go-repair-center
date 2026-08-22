<template>
  <div class="page-shell">
    <div class="page-header">
      <div>
        <h2 class="page-title">库存管控</h2>
        <p class="page-subtitle">监控低库存、生成补货建议，并查看库存变动时间线。</p>
      </div>
      <el-button @click="load">刷新</el-button>
    </div>

    <div class="card-grid">
      <KpiCard label="低库存配件" :value="summary?.low_stock_count ?? 0" caption="低于或等于安全库存的配件数量" />
      <KpiCard label="补货建议数" :value="summary?.restock_item_count ?? 0" caption="系统建议补货的配件条目数" />
      <KpiCard label="预计成本" :value="formatMoney(summary?.estimated_restock_cost)" caption="本轮建议补货的预估预算" />
    </div>

    <el-card shadow="never">
      <template #header>手工补货</template>
      <el-form inline>
        <el-form-item label="配件编号">
          <el-input v-model.number="restockForm.partId" type="number" style="width: 140px" />
        </el-form-item>
        <el-form-item label="补货数量">
          <el-input v-model.number="restockForm.quantity" type="number" style="width: 140px" />
        </el-form-item>
        <el-form-item label="补货原因">
          <el-input v-model="restockForm.reason" style="width: 280px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitRestock">提交</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <div class="two-col">
      <el-card shadow="never" v-loading="loading">
        <template #header>补货建议</template>
        <el-table :data="suggestions" stripe>
          <el-table-column prop="part_id" label="配件编号" />
          <el-table-column prop="code" label="编码" />
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="stock_quantity" label="当前库存" />
          <el-table-column prop="safety_stock" label="安全库存" />
          <el-table-column prop="suggested_quantity" label="建议补货" />
          <el-table-column prop="estimated_cost" label="预计成本" />
          <el-table-column label="操作" width="150">
            <template #default="{ row }">
              <el-button link type="primary" @click="restockSuggested(row)">一键补货</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <el-card shadow="never" v-loading="loading">
        <template #header>低库存列表</template>
        <el-table :data="lowStock" stripe>
          <el-table-column prop="id" label="配件编号" />
          <el-table-column prop="code" label="编码" />
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="stock_quantity" label="当前库存" />
          <el-table-column prop="safety_stock" label="安全库存" />
          <el-table-column prop="status" label="状态" />
        </el-table>
      </el-card>
    </div>

    <el-card shadow="never" v-loading="timelineLoading">
      <template #header>库存时间线</template>
      <el-space>
        <el-input v-model.number="timelinePartID" type="number" placeholder="请输入配件 ID" style="width: 180px" />
        <el-button @click="loadTimeline">加载时间线</el-button>
      </el-space>
      <el-table :data="timeline" stripe style="margin-top: 16px">
        <el-table-column prop="id" label="流水编号" />
        <el-table-column prop="change_type" label="变动类型" />
        <el-table-column prop="quantity_before" label="变动前" />
        <el-table-column prop="change_quantity" label="变动量" />
        <el-table-column prop="quantity_after" label="变动后" />
        <el-table-column prop="reference_type" label="关联类型" />
        <el-table-column prop="operator_id" label="操作人" />
        <el-table-column prop="created_at" label="创建时间" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import client from '../api/client'
import KpiCard from '../components/KpiCard.vue'

const loading = ref(false)
const timelineLoading = ref(false)
const summary = ref<any | null>(null)
const lowStock = ref<any[]>([])
const suggestions = ref<any[]>([])
const timeline = ref<any[]>([])
const timelinePartID = ref<number | null>(null)
const restockForm = reactive({
  partId: 0,
  quantity: 1,
  reason: '手工补货'
})

async function load() {
  loading.value = true
  try {
    const [summaryResponse, lowStockResponse, suggestionResponse] = await Promise.all([
      client.get('/inventory/summary'),
      client.get('/inventory/low-stock'),
      client.get('/inventory/restock-suggestions')
    ])
    summary.value = summaryResponse.data.data
    lowStock.value = lowStockResponse.data.data || []
    suggestions.value = suggestionResponse.data.data || []
  } finally {
    loading.value = false
  }
}

async function submitRestock() {
  await client.post('/inventory/restock', {
    part_id: restockForm.partId,
    quantity: restockForm.quantity,
    reason: restockForm.reason
  })
  ElMessage.success('补货已完成')
  await load()
  timelinePartID.value = restockForm.partId
  await loadTimeline()
}

async function restockSuggested(row: any) {
  restockForm.partId = row.part_id
  restockForm.quantity = row.suggested_quantity
  restockForm.reason = `系统建议补货：${row.code}`
  await submitRestock()
}

async function loadTimeline() {
  if (!timelinePartID.value) {
    ElMessage.warning('请先输入配件 ID')
    return
  }
  timelineLoading.value = true
  try {
    const { data } = await client.get(`/inventory/parts/${timelinePartID.value}/timeline`, {
      params: { limit: 20 }
    })
    timeline.value = data.data || []
  } finally {
    timelineLoading.value = false
  }
}

function formatMoney(value: number | undefined) {
  return `¥${Number(value || 0).toFixed(2)}`
}

load()
</script>
