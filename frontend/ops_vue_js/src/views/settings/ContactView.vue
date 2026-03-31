<script setup>
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { useValidation, isValidEmail } from '@/composables'
import { authApi } from '@/api/auth'
import SettingNav from '@/components/SettingNav.vue'

usePageTitle('settings.contact_information')
const { t } = useI18n()
const userStore = useUserStore()
const toast = useToastStore()
const { validate, errors, clearErrors } = useValidation()

const form = reactive({ email: '' })
const loading = ref(false)

async function handleChangeEmail() {
  clearErrors()

  const err = validate('email', form.email, t('message.please_enter_your_email'), isValidEmail)
  if (!err) return

  loading.value = true
  try {
    const { errCode } = await authApi.changeEmail(form.email)

    switch (errCode) {
      case 0:
        toast.success(t('message.change_ok'))
        await userStore.fetchUserInfo()
        break
      case -43:
        form.email = t('message.this_not_email')
        toast.error(t('message.this_not_email'))
        break
      default:
        toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-5xl px-6 py-6">
    <h2 class="mb-6 text-2xl font-bold text-gray-900 dark:text-white">{{ t('settings.my_account') }}</h2>
    <div class="flex flex-col gap-6 lg:flex-row lg:gap-8">
      <SettingNav />
      <div class="flex-1 space-y-6">
        <h3 class="mb-4 text-sm font-semibold uppercase text-gray-400 tracking-wider dark:text-gray-500">{{ t('settings.email') }}</h3>

        <div class="flex flex-col gap-4 sm:flex-row sm:items-start">
          <div class="flex-1">
            <input
              v-model="form.email"
              type="email"
              class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
              :class="errors.email ? 'border-red-500' : 'border-gray-300'"
              :placeholder="t('message.your_email_address')"
              @keydown.enter="handleChangeEmail"
            />
            <span v-if="errors.email" class="mt-1 block text-xs text-red-500">{{ errors.email }}</span>
          </div>
          <button
            class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500/20 focus:outline-none disabled:active:scale-100"
            :disabled="loading"
            @click="handleChangeEmail"
          >
            {{ t('settings.change_email') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
