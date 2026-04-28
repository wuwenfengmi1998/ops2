<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { useUsersStore } from '@/stores/users'
import { useToastStore } from '@/stores/toast'
import { authApi } from '@/api/auth'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { IconSearch, IconRefresh, IconChevronLeft, IconChevronRight, IconPlus, IconX } from '@tabler/icons-vue'

const { t } = useI18n()

const toast = useToastStore()

const usersStore = useUsersStore()

const userStore = useUserStore()
const activeTab = ref('users')
const sysAdmins = ref([])
const loading = ref(false)

// 用户列表相关
const users = ref([])
const usersLoading = ref(false)
const userSearch = ref('')
const userPage = ref(1)
const userPageSize = ref(20)
const userTotal = ref(0)

// 用户组相关
const groups = ref([])
const groupsLoading = ref(false)
const selectedGroup = ref(null)
const groupMembers = ref([])
const groupMembersLoading = ref(false)
const groupMemberPage = ref(1)
const groupMemberPageSize = ref(20)
const groupMemberTotal = ref(0)

// 登录失败日志相关
const loginFailLogs = ref([])
const loginFailLogsLoading = ref(false)
const loginFailLogSearch = ref('')
const loginFailLogPage = ref(1)
const loginFailLogPageSize = ref(20)
const loginFailLogTotal = ref(0)

// 用户详情相关
const showUserDetail = ref(false)
const userDetail = ref(null)
const userDetailInfo = ref(null)
const userDetailLoading = ref(false)
const newPassword = ref('')
const resetPasswordLoading = ref(false)

// 确认弹窗相关
const showConfirmDialog = ref(false)
const confirmDialogConfig = ref({
  title: '确认',
  message: '',
  confirmText: '确认',
  cancelText: '取消',
  danger: false,
  onConfirm: null,
})

const tabs = [
  { id: 'users', label: t('sysadmin.tab_users') },
  { id: 'groups', label: t('sysadmin.tab_groups') },
  { id: 'logs', label: t('sysadmin.tab_logs') },
]

async function fetchSysAdmins() {
  loading.value = true
  try {
    const res = await authApi.sysAdmins()
    if (res.errCode === 0 && Array.isArray(res.data.sysAdmins)) {
      sysAdmins.value = res.data.sysAdmins
      // 预加载管理员用户信息（头像/用户名）
      res.data.sysAdmins.forEach(id => usersStore.fetchUser(id))
    }
  } catch {
    // 错误已由拦截器处理
  } finally {
    loading.value = false
  }
}

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
      // 预加载用户头像/用户名信息
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

const totalPages = computed(() => Math.ceil(userTotal.value / userPageSize.value))
const groupMemberTotalPages = computed(() => Math.ceil(groupMemberTotal.value / groupMemberPageSize.value))

// 添加成员弹窗相关
const showAddMemberDialog = ref(false)
const addMemberSearch = ref('')
const addMemberSearchResults = ref([])
const addMemberLoading = ref(false)
const addMemberSearchLoading = ref(false)
const loginFailLogTotalPages = computed(() => Math.ceil(loginFailLogTotal.value / loginFailLogPageSize.value))

async function fetchGroups() {
  groupsLoading.value = true
  try {
    const res = await authApi.getGroups()
    if (res.errCode === 0) {
      groups.value = res.data.groups || []
      // 预加载用户组成员信息
      groups.value.forEach(g => {
        g.memberIDs?.forEach(id => usersStore.fetchUser(id))
      })
    }
  } catch {
    // 错误已由拦截器处理
  } finally {
    groupsLoading.value = false
  }
}

async function fetchGroupMembers() {
  if (!selectedGroup.value) return
  groupMembersLoading.value = true
  try {
    const res = await authApi.getGroupMembers(selectedGroup.value.id, {
      page: groupMemberPage.value,
      page_size: groupMemberPageSize.value,
    })
    if (res.errCode === 0) {
      groupMembers.value = res.data.members || []
      groupMemberTotal.value = res.data.total || 0
      groupMemberPage.value = res.data.page || 1
      groupMemberPageSize.value = res.data.page_size || 20
      // 预加载成员头像信息
      groupMembers.value.forEach(m => usersStore.fetchUser(m.id))
    }
  } catch {
    // 错误已由拦截器处理
  } finally {
    groupMembersLoading.value = false
  }
}

