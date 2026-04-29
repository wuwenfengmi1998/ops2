<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { useUsersStore } from '@/stores/users'
import { usePageTitle } from '@/composables/usePageTitle'
import { warehouseApi } from '@/api/warehouse'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import {
  IconPlus,
  IconChevronLeft,
  IconChevronRight,
  IconFolder,
  IconFolders,
  IconPackage,
  IconSearch,
  IconTrash,
  IconEdit,
  IconArrowRight,
  IconTool,
  IconUser,
} from '@tabler/icons-vue'

usePageTitle('warehouse.overview')
const { t, locale } = useI18n()
const router = useRouter()
const toast = useToastStore()
const usersStore = useUsersStore()

const isEn = computed(() => locale.value === 'en')

// ── 统计 ──
const stats = reactive({
  container_total: 0,
  item_total: 0,
  unstored_items: 0,
})

// ═══════════════════════════════════════════════════════
//  容器相关
// ═══════════════════════════════════════════════════════
const containers = ref([])
const canModifyContainers = ref([]) // 并行数组：与 containers 下标对应
const containerTotal = ref(0)
const containerPage = ref(1)
const containerPageSize = ref(10)
const containerSearch = ref('')
const containerLoading = ref(false)

// 新增/编辑弹窗
const showContainerForm = ref(false)
const containerFormTitle = ref('')
const editingContainerId = ref(null)
const containerForm = reactive({ title: '', remark: '' })
const submittingContainer = ref(false)

// 删除确认
const showDeleteConfirm = ref(false)
const deletingId = ref(null)
const deletingName = ref('')
const deleting = ref(false)

const containerTotalPages = computed(() => Math.ceil(containerTotal.value / containerPageSize.value) || 1)

function containerPageRange() {
  const total = containerTotalPages.value
  const cur = containerPage.value
  let start = Math.max(1, cur - 2)
  let end = Math.min(cur + 4, total)
  if (end - start < 4) start = Math.max(1, end - 4)
  return Array.from({ length: end - start + 1 }, (_, i) => start + i)
}

// ── 权限判断 ──
function canModifyContainer(index) {
  return canModifyContainers.value[index] === true
}
function canModifyItem(index) {
  return canModifyItems.value[index] === true
}

async function fetchContainerStats() {
  try {
    const { errCode, data } = await warehouseApi.getCount()
    if (errCode === 0) {
      stats.container_total = data.container_total ?? 0
      stats.item_total = data.item_total ?? 0
      stats.unstored_items = data.unstored_items ?? 0
    }
  } catch { /* silent */ }
}

async function fetchContainers() {
  containerLoading.value = true
  try {
    const { errCode, data } = await warehouseApi.getContainers({
      search: containerSearch.value,
      entries: containerPageSize.value,
      page: containerPage.value,
    })
    if (errCode === 0) {
      containers.value = data.containers ?? []
      canModifyContainers.value = data.canModifyContainers ?? []
      containerTotal.value = data.all_count ?? 0
    } else {
      toast.error(t('message.server_error'))
    }
  } catch { /* interceptor handled */ }
  finally { containerLoading.value = false }
}

function goContainerPage(page) {
  if (page < 1 || page > containerTotalPages.value) return
  containerPage.value = page
  fetchContainers()
}

function handleContainerPageSize(e) {
  let val = parseInt(e.target.value) || 10
  if (val > 300) val = 300
  if (val < 1) val = 1
  containerPageSize.value = val
  containerPage.value = 1
  fetchContainers()
}

function openAddContainer() {
  containerFormTitle.value = t('warehouse.add_container')
  editingContainerId.value = null
  containerForm.title = ''
  containerForm.remark = ''
  showContainerForm.value = true
}

async function openEditContainer(id, e) {
  e.stopPropagation()
  try {
    const { errCode, data } = await warehouseApi.getContainer(id)
    if (errCode === 0) {
      containerFormTitle.value = t('warehouse.edit_container')
      editingContainerId.value = id
      containerForm.title = data.container.title ?? ''
      containerForm.remark = data.container.remark ?? ''
      showContainerForm.value = true
    } else {
      toast.error(t('message.server_error'))
    }
  } catch { /* interceptor handled */ }
}

