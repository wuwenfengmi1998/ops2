<script setup>
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { IconUser, IconMail, IconLock } from '@tabler/icons-vue'

const { t } = useI18n()

const links = [
  { to: '/settings/account', label: 'account_information', icon: IconUser },
  { to: '/settings/contact', label: 'contact_information', icon: IconMail },
  { to: '/settings/security', label: 'security_settings', icon: IconLock },
]
</script>

<template>
  <nav class="flex flex-col gap-1 border-b-4 border-b-blue-600 w-56 shrink-0 py-6">
    <RouterLink
      v-for="link in links"
      :key="link.to"
      :to="link.to"
      class="flex items-center gap-2.5 rounded-lg px-3 py-2 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-100 dark:text-dk-subtle dark:hover:bg-dk-card"
      active-class="!bg-blue-50 !text-blue-600 dark:!bg-dk-card dark:!text-blue-400"
    >
      <component :is="link.icon" :size="18" />
      {{ t(`settings.${link.label}`) }}
    </RouterLink>
  </nav>
</template>
