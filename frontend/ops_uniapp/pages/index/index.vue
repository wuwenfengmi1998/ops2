<template>
  <view class="index-container">
    <!-- 顶部标题 -->
    <view class="header">
      <text class="title">{{ t('index.welcome') }}</text>
      <text class="subtitle">{{ t('index.subtitle') }}</text>
    </view>

    <!-- 订单统计卡片 -->
    <view class="stats-grid">
      <view class="stat-card" @click="goToOrders('pending')">
        <text class="stat-num">{{ stats.pending }}</text>
        <text class="stat-label">{{ t('index.pending') }}</text>
      </view>
      <view class="stat-card" @click="goToOrders('ordered')">
        <text class="stat-num">{{ stats.ordered }}</text>
        <text class="stat-label">{{ t('index.ordered') }}</text>
      </view>
      <view class="stat-card" @click="goToOrders('arrived')">
        <text class="stat-num">{{ stats.arrived }}</text>
        <text class="stat-label">{{ t('index.arrived') }}</text>
      </view>
      <view class="stat-card" @click="goToOrders('received')">
        <text class="stat-num">{{ stats.received }}</text>
        <text class="stat-label">{{ t('index.received') }}</text>
      </view>
    </view>

    <!-- 快捷操作 -->
    <view class="section">
      <text class="section-title">{{ t('index.quickActions') }}</text>
      <view class="action-list">
        <view class="action-item" @click="goToOrders('all')">
          <text class="action-icon">📋</text>
          <text class="action-text">{{ t('index.allOrders') }}</text>
          <text class="action-arrow">→</text>
        </view>
        <view class="action-item" @click="goToAddOrder">
          <text class="action-icon">➕</text>
          <text class="action-text">{{ t('index.addOrder') }}</text>
          <text class="action-arrow">→</text>
        </view>
        <view class="action-item" @click="goToSchedule">
          <text class="action-icon">📅</text>
          <text class="action-text">{{ t('index.schedule') }}</text>
          <text class="action-arrow">→</text>
        </view>
      </view>
    </view>

    <!-- 语言切换入口 -->
    <view class="section">
      <view class="lang-card" @click="switchLang">
        <text class="lang-icon">🌐</text>
        <view class="lang-info">
          <text class="lang-title">{{ t('index.language') }}</text>
          <text class="lang-value">{{ locale === 'zh' ? '中文' : 'English' }}</text>
        </view>
        <text class="action-arrow">→</text>
      </view>
    </view>

    <!-- 退出登录 -->
    <view class="logout-section">
      <button class="logout-btn" @click="handleLogout">
        {{ t('index.logout') }}
      </button>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { setLocale, getCurrentLocale } from '../../locales/index.js'

const { t, locale } = useI18n()

// 统计数据
const stats = reactive({
  pending: 0,
  ordered: 0,
  arrived: 0,
  received: 0
})

// 获取订单统计
const fetchStats = () => {
  uni.request({
    url: getApp().globalData.BASE_URL + '/purchase/getordercount',
    method: 'POST',
    header: {
      'Content-Type': 'application/json',
      'Cookie': uni.getStorageSync('sessionCookie') || ''
    },
    success: (res) => {
      if (res.data.code === 0 && res.data.data) {
        stats.pending = res.data.data.pending || 0
        stats.ordered = res.data.data.ordered || 0
        stats.arrived = res.data.data.arrived || 0
        stats.received = res.data.data.received || 0
      }
    }
  })
}

// 跳转订单列表
const goToOrders = (status) => {
  const url = status === 'all' ? '/pages/order/list' : `/pages/order/list?status=${status}`
  uni.navigateTo({ url })
}

// 跳转新增订单
const goToAddOrder = () => {
  uni.navigateTo({ url: '/pages/order/add' })
}

// 跳转日程
const goToSchedule = () => {
  uni.navigateTo({ url: '/pages/schedule/schedule' })
}

// 切换语言
const switchLang = () => {
  const newLang = locale.value === 'zh' ? 'en' : 'zh'
  setLocale(newLang)
  locale.value = newLang
}

// 退出登录
const handleLogout = () => {
  uni.showModal({
    title: t('index.logoutConfirm'),
    content: t('index.logoutMessage'),
    success: (res) => {
      if (res.confirm) {
        uni.removeStorageSync('sessionCookie')
        uni.removeStorageSync('userInfo')
        uni.reLaunch({
          url: '/pages/login/login'
        })
      }
    }
  })
}

// 生命周期
onMounted(() => {
  locale.value = getCurrentLocale()
  fetchStats()
})
</script>

<style lang="scss" scoped>
.index-container {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 20rpx 30rpx;
  padding-bottom: 140rpx;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 24rpx;
  padding: 40rpx 30rpx;
  margin-bottom: 30rpx;
}

.title {
  display: block;
  font-size: 44rpx;
  font-weight: 700;
  color: #fff;
}

.subtitle {
  display: block;
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 8rpx;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20rpx;
  margin-bottom: 30rpx;
}

.stat-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 24rpx 16rpx;
  text-align: center;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
}

.stat-num {
  display: block;
  font-size: 40rpx;
  font-weight: 700;
  color: #667eea;
}

.stat-label {
  display: block;
  font-size: 22rpx;
  color: #999;
  margin-top: 8rpx;
}

.section {
  margin-bottom: 30rpx;
}

.section-title {
  display: block;
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 20rpx;
}

.action-list {
  background-color: #fff;
  border-radius: 16rpx;
  overflow: hidden;
}

.action-item {
  display: flex;
  align-items: center;
  padding: 32rpx 30rpx;
  border-bottom: 1rpx solid #f0f0f0;
  
  &:last-child {
    border-bottom: none;
  }
}

.action-icon {
  font-size: 40rpx;
  margin-right: 20rpx;
}

.action-text {
  flex: 1;
  font-size: 30rpx;
  color: #333;
}

.action-arrow {
  font-size: 28rpx;
  color: #ccc;
}

.lang-card {
  display: flex;
  align-items: center;
  background-color: #fff;
  border-radius: 16rpx;
  padding: 32rpx 30rpx;
}

.lang-icon {
  font-size: 40rpx;
  margin-right: 20rpx;
}

.lang-info {
  flex: 1;
}

.lang-title {
  display: block;
  font-size: 30rpx;
  color: #333;
}

.lang-value {
  display: block;
  font-size: 24rpx;
  color: #999;
  margin-top: 4rpx;
}

.logout-section {
  margin-top: 40rpx;
  padding-bottom: 40rpx;
}

.logout-btn {
  width: 100%;
  height: 88rpx;
  background-color: #fff;
  color: #e53935;
  font-size: 32rpx;
  border-radius: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;

  &::after {
    border: none;
  }
}
</style>
