<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { warehouseApi } from '@/api/warehouse'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import {
  IconChevronLeft,
  IconChevronRight,
  IconPackage,
  IconSearch,
  IconTrash,
  IconArrowRight,
} from '@tabler/icons-vue'

usePageTitle('warehouse.item_list')
const { t, locale } = useI18n()
const router = useRouter()
const toast = useToastStore()

// ── 状态 ──
const items = ref([])
const canModifyItems = ref([]) // 并行数组：与 items 下标对应
const totalCount = ref(0)
const pageSize = ref(10)
const currentPage = ref(1)
const search = ref('')
const loading = ref(false)

// ── 权限判断 ──
function canModifyItem(idx) {
  return canModifyItems.value[idx] === true
}

// ── 获取物品列表 ──
async function fetchItems() {
  loading.value = true
  try {
    const { errCode, data } = await warehouseApi.getItems({
      search: search.value.trim(),
      entries: pageSize.value,
      page: currentPage.value,
    })
    if (errCode === 0 && data) {
      items.value = data.items || []
      canModifyItems.value = data.canModifyItems || []
      totalCount.value = data.all_count || 0
      stats.total = data.all_count || 0
      stats.inContainer = items.value.filter(i => i.ContainerID != null).length
      stats.unstored = items.value.filter(i => i.ContainerID == null).length
    }
  } catch {
    toast.error(t('message.server_error'))
  } finally {
    loading.value = false
  }
}

// ── 搜索 ──
let searchTimer = null
function onSearchInput() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    currentPage.value = 1
    fetchItems()
  }, 400)
}

// ── 分页 ──
function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
    fetchItems()
  }
}
function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    fetchItems()
  }
}
function goPage(p) {
  if (p >= 1 && p <= totalPages.value && p !== currentPage.value) {
    currentPage.value = p
    fetchItems()
  }
}

const pageNumbers = computed(() => {
  const total = totalPages.value
  const cur = currentPage.value
  if (total <= 5) return Array.from({ length: total }, (_, i) => i + 1)
  const pages = []
  if (cur <= 3) {
    pages.push(1, 2, 3, 4, '...', total)
  } else if (cur >= total - 2) {
    pages.push(1, '...', total - 3, total - 2, total - 1, total)
  } else {
    pages.push(1, '...', cur - 1, cur, cur + 1, '...', total)
  }
  return pages
})

// ── 跳转物品详情 ──
function goToDetail(item) {
  router.push(`/warehouse/item/${item.ID}`)
}

// ── 删除 ──
const deleteTarget = ref(null)
const confirmDelete = ref(false)
const deletingItem = ref(false)

function askDelete(item) {
  deleteTarget.value = item
  confirmDelete.value = true
}

