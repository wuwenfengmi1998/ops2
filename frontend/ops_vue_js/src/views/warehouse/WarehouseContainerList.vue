<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
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
  IconTool,
  IconUser,
} from '@tabler/icons-vue'

usePageTitle('warehouse.container_list')
const { t, locale } = useI18n()
const router = useRouter()
const toast = useToastStore()

// ── 状态 ──
const containers = ref([])
const totalCount = ref(0)
const pageSize = ref(10)
const currentPage = ref(1)
const loading = ref(false)
const search = ref('')

// 新增/编辑弹窗
const showForm = ref(false)
const formTitle = ref('')
const editingId = ref(null)
const form = reactive({
  title: '',
  remark: '',
  photos: [],
})
const submitting = ref(false)

// 删除确认
const showDeleteConfirm = ref(false)
const deletingId = ref(null)
const deletingName = ref('')

// 统计
const stats = reactive({ container_total: 0, item_total: 0, unstored_items: 0 })

// ── 分页 ──
const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value) || 1)

const pageRange = computed(() => {
  const total = totalPages.value
  const cur = currentPage.value
  let start = Math.max(1, cur - 2)
  let end = Math.min(cur + 4, total)
  if (end - start < 4) start = Math.max(1, end - 4)
  return Array.from({ length: end - start + 1 }, (_, i) => start + i)
})

