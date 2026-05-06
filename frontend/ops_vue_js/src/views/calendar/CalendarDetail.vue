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
import DatatimePickerForFullCalendar from "@/components/datatimePickerForFullCalendar.vue"

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
  isPublic: false,
  isEditing: false,
  isEditable: false,
})

const colorOptions = ref([
  { value: "#3788d9", label: t("schedule.work"), name: t("schedule.work"), type: "work" },
  { value: "#06d6a0", label: t("schedule.duty"), name: t("schedule.duty"), type: "duty" },
  { value: "#ff595e", label: t("schedule.exam"), name: t("schedule.exam"), type: "exam" },
  { value: "#ffca3a", label: t("schedule.standby"), name: t("schedule.standby"), type: "standby" },
  { value: "#D16C13", label: t("schedule.personal_holiday"), name: t("schedule.personal_holiday"), type: "personal_holiday" },
  { value: "#D10D21", label: t("schedule.public_holiday"), name: t("schedule.public_holiday"), type: "public_holiday" },
])

const pageData = ref({
  seleEventID: 0,
  lastClickTime: 0,
  lastClickTimeStr: "",
  lastEventClickTime: 0,
  lastEventClickID: 0,
  submitChecked: false,
  lastEventsSnapshot: null,
})

// 选中/取消选中事件
function unseleEvent(eventID) {
  const target = calendarOptions.value.events.find(item => item.id === eventID)
  if (target) {
    target.borderColor = "#F7F7F7"
  }
}

function unseleEventAll() {
  unseleEvent(pageData.value.seleEventID)
  pageData.value.seleEventID = 0
}

function closeEventModal() {
  showModal.value = false
}

function openEventModal(dateStr, dataEnd, id = 0, title = "", color = "#3788d9", isPublic = false, isEditing = false, isEditable = true) {
  eventData.value = {
    id: id,
    title: title,
    startDate: dateStr,
    endDate: dataEnd,
    color: color,
    isPublic: isPublic,
    isEditing: isEditing,
    isEditable: isEditable,
  }
  showModal.value = true
}

function editEvent(info) {
  openEventModal(
    info.event.startStr,
    info.event.end ? info.event.endStr : info.event.startStr,
    parseInt(info.event.id),
    info.event.title,
    info.event.backgroundColor,
    info.event.extendedProps?.isPublic || false,
    true,
    info.event.durationEditable,
  )
}

function selectColor(colorValue) {
  if (eventData.value.isEditable) {
    eventData.value.color = colorValue
  }
}

// 日期转后端格式：YYYY-MM-DD 00:00:00
function toDatetime(dateStr) {
  return dateStr ? dateStr + " 00:00:00" : ""
}

