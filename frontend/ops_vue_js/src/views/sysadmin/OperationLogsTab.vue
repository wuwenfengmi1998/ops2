<script setup>
import { ref, computed, onActivated } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUsersStore } from '@/stores/users'
import { authApi } from '@/api/auth'
import { IconSearch, IconRefresh, IconChevronLeft, IconChevronRight } from '@tabler/icons-vue'

const { t } = useI18n()
const usersStore = useUsersStore()

const logs = ref([])
const loading = ref(false)
const search = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

// 模块列表
const modules = [
  { id: 'all', label: 'operation_logs.all' },
  { id: 'customer', label: 'customer.title' },
  { id: 'purchase', label: 'purchase.title' },
  { id: 'schedule', label: 'schedule.title' },
  { id: 'warehouse', label: 'warehouse.title' },
  { id: 'work_order', label: 'work_order.title' },
]
const activeModule = ref('all')

async function fetchLogs() {
  loading.value = true
  try {
    const res = await authApi.getOperationLogs({
      page: page.value,
      page_size: pageSize.value,
      module: activeModule.value,
      search: search.value,
    })
    if (res.errCode === 0) {
      logs.value = res.data.logs || []
      total.value = res.data.total || 0
      page.value = res.data.page || 1
      pageSize.value = res.data.page_size || 20
      // 预加载用户信息
      logs.value.forEach(log => {
        if (log.user_id > 0) {
          usersStore.fetchUser(log.user_id)
        }
      })
    }
  } catch {
    // 错误已由拦截器处理
  } finally {
    loading.value = false
  }
}

function onSearch() {
  page.value = 1
  fetchLogs()
}

function onPageChange(newPage) {
  page.value = newPage
  fetchLogs()
}

function onModuleChange(moduleId) {
  activeModule.value = moduleId
  page.value = 1
  fetchLogs()
}

function formatActionType(actionType) {
  const key = `operation_logs.action_${actionType}`
  const text = t(key)
  return text === key ? actionType : text
}

function getModuleClass(module) {
  const map = {
    customer: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
    purchase: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
    schedule: 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400',
    warehouse: 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400',
    work_order: 'bg-cyan-100 text-cyan-800 dark:bg-cyan-900/30 dark:text-cyan-400',
  }
  return map[module] || 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400'
}

defineExpose({ fetchLogs })

