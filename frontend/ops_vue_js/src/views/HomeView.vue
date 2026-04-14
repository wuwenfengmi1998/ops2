<script setup>
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { useUsersStore } from '@/stores/users'
import { usePageTitle } from '@/composables/usePageTitle'
import { scheduleApi } from '@/api/schedule'
import { ref, computed, onMounted } from 'vue'

usePageTitle('appname.home')
const { t, locale } = useI18n()
const userStore = useUserStore()
const usersStore = useUsersStore()

// 今日日程数据
const todaySchedules = ref([])
const loadingSchedules = ref(false)

// 获取今日日期字符串
const todayStr = computed(() => {
  const today = new Date()
  return today.toISOString().split('T')[0]
})

// 格式化今日日期显示
const todayDisplay = computed(() => {
  const today = new Date()
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  const weekdaysEn = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
  const day = today.getDay()
  const dayNum = today.getDate()

  // 根据语言返回不同格式
  if (locale.value === 'en') {
    return `${dayNum} ${weekdaysEn[day]}`
  }
  return `${today.getMonth() + 1}月${dayNum}日 ${weekdays[day]}`
})

// 今日日程数量
const todayCount = computed(() => todaySchedules.value.length)

// 获取今日日程
async function fetchTodaySchedules() {
  loadingSchedules.value = true
  try {
    const { errCode, data } = await scheduleApi.getEvents({
      start: todayStr.value,
      end: todayStr.value
    })
    if (errCode === 0 && data?.list) {
      todaySchedules.value = data.list
    }
  } catch (e) {
    console.error('获取今日日程失败', e)
  } finally {
    loadingSchedules.value = false
  }
}

// 获取用户名
function getUsername(userId) {
  if (!userId) return ''
  return usersStore.getUsernameFromUserID(userId) || `用户${userId}`
}

// 格式化时间
function formatTime(dateStr) {
  const d = new Date(dateStr)
  return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
}

// 格式化日期
function formatDate(dateStr) {
  if (!dateStr) return ''
  return dateStr
}

// 格式化开始结束日期
function formatDateRange(startDate, endDate) {
  if (!startDate) return ''
  if (startDate === endDate) {
    return startDate
  }
  // 返回数组用于分行显示
  return [startDate, endDate]
}

// 获取星期几
function getWeekday(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const day = date.getDay()
  if (locale.value === 'en') {
    const weekdaysEn = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
    return weekdaysEn[day]
  }
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  return weekdays[day]
}

onMounted(() => {
  fetchTodaySchedules()
})
</script>

<template>
  <div class="mx-auto max-w-6xl px-6 py-6">
    <h2 class="mb-6 text-2xl font-bold text-gray-900 dark:text-white">{{ t('message.welcome') }}</h2>

    <!-- 日程卡片 -->
    <div class="mb-6 rounded-xl border border-gray-200 bg-white px-5 py-4 dark:border-dk-muted dark:bg-dk-card">
      <div class="mb-3 flex items-center justify-between">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('appname.schedule') }}</h3>
        <span class="text-sm text-gray-500 dark:text-gray-400">{{ t('home.today', { date: todayDisplay }) }}</span>
      </div>

      <!-- 加载状态 -->
      <div v-if="loadingSchedules" class="py-4 text-center text-gray-500">
        {{ t('home.loading') }}
      </div>

      <!-- 日程列表 -->
      <div v-else-if="todaySchedules.length > 0">
        <div class="mb-2 text-sm text-gray-600 dark:text-gray-400">
          {{ t('home.today_schedule_count', { count: todayCount }) }}
        </div>
        <ul class="space-y-2">
          <li
            v-for="schedule in todaySchedules"
            :key="schedule?.ID"
            class="flex items-start gap-3 rounded-lg bg-gray-50 px-3 py-2 dark:bg-dk-base"
          >
            <span
              class="flex flex-col whitespace-nowrap rounded px-2 py-0.5 text-sm font-medium text-white"
              :style="{ backgroundColor: schedule?.BgColor || '#999' }"
            >
              <span v-if="schedule?.StartDate !== schedule?.EndDate">
                <div>{{ schedule?.StartDate }} {{ getWeekday(schedule?.StartDate) }}</div>
                <div>{{ schedule?.EndDate }} {{ getWeekday(schedule?.EndDate) }}</div>
              </span>
              <span v-else>{{ schedule?.StartDate }} {{ getWeekday(schedule?.StartDate) }}</span>
            </span>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium text-gray-800 dark:text-gray-200">
                {{ schedule?.Title || '' }}
              </p>
              <p class="text-xs text-gray-500 dark:text-gray-400">
                {{ getUsername(schedule?.UserID) }}
              </p>
            </div>
          </li>
        </ul>
      </div>

      <!-- 无日程 -->
      <div v-else class="py-4 text-center text-gray-500 dark:text-gray-400">
        {{ t('home.today_no_schedule') }}
      </div>
    </div>

    <!-- 功能入口卡片 -->
    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      
      <div
        class="rounded-xl border border-gray-200 bg-white px-5 py-4 transition-shadow hover:shadow-md dark:border-dk-muted dark:bg-dk-card"
      >
        <p class="mb-1 text-sm text-gray-500">{{ t('appname.purchase') }}</p>
        <p class="text-lg font-bold text-gray-900 dark:text-white">—</p>
      </div>
    </div>
  </div>
</template>
