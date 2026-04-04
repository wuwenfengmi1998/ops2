<script setup>
// Vue 核心响应式 API
import { ref, watch, onMounted, onBeforeUnmount, reactive } from "vue";
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

// 模态框相关状态
const showModal = ref(false);
const eventData = ref({
  id: 0,
  title: "",
  startDate: "",
  endDate: "",
  color: "#066FD1", // 默认蓝色工作事件
  isEditing: false,   //是否处于编辑模式，false就是添加模式
  isEditable: false,//是否有权限编辑，无权限就让按钮灰掉
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

const pageData = reactive({//本页全局变量
  seleEventID: 0,//上次点击的eventid
  lastClickTime: 0,// 用于跟踪上次点击时间的响应式变量
  lastClickTimeStr: "",
  lastEventClickTime: 0,// 用于跟踪上次点击event时间的响应式变量
  lastEventClickID: 0,

  submitChecked: false,

  isCopy:false,
  copyTitle:"",
  copyColor:"",

})



function unseleEvent(eventID) {
  //寻找哪个event被单击了并修改边框
  const target = calendarOptions.value.events.find(item => item.id === eventID)
  if (target) {
    target.borderColor = target.backgroundColor;
  }
}

function unseleEventAll() {
  unseleEvent(pageData.seleEventID)
  pageData.seleEventID = 0;
}

function seleEvent(eventID) {
  //单击了event
  //取消上次选中
  if (pageData.seleEventID != 0) {
    unseleEvent(pageData.seleEventID);
  }
  //寻找哪个event被单击了并修改边框
  const target = calendarOptions.value.events.find(item => item.id === eventID)
  if (target) {
    target.borderColor = "#000000";
  }
  pageData.seleEventID = eventID;


}

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


  eventDidMount(info) {
    const titleEl = info.el.querySelector('.fc-event-title')
    if (titleEl && titleEl.scrollWidth > titleEl.clientWidth) {
      //titleEl.setAttribute('data-truncated', 'true')
      //日程过长 需要特殊处理
      //console.log("--",info)
    }
  },

  // 👇 加这个！日历渲染完成 / 切换年月都会触发
  datesSet(info) {
    calendarNowShow.value = info;
    //console.log(info);
  },

  // 日期点击事件处理函数
  dateClick(info) {
    const nowTime = new Date().getTime();
    const timeDifference = nowTime - pageData.lastClickTime;

    

    unseleEventAll();//点击了日期就取消event的选择

    //判断和上次点击的是不是同一天
    if (info.dateStr === pageData.lastClickTimeStr) {
      // 判断是否为双击（400ms 内连续点击）
      if (timeDifference < 400 && timeDifference > 0) {
        // 双击功能：快速添加事件

        //先判断是否登录
        if (userStore.isLoggedIn) {
          openEventModal(info.dateStr, info.dateStr);
        } else {
          toast.warning(t("message.login_to_your_account"));
          router.replace("/login?redirect=/schedule");
        }
      }
    }
    pageData.lastClickTimeStr = info.dateStr;


    // 更新上次点击时间
    pageData.lastClickTime = nowTime;
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
    const timeDifference = nowTime - pageData.lastEventClickTime;

    console.log(info)

    //判断event的title是否过长，如果是被截断的 就toast.info弹窗显示
    // const titleEl = info.el.querySelector('.fc-event-title');
    // if (titleEl && titleEl.scrollWidth > titleEl.clientWidth) {
    //   //title过长 得换一种方式显示
    //   //toast.info(info.event.title);
    // }

    // 单击功能：
    var eventid = parseInt(info.event.id)
    seleEvent(eventid);

    //判断和上次点击的是不是同一个event
    if (eventid === pageData.lastEventClickID) {
      // 判断是否为双击（400ms 内连续点击）
      if (timeDifference < 400 && timeDifference > 0) {
        //console.log("双击事件:", info);
        openEventModal(info.event.startStr, info.event.end ? info.event.endStr : info.event.startStr, parseInt(info.event.id), info.event.title, info.event.backgroundColor, true, info.event.durationEditable);
        // 双击功能：
        unseleEventAll()
      }
    }
    pageData.lastEventClickID = eventid;


    // 更新上次点击时间
    pageData.lastEventClickTime = nowTime;
  },

  //event拖动处理
  eventDrop(info) {
    updateEditData(parseInt(info.event.id),info.event.title,info.event.startStr,info.event.end === null ? info.event.startStr : DateUtils.toRealEnd(info.event.end),info.event.backgroundColor);
  },
});

