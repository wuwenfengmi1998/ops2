<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { customerApi } from '@/api/customer'
import AppToast from '@/components/AppToast.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import CustomerFormModal from './CustomerFormModal.vue'

const { t } = useI18n()

const customers = ref([])
const loading = ref(false)
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const showAddModal = ref(false)
const showEditModal = ref(false)
const showDeleteConfirm = ref(false)
const editingCustomer = ref(null)
const deletingCustomer = ref(null)

const toast = ref({ show: false, message: '', type: 'success' })

// 客户列表
async function fetchCustomers() {
  loading.value = true
  try {
    const res = await customerApi.list({
      page: currentPage.value,
      page_size: pageSize.value,
      search: searchQuery.value,
    })
    if (res.errCode === 0) {
      customers.value = res.data.customers || []
      total.value = res.data.total || 0
    }
  } catch {
    showToast(t('message.error'), 'error')
  } finally {
    loading.value = false
  }
}

function onSearch() {
  currentPage.value = 1
  fetchCustomers()
}

function onPageChange(page) {
  currentPage.value = page
  fetchCustomers()
}

// 新增客户
function openAddModal() {
  showAddModal.value = true
}

function onCustomerAdded() {
  showAddModal.value = false
  showToast(t('message.add_success'), 'success')
  fetchCustomers()
}

// 编辑客户
function openEditModal(customer) {
  editingCustomer.value = customer
  showEditModal.value = true
}

function onCustomerUpdated() {
  showEditModal.value = false
  editingCustomer.value = null
  showToast(t('message.update_success'), 'success')
  fetchCustomers()
}

// 删除客户
function confirmDelete(customer) {
  deletingCustomer.value = customer
  showDeleteConfirm.value = true
}

async function doDelete() {
  if (!deletingCustomer.value) return
  try {
    const res = await customerApi.delete({ id: deletingCustomer.value.id })
    if (res.errCode === 0) {
      showToast(t('message.delete_success'), 'success')
      fetchCustomers()
    } else {
      showToast(res.errMsg || t('message.error'), 'error')
    }
  } catch {
    showToast(t('message.error'), 'error')
  } finally {
    showDeleteConfirm.value = false
    deletingCustomer.value = null
  }
}

function showToast(message, type = 'success') {
  toast.value = { show: true, message, type }
  setTimeout(() => (toast.value.show = false), 3000)
}

