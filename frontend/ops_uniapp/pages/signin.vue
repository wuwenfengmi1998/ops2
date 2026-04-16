<script>
/**
 * 登录页
 * 对标 PC 前端 src/views/LoginView.vue
 */
import { authApi } from '../@api/auth.js'
import { userStore } from '../store/user.js'

export default {
  data() {
    return {
      form: {
        username: '',
        password: '',
        remember: false,
      },
      errors: {
        username: '',
        password: '',
      },
      showPassword: false,
      loading: false,
    }
  },

  methods: {
    validate() {
      let ok = true
      this.errors.username = ''
      this.errors.password = ''
      if (!this.form.username.trim()) {
        this.errors.username = '请输入用户名'
        ok = false
      }
      if (!this.form.password) {
        this.errors.password = '请输入密码'
        ok = false
      }
      return ok
    },

    async handleLogin() {
      if (!this.validate()) return
      this.loading = true
      try {
        const { errCode, data } = await authApi.login(
          this.form.username,
          this.form.password,
          this.form.remember,
        )
        switch (errCode) {
          case 0:
            userStore.login(data.cookie)
            uni.showToast({ title: '登录成功', icon: 'success' })
            setTimeout(() => {
              uni.reLaunch({ url: '/pages/index/index' })
            }, 800)
            break
          case -41:
          case -42:
            uni.showToast({ title: '用户名或密码错误', icon: 'none' })
            break
          default:
            uni.showToast({ title: '服务器错误，请稍后重试', icon: 'none' })
        }
      } catch {
        // request.js 已处理网络错误提示
      } finally {
        this.loading = false
      }
    },
  },
}
</script>

<template>
  <view class="signin-page">
    <!-- 顶部 Logo -->
    <view class="header">
      <image src="/static/logo.svg" class="logo" mode="aspectFit" />
      <text class="app-name">Operations</text>
    </view>

    <!-- 登录卡片 -->
    <view class="card">
      <text class="card-title">登录你的账号</text>

      <!-- 用户名 -->
      <view class="field">
        <text class="label">用户名</text>
        <input
          v-model="form.username"
          class="input"
          :class="{ 'input-error': errors.username }"
          placeholder="输入你的用户名"
          maxlength="25"
          @confirm="handleLogin"
        />
        <text v-if="errors.username" class="err-msg">{{ errors.username }}</text>
      </view>

      <!-- 密码 -->
      <view class="field">
        <text class="label">密码</text>
        <view class="input-wrap">
          <input
            v-model="form.password"
            class="input"
            :class="{ 'input-error': errors.password }"
            :password="!showPassword"
            placeholder="输入你的密码"
            maxlength="100"
            @confirm="handleLogin"
          />
          <view class="eye-btn" @tap="showPassword = !showPassword">
            <text class="eye-icon">{{ showPassword ? '👁' : '🙈' }}</text>
          </view>
        </view>
        <text v-if="errors.password" class="err-msg">{{ errors.password }}</text>
      </view>

      <!-- 记住我 -->
      <view class="remember-row" @tap="form.remember = !form.remember">
        <view class="checkbox" :class="{ checked: form.remember }">
          <text v-if="form.remember" class="check-mark">✓</text>
        </view>
        <text class="remember-text">在此设备上保持登录</text>
      </view>

      <!-- 登录按钮 -->
      <button
        class="submit-btn"
        :class="{ 'btn-loading': loading }"
        :disabled="loading"
        @tap="handleLogin"
      >
        <text v-if="loading" class="loading-icon">⟳</text>
        <text>{{ loading ? '登录中…' : '登 录' }}</text>
      </button>
    </view>
  </view>
</template>

<style scoped>
.signin-page {
  min-height: 100vh;
  background-color: #f3f4f6;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 80rpx 40rpx 60rpx;
}

/* ── 顶部 logo ── */
.header {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 60rpx;
}
.logo {
  width: 120rpx;
  height: 120rpx;
  border-radius: 24rpx;
  margin-bottom: 20rpx;
}
.app-name {
  font-size: 44rpx;
  font-weight: 700;
  color: #1f2937;
}

/* ── 卡片 ── */
.card {
  width: 100%;
  max-width: 680rpx;
  background: #fff;
  border-radius: 24rpx;
  padding: 60rpx 48rpx;
  box-shadow: 0 4rpx 24rpx rgba(0, 0, 0, 0.08);
}
.card-title {
  font-size: 40rpx;
  font-weight: 700;
  color: #111827;
  text-align: center;
  display: block;
  margin-bottom: 48rpx;
}

/* ── 表单字段 ── */
.field {
  margin-bottom: 36rpx;
}
.label {
  font-size: 28rpx;
  font-weight: 500;
  color: #374151;
  display: block;
  margin-bottom: 12rpx;
}
.input-wrap {
  position: relative;
  display: flex;
  align-items: center;
}
.input {
  width: 100%;
  height: 88rpx;
  border: 2rpx solid #d1d5db;
  border-radius: 16rpx;
  padding: 0 32rpx;
  font-size: 30rpx;
  color: #111827;
  background: #fff;
  box-sizing: border-box;
}
.input-error {
  border-color: #ef4444;
}
.eye-btn {
  position: absolute;
  right: 24rpx;
  height: 88rpx;
  display: flex;
  align-items: center;
  padding: 0 8rpx;
}
.eye-icon {
  font-size: 36rpx;
}
.err-msg {
  font-size: 24rpx;
  color: #ef4444;
  margin-top: 8rpx;
  display: block;
}

/* ── 记住我 ── */
.remember-row {
  display: flex;
  align-items: center;
  margin-bottom: 48rpx;
  gap: 16rpx;
}
.checkbox {
  width: 40rpx;
  height: 40rpx;
  border: 2rpx solid #d1d5db;
  border-radius: 8rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
}
.checkbox.checked {
  background: #2563eb;
  border-color: #2563eb;
}
.check-mark {
  color: #fff;
  font-size: 26rpx;
  line-height: 1;
}
.remember-text {
  font-size: 28rpx;
  color: #4b5563;
}

/* ── 登录按钮 ── */
.submit-btn {
  width: 100%;
  height: 96rpx;
  background: #2563eb;
  color: #fff;
  border-radius: 16rpx;
  font-size: 32rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  border: none;
}
.submit-btn[disabled] {
  background: #93c5fd;
}
.loading-icon {
  font-size: 36rpx;
  animation: spin 1s linear infinite;
  display: inline-block;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to   { transform: rotate(360deg); }
}
</style>
