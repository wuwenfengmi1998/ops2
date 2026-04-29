<script setup>
/**
 * 工单新增/编辑页面
 * - 路由有 :id 参数时为编辑模式，否则为新增模式
 * - 支持图片上传（复用 useDropzone 组件）
 * - 支持关联物品搜索
 */
import { reactive, ref, computed, onMounted, nextTick, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { useValidation } from '@/composables'
import { workOrderApi } from '@/api/work_order'
import { warehouseApi } from '@/api/warehouse'
import { customerApi } from '@/api/customer'
import useDropzone from '@/components/useDropzone.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { IconPackage, IconUser, IconX } from '@tabler/icons-vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToastStore()
const { validate, errors, clearErrors } = useValidation()

const isEdit = computed(() => !!route.params.id)
const orderId = computed(() => (isEdit.value ? Number(route.params.id) : null))

usePageTitle(isEdit.value ? 'work_order.edit_title' : 'work_order.add_title')

// ==================== 状态 ====================
const loading = ref(false)
const pageLoading = ref(false)
const pageError = ref('')
const showDeleteConfirm = ref(false)

// ==================== 关联物品搜索（多选） ====================
const itemSearchQuery = ref('')
const itemSearchResults = ref([])
const itemSearchLoading = ref(false)
const showItemDropdown = ref(false)
const selectedItems = ref([])

// ==================== 关联客户搜索（多选） ====================
const customerSearchQuery = ref('')
const customerSearchResults = ref([])
const customerSearchLoading = ref(false)
const showCustomerDropdown = ref(false)
const selectedCustomers = ref([])

let searchTimer = null
let customerSearchTimer = null

function onItemSearchInput() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(async () => {
    itemSearchLoading.value = true
    showItemDropdown.value = true
    try {
      let res
      if (itemSearchQuery.value.trim().length > 0) {
        // 有搜索词：搜索匹配物品
        res = await warehouseApi.getItems({ search: itemSearchQuery.value.trim() })
        if (res.errCode === 0 && res.data) {
          itemSearchResults.value = (res.data.items || []).slice(0, 10)
        } else {
          itemSearchResults.value = []
        }
      } else {
        // 无搜索词：显示最新5个物品
        res = await warehouseApi.getItems({ page: 1, page_size: 5 })
        if (res.errCode === 0 && res.data) {
          // 按创建时间倒序（最新在前）
          itemSearchResults.value = (res.data.items || []).sort((a, b) => {
            const tsA = parseInt(a.CreatedAt || '0', 10)
            const tsB = parseInt(b.CreatedAt || '0', 10)
            return tsB - tsA
          })
        } else {
          itemSearchResults.value = []
        }
      }
    } catch {
      itemSearchResults.value = []
    } finally {
      itemSearchLoading.value = false
    }
  }, 300)
}

function selectItem(item) {
  // 检查是否已选中
  if (!selectedItems.value.find(i => i.ID === item.ID)) {
    selectedItems.value.push(item)
  }
  itemSearchQuery.value = ''
  itemSearchResults.value = []
  showItemDropdown.value = false
}

function removeSelectedItem(itemId) {
  selectedItems.value = selectedItems.value.filter(i => i.ID !== itemId)
}

// ==================== 关联客户搜索 ====================
function onCustomerSearchInput() {
  clearTimeout(customerSearchTimer)
  customerSearchTimer = setTimeout(async () => {
    customerSearchLoading.value = true
    showCustomerDropdown.value = true
    try {
      let res
      if (customerSearchQuery.value.trim().length > 0) {
        res = await customerApi.list({ search: customerSearchQuery.value.trim(), page: 1, page_size: 10 })
        if (res.errCode === 0 && res.data) {
          customerSearchResults.value = (res.data.customers || []).slice(0, 10)
        } else {
          customerSearchResults.value = []
        }
      } else {
        res = await customerApi.list({ page: 1, page_size: 5 })
        if (res.errCode === 0 && res.data) {
          customerSearchResults.value = (res.data.customers || []).sort((a, b) => b.ID - a.ID)
        } else {
          customerSearchResults.value = []
        }
      }
    } catch {
      customerSearchResults.value = []
    } finally {
      customerSearchLoading.value = false
    }
  }, 300)
}

function selectCustomer(customer) {
  // 检查是否已选中
  if (!selectedCustomers.value.find(c => c.id === customer.id)) {
    selectedCustomers.value.push(customer)
  }
  customerSearchQuery.value = ''
  customerSearchResults.value = []
  showCustomerDropdown.value = false
}

function removeSelectedCustomer(customerId) {
  selectedCustomers.value = selectedCustomers.value.filter(c => c.id !== customerId)
}

