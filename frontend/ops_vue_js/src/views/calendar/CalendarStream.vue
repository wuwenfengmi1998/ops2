<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount, nextTick } from "vue"
import { useRoute, useRouter } from "vue-router"
import { useI18n } from "vue-i18n"
import { usePageTitle } from "@/composables/usePageTitle"
import { useToastStore } from "@/stores/toast"
import { useUserStore } from "@/stores/user"
import { useUsersStore } from "@/stores/users"
import { calendarApi } from "@/api/calendar"
import { useDateUtils } from "@/composables/useDateUtils"
import DatatimePickerForFullCalendar from "@/components/datatimePickerForFullCalendar.vue"
import ConfirmDialog from "@/components/ConfirmDialog.vue"

const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()
const toast = useToastStore()
const userStore = useUserStore()
const usersStore = useUsersStore()
const DateUtils = useDateUtils()

usePageTitle('appname.calendar_stream')

const calendarId = ref(parseInt(route.params.id))
const calendarInfo = ref({})
const loading = ref(false)

const MAX_WEEKS = 60
const INITIAL_WEEKS_BEFORE = 10
const INITIAL_WEEKS_AFTER = 10
const LOAD_BATCH = 5

const scheduleColors = {
  work: '#066fd1',
  duty: '#09d119',
  exam: '#ff00ff',
  standby: '#ffca3a',
  personal_holiday: '#d16c13',
  public_holiday: '#d10d21',
}

const colorOptions = ref([
  { value: "#066fd1", label: t("schedule.work"), type: "work" },
  { value: "#09d119", label: t("schedule.duty"), type: "duty" },
  { value: "#ff00ff", label: t("schedule.exam"), type: "exam" },
  { value: "#ffca3a", label: t("schedule.standby"), type: "standby" },
  { value: "#d16c13", label: t("schedule.personal_holiday"), type: "personal_holiday" },
  { value: "#d10d21", label: t("schedule.public_holiday"), type: "public_holiday" },
])

function getColorByScheduleType(scheduleType) {
  return scheduleColors[scheduleType] || "#3788d9"
}

const scrollContainerRef = ref(null)
const topSentinelRef = ref(null)
const bottomSentinelRef = ref(null)

const weeks = ref([])
const allEvents = ref([])
const eventBindUserID = ref([])

const loadingTop = ref(false)
const loadingBottom = ref(false)
let isLoadingMore = false

const selectedEventId = ref(0)
const selectedDate = ref('')
const lastClickTime = ref(0)
const lastClickDateStr = ref('')
const lastEventClickTime = ref(0)
const lastEventClickId = ref(0)
const submitChecked = ref(false)
let lastEventsSnapshot = null

const showModal = ref(false)
const showDeleteModal = ref(false)
const eventData = ref({
  id: 0,
  title: "",
  startDate: "",
  endDate: "",
  color: "#3788d9",
  scheduleType: "work",
  isPublic: false,
  isEditing: false,
  isEditable: false,
})

const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  eventInfo: null,
  targetDate: '',
})

const autoScrollToday = ref(localStorage.getItem('calendarStreamAutoScrollToday') === '1')

watch(autoScrollToday, v => {
  localStorage.setItem('calendarStreamAutoScrollToday', v ? '1' : '0')
})

const clipboard = ref(null)

const dragState = ref({ eventId: 0, originalStart: '', originalEnd: '', dragging: false })
const dragOverDate = ref('')

const weekdayShort = computed(() => {
  if (locale.value === 'zh-CN') return ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
  return ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
})

