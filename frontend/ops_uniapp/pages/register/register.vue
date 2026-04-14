<template>
  <view class="register-container">
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

    <!-- 返回按钮 -->
    <view class="back-btn" @click="goBack">
      <text>← {{ t('common.back') }}</text>
    </view>

    <view class="form-section">
      <text class="title">{{ t('register.title') }}</text>

      <view class="input-group">
        <text class="input-label">{{ t('register.username') }}</text>
        <input
          class="input-field"
          type="text"
          v-model="form.username"
          :placeholder="t('register.usernamePlaceholder')"
          placeholder-class="placeholder"
        />
      </view>

      <view class="input-group">
        <text class="input-label">{{ t('register.email') }}</text>
        <input
          class="input-field"
          type="email"
          v-model="form.email"
          :placeholder="t('register.emailPlaceholder')"
          placeholder-class="placeholder"
        />
      </view>

      <view class="input-group">
        <text class="input-label">{{ t('register.password') }}</text>
        <input
          class="input-field"
          :type="showPassword ? 'text' : 'password'"
          v-model="form.password"
          :placeholder="t('register.passwordPlaceholder')"
          placeholder-class="placeholder"
        />
      </view>

      <view class="input-group">
        <text class="input-label">{{ t('register.confirmPassword') }}</text>
        <input
          class="input-field"
          :type="showPassword ? 'text' : 'password'"
          v-model="form.confirmPassword"
          :placeholder="t('register.confirmPlaceholder')"
          placeholder-class="placeholder"
        />
      </view>

      <view class="password-toggle" @click="showPassword = !showPassword">
        <text>{{ showPassword ? '👁️ ' : '👁️‍🗨️ ' }}{{ showPassword ? t('register.showPassword') : t('register.hidePassword') }}</text>
      </view>

      <button
        class="register-btn"
        :loading="loading"
        :disabled="loading"
        @click="handleRegister"
      >
        {{ loading ? t('register.registering') : t('register.registerBtn') }}
      </button>

      <view class="error-tip" v-if="errorMsg">
        <text>{{ errorMsg }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { setLocale, getCurrentLocale } from '../../locales/index.js'

const { t, locale } = useI18n()

// 响应式数据
const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const showPassword = ref(false)
const loading = ref(false)
const errorMsg = ref('')

// 切换语言
const switchLang = (lang) => {
  setLocale(lang)
  locale.value = lang
}

// 返回上一页
const goBack = () => {
  uni.navigateBack()
}

// 表单验证
const validate = () => {
  if (!form.username.trim()) {
    errorMsg.value = t('register.usernameRequired')
    return false
  }
  if (!form.email.trim()) {
    errorMsg.value = t('register.emailRequired')
    return false
  }
  // 简单邮箱验证
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(form.email)) {
    errorMsg.value = t('register.emailInvalid')
    return false
  }
  if (!form.password) {
    errorMsg.value = t('register.passwordRequired')
    return false
  }
  if (form.password.length < 6) {
    errorMsg.value = t('register.passwordLength')
    return false
  }
  if (form.password !== form.confirmPassword) {
    errorMsg.value = t('register.passwordMismatch')
    return false
  }
  errorMsg.value = ''
  return true
}

// 注册处理
const handleRegister = () => {
  if (!validate()) return

  loading.value = true
  errorMsg.value = ''

  uni.request({
    url: getApp().globalData.BASE_URL + '/users/register',
    method: 'POST',
    data: {
      username: form.username,
      useremail: form.email,
      userpass: form.password
    },
    header: {
      'Content-Type': 'application/json'
    },
    success: (res) => {
      if (res.data.code === 0) {
        uni.showToast({
          title: t('register.registerSuccess'),
          icon: 'success',
          duration: 1500
        })
        setTimeout(() => {
          uni.navigateBack()
        }, 1500)
      } else {
        const msgMap = {
          userNameDup: t('register.usernameExists'),
          userEmailDup: t('register.emailUsed'),
          jsonErr: t('register.paramError'),
          postErr: t('register.requestFailed')
        }
        errorMsg.value = msgMap[res.data.code] || t('register.registerFailed')
      }
    },
    fail: (err) => {
      errorMsg.value = t('common.networkError')
      console.error('Register error:', err)
    },
    complete: () => {
      loading.value = false
    }
  })
}

// 生命周期
onMounted(() => {
  locale.value = getCurrentLocale()

  uni.$on('localeChanged', (lang) => {
    locale.value = lang
  })
})

onUnmounted(() => {
  uni.$off('localeChanged')
})
</script>

<style lang="scss" scoped>
.register-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 40rpx 60rpx;
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

.back-btn {
  margin-bottom: 40rpx;
  padding: 16rpx 0;

  text {
    color: #fff;
    font-size: 32rpx;
  }
}

.form-section {
  background-color: #fff;
  border-radius: 24rpx;
  padding: 48rpx 40rpx;
  box-shadow: 0 16rpx 48rpx rgba(0, 0, 0, 0.15);
}

.title {
  display: block;
  font-size: 44rpx;
  font-weight: 700;
  color: #333;
  text-align: center;
  margin-bottom: 48rpx;
}

.input-group {
  margin-bottom: 28rpx;
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
  margin-bottom: 32rpx;

  text {
    color: #666;
    font-size: 26rpx;
  }
}

.register-btn {
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
</style>