function handleClickOutside(e) {
  if (!e.target.closest('.item-search-wrapper') && !e.target.closest('.customer-search-wrapper')) {
    showItemDropdown.value = false
    showCustomerDropdown.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

// ==================== 表单数据 ====================
const form = reactive({
  title: '',
  description: '',
  photos: [],
})

// ==================== 图片上传 ====================
const dropzoneRef = ref(null)

function getPhotoHashes() {
  return dropzoneRef.value?.return_files().map((f) => f.hash) ?? []
}

// ==================== 加载编辑数据 ====================
const linkedItemIds = ref([])
const linkedCustomerIds = ref([])

onMounted(async () => {
  // 新增模式：检查预填数据
  if (!isEdit.value) {
    const prefillStr = localStorage.getItem('prefill_work_order')
    if (prefillStr) {
      try {
        const prefill = JSON.parse(prefillStr)
        form.title = prefill.title || ''
        form.description = prefill.description || ''

        // 如果有物品ID，获取物品详情并自动选中
        if (prefill.itemId) {
          try {
            const itemRes = await warehouseApi.getItem(prefill.itemId)
            if (itemRes.errCode === 0 && itemRes.data) {
              const item = itemRes.data.item
              if (!selectedItems.value.find(i => i.ID === item.ID)) {
                selectedItems.value.push(item)
              }
              if (!linkedItemIds.value.includes(prefill.itemId)) {
                linkedItemIds.value.push(prefill.itemId)
              }
            }
          } catch {
            // 获取物品详情失败，忽略
          }
        }

        // 如果有客户信息，自动选中
        if (prefill.customers && prefill.customers.length > 0) {
          selectedCustomers.value = prefill.customers.map(c => ({
            id: c.id,
            first_name: c.first_name,
            last_name: c.last_name,
            primary_phone: c.primary_phone,
          }))
          linkedCustomerIds.value = prefill.customers.map(c => c.id)
        }

        localStorage.removeItem('prefill_work_order')
      } catch {
        // ignore
      }
    }
    return
  }

  // 编辑模式
  pageLoading.value = true
  try {
    const res = await workOrderApi.get(orderId.value)
    if (res.errCode !== 0 || !res.data) {
      pageError.value = t('work_order.not_found')
      return
    }

    const { order, photos, linkedItems, linkedCustomers } = res.data
    form.title = order.Title ?? ''
    form.description = order.Description ?? ''

    // 回填图片
    await nextTick()
    if (photos && photos.length > 0) {
      form.photos = photos
    }

    // 回填关联物品
    if (linkedItems && linkedItems.length > 0) {
      selectedItems.value = linkedItems.map(item => ({
        ID: item.ID,
        Name: item.Name,
        SerialNumber: item.SerialNumber,
        ContainerID: item.ContainerID,
      }))
      linkedItemIds.value = linkedItems.map(item => item.ID)
    }

    // 回填关联客户
    if (linkedCustomers && linkedCustomers.length > 0) {
      selectedCustomers.value = linkedCustomers.map(c => ({
        id: c.id,
        first_name: c.first_name,
        last_name: c.last_name,
        primary_phone: c.primary_phone,
      }))
      linkedCustomerIds.value = linkedCustomers.map(c => c.id)
    }
  } catch {
    pageError.value = t('work_order.not_found')
  } finally {
    pageLoading.value = false
  }
})

// ==================== 删除 ====================
function handleDelete() {
  showDeleteConfirm.value = true
}

async function doDelete() {
  loading.value = true
  try {
    const res = await workOrderApi.delete(orderId.value)
    if (res.errCode === 0) {
      toast.success(t('message.delete_ok'))
      router.replace('/work_order')
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    toast.error(t('message.server_error'))
  } finally {
    loading.value = false
  }
}

// ==================== 提交 ====================
async function handleSubmit() {
  clearErrors()
  const ok = validate('title', form.title, t('work_order.title'))
  if (!ok) return

  const photos = getPhotoHashes()
  loading.value = true
  try {
    let res
    if (isEdit.value) {
      res = await workOrderApi.update({
        id: orderId.value,
        title: form.title,
        description: form.description,
        photos,
        item_ids: selectedItems.value.map(i => i.ID),
        customer_ids: selectedCustomers.value.map(c => c.id),
      })
    } else {
      res = await workOrderApi.add({
        title: form.title,
        description: form.description,
        photos,
        item_ids: selectedItems.value.map(i => i.ID),
        customer_ids: selectedCustomers.value.map(c => c.id),
      })
    }

    if (res.errCode === 0) {
      toast.success(t('message.save_ok'))
      const newId = isEdit.value ? orderId.value : res.data?.id
      setTimeout(() => {
        router.replace(newId ? `/work_order/show/${newId}` : '/work_order')
      }, 800)
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    toast.error(t('message.server_error'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-3xl px-6 py-6">
    <!-- 加载中 -->
    <div v-if="pageLoading" class="flex items-center justify-center py-20 text-gray-400">
      <svg class="mr-2 h-5 w-5 animate-spin" viewBox="0 0 24 24" fill="none">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
      </svg>
      {{ t('message.loading') }}
    </div>

    <!-- 工单不存在 -->
    <div
      v-else-if="pageError"
      class="rounded-xl border border-red-200 bg-red-50 px-6 py-10 text-center text-red-500 dark:border-red-800 dark:bg-red-900/20"
    >
      {{ pageError }}
    </div>

    <!-- 主卡片 -->
    <div v-else class="flex flex-col gap-0 rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <!-- 顶部标题栏 -->
      <div class="flex items-center justify-between border-b border-gray-200 px-6 py-4 dark:border-dk-muted">
        <h4 class="text-sm font-semibold text-gray-900 dark:text-white">
          {{ isEdit ? t('work_order.edit_title') : t('work_order.add_title') }}
        </h4>
        <div class="flex items-center gap-2">
          <!-- 删除按钮（编辑模式才显示） -->
          <button
            v-if="isEdit"
            class="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20"
            :disabled="loading"
            @click="handleDelete"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            {{ t('work_order.delete') }}
          </button>
          <!-- 返回按钮 -->
          <button
            class="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dk-base"
            @click="router.back()"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
            </svg>
            {{ t('purchase.back') }}
          </button>
        </div>
      </div>

      <!-- 字段验证错误提示 -->
      <div
        v-if="errors.title"
        class="mx-6 mt-4 rounded-lg bg-red-50 px-4 py-2 text-sm text-red-600 dark:bg-red-900/20 dark:text-red-400"
      >
        {{ errors.title }}
      </div>

      <!-- 表单区 -->
      <div class="space-y-5 px-6 py-5">
        <!-- 工单标题 -->
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('work_order.title') }}<span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.title"
            type="text"
            maxlength="200"
            :placeholder="t('work_order.title_placeholder')"
            class="w-full rounded-lg border px-3.5 py-2 text-sm outline-none transition-colors focus:ring-2 focus:ring-blue-500/20 dark:bg-dk-base dark:text-white"
            :class="errors.title ? 'border-red-500 focus:border-red-500' : 'border-gray-300 focus:border-blue-500 dark:border-dk-muted'"
          />
          <span v-if="errors.title" class="mt-1 block text-xs text-red-500">{{ errors.title }}</span>
        </div>

        <!-- 问题描述 -->
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('work_order.description') }}
          </label>
          <textarea
            v-model="form.description"
            rows="5"
            :placeholder="t('work_order.description_placeholder')"
            class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
          />
        </div>

        <!-- 关联物品搜索（多选） -->
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('work_order.linked_items') }}
          </label>

          <!-- 已选择物品列表 -->
          <div v-if="selectedItems.length > 0" class="mb-2 flex flex-wrap gap-2">
            <div
              v-for="item in selectedItems"
              :key="item.ID"
              class="inline-flex items-center gap-1 rounded-full border border-blue-200 bg-blue-50 px-2.5 py-1 text-xs font-medium text-blue-700 dark:border-blue-800 dark:bg-blue-900/30 dark:text-blue-300"
            >
              <IconPackage :size="12" />
              {{ item.Name }}{{ item.SerialNumber ? ' - ' + item.SerialNumber : '' }}
              <button
                type="button"
                class="ml-0.5 rounded-full p-0.5 transition-colors hover:bg-blue-200 dark:hover:bg-blue-800"
                @click="removeSelectedItem(item.ID)"
              >
                <IconX :size="12" />
              </button>
            </div>
          </div>

          <!-- 搜索框 -->
          <div class="item-search-wrapper relative" style="z-index: 9999;">
            <input
              v-model="itemSearchQuery"
              type="text"
              :placeholder="t('work_order.linked_item_placeholder')"
              class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
              @input="onItemSearchInput"
              @focus="itemSearchQuery || onItemSearchInput()"
            />
            <!-- 下拉结果 -->
            <div
              v-if="showItemDropdown && itemSearchResults.length > 0"
              class="absolute z-[9999] mt-1 max-h-60 w-full overflow-auto rounded-lg border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card"
            >
              <div
                v-for="item in itemSearchResults"
                :key="item.ID"
                class="cursor-pointer px-3 py-2 text-sm hover:bg-blue-50 dark:hover:bg-blue-900/30"
                @click="selectItem(item)"
              >
                <div class="font-medium text-gray-900 dark:text-white">{{ item.Name }}</div>
                <div v-if="item.SerialNumber" class="text-xs text-gray-500 dark:text-gray-400">{{ item.SerialNumber }}</div>
              </div>
            </div>
            <!-- 加载中 -->
            <div
              v-if="showItemDropdown && itemSearchLoading"
              class="absolute z-[9999] mt-1 w-full rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-500 shadow-lg dark:border-dk-muted dark:bg-dk-card"
            >
              {{ t('message.loading') }}
            </div>
            <!-- 无结果 -->
            <div
              v-if="showItemDropdown && !itemSearchLoading && itemSearchResults.length === 0 && itemSearchQuery.trim().length > 0"
              class="absolute z-[9999] mt-1 w-full rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-500 shadow-lg dark:border-dk-muted dark:bg-dk-card"
            >
              {{ t('work_order.linked_item_not_found') }}
            </div>
          </div>
        </div>

        <!-- 关联客户搜索（多选） -->
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('work_order.linked_customers') }}
          </label>

          <!-- 已选择客户列表 -->
          <div v-if="selectedCustomers.length > 0" class="mb-2 flex flex-wrap gap-2">
            <div
              v-for="customer in selectedCustomers"
              :key="customer.id"
              class="inline-flex items-center gap-1 rounded-full border border-blue-200 bg-blue-50 px-2.5 py-1 text-xs font-medium text-blue-700 dark:border-blue-800 dark:bg-blue-900/30 dark:text-blue-300"
            >
              <IconUser :size="12" />
              {{ (customer.last_name || '') + (customer.first_name ? ' ' + customer.first_name : '') }}
              <button
                type="button"
                class="ml-0.5 rounded-full p-0.5 transition-colors hover:bg-blue-200 dark:hover:bg-blue-800"
                @click="removeSelectedCustomer(customer.id)"
              >
                <IconX :size="12" />
              </button>
            </div>
          </div>

          <!-- 搜索框 -->
          <div class="customer-search-wrapper relative" style="z-index: 9999;">
            <input
              v-model="customerSearchQuery"
              type="text"
              :placeholder="t('work_order.linked_customer_placeholder')"
              class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
              @input="onCustomerSearchInput"
              @focus="customerSearchQuery || onCustomerSearchInput()"
            />
            <!-- 下拉结果 -->
            <div
              v-if="showCustomerDropdown && customerSearchResults.length > 0"
              class="absolute z-[9999] mt-1 max-h-60 w-full overflow-auto rounded-lg border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card"
            >
              <div
                v-for="customer in customerSearchResults"
                :key="customer.ID"
                class="cursor-pointer px-3 py-2 text-sm hover:bg-blue-50 dark:hover:bg-blue-900/30"
                @click="selectCustomer(customer)"
              >
                <div class="font-medium text-gray-900 dark:text-white">
                  {{ (customer.last_name || '') + (customer.first_name ? ' ' + customer.first_name : '') }}
                </div>
                <div v-if="customer.primary_phone || customer.primary_company" class="text-xs text-gray-500 dark:text-gray-400">
                  {{ customer.primary_phone }}{{ customer.primary_phone && customer.primary_company ? ' · ' : '' }}{{ customer.primary_company }}
                </div>
              </div>
            </div>
            <!-- 加载中 -->
            <div
              v-if="showCustomerDropdown && customerSearchLoading"
              class="absolute z-[9999] mt-1 w-full rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-500 shadow-lg dark:border-dk-muted dark:bg-dk-card"
            >
              {{ t('message.loading') }}
            </div>
            <!-- 无结果 -->
            <div
              v-if="showCustomerDropdown && !customerSearchLoading && customerSearchResults.length === 0 && customerSearchQuery.trim().length > 0"
              class="absolute z-[9999] mt-1 w-full rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-500 shadow-lg dark:border-dk-muted dark:bg-dk-card"
            >
              {{ t('work_order.linked_customer_not_found') }}
            </div>
          </div>
        </div>

        <!-- 图片上传 -->
        <div>
          <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('work_order.photos') }}
          </label>
          <useDropzone
            ref="dropzoneRef"
            :initial-files="form.photos"
          />
        </div>
      </div>

      <!-- 底部提交 -->
      <div class="border-t border-gray-200 px-6 py-4 dark:border-dk-muted">
        <button
          class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-5 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-60"
          :disabled="loading"
          @click="handleSubmit"
        >
          <svg v-if="loading" class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          {{ isEdit ? t('work_order.save_changes') : t('work_order.submit') }}
        </button>
      </div>
    </div>
  </div>

  <!-- 删除确认弹窗 -->
  <ConfirmDialog
    v-model="showDeleteConfirm"
    :message="t('work_order.confirm_delete')"
    danger
    @confirm="doDelete"
  />
</template>