function editSaveEvent(){
  //console.log(eventData)
  updateEditData(eventData.value.id,eventData.value.title,eventData.value.startDate,eventData.value.startDate===eventData.value.endDate?eventData.value.startDate:DateUtils.toRealEnd(eventData.value.endDate),eventData.value.color);
}

function updateEditData(id,title,start,end,color){
//需要发送到后端修改event
    //console.log(info.event)
    scheduleApi
      .editEvent({
        id: id,
        title: title,
        start: start,
        end: end,
        color: color,
      })
      .then((r) => {
        //console.log(r);
        if (r.errCode == 0) {
          //前端提交是否错误
          switch (
          r.raw.err_code //后端返回是否错误
          ) {
            case 0:
              closeEventModal();
              getEvents();//从新从后端获取最新数据
              break;
            default:
              toast.danger(t("message.server_error"));
              break;
          }
        }
      });
}

// 打开模态框
const openEventModal = (dateStr, dataEnd, id = 0, title = "", color = "#066FD1", isEditing = false, isEditable = true) => {
  eventData.value = {
    id: id,
    title: title,
    startDate: dateStr,
    endDate: dataEnd,
    color: color,
    isEditing: isEditing,
    isEditable: isEditable,
  };

  //console.log(eventData.value);

  showModal.value = true;
};

// 关闭模态框
const closeEventModal = () => {
  showModal.value = false;
};


