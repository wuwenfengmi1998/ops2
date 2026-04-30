<script setup>
import { ref, watch, onMounted, onBeforeUnmount, nextTick } from "vue"
import { useRoute, useRouter } from "vue-router"
import FullCalendar from "@fullcalendar/vue3"
import dayGridPlugin from "@fullcalendar/daygrid"
import timeGridPlugin from "@fullcalendar/timegrid"
import interactionPlugin from "@fullcalendar/interaction"
import listPlugin from "@fullcalendar/list"
import { useI18n } from "vue-i18n"
import { usePageTitle } from "@/composables/usePageTitle"
import { useToastStore } from "@/stores/toast"
import { useUserStore } from "@/stores/user"
import { calendarApi } from "@/api/calendar"
import { useDateUtils } from "@/composables/useDateUtils"
import { IconPlus, IconTrash, IconEdit, IconCalendar } from "@tabler/icons-vue"

const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()
const toast = useToastStore()
const userStore = useUserStore()
const DateUtils = useDateUtils()

const calendarId = ref(parseInt(route.params.id))
const calendarInfo = ref({})
const loading = ref(false)

usePageTitle('appname.calendar_detail')

const calendarRef = ref(null)
const calendarNowShow = ref()

const showModal = ref(false)
const eventData = ref({
  id: 0,
  title: "",
  startDate: "",
  endDate: "",
  color: "#3788d9",
  remark: "",
  isEditing: false,
  isEditable: false
})

const colorOptions = ref([
  { value: "#3788d9", label: t("schedule.work") },
  { value: "#06d6a0", label: t("schedule.duty") },
  { value: "#ff595e", label: t("schedule.exam") },
  { value: "#ffca3a", label: t("schedule.standby") },
])

const pageData = ref({
  seleEventID: 0,
  lastClickTime: 0,
  lastClickTimeStr: "",
  lastEventClickTime: 0,
  lastEventClickID: 0,
  submitChecked: false
})

function openEventModal(startDate, endDate) {
  eventData.value = {
    id: 0,
    title: "",
    startDate: startDate || "",
    endDate: endDate || "",
    color: calendarInfo.value.Color || "#3788d9",
    remark: "",
    isEditing: false,
    isEditable: true
  }
  showModal.value = true
}

function editEvent(event) {
  eventData.value = {
    id: event.id,
    title: event.title,
    startDate: event.start,
    endDate: event.end || event.start,
    color: event.backgroundColor,
    remark: event.extendedProps?.remark || "",
    isEditing: true,
    isEditable: true
  }
  showModal.value = true
}

