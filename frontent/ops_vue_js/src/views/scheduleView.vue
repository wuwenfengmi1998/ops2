<script setup>
import FullCalendar from "@fullcalendar/vue3";
import dayGridPlugin from "@fullcalendar/daygrid";
import timeGridPlugin from "@fullcalendar/timegrid";
import interactionPlugin from "@fullcalendar/interaction"; //拖动插件 需要用npm安装
import listPlugin from '@fullcalendar/list';

import { onMounted, watch, ref } from "vue";
import { useI18n } from "vue-i18n";

const { t, locale } = useI18n();

const calendarOptions = ref({
  height: "auto",
  locale: locale.value,
  plugins: [
    dayGridPlugin,
    timeGridPlugin,
    interactionPlugin, //导入拖动插件
    listPlugin,
  ],
  fixedWeekCount: false, //是否固定显示6行
  weekNumbers: true,
  initialView: "dayGridMonth", //默认月视图 dayGridMonth timeGridWeek listWeek
  editable: true,
  selectable: true,
  firstDay: 1,
  
  dayCellDidMount(info) {
    switch (info.dow) {
      case 0:
        info.el.style.backgroundColor = "#ffb5b5";
        break;
      case 6:
        info.el.style.backgroundColor = "#ffb5b5";
        break;
    }

    if (info.isToday) {
      //info.el.style.backgroundColor = '#ffff7f';
    }

    info.el.style.border = "1px solid #4b4b4b"; // 浅蓝色边框
  },

  events: [
    { title: "事件 1", start: "2025-11-10" },
    { title: "事件 2", start: "2025-11-15", end: "2024-06-17" },
    { title: "事件 3", start: "2025-11-20T10:30:00", end: "2024-06-20T12:30:00" },

  ],
});

function functionupdataTitle() {
  document.title = "Operations." + t("appname.schedule");
}

// 监听语言变化，更新标题
watch(locale, () => {
  functionupdataTitle();
  calendarOptions.value.locale = locale.value;
});

onMounted(() => {
  functionupdataTitle();
});
</script>
<template>
  <FullCalendar :options="calendarOptions" />
</template>
