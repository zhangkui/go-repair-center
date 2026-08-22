<template>
  <div class="page-shell">
    <div class="page-header">
      <div>
        <h2 class="page-title">派工中心</h2>
        <p class="page-subtitle">批量派工、校验技师时间冲突，并查看返修链路。</p>
      </div>
      <el-button @click="resetResult">清空结果</el-button>
    </div>

    <div class="card-grid">
      <KpiCard label="当前冲突数" :value="conflicts.length" caption="所选技师时间窗内的排期冲突" />
      <KpiCard label="已派工" :value="dispatchResult?.dispatched?.length ?? 0" caption="本次批量派工成功的工单数" />
      <KpiCard label="已跳过" :value="dispatchResult?.skipped?.length ?? 0" caption="因状态不允许而未派工的工单数" />
      <KpiCard label="待通知" :value="dispatchResult?.notifications?.length ?? 0" caption="派工后已加入通知队列的消息数" />
    </div>

    <el-card shadow="never">
      <template #header>批量派工</template>
      <el-form label-width="140px">
        <el-form-item label="工单编号列表">
          <el-input
            v-model="form.orderIdsText"
            type="textarea"
            :rows="4"
            placeholder="101,102,103"
          />
        </el-form-item>
        <el-form-item label="技师编号">
          <el-input v-model.number="form.technicianId" type="number" />
        </el-form-item>
        <el-form-item label="技师姓名">
          <el-input v-model="form.technicianName" />
        </el-form-item>
        <el-form-item label="预约开始时间">
          <el-input v-model="form.appointmentTime" type="datetime-local" />
        </el-form-item>
        <el-form-item label="预约结束时间">
          <el-input v-model="form.appointmentEnd" type="datetime-local" />
        </el-form-item>
        <el-form-item label="服务方式">
          <el-select v-model="form.serviceMethod" style="width: 240px">
            <el-option label="上门服务" value="ONSITE" />
            <el-option label="送修服务" value="DROPOFF" />
            <el-option label="远程支持" value="REMOTE" />
          </el-select>
        </el-form-item>
        <el-form-item label="通知技师">
          <el-switch v-model="form.notifyTechnician" />
        </el-form-item>
      </el-form>
      <el-space>
        <el-button @click="checkConflicts">检查冲突</el-button>
        <el-button type="primary" @click="runBatchDispatch">执行派工</el-button>
      </el-space>
    </el-card>

    <el-card shadow="never">
      <template #header>冲突结果</template>
      <el-table :data="conflicts" stripe>
        <el-table-column prop="id" label="工单编号" />
        <el-table-column prop="order_number" label="维修单号" />
        <el-table-column prop="appointment_time" label="开始时间" />
        <el-table-column prop="appointment_end" label="结束时间" />
        <el-table-column prop="status" label="状态" />
      </el-table>
    </el-card>

    <el-card shadow="never" v-if="dispatchResult">
      <template #header>派工结果</template>
      <div class="two-col">
        <el-card shadow="never">
          <template #header>派工成功</template>
          <el-table :data="dispatchResult.dispatched || []" stripe>
            <el-table-column prop="id" label="工单编号" />
            <el-table-column prop="order_number" label="维修单号" />
            <el-table-column prop="status" label="状态" />
            <el-table-column prop="notification_status" label="通知状态" />
          </el-table>
        </el-card>
        <el-card shadow="never">
          <template #header>跳过列表</template>
          <el-table :data="dispatchResult.skipped || []" stripe>
            <el-table-column prop="id" label="工单编号" />
            <el-table-column prop="status" label="状态" />
            <el-table-column prop="reason" label="原因" />
          </el-table>
        </el-card>
      </div>
    </el-card>

    <el-card shadow="never">
      <template #header>返修链路</template>
      <el-space>
        <el-input v-model.number="reworkOrderId" type="number" placeholder="请输入维修单 ID" style="width: 220px" />
        <el-button @click="loadReworkChain">加载链路</el-button>
      </el-space>
      <el-table :data="reworkChain" stripe style="margin-top: 16px">
        <el-table-column prop="id" label="工单编号" />
        <el-table-column prop="order_number" label="维修单号" />
        <el-table-column prop="original_order_id" label="原工单" />
        <el-table-column prop="status" label="状态" />
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

const form = reactive({
  orderIdsText: '101,102',
  technicianId: 1,
  technicianName: '张师傅',
  appointmentTime: '',
  appointmentEnd: '',
  serviceMethod: 'ONSITE',
  notifyTechnician: true
})

const conflicts = ref<any[]>([])
const dispatchResult = ref<any | null>(null)
const reworkOrderId = ref<number | null>(null)
const reworkChain = ref<any[]>([])

function parseOrderIDs() {
  return form.orderIdsText
    .split(/[\s,]+/)
    .map((item) => Number(item.trim()))
    .filter((item) => Number.isFinite(item) && item > 0)
}

function toISO(value: string) {
  if (!value) {
    return new Date().toISOString()
  }
  const normalized = value.length === 16 ? `${value}:00` : value
  return new Date(normalized).toISOString()
}

async function checkConflicts() {
  const { data } = await client.post('/dispatch/conflicts', {
    technician_id: form.technicianId,
    appointment_time: toISO(form.appointmentTime),
    appointment_end: toISO(form.appointmentEnd)
  })
  conflicts.value = data.data || []
}

async function runBatchDispatch() {
  const { data } = await client.post('/dispatch/batch', {
    order_ids: parseOrderIDs(),
    technician_id: form.technicianId,
    technician_name: form.technicianName,
    appointment_time: toISO(form.appointmentTime),
    appointment_end: toISO(form.appointmentEnd),
    service_method: form.serviceMethod,
    notify_technician: form.notifyTechnician
  })
  dispatchResult.value = data.data
  conflicts.value = data.data?.conflicts || []
  ElMessage.success(`成功派工 ${data.data?.dispatched?.length || 0} 单`)
}

async function loadReworkChain() {
  if (!reworkOrderId.value) {
    ElMessage.warning('请先输入维修单 ID')
    return
  }
  const { data } = await client.get(`/dispatch/rework-chain/${reworkOrderId.value}`)
  reworkChain.value = data.data || []
}

function resetResult() {
  conflicts.value = []
  dispatchResult.value = null
  reworkChain.value = []
}
</script>
