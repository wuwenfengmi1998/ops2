<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { usePageTitle } from '@/composables/usePageTitle'
import { useToastStore } from '@/stores/toast'
import { calendarApi } from '@/api/calendar'
import { useUsersStore } from '@/stores/users'
import { IconCalendar, IconClock, IconEdit, IconRestore } from '@tabler/icons-vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

usePageTitle('calendar.admin_title')
const { t } = useI18n()
const toast = useToastStore()
const usersStore = useUsersStore()

const calendars = ref([])
const loading = ref(false)
const eventCounts = ref({}) // calendarId -> event count

// 编辑相关
const showEditModal = ref(false)
const editingCalendar = ref(null)
const showDeleteModal = ref(false)
const deletingCalendar = ref(null)
const showRestoreModal = ref(false)
const restoringCalendar = ref(null)

const form = ref({
  name: '',
  description: '',
  color: '#3788d9',
  is_public: false
})

const colorOptions = [
  '#3788d9',
  '#06d6a0',
  '#ff595e',
  '#ffca3a',
  '#8b5cf6',
  '#ec4899'
]

async function fetchCalendars() {
  loading.value = true
  try {
    // 使用 getAllCalendars 获取所有日历（包括已删除的）
    const { errCode, data } = await calendarApi.getAllCalendars()
    if (errCode === 0) {
      calendars.value = data.list || []
      // 预加载创建者信息
      calendars.value.forEach(cal => {
        if (cal.UserID) {
          usersStore.fetchUser(cal.UserID)
        }
      })
      // 获取每个日历的事件数量
      fetchEventCounts()
    }
  } catch {
    // 拦截器已处理
  } finally {
    loading.value = false
  }
}

async function fetchEventCounts() {
  // 获取一年前到现在的事件，用于统计
  const now = new Date()
  const oneYearAgo = new Date()
  oneYearAgo.setFullYear(now.getFullYear() - 1)

  const startStr = oneYearAgo.toISOString().split('T')[0]
  const endStr = now.toISOString().split('T')[0]

  try {
    const { errCode, data } = await calendarApi.getEvents({
      start_date: startStr,
      end_date: endStr
    })
    if (errCode === 0 && data.list) {
      // 按日历ID统计事件数量
      const counts = {}
      data.list.forEach(event => {
        if (event.CalendarID) {
          counts[event.CalendarID] = (counts[event.CalendarID] || 0) + 1
        }
      })
      eventCounts.value = counts
    }
  } catch {
    // 忽略错误
  }
}

function getCreatorName(userID) {
  return usersStore.getUsernameFromUserID(userID) || '...'
}

function getCreatorAvatar(userID) {
  return usersStore.getAvatarUrlFromUserID(userID)
}

function formatDateTime(dateStr) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString()
}

// 判断日历是否已删除（deleted_at 不为 NULL）
function isDeleted(calendar) {
  return calendar.DeletedAt !== null && calendar.DeletedAt !== undefined
}

function openEditModal(calendar) {
  editingCalendar.value = calendar
  form.value = {
    name: calendar.Name,
    description: calendar.Description || '',
    color: calendar.Color,
    is_public: calendar.IsPublic
  }
  showEditModal.value = true
}

