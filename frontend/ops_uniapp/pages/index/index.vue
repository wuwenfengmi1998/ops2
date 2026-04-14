<template>
  <view class="min-h-screen" :style="{ paddingTop: statusBarHeight + 'px' }">
    <!-- 未登录状态 -->
    <view v-if="!isLoggedIn" class="flex flex-col items-center justify-center px-6" style="min-height: 80vh">
      <text class="text-6xl mb-6">🔐</text>
      <text class="block text-xl font-semibold text-center mb-8" style="color: #1f2937">{{ t('index.pleaseLogin') }}</text>
      <button
        class="w-full max-w-xs h-12 rounded-full flex items-center justify-center font-semibold text-base"
        :style="{ background: 'linear-gradient(to right, #667eea, #764ba2)', color: '#fff' }"
        @click="goToLogin"
      >
        {{ t('index.login') }}
      </button>
    </view>

    <!-- 已登录状态 -->
    <template v-else>
      <!-- 顶部标题 -->
      <view class="header mx-4 mt-4 rounded-3xl p-8 mb-6" :style="{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }">
        <text class="block text-3xl font-bold text-white">{{ t('index.welcome') }}</text>
        <text class="block text-base text-white/80 mt-2">{{ t('index.subtitle') }}</text>
      </view>

      <!-- 订单统计卡片 -->
      <view class="grid grid-cols-4 gap-3 mx-4 mb-6">
        <view class="stat-card rounded-2xl p-4 text-center shadow-sm" @click="goToOrders('pending')">
          <text class="block text-3xl font-bold" style="color: #667eea">{{ stats.pending }}</text>
          <text class="block text-sm mt-2" style="color: #6b7280">{{ t('index.pending') }}</text>
        </view>
        <view class="stat-card rounded-2xl p-4 text-center shadow-sm" @click="goToOrders('ordered')">
          <text class="block text-3xl font-bold" style="color: #667eea">{{ stats.ordered }}</text>
          <text class="block text-sm mt-2" style="color: #6b7280">{{ t('index.ordered') }}</text>
        </view>
        <view class="stat-card rounded-2xl p-4 text-center shadow-sm" @click="goToOrders('arrived')">
          <text class="block text-3xl font-bold" style="color: #667eea">{{ stats.arrived }}</text>
          <text class="block text-sm mt-2" style="color: #6b7280">{{ t('index.arrived') }}</text>
        </view>
        <view class="stat-card rounded-2xl p-4 text-center shadow-sm" @click="goToOrders('received')">
          <text class="block text-3xl font-bold" style="color: #667eea">{{ stats.received }}</text>
          <text class="block text-sm mt-2" style="color: #6b7280">{{ t('index.received') }}</text>
        </view>
      </view>

      <!-- 快捷操作 -->
      <view class="mx-4 mb-6">
        <text class="block text-xl font-semibold mb-4" style="color: #1f2937">{{ t('index.quickActions') }}</text>
        <view class="bg-white rounded-2xl overflow-hidden shadow-sm">
          <view class="action-item flex items-center p-6 border-b border-gray-100" @click="goToOrders('all')">
            <text class="text-2xl mr-4">📋</text>
            <text class="flex-1 text-base" style="color: #1f2937">{{ t('index.allOrders') }}</text>
            <text class="text-base" style="color: #d1d5db">→</text>
          </view>
          <view class="action-item flex items-center p-6 border-b border-gray-100" @click="goToAddOrder">
            <text class="text-2xl mr-4">➕</text>
            <text class="flex-1 text-base" style="color: #1f2937">{{ t('index.addOrder') }}</text>
            <text class="text-base" style="color: #d1d5db">→</text>
          </view>
          <view class="action-item flex items-center p-6" @click="goToSchedule">
            <text class="text-2xl mr-4">📅</text>
            <text class="flex-1 text-base" style="color: #1f2937">{{ t('index.schedule') }}</text>
            <text class="text-base" style="color: #d1d5db">→</text>
          </view>
        </view>
      </view>

      <!-- 语言切换入口 -->
      <view class="mx-4 mb-6">
        <view class="bg-white rounded-2xl p-6 flex items-center shadow-sm" @click="switchLang">
          <text class="text-2xl mr-4">🌐</text>
          <view class="flex-1">
            <text class="block text-base" style="color: #1f2937">{{ t('index.language') }}</text>
            <text class="block text-sm mt-1" style="color: #6b7280">{{ locale === 'zh' ? '中文' : 'English' }}</text>
          </view>
          <text class="text-base" style="color: #d1d5db">→</text>
        </view>
      </view>

      <!-- 退出登录 -->
      <view class="mx-4 pb-8">
        <button class="w-full h-11 rounded-full text-center text-lg" 
          style="background: #fff; color: #e53935" @click="handleLogout">
          {{ t('index.logout') }}
        </button>
      </view>
    </template>
  </view>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { setLocale, getCurrentLocale } from '../../locales/index.js'

