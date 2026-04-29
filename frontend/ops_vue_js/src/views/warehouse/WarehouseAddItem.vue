<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { warehouseApi } from '@/api/warehouse'
import { customerApi } from '@/api/customer'
import useDropzone from '@/components/useDropzone.vue'
import { IconUser, IconX } from '@tabler/icons-vue'

usePageTitle('warehouse.add_item')
const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const toast = useToastStore()

const containerId = ref(parseInt(route.params.id))

const form = reactive({
  name: '',
  serial_number: '',
  remark: '',
  quantity: 1,
})

const submitting = ref(false)
const loadingContainer = ref(true)
const containerName = ref('')

const dropzoneRef = ref(null)

function getPhotoHashes() {
  return dropzoneRef.value?.return_files().map((f) => f.hash) ?? []
}

// ==================== 关联客户搜索（多选） ====================
const customerSearchQuery = ref('')
const customerSearchResults = ref([])
const customerSearchLoading = ref(false)
const showCustomerDropdown = ref(false)
const selectedCustomers = ref([])

let customerSearchTimer = null

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
  if (!e.target.closest('.customer-search-wrapper')) {
    showCustomerDropdown.value = false
  }
}

onMounted(async () => {
  document.addEventListener('click', handleClickOutside)
  try {
    const { errCode, data } = await warehouseApi.getContainer(containerId.value)
    if (errCode === 0 && data?.container) {
      containerName.value = data.container.Title
    }
  } catch {
    // 找不到容器时仍允许提交
  } finally {
    loadingContainer.value = false
  }
})

async function submit() {
  if (!form.name.trim()) {
    toast.error(t('warehouse.item_name_required'))
    return
  }

  // 等待图片上传完成
  await new Promise((r) => setTimeout(r, 200))
  const hashes = getPhotoHashes()

  submitting.value = true
  try {
    const { errCode } = await warehouseApi.addItem({
      container_id: containerId.value,
      name: form.name.trim(),
      serial_number: form.serial_number.trim(),
      remark: form.remark.trim(),
      quantity: form.quantity > 0 ? form.quantity : 1,
      photos: hashes,
      customer_ids: selectedCustomers.value.map(c => c.id),
    })
    if (errCode === 0) {
      toast.success(t('message.save_success'))
      router.push(`/warehouse/container/${containerId.value}`)
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    toast.error(t('message.server_error'))
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="p-4 max-w-2xl mx-auto space-y-4">

    <!-- 面包屑 -->
    <div class="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
      <RouterLink to="/warehouse/container" class="text-blue-500 hover:underline">
        {{ t('warehouse.container_list') }}
      </RouterLink>
      <span>/</span>
      <RouterLink
        v-if="containerName"
        :to="`/warehouse/container/${containerId}`"
        class="text-blue-500 hover:underline"
      >
        {{ containerName }}
      </RouterLink>
      <svg v-if="loadingContainer" class="h-3.5 w-3.5 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
      </svg>
      <span v-else>/</span>
      <span class="text-gray-700 dark:text-gray-200">{{ t('warehouse.add_item') }}</span>
    </div>

    <!-- 表单卡片 -->
    <div class="rounded-xl border border-gray-200 bg-white px-5 py-5 shadow dark:border-dk-muted dark:bg-dk-card">
      <h2 class="mb-5 text-base font-semibold text-gray-900 dark:text-white">
        {{ t('warehouse.add_item') }}
        <span v-if="containerName" class="ml-2 text-sm font-normal text-gray-500">
          → {{ containerName }}
        </span>
      </h2>

      <div class="space-y-4">

        <!-- 物品名称 -->
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('warehouse.item_name') }} *
          </label>
          <input
            v-model="form.name"
            type="text"
            :placeholder="t('warehouse.item_name_placeholder')"
            class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
            @keyup.enter="submit"
          />
        </div>

        <!-- 序列号 -->
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('warehouse.serial_number') }}
          </label>
          <input
            v-model="form.serial_number"
            type="text"
            :placeholder="t('warehouse.serial_number_placeholder')"
            class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
          />
        </div>

        <!-- 数量 -->
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('warehouse.quantity') }}
          </label>
          <input
            v-model.number="form.quantity"
            type="number"
            min="1"
            class="w-28 rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
          />
        </div>

        <!-- 备注 -->
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('warehouse.remark') }}
          </label>
          <textarea
            v-model="form.remark"
            :placeholder="t('warehouse.remark_placeholder')"
            rows="3"
            class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
          ></textarea>
        </div>

        <!-- 关联客户搜索（多选） -->
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('warehouse.linked_customers') || '关联客户' }}
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
          <div class="customer-search-wrapper relative">
            <input
              v-model="customerSearchQuery"
              type="text"
              :placeholder="t('warehouse.linked_customer_placeholder') || '搜索客户...'"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
              @input="onCustomerSearchInput"
              @focus="customerSearchQuery || onCustomerSearchInput()"
            />
            <!-- 下拉结果 -->
            <div
              v-if="showCustomerDropdown && customerSearchResults.length > 0"
              class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-lg border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card"
            >
              <div
                v-for="customer in customerSearchResults"
                :key="customer.id"
                class="cursor-pointer px-3 py-2 text-sm hover:bg-blue-50 dark:hover:bg-blue-900/30"
                @click="selectCustomer(customer)"
              >
                <div class="font-medium text-gray-900 dark:text-white">{{ (customer.last_name || '') + (customer.first_name ? ' ' + customer.first_name : '') }}</div>
                <div v-if="customer.primary_phone || customer.primary_company" class="text-xs text-gray-500 dark:text-gray-400">
                  {{ customer.primary_phone }}{{ customer.primary_phone && customer.primary_company ? ' · ' : '' }}{{ customer.primary_company }}
                </div>
              </div>
            </div>
            <!-- 加载中 -->
            <div
              v-if="showCustomerDropdown && customerSearchLoading"
              class="absolute z-10 mt-1 w-full rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-500 shadow-lg dark:border-dk-muted dark:bg-dk-card"
            >
              {{ t('message.loading') }}
            </div>
            <!-- 无结果 -->
            <div
              v-if="showCustomerDropdown && !customerSearchLoading && customerSearchResults.length === 0 && customerSearchQuery.trim().length > 0"
              class="absolute z-10 mt-1 w-full rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm text-gray-500 shadow-lg dark:border-dk-muted dark:bg-dk-card"
            >
              {{ t('warehouse.linked_customer_not_found') || '未找到客户' }}
            </div>
          </div>
        </div>

        <!-- 图片上传 -->
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('purchase_addorder.upload_photos') }}
          </label>
          <useDropzone
            ref="dropzoneRef"
            uploadURL="/api/files/upload/image"
            :max-files="9"
          />
        </div>

      </div>

      <!-- 操作按钮 -->
      <div class="mt-6 flex justify-end gap-2">
        <button
          class="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white hover:bg-gray-50 dark:hover:bg-dk-muted"
          @click="router.push(`/warehouse/container/${containerId}`)"
        >
          {{ t('message.cancel') }}
        </button>
        <button
          class="inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
          :disabled="submitting"
          @click="submit"
        >
          <svg v-if="submitting" class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
          </svg>
          {{ t('message.save') }}
        </button>
      </div>
    </div>
  </div>
</template>
