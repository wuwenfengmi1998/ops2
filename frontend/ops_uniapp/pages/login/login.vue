<template>
  <view class="login-container">
    <!-- 语言切换 -->
    <view class="lang-switch">
      <text 
        class="lang-btn" 
        :class="{ active: locale === 'zh' }"
        @click="switchLang('zh')"
      >中文</text>
      <text class="lang-divider">|</text>
      <text 
        class="lang-btn" 
        :class="{ active: locale === 'en' }"
        @click="switchLang('en')"
      >EN</text>
    </view>

    <!-- Logo 区域 -->
    <view class="logo-section">
      <image class="logo" src="/static/logo.png" mode="aspectFit"></image>
      <text class="app-name">OPS</text>
      <text class="app-desc">{{ t('login.title') }}</text>
    </view>

    <!-- 登录表单 -->
    <view class="form-section">
      <view class="input-group">
        <text class="input-label">{{ t('login.username') }}</text>
        <input
          class="input-field"
          type="text"
          v-model="form.username"
          :placeholder="t('login.usernamePlaceholder')"
          placeholder-class="placeholder"
          @confirm="handleLogin"
        />
      </view>

      <view class="input-group">
        <text class="input-label">{{ t('login.password') }}</text>
        <input
          class="input-field"
          :type="showPassword ? 'text' : 'password'"
          v-model="form.password"
          :placeholder="t('login.passwordPlaceholder')"
          placeholder-class="placeholder"
          @confirm="handleLogin"
        />
        <view class="password-toggle" @click="showPassword = !showPassword">
          <text>{{ showPassword ? '👁️' : '👁️‍🗨️' }}</text>
        </view>
      </view>

      <view class="remember-row">
        <checkbox-group @change="onRememberChange">
          <label class="remember-label">
            <checkbox value="1" :checked="form.remember" color="#007AFF" />
            <text>{{ t('login.rememberMe') }}</text>
          </label>
        </checkbox-group>
      </view>

      <button
        class="login-btn"
        :loading="loading"
        :disabled="loading"
        @click="handleLogin"
      >
        {{ loading ? t('login.logging') : t('login.loginBtn') }}
      </button>

      <view class="error-tip" v-if="errorMsg">
        <text>{{ errorMsg }}</text>
      </view>
    </view>

    <!-- 底部信息 -->
    <view class="footer-section">
      <text class="footer-text">{{ t('login.registerLink') }}</text>
      <navigator url="/pages/register/register" class="link">{{ t('register.title') }}</navigator>
    </view>
	
	
  </view>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { setLocale, getCurrentLocale } from '../../locales/index.js'

const { t, locale } = useI18n()

// 响应式数据
const form = reactive({
  username: '',
  password: '',
  remember: false
})

const showPassword = ref(false)
const loading = ref(false)
const errorMsg = ref('')

// 切换语言
const switchLang = (lang) => {
  setLocale(lang)
  locale.value = lang
}

// 记住密码勾选
const onRememberChange = (e) => {
  form.remember = e.detail.value.length > 0
}

// 表单验证
const validate = () => {
  if (!form.username.trim()) {
    errorMsg.value = t('login.usernameRequired')
    return false
  }
  if (!form.password) {
    errorMsg.value = t('login.passwordRequired')
    return false
  }
  errorMsg.value = ''
  return true
}

