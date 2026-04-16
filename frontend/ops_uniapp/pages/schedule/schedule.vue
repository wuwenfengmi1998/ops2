<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { scheduleApi } from '../../@api/schedule.js'
import { userStore } from '../../store/user.js'

// ───────────── 状态 ─────────────
const loading = ref(false)
const events = ref([])
const showModal = ref(false)

// 当前选中月份（YYYY-MM）
const today = new Date()
const curYear = ref(today.getFullYear())
const curMonth = ref(today.getMonth() + 1)  // 1~12

// 事件表单
const form = reactive({
  id: 0,
  title: '',
  startDate: '',
  endDate: '',
  color: '#066FD1',
  isEditing: false,
  isEditable: true,
})

// 颜色选项（对标 PC 端 scheduleView.vue）
const colorOptions = [
  { value: '#066FD1', label: '工作' },
  { value: '#09D119', label: '值班' },
  { value: '#FF00FF', label: '考试' },
  { value: '#FFFF00', label: '待命' },
  { value: '#D16C13', label: '私假' },
  { value: '#D10D21', label: '公假' },
]

// ───────────── 工具函数 ─────────────
function pad(n) { return String(n).padStart(2, '0') }

function todayStr() {
  return `${today.getFullYear()}-${pad(today.getMonth() + 1)}-${pad(today.getDate())}`
}

function monthStart() {
  return `${curYear.value}-${pad(curMonth.value)}-01`
}

function monthEnd() {
  // 下个月第 0 天 = 本月最后一天
  const last = new Date(curYear.value, curMonth.value, 0)
  return `${curYear.value}-${pad(curMonth.value)}-${pad(last.getDate())}`
}

function formatDate(str) {
  if (!str) return ''
  return str.slice(0, 10)
}

function isToday(dateStr) {
  return dateStr && dateStr.slice(0, 10) === todayStr()
}

// 判断 event 是否在当天日期范围内
function isTodayEvent(item) {
  const t = todayStr()
  const s = (item.StartDate || '').slice(0, 10)
  const e = (item.EndDate || item.StartDate || '').slice(0, 10)
  return s <= t && t <= e
}

function colorLabel(hex) {
  return colorOptions.find(c => c.value === hex)?.label || ''
}

// ───────────── API 调用 ─────────────
async function loadEvents() {
  loading.value = true
  try {
    const { errCode, data } = await scheduleApi.getEvents({
      start: monthStart(),
      end: monthEnd(),
    })
    if (errCode === 0 && data?.list) {
      events.value = data.list
    }
  } finally {
    loading.value = false
  }
}

// ───────────── 月份切换 ─────────────
function prevMonth() {
  if (curMonth.value === 1) { curYear.value--; curMonth.value = 12 }
  else curMonth.value--
  loadEvents()
}

function nextMonth() {
  if (curMonth.value === 12) { curYear.value++; curMonth.value = 1 }
  else curMonth.value++
  loadEvents()
}

function goToday() {
  curYear.value = today.getFullYear()
  curMonth.value = today.getMonth() + 1
  loadEvents()
}

// ───────────── 日程按日期分组 ─────────────
function groupByDate(list) {
  const map = {}
  list.forEach(item => {
    const key = (item.StartDate || '').slice(0, 10)
    if (!map[key]) map[key] = []
    map[key].push(item)
  })
  return Object.keys(map).sort().map(date => ({ date, items: map[date] }))
}

// ───────────── 模态框操作 ─────────────
function openAdd() {
  if (!userStore.isLoggedIn) {
    uni.showToast({ title: '请先登录', icon: 'none' })
    return
  }
  Object.assign(form, {
    id: 0,
    title: '',
    startDate: todayStr(),
    endDate: todayStr(),
    color: '#066FD1',
    isEditing: false,
    isEditable: true,
  })
  showModal.value = true
}

