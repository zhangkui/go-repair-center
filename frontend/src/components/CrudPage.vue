<template>
  <div class="page-shell">
    <div class="page-header">
      <div>
        <h2 class="page-title">{{ title }}</h2>
        <p class="page-subtitle">{{ subtitle }}</p>
      </div>
      <el-space>
        <el-button @click="load">刷新</el-button>
        <el-button type="primary" @click="openCreate">新建</el-button>
      </el-space>
    </div>

    <el-card shadow="never">
      <el-space wrap>
        <el-input v-model="keyword" clearable placeholder="关键字搜索" style="width: 220px" />
        <el-input v-model="status" clearable placeholder="状态筛选" style="width: 180px" />
        <el-button type="primary" @click="load">查询</el-button>
      </el-space>
    </el-card>

    <el-card shadow="never">
      <el-table :data="items" v-loading="loading" stripe>
        <el-table-column
          v-for="column in columns"
          :key="column.key"
          :prop="column.key"
          :label="column.label"
          min-width="140"
        >
          <template #default="{ row }">
            {{ formatCell(column.key, row[column.key]) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="180">
          <template #default="{ row }">
            <el-space>
              <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
              <el-button link type="danger" @click="removeRow(row)">删除</el-button>
            </el-space>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editingId ? `编辑${title}` : `新建${title}`" width="680px">
      <el-form label-width="140px">
        <el-form-item v-for="field in fields" :key="field.key" :label="field.label">
          <el-input v-if="field.type !== 'number'" v-model="form[field.key]" />
          <el-input-number v-else v-model="form[field.key]" :min="0" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-space>
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="saving" @click="save">保存</el-button>
        </el-space>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import client from '../api/client'

export interface CrudColumn {
  key: string
  label: string
}

export interface CrudField extends CrudColumn {
  type?: 'text' | 'number'
}

const props = defineProps<{
  title: string
  subtitle: string
  endpoint: string
  columns: CrudColumn[]
  fields: CrudField[]
}>()

const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const keyword = ref('')
const status = ref('')
const editingId = ref<number | null>(null)
const items = ref<Record<string, unknown>[]>([])
const form = reactive<Record<string, any>>({})

const dictionaries: Record<string, Record<string, string>> = {
  status: {
    ACTIVE: '启用',
    INACTIVE: '停用',
    DISABLED: '停用',
    PENDING: '待处理',
    DISPATCHED: '已派工',
    WAITING_VISIT: '待上门',
    WAITING_DROPOFF: '待送修',
    ACCEPTED: '已接修',
    IN_PROGRESS: '处理中',
    WAITING_PICKUP: '待取机',
    COMPLETED: '已完成',
    DELIVERED: '已交付',
    CANCELLED: '已取消',
    DRAFT: '草稿',
    APPROVED: '已通过',
    REJECTED: '已驳回',
    EXPIRED: '已过期',
    VOID: '已作废',
    CLOSED: '已关闭',
    QUEUED: '已排队',
    SKIPPED: '已跳过',
    PAID: '已支付',
    PARTIAL: '部分支付'
  },
  level: {
    NORMAL: '普通',
    VIP: 'VIP',
    ENTERPRISE: '企业'
  },
  service_method: {
    ONSITE: '上门服务',
    DROPOFF: '送修服务',
    REMOTE: '远程支持'
  },
  urgency: {
    NORMAL: '一般',
    URGENT: '紧急',
    CRITICAL: '非常紧急'
  },
  method: {
    PHONE: '电话回访',
    ONLINE: '在线回访',
    ONSITE: '现场回访'
  },
  complaint_type: {
    SERVICE: '服务态度',
    QUALITY: '维修质量',
    PRICE: '收费问题',
    TIMELINESS: '时效问题',
    OTHER: '其他'
  },
  warranty_type: {
    FULL_MACHINE: '整机质保',
    COMPONENT: '部件质保'
  },
  approval_status: {
    NOT_REQUIRED: '无需审核',
    PENDING: '待审核',
    APPROVED: '已通过',
    REJECTED: '已驳回'
  },
  notification_status: {
    QUEUED: '待通知',
    SENT: '已通知',
    SKIPPED: '未通知'
  },
  value_type: {
    string: '字符串',
    number: '数值',
    boolean: '布尔值'
  }
}

watch(
  () => props.fields,
  () => {
    resetForm()
  },
  { immediate: true }
)

async function load() {
  loading.value = true
  try {
    const { data } = await client.get(props.endpoint, {
      params: { page: 1, page_size: 20, keyword: keyword.value, status: status.value }
    })
    items.value = data.data.items || []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  resetForm()
  dialogVisible.value = true
}

function openEdit(row: Record<string, any>) {
  editingId.value = Number(row.id)
  resetForm()
  props.fields.forEach((field) => {
    form[field.key] = row[field.key]
  })
  dialogVisible.value = true
}

async function save() {
  saving.value = true
  try {
    if (editingId.value) {
      await client.put(`${props.endpoint}/${editingId.value}`, form)
    } else {
      await client.post(props.endpoint, form)
    }
    dialogVisible.value = false
    await load()
    ElMessage.success('保存成功')
  } finally {
    saving.value = false
  }
}

async function removeRow(row: Record<string, any>) {
  await ElMessageBox.confirm(`确认删除记录 #${row.id} 吗？`, '提示', { type: 'warning' })
  await client.delete(`${props.endpoint}/${row.id}`)
  ElMessage.success('删除成功')
  await load()
}

function resetForm() {
  Object.keys(form).forEach((key) => delete form[key])
  props.fields.forEach((field) => {
    form[field.key] = field.type === 'number' ? 0 : ''
  })
}

function formatCell(key: string, value: unknown) {
  if (value === null || value === undefined || value === '') {
    return '-'
  }
  const dictionary = dictionaries[key]
  if (dictionary) {
    const mapped = dictionary[String(value)]
    if (mapped) {
      return mapped
    }
  }
  return value
}

load()
</script>