onActivated(() => fetchLogs())
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('sysadmin.tab_operation_logs') }}</h2>
      <span class="text-sm text-gray-500 dark:text-dk-subtle">
        {{ t('sysadmin.total_logs', { count: total }) }}
      </span>
    </div>

    <div class="flex gap-4">
      <!-- Left: Module Selector -->
      <div class="w-48 flex-shrink-0">
        <nav class="space-y-1">
          <button
            v-for="mod in modules"
            :key="mod.id"
            @click="onModuleChange(mod.id)"
            :class="[
              'w-full rounded-lg px-3 py-2 text-left text-sm font-medium transition-colors',
              activeModule === mod.id
                ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
                : 'text-gray-700 hover:bg-gray-100 dark:text-dk-subtle dark:hover:bg-dk-muted',
            ]"
          >
            {{ t(mod.label) }}
          </button>
        </nav>
      </div>

      <!-- Right: Logs Table -->
      <div class="flex-1 space-y-4">
        <!-- Search Bar -->
        <div class="flex gap-2">
          <div class="relative flex-1">
            <IconSearch class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
            <input
              v-model="search"
              type="text"
              :placeholder="t('operation_logs.search_placeholder')"
              class="w-full rounded-md border border-gray-300 py-2 pl-9 pr-4 text-sm focus:border-blue-500 focus:outline-none dark:border-dk-muted dark:bg-dk-base dark:text-dk-text"
              @keyup.enter="onSearch"
            />
          </div>
          <button
            @click="onSearch"
            class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
          >
            {{ t('common.search') }}
          </button>
          <button
            @click="fetchLogs"
            class="rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-600 hover:bg-gray-50 dark:border-dk-muted dark:text-dk-subtle dark:hover:bg-dk-card"
            :disabled="loading"
          >
            <IconRefresh :size="18" :class="{ 'animate-spin': loading }" />
          </button>
        </div>

        <!-- Logs Table -->
        <div class="overflow-hidden rounded-md border border-gray-200 dark:border-dk-muted">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-dk-muted">
            <thead class="bg-gray-50 dark:bg-dk-base">
              <tr>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('operation_logs.table_module') }}</th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('operation_logs.table_entity_id') }}</th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('operation_logs.table_user') }}</th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('operation_logs.table_action') }}</th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('operation_logs.table_ip') }}</th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('operation_logs.table_remark') }}</th>
                <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('operation_logs.table_created_at') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 bg-white dark:divide-dk-muted dark:bg-dk-card">
              <tr v-if="loading" class="text-center">
                <td colspan="7" class="py-8 text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.loading') }}</td>
              </tr>
              <tr v-else-if="logs.length === 0" class="text-center">
                <td colspan="7" class="py-8 text-gray-500 dark:text-dk-subtle">{{ t('operation_logs.no_logs') }}</td>
              </tr>
              <tr v-for="log in logs" :key="log.id" class="hover:bg-gray-50 dark:hover:bg-dk-base">
                <td class="whitespace-nowrap px-4 py-3">
                  <span :class="['rounded-full px-2 py-0.5 text-xs font-medium', getModuleClass(log.module)]">
                    {{ t(`operation_logs.module_${log.module}`) }}
                  </span>
                </td>
                <td class="whitespace-nowrap px-4 py-3 text-sm text-gray-900 dark:text-dk-text">{{ log.entity_id }}</td>
                <td class="whitespace-nowrap px-4 py-3">
                  <div class="flex items-center gap-2">
                    <img
                      v-if="log.user_id > 0"
                      :src="usersStore.getAvatarUrlFromUserID(log.user_id)"
                      class="h-6 w-6 rounded-full object-cover"
                      alt="avatar"
                    />
                    <span class="text-sm text-gray-900 dark:text-dk-text">
                      {{ usersStore.getUsernameFromUserID(log.user_id) || 'ID: ' + log.user_id }}
                    </span>
                  </div>
                </td>
                <td class="whitespace-nowrap px-4 py-3">
                  <span :class="['rounded-full px-2 py-0.5 text-xs font-medium', log.action_type === 'delete' ? 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400' : 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-400']">
                    {{ formatActionType(log.action_type) }}
                  </span>
                </td>
                <td class="whitespace-nowrap px-4 py-3 text-sm font-mono text-gray-500 dark:text-dk-subtle">{{ log.ip }}</td>
                <td class="px-4 py-3 text-sm text-gray-500 dark:text-dk-subtle max-w-xs truncate">{{ log.remark || '-' }}</td>
                <td class="whitespace-nowrap px-4 py-3 text-sm text-gray-500 dark:text-dk-subtle">{{ log.created_at ? new Date(log.created_at).toLocaleString() : '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-500 dark:text-dk-subtle">
            {{ t('sysadmin.pagination', { current: page, total: totalPages }) }}
          </div>
          <div class="flex gap-2">
            <button
              @click="onPageChange(page - 1)"
              :disabled="page <= 1 || loading"
              class="flex items-center gap-1 rounded-md border border-gray-300 px-3 py-1.5 text-sm disabled:opacity-50 dark:border-dk-muted dark:text-dk-text"
            >
              <IconChevronLeft :size="16" /> {{ t('sysadmin.prev_page') }}
            </button>
            <button
              @click="onPageChange(page + 1)"
              :disabled="page >= totalPages || loading"
              class="flex items-center gap-1 rounded-md border border-gray-300 px-3 py-1.5 text-sm disabled:opacity-50 dark:border-dk-muted dark:text-dk-text"
            >
              {{ t('sysadmin.next_page') }} <IconChevronRight :size="16" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
