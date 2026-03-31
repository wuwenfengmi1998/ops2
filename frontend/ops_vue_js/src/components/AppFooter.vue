<script setup>
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>

<template>
  <footer class="border-t border-gray-200 bg-white py-6 dark:border-dk-muted dark:bg-dk-base print:hidden">
    <div class="mx-auto flex max-w-6xl flex-col items-center justify-between gap-3 px-4 sm:flex-row sm:gap-0">
      <div class="flex items-center gap-3 text-sm text-gray-500">
        <a href="https://github.com/wuwenfengmi1998/ops2" target="_blank" class="hover:text-gray-700 dark:hover:text-dk-text" rel="noopener">{{ t('footer.github') }}</a>
        <span class="text-gray-300 dark:text-dk-muted">·</span>
        <a href="https://github.com/wuwenfengmi1998/ops2/blob/main/LICENSE" target="_blank" class="hover:text-gray-700 dark:hover:text-dk-text" rel="noopener">{{ t('footer.license') }}</a>
        <span class="text-gray-300 dark:text-dk-muted">·</span>
        <a href="https://github.com/wuwenfengmi1998" target="_blank" class="hover:text-gray-700 dark:hover:text-dk-text" rel="noopener">{{ t('footer.author_home') }}</a>
      </div>
      <div class="text-sm text-gray-400">{{ t('footer.copy') }}</div>
    </div>
  </footer>
</template>