// 保存日程事件
const saveEvent = () => {
  if (!eventData.value.title.trim()) {
    //alert("请输入日程内容");
    pageData.submitChecked = true;
    toast.warning(t("schedule.event_title_required"));
    return;
  }
  pageData.submitChecked = false;

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

  //console.log(newEvent)

  scheduleApi
    .addEvent({
      title: newEvent.title,
      start: newEvent.start,
      end: newEvent.end === newEvent.start ? newEvent.end : DateUtils.toRealEnd(newEvent.end),
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
      } else {
        toast.danger(t("message.server_error"));
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
      //console.log(r);
      if (r.errCode == 0) {
        //前端提交是否错误
        switch (
        r.raw.err_code //后端返回是否错误
        ) {
          case 0:
            calendarOptions.value.events = [];
            var events = r.raw.return.list;
            //console.log(events);
            var eventstemp = [];
            events?.forEach((item) => {

              calendarOptions.value.events.push({
                id: item.ID, // 后端 ID
                title: item.Title, // 标题
                start: item.StartDate, // 开始日期
                end: item.StartDate === item.EndDate ? item.EndDate : DateUtils.toCalendarEnd(item.EndDate), // 结束日期
                backgroundColor: item.BgColor, // 背景色
                borderColor: item.BgColor,      // 边框色（一般和背景一样）
                allDay: true, // 全天事件
                editable: item.edit,
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

//删除event
function delEvent(){
 
  scheduleApi
    .deleEvent({
      id:eventData.value.id
    }).then((r) => {
        //console.log(r);
        if (r.errCode == 0) {
          //前端提交是否错误
          switch (
          r.raw.err_code //后端返回是否错误
          ) {
            case 0:
              closeEventModal();
              getEvents();//从新从后端获取最新数据
              break;
            default:
              toast.danger(t("message.server_error"));
              break;
          }
        }else{
          toast.danger(t("message.server_error"));
        }
      });
    
}


// 颜色选择处理
const selectColor = (colorValue) => {
  if (eventData.value.isEditable) {
    eventData.value.color = colorValue;
  }

};

function copyEvent(){
    pageData.copyTitle=eventData.value.title;
    pageData.copyColor=eventData.value.color;
    pageData.isCopy=true;
    toast.info("已复制");
}

function pastEvent(){
  if (pageData.isCopy){
    if(eventData.value.isEditable){
      eventData.value.color=pageData.copyColor;
      eventData.value.title=pageData.copyTitle;
      toast.info("已粘贴");
    }else{
      toast.warning("这不是你的日程");
    }
  }

}

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
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-gray-800/20">
      <!-- 👇 最关键：给外层加最大高度 + 溢出隐藏 -->
      <div class="modal-content bg-white rounded-lg shadow-lg w-full max-w-2xl max-h-[95vh] flex flex-col">
        <!-- 模态框头部 -->
        <div class="modal-header border-b p-4 flex justify-between items-center flex-shrink-0">
          <h5 class="modal-title text-lg font-semibold">
            {{ userStore.isLoggedIn ? eventData.isEditing ? "修改日程" : t("schedule.add_event") : "查看日程" }}

          </h5>
          <h5 class="modal-title text-lg font-semibold absolute left-1/2 -translate-x-1/2">
            {{userStore.isLoggedIn ? eventData.isEditing ?"xxx的日程":"":""}}
            
          </h5>
          <button @click="closeEventModal" class="btn-close text-gray-500 hover:text-gray-700">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
              class="icon icon-tabler icons-tabler-outline icon-tabler-x">
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
            :isEditable="eventData.isEditable" />

          <!-- 内容输入区域 -->
          <div class="mb-4">
            <div class="uni-easyinput input relative">
              <div class="uni-easyinput__content is-input-border border border-gray-300 rounded-md bg-white relative"
                :class="{
                  'border-gray-300': eventData.title || !pageData.submitChecked,
                  'border-red-500': !eventData.title && pageData.submitChecked
                }">
                <input v-model="eventData.title" type="text" maxlength="140"
                  class="uni-easyinput__content-input w-full px-3 py-2 outline-none"
                  :placeholder="t('schedule.event_title_placeholder')" @keyup.enter="saveEvent"
                  :disabled="!eventData.isEditable" />
              </div>
            </div>
          </div>

          <!-- 颜色选择区域 -->
          <div class="mb-4">
            <div class="color_box grid grid-cols-3 gap-2">
              <div v-for="color in colorOptions" :key="color.value" class="color_box_item">
                <label class="uni-label-pointer form-colorinput flex items-center gap-2 cursor-pointer"
                  @click="selectColor(color.value)">
                  <div class="uni-radio-wrapper">
                    <div class="uni-radio-input flex items-center justify-center w-6 h-6 rounded-full transition-all"
                      :style="{
                        backgroundColor: color.value,
                        borderColor: color.value,
                      }">
                      <svg v-if="eventData.color === color.value" width="18" height="18" viewBox="0 0 32 32">
                        <path
                          d="M1.952 18.080q-0.32-0.352-0.416-0.88t0.128-0.976l0.16-0.352q0.224-0.416 0.64-0.528t0.8 0.176l6.496 4.704q0.384 0.288 0.912 0.272t0.88-0.336l17.312-14.272q0.352-0.288 0.848-0.256t0.848 0.352l-0.416-0.416q0.32 0.352 0.32 0.816t-0.32 0.816l-18.656 18.912q-0.32 0.352-0.8 0.352t-0.8-0.32l-7.936-8.064z"
                          fill="#ffffff"></path>
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
        <!-- 如果没登录直接不显示 -->
        <div v-if="userStore.isLoggedIn"
          class="modal-footer border-t p-4 flex justify-between items-center flex-shrink-0">

          <div class="flex gap-2">
            <button v-if="eventData.isEditing" @click="delEvent" class="btn px-4 py-2 text-white bg-red-500 hover:bg-red-600 rounded-md 
         disabled:bg-gray-400 disabled:cursor-not-allowed" :disabled="!eventData.isEditable">
              删除
            </button>
          </div>
          <div class="flex gap-2">
            <button class="btn px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md"
            @click="copyEvent"
            >
              {{ t("schedule.copy") }}
            </button>
            <button class="btn px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md disabled:cursor-not-allowed"
                :disabled="!pageData.isCopy"
                @click="pastEvent"
            >
              {{ t("schedule.paste") }}
            </button>
            <button 
              v-if="!eventData.isEditing"
              @click="saveEvent"
              class="btn btn-primary px-4 py-2 bg-cyan-600 text-white hover:bg-cyan-700 rounded-md disabled:bg-gray-400 disabled:cursor-not-allowed"
              :disabled="!eventData.isEditable">
              {{t("schedule.add_event_button") }}
            </button>
            <button 
              v-if="eventData.isEditing"
              @click="editSaveEvent"
              class="btn btn-primary px-4 py-2 bg-teal-600 text-white hover:bg-teal-700 rounded-md disabled:bg-gray-400 disabled:cursor-not-allowed"
              :disabled="!eventData.isEditable">
              修改日程
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 日历主体区域，带边框和阴影 -->
    <div class="flex-1 rounded-lg border border-gray-200 bg-white p-0.5 shadow dark:border-dk-muted dark:bg-dk-card">
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

<style scoped>
/* 基础：保持单行省略 */
/* :deep(.fc-daygrid-event .fc-event-title) {
  white-space: nowrap !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  display: block !important;
} */

/* 超长文字：自动滚动 */
:deep(.fc-daygrid-event .fc-event-title[data-truncated="true"]) {
  display: inline-block !important;
  width: auto !important;
  white-space: nowrap !important;
  overflow: hidden !important;
  padding-right: 10px !important;
  animation: textScroll 5s linear infinite !important;
}

@keyframes textScroll {
  0% {
    transform: translateX(0%);
  }

  100% {
    transform: translateX(-100%);
  }
}
</style>