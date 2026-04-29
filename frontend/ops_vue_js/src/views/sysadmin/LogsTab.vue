<script setup>
import { ref, computed, onActivated } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUsersStore } from '@/stores/users'
import { authApi } from '@/api/auth'
import { IconSearch, IconRefresh, IconChevronLeft, IconChevronRight } from '@tabler/icons-vue'

const { t } = useI18n()
const usersStore = useUsersStore()

const loginFailLogs = ref([])
const loginFailLogsLoading = ref(false)
const loginFailLogSearch = ref('')
const loginFailLogPage = ref(1)
const loginFailLogPageSize = ref(20)
const loginFailLogTotal = ref(0)
const loginFailLogTotalPages = computed(() => Math.ceil(loginFailLogTotal.value / loginFailLogPageSize.value))

async function fetchLoginFailLogs() {
  loginFailLogsLoading.value = true
  try {
    const res = await authApi.getLoginFailLogs({
      page: loginFailLogPage.value,
      page_size: loginFailLogPageSize.value,
      search: loginFailLogSearch.value,
    })
    if (res.errCode === 0) {
      loginFailLogs.value = res.data.logs || []
      loginFailLogTotal.value = res.data.total || 0
      loginFailLogPage.value = res.data.page || 1
      loginFailLogPageSize.value = res.data.page_size || 20
      loginFailLogs.value.forEach(log => {
        if (log.user_id > 0) {
          usersStore.fetchUser(log.user_id)
        }
      })
    }
  } catch {
    // 错误已由拦截器处理
  } finally {
    loginFailLogsLoading.value = false
  }
}

function onLoginFailLogSearch() {
  loginFailLogPage.value = 1
  fetchLoginFailLogs()
}

function onLoginFailLogPageChange(page) {
  loginFailLogPage.value = page
  fetchLoginFailLogs()
}

function formatReason(reason) {
  const reasonMap = {
    'password_error': t('sysadmin.reason_password_error'),
    'user_not_found': t('sysadmin.reason_user_not_found'),
  }
  return reasonMap[reason] || reason
}

function getReasonClass(reason) {
  if (reason === 'password_error') {
    return 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400'
  }
  if (reason === 'user_not_found') {
    return 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
  }
  return 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
}

defineExpose({ fetchLoginFailLogs })

onActivated(() => fetchLoginFailLogs())
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('sysadmin.tab_logs') }}</h2>
      <span class="text-sm text-gray-500 dark:text-dk-subtle">
        {{ t('sysadmin.total_logs', { count: loginFailLogTotal }) }}
      </span>
    </div>

    <!-- 搜索栏 -->
    <div class="flex gap-2">
      <div class="relative flex-1">
        <IconSearch class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
        <input
          v-model="loginFailLogSearch"
          type="text"
          :placeholder="t('sysadmin.search_log_placeholder')"
          class="w-full rounded-md border border-gray-300 py-2 pl-9 pr-4 text-sm focus:border-blue-500 focus:outline-none dark:border-dk-muted dark:bg-dk-base dark:text-dk-text"
          @keyup.enter="onLoginFailLogSearch"
        />
      </div>
      <button
        @click="onLoginFailLogSearch"
        class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
      >
        {{ t('sysadmin.search') }}
      </button>
      <button
        @click="fetchLoginFailLogs"
        class="rounded-md border border-gray-300 px-3 py-2 text-gray-600 hover:bg-gray-50 dark:border-dk-muted dark:text-dk-subtle dark:hover:bg-dk-card"
        :disabled="loginFailLogsLoading"
      >
        <IconRefresh :size="18" :class="{ 'animate-spin': loginFailLogsLoading }" />
      </button>
    </div>

    <!-- 日志列表 -->
    <div class="overflow-hidden rounded-md border border-gray-200 dark:border-dk-muted">
      <table class="min-w-full divide-y divide-gray-200 dark:divide-dk-muted">
        <thead class="bg-gray-50 dark:bg-dk-base">
          <tr>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_user') }}</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_reason') }}</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_count') }}</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_ip') }}</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_updated_at') }}</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_created') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 bg-white dark:divide-dk-muted dark:bg-dk-card">
          <tr v-if="loginFailLogsLoading" class="text-center">
            <td colspan="6" class="py-8 text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.loading') }}</td>
          </tr>
          <tr v-else-if="loginFailLogs.length === 0" class="text-center">
            <td colspan="6" class="py-8 text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.no_logs') }}</td>
          </tr>
          <tr v-for="log in loginFailLogs" :key="log.id" class="hover:bg-gray-50 dark:hover:bg-dk-base">
            <td class="whitespace-nowrap px-4 py-3">
              <div class="flex items-center gap-2">
                <img
                  v-if="log.user_id > 0"
                  :src="usersStore.getAvatarUrlFromUserID(log.user_id)"
                  class="h-7 w-7 rounded-full object-cover"
                  alt="avatar"
                />
                <div v-else class="h-7 w-7 rounded-full bg-gray-200 dark:bg-gray-700"></div>
                <span class="text-sm text-gray-900 dark:text-dk-text">
                  {{ log.username }}
                </span>
              </div>
            </td>
            <td class="whitespace-nowrap px-4 py-3">
              <span
                :class="[
                  'rounded-full px-2 py-0.5 text-xs font-medium',
                  getReasonClass(log.reason)
                ]"
              >
                {{ formatReason(log.reason) }}
              </span>
            </td>
            <td class="whitespace-nowrap px-4 py-3 text-sm">
              <span :class="[
                'font-medium',
                log.count >= 5 ? 'text-red-600 dark:text-red-400' : 'text-gray-900 dark:text-dk-text'
              ]">
                {{ log.count }}
              </span>
            </td>
            <td class="whitespace-nowrap px-4 py-3 text-sm font-mono text-gray-500 dark:text-dk-subtle">{{ log.ip }}</td>
            <td class="whitespace-nowrap px-4 py-3 text-sm text-gray-500 dark:text-dk-subtle">{{ new Date(log.updated_at).toLocaleString() }}</td>
            <td class="whitespace-nowrap px-4 py-3 text-sm text-gray-500 dark:text-dk-subtle">{{ new Date(log.created_at).toLocaleString() }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 分页 -->
    <div class="flex items-center justify-between">
      <div class="text-sm text-gray-500 dark:text-dk-subtle">
        {{ t('sysadmin.pagination', { current: loginFailLogPage, total: loginFailLogTotalPages }) }}
      </div>
      <div class="flex gap-2">
        <button
          @click="onLoginFailLogPageChange(loginFailLogPage - 1)"
          :disabled="loginFailLogPage <= 1 || loginFailLogsLoading"
          class="flex items-center gap-1 rounded-md border border-gray-300 px-3 py-1.5 text-sm disabled:opacity-50 dark:border-dk-muted dark:text-dk-text"
        >
          <IconChevronLeft :size="16" /> {{ t('sysadmin.prev_page') }}
        </button>
        <button
          @click="onLoginFailLogPageChange(loginFailLogPage + 1)"
          :disabled="loginFailLogPage >= loginFailLogTotalPages || loginFailLogsLoading"
          class="flex items-center gap-1 rounded-md border border-gray-300 px-3 py-1.5 text-sm disabled:opacity-50 dark:border-dk-muted dark:text-dk-text"
        >
          {{ t('sysadmin.next_page') }} <IconChevronRight :size="16" />
        </button>
      </div>
    </div>
  </div>
</template>
