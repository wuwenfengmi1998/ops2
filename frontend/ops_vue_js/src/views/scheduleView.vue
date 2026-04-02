<script setup>
// Vue 核心响应式 API
import { ref, watch, onMounted, onBeforeUnmount } from "vue";
// FullCalendar Vue 3 组件
import FullCalendar from "@fullcalendar/vue3";
// FullCalendar 插件：月视图
import dayGridPlugin from "@fullcalendar/daygrid";
// FullCalendar 插件：周视图（时间网格）
import timeGridPlugin from "@fullcalendar/timegrid";
// FullCalendar 插件：交互功能（拖拽、点击等）
import interactionPlugin from "@fullcalendar/interaction";
// FullCalendar 插件：列表视图
import listPlugin from "@fullcalendar/list";
// 国际化 hook
import { useI18n } from "vue-i18n";
// 页面标题 composable
import { usePageTitle } from "@/composables/usePageTitle";

import DatatimePickerForFullCalendar from "@/components/datatimePickerForFullCalendar.vue";

import { useToastStore } from "@/stores/toast";

// 用户状态管理
import { useUserStore } from "@/stores/user";

import { useRouter } from "vue-router";

import { scheduleApi } from "@/api/schedule";

import { useDateUtils } from "@/composables/useDateUtils";

const DateUtils = useDateUtils();

const router = useRouter();

// 获取用户 store 实例，用于访问和更新用户信息
const userStore = useUserStore();

const toast = useToastStore();

// 设置页面标题
usePageTitle("appname.schedule");

// 获取国际化翻译函数和当前语言
const { t, locale } = useI18n();

// FullCalendar 组件的引用，用于调用日历 API
const calendarRef = ref(null);
// 当前视图的年份
const calendarNowShow = ref();

// 用于跟踪上次点击时间的响应式变量
const lastClickTime = ref(0);
// 用于跟踪上次点击event时间的响应式变量
const lastEventClickTime = ref(0);

// 模态框相关状态
const showModal = ref(false);
const modalTitle = ref("添加日程");
const eventData = ref({
  title: "",
  startDate: "",
  endDate: "",
  color: "#066FD1", // 默认蓝色工作事件
});

// 颜色选项
const colorOptions = ref([
  { value: "#066FD1", label: t("schedule.work"), name: t("schedule.work") },
  { value: "#09D119", label: t("schedule.duty"), name: t("schedule.duty") },
  { value: "#FF00FF", label: t("schedule.exam"), name: t("schedule.exam") },
  {
    value: "#FFFF00",
    label: t("schedule.standby"),
    name: t("schedule.standby"),
  },
  {
    value: "#D16C13",
    label: t("schedule.personal_holiday"),
    name: t("schedule.personal_holiday"),
  },
  {
    value: "#D10D21",
    label: t("schedule.public_holiday"),
    name: t("schedule.public_holiday"),
  },
]);

// 日历配置选项
const calendarOptions = ref({
  // 日历高度：占满可用空间
  height: "100%",
  // 内容高度自适应
  contentHeight: "auto",
  // 使用当前应用语言
  locale: locale.value,
  // 注册使用的插件
  plugins: [dayGridPlugin, timeGridPlugin, interactionPlugin, listPlugin],
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
  expandRows: true,
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
    left: "prevYear,prev,today,next,nextYear",
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
        getEvents();
      },
    },
    // 下一年按钮
    nextYear: {
      text: t("schedule.next_year"),
      click() {
        calendarRef.value.getApi().nextYear();
        getEvents();
      },
    },
    // 上一个月按钮
    prev: {
      text: t("schedule.previous_month"),
      click() {
        calendarRef.value.getApi().prev();
        getEvents();
      },
    },
    // 下一个月按钮
    next: {
      text: t("schedule.next_month"),
      click() {
        calendarRef.value.getApi().next();
        getEvents();
      },
    },
    // 今天按钮：跳转到今天
    today: {
      text: t("schedule.today"),
      click() {
        calendarRef.value.getApi().today();
        getEvents();
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

  // 👇 加这个！日历渲染完成 / 切换年月都会触发
  datesSet(info) {
    calendarNowShow.value = info;
    //console.log(info);
  },

  // 日期点击事件处理函数
  dateClick(info) {
    const nowTime = new Date().getTime();
    const timeDifference = nowTime - lastClickTime.value;

    // 判断是否为双击（400ms 内连续点击）
    if (timeDifference < 400 && timeDifference > 0) {
      console.log("双击日期:", info.dateStr);
      // 双击功能：快速添加事件
      handleDoubleClick(info);
    } else {
      console.log("单击日期:", info.dateStr);
      // 单击功能：
      handleSingleClick(info);
    }

    // 更新上次点击时间
    lastClickTime.value = nowTime;
  },

  //选择日期
  select(info) {
    if (info.end - info.start > 86400000) {
      //选择了多日
      console.log("选择了多日:", info);
      openEventModal(info.startStr, info.endStr);
    } else {
      //选择单日 无功能
      //console.log("选择单日:", info);
    }
  },

  //事件event点击处理函数
  eventClick(info) {
    const nowTime = new Date().getTime();
    const timeDifference = nowTime - lastEventClickTime.value;

    // 判断是否为双击（400ms 内连续点击）
    if (timeDifference < 400 && timeDifference > 0) {
      console.log("双击事件:", info);
      // 双击功能：快速添加事件
    } else {
      console.log("单击事件:", info);
      // 单击功能：显示日期详情
    }
    // 更新上次点击时间
    lastEventClickTime.value = nowTime;
  },

  //event拖动处理
  eventDrop(info) {},
});