const { t, locale } = useI18n()

const statusBarHeight = ref(20)
const isLoggedIn = ref(false)

// API 请求封装，解决路径拼接问题
const request = (path, method = 'GET', data = {}) => {
  return new Promise((resolve, reject) => {
    // 确保 BASE_URL 不以 / 结尾
    const baseUrl = (getApp().globalData.BASE_URL || '').replace(/\/$/, '')
    // 确保 path 以 / 开头
    const fullPath = path.startsWith('/') ? path : '/' + path
    
    uni.request({
      url: baseUrl + fullPath,
      method,
      data,
      header: {
        'Content-Type': 'application/json',
        'Cookie': uni.getStorageSync('sessionCookie') || ''
      },
      success: (res) => resolve(res),
      fail: (err) => reject(err)
    })
  })
}

const stats = reactive({
  pending: 0,
  ordered: 0,
  arrived: 0,
  received: 0
})

const fetchStats = () => {
  request('/purchase/getordercount', 'POST').then((res) => {
    if (res.data.code === 0 && res.data.data) {
      stats.pending = res.data.data.pending || 0
      stats.ordered = res.data.data.ordered || 0
      stats.arrived = res.data.data.arrived || 0
      stats.received = res.data.data.received || 0
    }
  }).catch(console.error)
}

const goToLogin = () => {
  uni.navigateTo({
    url: '/pages/login/login'
  })
}

const goToOrders = (status) => {
  const url = status === 'all' ? '/pages/order/list' : `/pages/order/list?status=${status}`
  uni.navigateTo({ url })
}

const goToAddOrder = () => {
  uni.navigateTo({ url: '/pages/order/add' })
}

const goToSchedule = () => {
  uni.navigateTo({ url: '/pages/schedule/schedule' })
}

const switchLang = () => {
  const newLang = locale.value === 'zh' ? 'en' : 'zh'
  setLocale(newLang)
  locale.value = newLang
}

const handleLogout = () => {
  uni.showModal({
    title: t('index.logoutConfirm'),
    content: t('index.logoutMessage'),
    success: (res) => {
      if (res.confirm) {
        uni.removeStorageSync('sessionCookie')
        uni.removeStorageSync('userInfo')
        isLoggedIn.value = false
      }
    }
  })
}

onMounted(() => {
  const systemInfo = uni.getSystemInfoSync()
  statusBarHeight.value = systemInfo.statusBarHeight || 20
  
  locale.value = getCurrentLocale()
  
  // 检查登录状态
  const sessionCookie = uni.getStorageSync('sessionCookie')
  isLoggedIn.value = !!sessionCookie
  
  // 已登录才获取数据
  if (isLoggedIn.value) {
    fetchStats()
  }
})

onUnmounted(() => {
  uni.$off('localeChanged')
})
</script>

<style scoped>
.min-h-screen {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.flex-col { flex-direction: column; }
.items-center { align-items: center; }
.justify-center { justify-content: center; }
</style>
