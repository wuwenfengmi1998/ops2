<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { useUsersStore } from '@/stores/users'
import { warehouseApi } from '@/api/warehouse'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import {
  IconPackage,
  IconEdit,
  IconTrash,
  IconArrowRight,
  IconArrowLeft,
  IconSearch,
  IconPlus,
} from '@tabler/icons-vue'

usePageTitle('warehouse.item_detail')
const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const toast = useToastStore()
const usersStore = useUsersStore()

const itemId = computed(() => parseInt(route.params.id))

// ── 物品详情 ──
const item = ref(null)
const photos = ref([])
const commits = ref([])
const workOrders = ref([])
const canModifyItem = ref(false)
const loadingDetail = ref(true)
const notFound = ref(false)

// ── 容器名缓存 ──
const containerNames = reactive({})

// ── Tab ──
const activeTab = ref('work_orders')

// ── 移动弹窗 ──
const showMove = ref(false)
const moveTarget = ref(null)
const moveRemark = ref('')
const submittingMove = ref(false)

// 移动目标下拉
const targetContainers = ref([])
const targetSearch = ref('')
const targetLoading = ref(false)
const showTargetDropdown = ref(false)
let targetDropdownTimer = null

// ── 删除确认 ──
const showDeleteConfirm = ref(false)

// ── 工单状态颜色 ──
const statusColorMap = {
  pending: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/40 dark:text-yellow-400',
  checked: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-400',
  parts_ordered: 'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-400',
  repaired: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-400',
  returned: 'bg-gray-200 text-gray-600 dark:bg-gray-700 dark:text-gray-300',
  unrepairable: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400',
}

function getStatusClass(status) {
  return statusColorMap[status] || 'bg-gray-100 text-gray-600'
}

function getStatusLabel(status) {
  return t(`work_order.status_${status}`) || status
}

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



// ── 拉取物品详情 ──
async function fetchItem() {
  loadingDetail.value = true
  notFound.value = false
  try {
    const { errCode, data } = await warehouseApi.getItem(itemId.value)
    if (errCode === 0 && data) {
      item.value = data.item
      photos.value = data.photos ?? []
      commits.value = data.commits ?? []
      workOrders.value = data.work_orders ?? []
      canModifyItem.value = data.canModifyItem === true
      loadContainerNames()
    } else {
      notFound.value = true
    }
  } catch {
    notFound.value = true
  } finally {
    loadingDetail.value = false
  }
}

// ── 加载容器名 ──
async function loadContainerNames() {
  const ids = new Set()
  for (const c of commits.value) {
    if (c.OldContainer) ids.add(c.OldContainer)
    if (c.NewContainer) ids.add(c.NewContainer)
  }
  if (item.value?.ContainerID) ids.add(item.value.ContainerID)
  if (ids.size === 0) return

  try {
    const allContainers = []
    let page = 1
    let hasMore = true
    while (hasMore) {
      const { errCode, data } = await warehouseApi.getContainers({ entries: 300, page })
      if (errCode === 0) {
        allContainers.push(...(data.containers ?? []))
        hasMore = (data.containers ?? []).length === 300
        page++
      } else {
        break
      }
    }
    for (const c of allContainers) {
      if (ids.has(c.ID)) {
        containerNames[c.ID] = c.Title
      }
    }
  } catch {
    // ignore
  }
}

// ── 获取容器显示名 ──
function getContainerName(id) {
  if (!id) return t('warehouse.unstored')
  return containerNames[id] || `#${id}`
}

// ── 关联工单 ──
function openLinkWorkOrder() {
  if (!item.value) return
  // 存储预填数据到 localStorage
  const prefillData = {
    itemId: item.value.ID,
    title: item.value.SerialNumber
      ? `${item.value.Name}-${item.value.SerialNumber}`
      : item.value.Name,
    description: item.value.Remark || '',
  }
  localStorage.setItem('prefill_work_order', JSON.stringify(prefillData))
  router.push('/work_order/add')
}

// ── 移动 ──
async function openMove() {
  moveTarget.value = item.value?.ContainerID ?? null
  moveRemark.value = ''
  targetSearch.value = ''
  targetContainers.value = []
  showTargetDropdown.value = false
  showMove.value = true
}

function onTargetFocus() {
  if (targetDropdownTimer) clearTimeout(targetDropdownTimer)
  showTargetDropdown.value = true
  loadTargetContainers(targetSearch.value)
}

function onTargetInput() {
  if (targetDropdownTimer) clearTimeout(targetDropdownTimer)
  showTargetDropdown.value = true
  loadTargetContainers(targetSearch.value)
}

function closeTargetDropdown() {
  targetDropdownTimer = setTimeout(() => {
    showTargetDropdown.value = false
  }, 150)
}

