<script setup>
import { useI18n } from 'vue-i18n'
import { usePageTitle } from '@/composables/usePageTitle'

usePageTitle('appname.warehouse')
const { t } = useI18n()
</script>

<template>
  <div class="mx-auto max-w-6xl px-6 py-6">
    <h2 class="mb-6 text-2xl font-bold text-gray-900 dark:text-white">{{ t('appname.warehouse') }}</h2>
    <div class="rounded-xl border border-gray-200 bg-white px-12 py-12 text-center shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <p class="text-gray-400">{{ t('message.functionality_not_yet_developed') }}</p>
    </div>
  </div>
</template>
