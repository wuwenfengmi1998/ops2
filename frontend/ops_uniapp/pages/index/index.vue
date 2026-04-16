<script>
/**
 * 首页
 * 对标 PC 前端 src/views/HomeView.vue
 * 展示：今日日程 + 待处理采购订单数量
 */
import { userStore } from '../../store/user.js'
import { scheduleApi } from '../../@api/schedule.js'
import { purchaseApi } from '../../@api/purchase.js'

// 颜色类型映射（与 PC 端保持一致）
const COLOR_MAP = {
  '#066FD1': '工作',
  '#09D119': '值班',
  '#FF00FF': '考试',
  '#FFFF00': '待机',
  '#D16C13': '个人假期',
  '#D10D21': '公共假期',
}

const WEEKDAYS = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']

export default {
  data() {
    return {
      // 用户信息
      isLoggedIn: false,
      displayName: '',
      avatarUrl: '/static/ava.svg',

      // 今日日程
      todaySchedules: [],
      loadingSchedules: false,

      // 采购订单
      pendingOrderCount: 0,
      loadingOrders: false,

      // 今日显示
      todayDisplay: '',
      todayStr: '',
    }
  },

  methods: {
    // ── 工具方法 ──

    getTodayStr() {
      return new Date().toISOString().split('T')[0]
    },

    getTodayDisplay() {
      const today = new Date()
      const m = today.getMonth() + 1
      const d = today.getDate()
      const w = WEEKDAYS[today.getDay()]
      return `${m}月${d}日 ${w}`
    },

    getColorLabel(color) {
      return COLOR_MAP[color] || ''
    },

    formatDateRange(startDate, endDate) {
      if (!startDate) return ''
      if (startDate === endDate) return startDate
      return `${startDate} ~ ${endDate}`
    },

    getWeekday(dateStr) {
      if (!dateStr) return ''
      return WEEKDAYS[new Date(dateStr).getDay()]
    },

    // ── 数据加载 ──

    async fetchTodaySchedules() {
      this.loadingSchedules = true
      try {
        const today = this.getTodayStr()
        const { errCode, data } = await scheduleApi.getEvents({
          start: today,
          end: today,
        })
        if (errCode === 0 && data?.list) {
          this.todaySchedules = data.list
        }
      } catch {
        // 静默处理
      } finally {
        this.loadingSchedules = false
      }
    },

    async fetchPendingOrders() {
      if (!userStore.isLoggedIn) return
      this.loadingOrders = true
      try {
        const { errCode, data } = await purchaseApi.getOrderCount()
        if (errCode === 0 && data) {
          this.pendingOrderCount = data.pending || 0
        }
      } catch {
        // 静默处理
      } finally {
        this.loadingOrders = false
      }
    },

    refreshUserState() {
      this.isLoggedIn = userStore.isLoggedIn
      this.displayName = userStore.getDisplayName()
      this.avatarUrl = userStore.getAvatarUrl()
    },

    // ── 导航 ──

    goToSchedule() {
      uni.navigateTo({ url: '/pages/schedule/schedule' })
    },

    goToPurchase() {
      uni.navigateTo({ url: '/pages/purchase/list' })
    },

    goToLogin() {
      uni.navigateTo({ url: '/pages/signin' })
    },

    goToSettings() {
      uni.navigateTo({ url: '/pages/setting/my_info' })
    },
  },

  onShow() {
    // 每次回到首页都刷新用户状态（登录/登出后回来）
    this.refreshUserState()
  },

  onLoad() {
    this.todayStr = this.getTodayStr()
    this.todayDisplay = this.getTodayDisplay()
    this.refreshUserState()
    this.fetchTodaySchedules()
    this.fetchPendingOrders()
  },
}
</script>

<template>
  <view class="home-page">

    <!-- ── 顶部导航栏 ── -->
    <view class="navbar">
      <view class="navbar-left">
        <image src="/static/logo.svg" class="nav-logo" mode="aspectFit" />
        <text class="nav-title">Operations</text>
      </view>
      <view class="navbar-right">
        <view v-if="isLoggedIn" class="user-info" @tap="goToSettings">
          <image :src="avatarUrl" class="avatar" mode="aspectFill" />
          <text class="username">{{ displayName }}</text>
        </view>
        <view v-else class="login-btn" @tap="goToLogin">
          <text class="login-btn-text">登录</text>
        </view>
      </view>
    </view>

    <!-- ── 页面主体 ── -->
    <scroll-view scroll-y class="content">

      <!-- 欢迎语 -->
      <view class="welcome-row">
        <text class="welcome-text">
          {{ isLoggedIn ? `你好，${displayName}` : '欢迎使用' }}
        </text>
        <text class="date-text">{{ todayDisplay }}</text>
      </view>

      <!-- ── 今日日程卡片 ── -->
      <view class="card" @tap="goToSchedule">
        <view class="card-header">
          <view class="card-header-left">
            <view class="card-icon schedule-icon">
              <text class="icon-emoji">📅</text>
            </view>
            <text class="card-title">日程</text>
          </view>
          <text class="card-count">今日 {{ todaySchedules.length }} 项</text>
        </view>

        <!-- 加载中 -->
        <view v-if="loadingSchedules" class="empty-hint">
          <text>加载中…</text>
        </view>

        <!-- 有日程 -->
        <view v-else-if="todaySchedules.length > 0" class="schedule-list">
          <view
            v-for="item in todaySchedules"
            :key="item.ID"
            class="schedule-item"
          >
            <!-- 颜色标签 -->
            <view
              class="schedule-badge"
              :style="{ backgroundColor: item.BgColor || '#999' }"
            >
              <text class="badge-text">{{ getColorLabel(item.BgColor) }}</text>
            </view>
            <!-- 内容 -->
            <view class="schedule-info">
              <text class="schedule-title">{{ item.Title }}</text>
              <text class="schedule-date">
                {{ formatDateRange(item.StartDate, item.EndDate) }}
              </text>
            </view>
          </view>
        </view>

        <!-- 无日程 -->
        <view v-else class="empty-hint">
          <text>今日暂无日程</text>
        </view>

        <view class="card-arrow">
          <text class="arrow-text">查看全部 ›</text>
        </view>
      </view>

      <!-- ── 采购订单卡片（仅登录后显示待处理数） ── -->
      <view class="card" @tap="goToPurchase">
        <view class="card-header">
          <view class="card-header-left">
            <view class="card-icon purchase-icon">
              <text class="icon-emoji">🛒</text>
            </view>
            <text class="card-title">采购订单</text>
          </view>
        </view>

        <view v-if="isLoggedIn" class="stat-row">
          <view class="stat-item">
            <text class="stat-num" :class="{ 'num-warning': pendingOrderCount > 0 }">
              {{ loadingOrders ? '…' : pendingOrderCount || '—' }}
            </text>
            <text class="stat-label">待处理</text>
          </view>
        </view>
        <view v-else class="empty-hint">
          <text>登录后查看</text>
        </view>

        <view class="card-arrow">
          <text class="arrow-text">查看全部 ›</text>
        </view>
      </view>

    </scroll-view>
  </view>
