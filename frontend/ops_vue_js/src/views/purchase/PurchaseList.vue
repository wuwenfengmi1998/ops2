<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { purchaseApi } from '@/api/purchase'
import { IconPlus, IconChevronLeftPipe, IconChevronRightPipe, IconChevronsLeft, IconChevronsRight, IconSearch } from '@tabler/icons-vue'

usePageTitle('appname.purchase')
const { t, locale } = useI18n()
const router = useRouter()
const toast = useToastStore()

const orders = ref([])
const totalCount = ref(0)
const pageSize = ref(10)
const currentPage = ref(1)
const searchQuery = ref('')
const loading = ref(false)

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
    const { errCode, data } = await purchaseApi.getOrders({
      search: searchQuery.value,
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
  // const resolved = router.resolve({ path: `/purchase/showorder/${id}` })
  // window.open(resolved.href, '_blank')

  router.replace(`/purchase/showorder/${id}`);
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

function handleJumpPageInput(e) {
  const val = parseInt(e.target.value)
  if (val > 0 && val <= totalPages.value) {
    currentPage.value = val
    fetchOrders()
  }
}

onMounted(fetchOrders)
</script>

<template>
  <div class="mx-auto max-w-6xl px-6 py-6">
    <div class="flex flex-col gap-6 rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <!-- Header -->
      <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('purchase.purchase_list') }}</h3>
      </div>

      <!-- Toolbar -->
      <div class="flex flex-col gap-3 px-6 py-3 sm:flex-row sm:items-center">
        <div class="flex gap-2">
          <RouterLink to="/purchase/addorder" class="inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700">
            <IconPlus :size="16" />
            {{ t('purchase.add_part') }}
          </RouterLink>
          <button class="rounded-lg border border-gray-300 px-3 py-1.5 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:text-gray-400">{{ t('purchase.exp_report') }}</button>
        </div>
        <div class="flex items-center gap-2">
          <label class="text-sm text-gray-500">{{ t('purchase.search') }}</label>
          <input v-model="searchQuery" type="text" class="w-48 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white" @keydown.enter="currentPage=1;fetchOrders()" />
        </div>
      </div>

      <!-- Table -->
      <div class="overflow-x-auto px-0">
        <table class="w-full text-left text-sm text-gray-900">
          <thead>
            <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base">
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">No.</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('purchase.item_name') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('purchase.purpose') }}</th>
              
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap">{{ t('purchase.created_at') }}</th>
              <th class="px-6 py-3 font-medium text-gray-500 dark:text-gray-400">{{ t('purchase.status') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="6" class="px-6 py-8 text-center text-gray-400">
                <svg class="mx-auto mb-2 h-5 w-5 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                </svg>
                Loading...
              </td>
            </tr>
            <tr
              v-for="order in orders"
              :key="order.ID"
              class="cursor-pointer border-b border-gray-100 transition-colors hover:bg-blue-50/50 dark:border-dk-muted/50 dark:bg-dk-card dark:hover:bg-dk-base/50"
              @click="jumpToOrder(order.ID)"
            >
              <td class="px-6 py-3 text-gray-400">{{ order.ID }}</td>
              <td class="px-6 py-3 font-medium text-gray-900 dark:text-white">{{ order.Title }}</td>
              <td class="px-6 py-3 max-w-[200px] truncate text-gray-600 dark:text-gray-300">{{ order.Remark || '-' }}</td>
              <td class="px-6 py-3 whitespace-nowrap text-gray-500 dark:text-gray-400">{{ formatDate(order.CreatedAt) }}</td>
              <td class="px-6 py-3">
                <span
                  class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold"
                  :class="{
                    'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/40 dark:text-yellow-400': order.OrderStatus === 'pending',
                    'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-400': order.OrderStatus === 'ordered',
                    'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-400': order.OrderStatus === 'arrived',
                    'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400': order.OrderStatus === 'received',
                  }"
                >
                  {{ t('purchase.status_' + order.OrderStatus) }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="flex flex-col items-center justify-between gap-3 border-t border-gray-200 px-6 py-3 sm:flex-row dark:border-dk-muted">
        <div class="flex items-center gap-1.5 text-sm text-gray-500">
          <label>{{ t('purchase.show') }}</label>
          <input type="text" class="w-14 rounded border border-gray-300 px-2 py-1 text-center text-sm text-gray-900 dark:border-dk-muted dark:bg-dk-base dark:text-white" :value="pageSize" @change="handlePageSizeInput" />
          <label>{{ t('purchase.entries') }}</label>
          <span class="ml-1">{{ t('purchase.There_are_a_total_of') }} {{ totalCount }} {{ t('purchase.items') }}</span>
        </div>

        <div class="flex items-center gap-1">
          <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card" :disabled="currentPage <= 1" @click="goToPage(1)"><IconChevronsLeft :size="16" /></button>
          <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card" :disabled="currentPage <= 1" @click="goToPage(currentPage - 1)"><IconChevronLeftPipe :size="16" /></button>
          <template v-for="a in pageRange" :key="a">
            <button
              class="min-w-[32px] rounded px-2 py-1 text-sm font-medium transition-colors"
              :class="a === currentPage ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dk-card dark:text-gray-400 dark:hover:bg-dk-card'"
              @click="goToPage(a)"
            >{{ a }}</button>
          </template>
          <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card" :disabled="currentPage >= totalPages" @click="goToPage(currentPage + 1)"><IconChevronRightPipe :size="16" /></button>
          <button class="rounded p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 disabled:opacity-40 dark:hover:bg-dk-card" :disabled="currentPage >= totalPages" @click="goToPage(totalPages)"><IconChevronsRight :size="16" /></button>
          <input type="text" class="ml-2 w-14 rounded border border-gray-300 px-2 py-1 text-center text-sm text-gray-900 dark:border-dk-muted dark:bg-dk-base dark:text-white" @change="handleJumpPageInput" />
        </div>
      </div>
    </div>
  </div>
</template>
