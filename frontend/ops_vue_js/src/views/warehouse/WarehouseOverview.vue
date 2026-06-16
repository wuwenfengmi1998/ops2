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
//  合并列表（容器 + 物品，类似文件管理器）
// ═══════════════════════════════════════════════════════
//
// 策略：两个接口分别请求一较大窗口的数据（足以覆盖前端分页），
// 在前端合并为统一的 entries 列表（容器在前、物品在后），
// 共用一个搜索框 + 一个分页器。

// 一次性向后端取的最大条数（后端上限 300）
const FETCH_LIMIT = 300

const containers = ref([])           // 原始容器列表
const items = ref([])                // 原始物品列表
const containerTotal = ref(0)
const itemTotal = ref(0)

const search = ref('')               // 共用搜索词
const loading = ref(false)

const page = ref(1)
const pageSize = ref(20)

// 合并后的 entries：[{kind:'container', data:c}, {kind:'item', data:i}]
const entries = computed(() => {
  const cs = containers.value.map(c => ({ kind: 'container', id: `c-${c.ID}`, data: c }))
  const is = items.value.map(i => ({ kind: 'item', id: `i-${i.ID}`, data: i }))
  return [...cs, ...is]
})

const totalEntries = computed(() => entries.value.length)
const totalPages = computed(() => Math.ceil(totalEntries.value / pageSize.value) || 1)

const pagedEntries = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return entries.value.slice(start, start + pageSize.value)
})

function pageRange() {
  const total = totalPages.value
  const cur = page.value
  let start = Math.max(1, cur - 2)
  let end = Math.min(cur + 4, total)
  if (end - start < 4) start = Math.max(1, end - 4)
  return Array.from({ length: end - start + 1 }, (_, i) => start + i)
}

async function fetchStats() {
  try {
    const { errCode, data } = await warehouseApi.getCount()
    if (errCode === 0) {
      stats.container_total = data.container_total ?? 0
      stats.item_total = data.item_total ?? 0
      stats.unstored_items = data.unstored_items ?? 0
    }
  } catch { /* silent */ }
}

async function fetchAll() {
  loading.value = true
  try {
    const [cRes, iRes] = await Promise.all([
      warehouseApi.getContainers({
        search: search.value.trim(),
        entries: FETCH_LIMIT,
        page: 1,
      }),
      warehouseApi.getItems({
        search: search.value.trim(),
        unstored: true,
        entries: FETCH_LIMIT,
        page: 1,
      }),
    ])
    if (cRes.errCode === 0) {
      containers.value = cRes.data.containers ?? []
      containerTotal.value = cRes.data.all_count ?? 0
    } else {
      containers.value = []
      containerTotal.value = 0
    }
    if (iRes.errCode === 0) {
      // 根目录类比文件管理器：只显示「未入库」的物品（unstored=true 由后端过滤）；
      // 已归入容器的物品在对应容器详情页里显示。
      items.value = iRes.data.items ?? []
      itemTotal.value = iRes.data.all_count ?? 0
    } else {
      items.value = []
      itemTotal.value = 0
    }
    page.value = 1
  } catch {
    /* interceptor handled */
  } finally {
    loading.value = false
  }
}

let searchTimer = null
function onSearchInput() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => fetchAll(), 400)
}

function goPage(p) {
  if (p < 1 || p > totalPages.value) return
  page.value = p
}

function handlePageSize(e) {
  let val = parseInt(e.target.value) || 20
  if (val > 300) val = 300
  if (val < 1) val = 1
  pageSize.value = val
  page.value = 1
}

// ── 容器：新增 / 跳转 ──
const showContainerForm = ref(false)
const containerFormTitle = ref('')
const editingContainerId = ref(null)
const containerForm = reactive({ title: '', remark: '' })
const submittingContainer = ref(false)

function openAddContainer() {
  containerFormTitle.value = t('warehouse.add_container')
  editingContainerId.value = null
  containerForm.title = ''
  containerForm.remark = ''
  showContainerForm.value = true
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
      fetchAll()
      fetchStats()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch { /* interceptor handled */ }
  finally { submittingContainer.value = false }
}

function jumpToContainer(id) {
  router.push(`/warehouse/container/${id}`)
}

function jumpToItem(id) {
  router.push(`/warehouse/item/${id}`)
}

