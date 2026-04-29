<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { customerApi } from '@/api/customer'
import { useUsersStore } from '@/stores/users'
import AppToast from '@/components/AppToast.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const usersStore = useUsersStore()

const customer = ref(null)
const phones = ref([])
const emails = ref([])
const companies = ref([])
const loading = ref(false)

const toast = ref({ show: false, message: '', type: 'success' })

// 获取客户详情
async function fetchCustomerDetail() {
  const id = route.params.id
  if (!id) {
    router.push('/customer')
    return
  }

  loading.value = true
  try {
    const res = await customerApi.get({ id: parseInt(id) })
    if (res.errCode === 0) {
      customer.value = res.data.customer
      phones.value = res.data.phones || []
      emails.value = res.data.emails || []
      companies.value = res.data.companies || []
      // 预加载创建者信息
      if (customer.value?.created_by) {
        usersStore.fetchUser(customer.value.created_by)
      }
    } else {
      showToast(res.errMsg || t('message.error'), 'error')
      router.push('/customer')
    }
  } catch {
    showToast(t('message.error'), 'error')
    router.push('/customer')
  } finally {
    loading.value = false
  }
}

// 返回列表
function goBack() {
  router.push('/customer')
}

// 格式化日期
function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString()
}

// 获取标签文本
function getLabelText(label) {
  const key = `customer.label_${label}`
  const text = t(key)
  return text === key ? label : text
}

function showToast(message, type = 'success') {
  toast.value = { show: true, message, type }
  setTimeout(() => (toast.value.show = false), 3000)
}

onMounted(() => {
  fetchCustomerDetail()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 p-6 dark:bg-dk-base">
    <div class="mx-auto max-w-4xl">
      <!-- Header -->
      <div class="mb-6 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <button
            @click="goBack"
            class="flex items-center gap-1 text-gray-600 hover:text-gray-900 dark:text-dk-subtle dark:hover:text-dk-text"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            {{ t('common.back') }}
          </button>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-dk-text">{{ t('customer.detail_title') }}</h1>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-gray-300 border-t-blue-600"></div>
      </div>

      <!-- Content -->
      <div v-else-if="customer" class="space-y-6">
        <!-- Basic Info Card -->
        <div class="rounded-lg border border-gray-200 bg-white p-6 dark:border-dk-muted dark:bg-dk-card">
          <h2 class="mb-4 text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('customer.basic_info') }}</h2>
          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="text-sm text-gray-500 dark:text-dk-subtle">{{ t('customer.name') }}</label>
              <p class="text-lg font-medium text-gray-900 dark:text-dk-text">
                {{ (customer.last_name || '') + (customer.first_name ? ' ' + customer.first_name : '') }}
              </p>
            </div>
            <div>
              <label class="text-sm text-gray-500 dark:text-dk-subtle">{{ t('customer.salutation') }}</label>
              <p class="text-gray-900 dark:text-dk-text">
                {{ customer.title ? t(`customer.salutation_${customer.title.toLowerCase()}`) : '-' }}
              </p>
            </div>
            <div>
              <label class="text-sm text-gray-500 dark:text-dk-subtle">{{ t('customer.created_at') }}</label>
              <p class="text-gray-900 dark:text-dk-text">{{ formatDate(customer.created_at) }}</p>
            </div>
            <div>
              <label class="text-sm text-gray-500 dark:text-dk-subtle">{{ t('customer.created_by') }}</label>
              <div class="flex items-center gap-2">
                <img
                  :src="usersStore.getAvatarUrlFromUserID(customer.created_by)"
                  class="h-6 w-6 rounded-full"
                  alt="avatar"
                />
                <span class="text-gray-900 dark:text-dk-text">
                  {{ usersStore.getUsernameFromUserID(customer.created_by) || 'ID: ' + customer.created_by }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Phones Card -->
        <div v-if="phones.length > 0" class="rounded-lg border border-gray-200 bg-white p-6 dark:border-dk-muted dark:bg-dk-card">
          <h2 class="mb-4 text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('customer.phones') }}</h2>
          <div class="space-y-3">
            <div
              v-for="phone in phones"
              :key="phone.id"
              class="flex items-center justify-between rounded-lg border border-gray-100 bg-gray-50 p-3 dark:border-dk-muted dark:bg-dk-muted"
            >
              <div class="flex items-center gap-3">
                <span class="text-gray-500 dark:text-dk-subtle">+{{ phone.prefix }}</span>
                <span class="font-medium text-gray-900 dark:text-dk-text">{{ phone.phone }}</span>
                <span
                  class="rounded px-2 py-0.5 text-xs"
                  :class="phone.is_primary ? 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400' : 'bg-gray-200 text-gray-600 dark:bg-dk-muted dark:text-dk-subtle'"
                >
                  {{ getLabelText(phone.label) }}
                  <span v-if="phone.is_primary" class="ml-1">({{ t('customer.primary') }})</span>
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Emails Card -->
        <div v-if="emails.length > 0" class="rounded-lg border border-gray-200 bg-white p-6 dark:border-dk-muted dark:bg-dk-card">
          <h2 class="mb-4 text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('customer.emails') }}</h2>
          <div class="space-y-3">
            <div
              v-for="email in emails"
              :key="email.id"
              class="flex items-center justify-between rounded-lg border border-gray-100 bg-gray-50 p-3 dark:border-dk-muted dark:bg-dk-muted"
            >
              <div class="flex items-center gap-3">
                <span class="font-medium text-gray-900 dark:text-dk-text">{{ email.email }}</span>
                <span
                  class="rounded px-2 py-0.5 text-xs"
                  :class="email.is_primary ? 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400' : 'bg-gray-200 text-gray-600 dark:bg-dk-muted dark:text-dk-subtle'"
                >
                  {{ getLabelText(email.label) }}
                  <span v-if="email.is_primary" class="ml-1">({{ t('customer.primary') }})</span>
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Companies Card -->
        <div v-if="companies.length > 0" class="rounded-lg border border-gray-200 bg-white p-6 dark:border-dk-muted dark:bg-dk-card">
          <h2 class="mb-4 text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('customer.companies') }}</h2>
          <div class="space-y-3">
            <div
              v-for="company in companies"
              :key="company.id"
              class="rounded-lg border border-gray-100 bg-gray-50 p-3 dark:border-dk-muted dark:bg-dk-muted"
            >
              <div class="flex items-center justify-between">
                <span class="font-medium text-gray-900 dark:text-dk-text">{{ company.company_name }}</span>
                <span
                  v-if="company.is_primary"
                  class="rounded px-2 py-0.5 text-xs bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400"
                >
                  {{ t('customer.primary') }}
                </span>
              </div>
              <div v-if="company.department || company.position" class="mt-1 text-sm text-gray-500 dark:text-dk-subtle">
                <span v-if="company.department">{{ company.department }}</span>
                <span v-if="company.department && company.position"> · </span>
                <span v-if="company.position">{{ company.position }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Not Found -->
      <div v-else class="rounded-lg border border-gray-200 bg-white p-12 text-center dark:border-dk-muted dark:bg-dk-card">
        <p class="text-gray-500 dark:text-dk-subtle">{{ t('customer.not_found') }}</p>
        <button
          @click="goBack"
          class="mt-4 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >
          {{ t('common.back') }}
        </button>
      </div>
    </div>

    <!-- Toast -->
    <AppToast
      v-if="toast.show"
      :message="toast.message"
      :type="toast.type"
      @close="toast.show = false"
    />
  </div>
</template>
