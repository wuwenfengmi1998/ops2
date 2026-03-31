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
  height: 'auto',
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
  <div class="mx-auto max-w-6xl px-6 py-6">
    <h2 class="mb-6 text-2xl font-bold text-gray-900 dark:text-white">{{ t('schedule.my_schedule') }}</h2>
    <div class="rounded-xl border border-gray-200 bg-white px-4 py-4 shadow-lg dark:border-dk-muted dark:bg-dk-card">
      <FullCalendar ref="calendarRef" :options="calendarOptions" />
    </div>
  </div>
</template>
