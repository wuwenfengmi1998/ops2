<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { usePageTitle } from '@/composables/usePageTitle'
import { customerApi } from '@/api/customer'
import AppToast from '@/components/AppToast.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const isEdit = computed(() => !!route.params.id)
const customerId = computed(() => isEdit.value ? parseInt(route.params.id) : null)

usePageTitle(isEdit.value ? 'customer.edit_title' : 'customer.add_title')

const form = ref({
  first_name: '',
  last_name: '',
  title: 'Unit',
  phones: [{ prefix: '853', phone: '', label: 'mobile', is_primary: true }],
  emails: [{ email: '', label: 'work', is_primary: true }],
  companies: [{ company_name: '', department: '', position: '', is_primary: true }],
})

const loading = ref(false)
const pageLoading = ref(false)
const errors = ref({})
const toast = ref({ show: false, message: '', type: 'success' })

const titleOptions = ['Unit', 'Mr', 'Ms']
const phoneLabels = ['mobile', 'work', 'home', 'other']
const emailLabels = ['work', 'personal', 'other']
const prefixOptions = ['853', '852', '86']

// 编辑模式加载数据
onMounted(async () => {
  if (isEdit.value) {
    pageLoading.value = true
    try {
      const res = await customerApi.get({ id: customerId.value })
      if (res.errCode === 0) {
        const data = res.data
        form.value = {
          first_name: data.customer.first_name || '',
          last_name: data.customer.last_name || '',
          title: data.customer.title || 'Unit',
          phones: data.phones?.length > 0 ? data.phones.map(p => ({
            prefix: p.prefix || '853',
            phone: p.phone || '',
            label: p.label || 'mobile',
            is_primary: p.is_primary || false,
          })) : [{ prefix: '853', phone: '', label: 'mobile', is_primary: true }],
          emails: data.emails?.length > 0 ? data.emails.map(e => ({
            email: e.email || '',
            label: e.label || 'work',
            is_primary: e.is_primary || false,
          })) : [{ email: '', label: 'work', is_primary: true }],
          companies: data.companies?.length > 0 ? data.companies.map(c => ({
            company_name: c.company_name || '',
            department: c.department || '',
            position: c.position || '',
            is_primary: c.is_primary || false,
          })) : [{ company_name: '', department: '', position: '', is_primary: true }],
        }
      } else {
        showToast(res.errMsg || t('message.error'), 'error')
        router.push('/customer')
      }
    } catch {
      showToast(t('message.error'), 'error')
      router.push('/customer')
    } finally {
      pageLoading.value = false
    }
  }
})

function addPhone() {
  form.value.phones.push({ prefix: '853', phone: '', label: 'mobile', is_primary: false })
}

function removePhone(index) {
  form.value.phones.splice(index, 1)
  if (form.value.phones.length === 1) {
    form.value.phones[0].is_primary = true
  }
}

function setPrimaryPhone(index) {
  form.value.phones.forEach((p, i) => (p.is_primary = i === index))
}

function addEmail() {
  form.value.emails.push({ email: '', label: 'work', is_primary: false })
}

function removeEmail(index) {
  form.value.emails.splice(index, 1)
  if (form.value.emails.length === 1) {
    form.value.emails[0].is_primary = true
  }
}

function setPrimaryEmail(index) {
  form.value.emails.forEach((e, i) => (e.is_primary = i === index))
}

function addCompany() {
  form.value.companies.push({ company_name: '', department: '', position: '', is_primary: false })
}

function removeCompany(index) {
  form.value.companies.splice(index, 1)
  if (form.value.companies.length === 1) {
    form.value.companies[0].is_primary = true
  }
}

function setPrimaryCompany(index) {
  form.value.companies.forEach((c, i) => (c.is_primary = i === index))
}

