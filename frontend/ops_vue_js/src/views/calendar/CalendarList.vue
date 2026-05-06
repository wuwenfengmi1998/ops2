<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { calendarApi } from '@/api/calendar'
import { IconPlus, IconCalendar, IconTrash, IconEdit } from '@tabler/icons-vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

usePageTitle('appname.calendar')
const { t } = useI18n()
const router = useRouter()
const toast = useToastStore()

const calendars = ref([])
const loading = ref(false)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const editingCalendar = ref(null)
const deletingCalendar = ref(null)

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
    const { errCode, data } = await calendarApi.getCalendars()
    if (errCode === 0) {
      calendars.value = data.list || []
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  form.value = { name: '', description: '', color: '#3788d9', is_public: false }
  showCreateModal.value = true
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

async function createCalendar() {
  if (!form.value.name.trim()) {
    toast.error(t('calendar.name_required'))
    return
  }

  try {
    const { errCode, data } = await calendarApi.createCalendar(form.value)
    if (errCode === 0) {
      toast.success(t('calendar.create_success'))
      showCreateModal.value = false
      fetchCalendars()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  }
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

async function deleteCalendar(calendar) {
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

function goToCalendar(id) {
  router.push(`/calendar/${id}`)
}

onMounted(fetchCalendars)
</script>

<template>
  <div class="mx-auto max-w-6xl px-6 py-6">
    <div class="flex flex-col gap-6 rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <!-- Header -->
      <div class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-dk-muted">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('calendar.calendars') }}</h3>
        <button
          @click="openCreateModal"
          class="inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700"
        >
          <IconPlus :size="16" />
          {{ t('calendar.create_calendar') }}
        </button>
      </div>

      <!-- Calendar List -->
      <div class="px-6 py-3">
        <div v-if="loading" class="py-8 text-center text-gray-400">
          <svg class="mx-auto mb-2 h-5 w-5 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          Loading...
        </div>

        <div v-else-if="calendars.length === 0" class="py-8 text-center text-gray-400">
          {{ t('calendar.no_calendars') }}
        </div>

        <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <div
            v-for="calendar in calendars"
            :key="calendar.ID"
            class="cursor-pointer rounded-lg border border-gray-200 p-4 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:hover:bg-dk-base"
            @click="goToCalendar(calendar.ID)"
          >
            <div class="flex items-start justify-between">
              <div class="flex items-center gap-3">
                <div
                  class="h-4 w-4 rounded-full"
                  :style="{ backgroundColor: calendar.Color }"
                ></div>
                <div>
                  <h4 class="font-medium text-gray-900 dark:text-white">{{ calendar.Name }}</h4>
                  <p v-if="calendar.Description" class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                    {{ calendar.Description }}
                  </p>
                </div>
              </div>

              <div class="flex items-center gap-1" @click.stop>
                <button
                  v-if="calendar.canEdit"
                  @click="openEditModal(calendar)"
                  class="rounded p-1 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-dk-muted"
                >
                  <IconEdit :size="16" />
                </button>
                <button
                  v-if="calendar.canEdit"
                  @click="deleteCalendar(calendar)"
                  class="rounded p-1 text-gray-400 transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-dk-muted"
                >
                  <IconTrash :size="16" />
                </button>
              </div>
            </div>

            <div class="mt-3 flex items-center gap-2 text-xs text-gray-400">
              <IconCalendar :size="14" />
              <span>{{ calendar.IsPublic ? t('calendar.public') : t('calendar.private') }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Create Calendar Modal -->
  <div
    v-if="showCreateModal"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
    @click.self="showCreateModal = false"
  >
    <div class="w-full max-w-md rounded-xl bg-white p-6 shadow-xl dark:bg-dk-card">
      <h3 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">
        {{ t('calendar.create_calendar') }}
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
          id="is_public"
          v-model="form.is_public"
          type="checkbox"
          class="rounded border-gray-300"
        />
        <label for="is_public" class="text-sm text-gray-700 dark:text-gray-300">
          {{ t('calendar.is_public') }}
        </label>
      </div>

      <div class="flex justify-end gap-2">
        <button
          @click="showCreateModal = false"
          class="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:text-gray-300 dark:hover:bg-dk-muted"
        >
          {{ t('cancel') }}
        </button>
        <button
          @click="createCalendar"
          class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
        >
          {{ t('create') }}
        </button>
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
</template>