function selectGroup(group) {
  selectedGroup.value = group
  groupMemberPage.value = 1
  fetchGroupMembers()
}

function onGroupMemberPageChange(page) {
  groupMemberPage.value = page
  fetchGroupMembers()
}

async function openAddMemberDialog() {
  showAddMemberDialog.value = true
  addMemberSearch.value = ''
  addMemberSearchResults.value = []
}

function closeAddMemberDialog() {
  showAddMemberDialog.value = false
  addMemberSearch.value = ''
  addMemberSearchResults.value = []
}

async function searchUsersToAdd() {
  if (!addMemberSearch.value.trim()) {
    addMemberSearchResults.value = []
    return
  }
  addMemberSearchLoading.value = true
  try {
    const res = await authApi.getUsers({
      page: 1,
      page_size: 10,
      search: addMemberSearch.value,
    })
    if (res.errCode === 0) {
      // 过滤掉已经在组中的用户
      const existingMemberIds = new Set(groupMembers.value.map(m => m.id))
      addMemberSearchResults.value = (res.data.users || []).filter(u => !existingMemberIds.has(u.id))
    }
  } catch {
    // 错误已由拦截器处理
  } finally {
    addMemberSearchLoading.value = false
  }
}

async function addGroupMember(userId) {
  if (!selectedGroup.value) return
  addMemberLoading.value = true
  try {
    const res = await authApi.addGroupMember(selectedGroup.value.id, userId)
    if (res.errCode === 0) {
      toast.success(t('sysadmin.add_member') + t('message.save_ok'))
      fetchGroupMembers()
      // 从搜索结果中移除已添加的用户
      addMemberSearchResults.value = addMemberSearchResults.value.filter(u => u.id !== userId)
    } else {
      toast.error(res.raw?.err_msg || t('sysadmin.add_member') + t('message.save_ok'))
    }
  } catch {
    // 错误已由拦截器处理
  } finally {
    addMemberLoading.value = false
  }
}

function openConfirmDialog(config) {
  confirmDialogConfig.value = { ...confirmDialogConfig.value, ...config }
  showConfirmDialog.value = true
}

function handleConfirm() {
  if (confirmDialogConfig.value.onConfirm) {
    confirmDialogConfig.value.onConfirm()
  }
  showConfirmDialog.value = false
}

async function removeGroupMember(userId) {
  if (!selectedGroup.value) return
  openConfirmDialog({
    title: t('sysadmin.remove_member_title'),
    message: t('sysadmin.remove_member_confirm'),
    confirmText: t('sysadmin.remove_member'),
    danger: true,
    onConfirm: async () => {
      try {
        const res = await authApi.removeGroupMember(selectedGroup.value.id, userId)
        if (res.errCode === 0) {
          toast.success(t('sysadmin.remove_member_title') + t('message.delete_ok'))
          fetchGroupMembers()
        } else {
          toast.error(res.raw?.err_msg || t('sysadmin.remove_member') + t('message.delete_ok'))
        }
      } catch {
        // 错误已由拦截器处理
      }
    },
  })
}

async function fetchLoginFailLogs() {
  loginFailLogsLoading.value = true
  try {
    const res = await authApi.getLoginFailLogs({
      page: loginFailLogPage.value,
      page_size: loginFailLogPageSize.value,
      search: loginFailLogSearch.value,
    })
    if (res.errCode === 0) {
      loginFailLogs.value = res.data.logs || []
      loginFailLogTotal.value = res.data.total || 0
      loginFailLogPage.value = res.data.page || 1
      loginFailLogPageSize.value = res.data.page_size || 20
      // 预加载用户头像信息
      loginFailLogs.value.forEach(log => {
        if (log.user_id > 0) {
          usersStore.fetchUser(log.user_id)
        }
      })
    }
  } catch {
    // 错误已由拦截器处理
  } finally {
    loginFailLogsLoading.value = false
  }
}

function onLoginFailLogSearch() {
  loginFailLogPage.value = 1
  fetchLoginFailLogs()
}

function onLoginFailLogPageChange(page) {
  loginFailLogPage.value = page
  fetchLoginFailLogs()
}

