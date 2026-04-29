<script setup>
import { ref, computed, onActivated } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUsersStore } from '@/stores/users'
import { useToastStore } from '@/stores/toast'
import { authApi } from '@/api/auth'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import { IconRefresh, IconChevronLeft, IconChevronRight, IconPlus, IconX } from '@tabler/icons-vue'

const { t } = useI18n()
const toast = useToastStore()
const usersStore = useUsersStore()

// 用户组列表
const groups = ref([])
const groupsLoading = ref(false)

// 选中的组及其成员
const selectedGroup = ref(null)
const groupMembers = ref([])
const groupMembersLoading = ref(false)
const groupMemberPage = ref(1)
const groupMemberPageSize = ref(20)
const groupMemberTotal = ref(0)
const groupMemberTotalPages = computed(() => Math.ceil(groupMemberTotal.value / groupMemberPageSize.value))

// 添加成员弹窗
const showAddMemberDialog = ref(false)
const addMemberSearch = ref('')
const addMemberSearchResults = ref([])
const addMemberLoading = ref(false)
const addMemberSearchLoading = ref(false)

// 确认弹窗
const showConfirmDialog = ref(false)
const confirmDialogConfig = ref({
  title: '',
  message: '',
  confirmText: '',
  cancelText: '',
  danger: false,
  onConfirm: null,
})

async function fetchGroups() {
  groupsLoading.value = true
  try {
    const res = await authApi.getGroups()
    if (res.errCode === 0) {
      groups.value = res.data.groups || []
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

defineExpose({ fetchGroups })

onActivated(() => fetchGroups())
</script>

<template>
  <div class="space-y-4">
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