function dateToStr(date) {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

function strToDate(dateStr) {
  return new Date(dateStr)
}

function addDays(date, n) {
  const d = new Date(date)
  d.setDate(d.getDate() + n)
  return d
}

function isSameDay(a, b) {
  return dateToStr(a) === dateToStr(b)
}

function isToday(date) {
  return isSameDay(date, new Date())
}

function isWeekend(date) {
  const day = date.getDay()
  return day === 0 || day === 6
}

function getWeekStart(date) {
  const d = new Date(date)
  d.setHours(0, 0, 0, 0)
  const day = d.getDay()
  const diff = day === 0 ? -6 : 1 - day
  d.setDate(d.getDate() + diff)
  return d
}

function generateWeeks(startWeekStart, count) {
  const result = []
  for (let i = 0; i < count; i++) {
    const ws = new Date(startWeekStart)
    ws.setDate(ws.getDate() + i * 7)
    const days = []
    for (let j = 0; j < 7; j++) {
      days.push(addDays(ws, j))
    }
    result.push({ weekStart: ws, days })
  }
  return result
}

function generateInitialWeeks() {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const centerWeekStart = getWeekStart(today)
  const startWeekStart = addDays(centerWeekStart, -INITIAL_WEEKS_BEFORE * 7)
  weeks.value = generateWeeks(startWeekStart, INITIAL_WEEKS_BEFORE + 1 + INITIAL_WEEKS_AFTER)
}

function weekKey(week) {
  return dateToStr(week.weekStart)
}

function shouldShowMonthSeparator(index) {
  if (index === 0) return true
  const current = weeks.value[index].days[0]
  const prev = weeks.value[index - 1].days[0]
  return current.getMonth() !== prev.getMonth() || current.getFullYear() !== prev.getFullYear()
}

function getMonthLabel(week) {
  const d = week.days[0]
  if (locale.value === 'zh-CN') {
    return `${d.getFullYear()}年${d.getMonth() + 1}月`
  }
  const months = ['January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December']
  return `${months[d.getMonth()]} ${d.getFullYear()}`
}

function getLoadedRange() {
  if (weeks.value.length === 0) return { start: '', end: '' }
  return {
    start: dateToStr(weeks.value[0].days[0]),
    end: dateToStr(weeks.value[weeks.value.length - 1].days[6]),
  }
}

function transformEvent(item) {
  return {
    id: item.ID,
    title: item.Title,
    startDate: item.StartDate,
    endDate: item.EndDate,
    color: getColorByScheduleType(item.ScheduleType),
    scheduleType: item.ScheduleType || 'work',
    isPublic: item.IsPublic || false,
    canEdit: item.canEdit === true,
    userId: item.UserID,
  }
}

async function fetchEvents(start, end) {
  try {
    const { errCode, data } = await calendarApi.getEvents({
      calendar_id: calendarId.value,
      start: start,
      end: end,
    })
    if (errCode === 0 && data.list) {
      allEvents.value = allEvents.value.filter(e =>
        e.endDate < start || e.startDate > end
      )
      allEvents.value.push(...data.list.map(transformEvent))
      eventBindUserID.value = data.list.map(item => ({ eventID: item.ID, userID: item.UserID }))

      const newSnapshot = JSON.stringify(data.list)
      if (newSnapshot !== lastEventsSnapshot) {
        lastEventsSnapshot = newSnapshot
        nextTick(() => recalcScrollTitles())
      }
    }
  } catch {
  }
}

function getVisibleRange() {
  if (weeks.value.length === 0) return { start: '', end: '' }
  const container = scrollContainerRef.value
  if (!container) return getLoadedRange()

  const rect = container.getBoundingClientRect()
  let startDate = ''
  let endDate = ''

  const firstEl = document.elementFromPoint(rect.left + 10, rect.top + 60)
  const lastEl = document.elementFromPoint(rect.right - 10, rect.bottom - 10)

  const firstCell = firstEl?.closest('.day-cell')
  if (firstCell?.dataset?.date) startDate = firstCell.dataset.date

  const lastCell = lastEl?.closest('.day-cell')
  if (lastCell?.dataset?.date) endDate = lastCell.dataset.date

  if (!startDate) startDate = dateToStr(weeks.value[0].days[0])
  if (!endDate) endDate = dateToStr(weeks.value[weeks.value.length - 1].days[6])

  return { start: startDate, end: endDate }
}

async function refreshVisibleEvents() {
  const { start, end } = getVisibleRange()
  if (!start || !end) return
  try {
    const { errCode, data } = await calendarApi.getEvents({
      calendar_id: calendarId.value,
      start,
      end,
    })
    if (errCode === 0 && data.list) {
      const newSnapshot = JSON.stringify(data.list)
      if (newSnapshot !== lastEventsSnapshot) {
        lastEventsSnapshot = newSnapshot
        allEvents.value = allEvents.value.filter(e =>
          e.endDate < start || e.startDate > end
        )
        allEvents.value.push(...data.list.map(transformEvent))
        const existingMap = new Map(eventBindUserID.value.map(e => [e.eventID, e.userID]))
        data.list.forEach(item => existingMap.set(item.ID, item.UserID))
        eventBindUserID.value = Array.from(existingMap, ([eventID, userID]) => ({ eventID, userID }))
        nextTick(() => recalcScrollTitles())
      }
    }
  } catch {
  }
}

async function loadMoreWeeksTop() {
  if (isLoadingMore || loadingTop.value) return
  isLoadingMore = true
  loadingTop.value = true

  const firstWeek = weeks.value[0]
  const newStart = addDays(firstWeek.weekStart, -LOAD_BATCH * 7)
  const newWeeks = generateWeeks(newStart, LOAD_BATCH)

  const oldScrollHeight = scrollContainerRef.value?.scrollHeight || 0

  weeks.value.unshift(...newWeeks)

  const fetchStart = dateToStr(newStart)
  const fetchEnd = dateToStr(addDays(newStart, LOAD_BATCH * 7 - 1))
  await fetchEvents(fetchStart, fetchEnd)

  nextTick(() => {
    if (scrollContainerRef.value) {
      const newScrollHeight = scrollContainerRef.value.scrollHeight
      scrollContainerRef.value.scrollTop += (newScrollHeight - oldScrollHeight)
    }
    loadingTop.value = false
    isLoadingMore = false

    if (weeks.value.length > MAX_WEEKS) {
      const toRemove = weeks.value.length - MAX_WEEKS
      const oldH = scrollContainerRef.value?.scrollHeight || 0
      weeks.value.splice(-toRemove, toRemove)
      nextTick(() => {
        if (scrollContainerRef.value) {
          const newH = scrollContainerRef.value.scrollHeight
          scrollContainerRef.value.scrollTop += (newH - oldH)
        }
      })
    }
  })
}

async function loadMoreWeeksBottom() {
  if (isLoadingMore || loadingBottom.value) return
  isLoadingMore = true
  loadingBottom.value = true

  const lastWeek = weeks.value[weeks.value.length - 1]
  const newStart = addDays(lastWeek.weekStart, 7)
  const newWeeks = generateWeeks(newStart, LOAD_BATCH)

  weeks.value.push(...newWeeks)

  const fetchStart = dateToStr(newStart)
  const fetchEnd = dateToStr(addDays(newStart, LOAD_BATCH * 7 - 1))
  await fetchEvents(fetchStart, fetchEnd)

  loadingBottom.value = false
  isLoadingMore = false

  if (weeks.value.length > MAX_WEEKS) {
    const toRemove = weeks.value.length - MAX_WEEKS
    const oldH = scrollContainerRef.value?.scrollHeight || 0
    weeks.value.splice(0, toRemove)
    nextTick(() => {
      if (scrollContainerRef.value) {
        const newH = scrollContainerRef.value.scrollHeight
        scrollContainerRef.value.scrollTop += (newH - oldH)
      }
    })
  }
}

function weekdayIndex(dateStr) {
  return ((new Date(dateStr).getDay() + 6) % 7)
}

function buildWeekLayout(week) {
  const weekStart = dateToStr(week.days[0])
  const weekEnd = dateToStr(week.days[6])
  const segments = []

  for (const event of allEvents.value) {
    if (event.endDate < weekStart || event.startDate > weekEnd) continue
    const segStart = event.startDate < weekStart ? weekStart : event.startDate
    const segEnd = event.endDate > weekEnd ? weekEnd : event.endDate
    segments.push({
      event,
      segStart,
      segEnd,
      startCol: weekdayIndex(segStart),
      endCol: weekdayIndex(segEnd),
      continuesBefore: event.startDate < weekStart,
      continuesAfter: event.endDate > weekEnd,
    })
  }

  segments.sort((a, b) => {
    if (a.segStart !== b.segStart) return a.segStart < b.segStart ? -1 : 1
    if (a.event.startDate !== b.event.startDate) return a.event.startDate < b.event.startDate ? -1 : 1
    return a.event.title.localeCompare(b.event.title)
  })

  const laneEnds = []
  for (const seg of segments) {
    let lane = 0
    while (laneEnds[lane] !== undefined && laneEnds[lane] >= seg.segStart) lane++
    seg.lane = lane
    laneEnds[lane] = seg.segEnd
  }

  const maxLane = segments.reduce((m, s) => Math.max(m, s.lane), -1)
  const shownLanes = maxLane + 1
  const DATE_ROW_H = 34
  const LANE_H = 26

  return {
    segments,
    weekHeight: Math.max(90, DATE_ROW_H + shownLanes * LANE_H + 6),
  }
}

const weekLayouts = computed(() => {
  const map = new Map()
  for (const week of weeks.value) {
    map.set(weekKey(week), buildWeekLayout(week))
  }
  return map
})

function getEventBarClass(seg) {
  return {
    'event-bar-start': !seg.continuesBefore,
    'event-bar-end': !seg.continuesAfter,
    'event-bar-before': seg.continuesBefore,
    'event-bar-after': seg.continuesAfter,
    'event-selected': seg.event.id === selectedEventId.value,
    'event-draggable': seg.event.canEdit,
  }
}

function getUserIdFromEventID(eventID) {
  const target = eventBindUserID.value.find(item => item.eventID === eventID)
  return target ? target.userID : 0
}

function getUsernameFromUserID(userID) {
  if (userID == 0) return ""
  return usersStore.getUsernameFromUserID(userID)
}

function toDatetime(dateStr) {
  return dateStr ? dateStr + " 00:00:00" : ""
}

function openEventModal(dateStr, dataEnd, id = 0, title = "", color = "#3788d9", scheduleType = "work", isPublic = false, isEditing = false, isEditable = true) {
  eventData.value = {
    id, title,
    startDate: dateStr,
    endDate: dataEnd,
    color, scheduleType, isPublic, isEditing, isEditable,
  }
  showModal.value = true
}

function openEditModal(event) {
  const exclusiveEnd = event.startDate === event.endDate
    ? event.endDate
    : DateUtils.dateToStr(DateUtils.toCalendarEnd(event.endDate))
  openEventModal(
    event.startDate,
    exclusiveEnd,
    event.id,
    event.title,
    event.color,
    event.scheduleType,
    event.isPublic,
    true,
    event.canEdit,
  )
}

function closeEventModal() {
  showModal.value = false
}

function selectColor(colorValue) {
  if (eventData.value.isEditable) {
    eventData.value.color = colorValue
    const selectedColor = colorOptions.value.find(c => c.value === colorValue)
    if (selectedColor) {
      eventData.value.scheduleType = selectedColor.type
    }
  }
}

async function saveEvent() {
  if (!eventData.value.title.trim()) {
    submitChecked.value = true
    toast.warning(t('calendar.event_title_required'))
    return
  }
  submitChecked.value = false

  if (!eventData.value.startDate || !eventData.value.endDate) {
    toast.warning(t('schedule.date_required'))
    return
  }

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
        schedule_type: eventData.value.scheduleType,
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
        schedule_type: eventData.value.scheduleType,
        is_public: eventData.value.isPublic,
      })
    }

    if (result.errCode === 0) {
      toast.success(t('calendar.event_save_success'))
      closeEventModal()
      refreshVisibleEvents()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
  }
}

