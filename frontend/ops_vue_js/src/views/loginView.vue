<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()

function toggleLocale() {
  locale.value = locale.value === 'zh-CN' ? 'en' : 'zh-CN'
}
import { useUserStore } from '@/stores/user'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { useValidation } from '@/composables'
import { authApi } from '@/api/auth'
import { IconEye, IconEyeOff } from '@tabler/icons-vue'

usePageTitle('appname.login')
const router = useRouter()
const userStore = useUserStore()
const toast = useToastStore()
const { validate, errors, clearErrors } = useValidation()

const form = ref({
  username: '',
  password: '',
  remember: false,
})
const showPassword = ref(false)
const loading = ref(false)

async function handleLogin() {
  clearErrors()

  const err1 = validate('username', form.value.username, t('message.please_enter_username_and_password'))
  const err2 = validate('password', form.value.password, t('message.please_enter_username_and_password'))

  if (!err1 || !err2) return

  loading.value = true
  try {
    const { errCode, data } = await authApi.login(form.value.username, form.value.password, form.value.remember)

    switch (errCode) {
      case 0:
        userStore.login(data.cookie)
        toast.success(t('message.login_successful'))
        const redirectPath = router.query.redirect || '/'
        router.push(redirectPath)
        break
      case -42:
        toast.danger(t('message.username_or_password_incorrect'))
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
  <div class="mx-auto max-w-sm px-8">
    <div class="mb-8 flex items-start justify-between">
      <RouterLink to="/" class="inline-flex items-center">
        <img src="/logo.svg" class="h-10 w-10 rounded-lg" alt="Operations" />
        <span class="ml-2.5 text-2xl font-bold text-gray-800 dark:text-dk-text">Operations</span>
      </RouterLink>
      <button class="rounded-md border border-gray-200 px-2.5 py-1 text-xs font-semibold uppercase text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:border-dk-muted dark:text-gray-400 dark:hover:bg-dk-card dark:hover:text-dk-text" @click="toggleLocale">
        {{ locale === 'zh-CN' ? 'EN' : '中' }}
      </button>
    </div>

    <div class="rounded-xl border border-gray-200 bg-white px-8 py-8 shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <h2 class="mb-6 text-center text-xl font-bold text-gray-900 dark:text-white">{{ t('message.login_to_your_account') }}</h2>

      <!-- Username -->
      <div class="mb-4">
        <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('message.user_name') }}</label>
        <input
          v-model="form.username"
          type="text"
          class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2.5 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
          :class="errors.username ? 'border-red-500 focus:border-red-500 focus:ring-red-500/20' : 'border-gray-300'"
          :placeholder="t('message.your_user_name')"
          @keydown.enter="handleLogin"
        />
        <span v-if="errors.username" class="mt-1 block text-xs text-red-500">{{ errors.username }}</span>
      </div>

      <!-- Password -->
      <div class="mb-4">
        <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('message.password') }}</label>
        <div class="relative">
          <input
            v-model="form.password"
            :type="showPassword ? 'text' : 'password'"
            class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2.5 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            :class="errors.password ? 'border-red-500 focus:border-red-500 focus:ring-red-500/20' : 'border-gray-300'"
            :placeholder="t('message.your_password')"
            autocomplete="current-password"
            @keydown.enter="handleLogin"
          />
          <button type="button" class="absolute inset-y-0 right-0 flex items-center px-3 text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300" @click="showPassword = !showPassword">
            <IconEye v-if="!showPassword" :size="18" />
            <IconEyeOff v-else :size="18" />
          </button>
        </div>
        <span v-if="errors.password" class="mt-1 block text-xs text-red-500">{{ errors.password }}</span>
      </div>

      <!-- Remember -->
      <label class="mb-6 flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
        <input v-model="form.remember" type="checkbox" class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500" />
        {{ t('message.remember_me_on_this_device') }}
      </label>

      <!-- Submit -->
      <button
        class="flex w-full items-center justify-center gap-2 rounded-lg bg-blue-600 px-4 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500/20 focus:outline-none disabled:active:scale-100"
        :disabled="loading"
        @click="handleLogin"
      >
        <svg v-if="loading" class="h-4 w-4 animate-spin text-white" viewBox="0 0 24 24" fill="none">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
        {{ t('button.sign_in') }}
      </button>
    </div>

    <p class="mt-6 text-center text-sm text-gray-500">
      {{ t('message.dont_have_account_yet') }}
      <RouterLink to="/register" class="font-medium text-blue-600 hover:text-blue-500 dark:text-blue-400">{{ t('message.register_now') }}</RouterLink>
    </p>
  </div>
</template>