async function loadTargetContainers(search = '') {
  targetLoading.value = true
  try {
    const isSearch = search.trim().length > 0
    const { errCode, data } = await warehouseApi.getContainers({
      search,
      all_levels: true,
      entries: isSearch ? 50 : 10,
      page: 1,
    })
    if (errCode === 0) {
      const filtered = (data.containers ?? []).filter(
        (c) => c.ID !== item.value?.ContainerID
      )
      targetContainers.value = isSearch ? filtered : filtered.slice(0, 5)
    }
  } catch {
    targetContainers.value = []
  } finally {
    targetLoading.value = false
  }
}

function selectTarget(id, title) {
  moveTarget.value = id
  targetSearch.value = title
  showTargetDropdown.value = false
}

async function submitMove() {
  if (moveTarget.value === item.value?.ContainerID) {
    showMove.value = false
    return
  }
  submittingMove.value = true
  try {
    const { errCode } = await warehouseApi.moveItem({
      item_id: itemId.value,
      new_container: moveTarget.value,
      remark: moveRemark.value.trim(),
    })
    if (errCode === 0) {
      toast.success(t('message.save_success'))
      showMove.value = false
      fetchItem()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    toast.error(t('message.server_error'))
  } finally {
    submittingMove.value = false
  }
}

// ── 返回 ──
function goBack() {
  const cid = item.value?.ContainerID
  if (cid) {
    router.push(`/warehouse/container/${cid}`)
  } else {
    router.push('/warehouse/container')
  }
}

// ── 删除 ──
async function doDelete() {
  try {
    const { errCode } = await warehouseApi.deleteItem(itemId.value)
    if (errCode === 0) {
      toast.success(t('message.delete_ok'))
      showDeleteConfirm.value = false
      const cid = item.value?.ContainerID
      if (cid) {
        router.push(`/warehouse/container/${cid}`)
      } else {
        router.push('/warehouse/item')
      }
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    toast.error(t('message.server_error'))
  }
}

// ── 初始化 ──
onMounted(() => {
  fetchItem()
})
</script>

<template>
  <div class="p-4 max-w-5xl mx-auto space-y-4">

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
      <button class="mt-2 text-sm text-blue-500 hover:underline" @click="router.push('/warehouse/item')">
        {{ t('warehouse.back_to_list') }}
      </button>
    </div>

    <template v-else-if="item">

      <!-- 操作栏 -->
      <div class="flex items-center justify-between gap-2">
        <button
          class="inline-flex items-center gap-1.5 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:hover:bg-dk-muted"
          @click="goBack"
        >
          <IconArrowLeft :size="14" />
          返回
        </button>
        <div class="flex gap-2">
          <button
            class="inline-flex items-center gap-1.5 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:hover:bg-dk-muted"
            @click="openLinkWorkOrder"
          >
            <IconPlus :size="14" />
            关联工单
          </button>
          <button
            class="inline-flex items-center gap-1.5 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:hover:bg-dk-muted"
            @click="openMove"
          >
            <IconArrowRight :size="14" />
            {{ t('warehouse.move_item') }}
          </button>
          <button
            v-if="canModifyItem"
            class="inline-flex items-center gap-1.5 rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-base dark:text-white dark:hover:bg-dk-muted"
            @click="router.push(`/warehouse/item/edit/${itemId}`)"
          >
            <IconEdit :size="14" />
            {{ t('warehouse.edit') }}
          </button>
          <button
            v-if="canModifyItem"
            class="inline-flex items-center gap-1.5 rounded-lg border border-red-300 bg-white px-3 py-1.5 text-sm font-medium text-red-600 transition-colors hover:bg-red-50 dark:border-red-900 dark:bg-dk-base dark:text-red-400 dark:hover:bg-red-900/20"
            @click="showDeleteConfirm = true"
          >
            <IconTrash :size="14" />
            {{ t('warehouse.delete') }}
          </button>
        </div>
      </div>

      <!-- 物品信息卡 -->
      <div class="rounded-xl border border-gray-200 bg-white px-5 py-4 shadow dark:border-dk-muted dark:bg-dk-card">
        <div class="flex items-start gap-4 mb-3">
          <div class="flex-shrink-0 mt-1">
            <IconPackage :size="32" class="text-blue-500" />
          </div>
          <div class="flex-1 min-w-0">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ item.Name }}</h2>
            <div class="flex flex-wrap gap-x-6 gap-y-1 mt-1 text-sm text-gray-500">
              <span v-if="item.SerialNumber">{{ t('warehouse.serial_number') }}: {{ item.SerialNumber }}</span>
              <span>{{ t('warehouse.quantity') }}: {{ item.Quantity }}</span>
              <span>{{ t('warehouse.location') }}:
                <RouterLink
                  v-if="item.ContainerID"
                  :to="`/warehouse/container/${item.ContainerID}`"
                  class="text-blue-500 hover:underline ml-1"
                >
                  {{ containerNames[item.ContainerID] || `#${item.ContainerID}` }}
                </RouterLink>
                <span v-else class="text-orange-500 ml-1">{{ t('warehouse.unstored') }}</span>
              </span>
            </div>
          </div>
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
          v-if="item.Remark"
          class="text-sm text-gray-500 dark:text-gray-400 whitespace-pre-wrap mb-3"
        >
          {{ item.Remark }}
        </p>

        <!-- 元信息 -->
        <div class="flex flex-wrap gap-x-6 gap-y-1 text-xs text-gray-400 dark:text-gray-500">
          <span class="flex items-center gap-1">
            <span>{{ t('warehouse.created_by') }}:</span>
            <img
              :src="usersStore.getAvatarUrlFromUserID(item.CreatorID)"
              class="w-4 h-4 rounded-full object-cover"
            />
            {{ usersStore.getUsernameFromUserID(item.CreatorID) }}
          </span>
          <span>{{ t('warehouse.created_at') }}: {{ fmtTs(item.CreatedAt) }}</span>
          <span v-if="item.UpdatedAt">{{ t('warehouse.updated_at') }}: {{ fmtTs(item.UpdatedAt) }}</span>
        </div>
      </div>

      <!-- Tab 切换 -->
      <div class="flex gap-1 rounded-lg border border-gray-200 bg-gray-50 p-1 dark:border-dk-muted dark:bg-dk-base w-fit">
        <button
          class="px-4 py-1.5 text-sm rounded-md font-medium transition-colors"
          :class="activeTab === 'work_orders'
            ? 'bg-white text-gray-900 shadow-sm dark:bg-dk-card dark:text-white'
            : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
          @click="activeTab = 'work_orders'"
        >
          {{ t('warehouse.work_orders') }} ({{ workOrders.length }})
        </button>
        <button
          class="px-4 py-1.5 text-sm rounded-md font-medium transition-colors"
          :class="activeTab === 'history'
            ? 'bg-white text-gray-900 shadow-sm dark:bg-dk-card dark:text-white'
            : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
          @click="activeTab = 'history'"
        >
          {{ t('warehouse.move_history') }} ({{ commits.length }})
        </button>
      </div>

      <!-- 关联工单 -->
      <div v-if="activeTab === 'work_orders'">
        <div v-if="workOrders.length === 0" class="rounded-xl border border-gray-200 bg-white px-5 py-8 text-center text-sm text-gray-400 dark:border-dk-muted dark:bg-dk-card">
          {{ t('warehouse.no_work_orders') }}
        </div>
        <div v-else class="space-y-2">
          <RouterLink
            v-for="wo in workOrders"
            :key="wo.id"
            :to="`/work_order/show/${wo.id}`"
            class="rounded-xl border border-gray-200 bg-white px-4 py-3 flex items-center justify-between gap-3 hover:shadow transition-shadow dark:border-dk-muted dark:bg-dk-card dark:hover:shadow-none"
          >
            <div class="flex items-center gap-3 min-w-0">
              <IconPackage :size="16" class="text-blue-500 flex-shrink-0" />
              <span class="font-medium text-sm text-gray-900 truncate dark:text-white">{{ wo.title }}</span>
            </div>
            <span
              class="flex-shrink-0 rounded-full px-2.5 py-0.5 text-xs font-medium"
              :class="getStatusClass(wo.status)"
            >
              {{ getStatusLabel(wo.status) }}
            </span>
          </RouterLink>
        </div>
      </div>

      <!-- 移动历史 -->
      <div v-if="activeTab === 'history'">
        <div v-if="commits.length === 0" class="rounded-xl border border-gray-200 bg-white px-5 py-8 text-center text-sm text-gray-400 dark:border-dk-muted dark:bg-dk-card">
          {{ t('warehouse.no_move_history') }}
        </div>
        <div v-else class="space-y-2">
          <div
            v-for="commit in commits"
            :key="commit.ID"
            class="rounded-xl border border-gray-200 bg-white px-4 py-3 flex items-center gap-3 dark:border-dk-muted dark:bg-dk-card"
          >
            <!-- 操作人头像 -->
            <div class="flex-shrink-0">
              <img
                :src="usersStore.getAvatarUrlFromUserID(commit.UserID)"
                class="w-8 h-8 rounded-full object-cover"
              />
            </div>
            <!-- 路径 -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap text-xs text-gray-400">
                <span>{{ usersStore.getUsernameFromUserID(commit.UserID) || `User#${commit.UserID}` }}</span>
                <span>{{ fmtTs(commit.CreatedAt) }}</span>
              </div>
              <div class="flex items-center gap-1.5 mt-0.5 flex-wrap text-sm font-medium text-gray-700 dark:text-gray-200">
                <span>{{ getContainerName(commit.OldContainer) }}</span>
                <IconArrowRight :size="13" class="text-blue-500 flex-shrink-0" />
                <span>{{ getContainerName(commit.NewContainer) }}</span>
              </div>
              <p v-if="commit.Remark" class="text-xs text-gray-400 mt-0.5">{{ commit.Remark }}</p>
            </div>
          </div>
        </div>
      </div>

    </template>
  </div>

  <!-- 移动弹窗 -->
  <Transition name="fade">
    <div
      v-if="showMove"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/40"
      @click.self="showMove = false"
    >
      <div class="w-full max-w-md rounded-xl border border-gray-200 bg-white p-5 shadow-xl dark:border-dk-muted dark:bg-dk-card">
        <h3 class="mb-4 text-base font-semibold text-gray-900 dark:text-white">{{ t('warehouse.move_item') }}</h3>
        <div class="space-y-3">
          <div class="text-sm text-gray-500">
            {{ t('warehouse.current_location') }}:
            <span class="font-medium text-gray-700 dark:text-gray-300">
              <RouterLink
                v-if="item?.ContainerID"
                :to="`/warehouse/container/${item.ContainerID}`"
                class="text-blue-500 hover:underline"
              >
                {{ containerNames[item.ContainerID] || `#${item.ContainerID}` }}
              </RouterLink>
              <span v-else class="text-orange-500">{{ t('warehouse.unstored') }}</span>
            </span>
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('warehouse.target_container') }}</label>
            <div class="relative" @click.stop>
              <IconSearch class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" :size="15" />
              <input
                v-model="targetSearch"
                type="text"
                :placeholder="t('warehouse.search_container')"
                class="w-full rounded-lg border border-gray-300 bg-white py-2 pl-9 pr-3 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
                @focus="onTargetFocus"
                @input="onTargetInput"
                @blur="closeTargetDropdown"
              />
              <div
                v-if="showTargetDropdown"
                class="absolute z-10 mt-1 w-full rounded-lg border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card max-h-60 overflow-y-auto"
                @mousedown.prevent
              >
                <div v-if="targetLoading" class="px-3 py-2 text-xs text-gray-400">
                  <svg class="inline h-3.5 w-3.5 animate-spin mr-1" viewBox="0 0 24 24" fill="none">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                  </svg>
                  {{ t('message.loading') }}
                </div>
                <button
                  class="flex w-full items-center gap-2 px-3 py-2 text-sm text-left text-orange-500 hover:bg-gray-50 dark:text-orange-400 dark:hover:bg-dk-muted"
                  @click="selectTarget(null, t('warehouse.unstored')); showTargetDropdown = false"
                >
                  {{ t('warehouse.unstored') }}
                </button>
                <button
                  v-for="c in targetContainers"
                  :key="c.ID"
                  class="flex w-full items-center gap-2 px-3 py-2 text-sm text-left hover:bg-gray-50 dark:text-gray-200 dark:hover:bg-dk-muted"
                  @click="selectTarget(c.ID, c.Title); showTargetDropdown = false"
                >
                  <IconPackage :size="13" class="text-blue-500 flex-shrink-0" />
                  {{ c.Title }}
                </button>
              </div>
            </div>
            <div class="mt-1 text-xs">
              <span v-if="moveTarget === null" class="text-orange-500">→ {{ t('warehouse.unstored') }}</span>
              <span v-else-if="moveTarget" class="text-blue-500">→ {{ containerNames[moveTarget] || `#${moveTarget}` }}</span>
            </div>
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('warehouse.remark') }}</label>
            <textarea
              v-model="moveRemark"
              :placeholder="t('warehouse.move_remark_placeholder')"
              rows="2"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
            ></textarea>
          </div>
        </div>
        <div class="mt-4 flex justify-end gap-2">
          <button
            class="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white hover:bg-gray-50 dark:hover:bg-dk-muted"
            @click="showMove = false"
          >
            {{ t('message.cancel') }}
          </button>
          <button
            class="inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            :disabled="submittingMove"
            @click="submitMove"
          >
            <svg v-if="submittingMove" class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
            </svg>
            {{ t('warehouse.confirm_move') }}
          </button>
        </div>
      </div>
    </div>
  </Transition>

  <!-- 删除确认 -->
  <ConfirmDialog
    v-model="showDeleteConfirm"
    :title="t('warehouse.delete_item_title')"
    :message="t('warehouse.delete_item_msg', { name: item?.Name })"
    danger
    @confirm="doDelete"
  />

</template>