onMounted(() => {
  fetchCustomers()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 p-6 dark:bg-dk-base">
    <div class="mx-auto max-w-6xl">
      <!-- Header -->
      <div class="mb-6 flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-dk-text">{{ t('customer.title') }}</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-dk-subtle">{{ t('customer.subtitle') }}</p>
        </div>
        <button
          @click="openAddModal"
          class="flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          {{ t('customer.add') }}
        </button>
      </div>

      <!-- Search -->
      <div class="mb-4 flex gap-3">
        <div class="relative flex-1">
          <input
            v-model="searchQuery"
            @keyup.enter="onSearch"
            type="text"
            :placeholder="t('customer.search_placeholder')"
            class="w-full rounded-lg border border-gray-300 px-4 py-2 pl-10 text-sm focus:border-blue-500 focus:outline-none dark:border-dk-muted dark:bg-dk-card dark:text-dk-text"
          />
          <svg class="absolute left-3 top-2.5 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <button
          @click="onSearch"
          class="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-card dark:text-dk-text"
        >
          {{ t('common.search') }}
        </button>
      </div>

      <!-- Customer List -->
      <div class="rounded-lg border border-gray-200 bg-white dark:border-dk-muted dark:bg-dk-card">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="border-b border-gray-200 bg-gray-50 dark:border-dk-muted dark:bg-dk-muted">
              <tr>
                <th class="px-4 py-3 font-medium text-gray-700 dark:text-dk-subtle">{{ t('customer.name') }}</th>
                <th class="px-4 py-3 font-medium text-gray-700 dark:text-dk-subtle">{{ t('customer.salutation') }}</th>
                <th class="px-4 py-3 font-medium text-gray-700 dark:text-dk-subtle">{{ t('customer.phone') }}</th>
                <th class="px-4 py-3 font-medium text-gray-700 dark:text-dk-subtle">{{ t('customer.email') }}</th>
                <th class="px-4 py-3 font-medium text-gray-700 dark:text-dk-subtle">{{ t('customer.company') }}</th>
                <th class="px-4 py-3 font-medium text-gray-700 dark:text-dk-subtle">{{ t('customer.created_at') }}</th>
                <th class="px-4 py-3 font-medium text-gray-700 dark:text-dk-subtle">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="customer in customers"
                :key="customer.id"
                class="border-b border-gray-100 hover:bg-gray-50 dark:border-dk-muted dark:hover:bg-dk-muted"
              >
                <td class="px-4 py-3">
                  <router-link
                    :to="`/customer/detail/${customer.id}`"
                    class="font-medium text-blue-600 hover:text-blue-700 hover:underline dark:text-blue-400"
                  >
                    {{ (customer.last_name || '') + (customer.first_name ? ' ' + customer.first_name : '') }}
                  </router-link>
                </td>
                <td class="px-4 py-3 text-gray-600 dark:text-dk-subtle">{{ customer.title ? t(`customer.salutation_${customer.title.toLowerCase()}`) : '-' }}</td>
                <td class="px-4 py-3 text-gray-600 dark:text-dk-subtle">{{ customer.primary_phone || '-' }}</td>
                <td class="px-4 py-3 text-gray-600 dark:text-dk-subtle">{{ customer.primary_email || '-' }}</td>
                <td class="px-4 py-3 text-gray-600 dark:text-dk-subtle">{{ customer.primary_company || '-' }}</td>
                <td class="px-4 py-3 text-gray-500 dark:text-dk-subtle">{{ new Date(customer.created_at).toLocaleString() }}</td>
                <td class="px-4 py-3">
                  <div class="flex gap-2">
                    <button
                      v-if="customer.edit"
                      @click="openEditModal(customer)"
                      class="text-blue-600 hover:text-blue-700 dark:text-blue-400"
                    >
                      {{ t('common.edit') }}
                    </button>
                    <button
                      v-if="customer.edit"
                      @click="confirmDelete(customer)"
                      class="text-red-600 hover:text-red-700 dark:text-red-400"
                    >
                      {{ t('common.delete') }}
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="customers.length === 0 && !loading">
                <td colspan="7" class="px-4 py-8 text-center text-gray-500 dark:text-dk-subtle">
                  {{ t('customer.no_data') }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div v-if="total > pageSize" class="flex items-center justify-between border-t border-gray-200 px-4 py-3 dark:border-dk-muted">
          <div class="text-sm text-gray-500 dark:text-dk-subtle">
            {{ t('common.total') }} {{ total }} {{ t('common.items') }}
          </div>
          <div class="flex gap-1">
            <button
              v-for="page in Math.ceil(total / pageSize)"
              :key="page"
              @click="onPageChange(page)"
              :class="[
                'rounded px-3 py-1 text-sm',
                currentPage === page
                  ? 'bg-blue-600 text-white'
                  : 'border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 dark:border-dk-muted dark:bg-dk-card dark:text-dk-text',
              ]"
            >
              {{ page }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Modal -->
    <CustomerFormModal
      v-if="showAddModal"
      :title="t('customer.add_title')"
      @close="showAddModal = false"
      @submit="onCustomerAdded"
    />

    <!-- Edit Modal -->
    <CustomerFormModal
      v-if="showEditModal"
      :title="t('customer.edit_title')"
      :customer="editingCustomer"
      @close="showEditModal = false"
      @submit="onCustomerUpdated"
    />

    <!-- Delete Confirm -->
    <ConfirmDialog
      v-model="showDeleteConfirm"
      :title="t('customer.delete_title')"
      :message="t('customer.delete_confirm')"
      @confirm="doDelete"
      danger
    />

    <!-- Toast -->
    <AppToast
      v-if="toast.show"
      :message="toast.message"
      :type="toast.type"
      @close="toast.show = false"
    />
  </div>
</template>