// 打开模态框
const openEventModal = (dateStr, dataEnd) => {
  eventData.value = {
    title: "",
    startDate: dateStr,
    endDate: dataEnd,
    color: "#066FD1",
  };

  showModal.value = true;
};

// 关闭模态框
const closeEventModal = () => {
  showModal.value = false;
};

// 处理双击事件：打开模态框添加事件
const handleDoubleClick = (info) => {
  //先判断是否登录
  if (userStore.isLoggedIn) {
    openEventModal(info.dateStr, info.dateStr);
  } else {
    toast.warning(t("message.login_to_your_account"));
    router.replace("/login?redirect=/schedule");
  }
};

// 处理单机事件：显示日期详情
const handleSingleClick = (info) => {
  const dateEvents = calendarOptions.value.events.filter(
    (event) =>
      event.start === info.dateStr ||
      (event.start <= info.dateStr && event.end > info.dateStr),
  );

  if (dateEvents.length > 0) {
    //alert(`${info.dateStr} 有 ${dateEvents.length} 个事件`)
  } else {
    //alert(`${info.dateStr} 没有事件`)
  }
};

// 保存日程事件
const saveEvent = () => {
  if (!eventData.value.title.trim()) {
    //alert("请输入日程内容");
    toast.warning(t("schedule.event_title_required"));
    return;
  }

  if (!eventData.value.startDate || !eventData.value.endDate) {
    //alert("请选择日期");
    toast.warning(t("schedule.date_required"));
    return;
  }

  const selectedColor = colorOptions.value.find(
    (color) => color.value === eventData.value.color,
  );
  const colorName = selectedColor ? selectedColor.name : eventData.value.color;

  const newEvent = {
    title: eventData.value.title.trim(),
    start: eventData.value.startDate,
    end: eventData.value.endDate,
    allDay: true,
    backgroundColor: eventData.value.color,
    borderColor: eventData.value.color,
    textColor: "#ffffff",
    extendedProps: {
      type: colorName,
      description: eventData.value.title.trim(),
    },
  };

  // 添加到日历事件列表
  //提交到后端

  scheduleApi
    .addEvent({
      title: newEvent.title,
      start: newEvent.start,
      end: DateUtils.toRealEnd(newEvent.end),
      color: newEvent.backgroundColor,
    })
    .then((r) => {
      //console.log(r);
      if (r.errCode == 0) {
        //前端提交是否错误
        switch (
          r.raw.err_code //后端返回是否错误
        ) {
          case 0:
            //calendarOptions.value.events.push(newEvent);
            toast.success(t("schedule.event_added_successfully"));
            // 关闭模态框
            closeEventModal();
            getEvents();
            break;
          default:
            toast.danger(t("message.server_error"));
            break;
        }
      }
    });
};

//从后端获取events
const getEvents = () => {
  //console.log(calendarNowShow.value)
  scheduleApi
    .getEvents({
      start: DateUtils.dateToStr(calendarNowShow.value.start),
      end: DateUtils.toRealEnd(calendarNowShow.value.end),
    })
    .then((r) => {
      console.log(r);
      if (r.errCode == 0) {
        //前端提交是否错误
        switch (
          r.raw.err_code //后端返回是否错误
        ) {
          case 0:
            calendarOptions.value.events=[];
            var events = r.raw.return.list;
            console.log(events);
            var eventstemp = [];
            events.forEach((item) => {
              
              calendarOptions.value.events.push({
                id: item.ID, // 后端 ID
                title: item.Title, // 标题
                start: item.StartDate, // 开始日期
                end: DateUtils.toCalendarEnd(item.EndDate), // 结束日期
                backgroundColor: item.BgColor, // 背景色
                allDay: true, // 全天事件
              });
            });

            break;
          default:
            toast.danger(t("message.server_error"));
            break;
        }
      }
    });
};
// 清除日期选择
const clearDates = () => {
  eventData.value.startDate = "";
  eventData.value.endDate = "";
};

// 颜色选择处理
const selectColor = (colorValue) => {
  eventData.value.color = colorValue;
};