async function deleteEvent() {
  showDeleteModal.value = true
}

async function confirmDeleteEvent() {
  try {
    const result = await calendarApi.deleteEvent(eventData.value.id)
    if (result.errCode === 0) {
      toast.success(t("calendar.event_delete_success"))
      closeEventModal()
      refreshVisibleEvents()
    } else {
      toast.error(t("message.server_error"))
    }
  } catch {
  } finally {
    showDeleteModal.value = false
  }
}

function unselectAllEvents() {
  selectedEventId.value = 0
}

function handleDateClick(dateStr) {
  const now = Date.now()
  unselectAllEvents()
  selectedDate.value = dateStr

  if (dateStr === lastClickDateStr.value) {
    if (now - lastClickTime.value < 400 && now - lastClickTime.value > 0) {
      if (userStore.isLoggedIn) {
        openEventModal(dateStr, dateStr)
      } else {
        toast.warning(t("message.login_to_your_account"))
        router.replace("/login?redirect=/calendar/" + calendarId.value + "/stream")
      }
    }
  }
  lastClickDateStr.value = dateStr
  lastClickTime.value = now
}

function handleEventClick(event) {
  const now = Date.now()
  unselectAllEvents()
  selectedEventId.value = event.id
  selectedDate.value = ''

  if (event.id === lastEventClickId.value) {
    if (now - lastEventClickTime.value < 400 && now - lastEventClickTime.value > 0) {
      openEditModal(event)
      unselectAllEvents()
    }
  }
  lastEventClickId.value = event.id
  lastEventClickTime.value = now
}