async function doDeleteItem() {
  if (!deleteTarget.value) return
  deletingItem.value = true
  try {
    const { errCode } = await warehouseApi.deleteItem(deleteTarget.value.ID)
    if (errCode === 0) {
      toast.success(t('message.delete_success'))
      confirmDelete.value = false
      deleteTarget.value = null
      fetchItems()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    toast.error(t('message.server_error'))
  } finally {
    deletingItem.value = false
  }
}

// ── 工具函数 ──
function formatDate(dateStr) {
  if (!dateStr) return '—'
  try {
    const d = new Date(dateStr)
    return d.toLocaleDateString(isEn.value ? 'en-US' : 'zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
  } catch {
    return dateStr
  }
}

onMounted(fetchItems)
</script>

<template>
  <div class="p-4 max-w-6xl mx-auto space-y-4">

    <!-- 统计卡片 -->
    <div class="grid grid-cols-3 gap-3">
      <div class="rounded-xl border border-gray-200 bg-white px-5 py-4 shadow dark:border-dk-muted dark:bg-dk-card">
        <div class="flex items-center gap-2 text-gray-500 dark:text-gray-400">
          <IconPackage :size="18" />
          <span class="text-sm">{{ t('warehouse.item_count') }}</span>
        </div>
        <div class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ stats.total }}</div>
      </div>
      <div class="rounded-xl border border-gray-200 bg-white px-5 py-4 shadow dark:border-dk-muted dark:bg-dk-card">
        <div class="flex items-center gap-2 text-gray-500 dark:text-gray-400">
          <IconArrowRight :size="18" />
          <span class="text-sm">{{ t('warehouse.items_in_containers') }}</span>
        </div>
        <div class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ stats.inContainer }}</div>
      </div>
      <div class="rounded-xl border border-gray-200 bg-white px-5 py-4 shadow dark:border-dk-muted dark:bg-dk-card">
        <div class="flex items-center gap-2 text-gray-500 dark:text-gray-400">
          <IconPackage :size="18" />
          <span class="text-sm">{{ t('warehouse.unstored_items') }}</span>
        </div>
        <div class="mt-1 text-2xl font-bold text-orange-600 dark:text-orange-400">{{ stats.unstored }}</div>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="rounded-xl border border-gray-200 bg-white px-4 py-3 flex items-center gap-3 dark:border-dk-muted dark:bg-dk-card">
      <IconSearch :size="18" class="text-gray-400 flex-shrink-0" />
      <input
        v-model="search"
        class="flex-1 bg-transparent outline-none text-sm text-gray-900 dark:text-white"
        :placeholder="t('warehouse.search_item_placeholder')"
        @input="onSearchInput"
      />
    </div>

    <!-- 表格卡片 -->
    <div class="rounded-xl border border-gray-200 bg-white shadow dark:border-dk-muted dark:bg-dk-card">
      <!-- 加载 -->
      <div v-if="loading" class="flex justify-center py-12">
        <svg class="h-6 w-6 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
        </svg>
      </div>

      <!-- 空状态 -->
      <div v-else-if="items.length === 0" class="flex flex-col items-center justify-center py-16 text-gray-400">
        <IconPackage :size="40" class="mb-3 opacity-30" />
        <span class="text-sm">{{ t('warehouse.no_items') }}</span>
      </div>

      <!-- 表格 -->
      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-sm text-gray-900 dark:text-white">
          <thead>
            <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base dark:text-gray-400">
              <th class="px-5 py-3 font-medium">{{ t('warehouse.item_name') }}</th>
              <th class="px-5 py-3 font-medium">{{ t('warehouse.serial_number') }}</th>
              <th class="px-5 py-3 font-medium w-20 text-center">{{ t('warehouse.quantity') }}</th>
              <th class="px-5 py-3 font-medium">{{ t('warehouse.location') }}</th>
              <th class="px-5 py-3 font-medium whitespace-nowrap">{{ t('warehouse.created_at') }}</th>
              <th class="px-5 py-3 font-medium w-16 text-right">{{ t('warehouse.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(item, idx) in items"
              :key="item.ID"
              class="border-b border-gray-100 cursor-pointer transition-colors hover:bg-gray-50 dark:border-dk-muted dark:hover:bg-dk-base"
              @click="goToDetail(item)"
            >
              <td class="px-5 py-3 font-medium max-w-[200px] truncate">{{ item.Name }}</td>
              <td class="px-5 py-3 text-xs text-gray-500 dark:text-gray-400 max-w-[160px] truncate">{{ item.SerialNumber || '—' }}</td>
              <td class="px-5 py-3 text-center text-sm">{{ item.Quantity }}</td>
              <td class="px-5 py-3">
                <span v-if="item.ContainerBreadcrumb" class="inline-flex items-center gap-1 text-blue-600 text-sm">
                  <IconArrowRight :size="13" />
                  <span class="truncate max-w-[200px]">{{ item.ContainerBreadcrumb }}</span>
                </span>
                <span v-else class="inline-flex items-center gap-1 text-xs text-orange-500">
                  {{ t('warehouse.unstored_items') }}
                </span>
              </td>
              <td class="px-5 py-3 text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ formatDate(item.CreatedAt) }}</td>
              <td class="px-5 py-3 text-right" @click.stop>
                <button
                  v-if="canModifyItem(idx)"
                  class="inline-flex items-center justify-center w-7 h-7 rounded text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20"
                  @click="askDelete(item)"
                >
                  <IconTrash :size="14" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 分页 -->
      <div v-if="totalPages > 1" class="flex items-center justify-between px-5 py-3 border-t border-gray-100 dark:border-dk-muted">
        <div class="text-xs text-gray-400 dark:text-gray-500">
          {{ t('warehouse.total_items', { count: totalCount }) }}
        </div>
        <div class="flex items-center gap-1">
          <button
            class="flex items-center justify-center w-7 h-7 rounded text-sm text-gray-500 hover:bg-gray-100 disabled:opacity-30 disabled:cursor-not-allowed dark:text-gray-400 dark:hover:bg-dk-muted"
            :disabled="currentPage === 1"
            @click="prevPage"
          >
            <IconChevronLeft :size="15" />
          </button>
          <template v-for="p in pageNumbers" :key="p">
            <span v-if="p === '...'" class="px-1 text-gray-400 text-sm">…</span>
            <button
              v-else
              class="flex items-center justify-center w-7 h-7 rounded text-sm"
              :class="p === currentPage
                ? 'bg-blue-600 text-white font-medium'
                : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dk-muted'"
              @click="goPage(p)"
            >
              {{ p }}
            </button>
          </template>
          <button
            class="flex items-center justify-center w-7 h-7 rounded text-sm text-gray-500 hover:bg-gray-100 disabled:opacity-30 disabled:cursor-not-allowed dark:text-gray-400 dark:hover:bg-dk-muted"
            :disabled="currentPage === totalPages"
            @click="nextPage"
          >
            <IconChevronRight :size="15" />
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- 删除确认 -->
  <ConfirmDialog
    v-model="confirmDelete"
    :title="t('warehouse.delete_item_title')"
    :message="t('warehouse.delete_item_msg', { name: deleteTarget?.name })"
    @confirm="doDeleteItem"
  />
</template>
