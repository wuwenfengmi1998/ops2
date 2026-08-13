<script setup>
import { ref, computed, onActivated } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useToastStore } from '@/stores/toast'
import { purchaseApi } from '@/api/purchase'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import {
  IconLock,
  IconLockOpen,
  IconRefresh,
  IconChevronLeftPipe,
  IconChevronRightPipe,
  IconChevronsLeft,
  IconChevronsRight,
} from '@tabler/icons-vue'

const { t, locale } = useI18n()
const router = useRouter()
const toast = useToastStore()

const orders = ref([])
const totalCount = ref(0)
const pageSize = ref(20)
const currentPage = ref(1)
const statusFilter = ref('')
const loading = ref(false)

const showLockAllConfirm = ref(false)
const lockAllLoading = ref(false)

const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value) || 1)

const pageRange = computed(() => {
  const total = totalPages.value
  const cur = currentPage.value
  let start = Math.max(1, cur - 2)
  let end = Math.min(cur + 4, total)
  if (end - start < 4) start = Math.max(1, end - 4)
  return Array.from({ length: end - start + 1 }, (_, i) => start + i)
})

const statusOptions = [
  { value: '', labelKey: 'purchase.filter_all' },
  { value: 'pending', labelKey: 'purchase.status_pending' },
  { value: 'ordered', labelKey: 'purchase.status_ordered' },
  { value: 'arrived', labelKey: 'purchase.status_arrived' },
  { value: 'received', labelKey: 'purchase.status_received' },
  { value: 'lost', labelKey: 'purchase.status_lost' },
  { value: 'returned', labelKey: 'purchase.status_returned' },
]

async function fetchOrders() {
  loading.value = true
  try {
    const { errCode, data } = await purchaseApi.getOrders({
      status: statusFilter.value,
      entries: pageSize.value,
      page: currentPage.value,
    })
    if (errCode === 0) {
      orders.value = data.all_orders ?? []
      totalCount.value = data.all_count ?? 0
    }
  } catch {
    // handled by interceptor
  } finally {
    loading.value = false
  }
}

function onStatusChange() {
  currentPage.value = 1
  fetchOrders()
}

function goToPage(p) {
  if (p < 1 || p > totalPages.value) return
  currentPage.value = p
  fetchOrders()
}

function handlePageSizeInput(e) {
  const val = parseInt(e.target.value)
  if (val > 0 && val <= 300) {
    pageSize.value = val
    currentPage.value = 1
    fetchOrders()
  }
}

function handleJumpPageInput(e) {
  const val = parseInt(e.target.value)
  if (val >= 1 && val <= totalPages.value) {
    goToPage(val)
  }
}

function jumpToOrder(id) {
  router.push(`/purchase/showorder/${id}`)
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  return new Intl.DateTimeFormat(locale.value, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  }).format(d)
}

