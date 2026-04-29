<script setup>
import { ref, computed, onMounted, onActivated } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUsersStore } from '@/stores/users'
import { useToastStore } from '@/stores/toast'
import { authApi } from '@/api/auth'
import { IconSearch, IconRefresh, IconChevronLeft, IconChevronRight } from '@tabler/icons-vue'

const { t } = useI18n()
const toast = useToastStore()
const usersStore = useUsersStore()

// 用户列表
const users = ref([])
const usersLoading = ref(false)
const userSearch = ref('')
const userPage = ref(1)
const userPageSize = ref(20)
const userTotal = ref(0)
const totalPages = computed(() => Math.ceil(userTotal.value / userPageSize.value))

// 用户详情弹窗
const showUserDetail = ref(false)
const userDetail = ref(null)
const userDetailInfo = ref(null)
const userDetailLoading = ref(false)
const newPassword = ref('')
const resetPasswordLoading = ref(false)

async function fetchUsers() {
  usersLoading.value = true
  try {
    const res = await authApi.getUsers({
      page: userPage.value,
      page_size: userPageSize.value,
      search: userSearch.value,
    })
    if (res.errCode === 0) {
      users.value = res.data.users || []
      userTotal.value = res.data.total || 0
      userPage.value = res.data.page || 1
      userPageSize.value = res.data.page_size || 20
      users.value.forEach(u => usersStore.fetchUser(u.id))
    }
  } catch {
    // 错误已由拦截器处理
  } finally {
    usersLoading.value = false
  }
}

function onSearch() {
  userPage.value = 1
  fetchUsers()
}

function onPageChange(page) {
  userPage.value = page
  fetchUsers()
}

async function openUserDetail(user) {
  userDetail.value = user
  showUserDetail.value = true
  userDetailLoading.value = true
  try {
    const res = await authApi.getUserDetail(user.id)
    if (res.errCode === 0) {
      userDetail.value = res.data.user || user
      userDetailInfo.value = res.data.userinfo || null
    }
  } catch {
    // 错误已由拦截器处理
  } finally {
    userDetailLoading.value = false
  }
}

function closeUserDetail() {
  showUserDetail.value = false
  userDetail.value = null
  userDetailInfo.value = null
  newPassword.value = ''
}

async function resetUserPassword() {
  if (!newPassword.value || newPassword.value.length < 6) {
    toast.warning(t('sysadmin.new_password_placeholder'))
    return
  }
  if (!userDetail.value) return
  resetPasswordLoading.value = true
  try {
    const res = await authApi.resetUserPassword(userDetail.value.id, newPassword.value)
    if (res.errCode === 0) {
      toast.success(t('message.change_ok'))
      newPassword.value = ''
    } else {
      toast.error(res.raw?.err_msg || t('message.change_ok'))
    }
  } catch {
    // 错误已由拦截器处理
  } finally {
    resetPasswordLoading.value = false
  }
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString()
}

function formatGender(gender) {
  const map = { 'M': t('settings.male'), 'F': t('settings.female'), 'U': '-' }
  return map[gender] || '-'
}

defineExpose({ fetchUsers })

