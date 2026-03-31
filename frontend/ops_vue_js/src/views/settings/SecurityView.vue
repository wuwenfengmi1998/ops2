<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { useValidation } from '@/composables'
import { authApi } from '@/api/auth'
import { IconEye, IconEyeOff } from '@tabler/icons-vue'
import SettingNav from '@/components/SettingNav.vue'

usePageTitle('settings.security_settings')
const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()
const toast = useToastStore()
const { validate, errors, clearErrors } = useValidation()

const form = reactive({
  oldpass: '',
  newpass: '',
  confirm: '',
})
const showPassword = ref(false)
const loading = ref(false)

async function handleChangePassword() {
  clearErrors()

  const err1 = validate('oldpass', form.oldpass, t('message.type_old_pass'))
  const err2 = validate('newpass', form.newpass, t('message.type_new_pass'))
  const err3 = validate('confirm', form.confirm, t('message.type_cof_pass'))

  if (form.newpass !== form.confirm) {
    errors.confirm = t('message.confirm_password_incorrect')
    toast.warning(t('message.confirm_password_incorrect'))
  }

  if (!err1 || !err2 || !err3 || errors.confirm) return

  loading.value = true
  try {
    const { errCode } = await authApi.changePassword(form.oldpass, form.newpass)

    switch (errCode) {
      case 0:
        form.oldpass = ''
        form.newpass = ''
        form.confirm = ''
        toast.success(t('message.change_ok'), 2000)
        setTimeout(() => {
          userStore.logout()
          router.push('/')
        }, 2000)
        break
      case -42:
        form.oldpass = t('message.old_pass_incorrect')
        toast.error(t('message.old_pass_incorrect'))
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
        <div class="mb-4 flex items-center gap-2">
          <h3 class="text-sm font-semibold uppercase text-gray-400 tracking-wider dark:text-gray-500">{{ t('settings.password') }}</h3>
          <button type="button" class="rounded p-1 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 disabled:hover:bg-transparent dark:text-gray-500 dark:hover:text-gray-300" @click="showPassword = !showPassword">
            <IconEye v-if="!showPassword" :size="16" />
            <IconEyeOff v-else :size="16" />
          </button>
        </div>

        <div class="flex flex-col gap-4 sm:flex-row sm:items-start">
          <div>
            <input
              v-model="form.oldpass"
              :type="showPassword ? 'text' : 'password'"
              class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
              :class="errors.oldpass ? 'border-red-500' : 'border-gray-300'"
              :placeholder="t('message.type_old_pass')"
            />
            <span v-if="errors.oldpass" class="mt-1 block text-xs text-red-500">{{ errors.oldpass }}</span>
          </div>
          <div>
            <input
              v-model="form.newpass"
              :type="showPassword ? 'text' : 'password'"
              class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
              :class="errors.newpass ? 'border-red-500' : 'border-gray-300'"
              :placeholder="t('message.type_new_pass')"
            />
            <span v-if="errors.newpass" class="mt-1 block text-xs text-red-500">{{ errors.newpass }}</span>
          </div>
          <div>
            <input
              v-model="form.confirm"
              :type="showPassword ? 'text' : 'password'"
              class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
              :class="errors.confirm ? 'border-red-500' : 'border-gray-300'"
              :placeholder="t('message.type_cof_pass')"
              @keydown.enter="handleChangePassword"
            />
            <span v-if="errors.confirm" class="mt-1 block text-xs text-red-500">{{ errors.confirm }}</span>
          </div>
        </div>

        <div class="mt-4">
          <button
            class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500/20 focus:outline-none disabled:active:scale-100"
            :disabled="loading"
            @click="handleChangePassword"
          >
            {{ t('settings.set_new_password') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
