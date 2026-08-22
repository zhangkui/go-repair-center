<template>
  <div class="page-shell">
    <div class="page-header">
      <div>
        <h2 class="page-title">审批中心</h2>
        <p class="page-subtitle">在统一队列中处理报价审批和质保审批。</p>
      </div>
      <el-button @click="load">刷新</el-button>
    </div>

    <div class="card-grid">
      <KpiCard label="待审报价" :value="overview?.pending_quotation_count ?? 0" caption="等待复核人员处理的报价单" />
      <KpiCard label="待审质保" :value="overview?.pending_warranty_count ?? 0" caption="尚未确认的质保记录" />
      <KpiCard label="待处理返修" :value="overview?.pending_rework_count ?? 0" caption="已生成返修但尚未关闭的回访单" />
    </div>

    <div class="two-col">
      <el-card shadow="never" v-loading="loading">
        <template #header>报价审批队列</template>
        <el-table :data="quotations" stripe>
          <el-table-column prop="id" label="编号" />
          <el-table-column prop="repair_order_id" label="维修单" />
          <el-table-column prop="version" label="版本" />
          <el-table-column prop="total_amount" label="金额" />
          <el-table-column prop="approval_status" label="审批状态" />
          <el-table-column prop="valid_until" label="有效期至" />
          <el-table-column label="操作" width="170">
            <template #default="{ row }">
              <el-space>
                <el-button link type="success" @click="decide('quotations', row, true)">通过</el-button>
                <el-button link type="danger" @click="decide('quotations', row, false)">驳回</el-button>
              </el-space>
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <el-card shadow="never" v-loading="loading">
        <template #header>质保审批队列</template>
        <el-table :data="warranties" stripe>
          <el-table-column prop="id" label="编号" />
          <el-table-column prop="repair_order_id" label="维修单" />
          <el-table-column prop="warranty_type" label="质保类型" />
          <el-table-column prop="duration_days" label="时长（天）" />
          <el-table-column prop="status" label="状态" />
          <el-table-column label="操作" width="170">
            <template #default="{ row }">
              <el-space>
                <el-button link type="success" @click="decide('warranties', row, true)">通过</el-button>
                <el-button link type="danger" @click="decide('warranties', row, false)">驳回</el-button>
              </el-space>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import client from '../api/client'
import KpiCard from '../components/KpiCard.vue'

const loading = ref(false)
const overview = ref<any | null>(null)
const quotations = ref<any[]>([])
const warranties = ref<any[]>([])

async function load() {
  loading.value = true
  try {
    const [overviewResponse, quotationResponse, warrantyResponse] = await Promise.all([
      client.get('/approvals/overview'),
      client.get('/approvals/quotations', { params: { page: 1, page_size: 20 } }),
      client.get('/approvals/warranties', { params: { page: 1, page_size: 20 } })
    ])
    overview.value = overviewResponse.data.data
    quotations.value = quotationResponse.data.data?.items || []
    warranties.value = warrantyResponse.data.data?.items || []
  } finally {
    loading.value = false
  }
}

async function decide(resource: 'quotations' | 'warranties', row: any, approved: boolean) {
  let comment = ''
  const resourceText = resource === 'quotations' ? '报价单' : '质保单'
  try {
    const promptResult = await ElMessageBox.prompt(
      `${approved ? '确认通过' : '确认驳回'} ${resourceText} #${row.id}，审批意见可选填写。`,
      '审批处理',
      {
        inputPlaceholder: '请输入审批意见',
        confirmButtonText: approved ? '通过' : '驳回',
        cancelButtonText: '取消'
      }
    )
    comment = promptResult.value || ''
  } catch {
    return
  }
  await client.post(`/approvals/${resource}/${row.id}/decision`, {
    approved,
    comment
  })
  ElMessage.success(`${resourceText}已更新`)
  await load()
}

load()
</script>
