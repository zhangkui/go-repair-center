<template>
  <CrudPage
    :title="String(route.meta.title || '数据列表')"
    :subtitle="String(route.meta.subtitle || '维护业务数据')"
    :endpoint="String(route.meta.endpoint || '/')"
    :columns="(route.meta.columns as any[]) || []"
    :fields="(route.meta.fields as any[]) || []"
  />
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import CrudPage from '../components/CrudPage.vue'

const route = useRoute()
</script>