async function updateCalendar() {
  if (!form.value.name.trim()) {
    toast.error(t('calendar.name_required'))
    return
  }

  try {
    const { errCode } = await calendarApi.updateCalendar({
      id: editingCalendar.value.ID,
      ...form.value
    })
    if (errCode === 0) {
      toast.success(t('calendar.update_success'))
      showEditModal.value = false
      fetchCalendars()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  }
}

function openDeleteModal(calendar) {
  deletingCalendar.value = calendar
  showDeleteModal.value = true
}

async function confirmDelete() {
  if (!deletingCalendar.value) return
  try {
    const { errCode } = await calendarApi.deleteCalendar(deletingCalendar.value.ID)
    if (errCode === 0) {
      toast.success(t('calendar.delete_success'))
      fetchCalendars()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    showDeleteModal.value = false
    deletingCalendar.value = null
  }
}

function openRestoreModal(calendar) {
  restoringCalendar.value = calendar
  showRestoreModal.value = true
}

async function confirmRestore() {
  if (!restoringCalendar.value) return
  try {
    const { errCode } = await calendarApi.restoreCalendar(restoringCalendar.value.ID)
    if (errCode === 0) {
      toast.success(t('calendar.restore_success'))
      fetchCalendars()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    showRestoreModal.value = false
    restoringCalendar.value = null
  }
}

onMounted(fetchCalendars)
</script>

<template>
  <div class="min-h-screen bg-gray-50 p-6 dark:bg-dk-base">
    <div class="mx-auto max-w-6xl">
      <!-- Header -->
      <div class="mb-6 flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-dk-text">{{ t('calendar.admin_title') }}</h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-dk-subtle">
            {{ t('calendar.calendars') }}
          </p>
        </div>
        <button
          @click="fetchCalendars"
          class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
        >
          <IconClock :size="16" />
          {{ t('sysadmin.refresh') }}
        </button>
      </div>

      <!-- Calendar List Table -->
      <div class="rounded-lg border border-gray-200 bg-white dark:border-dk-muted dark:bg-dk-card">
        <!-- Loading -->
        <div v-if="loading" class="flex items-center justify-center py-12">
          <svg class="h-6 w-6 animate-spin text-blue-500" viewBox="0 0 24 24" fill="none">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          <span class="ml-2 text-gray-500 dark:text-dk-subtle">{{ t('calendar.loading') }}</span>
        </div>

        <!-- Empty -->
        <div v-else-if="calendars.length === 0" class="py-12 text-center">
          <IconCalendar :size="48" class="mx-auto mb-3 text-gray-300 dark:text-dk-muted" />
          <p class="text-gray-500 dark:text-dk-subtle">{{ t('calendar.no_calendars') }}</p>
        </div>

        <!-- Table -->
        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 dark:divide-dk-muted">
            <thead class="bg-gray-50 dark:bg-dk-base">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">
                  {{ t('calendar.name') }}
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">
                  {{ t('schedule.event_type') }}
                </th>
                <th class="px-6 py-3 text-center text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">
                  {{ t('calendar.event_count') }}
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">
                  {{ t('customer.created_by') }}
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">
                  {{ t('customer.created_at') }}
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-dk-subtle">
                  {{ t('common.actions') }}
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-dk-muted">
              <tr
                v-for="calendar in calendars"
                :key="calendar.ID"
                class="hover:bg-gray-50 dark:hover:bg-dk-base"
                :class="{ 'opacity-50': isDeleted(calendar) }"
              >
                <!-- Name + Color -->
                <td class="whitespace-nowrap px-6 py-4">
                  <div class="flex items-center gap-3">
                    <div
                      class="h-4 w-4 rounded-full flex-shrink-0"
                      :style="{ backgroundColor: calendar.Color || '#3788d9' }"
                    ></div>
                    <div>
                      <div class="flex items-center gap-2 font-medium text-gray-900 dark:text-dk-text">
                        {{ calendar.Name }}
                        <span
                          v-if="isDeleted(calendar)"
                          class="inline-flex rounded-full bg-red-100 px-2 py-0.5 text-xs font-medium text-red-800 dark:bg-red-900/30 dark:text-red-400"
                        >
                          {{ t('calendar.deleted') }}
                        </span>
                      </div>
                      <div v-if="calendar.Description" class="text-sm text-gray-500 dark:text-dk-subtle">
                        {{ calendar.Description }}
                      </div>
                    </div>
                  </div>
                </td>

                <!-- Public/Private -->
                <td class="whitespace-nowrap px-6 py-4">
                  <span
                    class="inline-flex rounded-full px-2 py-1 text-xs font-medium"
                    :class="calendar.IsPublic
                      ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
                      : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'"
                  >
                    {{ calendar.IsPublic ? t('calendar.public') : t('calendar.private') }}
                  </span>
                </td>

                <!-- Event Count -->
                <td class="whitespace-nowrap px-6 py-4 text-center">
                  <div class="inline-flex items-center gap-1 text-gray-900 dark:text-dk-text">
                    <IconCalendar :size="16" class="text-gray-400" />
                    <span class="font-medium">{{ eventCounts[calendar.ID] || 0 }}</span>
                  </div>
                </td>

                <!-- Creator -->
                <td class="whitespace-nowrap px-6 py-4">
                  <div class="flex items-center gap-2">
                    <img
                      :src="getCreatorAvatar(calendar.UserID)"
                      class="h-6 w-6 rounded-full"
                      alt="avatar"
                    />
                    <span class="text-sm text-gray-900 dark:text-dk-text">
                      {{ getCreatorName(calendar.UserID) }}
                    </span>
                  </div>
                </td>

                <!-- Created At -->
                <td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500 dark:text-dk-subtle">
                  {{ formatDateTime(calendar.CreatedAt) }}
                </td>

                <!-- Actions -->
                <td class="whitespace-nowrap px-6 py-4">
                  <div class="flex items-center gap-1">
                    <button
                      v-if="isDeleted(calendar)"
                      @click="openRestoreModal(calendar)"
                      class="rounded p-1.5 text-gray-400 transition-colors hover:bg-green-50 hover:text-green-600 dark:hover:bg-dk-muted"
                      :title="t('calendar.restore')"
                    >
                      <IconRestore :size="16" />
                    </button>
                    <button
                      v-else
                      @click="openEditModal(calendar)"
                      class="rounded p-1.5 text-gray-400 transition-colors hover:bg-blue-50 hover:text-blue-600 dark:hover:bg-dk-muted"
                      :title="t('common.edit')"
                    >
                      <IconEdit :size="16" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Summary -->
      <div v-if="calendars.length > 0" class="mt-4 text-sm text-gray-500 dark:text-dk-subtle">
        {{ t('calendar.calendars') }}: {{ calendars.length }}
      </div>
    </div>
  </div>

  <!-- Edit Calendar Modal -->
  <div
    v-if="showEditModal"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
    @click.self="showEditModal = false"
  >
    <div class="w-full max-w-md rounded-xl bg-white p-6 shadow-xl dark:bg-dk-card">
      <h3 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">
        {{ t('calendar.edit_calendar') }}
      </h3>

      <div class="mb-4">
        <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
          {{ t('calendar.name') }} *
        </label>
        <input
          v-model="form.name"
          type="text"
          class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white"
          :placeholder="t('calendar.name_placeholder')"
        />
      </div>

      <div class="mb-4">
        <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
          {{ t('calendar.description') }}
        </label>
        <textarea
          v-model="form.description"
          rows="3"
          class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white"
          :placeholder="t('calendar.description_placeholder')"
        ></textarea>
      </div>

      <div class="mb-4">
        <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
          {{ t('calendar.color') }}
        </label>
        <div class="flex gap-2">
          <button
            v-for="color in colorOptions"
            :key="color"
            @click="form.color = color"
            class="h-8 w-8 rounded-full transition-transform hover:scale-110"
            :class="{ 'ring-2 ring-blue-500 ring-offset-2': form.color === color }"
            :style="{ backgroundColor: color }"
          ></button>
        </div>
      </div>

      <div class="mb-6 flex items-center gap-2">
        <input
          id="edit_is_public"
          v-model="form.is_public"
          type="checkbox"
          class="rounded border-gray-300"
        />
        <label for="edit_is_public" class="text-sm text-gray-700 dark:text-gray-300">
          {{ t('calendar.is_public') }}
        </label>
      </div>

      <div class="flex justify-end gap-2">
        <button
          @click="showEditModal = false"
          class="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:text-gray-300 dark:hover:bg-dk-muted"
        >
          {{ t('cancel') }}
        </button>
        <button
          @click="updateCalendar"
          class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
        >
          {{ t('save') }}
        </button>
      </div>
    </div>
  </div>

  <!-- Delete Confirm Dialog -->
  <ConfirmDialog
    v-model="showDeleteModal"
    :title="t('calendar.confirm_delete')"
    :message="t('calendar.confirm_delete_message', { name: deletingCalendar?.Name || '' })"
    :confirm-text="t('delete')"
    :cancel-text="t('cancel')"
    danger
    @confirm="confirmDelete"
  />

  <!-- Restore Confirm Dialog -->
  <ConfirmDialog
    v-model="showRestoreModal"
    :title="t('calendar.confirm_restore')"
    :message="t('calendar.confirm_restore_message', { name: restoringCalendar?.Name || '' })"
    :confirm-text="t('calendar.restore')"
    :cancel-text="t('cancel')"
    @confirm="confirmRestore"
  />
</template>
