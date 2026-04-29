<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { useUsersStore } from '@/stores/users'
import { warehouseApi } from '@/api/warehouse'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import {
  IconChevronLeft,
  IconChevronRight,
  IconFolder,
  IconPackage,
  IconPlus,
  IconEdit,
  IconTrash,
  IconSearch,
  IconTool,
  IconUser,
} from '@tabler/icons-vue'

usePageTitle('warehouse.container_detail')
const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const toast = useToastStore()
const usersStore = useUsersStore()

const containerId = computed(() => parseInt(route.params.id))

// ── 路由参数变化时重新加载 ──
watch(containerId, async () => {
  subPage.value = 1
  itemPage.value = 1
  await fetchContainer()
  fetchSubContainers()
  fetchItems()
})

// ── 容器详情 ──
const container = ref(null)
const photos = ref([])
const parentChain = ref([])
const containerDepth = ref(0)
const canModifyContainer = ref(false)
const loadingDetail = ref(true)
const notFound = ref(false)

// ── 子容器列表 ──
const subContainers = ref([])
const subTotal = ref(0)
const subPage = ref(1)
const subPageSize = ref(10)
const subSearch = ref('')
const loadingSub = ref(false)

const subTotalPages = computed(() => Math.ceil(subTotal.value / subPageSize.value) || 1)

function subPageRange() {
  const total = subTotalPages.value
  const cur = subPage.value
  let start = Math.max(1, cur - 2)
  let end = Math.min(cur + 4, total)
  if (end - start < 4) start = Math.max(1, end - 4)
  return Array.from({ length: end - start + 1 }, (_, i) => start + i)
}

// ── 物品列表 ──
const items = ref([])
const itemTotal = ref(0)
const itemPage = ref(1)
const itemPageSize = ref(10)
const itemSearch = ref('')
const loadingItems = ref(false)

const itemTotalPages = computed(() => Math.ceil(itemTotal.value / itemPageSize.value) || 1)

function itemPageRange() {
  const total = itemTotalPages.value
  const cur = itemPage.value
  let start = Math.max(1, cur - 2)
  let end = Math.min(cur + 4, total)
  if (end - start < 4) start = Math.max(1, end - 4)
  return Array.from({ length: end - start + 1 }, (_, i) => start + i)
}

// ── 新增子容器弹窗 ──
const showAddSub = ref(false)
const addSubForm = reactive({ title: '', remark: '' })
const submittingSub = ref(false)

// ── 编辑容器弹窗 ──
const showEdit = ref(false)
const editForm = reactive({ title: '', remark: '' })
const submittingEdit = ref(false)

// ── 删除确认 ──
const showDeleteConfirm = ref(false)
const deletingName = ref('')
const deleting = ref(false)

// ── 时间格式化 ──
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

// ── 拉取容器详情 ──
async function fetchContainer() {
  container.value = null
  photos.value = []
  parentChain.value = []
  containerDepth.value = 0
  loadingDetail.value = true
  notFound.value = false
  try {
    const { errCode, data } = await warehouseApi.getContainer(containerId.value)
    if (errCode === 0 && data) {
      container.value = data.container
      photos.value = data.photos ?? []
      parentChain.value = data.parent_chain ?? []
      containerDepth.value = data.depth ?? 0
      canModifyContainer.value = data.canModifyContainer === true
    } else {
      notFound.value = true
    }
  } catch {
    notFound.value = true
  } finally {
    loadingDetail.value = false
  }
}

// ── 拉取子容器 ──
async function fetchSubContainers() {
  subContainers.value = []
  loadingSub.value = true
  try {
    const { errCode, data } = await warehouseApi.getContainers({
      search: subSearch.value,
      parent_id: containerId.value,
      entries: subPageSize.value,
      page: subPage.value,
    })
    if (errCode === 0) {
      subContainers.value = data.containers ?? []
      subTotal.value = data.all_count ?? 0
    }
  } catch {
    // 拦截器已处理
  } finally {
    loadingSub.value = false
  }
}