async function submit() {
  errors.value = {}

  // 基本验证（姓必填，名可选）
  if (!form.value.last_name?.trim()) {
    errors.value.last_name = t('validation.required')
  }

  // 过滤空数据
  const payload = {
    ...form.value,
    phones: form.value.phones.filter(p => p.phone.trim()),
    emails: form.value.emails.filter(e => e.email.trim()),
    companies: form.value.companies.filter(c => c.company_name.trim()),
  }

  if (Object.keys(errors.value).length > 0) return

  loading.value = true
  try {
    const api = isEdit.value ? customerApi.update : customerApi.add
    const params = isEdit.value ? { ...payload, id: customerId.value } : payload
    const res = await api(params)
    if (res.errCode === 0) {
      showToast(isEdit.value ? t('message.update_success') : t('message.add_success'), 'success')
      setTimeout(() => {
        router.push('/customer')
      }, 800)
    } else {
      errors.value.submit = res.errMsg || t('message.error')
    }
  } catch {
    errors.value.submit = t('message.error')
  } finally {
    loading.value = false
  }
}

function goBack() {
  router.back()
}

function showToast(message, type = 'success') {
  toast.value = { show: true, message, type }
  setTimeout(() => (toast.value.show = false), 3000)
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 p-6 dark:bg-dk-base">
    <div class="mx-auto max-w-3xl">
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
          <h1 class="text-2xl font-bold text-gray-900 dark:text-dk-text">
            {{ isEdit ? t('customer.edit_title') : t('customer.add_title') }}
          </h1>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="pageLoading" class="flex items-center justify-center py-12">
        <div class="h-8 w-8 animate-spin rounded-full border-2 border-gray-300 border-t-blue-600"></div>
      </div>

      <!-- Form -->
      <div v-else class="rounded-lg border border-gray-200 bg-white p-6 dark:border-dk-muted dark:bg-dk-card">
        <form @submit.prevent="submit" class="space-y-6">
          <!-- Basic Info -->
          <div>
            <h3 class="mb-4 text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('customer.basic_info') }}</h3>
            <div class="grid gap-4 sm:grid-cols-3">
              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-dk-subtle">{{ t('customer.salutation') }}</label>
                <select v-model="form.title" class="w-full rounded-lg border border-gray-300 px-3 py-2 dark:border-dk-muted dark:bg-dk-base dark:text-dk-text">
                  <option v-for="opt in titleOptions" :key="opt" :value="opt">{{ t(`customer.salutation_${opt.toLowerCase()}`) }}</option>
                </select>
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-dk-subtle">{{ t('customer.last_name') }} *</label>
                <input v-model="form.last_name" type="text" class="w-full rounded-lg border border-gray-300 px-3 py-2 dark:border-dk-muted dark:bg-dk-base dark:text-dk-text" />
                <p v-if="errors.last_name" class="mt-1 text-xs text-red-500">{{ errors.last_name }}</p>
              </div>
              <div>
                <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-dk-subtle">{{ t('customer.first_name') }}</label>
                <input v-model="form.first_name" type="text" class="w-full rounded-lg border border-gray-300 px-3 py-2 dark:border-dk-muted dark:bg-dk-base dark:text-dk-text" />
                <p v-if="errors.first_name" class="mt-1 text-xs text-red-500">{{ errors.first_name }}</p>
              </div>
            </div>
          </div>

          <!-- Phones -->
          <div>
            <div class="mb-2 flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('customer.phones') }}</h3>
              <button type="button" @click="addPhone" class="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400">+ {{ t('customer.add_phone') }}</button>
            </div>
            <div v-for="(phone, index) in form.phones" :key="index" class="mb-2 flex items-center gap-2">
              <select v-model="phone.prefix" class="w-20 rounded-lg border border-gray-300 px-2 py-2 dark:border-dk-muted dark:bg-dk-base dark:text-dk-text">
                <option v-for="p in prefixOptions" :key="p" :value="p">+{{ p }}</option>
              </select>
              <input v-model="phone.phone" type="text" :placeholder="t('customer.phone_number')" class="flex-1 rounded-lg border border-gray-300 px-3 py-2 dark:border-dk-muted dark:bg-dk-base dark:text-dk-text" />
              <select v-model="phone.label" class="w-24 rounded-lg border border-gray-300 px-2 py-2 dark:border-dk-muted dark:bg-dk-base dark:text-dk-text">
                <option v-for="l in phoneLabels" :key="l" :value="l">{{ t(`customer.label_${l}`) }}</option>
              </select>
              <button type="button" @click="setPrimaryPhone(index)" :class="['rounded px-2 py-1 text-xs', phone.is_primary ? 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400' : 'text-gray-500']">
                {{ phone.is_primary ? t('customer.primary') : t('customer.set_primary') }}
              </button>
              <button v-if="form.phones.length > 1" type="button" @click="removePhone(index)" class="text-red-500 hover:text-red-600">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Emails -->
          <div>
            <div class="mb-2 flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('customer.emails') }}</h3>
              <button type="button" @click="addEmail" class="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400">+ {{ t('customer.add_email') }}</button>
            </div>
            <div v-for="(email, index) in form.emails" :key="index" class="mb-2 flex items-center gap-2">
              <input v-model="email.email" type="email" :placeholder="t('customer.email_address')" class="flex-1 rounded-lg border border-gray-300 px-3 py-2 dark:border-dk-muted dark:bg-dk-base dark:text-dk-text" />
              <select v-model="email.label" class="w-24 rounded-lg border border-gray-300 px-2 py-2 dark:border-dk-muted dark:bg-dk-base dark:text-dk-text">
                <option v-for="l in emailLabels" :key="l" :value="l">{{ t(`customer.label_${l}`) }}</option>
              </select>
              <button type="button" @click="setPrimaryEmail(index)" :class="['rounded px-2 py-1 text-xs', email.is_primary ? 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400' : 'text-gray-500']">
                {{ email.is_primary ? t('customer.primary') : t('customer.set_primary') }}
              </button>
              <button v-if="form.emails.length > 1" type="button" @click="removeEmail(index)" class="text-red-500 hover:text-red-600">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Companies -->
          <div>
            <div class="mb-2 flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('customer.companies') }}</h3>
              <button type="button" @click="addCompany" class="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400">+ {{ t('customer.add_company') }}</button>
            </div>
            <div v-for="(company, index) in form.companies" :key="index" class="mb-3 rounded-lg border border-gray-200 p-3 dark:border-dk-muted">
              <div class="mb-2 flex items-center justify-between">
                <button type="button" @click="setPrimaryCompany(index)" :class="['rounded px-2 py-1 text-xs', company.is_primary ? 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400' : 'text-gray-500']">
                  {{ company.is_primary ? t('customer.primary') : t('customer.set_primary') }}
                </button>
                <button v-if="form.companies.length > 1" type="button" @click="removeCompany(index)" class="text-red-500 hover:text-red-600">
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
              <div class="grid gap-2 sm:grid-cols-3">
                <input v-model="company.company_name" type="text" :placeholder="t('customer.company_name')" class="rounded-lg border border-gray-300 px-3 py-2 dark:border-dk-muted dark:bg-dk-base dark:text-dk-text" />
                <input v-model="company.department" type="text" :placeholder="t('customer.department')" class="rounded-lg border border-gray-300 px-3 py-2 dark:border-dk-muted dark:bg-dk-base dark:text-dk-text" />
                <input v-model="company.position" type="text" :placeholder="t('customer.position')" class="rounded-lg border border-gray-300 px-3 py-2 dark:border-dk-muted dark:bg-dk-base dark:text-dk-text" />
              </div>
            </div>
          </div>

          <p v-if="errors.submit" class="text-sm text-red-500">{{ errors.submit }}</p>

          <div class="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-dk-muted">
            <button type="button" @click="goBack" class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-dk-muted dark:text-dk-text">
              {{ t('common.cancel') }}
            </button>
            <button type="submit" :disabled="loading" class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50">
              {{ loading ? t('common.saving') : t('common.save') }}
            </button>
          </div>
        </form>
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
