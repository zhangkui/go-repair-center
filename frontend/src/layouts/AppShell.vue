<template>
  <div class="shell">
    <aside class="sidebar">
      <div class="brand">
        <div class="brand-tag">维修中心</div>
        <h1>go-repair-center</h1>
        <p>家电维修管理控制台</p>
      </div>
      <nav class="nav">
        <RouterLink
          v-for="item in menuItems"
          :key="item.path"
          :to="item.path"
          class="nav-item"
        >
          <span>{{ item.meta?.title }}</span>
        </RouterLink>
      </nav>
    </aside>
    <main class="content">
      <header class="content-header">
        <div>
          <strong>{{ auth.user?.display_name || auth.user?.username || '用户' }}</strong>
          <div class="hint">默认账号仅供本地开发验收使用，生产环境请立即修改。</div>
        </div>
        <el-space>
          <el-button @click="$router.push('/profile')">个人资料</el-button>
          <el-button type="primary" plain @click="handleLogout">退出登录</el-button>
        </el-space>
      </header>
      <section class="content-body">
        <RouterView />
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute, RouterLink, RouterView } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const menuItems = computed(() =>
  router.getRoutes().filter((item) => item.meta?.menu && item.path !== route.path)
)

async function handleLogout() {
  await auth.logout()
  await router.push('/login')
}
</script>

<style scoped>
.shell {
  display: grid;
  grid-template-columns: 280px 1fr;
  min-height: 100vh;
}

.sidebar {
  padding: 28px;
  background: linear-gradient(180deg, #0f172a 0%, #17395e 100%);
  color: #eef7fb;
}

.brand h1 {
  margin: 10px 0 6px;
  font-size: 28px;
}

.brand p {
  margin: 0;
  color: rgba(238, 247, 251, 0.72);
}

.brand-tag {
  display: inline-flex;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.16);
  font-size: 12px;
  letter-spacing: 0.14em;
}

.nav {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 28px;
}

.nav-item {
  padding: 12px 14px;
  border-radius: 14px;
  color: rgba(238, 247, 251, 0.86);
  transition: 0.2s ease;
}

.nav-item:hover,
.router-link-active {
  background: rgba(255, 255, 255, 0.12);
  color: #ffffff;
}

.content {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 28px;
  background: rgba(255, 255, 255, 0.82);
  backdrop-filter: blur(18px);
  border-bottom: 1px solid var(--border);
}

.hint {
  margin-top: 6px;
  color: var(--danger);
  font-size: 12px;
}

.content-body {
  padding: 24px 28px 32px;
}

@media (max-width: 960px) {
  .shell {
    grid-template-columns: 1fr;
  }
  .sidebar {
    padding: 20px;
  }
}
</style>
