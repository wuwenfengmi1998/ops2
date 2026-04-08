<script setup>
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { usePageTitle } from '@/composables/usePageTitle'

usePageTitle('appname.home')
const { t } = useI18n()
const userStore = useUserStore()

const features = computed(() => [
  // { title: t('appname.schedule'), desc: '—' },
  // { title: t('appname.purchase'), desc: '—' },
  // { title: t('appname.warehouse'), desc: '—' },
])

import { computed } from 'vue'
</script>

<template>
  <div class="mx-auto max-w-6xl px-6 py-6">
    <h2 class="mb-6 text-2xl font-bold text-gray-900 dark:text-white">{{ t('message.welcome') }}</h2>

    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div
        v-for="feature in features"
        :key="feature.title"
        class="rounded-xl border border-gray-200 bg-white px-5 py-4 transition-shadow hover:shadow-md dark:border-dk-muted dark:bg-dk-card"
      >
        <p class="mb-1 text-sm text-gray-500">{{ feature.title }}</p>
        <p class="text-lg font-bold text-gray-900 dark:text-white">{{ feature.desc }}</p>
      </div>
    </div>
  </div>
</template>
