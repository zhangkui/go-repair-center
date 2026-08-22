<template>
  <div class="page-shell">
    <div class="page-header">
      <div>
        <h2 class="page-title">数据工具</h2>
        <p class="page-subtitle">预览 CSV 导入、应用校验通过的数据，并导出业务数据集。</p>
      </div>
    </div>

    <el-card shadow="never">
      <el-space wrap>
        <el-select v-model="resource" style="width: 220px">
          <el-option label="客户" value="customers" />
          <el-option label="配件" value="parts" />
          <el-option label="设备" value="devices" />
        </el-select>
        <el-button @click="downloadTemplate">下载模板</el-button>
        <el-button @click="downloadExport('repair-orders')">导出维修单</el-button>
        <el-button @click="downloadExport('parts')">导出配件</el-button>
        <el-button @click="downloadExport('feedbacks')">导出回访</el-button>
        <el-button @click="downloadExport('audit-logs')">导出审计日志</el-button>
      </el-space>
    </el-card>

    <el-card shadow="never">
      <template #header>CSV 导入</template>
      <el-input
        v-model="csvContent"
        type="textarea"
        :rows="14"
        placeholder="请在这里粘贴 CSV 内容"
      />
      <el-space style="margin-top: 16px">
        <el-button type="primary" @click="previewImport">预览</el-button>
        <el-button type="success" @click="applyImport">导入</el-button>
      </el-space>
    </el-card>

    <el-card shadow="never" v-if="preview">
      <template #header>导入预览</template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="表头">{{ (preview.summary?.headers || []).join(', ') }}</el-descriptions-item>
        <el-descriptions-item label="行数">{{ preview.summary?.row_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="有效行">{{ preview.summary?.valid_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="错误数">{{ preview.summary?.error_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="已导入">{{ preview.summary?.inserted_count || 0 }}</el-descriptions-item>
      </el-descriptions>
      <el-table :data="preview.result?.errors || []" style="margin-top: 16px" stripe>
        <el-table-column prop="line_number" label="行号" width="100" />
        <el-table-column prop="message" label="错误信息" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import client from '../api/client'

const resource = ref<'customers' | 'parts' | 'devices'>('customers')
const csvContent = ref('name,phone,address,level,remark\n示例客户,13800000000,示例路 100 号,NORMAL,导入示例')
const preview = ref<any>(null)

async function previewImport() {
  const { data } = await client.post(`/tools/import/${resource.value}/preview`, {
    csv_content: csvContent.value
  })
  preview.value = data.data
}

async function applyImport() {
  const { data } = await client.post(`/tools/import/${resource.value}/apply`, {
    csv_content: csvContent.value
  })
  preview.value = data.data
  ElMessage.success(`成功导入 ${data.data.summary.inserted_count || 0} 行`)
}

function downloadTemplate() {
  window.open(`/api/v1/tools/import/templates/${resource.value}`, '_blank')
}

function downloadExport(resourceName: string) {
  window.open(`/api/v1/tools/export/${resourceName}.csv`, '_blank')
}
</script>
