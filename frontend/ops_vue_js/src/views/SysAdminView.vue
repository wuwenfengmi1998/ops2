<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { authApi } from '@/api/auth'

const userStore = useUserStore()
const activeTab = ref('users')
const sysAdmins = ref([])
const loading = ref(false)

const tabs = [
  { id: 'users', label: '用户管理' },
  { id: 'groups', label: '用户组' },
  { id: 'logs', label: '登录日志' },
  { id: 'config', label: '系统配置' },
]

async function fetchSysAdmins() {
  loading.value = true
  try {
    const res = await authApi.sysAdmins()
    if (res.errCode === 0 && Array.isArray(res.data.sysAdmins)) {
      sysAdmins.value = res.data.sysAdmins
    }
  } catch {
    // 错误已由拦截器处理
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchSysAdmins()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 p-6 dark:bg-dk-base">
    <div class="mx-auto max-w-6xl">
      <!-- Header -->
      <div class="mb-6 flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-dk-text">系统管理</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-dk-subtle">
            系统管理员专用页面
          </p>
        </div>
        <div class="flex items-center gap-2 rounded-lg bg-amber-100 px-3 py-1.5 dark:bg-amber-900/30">
          <span class="text-amber-700 dark:text-amber-400">管理员: {{ userStore.user?.Username }}</span>
        </div>
      </div>

      <!-- Tabs -->
      <div class="mb-6 border-b border-gray-200 dark:border-dk-muted">
        <nav class="-mb-px flex gap-6">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'border-b-2 px-1 pb-3 text-sm font-medium transition-colors',
              activeTab === tab.id
                ? 'border-blue-500 text-blue-600 dark:border-blue-400 dark:text-blue-400'
                : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 dark:text-dk-subtle dark:hover:text-dk-text',
            ]"
          >
            {{ tab.label }}
          </button>
        </nav>
      </div>

      <!-- Content -->
      <div class="rounded-lg border border-gray-200 bg-white p-6 dark:border-dk-muted dark:bg-dk-card">
        <!-- 用户管理 -->
        <div v-if="activeTab === 'users'" class="space-y-4">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-dk-text">用户管理</h2>
          <div class="rounded-md bg-gray-50 p-8 text-center dark:bg-dk-base">
            <p class="text-gray-500 dark:text-dk-subtle">sysadmin_users_placeholder</p>
          </div>
        </div>

        <!-- 用户组 -->
        <div v-if="activeTab === 'groups'" class="space-y-4">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-dk-text">用户组管理</h2>
          <div class="rounded-md bg-gray-50 p-8 text-center dark:bg-dk-base">
            <p class="text-gray-500 dark:text-dk-subtle">sysadmin_groups_placeholder</p>
          </div>
        </div>

        <!-- 登录日志 -->
        <div v-if="activeTab === 'logs'" class="space-y-4">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-dk-text">登录失败日志</h2>
          <div class="rounded-md bg-gray-50 p-8 text-center dark:bg-dk-base">
            <p class="text-gray-500 dark:text-dk-subtle">sysadmin_logs_placeholder</p>
          </div>
        </div>

        <!-- 系统配置 -->
        <div v-if="activeTab === 'config'" class="space-y-4">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-dk-text">系统配置</h2>
          <div class="rounded-md bg-gray-50 p-8 text-center dark:bg-dk-base">
            <p class="text-gray-500 dark:text-dk-subtle">sysadmin_config_placeholder</p>
          </div>
        </div>
      </div>

      <!-- SysAdmins List -->
      <div class="mt-6 rounded-lg border border-gray-200 bg-white p-4 dark:border-dk-muted dark:bg-dk-card">
        <div class="mb-2 flex items-center justify-between">
          <h3 class="text-sm font-medium text-gray-700 dark:text-dk-subtle">当前系统管理员列表</h3>
          <button
            @click="fetchSysAdmins"
            class="text-xs text-blue-600 hover:text-blue-700 dark:text-blue-400"
            :disabled="loading"
          >
            {{ loading ? '加载中...' : '刷新' }}
          </button>
        </div>
        <div class="flex flex-wrap gap-2">
          <span
            v-for="adminId in sysAdmins"
            :key="adminId"
            class="rounded-full bg-amber-100 px-2.5 py-0.5 text-xs font-medium text-amber-800 dark:bg-amber-900/30 dark:text-amber-400"
          >
            ID: {{ adminId }}
          </span>
          <span v-if="sysAdmins.length === 0" class="text-sm text-gray-400 dark:text-dk-muted">
            暂无系统管理员
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