function handleDayContextMenu(e, dateStr) {
  contextMenu.value = {
    visible: true,
    x: e.clientX,
    y: e.clientY,
    eventInfo: null,
    targetDate: dateStr,
  }
}

function handleEventContextMenu(e, event) {
  contextMenu.value = {
    visible: true,
    x: e.clientX,
    y: e.clientY,
    eventInfo: {
      id: event.id,
      title: event.title,
      start: event.startDate,
      end: event.endDate,
      color: event.color,
      scheduleType: event.scheduleType,
      isPublic: event.isPublic,
      canEdit: event.canEdit,
    },
    targetDate: event.startDate,
  }
}

function closeContextMenu() {
  contextMenu.value.visible = false
}

function copyEvent() {
  if (!contextMenu.value.eventInfo) return
  clipboard.value = { ...contextMenu.value.eventInfo }
  toast.success(t('calendar.copy_success'))
  closeContextMenu()
}

async function pasteEvent() {
  if (!clipboard.value) {
    toast.warning(t('calendar.no_event_to_paste'))
    closeContextMenu()
    return
  }

  let targetStart = contextMenu.value.targetDate
  if (!targetStart) {
    closeContextMenu()
    return
  }
  closeContextMenu()

  await pasteToDate(targetStart)
}

async function pasteToDate(targetStart) {
  if (!clipboard.value) return

  const origStart = clipboard.value.start
  const origEnd = clipboard.value.end

  const origStartDate = strToDate(origStart)
  const origEndDate = strToDate(origEnd)
  const diffDays = Math.round((origEndDate - origStartDate) / 86400000)
  const isSameDay = diffDays === 0

  let targetEnd = targetStart
  if (!isSameDay) {
    const targetEndDate = addDays(strToDate(targetStart), diffDays)
    targetEnd = dateToStr(targetEndDate)
  }

  try {
    const result = await calendarApi.addEvent({
      calendar_id: calendarId.value,
      title: clipboard.value.title,
      start: toDatetime(targetStart),
      end: toDatetime(
        isSameDay ? targetEnd : DateUtils.toRealEnd(targetEnd),
      ),
      schedule_type: clipboard.value.scheduleType,
      is_public: clipboard.value.isPublic,
    })
    if (result.errCode === 0) {
      toast.success(t('calendar.paste_success'))
      refreshVisibleEvents()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
  }
}

function handleKeyDown(e) {
  if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return

  if (e.ctrlKey || e.metaKey) {
    if (e.key === 'c') {
      if (selectedEventId.value) {
        const event = allEvents.value.find(ev => ev.id === selectedEventId.value)
        if (event) {
          clipboard.value = {
            id: event.id,
            title: event.title,
            start: event.startDate,
            end: event.endDate,
            color: event.color,
            scheduleType: event.scheduleType,
            isPublic: event.isPublic,
          }
          toast.success(t('calendar.copy_success'))
        }
      }
    } else if (e.key === 'v') {
      if (clipboard.value && selectedDate.value) {
        e.preventDefault()
        pasteToDate(selectedDate.value)
      }
    }
  }
}

function handleDragStart(seg, e) {
  const event = seg.event
  if (!event.canEdit) {
    e.preventDefault()
    return
  }
  e.dataTransfer.effectAllowed = 'move'
  e.dataTransfer.setData('text/plain', String(event.id))

  const barRect = e.currentTarget.getBoundingClientRect()
  const spanDays = seg.endCol - seg.startCol + 1
  const grabDay = Math.max(0, Math.min(spanDays - 1,
    Math.floor(((e.clientX - barRect.left) / barRect.width) * spanDays)))
  const grabDate = addDays(strToDate(seg.segStart), grabDay)
  const grabOffset = Math.round((grabDate - strToDate(event.startDate)) / 86400000)

  dragState.value = {
    eventId: event.id,
    originalStart: event.startDate,
    originalEnd: event.endDate,
    grabOffset,
    dragging: true,
  }
}

function getDateFromEvent(e, week) {
  const gridEl = e.currentTarget
  const rect = gridEl.getBoundingClientRect()
  const col = Math.max(0, Math.min(6, Math.floor((e.clientX - rect.left) / (rect.width / 7))))
  return dateToStr(addDays(week.days[0], col))
}

function handleWeekDragOver(e, week) {
  dragOverDate.value = getDateFromEvent(e, week)
}

function handleWeekDragLeave(e) {
  if (!e.currentTarget.contains(e.relatedTarget)) {
    dragOverDate.value = ''
  }
}

function handleWeekDrop(e, week) {
  if (!dragState.value.dragging) return
  handleDrop(getDateFromEvent(e, week))
}

function handleDrop(dateStr) {
  if (!dragState.value.dragging) return

  const event = allEvents.value.find(ev => ev.id === dragState.value.eventId)
  if (!event) {
    dragState.value.dragging = false
    dragOverDate.value = ''
    return
  }

  const origStart = strToDate(event.startDate)
  const targetDate = strToDate(dateStr)
  const diff = Math.round((targetDate - origStart) / 86400000)
  const grabOffset = dragState.value.grabOffset ?? 0

  const newStart = dateToStr(addDays(targetDate, -grabOffset))
  const newEnd = dateToStr(addDays(strToDate(event.endDate), diff - grabOffset))

  calendarApi.updateEvent({
    id: event.id,
    title: event.title,
    start: toDatetime(newStart),
    end: toDatetime(newEnd),
    schedule_type: event.scheduleType,
  }).then(r => {
    if (r.errCode === 0) {
      toast.success(t('calendar.event_save_success'))
      refreshVisibleEvents()
    } else {
      toast.error(t('message.server_error'))
    }
  })

  dragState.value.dragging = false
  dragOverDate.value = ''
}

function handleDragEnd() {
  dragState.value.dragging = false
  dragOverDate.value = ''
}

function goToToday() {
  generateInitialWeeks()
  lastEventsSnapshot = null
  if (scrollContainerRef.value) scrollContainerRef.value.scrollTop = 0
  const { start, end } = getLoadedRange()
  fetchEvents(start, end).then(() => scrollToToday())
}

function getMsUntilMidnight() {
  const now = new Date()
  const next = new Date(now)
  next.setHours(24, 0, 0, 0)
  return Math.max(1000, next - now)
}

function onMidnight() {
  if (autoScrollToday.value) goToToday()
  scheduleMidnightTimer()
}

function scheduleMidnightTimer() {
  if (midnightTimer) clearTimeout(midnightTimer)
  midnightTimer = setTimeout(onMidnight, getMsUntilMidnight())
}

function scrollToToday() {
  nextTick(() => {
    requestAnimationFrame(() => {
      applyTodayScroll()
      requestAnimationFrame(() => {
        applyTodayScroll()
      })
    })
  })
}

function applyTodayScroll() {
  const container = scrollContainerRef.value
  if (!container) return
  const todayStr = dateToStr(new Date())
  const todayCell = container.querySelector(`.day-cell[data-date="${todayStr}"]`)
  if (!todayCell) return

  const idx = weeks.value.findIndex(w => weekKey(w) === dateToStr(getWeekStart(new Date())))
  let desiredTop = idx > 0
    ? (weekLayouts.get(weekKey(weeks.value[idx - 1]))?.weekHeight ?? 90)
    : 0

  const todayWeekEl = todayCell.closest('.week-row')
  const prevSib = todayWeekEl?.previousElementSibling
  if (prevSib?.classList.contains('month-separator')) {
    desiredTop += prevSib.getBoundingClientRect().height
  }

  const contRect = container.getBoundingClientRect()
  const cellRect = todayCell.getBoundingClientRect()
  const currentTop = cellRect.top - contRect.top

  container.scrollTop += (currentTop - desiredTop)
}

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
      const container = scrollContainerRef.value
      if (!container) return
      container.querySelectorAll(".event-title").forEach(applyScrollToTitle)
    })
  })
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
  } finally {
    loading.value = false
  }
}

