<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { useUserStore } from '@/stores/user'
import { useUsersStore } from '@/stores/users'
import { workOrderApi } from '@/api/work_order'
import useDropzone from '@/components/useDropzone.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import {
  IconChevronLeft,
  IconCheck,
  IconLoader2,
  IconTrash,
  IconX,
  IconSearch,
  IconExternalLink,
} from '@tabler/icons-vue'

usePageTitle('work_order.detail_title')
const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const toast = useToastStore()
const userStore = useUserStore()
const usersStore = useUsersStore()

const orderId = computed(() => parseInt(route.params.id))

const order = ref(null)
const photos = ref([])
const commits = ref([])
const canModify = ref(false)
const canCommit = ref(false)
const loading = ref(true)
const notFound = ref(false)

// 提交进度相关
const submittingCommit = ref(false)
const commitStatus = ref('pending')
const commitComment = ref('')
const commitDropzoneRef = ref(null)

// 采购订单关联相关
const purchaseSearchQuery = ref('')
const purchaseSearchResults = ref([])
const selectedPurchaseOrders = ref([])
const purchaseSearchLoading = ref(false)
const purchaseDropdownVisible = ref(false)
let purchaseSearchTimer = null
const purchaseDropdownRef = ref(null)

// 是否可以提交（订单、备注、上传图片都为空时才禁止）
const canSubmit = computed(() => {
  const hasSelectedOrders = selectedPurchaseOrders.value.length > 0
  const hasComment = !!commitComment.value
  const hasPhotos = commitDropzoneRef.value?.return_files().filter(f => f.is_upload).length > 0
  // 订单、备注、上传图片都为空时才禁止提交
  return hasSelectedOrders || hasComment || hasPhotos
})

// 所有 commits 中关联的采购订单（去重）
const allPurchaseOrders = computed(() => {
  const map = new Map()
  for (const commit of commits.value) {
    if (commit.purchaseOrders) {
      for (const po of commit.purchaseOrders) {
        if (!map.has(po.id)) {
          map.set(po.id, po)
        }
      }
    }
  }
  return [...map.values()]
})

// 点击外部关闭下拉框
function onDocumentClick(e) {
  if (purchaseDropdownRef.value && !purchaseDropdownRef.value.contains(e.target)) {
    purchaseDropdownVisible.value = false
  }
}

// 状态选项
const statusOptions = [
  { value: 'pending',       labelKey: 'work_order.status_pending',       color: 'yellow' },
  { value: 'checked',       labelKey: 'work_order.status_checked',       color: 'blue' },
  { value: 'parts_ordered', labelKey: 'work_order.status_parts_ordered', color: 'purple' },
  { value: 'repaired',      labelKey: 'work_order.status_repaired',      color: 'green' },
  { value: 'returned',      labelKey: 'work_order.status_returned',      color: 'gray' },
  { value: 'unrepairable',  labelKey: 'work_order.status_unrepairable',  color: 'red' },
]

const statusColorMap = {
  pending:       'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/40 dark:text-yellow-400',
  checked:       'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-400',
  parts_ordered: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-400',
  repaired:      'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400',
  returned:      'bg-gray-200 text-gray-600 dark:bg-gray-700 dark:text-gray-300',
  unrepairable:  'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400',
}

function getStatusLabel(status) {
  const opt = statusOptions.find((o) => o.value === status)
  return opt ? t(opt.labelKey) : status
}

function getStatusColorClass(status) {
  return statusColorMap[status] || 'bg-gray-100 text-gray-600'
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  return new Intl.DateTimeFormat(locale.value, {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false,
  }).format(d)
}

function getPhotoUrl(file) {
  return `/api/files/get/${file.Sha256}`
}

// ==================== 加载工单 ====================
async function fetchOrder() {
  loading.value = true
  try {
    const { errCode, data } = await workOrderApi.get(orderId.value)
    if (errCode === 0 && data) {
      order.value = data.order ?? null
      canModify.value = data.canModify ?? false
      canCommit.value = data.canCommit ?? false
      photos.value = data.photos ?? []
      commits.value = data.commits ?? []
      // 初始化进度提交状态为当前状态
      if (order.value?.CurrentStatus) {
        commitStatus.value = order.value.CurrentStatus
      }
    } else {
      notFound.value = true
    }
  } catch {
    notFound.value = true
  } finally {
    loading.value = false
  }
}

