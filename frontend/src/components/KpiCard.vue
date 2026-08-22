<template>
  <div class="kpi-card">
    <span class="label">{{ label }}</span>
    <strong class="value">{{ value }}</strong>
    <small class="caption">{{ caption }}</small>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label: string
  value: string | number
  caption?: string
}>()
</script>

<style scoped>
.kpi-card {
  padding: 20px;
  border-radius: 20px;
  background: linear-gradient(180deg, #ffffff, #f6fbfb);
  box-shadow: var(--shadow);
  border: 1px solid rgba(15, 118, 110, 0.1);
}

.label {
  color: var(--muted);
  font-size: 13px;
}

.value {
  display: block;
  margin-top: 10px;
  font-size: 32px;
}

.caption {
  display: block;
  margin-top: 8px;
  color: var(--muted);
}
</style>