async function saveEvent() {
  if (!eventData.value.title.trim()) {
    pageData.value.submitChecked = true
    toast.warning(t('calendar.event_title_required'))
    return
  }
  pageData.value.submitChecked = false

  if (!eventData.value.startDate || !eventData.value.endDate) {
    toast.warning(t('schedule.date_required'))
    return
  }

  const selectedColor = colorOptions.value.find(c => c.value === eventData.value.color)
  const scheduleType = selectedColor ? selectedColor.type : "work"

  try {
    let result
    if (eventData.value.isEditing) {
      result = await calendarApi.updateEvent({
        id: eventData.value.id,
        title: eventData.value.title.trim(),
        start: toDatetime(eventData.value.startDate),
        end: toDatetime(
          eventData.value.startDate === eventData.value.endDate
            ? eventData.value.endDate
            : DateUtils.toRealEnd(eventData.value.endDate),
        ),
        schedule_type: scheduleType,
        is_public: eventData.value.isPublic,
      })
    } else {
      result = await calendarApi.addEvent({
        calendar_id: calendarId.value,
        title: eventData.value.title.trim(),
        start: toDatetime(eventData.value.startDate),
        end: toDatetime(
          eventData.value.startDate === eventData.value.endDate
            ? eventData.value.endDate
            : DateUtils.toRealEnd(eventData.value.endDate),
        ),
        schedule_type: scheduleType,
        is_public: eventData.value.isPublic,
      })
    }

    if (result.errCode === 0) {
      toast.success(t('calendar.event_save_success'))
      closeEventModal()
      getEvents()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  }
}

async function deleteEvent() {
  if (!confirm(t('calendar.confirm_delete_event'))) return

  try {
    const result = await calendarApi.deleteEvent(eventData.value.id)
    if (result.errCode === 0) {
      toast.success(t('calendar.event_delete_success'))
      closeEventModal()
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

  const calApi = calendarRef.value.getApi()
  const start = DateUtils.dateToStr(calendarNowShow.value.start)
  const end = DateUtils.toRealEnd(calendarNowShow.value.end)

  try {
    const { errCode, data } = await calendarApi.getEvents({
      calendar_id: calendarId.value,
      start: start,
      end: end,
    })

    if (errCode === 0) {
      calendarOptions.value.events = []
      ;(data.list || []).forEach(item => {
        calendarOptions.value.events.push({
          id: item.ID,
          title: item.Title,
          start: item.StartDate,
          end: item.StartDate === item.EndDate
            ? item.EndDate
            : DateUtils.toCalendarEnd(item.EndDate),
          backgroundColor: item.BgColor,
          borderColor: item.ID === pageData.value.seleEventID ? "#000000" : "#F7F7F7",
          allDay: true,
        })
      })
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

// ─── 滚动标题工具函数 ───────────────────────────────────────────────────────
function applyScrollToTitle(titleEl) {
  titleEl.removeAttribute("data-truncated")
  titleEl.style.removeProperty("--scroll-distance")
  const overflow = titleEl.scrollWidth - titleEl.clientWidth
  if (overflow > 0) {
    titleEl.style.setProperty("--scroll-distance", `-${overflow}px`)
    titleEl.setAttribute("data-truncated", "true")
  }
}

function recalcScrollTitles() {
  nextTick(() => {
    requestAnimationFrame(() => {
      const calendarEl = calendarRef.value?.$el
      if (!calendarEl) return
      calendarEl.querySelectorAll(".fc-event-title").forEach(applyScrollToTitle)
    })
  })
}

// ─────────────────────────────────────────────────────────────────────────────

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

  dayCellDidMount(info) {
    if (info.date.getDay() === 0 || info.date.getDay() === 6) {
      info.el.style.backgroundColor = "#f5f5f5"
    }
    info.el.style.border = "1px solid #e5e7eb"
  },

  headerToolbar: {
    left: "prevYear,prev,today,next,nextYear",
    center: "title",
    right: "",
  },

  customButtons: {
    prevYear: {
      text: t("schedule.previous_year"),
      click() { calendarRef.value.getApi().prevYear(); getEvents() },
    },
    nextYear: {
      text: t("schedule.next_year"),
      click() { calendarRef.value.getApi().nextYear(); getEvents() },
    },
    prev: {
      text: t("schedule.previous_month"),
      click() { calendarRef.value.getApi().prev(); getEvents() },
    },
    next: {
      text: t("schedule.next_month"),
      click() { calendarRef.value.getApi().next(); getEvents() },
    },
    today: {
      text: t("schedule.today"),
      click() { calendarRef.value.getApi().today(); getEvents() },
    },
    week: {
      text: t("schedule.week"),
      click() { calendarRef.value.getApi().changeView("timeGridWeek") },
    },
  },

  events: [],

  eventDidMount(info) {
    const titleEl = info.el.querySelector(".fc-event-title")
    if (titleEl) {
      requestAnimationFrame(() => applyScrollToTitle(titleEl))
    }
  },

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
    const eventid = parseInt(info.event.id)

    unseleEventAll()
    const target = calendarOptions.value.events.find(item => String(item.id) === String(info.event.id))
    if (target) {
      target.borderColor = "#000000"
      pageData.value.seleEventID = target.id
    }

    if (eventid === pageData.value.lastEventClickID) {
      if (timeDifference < 400 && timeDifference > 0) {
        editEvent(info)
        unseleEventAll()
      }
    }
    pageData.value.lastEventClickID = eventid
    pageData.value.lastEventClickTime = nowTime
  },

  eventDrop(info) {
    // 拖拽后直接更新
    const selectedColor = colorOptions.value.find(c => c.value === info.event.backgroundColor)
    const scheduleType = selectedColor ? selectedColor.type : "work"
    const startStr = info.event.startStr
    const endStr = info.event.end ? info.event.endStr : startStr
    calendarApi.updateEvent({
      id: parseInt(info.event.id),
      title: info.event.title,
      start: toDatetime(startStr),
      end: toDatetime(startStr === endStr ? endStr : DateUtils.toRealEnd(endStr)),
      schedule_type: scheduleType,
    }).then(r => {
      if (r.errCode !== 0) toast.error(t('message.server_error'))
      else getEvents()
    })
  },
})

// 监听语言变化
watch(locale, () => {
  calendarOptions.value.locale = locale.value
  calendarOptions.value.customButtons.prevYear.text = t("schedule.previous_year")
  calendarOptions.value.customButtons.nextYear.text = t("schedule.next_year")
  calendarOptions.value.customButtons.prev.text = t("schedule.previous_month")
  calendarOptions.value.customButtons.next.text = t("schedule.next_month")
  calendarOptions.value.customButtons.today.text = t("schedule.today")
  calendarOptions.value.customButtons.week.text = t("schedule.week")
  colorOptions.value = [
    { value: "#3788d9", label: t("schedule.work"), name: t("schedule.work"), type: "work" },
    { value: "#06d6a0", label: t("schedule.duty"), name: t("schedule.duty"), type: "duty" },
    { value: "#ff595e", label: t("schedule.exam"), name: t("schedule.exam"), type: "exam" },
    { value: "#ffca3a", label: t("schedule.standby"), name: t("schedule.standby"), type: "standby" },
    { value: "#D16C13", label: t("schedule.personal_holiday"), name: t("schedule.personal_holiday"), type: "personal_holiday" },
    { value: "#D10D21", label: t("schedule.public_holiday"), name: t("schedule.public_holiday"), type: "public_holiday" },
  ]
})

let resizeObserver = null

onMounted(() => {
  fetchCalendarInfo()
  // 监听日历容器宽度变化
  let resizeTimer = null
  resizeObserver = new ResizeObserver(() => {
    clearTimeout(resizeTimer)
    resizeTimer = setTimeout(() => recalcScrollTitles(), 150)
  })
  if (calendarRef.value?.$el) {
    resizeObserver.observe(calendarRef.value.$el)
  }
  onBeforeUnmount(() => {
    if (resizeObserver) {
      resizeObserver.disconnect()
      resizeObserver = null
    }
    clearTimeout(resizeTimer)
  })
})
</script>

<template>
  <div class="flex w-full flex-col relative">
    <!-- 事件编辑模态框 -->
    <div
      v-if="showModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-gray-800/20"
    >
      <div
        class="modal-content bg-white rounded-lg shadow-lg w-full max-w-2xl max-h-[95vh] flex flex-col"
      >
        <!-- 模态框头部 -->
        <div
          class="modal-header border-b p-4 flex justify-between items-center flex-shrink-0"
        >
          <h5 class="modal-title text-lg font-semibold">
            {{
              userStore.isLoggedIn
                ? eventData.isEditing
                  ? t("calendar.edit_event")
                  : t("calendar.add_event")
                : t("calendar.view_event")
            }}
          </h5>
          <button
            @click="closeEventModal"
            class="btn-close text-gray-500 hover:text-gray-700"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="icon icon-tabler icons-tabler-outline icon-tabler-x"
            >
              <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
              <path d="M18 6l-12 12"></path>
              <path d="M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <!-- 主体区域 -->
        <div class="modal-body p-4 flex-1 overflow-y-auto">
          <!-- 日期选择区域 -->
          <DatatimePickerForFullCalendar
            v-model:startDate="eventData.startDate"
            v-model:endDate="eventData.endDate"
            :color="eventData.color"
            :title="eventData.title"
            :isEditable="eventData.isEditable"
          />

          <!-- 内容输入区域 -->
          <div class="mb-4">
            <div class="uni-easyinput input relative">
              <div
                class="uni-easyinput__content is-input-border border border-gray-300 rounded-md bg-white relative"
                :class="{
                  'border-gray-300': eventData.title || !pageData.submitChecked,
                  'border-red-500': !eventData.title && pageData.submitChecked,
                }"
              >
                <input
                  v-model="eventData.title"
                  type="text"
                  maxlength="140"
                  class="uni-easyinput__content-input w-full px-3 py-2 outline-none"
                  :placeholder="t('calendar.event_title_placeholder')"
                  @keyup.enter="saveEvent"
                  :disabled="!eventData.isEditable"
                />
              </div>
            </div>
          </div>

          <!-- 颜色选择区域 -->
          <div class="mb-4">
            <div class="color_box grid grid-cols-3 gap-2">
              <div
                v-for="color in colorOptions"
                :key="color.value"
                class="color_box_item"
              >
                <label
                  class="uni-label-pointer form-colorinput flex items-center gap-2 cursor-pointer"
                  @click="selectColor(color.value)"
                >
                  <div class="uni-radio-wrapper">
                    <div
                      class="uni-radio-input flex items-center justify-center w-6 h-6 rounded-full transition-all"
                      :style="{
                        backgroundColor: color.value,
                        borderColor: color.value,
                      }"
                    >
                      <svg
                        v-if="eventData.color === color.value"
                        width="18"
                        height="18"
                        viewBox="0 0 32 32"
                      >
                        <path
                          d="M1.952 18.080q-0.32-0.352-0.416-0.88t0.128-0.976l0.16-0.352q0.224-0.416 0.64-0.528t0.8 0.176l6.496 4.704q0.384 0.288 0.912 0.272t0.88-0.336l17.312-14.272q0.352-0.288 0.848-0.256t0.848 0.352l-0.416-0.416q0.32 0.352 0.32 0.816t-0.32 0.816l-18.656 18.912q-0.32 0.352-0.8 0.352t-0.8-0.32l-7.936-8.064z"
                          fill="#ffffff"
                        ></path>
                      </svg>
                    </div>
                  </div>
                  <span class="text-gray-700">{{ color.label }}</span>
                </label>
              </div>
            </div>
          </div>

          <!-- 公共日程开关 -->
          <div class="mb-4 flex items-center justify-between">
            <span class="text-gray-700">{{ t('calendar.is_public_event') }}</span>
            <label class="relative inline-flex items-center cursor-pointer">
              <input
                v-model="eventData.isPublic"
                type="checkbox"
                class="sr-only peer"
                :disabled="!eventData.isEditable"
              />
              <div class="w-11 h-6 bg-gray-300 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-cyan-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-cyan-600 peer-disabled:opacity-50 peer-disabled:cursor-not-allowed"></div>
            </label>
          </div>
        </div>

        <!-- 底部固定 -->
        <div
          v-if="userStore.isLoggedIn"
          class="modal-footer border-t p-4 flex justify-end items-center flex-shrink-0"
        >
          <div class="flex gap-2">
            <button
              v-if="eventData.isEditing"
              @click="deleteEvent"
              class="btn px-4 py-2 text-white bg-red-500 hover:bg-red-600 rounded-md disabled:bg-gray-400 disabled:cursor-not-allowed"
              :disabled="!eventData.isEditable"
            >
              {{ t('delete') }}
            </button>
            <button
              v-if="!eventData.isEditing"
              @click="saveEvent"
              class="btn btn-primary px-4 py-2 bg-cyan-600 text-white hover:bg-cyan-700 rounded-md disabled:bg-gray-400 disabled:cursor-not-allowed"
              :disabled="!eventData.isEditable"
            >
              {{ t('calendar.add_event') }}
            </button>
            <button
              v-if="eventData.isEditing"
              @click="saveEvent"
              class="btn btn-primary px-4 py-2 bg-teal-600 text-white hover:bg-teal-700 rounded-md disabled:bg-gray-400 disabled:cursor-not-allowed"
              :disabled="!eventData.isEditable"
            >
              {{ t('calendar.edit_event') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 日历主体区域 -->
    <div
      class="flex-1 rounded-lg border border-gray-200 bg-white p-0.5 shadow dark:border-dk-muted dark:bg-dk-card"
    >
      <div class="h-full w-full overflow-hidden rounded-md">
        <FullCalendar ref="calendarRef" :options="calendarOptions" />
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 父容器作为裁剪视口 */
:deep(.fc-daygrid-event .fc-event-title-container) {
  overflow: hidden !important;
}

/* 默认状态：单行截断省略 */
:deep(.fc-daygrid-event .fc-event-title) {
  white-space: nowrap !important;
  overflow: visible !important;
  text-overflow: ellipsis !important;
  display: block !important;
  will-change: transform;
}

/* 需要滚动的标题 */
:deep(.fc-daygrid-event .fc-event-title[data-truncated="true"]) {
  text-overflow: clip !important;
  display: inline-block !important;
  animation: marquee-bounce 6s ease-in-out infinite !important;
}

@keyframes marquee-bounce {
  0%, 20% { transform: translateX(0); }
  60% { transform: translateX(var(--scroll-distance, 0px)); }
  80% { transform: translateX(var(--scroll-distance, 0px)); }
  100% { transform: translateX(0); }
}
</style>