// ==================== 提交进度 ====================
async function handleCommit() {
  if (submittingCommit.value) return
  submittingCommit.value = true
  try {
    const purchaseOrderIds = selectedPurchaseOrders.value.map(p => p.id)
    // 从 dropzone 获取已上传的文件 hash
    const uploadedPhotos = commitDropzoneRef.value?.return_files()
      .filter(f => f.is_upload)
      .map(f => f.hash) ?? []
    const { errCode } = await workOrderApi.commit(
      orderId.value,
      commitStatus.value,
      commitComment.value,
      uploadedPhotos,
      purchaseOrderIds,
    )
    if (errCode === 0) {
      toast.success(t('message.save_ok'))
      commitComment.value = ''
      selectedPurchaseOrders.value = []
      // 清空 dropzone（刷新组件即可）
      await fetchOrder()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    toast.error(t('message.server_error'))
  } finally {
    submittingCommit.value = false
  }
}

// ==================== 删除进度 ====================
const showDeleteCommitConfirm = ref(false)
const pendingDeleteCommitId = ref(null)

function handleDeleteCommit(commitId) {
  pendingDeleteCommitId.value = commitId
  showDeleteCommitConfirm.value = true
}

async function confirmDeleteCommit() {
  if (!pendingDeleteCommitId.value) return
  try {
    const { errCode } = await workOrderApi.deleteCommit(orderId.value, pendingDeleteCommitId.value)
    if (errCode === 0) {
      toast.success(t('message.delete_ok'))
      // 前端直接移除该 commit，保持滚动位置
      commits.value = commits.value.filter(c => c.ID !== pendingDeleteCommitId.value)
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    toast.error(t('message.server_error'))
  } finally {
    pendingDeleteCommitId.value = null
    showDeleteCommitConfirm.value = false
  }
}

// 判断是否可以删除进度
function canDeleteCommit(commit, index) {
  // 最新状态（第0条）不显示删除按钮
  if (index === 0) return false
  // 订单创建者
  if (order.value?.UserID === userStore.user?.ID) return true
  // 进度创建者
  if (commit.UserID === userStore.user?.ID) return true
  // 管理员
  if (userStore.user?.Type === 'admin') return true
  return false
}

// ==================== 快捷切换状态 ====================
async function quickChangeStatus(newStatus) {
  if (newStatus === order.value?.CurrentStatus) return
  submittingCommit.value = true
  try {
    const { errCode } = await workOrderApi.commit(orderId.value, newStatus, '')
    if (errCode === 0) {
      toast.success(t('message.save_ok'))
      await fetchOrder()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    toast.error(t('message.server_error'))
  } finally {
    submittingCommit.value = false
  }
}

// ==================== 采购订单搜索 ====================
async function searchPurchaseOrders() {
  purchaseSearchLoading.value = true
  try {
    const { errCode, data } = await workOrderApi.searchPurchaseOrders(purchaseSearchQuery.value, 10)
    if (errCode === 0) {
      purchaseSearchResults.value = data.orders || []
    }
  } catch {
    purchaseSearchResults.value = []
  } finally {
    purchaseSearchLoading.value = false
  }
}

function onPurchaseSearchInput() {
  purchaseDropdownVisible.value = true
  clearTimeout(purchaseSearchTimer)
  purchaseSearchTimer = setTimeout(() => {
    searchPurchaseOrders()
  }, 300)
}

function selectPurchaseOrder(po) {
  // 检查是否已选中
  if (!selectedPurchaseOrders.value.find(p => p.id === po.id)) {
    selectedPurchaseOrders.value.push(po)
  }
  // 清空搜索框并重新搜索，保持下拉框显示
  purchaseSearchQuery.value = ''
  searchPurchaseOrders()
}

function removePurchaseOrder(poId) {
  selectedPurchaseOrders.value = selectedPurchaseOrders.value.filter(p => p.id !== poId)
}

function goToPurchaseOrder(poId) {
  router.push(`/purchase/showorder/${poId}`)
}

// 采购订单状态颜色
function getPurchaseStatusClass(status) {
  const map = {
    pending: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/40 dark:text-yellow-400',
    ordered: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-400',
    arrived: 'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-400',
    received: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400',
    lost: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400',
    returned: 'bg-gray-200 text-gray-600 dark:bg-gray-700 dark:text-gray-300',
  }
  return map[status] || 'bg-gray-100 text-gray-600'
}

const purchaseStatusLabels = {
  pending: '待处理',
  ordered: '已下单',
  arrived: '已到达',
  received: '已收件',
  lost: '丢件',
  returned: '退件',
}

function getPurchaseStatusLabel(status) {
  return purchaseStatusLabels[status] || status
}

onMounted(() => {
  fetchOrder()
  document.addEventListener('click', onDocumentClick)
})

onUnmounted(() => {
  document.removeEventListener('click', onDocumentClick)
})
</script>

<template>
  <div class="mx-auto max-w-6xl px-6 py-6">
    <!-- 顶部操作栏 -->
    <div class="mb-4 flex items-center justify-between">
      <button
        @click="$router.back()"
        class="inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-dk-card dark:hover:text-gray-200"
      >
        <IconChevronLeft :size="16" />
        {{ t('work_order.back_to_list') }}
      </button>
      <RouterLink
        v-if="canModify && order"
        :to="`/work_order/edit/${order.ID}`"
        class="inline-flex items-center gap-1.5 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-card dark:text-gray-300 dark:hover:bg-dk-base"
      >
        <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
        </svg>
        {{ t('work_order.edit') }}
      </RouterLink>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-24">
      <svg class="h-8 w-8 animate-spin text-blue-500" viewBox="0 0 24 24" fill="none">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
      </svg>
      <span class="ml-3 text-gray-500">{{ t('message.loading') }}</span>
    </div>

    <!-- Not Found -->
    <div
      v-else-if="notFound"
      class="rounded-xl border border-gray-200 bg-white py-16 text-center shadow-lg dark:border-dk-muted dark:bg-dk-card"
    >
      <p class="text-gray-400">{{ t('work_order.not_found') }}</p>
    </div>

    <!-- 工单详情 -->
    <div v-else class="flex flex-col gap-6">

      <!-- ===== 主信息卡 ===== -->
      <div class="rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
        <!-- 卡头：标题 + 状态 + 创建者 + 时间 -->
        <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
          <div class="flex flex-wrap items-center gap-3">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('work_order.detail_title') }} #{{ orderId }}
            </h2>
            <span
              v-if="order?.CurrentStatus"
              class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-semibold"
              :class="getStatusColorClass(order.CurrentStatus)"
            >
              <IconLoader2 v-if="submittingCommit" :size="10" class="animate-spin" />
              {{ getStatusLabel(order.CurrentStatus) }}
            </span>
            <span
              v-if="order?.UserID"
              class="flex items-center gap-1.5 rounded-full border border-gray-200 bg-gray-50 px-2 py-0.5 text-xs text-gray-500 dark:border-dk-muted dark:bg-dk-base dark:text-gray-400"
            >
              <img
                :src="usersStore.getAvatarUrlFromUserID(order.UserID)"
                class="h-4 w-4 rounded-full object-cover"
              />
              {{ usersStore.getUsernameFromUserID(order.UserID) }}
            </span>
          </div>
          <span class="text-sm text-gray-400">{{ formatDate(order?.CreatedAt) }}</span>
        </div>

        <!-- 状态快捷切换（所有登录用户可见） -->
        <div
          v-if="canCommit"
          class="flex flex-wrap items-center gap-2 border-b border-gray-100 px-6 py-3 dark:border-dk-muted"
        >
          <span class="text-sm text-gray-500 dark:text-gray-400">{{ t('purchase.change_status') }}:</span>
          <button
            v-for="opt in statusOptions"
            :key="opt.value"
            class="inline-flex items-center gap-1 rounded-full border px-3 py-1 text-xs font-medium transition-all"
            :class="order?.CurrentStatus === opt.value
              ? [getStatusColorClass(opt.value), 'border-transparent']
              : 'border-gray-200 text-gray-500 hover:border-gray-300 hover:bg-gray-50 dark:border-dk-muted dark:text-gray-400 dark:hover:bg-dk-base'"
            :disabled="submittingCommit"
            @click="quickChangeStatus(opt.value)"
          >
            <IconCheck v-if="order?.CurrentStatus === opt.value" :size="12" />
            {{ t(opt.labelKey) }}
          </button>
        </div>

        <!-- 工单基本信息 -->
        <div class="space-y-4 px-6 py-5">
          <!-- 标题 -->
          <div>
            <label class="mb-1 block text-xs font-medium text-gray-400">{{ t('work_order.title') }}</label>
            <p class="font-medium text-gray-900 dark:text-white">{{ order?.Title || '-' }}</p>
          </div>
          <!-- 描述 -->
          <div v-if="order?.Description">
            <label class="mb-1 block text-xs font-medium text-gray-400">{{ t('work_order.description') }}</label>
            <p class="whitespace-pre-wrap text-sm text-gray-700 dark:text-gray-300">{{ order.Description }}</p>
          </div>
          <!-- 关联采购订单汇总（去重） -->
          <div v-if="allPurchaseOrders.length > 0">
            <label class="mb-1 block text-xs font-medium text-gray-400">关联采购订单</label>
            <div class="flex flex-wrap gap-2">
              <RouterLink
                v-for="po in allPurchaseOrders"
                :key="po.id"
                :to="`/purchase/showorder/${po.id}`"
                class="inline-flex items-center gap-1 rounded-full border border-blue-200 bg-blue-50 px-2.5 py-1 text-xs font-medium text-blue-700 transition-colors hover:bg-blue-100 dark:border-blue-800 dark:bg-blue-900/30 dark:text-blue-300 dark:hover:bg-blue-900/50"
              >
                #{{ po.id }} {{ po.title || '' }}
                <span
                  class="ml-1 rounded px-1.5 py-0.5 text-[10px]"
                  :class="getPurchaseStatusClass(po.status)"
                >
                  {{ getPurchaseStatusLabel(po.status) }}
                </span>
              </RouterLink>
            </div>
          </div>
        </div>

        <!-- 图片区 -->
        <div class="border-t border-gray-100 px-6 py-5 dark:border-dk-muted">
          <h4 class="mb-3 text-sm font-semibold text-gray-500 dark:text-gray-400">{{ t('work_order.photos') }}</h4>
          <div v-if="photos.length === 0" class="text-sm text-gray-400">{{ t('work_order.no_photos') }}</div>
          <div v-else class="flex flex-wrap gap-3">
            <a
              v-for="file in photos"
              :key="file.ID"
              :href="getPhotoUrl(file)"
              target="_blank"
              class="group relative block h-24 w-24 overflow-hidden rounded-lg border border-gray-200 bg-gray-50 transition hover:border-blue-400 dark:border-dk-muted dark:bg-dk-base"
            >
              <img
                :src="getPhotoUrl(file)"
                :alt="file.Name"
                class="h-full w-full object-cover transition group-hover:scale-105"
              />
            </a>
          </div>
        </div>
      </div>

      <!-- ===== 工作进度时间线 ===== -->
      <div class="rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
        <div class="border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
          <h3 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('work_order.commit_history') }}</h3>
        </div>

        <!-- 新增进度表单（所有登录用户可见） -->
        <div v-if="canCommit" class="border-t border-gray-100 px-6 py-5 dark:border-dk-muted">
          <h4 class="mb-3 text-sm font-semibold text-gray-700 dark:text-gray-300">{{ t('work_order.add_commit') }}</h4>

          <!-- 第一行：进度状态、关联采购订单 -->
          <div class="mb-3 flex flex-wrap items-start gap-3">
            <!-- 状态选择 -->
            <div class="min-w-40">
              <label class="mb-1 block text-xs font-medium text-gray-400">{{ t('work_order.commit_status_label') }}</label>
              <select
                v-model="commitStatus"
                class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
              >
                <option v-for="opt in statusOptions" :key="opt.value" :value="opt.value">
                  {{ t(opt.labelKey) }}
                </option>
              </select>
            </div>

            <!-- 采购订单关联（仅在已下单状态显示） -->
            <div v-if="commitStatus === 'parts_ordered'" class="flex-1">
              <label class="mb-1 block text-xs font-medium text-gray-400">关联采购订单</label>
              <!-- 已选中的订单 -->
              <div v-if="selectedPurchaseOrders.length > 0" class="mb-2 flex flex-wrap gap-2">
                <div
                  v-for="po in selectedPurchaseOrders"
                  :key="po.id"
                  class="inline-flex items-center gap-1 rounded-full border border-blue-200 bg-blue-50 px-2.5 py-1 text-xs font-medium text-blue-700 dark:border-blue-800 dark:bg-blue-900/30 dark:text-blue-300"
                >
                  #{{ po.id }} {{ po.title || '' }}
                  <button
                    class="ml-0.5 rounded-full p-0.5 transition-colors hover:bg-blue-200 dark:hover:bg-blue-800"
                    @click="removePurchaseOrder(po.id)"
                  >
                    <IconX :size="12" />
                  </button>
                </div>
              </div>
              <!-- 搜索框 -->
              <div ref="purchaseDropdownRef" class="relative">
                <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                  <IconSearch :size="14" class="text-gray-400" />
                </div>
                <input
                  v-model="purchaseSearchQuery"
                  type="text"
                  placeholder="搜索采购订单..."
                  class="w-full rounded-lg border border-gray-300 bg-white py-2 pl-9 pr-3 text-sm outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white"
                  @input="onPurchaseSearchInput"
                  @focus="purchaseDropdownVisible = true; searchPurchaseOrders()"
                />
                <!-- 搜索结果下拉 -->
                <div
                  v-if="purchaseDropdownVisible && purchaseSearchResults.length > 0"
                  class="absolute z-10 mt-1 max-h-48 w-full overflow-auto rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-dk-muted dark:bg-dk-card"
                >
                  <button
                    v-for="po in purchaseSearchResults"
                    :key="po.id"
                    class="flex w-full items-center justify-between px-3 py-2 text-left text-sm hover:bg-gray-50 dark:hover:bg-dk-muted"
                    @mousedown.prevent="selectPurchaseOrder(po)"
                  >
                    <span class="flex items-center gap-2">
                      <span class="font-medium text-gray-900 dark:text-white">#{{ po.id }}</span>
                      <span class="text-gray-600 dark:text-gray-300">{{ po.title || '-' }}</span>
                    </span>
                    <span
                      class="rounded-full px-2 py-0.5 text-xs"
                      :class="getPurchaseStatusClass(po.status)"
                    >
                      {{ getPurchaseStatusLabel(po.status) }}
                    </span>
                  </button>
                </div>
                <div
                  v-else-if="purchaseDropdownVisible && purchaseSearchQuery && purchaseSearchResults.length === 0 && !purchaseSearchLoading"
                  class="absolute z-10 mt-1 w-full rounded-lg border border-gray-200 bg-white py-3 text-center text-sm text-gray-400 dark:border-dk-muted dark:bg-dk-card"
                >
                  未找到匹配的订单
                </div>
              </div>
            </div>
          </div>

          <!-- 第二行：备注 -->
          <div class="mb-3">
            <label class="mb-1 block text-xs font-medium text-gray-400">{{ t('work_order.commit_comment_label') }}</label>
            <textarea
              v-model="commitComment"
              rows="2"
              :placeholder="t('work_order.commit_comment_placeholder')"
              class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            />
          </div>

          <!-- 第三行：图片上传 -->
          <div class="mb-3">
            <label class="mb-1 block text-xs font-medium text-gray-400">{{ t('work_order.commit_photos_label') }}</label>
            <useDropzone
              ref="commitDropzoneRef"
              :maxFiles="10"
              :maxSize="10 * 1024 * 1024"
              accept="image/*"
            />
          </div>

          <!-- 第四行：提交 -->
          <div class="flex justify-end">
            <button
              :disabled="submittingCommit || !canSubmit"
              class="rounded-lg bg-blue-600 px-6 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
              @click="handleCommit"
            >
              {{ submittingCommit ? t('message.submitting') : t('work_order.commit_submit') }}
            </button>
          </div>
        </div>

        <!-- 进度列表 -->
        <div class="px-6 py-4">
          <div v-if="commits.length === 0" class="py-4 text-sm text-gray-400">{{ t('work_order.no_commits') }}</div>
          <ol v-else class="relative border-l border-gray-200 dark:border-dk-muted">
            <li
              v-for="(commit, index) in commits"
              :key="commit.ID"
              class="mb-6 ml-4 rounded-lg border border-gray-100 bg-gray-50/50 px-4 py-3 dark:border-dk-muted dark:bg-dk-base/30"
            >
              <!-- 时间线圆点 -->
              <div
                class="absolute -left-1.5 mt-1.5 h-3 w-3 rounded-full border-2 border-white dark:border-dk-card"
                :class="commit.Status ? getStatusColorClass(commit.Status).split(' ')[0] : 'bg-gray-400'"
              ></div>

              <div class="flex flex-wrap items-center gap-2">
                <!-- 状态标签 -->
                <span
                  v-if="commit.Status"
                  class="inline-block rounded-full px-2.5 py-0.5 text-xs font-semibold"
                  :class="getStatusColorClass(commit.Status)"
                >
                  {{ getStatusLabel(commit.Status) }}
                </span>
                <!-- 操作人 -->
                <span class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400">
                  <img
                    :src="usersStore.getAvatarUrlFromUserID(commit.UserID)"
                    class="h-4 w-4 rounded-full object-cover"
                  />
                  {{ usersStore.getUsernameFromUserID(commit.UserID) }}
                </span>
                <!-- 时间 -->
                <time class="text-xs text-gray-400">{{ formatDate(commit.CreatedAt) }}</time>
                <!-- 删除按钮 -->
                <button
                  v-if="canDeleteCommit(commit, index)"
                  class="ml-auto rounded-lg border border-red-200 bg-red-50 px-3 py-1.5 text-xs font-medium text-red-600 transition-colors hover:bg-red-100 hover:border-red-300 dark:border-red-900/50 dark:bg-red-900/30 dark:text-red-400 dark:hover:bg-red-900/50"
                  @click="handleDeleteCommit(commit.ID)"
                >
                  <IconTrash :size="14" class="mr-1 inline align-middle" />
                  删除
                </button>
              </div>

              <!-- 备注文字 -->
              <p
                v-if="commit.Comment && commit.Comment !== '状态变更为: ' + commit.Status && commit.Action !== 'create'"
                class="mt-1.5 text-sm text-gray-700 dark:text-gray-300"
              >
                {{ commit.Comment }}
              </p>
              <p v-else-if="commit.Action === 'create'" class="mt-1.5 text-sm text-gray-500 dark:text-gray-400">
                {{ t('work_order.commit_create') }}
              </p>

              <!-- 进度图片 -->
              <div
                v-if="commit.photos && commit.photos.length > 0"
                class="mt-2 flex flex-wrap gap-2"
              >
                <a
                  v-for="file in commit.photos"
                  :key="file.ID"
                  :href="getPhotoUrl(file)"
                  target="_blank"
                  class="group relative block h-16 w-16 overflow-hidden rounded-lg border border-gray-200 bg-gray-50 transition hover:border-blue-400 dark:border-dk-muted dark:bg-dk-base"
                >
                  <img
                    :src="getPhotoUrl(file)"
                    :alt="file.Name"
                    class="h-full w-full object-cover transition group-hover:scale-105"
                  />
                </a>
              </div>

              <!-- 关联的采购订单 -->
              <div
                v-if="commit.purchaseOrders && commit.purchaseOrders.length > 0"
                class="mt-2 flex flex-wrap items-center gap-2"
              >
                <span class="text-xs text-gray-400">关联采购订单:</span>
                <RouterLink
                  v-for="po in commit.purchaseOrders"
                  :key="po.id"
                  :to="`/purchase/showorder/${po.id}`"
                  class="inline-flex items-center gap-1 rounded-full border border-blue-200 bg-blue-50 px-2.5 py-1 text-xs font-medium text-blue-700 transition-colors hover:bg-blue-100 dark:border-blue-800 dark:bg-blue-900/30 dark:text-blue-300 dark:hover:bg-blue-900/50"
                >
                  #{{ po.id }} {{ po.title || '' }}
                  <span
                    class="rounded-full px-1.5 py-0.5 text-[10px]"
                    :class="getPurchaseStatusClass(po.status)"
                  >
                    {{ getPurchaseStatusLabel(po.status) }}
                  </span>
                </RouterLink>
              </div>
            </li>
          </ol>
        </div>
      </div>

    </div>
  </div>

  <!-- 删除进度确认弹窗 -->
  <ConfirmDialog
    v-model="showDeleteCommitConfirm"
    :message="t('work_order.confirm_delete_commit')"
    danger
    @confirm="confirmDeleteCommit"
  />
</template>
