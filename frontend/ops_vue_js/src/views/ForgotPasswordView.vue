<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { useValidation } from '@/composables'
import { IconMail } from '@tabler/icons-vue'

usePageTitle('appname.forgot_password')
const { t } = useI18n()
const toast = useToastStore()
const { validate, errors, clearErrors } = useValidation()

const form = ref({ username: '' })
const loading = ref(false)

async function handleReset() {
  const err = validate('username', form.value.username, t('message.please_enter_your_username'))
  if (!err) return

  // 功能未开发
  toast.warning(t('message.functionality_not_yet_developed'))
}
</script>

<template>
  <div class="mx-auto max-w-sm px-8">
    <div class="mb-8 text-center">
      <RouterLink to="/">
        <img src="/static/logo.svg" width="110" height="32" alt="Operations" class="mx-auto" />
      </RouterLink>
    </div>

    <div class="rounded-xl border border-gray-200 bg-white px-8 py-8 shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <h2 class="mb-2 text-center text-xl font-bold text-gray-900 dark:text-white">{{ t('message.forgot_password') }}</h2>
      <p class="mb-6 text-center text-sm text-gray-500">{{ t('message.enter_your_username_to_reset_password') }}</p>

      <div class="mb-6">
        <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('message.user_name') }}</label>
        <input
          v-model="form.username"
          type="text"
          class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2.5 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
          :class="errors.username ? 'border-red-500 focus:border-red-500 focus:ring-red-500/20' : 'border-gray-300'"
          :placeholder="t('message.your_user_name')"
          @keydown.enter="handleReset"
        />
        <span v-if="errors.username" class="mt-1 block text-xs text-red-500">{{ errors.username }}</span>
      </div>

      <button
        class="flex w-full items-center justify-center gap-2 rounded-lg bg-blue-600 px-4 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500/20 focus:outline-none disabled:active:scale-60"
        :disabled="loading"
        @click="handleReset"
      >
        <IconMail :size="18" />
        {{ t('button.send_me_new_password') }}
      </button>
    </div>

    <p class="mt-6 text-center text-sm text-gray-500">
      <RouterLink to="/login" class="font-medium text-blue-600 hover:text-blue-500 dark:text-blue-400">{{ t('message.back_to_login') }}</RouterLink>
    </p>
  </div>
</template>
