<template>
  <view class="container">
    <!-- 未登录状态 -->
    <view v-if="!isLoggedIn" class="login-prompt">
      <view class="prompt-icon">👤</view>
      <text class="prompt-text">{{ t('index.pleaseLogin') }}</text>
      <button class="login-btn" @click="goLogin">{{ t('index.login') }}</button>
    </view>

    <!-- 已登录状态 -->
    <view v-else class="user-info">
      <!-- 用户卡片 -->
      <view class="user-card">
        <view class="avatar-box">
          <view class="avatar-placeholder">👤</view>
        </view>
        <view class="user-details">
          <text class="username">{{ userInfo.username || userInfo.Name || 'User' }}</text>
          <text class="email">{{ userInfo.email || userInfo.Email || '' }}</text>
        </view>
      </view>

      <!-- 加载状态 -->
      <view v-if="loading" class="loading-box">
        <text>{{ t('common.loading') }}</text>
      </view>

      <!-- 用户详细信息 -->
      <view v-else-if="userDetail" class="detail-section">
        <view class="section-title">{{ t('user.profile') }}</view>
        <view class="info-list">
          <view class="info-item">
            <text class="info-label">{{ t('user.username') }}</text>
            <text class="info-value">{{ userDetail.userInfo.Username || '-' }}</text>
          </view>
          <view class="info-item">
            <text class="info-label">{{ t('user.firstName') }}</text>
            <text class="info-value">{{ userDetail.userInfo.FirstName || '-' }}</text>
          </view>
          <view class="info-item">
            <text class="info-label">{{ t('user.birthday') }}</text>
            <text class="info-value">{{ userDetail.userInfo.Birthdate || '-' }}</text>
          </view>
          <view class="info-item">
            <text class="info-label">{{ t('user.gender') }}</text>
            <text class="info-value">{{ genderText }}</text>
          </view>
          <view class="info-item">
            <text class="info-label">{{ t('user.region') }}</text>
            <text class="info-value">{{ userDetail.userInfo.Region || '-' }}</text>
          </view>
        </view>
      </view>

      <!-- 操作按钮 -->
      <view class="action-section">
        <button class="action-btn logout" @click="handleLogout">
          {{ t('index.logout') }}
        </button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { onShow } from '@dcloudio/uni-app'

const { t, locale } = useI18n()

const isLoggedIn = ref(false)
const loading = ref(false)
const userInfo = ref({})
const userDetail = ref(null)

// 性别显示文本
const genderText = computed(() => {
  if (!userDetail.value?.userInfo?.Gender) return '-'
  const map = { M: t('user.male'), F: t('user.female'), U: t('user.unknown') }
  return map[userDetail.value.userInfo.Gender] || '-'
})

// 检查登录状态
const checkLogin = () => {
  const cookie = uni.getStorageSync('sessionCookie')
  const info = uni.getStorageSync('userInfo')
  isLoggedIn.value = !!cookie && !!info
  if (info) userInfo.value = info
}

// 获取用户详细信息
const fetchUserDetail = () => {
  if (!isLoggedIn.value) return

  loading.value = true
  const cookie = uni.getStorageSync('sessionCookie')

  uni.request({
    url: getApp().globalData.BASE_URL + '/users/getinfo',
    method: 'POST',
    data: {
      userCookieValue: cookie,
      data: {}
    },
    header: {
      'Content-Type': 'application/json'
    },
    success: (res) => {
      if (res.data.err_code === 0 && res.data.return) {
        userDetail.value = res.data.return
        // 更新本地存储的用户信息
        const newInfo = {
          userId: res.data.return.user?.ID,
          username: res.data.return.user?.Name,
          email: res.data.return.user?.Email,
          ...res.data.return.userInfo
        }
        uni.setStorageSync('userInfo', newInfo)
        userInfo.value = newInfo
      } else if (res.data.err_code === -44) {
        // Cookie 过期，清除登录状态
        handleLogout()
      }
    },
    fail: (err) => {
      console.error('Get user info error:', err)
    },
    complete: () => {
      loading.value = false
    }
  })
}

// 跳转到登录页
const goLogin = () => {
  uni.navigateTo({
    url: '/pages/login/login'
  })
}

// 退出登录
const handleLogout = () => {
  uni.showModal({
    title: t('index.logoutConfirm'),
    content: t('index.logoutMessage'),
    success: (res) => {
      if (res.confirm) {
        uni.removeStorageSync('sessionCookie')
        uni.removeStorageSync('cookieExpires')
        uni.removeStorageSync('userInfo')
        isLoggedIn.value = false
        userInfo.value = {}
        userDetail.value = null
        uni.showToast({
          title: t('index.logout'),
          icon: 'success'
        })
      }
    }
  })
}

// 页面显示时刷新数据
onShow(() => {
  checkLogin()
  if (isLoggedIn.value) {
    fetchUserDetail()
  }
})
</script>

<style>
.container {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 20rpx;
}

.login-prompt {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding-top: 200rpx;
}

.prompt-icon {
  font-size: 120rpx;
  margin-bottom: 30rpx;
}

.prompt-text {
  font-size: 32rpx;
  color: #666;
  margin-bottom: 40rpx;
}

.login-btn {
  width: 300rpx;
  height: 88rpx;
  line-height: 88rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border-radius: 44rpx;
  font-size: 32rpx;
  border: none;
}

.user-card {
  display: flex;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20rpx;
  padding: 40rpx;
  margin-bottom: 30rpx;
}

.avatar-box {
  width: 120rpx;
  height: 120rpx;
  border-radius: 60rpx;
  background: rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 30rpx;
}

.avatar-placeholder {
  font-size: 60rpx;
}

.user-details {
  flex: 1;
}

.username {
  display: block;
  font-size: 36rpx;
  font-weight: bold;
  color: #fff;
  margin-bottom: 10rpx;
}

.email {
  display: block;
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.8);
}

.loading-box {
  text-align: center;
  padding: 60rpx;
  color: #999;
}

.detail-section {
  background: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 30rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 20rpx;
  padding-bottom: 20rpx;
  border-bottom: 1rpx solid #eee;
}

.info-list {
  display: flex;
  flex-direction: column;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 24rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 28rpx;
  color: #666;
}

.info-value {
  font-size: 28rpx;
  color: #333;
}

.action-section {
  margin-top: 40rpx;
}

.action-btn {
  width: 100%;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 44rpx;
  font-size: 32rpx;
  border: none;
  margin-bottom: 20rpx;
}

.action-btn.logout {
  background: #fff;
  color: #ff4d4f;
  border: 2rpx solid #ff4d4f;
}
</style>