</template>

<style scoped>
/* ── 页面容器 ── */
.home-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f3f4f6;
}

/* ── 顶部导航 ── */
.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 88rpx 32rpx 24rpx;
  background: #fff;
  border-bottom: 1rpx solid #e5e7eb;
}
.navbar-left {
  display: flex;
  align-items: center;
  gap: 16rpx;
}
.nav-logo {
  width: 64rpx;
  height: 64rpx;
  border-radius: 12rpx;
}
.nav-title {
  font-size: 36rpx;
  font-weight: 700;
  color: #1f2937;
}
.navbar-right {
  display: flex;
  align-items: center;
}
.user-info {
  display: flex;
  align-items: center;
  gap: 12rpx;
}
.avatar {
  width: 68rpx;
  height: 68rpx;
  border-radius: 50%;
  background: #e5e7eb;
}
.username {
  font-size: 28rpx;
  color: #374151;
  max-width: 160rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.login-btn {
  background: #2563eb;
  border-radius: 40rpx;
  padding: 12rpx 32rpx;
}
.login-btn-text {
  font-size: 28rpx;
  color: #fff;
  font-weight: 600;
}

/* ── 内容区 ── */
.content {
  flex: 1;
  padding: 32rpx 32rpx 40rpx;
}

/* ── 欢迎语 ── */
.welcome-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: 32rpx;
}
.welcome-text {
  font-size: 40rpx;
  font-weight: 700;
  color: #111827;
}
.date-text {
  font-size: 26rpx;
  color: #6b7280;
}

/* ── 通用卡片 ── */
.card {
  background: #fff;
  border-radius: 24rpx;
  padding: 36rpx;
  margin-bottom: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);
}

/* ── 卡片头部 ── */
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 28rpx;
}
.card-header-left {
  display: flex;
  align-items: center;
  gap: 16rpx;
}
.card-icon {
  width: 72rpx;
  height: 72rpx;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}
.schedule-icon { background: #eff6ff; }
.purchase-icon { background: #fefce8; }
.icon-emoji {
  font-size: 36rpx;
}
.card-title {
  font-size: 34rpx;
  font-weight: 600;
  color: #111827;
}
.card-count {
  font-size: 26rpx;
  color: #6b7280;
}

/* ── 日程列表 ── */
.schedule-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
  margin-bottom: 20rpx;
}
.schedule-item {
  display: flex;
  align-items: flex-start;
  gap: 20rpx;
  background: #f9fafb;
  border-radius: 16rpx;
  padding: 20rpx 24rpx;
}
.schedule-badge {
  border-radius: 8rpx;
  padding: 6rpx 14rpx;
  flex-shrink: 0;
}
.badge-text {
  font-size: 22rpx;
  color: #fff;
  font-weight: 600;
}
.schedule-info {
  flex: 1;
}
.schedule-title {
  font-size: 30rpx;
  font-weight: 500;
  color: #1f2937;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.schedule-date {
  font-size: 24rpx;
  color: #9ca3af;
  display: block;
  margin-top: 6rpx;
}

/* ── 统计行 ── */
.stat-row {
  display: flex;
  gap: 40rpx;
  margin-bottom: 20rpx;
}
.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.stat-num {
  font-size: 56rpx;
  font-weight: 700;
  color: #111827;
  line-height: 1;
}
.num-warning {
  color: #d97706;
}
.stat-label {
  font-size: 24rpx;
  color: #9ca3af;
  margin-top: 8rpx;
}

/* ── 空提示 & 箭头 ── */
.empty-hint {
  padding: 24rpx 0;
  text-align: center;
  font-size: 28rpx;
  color: #9ca3af;
}
.card-arrow {
  display: flex;
  justify-content: flex-end;
  margin-top: 8rpx;
}
.arrow-text {
  font-size: 26rpx;
  color: #2563eb;
}
</style>