function formatReason(reason) {
  const reasonMap = {
    'password_error': t('sysadmin.reason_password_error'),
    'user_not_found': t('sysadmin.reason_user_not_found'),
  }
  return reasonMap[reason] || reason
}

function getReasonClass(reason) {
  if (reason === 'password_error') {
    return 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400'
  }
  if (reason === 'user_not_found') {
    return 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
  }
  return 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
}

async function openUserDetail(user) {
  userDetail.value = user
  showUserDetail.value = true
  userDetailLoading.value = true
  try {
    // 获取用户详细信息
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

// 监听 Tab 切换
watch(activeTab, (tab) => {
  if (tab === 'users') {
    fetchUsers()
  } else if (tab === 'groups') {
    fetchGroups()
    if (!selectedGroup.value && groups.value.length > 0) {
      selectGroup(groups.value[0])
    }
  } else if (tab === 'logs') {
    fetchLoginFailLogs()
  }
})

onMounted(() => {
  fetchSysAdmins()
  fetchUsers()
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
            @click="activeTab = tab.id"
            :class="[
              'border-b-2 px-1 pb-3 text-sm font-medium transition-colors',
              activeTab === tab.id
                ? 'border-blue-500 text-blue-600 dark:border-blue-400 dark:text-blue-400'
                : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 dark:text-dk-subtle dark:hover:text-dk-text',
            ]"
          >
            {{ tab.label }}
          </button>
        </nav>
      </div>

      <!-- Content -->
      <div class="rounded-lg border border-gray-200 bg-white p-6 dark:border-dk-muted dark:bg-dk-card">
        <!-- 用户管理 -->
        <div v-if="activeTab === 'users'" class="space-y-4">
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

        <!-- 用户组 -->
        <div v-if="activeTab === 'groups'" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('sysadmin.tab_groups') }}</h2>
            <button
              @click="fetchGroups"
              class="rounded-md border border-gray-300 px-3 py-1.5 text-gray-600 hover:bg-gray-50 dark:border-dk-muted dark:text-dk-subtle dark:hover:bg-dk-card"
              :disabled="groupsLoading"
            >
              <IconRefresh :size="18" :class="{ 'animate-spin': groupsLoading }" />
            </button>
          </div>

          <div class="grid gap-4 lg:grid-cols-3">
            <!-- 用户组列表 -->
            <div class="lg:col-span-1">
              <div class="overflow-hidden rounded-md border border-gray-200 dark:border-dk-muted">
                <div class="bg-gray-50 px-4 py-2 text-sm font-medium text-gray-700 dark:bg-dk-base dark:text-dk-subtle">
                  {{ t('sysadmin.group_list') }}
                </div>
                <div v-if="groupsLoading" class="p-4 text-center text-gray-500 dark:text-dk-subtle">
                  {{ t('sysadmin.loading') }}
                </div>
                <div v-else-if="groups.length === 0" class="p-4 text-center text-gray-500 dark:text-dk-subtle">
                  {{ t('sysadmin.no_groups') }}
                </div>
                <div v-else class="divide-y divide-gray-200 dark:divide-dk-muted">
                  <button
                    v-for="group in groups"
                    :key="group.id"
                    @click="selectGroup(group)"
                    :class="[
                      'w-full px-4 py-3 text-left transition-colors',
                      selectedGroup?.id === group.id
                        ? 'bg-blue-50 dark:bg-blue-900/20'
                        : 'hover:bg-gray-50 dark:hover:bg-dk-base'
                    ]"
                  >
                    <div class="flex items-center">
                      <div>
                        <div class="font-medium text-gray-900 dark:text-dk-text">{{ group.name }}</div>
                      </div>
                    </div>
                  </button>
                </div>
              </div>
            </div>

            <!-- 组成员详情 -->
            <div class="lg:col-span-2">
              <div v-if="!selectedGroup" class="rounded-md bg-gray-50 p-8 text-center dark:bg-dk-base">
                <p class="text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.select_group_hint') }}</p>
              </div>
              <div v-else class="space-y-4">
                <div class="flex items-center justify-between">
                  <div>
                    <h3 class="font-semibold text-gray-900 dark:text-dk-text">{{ selectedGroup.name }}</h3>
                    <p class="text-sm text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.total_members', { count: groupMemberTotal }) }}</p>
                  </div>
                  <button
                    @click="openAddMemberDialog"
                    class="flex items-center gap-1 rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-blue-700"
                  >
                    <IconPlus :size="16" /> {{ t('sysadmin.add_member') }}
                  </button>
                </div>

                <div class="overflow-hidden rounded-md border border-gray-200 dark:border-dk-muted">
                  <table class="min-w-full divide-y divide-gray-200 dark:divide-dk-muted">
                    <thead class="bg-gray-50 dark:bg-dk-base">
                      <tr>
                        <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_user') }}</th>
                        <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_email') }}</th>
                        <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_type') }}</th>
                        <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_action') }}</th>
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-200 bg-white dark:divide-dk-muted dark:bg-dk-card">
                      <tr v-if="groupMembersLoading" class="text-center">
                        <td colspan="4" class="py-8 text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.loading') }}</td>
                      </tr>
                      <tr v-else-if="groupMembers.length === 0" class="text-center">
                        <td colspan="4" class="py-8 text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.no_members') }}</td>
                      </tr>
                      <tr v-for="member in groupMembers" :key="member.id" class="hover:bg-gray-50 dark:hover:bg-dk-base">
                        <td class="whitespace-nowrap px-4 py-3">
                          <div class="flex items-center gap-2">
                            <img
                              :src="usersStore.getAvatarUrlFromUserID(member.id)"
                              class="h-7 w-7 rounded-full object-cover"
                              alt="avatar"
                            />
                            <span class="text-sm text-gray-900 dark:text-dk-text">
                              {{ usersStore.getUsernameFromUserID(member.id) || member.name }}
                            </span>
                          </div>
                        </td>
                        <td class="whitespace-nowrap px-4 py-3 text-sm text-gray-500 dark:text-dk-subtle">{{ member.email }}</td>
                        <td class="whitespace-nowrap px-4 py-3 text-sm">
                          <span
                            :class="[
                              'rounded-full px-2 py-0.5 text-xs font-medium',
                              member.type === 'admin' ? 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400' : 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
                            ]"
                          >
                            {{ member.type }}
                          </span>
                        </td>
                        <td class="whitespace-nowrap px-4 py-3 text-sm">
                          <button
                            @click="removeGroupMember(member.id)"
                            class="text-red-600 hover:text-red-700 dark:text-red-400"
                          >
                            {{ t('sysadmin.remove_member') }}
                          </button>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <!-- 分页 -->
                <div class="flex items-center justify-between">
                  <div class="text-sm text-gray-500 dark:text-dk-subtle">
                    {{ t('sysadmin.pagination', { current: groupMemberPage, total: groupMemberTotalPages }) }}
                  </div>
                  <div class="flex gap-2">
                    <button
                      @click="onGroupMemberPageChange(groupMemberPage - 1)"
                      :disabled="groupMemberPage <= 1 || groupMembersLoading"
                      class="flex items-center gap-1 rounded-md border border-gray-300 px-3 py-1.5 text-sm disabled:opacity-50 dark:border-dk-muted dark:text-dk-text"
                    >
                      <IconChevronLeft :size="16" /> {{ t('sysadmin.prev_page') }}
                    </button>
                    <button
                      @click="onGroupMemberPageChange(groupMemberPage + 1)"
                      :disabled="groupMemberPage >= groupMemberTotalPages || groupMembersLoading"
                      class="flex items-center gap-1 rounded-md border border-gray-300 px-3 py-1.5 text-sm disabled:opacity-50 dark:border-dk-muted dark:text-dk-text"
                    >
                      {{ t('sysadmin.next_page') }} <IconChevronRight :size="16" />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 登录日志 -->
        <div v-if="activeTab === 'logs'" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('sysadmin.tab_logs') }}</h2>
            <span class="text-sm text-gray-500 dark:text-dk-subtle">
              {{ t('sysadmin.total_logs', { count: loginFailLogTotal }) }}
            </span>
          </div>

          <!-- 搜索栏 -->
          <div class="flex gap-2">
            <div class="relative flex-1">
              <IconSearch class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
              <input
                v-model="loginFailLogSearch"
                type="text"
                :placeholder="t('sysadmin.search_log_placeholder')"
                class="w-full rounded-md border border-gray-300 py-2 pl-9 pr-4 text-sm focus:border-blue-500 focus:outline-none dark:border-dk-muted dark:bg-dk-base dark:text-dk-text"
                @keyup.enter="onLoginFailLogSearch"
              />
            </div>
            <button
              @click="onLoginFailLogSearch"
              class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
            >
              {{ t('sysadmin.search') }}
            </button>
            <button
              @click="fetchLoginFailLogs"
              class="rounded-md border border-gray-300 px-3 py-2 text-gray-600 hover:bg-gray-50 dark:border-dk-muted dark:text-dk-subtle dark:hover:bg-dk-card"
              :disabled="loginFailLogsLoading"
            >
              <IconRefresh :size="18" :class="{ 'animate-spin': loginFailLogsLoading }" />
            </button>
          </div>

          <!-- 日志列表 -->
          <div class="overflow-hidden rounded-md border border-gray-200 dark:border-dk-muted">
            <table class="min-w-full divide-y divide-gray-200 dark:divide-dk-muted">
              <thead class="bg-gray-50 dark:bg-dk-base">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_user') }}</th>
                  <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_reason') }}</th>
                  <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_count') }}</th>
                  <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_ip') }}</th>
                  <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_updated_at') }}</th>
                  <th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.table_created') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200 bg-white dark:divide-dk-muted dark:bg-dk-card">
                <tr v-if="loginFailLogsLoading" class="text-center">
                  <td colspan="6" class="py-8 text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.loading') }}</td>
                </tr>
                <tr v-else-if="loginFailLogs.length === 0" class="text-center">
                  <td colspan="6" class="py-8 text-gray-500 dark:text-dk-subtle">{{ t('sysadmin.no_logs') }}</td>
                </tr>
                <tr v-for="log in loginFailLogs" :key="log.id" class="hover:bg-gray-50 dark:hover:bg-dk-base">
                  <td class="whitespace-nowrap px-4 py-3">
                    <div class="flex items-center gap-2">
                      <img
                        v-if="log.user_id > 0"
                        :src="usersStore.getAvatarUrlFromUserID(log.user_id)"
                        class="h-7 w-7 rounded-full object-cover"
                        alt="avatar"
                      />
                      <div v-else class="h-7 w-7 rounded-full bg-gray-200 dark:bg-gray-700"></div>
                      <span class="text-sm text-gray-900 dark:text-dk-text">
                        {{ log.username }}
                      </span>
                    </div>
                  </td>
                  <td class="whitespace-nowrap px-4 py-3">
                    <span
                      :class="[
                        'rounded-full px-2 py-0.5 text-xs font-medium',
                        getReasonClass(log.reason)
                      ]"
                    >
                      {{ formatReason(log.reason) }}
                    </span>
                  </td>
                  <td class="whitespace-nowrap px-4 py-3 text-sm">
                    <span :class="[
                      'font-medium',
                      log.count >= 5 ? 'text-red-600 dark:text-red-400' : 'text-gray-900 dark:text-dk-text'
                    ]">
                      {{ log.count }}
                    </span>
                  </td>
                  <td class="whitespace-nowrap px-4 py-3 text-sm font-mono text-gray-500 dark:text-dk-subtle">{{ log.ip }}</td>
                  <td class="whitespace-nowrap px-4 py-3 text-sm text-gray-500 dark:text-dk-subtle">{{ new Date(log.updated_at).toLocaleString() }}</td>
                  <td class="whitespace-nowrap px-4 py-3 text-sm text-gray-500 dark:text-dk-subtle">{{ new Date(log.created_at).toLocaleString() }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- 分页 -->
          <div class="flex items-center justify-between">
            <div class="text-sm text-gray-500 dark:text-dk-subtle">
              {{ t('sysadmin.pagination', { current: loginFailLogPage, total: loginFailLogTotalPages }) }}
            </div>
            <div class="flex gap-2">
              <button
                @click="onLoginFailLogPageChange(loginFailLogPage - 1)"
                :disabled="loginFailLogPage <= 1 || loginFailLogsLoading"
                class="flex items-center gap-1 rounded-md border border-gray-300 px-3 py-1.5 text-sm disabled:opacity-50 dark:border-dk-muted dark:text-dk-text"
              >
                <IconChevronLeft :size="16" /> {{ t('sysadmin.prev_page') }}
              </button>
              <button
                @click="onLoginFailLogPageChange(loginFailLogPage + 1)"
                :disabled="loginFailLogPage >= loginFailLogTotalPages || loginFailLogsLoading"
                class="flex items-center gap-1 rounded-md border border-gray-300 px-3 py-1.5 text-sm disabled:opacity-50 dark:border-dk-muted dark:text-dk-text"
              >
                {{ t('sysadmin.next_page') }} <IconChevronRight :size="16" />
              </button>
            </div>
          </div>
        </div>

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

  <!-- 添加成员弹窗 -->
  <div
    v-if="showAddMemberDialog"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
    @click.self="closeAddMemberDialog"
  >
    <div class="w-full max-w-md rounded-lg bg-white p-6 shadow-xl dark:bg-dk-card">
        <div class="mb-4 flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-dk-text">{{ t('sysadmin.add_member_title', { name: selectedGroup?.name }) }}</h3>
        <button
          @click="closeAddMemberDialog"
          class="text-gray-400 hover:text-gray-600 dark:text-dk-subtle dark:hover:text-dk-text"
        >
          <IconX :size="20" />
        </button>
      </div>

        <!-- 搜索框 -->
        <div class="mb-4">
          <div class="flex gap-2">
            <input
              v-model="addMemberSearch"
              type="text"
              :placeholder="t('sysadmin.search_user_placeholder')"
              class="flex-1 rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none dark:border-dk-muted dark:bg-dk-base dark:text-dk-text"
              @keyup.enter="searchUsersToAdd"
            />
            <button
              @click="searchUsersToAdd"
              :disabled="addMemberSearchLoading"
              class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            >
              {{ addMemberSearchLoading ? t('sysadmin.searching') : t('sysadmin.search') }}
            </button>
          </div>
        </div>

        <!-- 搜索结果 -->
        <div class="max-h-64 overflow-y-auto">
          <div v-if="addMemberSearchLoading" class="py-4 text-center text-gray-500 dark:text-dk-subtle">
            {{ t('sysadmin.searching') }}
          </div>
          <div v-else-if="addMemberSearchResults.length === 0 && addMemberSearch" class="py-4 text-center text-gray-500 dark:text-dk-subtle">
            {{ t('sysadmin.no_search_results') }}
          </div>
          <div v-else-if="!addMemberSearch" class="py-4 text-center text-gray-500 dark:text-dk-subtle">
            {{ t('sysadmin.search_hint') }}
          </div>
        <div v-else class="space-y-2">
          <div
            v-for="user in addMemberSearchResults"
            :key="user.id"
            class="flex items-center justify-between rounded-md border border-gray-200 p-3 dark:border-dk-muted"
          >
            <div class="flex items-center gap-3">
              <img
                :src="usersStore.getAvatarUrlFromUserID(user.id)"
                class="h-8 w-8 rounded-full object-cover"
                alt="avatar"
              />
              <div>
                <div class="text-sm font-medium text-gray-900 dark:text-dk-text">
                  {{ usersStore.getUsernameFromUserID(user.id) || user.name }}
                </div>
                <div class="text-xs text-gray-500 dark:text-dk-subtle">{{ user.email }}</div>
              </div>
            </div>
            <button
              @click="addGroupMember(user.id)"
              :disabled="addMemberLoading"
              class="rounded-md bg-green-600 px-3 py-1 text-xs font-medium text-white hover:bg-green-700 disabled:opacity-50"
            >
              {{ t('sysadmin.add') }}
            </button>
          </div>
        </div>
      </div>

      <div class="mt-4 flex justify-end">
        <button
          @click="closeAddMemberDialog"
          class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-dk-muted dark:text-dk-text dark:hover:bg-dk-base"
        >
          {{ t('sysadmin.close') }}
        </button>
      </div>
    </div>
  </div>

  <!-- 确认弹窗 -->
  <ConfirmDialog
    v-model="showConfirmDialog"
    :title="confirmDialogConfig.title"
    :message="confirmDialogConfig.message"
    :confirm-text="confirmDialogConfig.confirmText"
    :cancel-text="confirmDialogConfig.cancelText"
    :danger="confirmDialogConfig.danger"
    @confirm="handleConfirm"
  />
</template>
