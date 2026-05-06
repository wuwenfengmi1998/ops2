<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { usePageTitle } from '@/composables/usePageTitle'
import { useUserStore } from '@/stores/user'

usePageTitle('message.profile_information')
const { t } = useI18n()
const userStore = useUserStore()

onMounted(() => {
  userStore.fetchUserInfo()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 p-6 dark:bg-dk-base">
    <div class="mx-auto max-w-4xl">
      <!-- Header -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-dk-text">
          {{ t('message.profile_information') }}
        </h1>
      </div>

      <!-- Profile Card -->
      <div class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dk-muted dark:bg-dk-card">
        <div class="flex flex-col items-center gap-6 md:flex-row md:items-start">
          <!-- Avatar -->
          <div class="flex-shrink-0">
            <img
              :src="userStore.avatarUrl"
              class="h-24 w-24 rounded-full border-4 border-white shadow-lg dark:border-dk-muted"
              alt="avatar"
            />
          </div>

          <!-- User Info -->
          <div class="flex-1 text-center md:text-left">
            <h2 class="text-xl font-semibold text-gray-900 dark:text-dk-text">
              {{ userStore.user?.Name || '-' }}
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-dk-subtle">
              {{ userStore.userInfo?.Username || '-' }}
            </p>

            <!-- Additional Info -->
            <div class="mt-4 grid grid-cols-1 gap-3 text-sm md:grid-cols-2">
              <div class="rounded-lg bg-gray-50 p-3 dark:bg-dk-base">
                <div class="text-xs text-gray-500 dark:text-dk-subtle">{{ t('settings.remark') }}</div>
                <div class="mt-1 font-medium text-gray-900 dark:text-dk-text">
                  {{ userStore.userInfo?.FirstName || '-' }}
                </div>
              </div>
              <div class="rounded-lg bg-gray-50 p-3 dark:bg-dk-base">
                <div class="text-xs text-gray-500 dark:text-dk-subtle">{{ t('settings.birthday') }}</div>
                <div class="mt-1 font-medium text-gray-900 dark:text-dk-text">
                  {{ userStore.birthday || '-' }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="mt-6 grid gap-4 md:grid-cols-2">
        <RouterLink
          to="/settings/account"
          class="flex items-center gap-3 rounded-xl border border-gray-200 bg-white p-4 shadow-sm transition-shadow hover:shadow-md dark:border-dk-muted dark:bg-dk-card"
        >
          <div class="rounded-full bg-blue-100 p-3 dark:bg-blue-900/30">
            <svg class="h-6 w-6 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </div>
          <div>
            <div class="font-medium text-gray-900 dark:text-dk-text">{{ t('message.user_settings') }}</div>
            <div class="text-sm text-gray-500 dark:text-dk-subtle">{{ t('settings.edit_profile') }}</div>
          </div>
        </RouterLink>

        <RouterLink
          to="/settings/security"
          class="flex items-center gap-3 rounded-xl border border-gray-200 bg-white p-4 shadow-sm transition-shadow hover:shadow-md dark:border-dk-muted dark:bg-dk-card"
        >
          <div class="rounded-full bg-green-100 p-3 dark:bg-green-900/30">
            <svg class="h-6 w-6 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <div>
            <div class="font-medium text-gray-900 dark:text-dk-text">{{ t('settings.security') }}</div>
            <div class="text-sm text-gray-500 dark:text-dk-subtle">{{ t('settings.security_description') }}</div>
          </div>
        </RouterLink>
      </div>
    </div>
  </div>
</template>
