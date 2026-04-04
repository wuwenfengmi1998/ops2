<script setup>
import { ref, watch, defineProps, onMounted } from "vue";

// FullCalendar Vue 3 组件
import FullCalendar from "@fullcalendar/vue3";
// FullCalendar 插件：月视图
import dayGridPlugin from "@fullcalendar/daygrid";
// FullCalendar 插件：交互功能（拖拽、点击等）
import interactionPlugin from "@fullcalendar/interaction";

import { useDateUtils } from "@/composables/useDateUtils";

const DateUtils = useDateUtils();

// FullCalendar 组件的引用，用于调用日历 API
const calendarRef = ref(null);

// 用于跟踪上次点击event时间的响应式变量
const lastEventClickTime = ref(0);

var firstClickOnDate = false;
var dataStartTemp = "";

// 国际化 hook
import { useI18n } from "vue-i18n";
// 获取国际化翻译函数和当前语言
const { t, locale } = useI18n();

const isShow = ref(false);

// 定义props：从父组件接收日期数据
const props = defineProps({
  startDate: {
    type: String,
    required: false,
    default: "",
  },
  endDate: {
    type: String,
    required: false,
    default: "",
  },
  title: {
    type: String,
    required: false,
    default: "",
  },
  color: {
    type: String,
    required: false,
    default: "",
  },
  isEditable:{
    type:Boolean,
    required: false,
    default: true,
  }
});

const eventData = ref({
  id: "0",
  title: props.title,
  start: props.startDate,
  end: props.endDate,
  backgroundColor: props.color,
  borderColor: props.color,
  allDay: true,
  editable: true,
});

function passing_date_characters(startDate, endDate) {
  //需要先判断大小确定哪个是开始日期哪个是结束日期
  var tempStart = DateUtils.strToDate(startDate);
  var tempEnd = DateUtils.strToDate(endDate);

  if (tempStart > tempEnd) {
    //反了
    eventData.value.start = endDate;
    //end 需要加一天
    eventData.value.end = DateUtils.dateToStr(
      DateUtils.toCalendarEnd(startDate),
    );
  } else {
    eventData.value.start = startDate;
    //end 需要加一天
    eventData.value.end = DateUtils.dateToStr(DateUtils.toCalendarEnd(endDate));
  }
}
function passing_date_characters_Select(startDate, endDate) {
  //滑动选择日期的参数来自FullCalendar，不需要加一天也不需要判断大小，而且直接就是字符串类型
 
  eventData.value.start = startDate;
  eventData.value.end = endDate===""?startDate:endDate;
}

// 监听props变化，更新本地eventData
watch(
  () => props.start,
  (newVal) => {
    eventData.value.start = newVal;
  },
);

watch(
  () => props.end,
  (newVal) => {
    eventData.value.end = newVal;
  },
);

watch(
  () => props.title,
  (newVal) => {
    eventData.value.title = newVal;
  },
);

watch(
  () => props.color,
  (newVal) => {
    eventData.value.backgroundColor = newVal;
    eventData.value.borderColor = newVal;
  },
);

// 定义事件发射：通知父组件日期变化
const emit = defineEmits(["update:startDate", "update:endDate", "clearDates"]);


// 监听本地eventData变化，同步更新到父组件
watch(
  () => eventData.value.start,
  (newVal) => {
    emit("update:startDate", newVal);
  },
);

watch(
  () => eventData.value.end,
  (newVal) => {
    emit("update:endDate", newVal);
  },
);

