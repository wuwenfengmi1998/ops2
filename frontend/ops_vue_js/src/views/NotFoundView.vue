<script setup>
import { useI18n } from 'vue-i18n'
import { usePageTitle } from '@/composables/usePageTitle'
import { IconAlertTriangle } from '@tabler/icons-vue'

usePageTitle('errorpage.404_title')
const { t } = useI18n()
</script>

<template>
  <div class="flex min-h-[60vh] flex-col items-center justify-center px-4">
    <div class="text-center">
      <IconAlertTriangle :size="64" class="mx-auto mb-4 text-yellow-400" stroke-width="1.5" />
      <h1 class="mb-2 text-5xl font-bold text-blue-600">404</h1>
      <h2 class="mb-3 text-2xl font-bold text-gray-900 dark:text-white">{{ t('errorpage.404_msg_title') }}</h2>
      <p class="mb-6 text-gray-500">{{ t('errorpage.404_msg') }}</p>
      <RouterLink to="/" class="inline-block rounded-lg bg-blue-600 px-5 py-2.5 font-medium text-white transition-colors hover:bg-blue-700">
        {{ t('errorpage.404_back_home') }}
      </RouterLink>
    </div>
  </div>
</template>
