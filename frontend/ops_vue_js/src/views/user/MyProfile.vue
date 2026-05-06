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

      <!-- My Groups -->
      <div class="mt-6">
        <h2 class="mb-3 text-lg font-semibold text-gray-900 dark:text-dk-text">
          {{ t('settings.my_groups') }}
        </h2>
        <div v-if="userStore.groups.length > 0" class="flex flex-wrap gap-2">
          <span
            v-for="group in userStore.groups"
            :key="group.id"
            class="inline-flex items-center rounded-full bg-purple-100 px-3 py-1 text-sm font-medium text-purple-800 dark:bg-purple-900/30 dark:text-purple-300"
          >
            <svg class="mr-1.5 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
            {{ group.name }}
          </span>
        </div>
        <div v-else class="rounded-lg border border-gray-200 bg-white p-4 text-center text-sm text-gray-500 dark:border-dk-muted dark:bg-dk-card dark:text-dk-subtle">
          {{ t('settings.no_groups') }}
        </div>
      </div>
    </div>
  </div>
</template>
