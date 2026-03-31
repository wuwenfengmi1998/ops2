<script setup>
import { ref, watch } from 'vue'
import FullCalendar from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import timeGridPlugin from '@fullcalendar/timegrid'
import interactionPlugin from '@fullcalendar/interaction'
import listPlugin from '@fullcalendar/list'
import { useI18n } from 'vue-i18n'
import { usePageTitle } from '@/composables/usePageTitle'

usePageTitle('appname.schedule')
const { t, locale } = useI18n()

const calendarRef = ref(null)

const calendarOptions = ref({
  height: '100%',
  contentHeight: 'auto',
  locale: locale.value,
  plugins: [dayGridPlugin, timeGridPlugin, interactionPlugin, listPlugin],
  nowIndicator: true,
  weekends: true,
  initialView: 'dayGridMonth',
  selectable: true,
  editable: true,
  dayMaxEvents: true,
  navLinks: true,
  firstDay: 1,
  expandRows: true,
  stickyHeaderDates: true,

  dayCellDidMount(info) {
    if (info.date.getDay() === 0 || info.date.getDay() === 6) {
      info.el.style.backgroundColor = '#f5f5f5'
    }
    info.el.style.border = '1px solid #e5e7eb'
  },

  headerToolbar: {
    left: 'prevYear,prev,next,nextYear',
    center: 'title',
    right: '',
  },

  customButtons: {
    prevYear: {
      text: t('schedule.previous_year'),
      click() { calendarRef.value.getApi().prevYear() },
    },
    nextYear: {
      text: t('schedule.next_year'),
      click() { calendarRef.value.getApi().nextYear() },
    },
    prevMonth: {
      text: t('schedule.previous_month'),
      click() { calendarRef.value.getApi().prev() },
    },
    nextMonth: {
      text: t('schedule.next_month'),
      click() { calendarRef.value.getApi().next() },
    },
    week: {
      text: t('schedule.week'),
      click() { calendarRef.value.getApi().changeView('timeGridWeek') },
    },
  },

  events: [
    { title: 'Event1', date: '2025-11-10' },
    { title: 'Event2', date: '2025-11-15', end: '2025-11-17' },
    { title: 'Event3', date: '2025-11-20T10:30:00', end: '2025-11-20T12:30:00' },
  ],
})

watch(locale, () => {
  calendarOptions.value.locale = locale.value
  calendarOptions.value.headerToolbar.customButtons.prevYear.text = t('schedule.previous_year')
  calendarOptions.value.headerToolbar.customButtons.nextYear.text = t('schedule.next_year')
  calendarOptions.value.headerToolbar.customButtons.prevMonth.text = t('schedule.previous_month')
  calendarOptions.value.headerToolbar.customButtons.nextMonth.text = t('schedule.next_month')
  calendarOptions.value.headerToolbar.customButtons.week.text = t('schedule.week')
})
</script>

<template>
  <div class="flex h-[calc(100vh-3.5rem)] w-full flex-col">
    <div class="flex-1 rounded-lg border border-gray-200 bg-white p-0.5 shadow dark:border-dk-muted dark:bg-dk-card">
      <div class="h-full w-full overflow-hidden rounded-md">
        <FullCalendar ref="calendarRef" :options="calendarOptions" />
      </div>
    </div>
  </div>
</template>