// ── 工具函数 ──
function fmtTs(ts) {
  if (!ts) return '—'
  let d
  if (typeof ts === 'number') {
    d = new Date(ts * 1000)
  } else if (typeof ts === 'string') {
    if (/^\d+$/.test(ts)) {
      d = new Date(parseInt(ts, 10) * 1000)
    } else {
      d = new Date(ts)
    }
  } else {
    d = new Date(ts)
  }
  if (isNaN(d.getTime())) return '—'
  return d.toLocaleString(isEn.value ? 'en-US' : 'zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', hour12: false
  })
}

// ── 初始化 ──
onMounted(() => {
  fetchStats()
  fetchAll()
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

    <!-- 主卡片：容器 + 物品（合并显示） -->
    <div class="rounded-xl border border-gray-200 bg-white shadow dark:border-dk-muted dark:bg-dk-card">
      <!-- Header -->
      <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('warehouse.contents') }}</h3>
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
        <div class="relative flex-1 max-w-md">
          <IconSearch class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" :size="16" />
          <input
            v-model="search"
            type="text"
            :placeholder="t('warehouse.search_contents_placeholder')"
            class="w-full rounded-lg border border-gray-300 bg-white py-1.5 pl-9 pr-3 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
            @input="onSearchInput"
            @keyup.enter="fetchAll"
          />
        </div>
        <div class="text-xs text-gray-400 dark:text-gray-500">
          {{ t('warehouse.containers') }}: {{ containerTotal }} · {{ t('warehouse.items') }}: {{ itemTotal }}
        </div>
      </div>

      <!-- 表格 -->
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm text-gray-900 dark:text-white">
          <thead>
            <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base dark:text-gray-400">
              <th class="px-5 py-3 font-medium w-80">{{ t('warehouse.container_name') }} / {{ t('warehouse.item_name') }}</th>
              <th class="px-5 py-3 font-medium">{{ t('warehouse.remark') }}</th>
              <th class="px-5 py-3 font-medium w-32 text-center">{{ t('warehouse.child_containers') }} / {{ t('warehouse.quantity') }}</th>
              <th class="px-5 py-3 font-medium w-28 text-center">{{ t('work_order.work_order_count') }}</th>
              <th class="px-5 py-3 font-medium">{{ t('customer.related_customers') }}</th>
              <th class="px-5 py-3 font-medium whitespace-nowrap w-44">{{ t('warehouse.updated_at') }}</th>
              <th class="px-5 py-3 font-medium">{{ t('warehouse.created_by') }}</th>
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
            <tr v-else-if="totalEntries === 0">
              <td colspan="7" class="px-6 py-8 text-center text-gray-400 dark:text-gray-500">
                {{ t('warehouse.no_contents') }}
              </td>
            </tr>
            <template v-else>
              <tr
                v-for="entry in pagedEntries"
                :key="entry.id"
                class="cursor-pointer border-b border-gray-100 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:hover:bg-dk-base"
                @click="entry.kind === 'container' ? jumpToContainer(entry.data.ID) : jumpToItem(entry.data.ID)"
              >
                <!-- 名称（容器：文件夹图标 / 物品：包裹图标） -->
                <td class="px-5 py-3">
                  <div class="flex items-center gap-2">
                    <IconFolder
                      v-if="entry.kind === 'container'"
                      class="flex-shrink-0 text-blue-500"
                      :size="18"
                    />
                    <IconPackage
                      v-else
                      class="flex-shrink-0 text-green-500"
                      :size="18"
                    />
                    <span class="max-w-xs truncate font-medium text-gray-900 dark:text-white">
                      {{ entry.kind === 'container' ? entry.data.Title : entry.data.Name }}
                    </span>
                    <span
                      v-if="entry.kind === 'item' && entry.data.SerialNumber"
                      class="text-xs text-gray-400 dark:text-gray-500 truncate"
                    >· {{ entry.data.SerialNumber }}</span>
                  </div>
                </td>

                <!-- 备注 -->
                <td class="px-5 py-3 max-w-xs truncate text-gray-500 dark:text-gray-400">
                  {{ entry.data.Remark || '—' }}
                </td>

                <!-- 容量 / 数量（容器显示子容器+物品数；物品显示数量） -->
                <td class="px-5 py-3 text-center">
                  <template v-if="entry.kind === 'container'">
                    <span class="inline-flex items-center gap-1 rounded-full bg-purple-100 px-2 py-0.5 text-xs font-medium text-purple-700 dark:bg-purple-900/40 dark:text-purple-400">
                      <IconFolders :size="12" />
                      {{ entry.data.ChildCount }}
                    </span>
                    <span class="ml-1 inline-flex items-center gap-1 rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-700 dark:bg-green-900/40 dark:text-green-400">
                      <IconPackage :size="12" />
                      {{ entry.data.ItemCount }}
                    </span>
                  </template>
                  <span v-else class="text-sm text-gray-700 dark:text-gray-300">{{ entry.data.Quantity }}</span>
                </td>

                <!-- 工单数（仅物品） -->
                <td class="px-5 py-3 text-center">
                  <template v-if="entry.kind === 'item'">
                    <span
                      v-if="entry.data.WorkOrderCount > 0"
                      class="inline-flex items-center gap-1 rounded-full bg-orange-100 px-2 py-0.5 text-xs font-medium text-orange-700 dark:bg-orange-900/40 dark:text-orange-400"
                    >
                      <IconTool :size="12" />
                      {{ entry.data.WorkOrderCount }}
                    </span>
                    <span v-else class="text-gray-400">—</span>
                  </template>
                  <span v-else class="text-gray-400">—</span>
                </td>

                <!-- 关联客户（仅物品） -->
                <td class="px-5 py-3">
                  <template v-if="entry.kind === 'item' && entry.data.Customers && entry.data.Customers.length > 0">
                    <div class="flex flex-wrap gap-1">
                      <RouterLink
                        v-for="customer in entry.data.Customers.slice(0, 3)"
                        :key="customer.id"
                        :to="`/customer/detail/${customer.id}`"
                        class="inline-flex items-center gap-1 rounded-full bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-700 hover:bg-blue-200 dark:bg-blue-900/40 dark:text-blue-400 dark:hover:bg-blue-900/60 whitespace-nowrap"
                        @click.stop
                      >
                        <IconUser :size="10" />
                        {{ customer.first_name }} {{ customer.last_name }}
                      </RouterLink>
                      <span v-if="entry.data.Customers.length > 3" class="text-xs text-gray-400">
                        +{{ entry.data.Customers.length - 3 }}
                      </span>
                    </div>
                  </template>
                  <span v-else class="text-gray-400">—</span>
                </td>

                <!-- 更新时间 -->
                <td class="px-5 py-3 whitespace-nowrap text-xs text-gray-500 dark:text-gray-400">
                  {{ fmtTs(entry.data.UpdatedAt) }}
                </td>

                <!-- 创建人 -->
                <td class="px-5 py-3">
                  <div class="flex items-center gap-1.5">
                    <img
                      :src="usersStore.getAvatarUrlFromUserID(entry.data.CreatorID)"
                      class="w-5 h-5 rounded-full object-cover flex-shrink-0"
                    />
                    <span class="truncate text-gray-600 dark:text-gray-400">
                      {{ usersStore.getUsernameFromUserID(entry.data.CreatorID) }}
                    </span>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>

      <!-- 分页 -->
      <div class="flex flex-col items-center gap-3 border-t border-gray-100 px-6 py-4 sm:flex-row sm:justify-between dark:border-dk-muted">
        <div class="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
          <span>{{ t('warehouse.total_items', { count: totalEntries }) }}</span>
          <span>·</span>
          <span>{{ t('warehouse.containers') }}: {{ containers.length }}</span>
          <span>·</span>
          <span>{{ t('warehouse.items') }}: {{ items.length }}</span>
          <span class="ml-3">{{ isEn ? 'Per page' : '每页' }}</span>
          <input
            type="number"
            :value="pageSize"
            min="1"
            max="300"
            class="w-14 rounded border border-gray-300 px-1.5 py-0.5 text-center text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
            @change="handlePageSize"
          />
        </div>
        <div class="flex items-center gap-1">
          <button @click="goPage(1)" :disabled="page === 1" class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted">
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 17l-5-5 5-5M18 17l-5-5 5-5"/></svg>
          </button>
          <button @click="goPage(page - 1)" :disabled="page === 1" class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted">
            <IconChevronLeft :size="16" />
          </button>
          <button
            v-for="p in pageRange()" :key="p"
            @click="goPage(p)"
            :class="['rounded px-2.5 py-1 text-sm', p === page ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dk-muted']"
          >{{ p }}</button>
          <button @click="goPage(page + 1)" :disabled="page === totalPages" class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted">
            <IconChevronRight :size="16" />
          </button>
          <button @click="goPage(totalPages)" :disabled="page === totalPages" class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted">
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M13 17l5-5-5-5M6 17l5-5-5-5"/></svg>
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
</template>