async function saveEvent() {
  if (!eventData.value.title.trim()) {
    toast.error(t('calendar.event_title_required'))
    return
  }

  pageData.value.submitChecked = true

  try {
    let result
    if (eventData.value.isEditing) {
      result = await calendarApi.updateEvent({
        id: eventData.value.id,
        title: eventData.value.title,
        start: eventData.value.startDate,
        end: eventData.value.endDate,
        color: eventData.value.color,
        remark: eventData.value.remark
      })
    } else {
      result = await calendarApi.addEvent({
        calendar_id: calendarId.value,
        title: eventData.value.title,
        start: eventData.value.startDate,
        end: eventData.value.endDate,
        color: eventData.value.color,
        remark: eventData.value.remark
      })
    }

    if (result.errCode === 0) {
      toast.success(t('calendar.event_save_success'))
      showModal.value = false
      getEvents()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    pageData.value.submitChecked = false
  }
}

async function deleteEvent() {
  if (!confirm(t('calendar.confirm_delete_event'))) {
    return
  }

  try {
    const result = await calendarApi.deleteEvent(eventData.value.id)
    if (result.errCode === 0) {
      toast.success(t('calendar.event_delete_success'))
      showModal.value = false
      getEvents()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  }
}

async function getEvents() {
  if (!calendarRef.value) return

  const calendarApi2 = calendarRef.value.getApi()
  const view = calendarApi2.view
  const start = DateUtils.formatDate(view.activeStart)
  const end = DateUtils.formatDate(view.activeEnd)

  try {
    const { errCode, data } = await calendarApi.getEvents({
      calendar_id: calendarId.value,
      start: start,
      end: end
    })

    if (errCode === 0) {
      const events = (data.list || []).map(event => ({
        id: event.ID,
        title: event.Title,
        start: event.StartDate,
        end: event.EndDate,
        backgroundColor: event.BgColor,
        borderColor: event.BgColor,
        extendedProps: {
          remark: event.Remark
        }
      }))
      calendarOptions.value.events = events
    }
  } catch {
    // 拦截器已处理
  }
}

async function fetchCalendarInfo() {
  loading.value = true
  try {
    const { errCode, data } = await calendarApi.getCalendars()
    if (errCode === 0) {
      const calendar = (data.list || []).find(c => c.ID === calendarId.value)
      if (calendar) {
        calendarInfo.value = calendar
      } else {
        toast.error(t('calendar.calendar_not_found'))
        router.push('/calendars')
      }
    }
  } catch {
    // 拦截器已处理
  } finally {
    loading.value = false
  }
}

function unseleEvent(eventID) {
  const target = calendarOptions.value.events.find(item => item.id === eventID)
  if (target) {
    target.borderColor = target.backgroundColor
  }
}

function unseleEventAll() {
  unseleEvent(pageData.value.seleEventID)
  pageData.value.seleEventID = 0
}

const calendarOptions = ref({
  height: "100%",
  contentHeight: "auto",
  locale: locale.value,
  plugins: [dayGridPlugin, timeGridPlugin, interactionPlugin, listPlugin],
  nowIndicator: true,
  weekends: true,
  initialView: "dayGridMonth",
  selectable: true,
  editable: true,
  dayMaxEvents: true,
  firstDay: 1,
  expandRows: true,
  stickyHeaderDates: true,

  headerToolbar: {
    left: "prevYear,prev,today,next,nextYear",
    center: "title",
    right: ""
  },

  customButtons: {
    prevYear: {
      text: t("schedule.previous_year"),
      click() {
        calendarRef.value.getApi().prevYear()
        getEvents()
      }
    },
    nextYear: {
      text: t("schedule.next_year"),
      click() {
        calendarRef.value.getApi().nextYear()
        getEvents()
      }
    },
    prev: {
      text: t("schedule.previous_month"),
      click() {
        calendarRef.value.getApi().prev()
        getEvents()
      }
    },
    next: {
      text: t("schedule.next_month"),
      click() {
        calendarRef.value.getApi().next()
        getEvents()
      }
    },
    today: {
      text: t("schedule.today"),
      click() {
        calendarRef.value.getApi().today()
        getEvents()
      }
    }
  },

  events: [],

  datesSet(info) {
    calendarNowShow.value = info
    getEvents()
  },

  dateClick(info) {
    const nowTime = new Date().getTime()
    const timeDifference = nowTime - pageData.value.lastClickTime

    unseleEventAll()

    if (info.dateStr === pageData.value.lastClickTimeStr) {
      if (timeDifference < 400 && timeDifference > 0) {
        if (userStore.isLoggedIn) {
          openEventModal(info.dateStr, info.dateStr)
        } else {
          toast.warning(t("message.login_to_your_account"))
          router.replace("/login?redirect=/calendar/" + calendarId.value)
        }
      }
    }
    pageData.value.lastClickTimeStr = info.dateStr
    pageData.value.lastClickTime = nowTime
  },

  select(info) {
    if (info.end - info.start > 86400000) {
      if (userStore.isLoggedIn) {
        openEventModal(info.startStr, info.endStr)
      } else {
        toast.warning(t("message.login_to_your_account"))
      }
    }
  },

  eventClick(info) {
    const nowTime = new Date().getTime()
    const timeDifference = nowTime - pageData.value.lastEventClickTime

    if (info.event.id === pageData.value.lastEventClickID) {
      if (timeDifference < 400 && timeDifference > 0) {
        editEvent({
          id: info.event.id,
          title: info.event.title,
          start: info.event.startStr?.split('T')[0] || info.event.start?.toISOString().split('T')[0],
          end: info.event.endStr?.split('T')[0] || info.event.end?.toISOString().split('T')[0],
          backgroundColor: info.event.backgroundColor
        })
      }
    }
    pageData.value.lastEventClickID = info.event.id
    pageData.value.lastEventClickTime = nowTime

    // 选中效果
    unseleEventAll()
    const target = calendarOptions.value.events.find(item => String(item.id) === String(info.event.id))
    if (target) {
      target.borderColor = "#000000"
      pageData.value.seleEventID = target.id
    }
  }
})

onMounted(() => {
  fetchCalendarInfo()
})
</script>

<template>
  <div class="flex flex-col gap-6 px-6 py-6">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <button
          @click="router.push('/calendars')"
          class="rounded-lg p-2 text-gray-600 transition-colors hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dk-muted"
        >
          ←
        </button>
        <div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ calendarInfo.Name || t('calendar.loading') }}
          </h3>
          <p v-if="calendarInfo.Description" class="text-sm text-gray-500 dark:text-gray-400">
            {{ calendarInfo.Description }}
          </p>
        </div>
      </div>
      <div
        v-if="userStore.isLoggedIn"
        class="flex items-center gap-2"
      >
        <button
          @click="openEventModal()"
          class="inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700"
        >
          <IconPlus :size="16" />
          {{ t('calendar.add_event') }}
        </button>
      </div>
    </div>

    <div class="rounded-xl border border-gray-200 bg-white shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <FullCalendar
        ref="calendarRef"
        :options="calendarOptions"
      />
    </div>

    <!-- Event Modal -->
    <div
      v-if="showModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
      @click.self="showModal = false"
    >
      <div class="w-full max-w-md rounded-xl bg-white p-6 shadow-xl dark:bg-dk-card">
        <h3 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">
          {{ eventData.isEditing ? t('calendar.edit_event') : t('calendar.add_event') }}
        </h3>

        <div class="mb-4">
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('calendar.event_title') }} *
          </label>
          <input
            v-model="eventData.title"
            type="text"
            class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            :placeholder="t('calendar.event_title_placeholder')"
          />
        </div>

        <div class="mb-4 grid grid-cols-2 gap-3">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
              {{ t('calendar.start_date') }} *
            </label>
            <input
              v-model="eventData.startDate"
              type="date"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
              {{ t('calendar.end_date') }} *
            </label>
            <input
              v-model="eventData.endDate"
              type="date"
              class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            />
          </div>
        </div>

        <div class="mb-4">
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('calendar.color') }}
          </label>
          <div class="flex gap-2">
            <button
              v-for="color in colorOptions"
              :key="color.value"
              @click="eventData.color = color.value"
              class="h-8 w-8 rounded-full transition-transform hover:scale-110"
              :class="{ 'ring-2 ring-blue-500 ring-offset-2': eventData.color === color.value }"
              :style="{ backgroundColor: color.value }"
              :title="color.label"
            ></button>
          </div>
        </div>

        <div class="mb-4">
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('calendar.remark') }}
          </label>
          <textarea
            v-model="eventData.remark"
            rows="3"
            class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition-colors focus:border-blue-500 dark:border-dk-muted dark:bg-dk-base dark:text-white"
            :placeholder="t('calendar.remark_placeholder')"
          ></textarea>
        </div>

        <div class="flex justify-between">
          <button
            v-if="eventData.isEditing"
            @click="deleteEvent"
            class="inline-flex items-center gap-1 rounded-lg bg-red-600 px-3 py-2 text-sm font-medium text-white transition-colors hover:bg-red-700"
          >
            <IconTrash :size="16" />
            {{ t('delete') }}
          </button>
          <div v-else></div>

          <div class="flex gap-2">
            <button
              @click="showModal = false"
              class="rounded-lg border border-gray-300 px-4 py-2 text-sm text-gray-600 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:text-gray-300 dark:hover:bg-dk-muted"
            >
              {{ t('cancel') }}
            </button>
            <button
              @click="saveEvent"
              :disabled="pageData.submitChecked"
              class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
            >
              {{ t('save') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
