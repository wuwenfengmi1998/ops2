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
    left: 'prevYear,prev,today,next,nextYear',
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
    prev: {
      text: t('schedule.previous_month'),
      click() { calendarRef.value.getApi().prev() },
    },
    next: {
      text: t('schedule.next_month'),
      click() { calendarRef.value.getApi().next() },
    },
    today: {
      text: t('schedule.today'),
      click() { calendarRef.value.getApi().today() },
    },
    week: {
      text: t('schedule.week'),
      click() { calendarRef.value.getApi().changeView('timeGridWeek') },
    },
  },

  events: [

  ],
})

watch(locale, () => {
  calendarOptions.value.locale = locale.value
  calendarOptions.value.customButtons.prevYear.text = t('schedule.previous_year')
  calendarOptions.value.customButtons.nextYear.text = t('schedule.next_year')
  calendarOptions.value.customButtons.prev.text = t('schedule.previous_month')
  calendarOptions.value.customButtons.next.text = t('schedule.next_month')
  calendarOptions.value.customButtons.today.text = t('schedule.today')
  calendarOptions.value.customButtons.week.text = t('schedule.week')
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