let intersectionObserver = null
let refreshTimer = null
let resizeObserver = null
let midnightTimer = null

onMounted(() => {
  fetchCalendarInfo()
  generateInitialWeeks()
  scheduleMidnightTimer()

  const { start, end } = getLoadedRange()
  fetchEvents(start, end).then(() => scrollToToday())

  nextTick(() => {
    if (topSentinelRef.value && bottomSentinelRef.value && scrollContainerRef.value) {
      intersectionObserver = new IntersectionObserver((entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            if (entry.target === topSentinelRef.value) {
              loadMoreWeeksTop()
            } else if (entry.target === bottomSentinelRef.value) {
              loadMoreWeeksBottom()
            }
          }
        }
      }, {
        root: scrollContainerRef.value,
        rootMargin: '500px 0px 500px 0px',
      })
      intersectionObserver.observe(topSentinelRef.value)
      intersectionObserver.observe(bottomSentinelRef.value)
    }
  })

  document.addEventListener('click', closeContextMenu)
  document.addEventListener('keydown', handleKeyDown)

  refreshTimer = setInterval(() => {
    refreshVisibleEvents()
  }, 5000)

  let resizeTimer = null
  resizeObserver = new ResizeObserver(() => {
    clearTimeout(resizeTimer)
    resizeTimer = setTimeout(() => recalcScrollTitles(), 150)
  })
  if (scrollContainerRef.value) {
    resizeObserver.observe(scrollContainerRef.value)
  }

  onBeforeUnmount(() => {
    if (midnightTimer) {
      clearTimeout(midnightTimer)
      midnightTimer = null
    }
    if (refreshTimer) {
      clearInterval(refreshTimer)
      refreshTimer = null
    }
    if (intersectionObserver) {
      intersectionObserver.disconnect()
      intersectionObserver = null
    }
    if (resizeObserver) {
      resizeObserver.disconnect()
      resizeObserver = null
    }
    document.removeEventListener('click', closeContextMenu)
    document.removeEventListener('keydown', handleKeyDown)
    clearTimeout(resizeTimer)
  })
})