onMounted(() => fetchUsers())
onActivated(() => fetchUsers())
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('sysadmin.tab_users') }}</h2>
      <span class="text-sm text-gray-500 dark:text-dk-subtle">
        {{ t('sysadmin.total_users', { count: userTotal }) }}
      </span>
    </div>

    <!-- 搜索栏 -->
    <div class="flex gap-2">
      <div class="relative flex-1">
        <IconSearch class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
        <input
          v-model="userSearch"
          type="text"
          :placeholder="t('sysadmin.search_placeholder')"
          class="w-full rounded-md border border-gray-300 py-2 pl-9 pr-4 text-sm focus:border-blue-500 focus:outline-none dark:border-dk-muted dark:bg-dk-base dark:text-dk-text"
          @keyup.enter="onSearch"
        />
      </div>
      <button
        @click="onSearch"
        class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
      >
        {{ t('sysadmin.search') }}
      </button>
      <button
        @click="fetchUsers"
        class="rounded-md border border-gray-300 px-3 py-2 text-gray-600 hover:bg-gray-50 dark:border-dk-muted dark:text-dk-subtle dark:hover:bg-dk-card"
        :disabled="usersLoading"
      >
        <IconRefresh :size="18" :class="{ 'animate-spin': usersLoading }" />
      </button>
    </div>

    <!-- 用户列表 -->
    <div class="overflow-hidden rounded-md border border-gray-200 dark:border-dk-muted">
      <table class="min-w-full divide-y divide-gray-200 dark:divide-dk-muted">
        <thead class="bg-gray-50 dark:bg-dk-base">
          <tr>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_id') }}</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_username') }}</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_email') }}</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_type') }}</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_created_at') }}</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_action') }}</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 bg-white dark:divide-dk-muted dark:bg-dk-card">
          <tr v-if="usersLoading" class="text-center">
            <td colspan="6" class="py-8 text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.loading') }}</td>
          </tr>
          <tr v-else-if="users.length === 0" class="text-center">
            <td colspan="6" class="py-8 text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.no_users') }}</td>
          </tr>
          <tr v-for="user in users" :key="user.id" class="hover:bg-gray-50 dark:hover:bg-dk-base">
            <td class="whitespace-nowrap px-4 py-3 text-sm text-gray-900 dark:text-dk-text">{{ user.id }}</td>
            <td class="whitespace-nowrap px-4 py-3">
              <div class="flex items-center gap-2">
                <img
                  :src="usersStore.getAvatarUrlFromUserID(user.id)"
                  class="h-7 w-7 rounded-full object-cover"
                  alt="avatar"
                />
                <span class="text-sm text-gray-900 dark:text-dk-text">
                  {{ usersStore.getUsernameFromUserID(user.id) || user.name }}
                </span>
              </div>
            </td>
            <td class="whitespace-nowrap px-4 py-3 text-sm text-gray-500 dark:text-dk-subtle">{{ user.email }}</td>
            <td class="whitespace-nowrap px-4 py-3 text-sm">
              <span
                :class="[
                  'rounded-full px-2 py-0.5 text-xs font-medium',
                  user.type === 'admin' ? 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400' : 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
                ]"
              >
                {{ user.type }}
              </span>
            </td>
            <td class="whitespace-nowrap px-4 py-3 text-sm text-gray-500 dark:text-dk-subtle">{{ new Date(user.date).toLocaleString() }}</td>
            <td class="whitespace-nowrap px-4 py-3 text-sm">
              <button @click="openUserDetail(user)" class="text-blue-600 hover:text-blue-700 dark:text-blue-400">{{ t('sysadmin.detail') }}</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 分页 -->
    <div class="flex items-center justify-between">
      <div class="text-sm text-gray-500 dark:text-dk-subtle">
        {{ t('sysadmin.pagination', { current: userPage, total: totalPages }) }}
      </div>
      <div class="flex gap-2">
        <button
          @click="onPageChange(userPage - 1)"
          :disabled="userPage <= 1 || usersLoading"
          class="flex items-center gap-1 rounded-md border border-gray-300 px-3 py-1.5 text-sm disabled:opacity-50 dark:border-dk-muted dark:text-dk-text"
        >
          <IconChevronLeft :size="16" /> {{ t('sysadmin.prev_page') }}
        </button>
        <button
          @click="onPageChange(userPage + 1)"
          :disabled="userPage >= totalPages || usersLoading"
          class="flex items-center gap-1 rounded-md border border-gray-300 px-3 py-1.5 text-sm disabled:opacity-50 dark:border-dk-muted dark:text-dk-text"
        >
          {{ t('sysadmin.next_page') }} <IconChevronRight :size="16" />
        </button>
      </div>
    </div>
  </div>

  <!-- 用户详情弹窗 -->
  <div
    v-if="showUserDetail"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
    @click.self="closeUserDetail"
  >
    <div class="w-full max-w-lg rounded-lg bg-white p-6 shadow-xl dark:bg-dk-card">
      <div class="mb-4 flex items-center justify-between">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('sysadmin.user_detail') }}</h3>
        <button
          @click="closeUserDetail"
          class="text-gray-400 hover:text-gray-600 dark:text-dk-subtle dark:hover:text-dk-text"
        >
          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div v-if="userDetailLoading" class="py-8 text-center text-gray-500 dark:text-dk-subtle">
        {{ t('sysadmin.loading') }}
      </div>

      <div v-else-if="userDetail" class="space-y-4">
        <!-- 用户头像和基本信息 -->
        <div class="flex items-center gap-4">
          <img
            :src="usersStore.getAvatarUrlFromUserID(userDetail.id)"
            class="h-16 w-16 rounded-full object-cover"
            alt="avatar"
          />
          <div>
            <div class="text-lg font-semibold text-gray-900 dark:text-dk-text">
              {{ usersStore.getUsernameFromUserID(userDetail.id) || userDetail.name }}
            </div>
            <div class="text-sm text-gray-500 dark:text-dk-subtle">{{ userDetail.email }}</div>
            <span
              :class="[
                'mt-1 inline-block rounded-full px-2 py-0.5 text-xs font-medium',
                userDetail.type === 'admin' ? 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400' : 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
              ]"
            >
              {{ userDetail.type }}
            </span>
          </div>
        </div>

        <hr class="border-gray-200 dark:border-dk-muted" />

        <!-- 详细信息 -->
        <div class="space-y-2 text-sm">
          <div class="flex justify-between">
            <span class="text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.user_id') }}</span>
            <span class="text-gray-900 dark:text-dk-text">{{ userDetail.id }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.name') }}</span>
            <span class="text-gray-900 dark:text-dk-text">{{ userDetail.name }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_created_at') }}</span>
            <span class="text-gray-900 dark:text-dk-text">{{ new Date(userDetail.date).toLocaleString() }}</span>
          </div>

          <!-- 用户扩展信息 -->
          <template v-if="userDetailInfo">
            <hr class="border-gray-200 dark:border-dk-muted" />
            <div class="text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.extended_info') }}</div>
            <div class="flex justify-between">
              <span class="text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.info_nickname') }}</span>
              <span class="text-gray-900 dark:text-dk-text">{{ userDetailInfo.username || '-' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.info_remark') }}</span>
              <span class="text-gray-900 dark:text-dk-text">{{ userDetailInfo.firstname || '-' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.birthday') }}</span>
              <span class="text-gray-900 dark:text-dk-text">{{ formatDate(userDetailInfo.birthdate) }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.gender') }}</span>
              <span class="text-gray-900 dark:text-dk-text">{{ formatGender(userDetailInfo.gender) }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.region') }}</span>
              <span class="text-gray-900 dark:text-dk-text">{{ userDetailInfo.region || '-' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.language') }}</span>
              <span class="text-gray-900 dark:text-dk-text">{{ userDetailInfo.language || '-' }}</span>
            </div>
          </template>
        </div>
      </div>

      <!-- 修改密码区域 -->
      <div class="mt-4 space-y-3 border-t border-gray-200 pt-4 dark:border-dk-muted">
        <div class="text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.reset_password') }}</div>
        <div class="flex gap-2">
          <input
            v-model="newPassword"
            type="password"
            :placeholder="t('sysadmin.new_password_placeholder')"
            class="flex-1 rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none dark:border-dk-muted dark:bg-dk-base dark:text-dk-text"
          />
          <button
            @click="resetUserPassword"
            :disabled="resetPasswordLoading || !newPassword"
            class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
          >
            {{ resetPasswordLoading ? t('sysadmin.resetting') : t('sysadmin.reset_password_btn') }}
          </button>
        </div>
      </div>

      <div class="mt-6 flex justify-end">
        <button
          @click="closeUserDetail"
          class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-dk-muted dark:text-dk-text dark:hover:bg-dk-base"
        >
          {{ t('sysadmin.close') }}
        </button>
      </div>
    </div>
  </div>
</template>