// 日历配置选项
const calendarOptions = ref({
  // 日历高度：占满可用空间
  //height: "300px",
  // 内容高度自适应
  //contentHeight: "auto",
  // 使用当前应用语言
  locale: locale.value,
  // 注册使用的插件
  plugins: [dayGridPlugin, interactionPlugin],
  // 显示当前时间指示线
  nowIndicator: true,
  // 显示周末
  weekends: true,
  // 初始视图：月视图
  initialView: "dayGridMonth",
  // 允许选择日期/时间段
  selectable: true,
  // 允许拖拽调整事件
  editable: true,
  // 日期格中事件过多时显示"+N more"
  dayMaxEvents: true,

  // 日期标题可点击跳转 //不跳转
  //navLinks: true,

  // 一周的第一天：1 = 周一
  firstDay: 1,
  // 自动展开行高
  //expandRows: true,
  // 固定头部日期
  stickyHeaderDates: true,

  // 日期格子挂载时的样式处理
  dayCellDidMount(info) {
    // 周六周日显示灰色背景
    if (info.date.getDay() === 0 || info.date.getDay() === 6) {
      info.el.style.backgroundColor = "#f5f5f5";
    }
    // 添加边框样式
    info.el.style.border = "1px solid #e5e7eb";
  },

  // 顶部工具栏配置
  headerToolbar: {
    // 左侧：年份和月份导航按钮
    left: "prev,today,next",
    // 中间：标题（显示当前月份/年份）
    center: "title",
    // 右侧：留空（可通过 customButtons 扩展）
    right: "",
  },

  // 自定义按钮：扩展工具栏功能
  customButtons: {
    // 上一年按钮
    prevYear: {
      text: t("schedule.previous_year"),
      click() {
        calendarRef.value.getApi().prevYear();
      },
    },
    // 下一年按钮
    nextYear: {
      text: t("schedule.next_year"),
      click() {
        calendarRef.value.getApi().nextYear();
      },
    },
    // 上一个月按钮
    prev: {
      text: t("schedule.previous_month"),
      click() {
        calendarRef.value.getApi().prev();
      },
    },
    // 下一个月按钮
    next: {
      text: t("schedule.next_month"),
      click() {
        calendarRef.value.getApi().next();
      },
    },
    // 今天按钮：跳转到今天
    today: {
      text: t("schedule.today"),
      click() {
        calendarRef.value.getApi().today();
      },
    },
    // 周视图按钮：切换到周视图
    week: {
      text: t("schedule.week"),
      click() {
        calendarRef.value.getApi().changeView("timeGridWeek");
      },
    },
  },

  // 日历事件列表（目前为空，后续可接入数据源）
  events: [],

  // 日期点击事件处理函数
  dateClick(info) {
    //console.log(info);

    if (firstClickOnDate) {
      firstClickOnDate = false;
      passing_date_characters(dataStartTemp, info.dateStr);
    } else {
      firstClickOnDate = true;
      dataStartTemp = info.dateStr;
    }
  },

  //选择日期
  select(info) {
    if (info.end - info.start > 86400000) {
      //选择了多日
      //console.log("选择了多日:", info);
      passing_date_characters_Select(info.startStr, info.endStr);
    } else {
      //选择单日
      //console.log("选择单日:", info);
    }
  },

  //事件event点击处理函数
  // eventClick(info) {
  //   const nowTime = new Date().getTime();
  //   const timeDifference = nowTime - lastEventClickTime.value;

  //   // 判断是否为双击（400ms 内连续点击）
  //   if (timeDifference < 400 && timeDifference > 0) {
  //     console.log("双击事件:", info);
  //     // 双击功能：快速添加事件
  //   } else {
  //     console.log("单击事件:", info);
  //     // 单击功能：显示日期详情
  //   }
  //   // 更新上次点击时间
  //   lastEventClickTime.value = nowTime;
  // },

  //event拖动处理
  eventDrop(info) {
    //console.log(info);
    passing_date_characters_Select(info.event.startStr, info.event.endStr);
  },
});

function switchShow() {

  if(props.isEditable){
    if (isShow.value) {
        isShow.value = false;
      } else {
        isShow.value = true;
      }
  }
  
}

function splicingDataWeek(data) {
  return data + " " + DateUtils.getI18nWeekday(data);
}

onMounted(() => {
  //console.log(eventData.value)
  calendarOptions.value.events.push(eventData.value);
});
</script>
<template>
  <div class="mb-4">
    <div
      @click="switchShow"
      class="flex items-center gap-2 border border-gray-200 rounded-xl px-3 py-1.5 shadow-sm bg-white mb-4 sticky top-0 z-50"
    >
      <!-- 日历图标 -->
      <div class="date-icon">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="icon icon-tabler icons-tabler-outline icon-tabler-calendar-week"
        >
          <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
          <path
            d="M4 7a2 2 0 0 1 2 -2h12a2 2 0 0 1 2 2v12a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2v-12z"
          ></path>
          <path d="M16 3v4"></path>
          <path d="M8 3v4"></path>
          <path d="M4 11h16"></path>
          <path d="M7 14h.013"></path>
          <path d="M10.01 14h.005"></path>
          <path d="M13.01 14h.005"></path>
          <path d="M16.015 14h.005"></path>
          <path d="M13.015 17h.005"></path>
          <path d="M7.01 17h.005"></path>
          <path d="M10.01 17h.005"></path>
        </svg>
      </div>

      <!-- 日期显示 -->
      <div class="date-display flex items-center justify-between gap-2 flex-1">
        <div class="start-date text-gray-700 font-medium">
          {{ splicingDataWeek(eventData.start) }}
        </div>
        <div class="text-gray-500">{{ t("schedule.to") }}</div>
        <div class="end-date text-gray-700 font-medium">
          {{
            splicingDataWeek(
              eventData.start === eventData.end
                ? eventData.end
                : DateUtils.toRealEnd(eventData.end),
            )
          }}
        </div>
      </div>
    </div>

    <FullCalendar v-if="isShow" ref="calendarRef" :options="calendarOptions" />
  </div>
</template>