async function submitContainerForm() {
  if (!containerForm.title.trim()) {
    toast.warning(t('warehouse.title_required'))
    return
  }
  submittingContainer.value = true
  try {
    const payload = { title: containerForm.title.trim(), remark: containerForm.remark.trim() }
    const { errCode } = editingContainerId.value
      ? await warehouseApi.updateContainer({ id: editingContainerId.value, ...payload })
      : await warehouseApi.addContainer(payload)
    if (errCode === 0) {
      showContainerForm.value = false
      toast.success(t('message.save_success'))
      fetchContainers()
      fetchContainerStats()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch { /* interceptor handled */ }
  finally { submittingContainer.value = false }
}

function confirmDeleteContainer(id, title, e) {
  e.stopPropagation()
  deletingId.value = id
  deletingName.value = title
  showDeleteConfirm.value = true
}

async function doDeleteContainer() {
  deleting.value = true
  try {
    const { errCode } = await warehouseApi.deleteContainer(deletingId.value)
    if (errCode === 0) {
      toast.success(t('message.delete_ok'))
      showDeleteConfirm.value = false
      fetchContainers()
      fetchContainerStats()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch { /* silent */ }
  finally { deleting.value = false }
}

function jumpToContainer(id) {
  router.push(`/warehouse/container/${id}`)
}

// ═══════════════════════════════════════════════════════
//  物品相关
// ═══════════════════════════════════════════════════════
const items = ref([])
const canModifyItems = ref([]) // 并行数组：与 items 下标对应
const itemTotal = ref(0)
const itemPage = ref(1)
const itemPageSize = ref(10)
const itemSearch = ref('')
const itemLoading = ref(false)
const itemTotalPages = computed(() => Math.ceil(itemTotal.value / itemPageSize.value) || 1)

const itemStats = reactive({ total: 0, inContainer: 0, unstored: 0 })

async function fetchItems() {
  itemLoading.value = true
  try {
    const { errCode, data } = await warehouseApi.getItems({
      search: itemSearch.value.trim(),
      entries: itemPageSize.value,
      page: itemPage.value,
    })
    if (errCode === 0 && data) {
      items.value = data.items || []
      canModifyItems.value = data.canModifyItems || []
      itemTotal.value = data.all_count || 0
      itemStats.total = data.all_count || 0
      itemStats.inContainer = items.value.filter(i => i.ContainerID != null).length
      itemStats.unstored = items.value.filter(i => i.ContainerID == null).length
    }
  } catch { toast.error(t('message.server_error')) }
  finally { itemLoading.value = false }
}

let searchTimer = null
function onItemSearchInput() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    itemPage.value = 1
    fetchItems()
  }, 400)
}

function prevItemPage() {
  if (itemPage.value > 1) { itemPage.value--; fetchItems() }
}
function nextItemPage() {
  if (itemPage.value < itemTotalPages.value) { itemPage.value++; fetchItems() }
}
function goItemPage(p) {
  if (p >= 1 && p <= itemTotalPages.value && p !== itemPage.value) { itemPage.value = p; fetchItems() }
}

const itemPageNumbers = computed(() => {
  const total = itemTotalPages.value
  const cur = itemPage.value
  if (total <= 5) return Array.from({ length: total }, (_, i) => i + 1)
  const pages = []
  if (cur <= 3) pages.push(1, 2, 3, 4, '...', total)
  else if (cur >= total - 2) pages.push(1, '...', total - 3, total - 2, total - 1, total)
  else pages.push(1, '...', cur - 1, cur, cur + 1, '...', total)
  return pages
})

function goToItemDetail(item) {
  router.push(`/warehouse/item/${item.ID}`)
}

const deleteItemTarget = ref(null)
const showDeleteItemConfirm = ref(false)
const deletingItem = ref(false)

function askDeleteItem(item) {
  deleteItemTarget.value = item
  showDeleteItemConfirm.value = true
}

async function doDeleteItem() {
  if (!deleteItemTarget.value) return
  deletingItem.value = true
  try {
    const { errCode } = await warehouseApi.deleteItem(deleteItemTarget.value.ID)
    if (errCode === 0) {
      toast.success(t('message.delete_success'))
      showDeleteItemConfirm.value = false
      deleteItemTarget.value = null
      fetchItems()
      fetchContainerStats()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch { toast.error(t('message.server_error')) }
  finally { deletingItem.value = false }
}

// ── 工具函数 ──
function formatDate(dateStr) {
  if (!dateStr) return '—'
  try {
    let d
    if (typeof dateStr === 'string' && /^\d+$/.test(dateStr)) {
      // Unix timestamp string like "1712345678"
      d = new Date(parseInt(dateStr, 10) * 1000)
    } else {
      d = new Date(dateStr)
    }
    if (isNaN(d.getTime())) return '—'
    return d.toLocaleString(isEn.value ? 'en-US' : 'zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false })
  } catch { return dateStr }
}

function fmtTs(ts) {
  if (!ts) return '—'
  let d
  if (typeof ts === 'number') {
    // Unix timestamp in seconds (number)
    d = new Date(ts * 1000)
  } else if (typeof ts === 'string') {
    // Check if it's a Unix timestamp string like "1712345678"
    if (/^\d+$/.test(ts)) {
      d = new Date(parseInt(ts, 10) * 1000)
    } else {
      // ISO 8601 or other string format
      d = new Date(ts)
    }
  } else {
    d = new Date(ts)
  }
  if (isNaN(d.getTime())) return '—'
  return d.toLocaleString()
}

// ── 初始化 ──
onMounted(() => {
  fetchContainerStats()
  fetchContainers()
  fetchItems()
})
</script>

<template>
  <div class="mx-auto max-w-7xl px-6 py-6 space-y-4">

    <!-- 统计卡片 -->
    <div class="grid grid-cols-3 gap-4">
      <div class="rounded-xl border border-gray-200 bg-white px-5 py-4 shadow dark:border-dk-muted dark:bg-dk-card">
        <div class="flex items-center gap-2 text-gray-500 dark:text-gray-400">
          <IconFolders :size="18" />
          <span class="text-sm">{{ t('warehouse.container_count') }}</span>
        </div>
        <div class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ stats.container_total }}</div>
      </div>
      <div class="rounded-xl border border-gray-200 bg-white px-5 py-4 shadow dark:border-dk-muted dark:bg-dk-card">
        <div class="flex items-center gap-2 text-gray-500 dark:text-gray-400">
          <IconPackage :size="18" />
          <span class="text-sm">{{ t('warehouse.item_count') }}</span>
        </div>
        <div class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ stats.item_total }}</div>
      </div>
      <div class="rounded-xl border border-gray-200 bg-white px-5 py-4 shadow dark:border-dk-muted dark:bg-dk-card">
        <div class="flex items-center gap-2 text-gray-500 dark:text-gray-400">
          <IconPackage :size="18" />
          <span class="text-sm">{{ t('warehouse.unstored_items') }}</span>
        </div>
        <div class="mt-1 text-2xl font-bold text-orange-600 dark:text-orange-400">{{ stats.unstored_items }}</div>
      </div>
    </div>

    <!-- 主卡片：容器 -->
    <div class="rounded-xl border border-gray-200 bg-white shadow dark:border-dk-muted dark:bg-dk-card">
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('warehouse.container_list') }}</h3>
          <button
            @click="openAddContainer"
            class="inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700"
          >
            <IconPlus :size="16" />
            {{ t('warehouse.add_container') }}
          </button>
        </div>

        <!-- 搜索栏 -->
        <div class="flex items-center gap-3 border-b border-gray-100 px-6 py-3 dark:border-dk-muted">
          <div class="relative flex-1 max-w-xs">
            <IconSearch class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" :size="16" />
            <input
              v-model="containerSearch"
              type="text"
              :placeholder="t('warehouse.search_placeholder')"
              class="w-full rounded-lg border border-gray-300 bg-white py-1.5 pl-9 pr-3 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
              @keyup.enter="containerPage = 1; fetchContainers()"
            />
          </div>
          <button
            @click="containerPage = 1; fetchContainers()"
            class="rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white hover:bg-gray-50 dark:hover:bg-dk-muted"
          >
            {{ t('purchase.search') }}
          </button>
        </div>

        <!-- 表格 -->
        <div class="overflow-x-auto">
          <table class="w-full text-left text-sm text-gray-900">
            <thead>
              <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base dark:text-gray-400">
                <th class="px-6 py-3 font-medium w-16">ID</th>
                <th class="px-6 py-3 font-medium">{{ t('warehouse.container_name') }}</th>
                <th class="px-6 py-3 font-medium">{{ t('warehouse.remark') }}</th>
                <th class="px-6 py-3 font-medium w-24 text-center">{{ t('warehouse.child_containers') }}</th>
                <th class="px-6 py-3 font-medium w-24 text-center">{{ t('warehouse.items') }}</th>
                <th class="px-6 py-3 font-medium whitespace-nowrap w-44">{{ t('warehouse.created_at') }}</th>
                <th class="px-6 py-3 font-medium whitespace-nowrap w-44">{{ t('warehouse.updated_at') }}</th>
                <th class="px-6 py-3 font-medium">{{ t('warehouse.created_by') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="containerLoading">
                <td colspan="8" class="px-6 py-8 text-center text-gray-400">
                  <svg class="mx-auto mb-2 h-5 w-5 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                  </svg>
                  {{ t('message.loading') }}
                </td>
              </tr>
              <tr v-else-if="containers.length === 0">
                <td colspan="8" class="px-6 py-8 text-center text-gray-400 dark:text-gray-500">
                  {{ t('warehouse.no_containers') }}
                </td>
              </tr>
              <tr
                v-else
                v-for="(c, idx) in containers"
                :key="c.ID"
                class="cursor-pointer border-b border-gray-100 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:hover:bg-dk-base"
                @click="jumpToContainer(c.ID)"
              >
                <td class="px-6 py-3 text-gray-500 dark:text-gray-400">{{ c.ID }}</td>
                <td class="px-6 py-3">
                  <div class="flex items-center gap-2">
                    <IconFolder class="flex-shrink-0 text-blue-500" :size="18" />
                    <span class="max-w-xs truncate font-medium text-gray-900 dark:text-white">{{ c.Title }}</span>
                  </div>
                </td>
                <td class="px-6 py-3 max-w-xs truncate text-gray-500 dark:text-gray-400">{{ c.Remark || '—' }}</td>
                <td class="px-6 py-3 text-center">
                  <span class="inline-flex items-center gap-1 rounded-full bg-purple-100 px-2 py-0.5 text-xs font-medium text-purple-700 dark:bg-purple-900/40 dark:text-purple-400">
                    <IconFolders :size="12" />
                    {{ c.ChildCount }}
                  </span>
                </td>
                <td class="px-6 py-3 text-center">
                  <span class="inline-flex items-center gap-1 rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-700 dark:bg-green-900/40 dark:text-green-400">
                    <IconPackage :size="12" />
                    {{ c.ItemCount }}
                  </span>
                </td>
                <td class="px-6 py-3 whitespace-nowrap text-gray-500 dark:text-gray-400">{{ fmtTs(c.CreatedAt) }}</td>
                <td class="px-6 py-3 whitespace-nowrap text-gray-500 dark:text-gray-400">{{ fmtTs(c.UpdatedAt) }}</td>
                <td class="px-6 py-3">
                  <div class="flex items-center gap-1.5">
                    <img
                      :src="usersStore.getAvatarUrlFromUserID(c.CreatorID)"
                      class="w-5 h-5 rounded-full object-cover flex-shrink-0"
                    />
                    <span class="truncate text-gray-600 dark:text-gray-400">{{ usersStore.getUsernameFromUserID(c.CreatorID) }}</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- 分页 -->
        <div class="flex flex-col items-center gap-3 border-t border-gray-100 px-6 py-4 sm:flex-row sm:justify-between dark:border-dk-muted">
          <div class="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
            <span>共 {{ containerTotal }} 条</span>
            <span>每页</span>
            <input
              type="number"
              :value="containerPageSize"
              min="1"
              max="300"
              class="w-14 rounded border border-gray-300 px-1.5 py-0.5 text-center text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
              @change="handleContainerPageSize"
            />
            <span>条</span>
          </div>
          <div class="flex items-center gap-1">
            <button @click="goContainerPage(1)" :disabled="containerPage === 1" class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted">
              <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 17l-5-5 5-5M18 17l-5-5 5-5"/></svg>
            </button>
            <button @click="goContainerPage(containerPage - 1)" :disabled="containerPage === 1" class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted">
              <IconChevronLeft :size="16" />
            </button>
            <button
              v-for="p in containerPageRange()" :key="p"
              @click="goContainerPage(p)"
              :class="['rounded px-2.5 py-1 text-sm', p === containerPage ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dk-muted']"
            >{{ p }}</button>
            <button @click="goContainerPage(containerPage + 1)" :disabled="containerPage === containerTotalPages" class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted">
              <IconChevronRight :size="16" />
            </button>
            <button @click="goContainerPage(containerTotalPages)" :disabled="containerPage === containerTotalPages" class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted">
              <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M13 17l5-5-5-5M6 17l5-5-5-5"/></svg>
            </button>
          </div>
        </div>
      </div>

    <!-- 主卡片：物品 -->
    <div class="rounded-xl border border-gray-200 bg-white shadow dark:border-dk-muted dark:bg-dk-card">
      <!-- Header -->
      <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('warehouse.item_list') }}</h3>
      </div>

      <!-- 搜索栏 -->
      <div class="flex items-center gap-3 border-b border-gray-100 px-6 py-3 dark:border-dk-muted">
        <IconSearch class="flex-shrink-0 text-gray-400" :size="18" />
        <input
          v-model="itemSearch"
          class="flex-1 bg-transparent text-sm text-gray-900 outline-none dark:text-white"
          :placeholder="t('warehouse.search_item_placeholder')"
          @input="onItemSearchInput"
        />
      </div>

      <!-- 表格 -->
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm text-gray-900 dark:text-white">
          <thead>
            <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base dark:text-gray-400">
              <th class="px-5 py-3 font-medium w-80">{{ t('warehouse.item_name') }}</th>
              <th class="px-5 py-3 font-medium w-20">{{ t('warehouse.serial_number') }}</th>
              <th class="px-5 py-3 font-medium w-50">{{ t('warehouse.remark') }}</th>
              <th class="px-5 py-3 font-medium  text-center">{{ t('warehouse.quantity') }}</th>
              <th class="px-5 py-3 font-medium  text-center">{{ t('work_order.work_order_count') }}</th>
              <th class="px-5 py-3 font-medium ">{{ t('customer.related_customers') }}</th>
              <th class="px-5 py-3 font-medium whitespace-nowrap">{{ t('warehouse.updated_at') }}</th>
              <th class="px-5 py-3 font-medium">{{ t('warehouse.created_by') }}</th>
              <th class="px-5 py-3 font-medium w-20 text-right">{{ t('warehouse.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="itemLoading">
              <td colspan="9" class="px-6 py-8 text-center text-gray-400">
                <svg class="mx-auto mb-2 h-5 w-5 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                </svg>
                {{ t('message.loading') }}
              </td>
            </tr>
            <tr v-else-if="items.length === 0">
              <td colspan="9" class="px-6 py-8 text-center text-gray-400 dark:text-gray-500">
                {{ t('warehouse.no_items') }}
              </td>
            </tr>
            <tr
              v-else
              v-for="(item, idx) in items" :key="item.ID"
              class="cursor-pointer border-b border-gray-100 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:hover:bg-dk-base"
              @click="goToItemDetail(item)"
            >
              <td class="px-5 py-3 font-medium max-w-[13rem] truncate">{{ item.Name }}</td>
              <td class="px-5 py-3 text-xs text-gray-500 dark:text-gray-400 max-w-[140px] truncate">{{ item.SerialNumber || '—' }}</td>
              <td class="px-5 py-3 text-xs text-gray-500 dark:text-gray-400 max-w-[200px] truncate">{{ item.Remark || '—' }}</td>
              <td class="px-5 py-3 text-center text-sm">{{ item.Quantity }}</td>
              <td class="px-5 py-3 text-center">
                <span v-if="item.WorkOrderCount > 0" class="inline-flex items-center gap-1 rounded-full bg-orange-100 px-2 py-0.5 text-xs font-medium text-orange-700 dark:bg-orange-900/40 dark:text-orange-400">
                  <IconTool :size="12" />
                  {{ item.WorkOrderCount }}
                </span>
                <span v-else class="text-gray-400">—</span>
              </td>
              <td class="px-5 py-3">
                <div v-if="item.Customers && item.Customers.length > 0" class="flex flex-wrap gap-1">
                  <RouterLink
                    v-for="customer in item.Customers.slice(0, 3)"
                    :key="customer.id"
                    :to="`/customer/detail/${customer.id}`"
                    class="inline-flex items-center gap-1 rounded-full bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-700 hover:bg-blue-200 dark:bg-blue-900/40 dark:text-blue-400 dark:hover:bg-blue-900/60 whitespace-nowrap"
                    @click.stop
                  >
                    <IconUser :size="10" />
                    {{ customer.first_name }} {{ customer.last_name }}
                  </RouterLink>
                  <span v-if="item.Customers.length > 3" class="text-xs text-gray-400">+{{ item.Customers.length - 3 }}</span>
                </div>
                <span v-else class="text-gray-400">—</span>
              </td>
              <td class="px-5 py-3 text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ formatDate(item.UpdatedAt) }}</td>
              <td class="px-5 py-3">
                <div class="flex items-center gap-1.5">
                  <img
                    :src="usersStore.getAvatarUrlFromUserID(item.CreatorID)"
                    class="w-5 h-5 rounded-full object-cover flex-shrink-0"
                  />
                  <span class="truncate text-gray-600 dark:text-gray-400">{{ usersStore.getUsernameFromUserID(item.CreatorID) }}</span>
                </div>
              </td>
              <td class="px-5 py-3 text-right">
                <button
                  class="text-xs text-blue-500 hover:underline"
                  @click.stop="goToItemDetail(item)"
                >
                  {{ t('warehouse.view_items') }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 分页 -->
      <div v-if="itemTotalPages > 1" class="flex items-center justify-between border-t border-gray-100 px-6 py-3 dark:border-dk-muted">
        <div class="text-xs text-gray-400 dark:text-gray-500">
          {{ t('warehouse.total_items', { count: itemTotal }) }}
        </div>
        <div class="flex items-center gap-1">
          <button
            class="flex h-7 w-7 items-center justify-center rounded text-sm text-gray-500 hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-30 dark:text-gray-400 dark:hover:bg-dk-muted"
            :disabled="itemPage === 1"
            @click="prevItemPage"
          >
            <IconChevronLeft :size="15" />
          </button>
          <template v-for="p in itemPageNumbers" :key="p">
            <span v-if="p === '...'" class="px-1 text-sm text-gray-400">…</span>
            <button
              v-else
              class="flex h-7 w-7 items-center justify-center rounded text-sm"
              :class="p === itemPage
                ? 'bg-blue-600 font-medium text-white'
                : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dk-muted'"
              @click="goItemPage(p)"
            >{{ p }}</button>
          </template>
          <button
            class="flex h-7 w-7 items-center justify-center rounded text-sm text-gray-500 hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-30 dark:text-gray-400 dark:hover:bg-dk-muted"
            :disabled="itemPage === itemTotalPages"
            @click="nextItemPage"
          >
            <IconChevronRight :size="15" />
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- 容器 新增/编辑弹窗 -->
  <Transition name="fade">
    <div
      v-if="showContainerForm"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/40"
      @click.self="showContainerForm = false"
    >
      <div class="w-full max-w-md rounded-xl border border-gray-200 bg-white p-5 shadow-xl dark:border-dk-muted dark:bg-dk-card">
        <h3 class="mb-4 text-base font-semibold text-gray-900 dark:text-white">{{ containerFormTitle }}</h3>
        <div class="space-y-3">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
              {{ t('warehouse.container_name') }} <span class="text-red-500">*</span>
            </label>
            <input
              v-model="containerForm.title"
              type="text"
              :placeholder="t('warehouse.title_placeholder')"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
              @keyup.enter="submitContainerForm"
            />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('warehouse.remark') }}</label>
            <textarea
              v-model="containerForm.remark"
              rows="3"
              :placeholder="t('warehouse.remark_placeholder')"
              class="w-full resize-none rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
            ></textarea>
          </div>
        </div>
        <div class="mt-4 flex justify-end gap-2">
          <button
            @click="showContainerForm = false"
            class="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white hover:bg-gray-50 dark:hover:bg-dk-muted"
          >
            {{ t('message.cancel') }}
          </button>
          <button
            @click="submitContainerForm"
            :disabled="submittingContainer"
            class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-60"
          >
            {{ submittingContainer ? t('message.submitting') : t('message.submit') }}
          </button>
        </div>
      </div>
    </div>
  </Transition>

  <!-- 容器 删除确认 -->
  <ConfirmDialog
    v-model="showDeleteConfirm"
    :title="t('warehouse.delete_confirm_title')"
    :message="t('warehouse.delete_confirm_msg', { name: deletingName })"
    :confirm-text="t('warehouse.delete')"
    :cancel-text="t('message.cancel')"
    danger
    @confirm="doDeleteContainer"
  />

  <!-- 物品 删除确认 -->
  <ConfirmDialog
    v-model="showDeleteItemConfirm"
    :title="t('warehouse.delete_item_title')"
    :message="t('warehouse.delete_item_msg', { name: deleteItemTarget?.name })"
    @confirm="doDeleteItem"
  />
</template>