function openEdit(item) {
  Object.assign(form, {
    id: item.ID,
    title: item.Title,
    startDate: (item.StartDate || '').slice(0, 10),
    endDate: (item.EndDate || item.StartDate || '').slice(0, 10),
    color: item.BgColor || '#066FD1',
    isEditing: true,
    isEditable: item.edit !== false,
  })
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function submitForm() {
  if (!form.title.trim()) {
    uni.showToast({ title: '请输入日程内容', icon: 'none' })
    return
  }
  try {
    let res
    if (form.isEditing) {
      res = await scheduleApi.editEvent({
        id: form.id,
        title: form.title.trim(),
        start: form.startDate,
        end: form.startDate === form.endDate ? form.endDate : form.endDate,
        color: form.color,
      })
    } else {
      res = await scheduleApi.addEvent({
        title: form.title.trim(),
        start: form.startDate,
        end: form.startDate === form.endDate ? form.endDate : form.endDate,
        color: form.color,
      })
    }
    if (res.errCode === 0 && res.data?.err_code === 0 || res.raw?.err_code === 0) {
      uni.showToast({ title: form.isEditing ? '修改成功' : '添加成功', icon: 'success' })
      closeModal()
      loadEvents()
    } else {
      uni.showToast({ title: '操作失败', icon: 'none' })
    }
  } catch {
    uni.showToast({ title: '网络错误', icon: 'none' })
  }
}

async function deleteEvent() {
  if (!form.isEditable) return
  uni.showModal({
    title: '确认删除',
    content: '确定要删除这个日程吗？',
    success: async (res) => {
      if (!res.confirm) return
      try {
        const r = await scheduleApi.deleEvent({ id: form.id })
        if (r.raw?.err_code === 0) {
          uni.showToast({ title: '已删除', icon: 'success' })
          closeModal()
          loadEvents()
        }
      } catch {
        uni.showToast({ title: '删除失败', icon: 'none' })
      }
    }
  })
}

// ───────────── 生命周期 ─────────────
let timer = null
onMounted(() => {
  loadEvents()
  timer = setInterval(loadEvents, 30000)
})
onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <view class="page-wrap">
    <!-- 顶部导航 -->
    <view class="nav-bar">
      <view class="nav-left" @tap="prevMonth">
        <text class="nav-arrow">‹</text>
      </view>
      <view class="nav-center" @tap="goToday">
        <text class="nav-title">{{ curYear }} 年 {{ curMonth }} 月</text>
        <text class="nav-sub">点击回今天</text>
      </view>
      <view class="nav-right" @tap="nextMonth">
        <text class="nav-arrow">›</text>
      </view>
    </view>

    <!-- 加载中 -->
    <view v-if="loading" class="loading-box">
      <text class="loading-text">加载中…</text>
    </view>

    <!-- 空状态 -->
    <view v-else-if="!events.length" class="empty-box">
      <text class="empty-icon">📅</text>
      <text class="empty-text">本月暂无日程</text>
    </view>

    <!-- 日程列表（按日期分组） -->
    <scroll-view v-else scroll-y class="list-scroll">
      <view v-for="group in groupByDate(events)" :key="group.date" class="date-group">
        <!-- 日期标题 -->
        <view class="date-header" :class="{ 'date-today': isToday(group.date) }">
          <text class="date-text">{{ group.date }}</text>
          <text v-if="isToday(group.date)" class="today-badge">今天</text>
        </view>
        <!-- 该日期下的事件 -->
        <view
          v-for="item in group.items"
          :key="item.ID"
          class="event-card"
          @tap="openEdit(item)"
        >
          <view class="event-color-bar" :style="{ backgroundColor: item.BgColor || '#066FD1' }" />
          <view class="event-body">
            <text class="event-title">{{ item.Title }}</text>
            <text class="event-meta">
              {{ colorLabel(item.BgColor) }}
              <text v-if="item.StartDate !== item.EndDate">
                · {{ formatDate(item.StartDate) }} ~ {{ formatDate(item.EndDate) }}
              </text>
            </text>
          </view>
          <text v-if="item.edit !== false" class="event-edit-icon">›</text>
        </view>
      </view>
      <view style="height: 80px;" />
    </scroll-view>

    <!-- 悬浮添加按钮（已登录才显示） -->
    <view v-if="userStore.isLoggedIn" class="fab" @tap="openAdd">
      <text class="fab-icon">＋</text>
    </view>

    <!-- 弹出模态框 -->
    <view v-if="showModal" class="modal-mask" @tap.self="closeModal">
      <view class="modal-box">
        <view class="modal-header">
          <text class="modal-title">{{ form.isEditing ? '编辑日程' : '添加日程' }}</text>
          <text class="modal-close" @tap="closeModal">✕</text>
        </view>

        <view class="modal-body">
          <!-- 日程内容 -->
          <view class="form-item">
            <text class="form-label">日程内容</text>
            <input
              v-model="form.title"
              class="form-input"
              placeholder="请输入日程内容"
              maxlength="140"
              :disabled="!form.isEditable"
            />
          </view>

          <!-- 开始日期 -->
          <view class="form-item">
            <text class="form-label">开始日期</text>
            <picker mode="date" :value="form.startDate" @change="e => form.startDate = e.detail.value">
              <view class="form-picker">
                <text>{{ form.startDate || '请选择' }}</text>
                <text class="picker-arrow">›</text>
              </view>
            </picker>
          </view>

          <!-- 结束日期 -->
          <view class="form-item">
            <text class="form-label">结束日期</text>
            <picker mode="date" :value="form.endDate" @change="e => form.endDate = e.detail.value">
              <view class="form-picker">
                <text>{{ form.endDate || '请选择' }}</text>
                <text class="picker-arrow">›</text>
              </view>
            </picker>
          </view>

          <!-- 颜色选择 -->
          <view class="form-item">
            <text class="form-label">类型</text>
            <view class="color-grid">
              <view
                v-for="c in colorOptions"
                :key="c.value"
                class="color-item"
                :class="{ 'color-selected': form.color === c.value }"
                @tap="form.isEditable && (form.color = c.value)"
              >
                <view class="color-dot" :style="{ backgroundColor: c.value }" />
                <text class="color-label">{{ c.label }}</text>
              </view>
            </view>
          </view>
        </view>

        <view class="modal-footer">
          <button
            v-if="form.isEditing && form.isEditable"
            class="btn-danger"
            @tap="deleteEvent"
          >删除</button>
          <view class="footer-right">
            <button class="btn-cancel" @tap="closeModal">取消</button>
            <button
              v-if="form.isEditable"
              class="btn-primary"
              @tap="submitForm"
            >{{ form.isEditing ? '保存' : '添加' }}</button>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page-wrap {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f3f4f6;
}

/* 顶部导航 */
.nav-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  padding: 48rpx 32rpx 24rpx;
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.06);
}
.nav-arrow {
  font-size: 48rpx;
  color: #374151;
  padding: 0 24rpx;
}
.nav-center {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.nav-title {
  font-size: 34rpx;
  font-weight: 600;
  color: #111827;
}
.nav-sub {
  font-size: 22rpx;
  color: #9ca3af;
  margin-top: 4rpx;
}

/* 加载/空状态 */
.loading-box, .empty-box {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16rpx;
}
.loading-text { color: #9ca3af; font-size: 28rpx; }
.empty-icon { font-size: 80rpx; }
.empty-text { color: #9ca3af; font-size: 28rpx; }

/* 列表滚动区 */
.list-scroll {
  flex: 1;
  padding: 24rpx 24rpx 0;
}

/* 日期分组 */
.date-group { margin-bottom: 24rpx; }
.date-header {
  display: flex;
  align-items: center;
  gap: 12rpx;
  padding: 8rpx 0 12rpx;
}
.date-text { font-size: 26rpx; color: #6b7280; font-weight: 500; }
.date-today .date-text { color: #2563eb; }
.today-badge {
  background: #2563eb;
  color: #fff;
  font-size: 20rpx;
  padding: 2rpx 12rpx;
  border-radius: 20rpx;
}

/* 事件卡片 */
.event-card {
  display: flex;
  align-items: center;
  background: #fff;
  border-radius: 16rpx;
  margin-bottom: 12rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.06);
}
.event-color-bar {
  width: 8rpx;
  align-self: stretch;
  flex-shrink: 0;
}
.event-body {
  flex: 1;
  padding: 20rpx 20rpx 20rpx 16rpx;
  min-width: 0;
}
.event-title {
  font-size: 30rpx;
  font-weight: 500;
  color: #111827;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.event-meta {
  font-size: 24rpx;
  color: #9ca3af;
  margin-top: 6rpx;
  display: block;
}
.event-edit-icon {
  color: #d1d5db;
  font-size: 36rpx;
  padding-right: 20rpx;
}

/* FAB 悬浮按钮 */
.fab {
  position: fixed;
  right: 40rpx;
  bottom: 120rpx;
  width: 96rpx;
  height: 96rpx;
  border-radius: 48rpx;
  background: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 24rpx rgba(37,99,235,0.4);
}
.fab-icon { color: #fff; font-size: 52rpx; line-height: 1; }

/* 模态框 */
.modal-mask {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  z-index: 100;
  display: flex;
  align-items: flex-end;
}
.modal-box {
  width: 100%;
  background: #fff;
  border-radius: 32rpx 32rpx 0 0;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32rpx 32rpx 24rpx;
  border-bottom: 1rpx solid #f3f4f6;
}
.modal-title { font-size: 32rpx; font-weight: 600; color: #111827; }
.modal-close { font-size: 36rpx; color: #9ca3af; padding: 8rpx; }

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 24rpx 32rpx;
}

/* 表单 */
.form-item { margin-bottom: 28rpx; }
.form-label {
  display: block;
  font-size: 26rpx;
  color: #6b7280;
  margin-bottom: 12rpx;
}
.form-input {
  width: 100%;
  border: 2rpx solid #e5e7eb;
  border-radius: 12rpx;
  padding: 16rpx 20rpx;
  font-size: 30rpx;
  color: #111827;
  background: #fff;
  box-sizing: border-box;
}
.form-picker {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border: 2rpx solid #e5e7eb;
  border-radius: 12rpx;
  padding: 16rpx 20rpx;
  font-size: 30rpx;
  color: #111827;
  background: #fff;
}
.picker-arrow { color: #d1d5db; font-size: 32rpx; }

/* 颜色网格 */
.color-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16rpx;
}
.color-item {
  display: flex;
  align-items: center;
  gap: 12rpx;
  padding: 14rpx 16rpx;
  border: 2rpx solid #e5e7eb;
  border-radius: 12rpx;
  background: #fff;
}
.color-selected {
  border-color: #2563eb;
  background: #eff6ff;
}
.color-dot {
  width: 28rpx;
  height: 28rpx;
  border-radius: 50%;
  flex-shrink: 0;
}
.color-label { font-size: 24rpx; color: #374151; }

/* 底部按钮 */
.modal-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 32rpx 48rpx;
  border-top: 1rpx solid #f3f4f6;
}
.footer-right {
  display: flex;
  gap: 16rpx;
  margin-left: auto;
}
.btn-danger {
  background: #fee2e2;
  color: #dc2626;
  border: none;
  border-radius: 12rpx;
  padding: 16rpx 32rpx;
  font-size: 28rpx;
}
.btn-cancel {
  background: #f3f4f6;
  color: #374151;
  border: none;
  border-radius: 12rpx;
  padding: 16rpx 32rpx;
  font-size: 28rpx;
}
.btn-primary {
  background: #2563eb;
  color: #fff;
  border: none;
  border-radius: 12rpx;
  padding: 16rpx 32rpx;
  font-size: 28rpx;
}
</style>
