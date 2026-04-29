<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { useUsersStore } from '@/stores/users'
import { authApi } from '@/api/auth'
import { useRouter } from 'vue-router'
import UsersTab from '@/views/sysadmin/UsersTab.vue'
import GroupsTab from '@/views/sysadmin/GroupsTab.vue'
import LogsTab from '@/views/sysadmin/LogsTab.vue'

const { t } = useI18n()
const userStore = useUserStore()
const usersStore = useUsersStore()
const router = useRouter()

const activeTab = ref('users')
const sysAdmins = ref([])
const loading = ref(false)

// 子组件 ref，用于主动调用其刷新方法
const usersTabRef = ref(null)
const groupsTabRef = ref(null)
const logsTabRef = ref(null)

const tabs = [
  { id: 'users', label: t('sysadmin.tab_users') },
  { id: 'groups', label: t('sysadmin.tab_groups') },
  { id: 'logs', label: t('sysadmin.tab_logs') },
  { id: 'customer', label: t('customer.title'), to: '/customer' },
]

async function fetchSysAdmins() {
  loading.value = true
  try {
    const res = await authApi.sysAdmins()
    if (res.errCode === 0 && Array.isArray(res.data.sysAdmins)) {
      sysAdmins.value = res.data.sysAdmins
      res.data.sysAdmins.forEach(id => usersStore.fetchUser(id))
    }
  } catch {
    // 错误已由拦截器处理
  } finally {
    loading.value = false
  }
}



onMounted(() => {
  fetchSysAdmins()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 p-6 dark:bg-dk-base">
    <div class="mx-auto max-w-6xl">
      <!-- Header -->
      <div class="mb-6 flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-dk-text">{{ t('sysadmin.title') }}</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-dk-subtle">
            {{ t('sysadmin.subtitle') }}
          </p>
        </div>
        <div class="flex items-center gap-2 rounded-lg bg-amber-100 px-3 py-1.5 dark:bg-amber-900/30">
          <span class="text-amber-700 dark:text-amber-400">{{ t('sysadmin.admin_label') }}: {{ userStore.user?.Username }}</span>
        </div>
      </div>

      <!-- Tabs -->
      <div class="mb-6 border-b border-gray-200 dark:border-dk-muted">
        <nav class="-mb-px flex gap-6">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="tab.to ? router.push(tab.to) : (activeTab = tab.id)"
            :class="[
              'border-b-2 px-1 pb-3 text-sm font-medium transition-colors',
              activeTab === tab.id && !tab.to
                ? 'border-blue-500 text-blue-600 dark:border-blue-400 dark:text-blue-400'
                : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 dark:text-dk-subtle dark:hover:text-dk-text',
            ]"
          >
            {{ tab.label }}
          </button>
        </nav>
      </div>

      <!-- Tab Content -->
      <div class="rounded-lg border border-gray-200 bg-white p-6 dark:border-dk-muted dark:bg-dk-card">
        <KeepAlive>
          <UsersTab v-if="activeTab === 'users'" ref="usersTabRef" />
          <GroupsTab v-else-if="activeTab === 'groups'" ref="groupsTabRef" />
          <LogsTab v-else-if="activeTab === 'logs'" ref="logsTabRef" />
        </KeepAlive>
      </div>

      <!-- SysAdmins List -->
      <div class="mt-6 rounded-lg border border-gray-200 bg-white p-4 dark:border-dk-muted dark:bg-dk-card">
        <div class="mb-2 flex items-center justify-between">
          <h3 class="text-sm font-medium text-gray-700 dark:text-dk-subtle">{{ t('sysadmin.current_admins') }}</h3>
          <button
            @click="fetchSysAdmins"
            class="text-xs text-blue-600 hover:text-blue-700 dark:text-blue-400"
            :disabled="loading"
          >
            {{ loading ? t('sysadmin.loading') : t('sysadmin.refresh') }}
          </button>
        </div>
        <div class="flex flex-wrap gap-2">
          <div
            v-for="adminId in sysAdmins"
            :key="adminId"
            class="flex items-center gap-2 rounded-full bg-amber-100 px-3 py-1 dark:bg-amber-900/30"
          >
            <img :src="usersStore.getAvatarUrlFromUserID(adminId)" class="w-5 h-5 rounded-full" alt="avatar" />
            <span class="text-xs font-medium text-amber-800 dark:text-amber-400">
              {{ usersStore.getUsernameFromUserID(adminId) || 'ID: ' + adminId }}
            </span>
          </div>
          <span v-if="sysAdmins.length === 0" class="text-sm text-gray-400 dark:text-dk-muted">
            {{ t('sysadmin.no_admins') }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
