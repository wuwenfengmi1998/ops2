<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { workOrderApi } from '@/api/work_order'
import {
  IconPlus,
  IconChevronLeftPipe,
  IconChevronRightPipe,
  IconChevronsLeft,
  IconChevronsRight,
} from '@tabler/icons-vue'

usePageTitle('work_order.list_title')
const { t, locale } = useI18n()
const router = useRouter()
const toast = useToastStore()

const orders = ref([])
const totalCount = ref(0)
const pageSize = ref(10)
const currentPage = ref(1)
const statusFilter = ref('')
const loading = ref(false)

const statusOptions = [
  { value: '', labelKey: 'work_order.filter_all' },
  { value: 'pending', labelKey: 'work_order.status_pending' },
  { value: 'checked', labelKey: 'work_order.status_checked' },
  { value: 'parts_ordered', labelKey: 'work_order.status_parts_ordered' },
  { value: 'repaired', labelKey: 'work_order.status_repaired' },
  { value: 'returned', labelKey: 'work_order.status_returned' },
  { value: 'unrepairable', labelKey: 'work_order.status_unrepairable' },
]

const statusColorMap = {
  pending: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/40 dark:text-yellow-400',
  checked: 'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-400',
  parts_ordered: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-400',
  repaired: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400',
  returned: 'bg-gray-200 text-gray-600 dark:bg-gray-700 dark:text-gray-300',
  unrepairable: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400',
}

const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value) || 1)

const pageRange = computed(() => {
  const total = totalPages.value
  const cur = currentPage.value
  let start = Math.max(1, cur - 2)
  let end = Math.min(cur + 4, total)
  if (end - start < 4) start = Math.max(1, end - 4)
  return Array.from({ length: end - start + 1 }, (_, i) => start + i)
})

async function fetchOrders() {
  loading.value = true
  try {
    const { errCode, data } = await workOrderApi.getList({
      status: statusFilter.value,
      entries: pageSize.value,
      page: currentPage.value,
    })
    if (errCode === 0) {
      orders.value = data.all_orders ?? []
      totalCount.value = data.all_count ?? 0
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    loading.value = false
  }
}

function goToPage(page) {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchOrders()
}

function jumpToOrder(id) {
  router.push(`/work_order/show/${id}`)
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Intl.DateTimeFormat(locale.value, {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false,
  }).format(new Date(dateStr))
}

function handlePageSizeInput(e) {
  let val = parseInt(e.target.value) || 10
  if (val > 300) val = 300
  if (val < 1) val = 1
  pageSize.value = val
  currentPage.value = 1
  fetchOrders()
}

onMounted(fetchOrders)
</script>

<template>
  <div class="mx-auto max-w-7xl px-6 py-6">
    <div class="flex flex-col gap-6 rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <!-- Header -->
      <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('work_order.list_title') }}</h3>
      </div>

      <!-- Toolbar -->
      <div class="flex flex-col gap-3 px-6 py-3 sm:flex-row sm:items-center">
        <div class="flex gap-2">
          <RouterLink
            to="/work_order/add"
            class="inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700"
          >
            <IconPlus :size="16" />
            {{ t('work_order.add') }}
          </RouterLink>
          <select
            v-model="statusFilter"
            class="rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
            @change="currentPage = 1; fetchOrders()"
          >
            <option v-for="opt in statusOptions" :key="opt.value" :value="opt.value">
              {{ t(opt.labelKey) }}
            </option>
          </select>
        </div>
      </div>

      <!-- Table -->
      <div class="overflow-x-auto px-0">
        <table class="w-full text-left text-sm text-gray-900">
          <thead>
            <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base">
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400 w-16">No.</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400 w-50">{{ t('work_order.title') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400 w-50">描述</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">关联客户</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap w-44">{{ t('work_order.created_at') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap w-44">{{ t('work_order.updated_at') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400 w-36">{{ t('work_order.status') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="7" class="px-6 py-8 text-center text-gray-400">
                <svg class="mx-auto mb-2 h-5 w-5 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                </svg>
                {{ t('message.loading') }}
              </td>
            </tr>
            <tr v-else-if="orders.length === 0">
              <td colspan="7" class="px-6 py-8 text-center text-gray-400 dark:text-gray-500">
                暂无工单
              </td>
            </tr>
            <tr
              v-else
              v-for="order in orders"
              :key="order.ID"
              class="border-b border-gray-100 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:hover:bg-dk-base cursor-pointer"
              @click="jumpToOrder(order.ID)"
            >
              <td class="px-6 py-3 text-gray-500 dark:text-gray-400">{{ order.ID }}</td>
              <td class="px-6 py-3 font-medium text-gray-900 dark:text-white max-w-[12rem] truncate">{{ order.Title }}</td>
              <td class="px-6 py-3 text-gray-500 dark:text-gray-400 max-w-xs truncate">{{ order.Description || '—' }}</td>
              <td class="px-6 py-3">
                <div v-if="order.customers && order.customers.length > 0" class="flex flex-wrap gap-1">
                  <RouterLink
                    v-for="c in order.customers"
                    :key="c.id"
                    :to="`/customer/detail/${c.id}`"
                    class="inline-flex items-center gap-1 rounded-full border border-blue-200 bg-blue-50 px-2 py-0.5 text-xs text-blue-700 hover:bg-blue-100 dark:border-blue-800 dark:bg-blue-900/30 dark:text-blue-300 whitespace-nowrap"
                    @click.stop
                  >
                    {{ (c.last_name || '') + (c.first_name ? ' ' + c.first_name : '') }}
                  </RouterLink>
                </div>
                <span v-else class="text-gray-400">—</span>
              </td>
              <td class="px-6 py-3 whitespace-nowrap text-gray-500 dark:text-gray-400">{{ formatDate(order.CreatedAt) }}</td>
              <td class="px-6 py-3 whitespace-nowrap text-gray-500 dark:text-gray-400">{{ formatDate(order.UpdatedAt) }}</td>
              <td class="px-6 py-3">
                <span
                  class="inline-block rounded-full px-2.5 py-0.5 text-xs font-medium"
                  :class="statusColorMap[order.CurrentStatus] || 'bg-gray-100 text-gray-600'"
                >
                  {{ t('work_order.status_' + order.CurrentStatus) || order.CurrentStatus }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="flex flex-col items-center gap-3 border-t border-gray-100 px-6 py-4 sm:flex-row sm:justify-between dark:border-dk-muted">
        <div class="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
          <span>共 {{ totalCount }} 条</span>
          <span>每页</span>
          <input
            type="number"
            :value="pageSize"
            min="1"
            max="300"
            class="w-14 rounded border border-gray-300 px-1.5 py-0.5 text-center text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
            @change="handlePageSizeInput"
          />
          <span>条</span>
        </div>
        <div class="flex items-center gap-1">
          <button @click="goToPage(1)" :disabled="currentPage === 1" class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted">
            <IconChevronsLeft :size="16" />
          </button>
          <button @click="goToPage(currentPage - 1)" :disabled="currentPage === 1" class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted">
            <IconChevronLeftPipe :size="16" />
          </button>
          <button
            v-for="p in pageRange"
            :key="p"
            @click="goToPage(p)"
            :class="['rounded px-2.5 py-1 text-sm', p === currentPage ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dk-muted']"
          >{{ p }}</button>
          <button @click="goToPage(currentPage + 1)" :disabled="currentPage === totalPages" class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted">
            <IconChevronRightPipe :size="16" />
          </button>
          <button @click="goToPage(totalPages)" :disabled="currentPage === totalPages" class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted">
            <IconChevronsRight :size="16" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