watch(locale, () => {
  colorOptions.value = [
    { value: "#066fd1", label: t("schedule.work"), type: "work" },
    { value: "#09d119", label: t("schedule.duty"), type: "duty" },
    { value: "#ff00ff", label: t("schedule.exam"), type: "exam" },
    { value: "#ffca3a", label: t("schedule.standby"), type: "standby" },
    { value: "#d16c13", label: t("schedule.personal_holiday"), type: "personal_holiday" },
    { value: "#d10d21", label: t("schedule.public_holiday"), type: "public_holiday" },
  ]
})
</script>

<template>
  <div class="flex w-full flex-col" style="height: calc(100vh - 3.5rem)">
    <!-- Event Modal -->
    <div
      v-if="showModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-gray-800/20"
    >
      <div
        class="modal-content bg-white rounded-lg shadow-lg w-full max-w-2xl max-h-[95vh] flex flex-col"
      >
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
          <div
            v-if="eventData.isEditing && userStore.isLoggedIn && getUserIdFromEventID(eventData.id)"
            class="absolute left-1/2 -translate-x-1/2 flex items-center gap-2"
          >
            <img
              :src="usersStore.getAvatarUrlFromUserID(getUserIdFromEventID(eventData.id))"
              class="h-6 w-6 rounded-full"
              alt="avatar"
            />
            <span class="text-sm text-gray-500">
              {{ t("calendar.created_by", { name: getUsernameFromUserID(getUserIdFromEventID(eventData.id)) }) }}
            </span>
          </div>
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
            >
              <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
              <path d="M18 6l-12 12"></path>
              <path d="M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <div class="modal-body p-4 flex-1 overflow-y-auto">
          <DatatimePickerForFullCalendar
            v-model:startDate="eventData.startDate"
            v-model:endDate="eventData.endDate"
            :color="eventData.color"
            :title="eventData.title"
            :isEditable="eventData.isEditable"
          />

          <div class="mb-4">
            <div class="uni-easyinput input relative">
              <div
                class="uni-easyinput__content is-input-border border border-gray-300 rounded-md bg-white relative"
                :class="{
                  'border-gray-300': eventData.title || !submitChecked,
                  'border-red-500': !eventData.title && submitChecked,
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

          <div class="mb-4 px-1">
            <span class="text-sm text-gray-600">{{ t('schedule.event_type') }}: </span>
            <span class="text-sm font-medium" :style="{ color: eventData.color }">{{ t('schedule.' + eventData.scheduleType) || eventData.scheduleType }}</span>
          </div>

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

    <!-- Header Bar -->
    <div class="grid grid-cols-3 items-center border-b border-gray-200 px-4 py-2 bg-white dark:bg-dk-card dark:border-dk-muted flex-shrink-0">
      <div class="flex items-center gap-2 justify-start">
        <button
          @click="router.push('/calendars')"
          class="inline-flex items-center gap-1.5 rounded-lg border border-gray-300 px-3 py-1.5 text-base font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:text-gray-300 dark:hover:bg-dk-base"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
          {{ t('calendar.calendars') }}
        </button>
        <button
          @click="goToToday"
          class="inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-3 py-1.5 text-base font-medium text-white transition-colors hover:bg-blue-700"
        >
          {{ t('schedule.today') }}
        </button>
        <RouterLink
          :to="`/calendar/${calendarId}`"
          class="inline-flex items-center gap-1.5 rounded-lg border border-gray-300 px-3 py-1.5 text-base font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-dk-muted dark:text-gray-300 dark:hover:bg-dk-base"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
          {{ t('calendar.classic_view') }}
        </RouterLink>
      </div>
      <div class="text-center">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white">{{ calendarInfo?.Name || '' }}</h2>
      </div>
      <div class="flex items-center justify-end">
        <label class="flex items-center gap-1.5 text-base text-gray-600 dark:text-gray-300 cursor-pointer select-none">
          <input
            v-model="autoScrollToday"
            type="checkbox"
            class="rounded border-gray-300"
          />
          {{ t('calendar.auto_scroll_today') }}
        </label>
      </div>
    </div>

    <!-- Weekday Header -->
    <div class="grid grid-cols-7 border-b border-gray-200 bg-gray-50 dark:bg-dk-base dark:border-dk-muted flex-shrink-0">
      <div
        v-for="(label, i) in weekdayShort"
        :key="i"
        class="py-1.5 text-center text-base font-bold text-gray-500 dark:text-gray-400"
      >
        {{ label }}
      </div>
    </div>

    <!-- Scrollable Container -->
    <div
      ref="scrollContainerRef"
      class="flex-1 overflow-y-auto bg-white dark:bg-dk-card"
    >
      <!-- Top Sentinel -->
      <div ref="topSentinelRef" class="h-px"></div>

      <!-- Loading indicator (top) -->
      <div v-if="loadingTop" class="py-2 text-center text-base text-gray-400">
        <svg class="inline animate-spin h-4 w-4 mr-1" viewBox="0 0 24 24" fill="none">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
        </svg>
        {{ t('calendar.loading_more') }}
      </div>

      <!-- Week Rows -->
      <template v-for="(week, i) in weeks" :key="weekKey(week)">
        <!-- Month Separator -->
        <div
          v-if="shouldShowMonthSeparator(i)"
          class="month-separator sticky top-0 z-[5] py-1 px-4 text-base font-medium text-gray-500 bg-gray-100 border-b border-gray-200 dark:bg-dk-base dark:text-gray-400 dark:border-dk-muted"
        >
          {{ getMonthLabel(week) }}
        </div>

        <!-- Week Row -->
        <div class="week-row border-b border-gray-100 dark:border-dk-muted">
          <div
            class="relative"
            :style="{ height: weekLayouts.get(weekKey(week)).weekHeight + 'px' }"
            @dragover.prevent="handleWeekDragOver($event, week)"
            @dragleave="handleWeekDragLeave"
            @drop.prevent="handleWeekDrop($event, week)"
          >
            <!-- Day cells (background layer) -->
            <div class="absolute inset-0 grid grid-cols-7">
              <div
                v-for="day in week.days"
                :key="dateToStr(day)"
                class="day-cell border-r border-gray-100 p-1 dark:border-dk-muted relative"
                :class="{
                  'bg-gray-50 dark:bg-dk-base/50': isWeekend(day),
                  'day-today': isToday(day),
                  'drag-over': dragOverDate === dateToStr(day),
                }"
                :data-date="dateToStr(day)"
                @click="handleDateClick(dateToStr(day))"
                @contextmenu.prevent="handleDayContextMenu($event, dateToStr(day))"
              >
                <!-- Date number -->
                <div class="flex justify-center mb-0.5">
                  <span
                    class="inline-flex items-center justify-center text-base leading-none w-6 h-6 rounded-full text-gray-600 dark:text-gray-300"
                  >
                    {{ day.getDate() }}
                  </span>
                </div>
              </div>
            </div>

            <!-- Event bars (overlay layer) -->
            <div class="absolute inset-0" style="pointer-events: none">
              <div
                v-for="seg in weekLayouts.get(weekKey(week)).segments"
                :key="seg.event.id + '-' + seg.segStart"
                class="event-chip absolute text-sm px-1.5 cursor-pointer select-none"
                :class="getEventBarClass(seg)"
                :style="{
                  backgroundColor: seg.event.color,
                  color: '#fff',
                  left: seg.startCol * 14.2857 + '%',
                  width: (seg.endCol - seg.startCol + 1) * 14.2857 + '%',
                  top: 34 + seg.lane * 26 + 'px',
                  pointerEvents: 'auto',
                }"
                draggable="true"
                @click.stop="handleEventClick(seg.event)"
                @contextmenu.prevent.stop="handleEventContextMenu($event, seg.event)"
                @dragstart="handleDragStart(seg, $event)"
                @dragend="handleDragEnd"
              >
                <span v-if="seg.continuesBefore" class="event-continue-left">«</span>
                <span class="event-title" :class="{ 'pl-3': seg.continuesBefore, 'pr-3': seg.continuesAfter }">{{ seg.event.title }}</span>
                <span v-if="seg.continuesAfter" class="event-continue-right">»</span>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- Loading indicator (bottom) -->
      <div v-if="loadingBottom" class="py-2 text-center text-base text-gray-400">
        <svg class="inline animate-spin h-4 w-4 mr-1" viewBox="0 0 24 24" fill="none">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
        </svg>
        {{ t('calendar.loading_more') }}
      </div>

      <!-- Bottom Sentinel -->
      <div ref="bottomSentinelRef" class="h-px"></div>
    </div>

    <!-- Context Menu -->
    <div
      v-if="contextMenu.visible"
      class="fixed z-[60] min-w-[140px] rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-dk-muted dark:bg-dk-card"
      :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
    >
      <button
        v-if="contextMenu.eventInfo"
        @click="copyEvent"
        class="flex w-full items-center gap-2 px-3 py-2 text-base text-gray-700 hover:bg-gray-100 dark:text-dk-text dark:hover:bg-dk-base"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
        {{ t('calendar.copy_event') }}
      </button>
      <button
        @click="pasteEvent"
        class="flex w-full items-center gap-2 px-3 py-2 text-base text-gray-700 hover:bg-gray-100 dark:text-dk-text dark:hover:bg-dk-base"
        :class="{ 'opacity-40 cursor-not-allowed': !clipboard }"
        :disabled="!clipboard"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path><rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect></svg>
        {{ t('calendar.paste_event') }}
      </button>
    </div>

    <!-- Delete Confirm Dialog -->
    <ConfirmDialog
      v-model="showDeleteModal"
      :title="t('calendar.delete_event')"
      :message="t('calendar.confirm_delete_event')"
      :confirm-text="t('delete')"
      :cancel-text="t('cancel')"
      danger
      @confirm="confirmDeleteEvent"
    />
  </div>