// 监听语言变化，更新日历的本地化和按钮文字
watch(locale, () => {
  // 更新日历语言
  calendarOptions.value.locale = locale.value;
  // 更新自定义按钮的文字
  calendarOptions.value.customButtons.prevYear.text = t(
    "schedule.previous_year",
  );
  calendarOptions.value.customButtons.nextYear.text = t("schedule.next_year");
  calendarOptions.value.customButtons.prev.text = t("schedule.previous_month");
  calendarOptions.value.customButtons.next.text = t("schedule.next_month");
  calendarOptions.value.customButtons.today.text = t("schedule.today");
  calendarOptions.value.customButtons.week.text = t("schedule.week");

  colorOptions.value = [
    { value: "#066FD1", label: t("schedule.work"), name: t("schedule.work") },
    { value: "#09D119", label: t("schedule.duty"), name: t("schedule.duty") },
    { value: "#FF00FF", label: t("schedule.exam"), name: t("schedule.exam") },
    {
      value: "#FFFF00",
      label: t("schedule.standby"),
      name: t("schedule.standby"),
    },
    {
      value: "#D16C13",
      label: t("schedule.personal_holiday"),
      name: t("schedule.personal_holiday"),
    },
    {
      value: "#D10D21",
      label: t("schedule.public_holiday"),
      name: t("schedule.public_holiday"),
    },
  ];
});

onMounted(() => {
  getEvents();
  // const handleKeydown = (event) => {
  //   // Ctrl+C 事件
  //   if (event.ctrlKey && event.key === "c") {
  //     event.preventDefault(); // 可选：阻止默认复制行为
  //     console.log("Ctrl+C 被按下");
  //     // 你的业务逻辑
  //   }
  //   // Ctrl+V 事件
  //   if (event.ctrlKey && event.key === "v") {
  //     event.preventDefault(); // 可选：阻止默认粘贴行为
  //     console.log("Ctrl+V 被按下");
  //     // 你的业务逻辑
  //   }
  // };
  // document.addEventListener("keydown", handleKeydown);
  // // 清理事件监听器
  // onBeforeUnmount(() => {
  //   document.removeEventListener("keydown", handleKeydown);
  // });
});
</script>

<template>
  <!-- {{userStore.userCookie.Value}} -->
  <!-- 日历容器：占满视口高度减去顶部导航高度 -->
  <div class="flex w-full flex-col relative">
    <!-- 事件编辑模态框 -->
    <div
      v-if="showModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-gray-800/20"
    >
      <!-- 👇 最关键：给外层加最大高度 + 溢出隐藏 -->
      <div
        class="modal-content bg-white rounded-lg shadow-lg w-full max-w-2xl max-h-[95vh] flex flex-col"
      >
        <!-- 模态框头部 -->
        <div
          class="modal-header border-b p-4 flex justify-between items-center flex-shrink-0"
        >
          <h5 class="modal-title text-lg font-semibold">
            {{ t("schedule.add_event") }}
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

        <!-- 👇 主体区域：允许内部滚动 -->
        <div class="modal-body p-4 flex-1 overflow-y-auto">
          <!-- 日期选择区域 -->
          <DatatimePickerForFullCalendar
            v-model:startDate="eventData.startDate"
            v-model:endDate="eventData.endDate"
            :color="eventData.color"
            :title="eventData.title"
          />

          <!-- 内容输入区域 -->
          <div class="mb-4">
            <div class="uni-easyinput input relative">
              <div
                class="uni-easyinput__content is-input-border border border-gray-300 rounded-md bg-white relative"
              >
                <input
                  v-model="eventData.title"
                  type="text"
                  maxlength="140"
                  class="uni-easyinput__content-input w-full px-3 py-2 outline-none"
                  :placeholder="t('schedule.event_title_placeholder')"
                  @keyup.enter="saveEvent"
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
        </div>

        <!-- 👇 底部固定：不被滚动、不压缩 -->
        <div
          class="modal-footer border-t p-4 flex justify-between items-center flex-shrink-0"
        >
          <button
            @click="closeEventModal"
            class="btn px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md"
          >
            {{ t("schedule.close") }}
          </button>
          <div class="flex gap-2">
            <button
              class="btn px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md"
              disabled
            >
              {{ t("schedule.copy") }}
            </button>
            <button
              class="btn px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md"
              disabled
            >
              {{ t("schedule.paste") }}
            </button>
            <button
              @click="saveEvent"
              class="btn btn-primary px-4 py-2 bg-blue-600 text-white hover:bg-blue-700 rounded-md"
            >
              {{ t("schedule.add_event_button") }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 日历主体区域，带边框和阴影 -->
    <div
      class="flex-1 rounded-lg border border-gray-200 bg-white p-0.5 shadow dark:border-dk-muted dark:bg-dk-card"
    >
      <!-- 内层容器：隐藏溢出内容 -->
      <div class="h-full w-full overflow-hidden rounded-md">
        <!-- FullCalendar 日历组件 -->
        <!-- ref="calendarRef" 用于获取组件实例调用日历 API -->
        <!-- :options 绑定日历配置对象 -->
        <FullCalendar ref="calendarRef" :options="calendarOptions" />
      </div>
    </div>
  </div>
</template>