// ── 拉取物品 ──
async function fetchItems() {
  items.value = []
  loadingItems.value = true
  try {
    const { errCode, data } = await warehouseApi.getItems({
      search: itemSearch.value,
      container_id: containerId.value,
      entries: itemPageSize.value,
      page: itemPage.value,
    })
    if (errCode === 0) {
      items.value = data.items ?? []
      itemTotal.value = data.all_count ?? 0
    }
  } catch {
    // 拦截器已处理
  } finally {
    loadingItems.value = false
  }
}

// ── 新增子容器 ──
async function submitAddSub() {
  if (!addSubForm.title.trim()) {
    toast.error(t('warehouse.title_required'))
    return
  }
  submittingSub.value = true
  try {
    const { errCode } = await warehouseApi.addContainer({
      parent_id: containerId.value,
      title: addSubForm.title.trim(),
      remark: addSubForm.remark.trim(),
    })
    if (errCode === 0) {
      toast.success(t('message.save_success'))
      showAddSub.value = false
      Object.assign(addSubForm, { title: '', remark: '' })
      fetchSubContainers()
      fetchContainer()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    toast.error(t('message.server_error'))
  } finally {
    submittingSub.value = false
  }
}

// ── 编辑容器 ──
function openEdit() {
  editForm.title = container.value?.Title ?? ''
  editForm.remark = container.value?.Remark ?? ''
  showEdit.value = true
}

async function submitEdit() {
  if (!editForm.title.trim()) {
    toast.error(t('warehouse.title_required'))
    return
  }
  submittingEdit.value = true
  try {
    const { errCode } = await warehouseApi.updateContainer({
      id: containerId.value,
      title: editForm.title.trim(),
      remark: editForm.remark.trim(),
    })
    if (errCode === 0) {
      toast.success(t('message.save_success'))
      showEdit.value = false
      fetchContainer()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    toast.error(t('message.server_error'))
  } finally {
    submittingEdit.value = false
  }
}

// ── 删除容器 ──
function confirmDelete() {
  deletingName.value = container.value?.Title ?? ''
  showDeleteConfirm.value = true
}

async function doDelete() {
  deleting.value = true
  try {
    const { errCode } = await warehouseApi.deleteContainer(containerId.value)
    if (errCode === 0) {
      toast.success(t('message.delete_ok'))
      showDeleteConfirm.value = false
      if (container.value.ParentID) {
        router.push(`/warehouse/container/${container.value.ParentID}`)
      } else {
        router.push('/warehouse/container')
      }
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    toast.error(t('message.server_error'))
  } finally {
    deleting.value = false
  }
}

// ── 初始化 ──
onMounted(async () => {
  await fetchContainer()
  fetchSubContainers()
  fetchItems()
})
</script>

<template>
  <div class="p-4 max-w-7xl mx-auto space-y-4">

    <!-- 加载 -->
    <div v-if="loadingDetail" class="flex justify-center py-16">
      <svg class="h-6 w-6 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
      </svg>
    </div>

    <!-- 未找到 -->
    <div v-else-if="notFound" class="flex flex-col items-center justify-center py-16 text-gray-400">
      <p>{{ t('message.not_found') }}</p>
      <button
        class="mt-2 text-sm text-blue-500 hover:underline"
        @click="router.push('/warehouse/container')"
      >
        {{ t('warehouse.back_to_list') }}
      </button>
    </div>

    <template v-else-if="container">

      <!-- 面包屑 + 操作栏 -->
      <div class="flex items-center justify-between gap-3">
        <div class="flex items-center gap-1 text-sm flex-wrap">
          <RouterLink to="/warehouse/container" class="text-blue-500 hover:underline">
            {{ t('warehouse.container_list') }}
          </RouterLink>
          <template v-for="p in parentChain" :key="p.id">
            <span class="text-gray-400">/</span>
            <RouterLink
              :to="`/warehouse/container/${p.id}`"
              class="text-blue-500 hover:underline"
            >{{ p.title }}</RouterLink>
          </template>
          <span class="text-gray-400">/</span>
          <span class="font-medium text-gray-900 dark:text-white">{{ container.Title }}</span>
        </div>
        <div class="flex gap-2">
          <button
            v-if="canModifyContainer"
            class="inline-flex items-center gap-1.5 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:hover:bg-dk-muted"
            @click="openEdit"
          >
            <IconEdit :size="14" />
            {{ t('warehouse.edit') }}
          </button>
          <button
            v-if="canModifyContainer && container.ChildCount === 0 && container.ItemCount === 0"
            class="inline-flex items-center gap-1.5 rounded-lg border border-red-300 bg-white px-3 py-1.5 text-sm font-medium text-red-600 transition-colors hover:bg-red-50 dark:border-red-900 dark:bg-dk-base dark:text-red-400 dark:hover:bg-red-900/20"
            @click="confirmDelete"
          >
            <IconTrash :size="14" />
            {{ t('warehouse.delete') }}
          </button>
        </div>
      </div>

      <!-- 容器信息卡片 -->
      <div class="rounded-xl border border-gray-200 bg-white px-5 py-4 shadow dark:border-dk-muted dark:bg-dk-card">
        <div class="flex items-center gap-2 mb-2">
          <IconFolder class="text-blue-500 flex-shrink-0" :size="20" />
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ container.Title }}</h2>
        </div>

        <!-- 图片 -->
        <div v-if="photos.length" class="flex gap-2 flex-wrap mb-3">
          <a
            v-for="photo in photos"
            :key="photo.ID"
            :href="`/api/files/get/${photo.Sha256}`"
            target="_blank"
          >
            <img
              :src="`/api/files/get/${photo.Sha256}`"
              class="w-16 h-16 object-cover rounded border border-gray-200 dark:border-dk-muted"
              :alt="photo.Name"
            />
          </a>
        </div>

        <!-- 备注 -->
        <p
          v-if="container.Remark"
          class="text-sm text-gray-500 dark:text-gray-400 whitespace-pre-wrap mb-3"
        >
          {{ container.Remark }}
        </p>

        <!-- 元信息 -->
        <div class="flex flex-wrap gap-x-6 gap-y-1 text-xs text-gray-400 dark:text-gray-500">
          <span class="flex items-center gap-1">
            <span>{{ t('warehouse.created_by') }}:</span>
            <img
              :src="usersStore.getAvatarUrlFromUserID(container.CreatorID)"
              class="w-4 h-4 rounded-full object-cover"
            />
            {{ usersStore.getUsernameFromUserID(container.CreatorID) }}
          </span>
          <span>{{ t('warehouse.created_at') }}: {{ fmtTs(container.CreatedAt) }}</span>
          <span v-if="container.UpdatedAt">{{ t('warehouse.updated_at') }}: {{ fmtTs(container.UpdatedAt) }}</span>
          <span>{{ t('warehouse.child_containers') }}: {{ container.ChildCount }}</span>
          <span>{{ t('warehouse.items') }}: {{ container.ItemCount }}</span>
        </div>
      </div>

      <!-- 子容器列表 -->
      <div class="rounded-xl border border-gray-200 bg-white shadow dark:border-dk-muted dark:bg-dk-card">
        <!-- 工具栏 -->
        <div class="flex items-center gap-2 px-5 py-3 border-b border-gray-100 dark:border-dk-muted">
          <div class="relative flex-1 max-w-xs">
            <IconSearch class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" :size="15" />
            <input
              v-model="subSearch"
              type="text"
              :placeholder="t('warehouse.search_placeholder')"
              class="w-full rounded-lg border border-gray-300 bg-white py-1.5 pl-9 pr-3 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
              @keyup.enter="subPage = 1; fetchSubContainers()"
            />
          </div>
          <button
            v-if="containerDepth < 4"
            class="inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700"
            @click="showAddSub = true"
          >
            <IconPlus :size="15" />
            {{ t('warehouse.add_container') }}
          </button>
        </div>

        <!-- 表格 -->
        <div class="overflow-x-auto">
          <table class="w-full text-left text-sm text-gray-900 dark:text-white">
            <thead>
              <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base dark:text-gray-400">
                <th class="px-5 py-3 font-medium">{{ t('warehouse.container_name') }}</th>
                <th class="px-5 py-3 font-medium w-24 text-center">{{ t('warehouse.child_containers') }}</th>
                <th class="px-5 py-3 font-medium w-24 text-center">{{ t('warehouse.items') }}</th>
                <th class="px-5 py-3 font-medium whitespace-nowrap">{{ t('warehouse.created_at') }}</th>
                <th class="px-5 py-3 font-medium whitespace-nowrap">{{ t('warehouse.updated_at') }}</th>
                <th class="px-5 py-3 font-medium">{{ t('warehouse.created_by') }}</th>
                <th class="px-5 py-3 font-medium w-24 text-right">{{ t('warehouse.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="loadingSub">
                <td colspan="7" class="px-5 py-8 text-center">
                  <svg class="mx-auto h-5 w-5 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                  </svg>
                </td>
              </tr>
              <tr v-else-if="subContainers.length === 0">
                <td colspan="7" class="px-5 py-8 text-center text-gray-400 dark:text-gray-500">
                  {{ t('warehouse.no_containers') }}
                </td>
              </tr>
              <tr
                v-else
                v-for="c in subContainers"
                :key="c.ID"
                class="border-b border-gray-100 cursor-pointer transition-colors hover:bg-gray-50 dark:border-dk-muted dark:hover:bg-dk-base"
                @click="router.push(`/warehouse/container/${c.ID}`)"
              >
                <td class="px-5 py-3">
                  <div class="flex items-center gap-2">
                    <IconFolder class="text-blue-500 flex-shrink-0" :size="16" />
                    <span class="font-medium max-w-xs truncate">{{ c.Title }}</span>
                  </div>
                </td>
                <td class="px-5 py-3 text-center">
                  <span class="inline-flex items-center gap-1 rounded-full bg-purple-100 px-2 py-0.5 text-xs font-medium text-purple-700 dark:bg-purple-900/40 dark:text-purple-400">
                    {{ c.ChildCount }}
                  </span>
                </td>
                <td class="px-5 py-3 text-center">
                  <span class="inline-flex items-center gap-1 rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-700 dark:bg-green-900/40 dark:text-green-400">
                    {{ c.ItemCount }}
                  </span>
                </td>
                <td class="px-5 py-3 text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ fmtTs(c.CreatedAt) }}</td>
                <td class="px-5 py-3 text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ fmtTs(c.UpdatedAt) }}</td>
                <td class="px-5 py-3">
                  <div class="flex items-center gap-1.5">
                    <img
                      :src="usersStore.getAvatarUrlFromUserID(c.CreatorID)"
                      class="w-5 h-5 rounded-full object-cover flex-shrink-0"
                    />
                    <span class="truncate text-gray-600 dark:text-gray-400">{{ usersStore.getUsernameFromUserID(c.CreatorID) }}</span>
                  </div>
                </td>
                <td class="px-5 py-3 text-right">
                  <button
                    class="text-xs text-blue-500 hover:underline"
                    @click.stop="router.push(`/warehouse/container/${c.ID}`)"
                  >
                    {{ t('warehouse.view_items') }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- 分页 -->
        <div v-if="subTotalPages > 1" class="flex items-center justify-center gap-1 px-5 py-3 border-t border-gray-100 dark:border-dk-muted">
          <button
            class="flex items-center justify-center w-7 h-7 rounded text-sm text-gray-500 hover:bg-gray-100 disabled:opacity-30 disabled:cursor-not-allowed dark:text-gray-400 dark:hover:bg-dk-muted"
            :disabled="subPage <= 1"
            @click="subPage--; fetchSubContainers()"
          >
            <IconChevronLeft :size="15" />
          </button>
          <button
            v-for="p in subPageRange()"
            :key="p"
            class="flex items-center justify-center w-7 h-7 rounded text-sm"
            :class="p === subPage
              ? 'bg-blue-600 text-white font-medium'
              : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dk-muted'"
            @click="subPage = p; fetchSubContainers()"
          >
            {{ p }}
          </button>
          <button
            class="flex items-center justify-center w-7 h-7 rounded text-sm text-gray-500 hover:bg-gray-100 disabled:opacity-30 disabled:cursor-not-allowed dark:text-gray-400 dark:hover:bg-dk-muted"
            :disabled="subPage >= subTotalPages"
            @click="subPage++; fetchSubContainers()"
          >
            <IconChevronRight :size="15" />
          </button>
        </div>
      </div>

      <!-- 物品列表 -->
      <div class="rounded-xl border border-gray-200 bg-white shadow dark:border-dk-muted dark:bg-dk-card">
        <!-- 工具栏 -->
        <div class="flex items-center gap-2 px-5 py-3 border-b border-gray-100 dark:border-dk-muted">
          <div class="relative flex-1 max-w-xs">
            <IconSearch class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" :size="15" />
            <input
              v-model="itemSearch"
              type="text"
              :placeholder="t('warehouse.search_placeholder')"
              class="w-full rounded-lg border border-gray-300 bg-white py-1.5 pl-9 pr-3 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
              @keyup.enter="itemPage = 1; fetchItems()"
            />
          </div>
          <button
            class="inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700"
            @click="router.push(`/warehouse/container/${containerId}/add-item`)"
          >
            <IconPlus :size="15" />
            {{ t('warehouse.add_item') }}
          </button>
        </div>

        <!-- 表格 -->
        <div class="overflow-x-auto">
          <table class="w-full text-left text-sm text-gray-900 dark:text-white">
            <thead>
              <tr class="border-b border-gray-200 bg-gray-50 text-gray-500 dark:border-dk-muted dark:bg-dk-base dark:text-gray-400">
                <th class="px-5 py-3 font-medium">{{ t('warehouse.item_name') }}</th>
                <th class="px-5 py-3 font-medium">{{ t('warehouse.serial_number') }}</th>
                <th class="px-5 py-3 font-medium">{{ t('warehouse.remark') }}</th>
                <th class="px-5 py-3 font-medium w-20 text-center">{{ t('warehouse.quantity') }}</th>
                <th class="px-5 py-3 font-medium w-24 text-center">{{ t('work_order.work_order_count') }}</th>
                <th class="px-5 py-3 font-medium">{{ t('customer.related_customers') }}</th>
                <th class="px-5 py-3 font-medium whitespace-nowrap">{{ t('warehouse.created_at') }}</th>
                <th class="px-5 py-3 font-medium whitespace-nowrap">{{ t('warehouse.updated_at') }}</th>
                <th class="px-5 py-3 font-medium">{{ t('warehouse.created_by') }}</th>
                <th class="px-5 py-3 font-medium w-20 text-right">{{ t('warehouse.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="loadingItems">
                <td colspan="10" class="px-5 py-8 text-center">
                  <svg class="mx-auto h-5 w-5 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                  </svg>
                </td>
              </tr>
              <tr v-else-if="items.length === 0">
                <td colspan="10" class="px-5 py-8 text-center text-gray-400 dark:text-gray-500">
                  {{ t('warehouse.no_items') }}
                </td>
              </tr>
              <tr
                v-else
                v-for="item in items"
                :key="item.ID"
                class="border-b border-gray-100 cursor-pointer transition-colors hover:bg-gray-50 dark:border-dk-muted dark:hover:bg-dk-base"
                @click="router.push(`/warehouse/item/${item.ID}`)"
              >
                <td class="px-5 py-3">
                  <div class="flex items-center gap-2">
                    <IconPackage class="text-green-500 flex-shrink-0" :size="16" />
                    <span class="font-medium max-w-xs truncate">{{ item.Name }}</span>
                  </div>
                </td>
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
                      class="inline-flex items-center gap-1 rounded-full bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-700 hover:bg-blue-200 dark:bg-blue-900/40 dark:text-blue-400 dark:hover:bg-blue-900/60"
                      @click.stop
                    >
                      <IconUser :size="10" />
                      {{ customer.first_name }} {{ customer.last_name }}
                    </RouterLink>
                    <span v-if="item.Customers.length > 3" class="text-xs text-gray-400">+{{ item.Customers.length - 3 }}</span>
                  </div>
                  <span v-else class="text-gray-400">—</span>
                </td>
                <td class="px-5 py-3 text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ fmtTs(item.CreatedAt) }}</td>
                <td class="px-5 py-3 text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ fmtTs(item.UpdatedAt) }}</td>
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
                    @click.stop="router.push(`/warehouse/item/${item.ID}`)"
                  >
                    {{ t('warehouse.view_items') }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- 分页 -->
        <div v-if="itemTotalPages > 1" class="flex items-center justify-center gap-1 px-5 py-3 border-t border-gray-100 dark:border-dk-muted">
          <button
            class="flex items-center justify-center w-7 h-7 rounded text-sm text-gray-500 hover:bg-gray-100 disabled:opacity-30 disabled:cursor-not-allowed dark:text-gray-400 dark:hover:bg-dk-muted"
            :disabled="itemPage <= 1"
            @click="itemPage--; fetchItems()"
          >
            <IconChevronLeft :size="15" />
          </button>
          <button
            v-for="p in itemPageRange()"
            :key="p"
            class="flex items-center justify-center w-7 h-7 rounded text-sm"
            :class="p === itemPage
              ? 'bg-blue-600 text-white font-medium'
              : 'text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dk-muted'"
            @click="itemPage = p; fetchItems()"
          >
            {{ p }}
          </button>
          <button
            class="flex items-center justify-center w-7 h-7 rounded text-sm text-gray-500 hover:bg-gray-100 disabled:opacity-30 disabled:cursor-not-allowed dark:text-gray-400 dark:hover:bg-dk-muted"
            :disabled="itemPage >= itemTotalPages"
            @click="itemPage++; fetchItems()"
          >
            <IconChevronRight :size="15" />
          </button>
        </div>
      </div>

    </template>
  </div>

  <!-- 新增子容器弹窗 -->
  <Transition name="fade">
    <div
      v-if="showAddSub"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/40"
      @click.self="showAddSub = false"
    >
      <div class="w-full max-w-md rounded-xl border border-gray-200 bg-white p-5 shadow-xl dark:border-dk-muted dark:bg-dk-card">
        <h3 class="mb-4 text-base font-semibold text-gray-900 dark:text-white">{{ t('warehouse.add_container') }}</h3>
        <div class="space-y-3">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('warehouse.container_name') }} *</label>
            <input
              v-model="addSubForm.title"
              type="text"
              :placeholder="t('warehouse.title_placeholder')"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
              @keyup.enter="submitAddSub"
            />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('warehouse.remark') }}</label>
            <textarea
              v-model="addSubForm.remark"
              :placeholder="t('warehouse.remark_placeholder')"
              rows="3"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
            ></textarea>
          </div>
        </div>
        <div class="mt-4 flex justify-end gap-2">
          <button
            class="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white hover:bg-gray-50 dark:hover:bg-dk-muted"
            @click="showAddSub = false"
          >
            {{ t('message.cancel') }}
          </button>
          <button
            class="inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            :disabled="submittingSub"
            @click="submitAddSub"
          >
            <svg v-if="submittingSub" class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
            </svg>
            {{ t('message.save') }}
          </button>
        </div>
      </div>
    </div>
  </Transition>

  <!-- 编辑容器弹窗 -->
  <Transition name="fade">
    <div
      v-if="showEdit"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/40"
      @click.self="showEdit = false"
    >
      <div class="w-full max-w-md rounded-xl border border-gray-200 bg-white p-5 shadow-xl dark:border-dk-muted dark:bg-dk-card">
        <h3 class="mb-4 text-base font-semibold text-gray-900 dark:text-white">{{ t('warehouse.edit_container') }}</h3>
        <div class="space-y-3">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('warehouse.container_name') }} *</label>
            <input
              v-model="editForm.title"
              type="text"
              :placeholder="t('warehouse.title_placeholder')"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
              @keyup.enter="submitEdit"
            />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('warehouse.remark') }}</label>
            <textarea
              v-model="editForm.remark"
              :placeholder="t('warehouse.remark_placeholder')"
              rows="3"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
            ></textarea>
          </div>
        </div>
        <div class="mt-4 flex justify-end gap-2">
          <button
            class="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white hover:bg-gray-50 dark:hover:bg-dk-muted"
            @click="showEdit = false"
          >
            {{ t('message.cancel') }}
          </button>
          <button
            class="inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            :disabled="submittingEdit"
            @click="submitEdit"
          >
            <svg v-if="submittingEdit" class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
            </svg>
            {{ t('message.save') }}
          </button>
        </div>
      </div>
    </div>
  </Transition>

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