</template>

<style scoped>
.event-chip {
  border: 1px solid rgba(0, 0, 0, 0.1);
  white-space: nowrap;
  overflow: hidden;
  height: 22px;
  line-height: 20px;
}

.event-bar-start {
  border-top-left-radius: 4px;
  border-bottom-left-radius: 4px;
}

.event-bar-end {
  border-top-right-radius: 4px;
  border-bottom-right-radius: 4px;
}

.event-bar-before {
  border-left: none;
}

.event-bar-after {
  border-right: none;
}

.event-selected {
  outline: 2px solid #000;
  outline-offset: -1px;
}

.event-draggable {
  cursor: grab;
}

.event-draggable:active {
  cursor: grabbing;
}

.event-title {
  display: block;
  white-space: nowrap;
  overflow: visible;
  will-change: transform;
}

.event-title[data-truncated="true"] {
  display: inline-block;
  animation: marquee-bounce 6s ease-in-out infinite;
}

@keyframes marquee-bounce {
  0%, 20% { transform: translateX(0); }
  60% { transform: translateX(var(--scroll-distance, 0px)); }
  80% { transform: translateX(var(--scroll-distance, 0px)); }
  100% { transform: translateX(0); }
}

.event-continue-left,
.event-continue-right {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  font-size: 10px;
  line-height: 1;
  opacity: 0.85;
}

.event-continue-left {
  left: 2px;
}

.event-continue-right {
  right: 2px;
}

.day-cell {
  transition: background-color 0.1s;
}

.day-cell.day-today {
  background-color: #fefce8;
}

.dark .day-cell.day-today {
  background-color: rgba(255, 220, 40, 0.15);
}

.day-cell.drag-over {
  background-color: rgb(219 234 254);
}

.dark .day-cell.drag-over {
  background-color: rgba(59, 130, 246, 0.3);
}

.day-cell:hover {
  background-color: rgba(59, 130, 246, 0.03);
}

.week-row {
  contain: layout style;
}

:deep(.dark) .event-chip {
  border-color: rgba(255, 255, 255, 0.15);
}

:deep(.dark) .event-selected {
  outline-color: #fff;
}
</style>
