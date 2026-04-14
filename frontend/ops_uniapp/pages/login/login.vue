<template>
  <view class="min-h-screen bg-gradient-to-br from-purple-500 to-purple-800 px-6" :style="{ paddingTop: statusBarHeight + 'px' }">
    <!-- 语言切换 -->
    <view class="flex justify-center items-center py-4">
      <text 
        class="text-sm px-2"
        :class="locale === 'zh' ? 'text-white font-semibold' : 'text-white/60'"
        @click="switchLang('zh')"
      >中文</text>
      <text class="text-white/40 mx-1">|</text>
      <text 
        class="text-sm px-2"
        :class="locale === 'en' ? 'text-white font-semibold' : 'text-white/60'"
        @click="switchLang('en')"
      >EN</text>
    </view>

    <!-- Logo 区域 -->
    <view class="flex flex-col items-center mb-16">
      <image class="w-20 h-20 rounded-2xl bg-white shadow-lg" src="/static/logo.png" mode="aspectFit"></image>
      <text class="text-4xl font-bold text-white mt-4 tracking-wider">OPS</text>
      <text class="text-sm text-white/80 mt-2">{{ t('login.title') }}</text>
    </view>

    <!-- 登录表单 -->
    <view class="bg-white rounded-3xl p-8 shadow-xl">
      <view class="mb-6">
        <text class="block text-sm font-medium mb-3" style="color: #1f2937">{{ t('login.username') }}</text>
        <input
          class="w-full h-11 rounded-xl px-4 text-base"
          style="background: #f3f4f6; color: #1f2937"
          type="text"
          v-model="form.username"
          :placeholder="t('login.usernamePlaceholder')"
          placeholder-class="placeholder-gray"
          @confirm="handleLogin"
        />
      </view>

      <view class="mb-6 relative">
        <text class="block text-sm font-medium mb-3" style="color: #1f2937">{{ t('login.password') }}</text>
        <input
          class="w-full h-11 rounded-xl px-4 text-base pr-10"
          style="background: #f3f4f6; color: #1f2937"
          :type="showPassword ? 'text' : 'password'"
          v-model="form.password"
          :placeholder="t('login.passwordPlaceholder')"
          placeholder-class="placeholder-gray"
          @confirm="handleLogin"
        />
        <view class="absolute right-3 bottom-2.5 text-xl" @click="showPassword = !showPassword">
          <text>{{ showPassword ? '👁️' : '👁️‍🗨️' }}</text>
        </view>
      </view>

      <view class="mb-8 flex items-center">
        <checkbox-group @change="onRememberChange">
          <label class="flex items-center text-sm" style="color: #6b7280">
            <checkbox value="1" :checked="form.remember" color="#667eea" />
            <text class="ml-2">{{ t('login.rememberMe') }}</text>
          </label>
        </checkbox-group>
      </view>

      <button
        class="w-full h-12 rounded-full flex items-center justify-center font-semibold text-lg"
        :style="{ background: 'linear-gradient(to right, #667eea, #764ba2)', color: '#fff' }"
        :loading="loading"
        :disabled="loading"
        @click="handleLogin"
      >
        {{ loading ? t('login.logging') : t('login.loginBtn') }}
      </button>

      <view class="mt-4 p-4 rounded-xl text-center" style="background: #fef2f2" v-if="errorMsg">
        <text class="text-sm" style="color: #ef4444">{{ errorMsg }}</text>
      </view>
    </view>

    <!-- 底部信息 -->
    <view class="flex justify-center items-end pb-8 mt-8">
      <text class="text-sm text-white/80">{{ t('login.registerLink') }}</text>
      <navigator url="/pages/register/register" class="text-sm text-white font-semibold ml-1">{{ t('register.title') }}</navigator>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { setLocale, getCurrentLocale } from '../../locales/index.js'

const { t, locale } = useI18n()

const statusBarHeight = ref(20)
const form = reactive({
  username: '',
  password: '',
  remember: false
})

const showPassword = ref(false)
const loading = ref(false)
const errorMsg = ref('')

const switchLang = (lang) => {
  setLocale(lang)
  locale.value = lang
}

const onRememberChange = (e) => {
  form.remember = e.detail.value.length > 0
}

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

const handleLogin = () => {
  if (!validate()) return

  loading.value = true
  errorMsg.value = ''

  uni.request({
    url: getApp().globalData.BASE_URL + '/users/login',
    method: 'POST',
    data: {
      userCookieValue: '',
      data: {
        username: form.username,
        password: form.password,
        remember: form.remember
      }
    },
    header: {
      'Content-Type': 'application/json'
    },
    success: (res) => {
      if (res.data.err_code === 0 && res.data.return && res.data.return.cookie) {
        const cookieData = res.data.return.cookie
        // 存储 cookie Value 作为 session 标识
        uni.setStorageSync('sessionCookie', cookieData.Value)
        uni.setStorageSync('cookieExpires', cookieData.ExpiresAt)
        uni.setStorageSync('userInfo', {
          userId: cookieData.ID,
          username: form.username
        })

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

        setTimeout(() => {
          uni.switchTab({
            url: '/pages/index/index'
          })
        }, 1500)
      } else {
        // 根据 err_code 显示错误信息
        const errCode = res.data.err_code
        const msgMap = {
          '-41': t('login.usernameNotFound'),
          '-42': t('login.passwordIncorrect'),
          '-3': t('login.paramError'),
          '-2': t('login.requestFailed')
        }
        errorMsg.value = msgMap[errCode] || res.data.err_msg || t('login.loginFailed')
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

onMounted(() => {
  const systemInfo = uni.getSystemInfoSync()
  statusBarHeight.value = systemInfo.statusBarHeight || 20
  
  locale.value = getCurrentLocale()
  
  const savedUsername = uni.getStorageSync('savedUsername')
  const savedRemember = uni.getStorageSync('savedRemember')
  if (savedUsername && savedRemember) {
    form.username = savedUsername
    form.remember = true
  }

  uni.$on('localeChanged', (lang) => {
    locale.value = lang
  })
})

onUnmounted(() => {
  uni.$off('localeChanged')
})
</script>

<style>
.placeholder-gray {
  color: #9ca3af;
}
</style>