// ── 拉取列表 ──
async function fetchContainers() {
  loading.value = true
  try {
    const { errCode, data } = await warehouseApi.getContainers({
      search: search.value,
      entries: pageSize.value,
      page: currentPage.value,
    })
    if (errCode === 0) {
      containers.value = data.containers ?? []
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

// ── 拉取统计 ──
async function fetchStats() {
  try {
    const { errCode, data } = await warehouseApi.getCount()
    if (errCode === 0) {
      stats.container_total = data.container_total ?? 0
      stats.item_total = data.item_total ?? 0
      stats.unstored_items = data.unstored_items ?? 0
    }
  } catch {
    // 静默
  }
}

// ── 翻页 ──
function goToPage(page) {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  fetchContainers()
}

function handlePageSizeInput(e) {
  let val = parseInt(e.target.value) || 10
  if (val > 300) val = 300
  if (val < 1) val = 1
  pageSize.value = val
  currentPage.value = 1
  fetchContainers()
}

function handleSearch() {
  currentPage.value = 1
  fetchContainers()
}

// ── 跳转到子容器 ──
function jumpToContainer(id) {
  router.push(`/warehouse/container/${id}`)
}

// ── 打开新增 ──
function openAdd() {
  formTitle.value = t('warehouse.add_container')
  editingId.value = null
  form.title = ''
  form.remark = ''
  form.photos = []
  showForm.value = true
}

// ── 打开编辑 ──
async function openEdit(id, e) {
  e.stopPropagation()
  try {
    const { errCode, data } = await warehouseApi.getContainer(id)
    if (errCode === 0) {
      formTitle.value = t('warehouse.edit_container')
      editingId.value = id
      form.title = data.container.title ?? ''
      form.remark = data.container.remark ?? ''
      form.photos = (data.photos ?? []).map((p) => p.sha256 ?? '')
      showForm.value = true
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  }
}

// ── 提交表单 ──
async function submitForm() {
  if (!form.title.trim()) {
    toast.warning(t('warehouse.title_required'))
    return
  }
  submitting.value = true
  try {
    const payload = { title: form.title.trim(), remark: form.remark, photos: form.photos }
    const { errCode } = editingId.value
      ? await warehouseApi.updateContainer({ id: editingId.value, ...payload })
      : await warehouseApi.addContainer(payload)
    if (errCode === 0) {
      showForm.value = false
      toast.success(t('message.save_success'))
      fetchContainers()
      fetchStats()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    submitting.value = false
  }
}

// ── 删除确认 ──
function confirmDelete(id, title, e) {
  e.stopPropagation()
  deletingId.value = id
  deletingName.value = title
  showDeleteConfirm.value = true
}

async function doDelete() {
  try {
    const { errCode } = await warehouseApi.deleteContainer(deletingId.value)
    if (errCode === 0) {
      toast.success(t('message.delete_ok'))
      showDeleteConfirm.value = false
      fetchContainers()
      fetchStats()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  }
}

// ── 格式化时间 ──
function formatDate(dateStr) {
  if (!dateStr) return ''
  const ts = parseInt(dateStr) * 1000
  if (isNaN(ts)) return ''
  return new Intl.DateTimeFormat(locale.value, {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', hour12: false,
  }).format(new Date(ts))
}

onMounted(() => {
  fetchContainers()
  fetchStats()
})
</script>

<template>
  <div class="mx-auto max-w-7xl px-6 py-6">

    <!-- 统计卡片 -->
    <div class="mb-6 grid grid-cols-3 gap-4">
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
        <div class="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{{ stats.unstored_items }}</div>
      </div>
    </div>

    <!-- 主列表卡片 -->
    <div class="rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">

      <!-- Header -->
      <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('warehouse.container_list') }}</h3>
        <button
          @click="openAdd"
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
            v-model="search"
            type="text"
            :placeholder="t('warehouse.search_placeholder')"
            class="w-full rounded-lg border border-gray-300 bg-white py-1.5 pl-9 pr-3 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
            @keyup.enter="handleSearch"
          />
        </div>
        <button
          @click="handleSearch"
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
              <th class="px-6 py-3 font-medium w-24 text-center">{{ t('work_order.work_order_count') }}</th>
              <th class="px-6 py-3 font-medium">{{ t('customer.related_customers') }}</th>
              <th class="px-6 py-3 font-medium whitespace-nowrap w-44">{{ t('warehouse.created_at') }}</th>
              <th class="px-6 py-3 font-medium w-28 text-right">{{ t('warehouse.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <!-- 加载中 -->
            <tr v-if="loading">
              <td colspan="9" class="px-6 py-8 text-center text-gray-400">
                <svg class="mx-auto mb-2 h-5 w-5 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                </svg>
                {{ t('message.loading') }}
              </td>
            </tr>
            <!-- 空状态 -->
            <tr v-else-if="containers.length === 0">
              <td colspan="9" class="px-6 py-8 text-center text-gray-400 dark:text-gray-500">
                {{ t('warehouse.no_containers') }}
              </td>
            </tr>
            <!-- 列表 -->
            <tr
              v-else
              v-for="c in containers"
              :key="c.ID"
              class="border-b border-gray-100 cursor-pointer transition-colors hover:bg-gray-50 dark:border-dk-muted dark:hover:bg-dk-base"
              @click="jumpToContainer(c.ID)"
            >
              <td class="px-6 py-3 text-gray-500 dark:text-gray-400">{{ c.ID }}</td>
              <td class="px-6 py-3">
                <div class="flex items-center gap-2">
                  <IconFolder class="text-blue-500 flex-shrink-0" :size="18" />
                  <span class="font-medium text-gray-900 dark:text-white max-w-xs truncate">{{ c.Title }}</span>
                </div>
              </td>
              <td class="px-6 py-3 text-gray-500 dark:text-gray-400 max-w-xs truncate">{{ c.Remark || '—' }}</td>
              <td class="px-6 py-3 text-center">
                <span
                  class="inline-flex items-center gap-1 rounded-full bg-purple-100 px-2 py-0.5 text-xs font-medium text-purple-700 dark:bg-purple-900/40 dark:text-purple-400"
                >
                  <IconFolders :size="12" />
                  {{ c.ChildCount }}
                </span>
              </td>
              <td class="px-6 py-3 text-center">
                <span
                  class="inline-flex items-center gap-1 rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-700 dark:bg-green-900/40 dark:text-green-400"
                >
                  <IconPackage :size="12" />
                  {{ c.ItemCount }}
                </span>
              </td>
              <td class="px-6 py-3 text-center">
                <span
                  class="inline-flex items-center gap-1 rounded-full bg-orange-100 px-2 py-0.5 text-xs font-medium text-orange-700 dark:bg-orange-900/40 dark:text-orange-400"
                >
                  <IconTool :size="12" />
                  {{ c.WorkOrderCount }}
                </span>
              </td>
              <td class="px-6 py-3">
                <div v-if="c.Customers && c.Customers.length > 0" class="flex flex-wrap gap-1">
                  <RouterLink
                    v-for="customer in c.Customers.slice(0, 3)"
                    :key="customer.id"
                    :to="`/customer/detail/${customer.id}`"
                    class="inline-flex items-center gap-1 rounded-full bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-700 hover:bg-blue-200 dark:bg-blue-900/40 dark:text-blue-400 dark:hover:bg-blue-900/60"
                    @click.stop
                  >
                    <IconUser :size="10" />
                    {{ customer.first_name }} {{ customer.last_name }}
                  </RouterLink>
                  <span v-if="c.Customers.length > 3" class="text-xs text-gray-400">+{{ c.Customers.length - 3 }}</span>
                </div>
                <span v-else class="text-gray-400">—</span>
              </td>
              <td class="px-6 py-3 whitespace-nowrap text-gray-500 dark:text-gray-400">{{ formatDate(c.CreatedAt) }}</td>
              <td class="px-6 py-3 text-right" @click.stop>
                <div class="flex items-center justify-end gap-1">
                  <button
                    v-if="c.ChildCount === 0 && c.ItemCount === 0"
                    class="rounded p-1.5 text-gray-400 hover:bg-gray-100 hover:text-red-500 dark:hover:bg-dk-muted"
                    :title="t('warehouse.delete')"
                    @click="confirmDelete(c.ID, c.Title, $event)"
                  >
                    <IconTrash :size="15" />
                  </button>
                  <button
                    class="rounded p-1.5 text-gray-400 hover:bg-gray-100 hover:text-blue-500 dark:hover:bg-dk-muted"
                    :title="t('warehouse.edit')"
                    @click="openEdit(c.ID, $event)"
                  >
                    <IconEdit :size="15" />
                  </button>
                  <button
                    class="rounded p-1.5 text-gray-400 hover:bg-gray-100 hover:text-blue-500 dark:hover:bg-dk-muted"
                    :title="t('warehouse.view_items')"
                    @click="jumpToContainer(c.ID)"
                  >
                    <IconChevronRight :size="15" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 分页 -->
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
          <button
            @click="goToPage(1)"
            :disabled="currentPage === 1"
            class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted"
          >
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 17l-5-5 5-5M18 17l-5-5 5-5"/></svg>
          </button>
          <button
            @click="goToPage(currentPage - 1)"
            :disabled="currentPage === 1"
            class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted"
          >
            <IconChevronLeft :size="16" />
          </button>
          <button
            v-for="p in pageRange"
            :key="p"
            @click="goToPage(p)"
            :class="['rounded px-2.5 py-1 text-sm', p === currentPage ? 'bg-blue-600 text-white' : 'text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dk-muted']"
          >{{ p }}</button>
          <button
            @click="goToPage(currentPage + 1)"
            :disabled="currentPage === totalPages"
            class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted"
          >
            <IconChevronRight :size="16" />
          </button>
          <button
            @click="goToPage(totalPages)"
            :disabled="currentPage === totalPages"
            class="rounded p-1.5 text-gray-500 hover:bg-gray-100 disabled:opacity-30 dark:hover:bg-dk-muted"
          >
            <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M13 17l5-5-5-5M6 17l5-5-5-5"/></svg>
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- 新增/编辑弹窗 -->
  <Teleport to="body">
    <div
      v-if="showForm"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4"
      @click.self="showForm = false"
    >
      <div class="w-full max-w-lg rounded-xl bg-white p-6 shadow-xl dark:bg-dk-card">
        <h3 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">{{ formTitle }}</h3>
        <div class="space-y-4">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
              {{ t('warehouse.container_name') }} <span class="text-red-500">*</span>
            </label>
            <input
              v-model="form.title"
              type="text"
              :placeholder="t('warehouse.title_placeholder')"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
            />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('warehouse.remark') }}</label>
            <textarea
              v-model="form.remark"
              rows="3"
              :placeholder="t('warehouse.remark_placeholder')"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white resize-none"
            ></textarea>
          </div>
        </div>
        <div class="mt-5 flex justify-end gap-2">
          <button
            @click="showForm = false"
            class="rounded-lg border border-gray-300 px-4 py-2 text-sm dark:border-dk-muted dark:text-white hover:bg-gray-50 dark:hover:bg-dk-muted"
          >
            {{ t('message.cancel') }}
          </button>
          <button
            @click="submitForm"
            :disabled="submitting"
            class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-60"
          >
            {{ submitting ? t('message.submitting') : t('message.submit') }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>

  <!-- 删除确认 -->
  <ConfirmDialog
    v-model="showDeleteConfirm"
    :title="t('warehouse.delete_confirm_title')"
    :message="t('warehouse.delete_confirm_msg', { name: deletingName })"
    :confirm-text="t('warehouse.delete')"
    :cancel-text="t('message.cancel')"
    danger
    @confirm="doDelete"
  />
</template>