// 登录处理
const handleLogin = () => {
  if (!validate()) return

  loading.value = true
  errorMsg.value = ''

  uni.request({
    url: getApp().globalData.BASE_URL + '/users/login',
    method: 'POST',
    data: {
      username: form.username,
      password: form.password,
      remember: form.remember
    },
    header: {
      'Content-Type': 'application/json'
    },
    success: (res) => {
      if (res.data.code === 0 && res.data.data && res.data.data.cookie) {
        // 登录成功
        const cookie = res.data.data.cookie

        // 保存 cookie 到本地
        uni.setStorageSync('sessionCookie', cookie.Value)
        uni.setStorageSync('userInfo', res.data.data)

        // 处理记住密码
        if (form.remember) {
          uni.setStorageSync('savedUsername', form.username)
          uni.setStorageSync('savedRemember', true)
        } else {
          uni.removeStorageSync('savedUsername')
          uni.removeStorageSync('savedRemember')
        }

        uni.showToast({
          title: t('login.loginSuccess'),
          icon: 'success',
          duration: 1500
        })

        // 跳转到首页
        setTimeout(() => {
          uni.switchTab({
            url: '/pages/index/index'
          })
        }, 1500)
      } else {
        // 登录失败
        const msgMap = {
          userNameNoFund: t('login.usernameNotFound'),
          userPassIncorrect: t('login.passwordIncorrect'),
          jsonErr: t('login.paramError'),
          postErr: t('login.requestFailed')
        }
        errorMsg.value = msgMap[res.data.code] || t('login.loginFailed')
      }
    },
    fail: (err) => {
      errorMsg.value = t('login.networkError')
      console.error('Login error:', err)
    },
    complete: () => {
      loading.value = false
    }
  })
}

// 生命周期
onMounted(() => {
  // 同步当前语言
  locale.value = getCurrentLocale()
  
  // 从本地存储读取记住的登录状态
  const savedUsername = uni.getStorageSync('savedUsername')
  const savedRemember = uni.getStorageSync('savedRemember')
  if (savedUsername && savedRemember) {
    form.username = savedUsername
    form.remember = true
  }

  // 监听语言切换
  uni.$on('localeChanged', (lang) => {
    locale.value = lang
  })
})

onUnmounted(() => {
  uni.$off('localeChanged')
})
</script>

<style lang="scss" scoped>
.login-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 80rpx 60rpx;
  display: flex;
  flex-direction: column;
}

.lang-switch {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20rpx;
}

.lang-btn {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.6);
  padding: 8rpx 16rpx;
  
  &.active {
    color: #fff;
    font-weight: 600;
  }
}

.lang-divider {
  color: rgba(255, 255, 255, 0.4);
  margin: 0 8rpx;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 80rpx;
}

.logo {
  width: 160rpx;
  height: 160rpx;
  border-radius: 32rpx;
  background-color: #fff;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.15);
}

.app-name {
  font-size: 56rpx;
  font-weight: 700;
  color: #fff;
  margin-top: 24rpx;
  letter-spacing: 4rpx;
}

.app-desc {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 8rpx;
}

.form-section {
  background-color: #fff;
  border-radius: 24rpx;
  padding: 48rpx 40rpx;
  box-shadow: 0 16rpx 48rpx rgba(0, 0, 0, 0.15);
}

.input-group {
  margin-bottom: 32rpx;
  position: relative;
}

.input-label {
  display: block;
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  margin-bottom: 16rpx;
}

.input-field {
  width: 100%;
  height: 88rpx;
  background-color: #f5f5f5;
  border-radius: 16rpx;
  padding: 0 32rpx;
  font-size: 30rpx;
  color: #333;
}

.placeholder {
  color: #999;
}

.password-toggle {
  position: absolute;
  right: 24rpx;
  bottom: 24rpx;
  font-size: 36rpx;
}

.remember-row {
  margin-bottom: 40rpx;
}

.remember-label {
  display: flex;
  align-items: center;
  font-size: 28rpx;
  color: #666;
}

.login-btn {
  width: 100%;
  height: 96rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-size: 34rpx;
  font-weight: 600;
  border-radius: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;

  &::after {
    border: none;
  }

  &:active {
    opacity: 0.9;
  }

  &[disabled] {
    opacity: 0.6;
  }
}

.error-tip {
  margin-top: 24rpx;
  padding: 20rpx;
  background-color: #fff5f5;
  border-radius: 12rpx;
  text-align: center;

  text {
    color: #e53935;
    font-size: 26rpx;
  }
}

.footer-section {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: flex-end;
  padding-bottom: 40rpx;
}

.footer-text {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
}

.link {
  font-size: 28rpx;
  color: #fff;
  font-weight: 600;
  margin-left: 8rpx;
}
</style>