async function confirmLockAll() {
  lockAllLoading.value = true
  try {
    const { errCode, data } = await purchaseApi.lockAllCompleted()
    if (errCode === 0) {
      const count = data?.locked_count ?? 0
      if (count > 0) {
        toast.success(t('sysadmin.lock_all_success', { count }))
      } else {
        toast.info(t('sysadmin.lock_all_none'))
      }
      await fetchOrders()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    toast.error(t('message.server_error'))
  } finally {
    lockAllLoading.value = false
    showLockAllConfirm.value = false
  }
}

defineExpose({ fetchOrders })

onActivated(() => fetchOrders())
</script>

<template>
  <div class="space-y-4">
    <!-- 顶部操作栏 -->
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-dk-text">
        {{ t('sysadmin.tab_purchase_orders') }}
      </h2>
      <div class="flex items-center gap-2">
        <span class="text-sm text-gray-500 dark:text-dk-subtle">
          {{ t('purchase.There_are_a_total_of') }} {{ totalCount }} {{ t('purchase.items') }}
        </span>
        <button
          class="flex items-center gap-1.5 rounded-md border border-gray-300 px-3 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:text-dk-subtle dark:hover:bg-dk-card"
          :disabled="loading"
          @click="fetchOrders"
        >
          <IconRefresh :size="16" :class="{ 'animate-spin': loading }" />
        </button>
        <button
          class="flex items-center gap-1.5 rounded-md bg-orange-500 px-4 py-1.5 text-sm font-medium text-white transition-colors hover:bg-orange-600 disabled:opacity-50"
          :disabled="lockAllLoading"
          @click="showLockAllConfirm = true"
        >
          <IconLock :size="16" />
          {{ t('sysadmin.lock_all_completed') }}
        </button>
      </div>
    </div>

    <!-- 状态筛选 -->
    <div class="flex flex-wrap gap-2">
      <button
        v-for="opt in statusOptions"
        :key="opt.value"
        class="rounded-full border px-3 py-1 text-xs font-medium transition-all"
        :class="
          statusFilter === opt.value
            ? 'border-blue-500 bg-blue-500 text-white'
            : 'border-gray-200 text-gray-500 hover:border-gray-300 hover:bg-gray-50 dark:border-dk-muted dark:text-dk-subtle dark:hover:bg-dk-card'
        "
        @click="statusFilter = opt.value; onStatusChange()"
      >
        {{ t('purchase.' + opt.labelKey) }}
      </button>
    </div>

    <!-- 订单表格 -->
    <div class="overflow-hidden rounded-md border border-gray-200 dark:border-dk-muted">
      <table class="min-w-full divide-y divide-gray-200 dark:divide-dk-muted">
        <thead class="bg-gray-50 dark:bg-dk-base">
          <tr>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">ID</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('purchase_addorder.part_name') }}</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('purchase.status') }}</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('purchase.created_at') }}</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('purchase.updated_at') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 bg-white dark:divide-dk-muted dark:bg-dk-card">
          <tr v-if="loading" class="text-center">
            <td colspan="5" class="py-8 text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.loading') }}</td>
          </tr>
          <tr v-else-if="orders.length === 0" class="text-center">
            <td colspan="5" class="py-8 text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.no_data') }}</td>
          </tr>
          <tr
            v-for="order in orders"
            :key="order.ID"
            class="cursor-pointer transition-colors hover:bg-blue-50/50 dark:hover:bg-dk-base/50"
            @click="jumpToOrder(order.ID)"
          >
            <td class="whitespace-nowrap px-4 py-3 text-sm text-gray-400">{{ order.ID }}</td>
            <td class="px-4 py-3 text-sm font-medium text-gray-900 dark:text-dk-text max-w-[200px] truncate">{{ order.Title || '-' }}</td>
            <td class="whitespace-nowrap px-4 py-3">
              <span
                class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-semibold"
                :class="{
                  'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/40 dark:text-yellow-400': order.OrderStatus === 'pending',
                  'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-400': order.OrderStatus === 'ordered',
                  'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-400': order.OrderStatus === 'arrived',
                  'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400': order.OrderStatus === 'received',
                  'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400': order.OrderStatus === 'lost',
                  'bg-gray-200 text-gray-600 dark:bg-gray-700 dark:text-gray-300': order.OrderStatus === 'returned',
                }"
              >
                {{ t('purchase.status_' + order.OrderStatus) }}
              </span>
              <span
                v-if="order.Locked"
                class="ml-1 inline-flex items-center text-red-500 dark:text-red-400"
                :title="t('purchase.locked')"
              >
                <IconLock :size="14" />
              </span>
            </td>
            <td class="whitespace-nowrap px-4 py-3 text-sm text-gray-500 dark:text-dk-subtle">{{ formatDate(order.CreatedAt) }}</td>
            <td class="whitespace-nowrap px-4 py-3 text-sm text-gray-500 dark:text-dk-subtle">{{ formatDate(order.UpdatedAt) }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 分页 -->
    <div class="flex flex-col items-center justify-between gap-3 sm:flex-row">
      <div class="flex items-center gap-1.5 text-sm text-gray-500 dark:text-dk-subtle">
        <label>{{ t('purchase.show') }}</label>
        <input
          type="text"
          class="w-14 rounded border border-gray-300 px-2 py-1 text-center text-sm text-gray-900 dark:border-dk-muted dark:bg-dk-base dark:text-white"
          :value="pageSize"
          @change="handlePageSizeInput"
        />
        <label>{{ t('purchase.entries') }}</label>
      </div>
      <div class="flex items-center gap-1">
        <button
          class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card"
          :disabled="currentPage <= 1"
          @click="goToPage(1)"
        >
          <IconChevronsLeft :size="16" />
        </button>
        <button
          class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card"
          :disabled="currentPage <= 1"
          @click="goToPage(currentPage - 1)"
        >
          <IconChevronLeftPipe :size="16" />
        </button>
        <template v-for="a in pageRange" :key="a">
          <button
            class="min-w-[32px] rounded px-2 py-1 text-sm font-medium transition-colors"
            :class="a === currentPage ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100 dark:text-dk-subtle dark:hover:bg-dk-card'"
            @click="goToPage(a)"
          >
            {{ a }}
          </button>
        </template>
        <button
          class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card"
          :disabled="currentPage >= totalPages"
          @click="goToPage(currentPage + 1)"
        >
          <IconChevronRightPipe :size="16" />
        </button>
        <button
          class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card"
          :disabled="currentPage >= totalPages"
          @click="goToPage(totalPages)"
        >
          <IconChevronsRight :size="16" />
        </button>
        <input
          type="text"
          class="ml-2 w-14 rounded border border-gray-300 px-2 py-1 text-center text-sm text-gray-900 dark:border-dk-muted dark:bg-dk-base dark:text-white"
          @change="handleJumpPageInput"
        />
      </div>
    </div>

    <!-- 一键锁定确认弹窗 -->
    <ConfirmDialog
      v-model="showLockAllConfirm"
      :message="t('sysadmin.lock_all_confirm')"
      danger
      @confirm="confirmLockAll"
    />
  </div>
</template>
