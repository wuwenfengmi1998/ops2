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

    <!-- 返回按钮 -->
    <view class="mb-8 py-3" @click="goBack">
      <text class="text-lg text-white">← {{ t('common.back') }}</text>
    </view>

    <!-- 注册表单 -->
    <view class="bg-white rounded-3xl p-8 shadow-xl">
      <text class="block text-3xl font-bold text-center mb-8" style="color: #1f2937">{{ t('register.title') }}</text>

      <view class="mb-5">
        <text class="block text-sm font-medium mb-3" style="color: #1f2937">{{ t('register.username') }}</text>
        <input
          class="w-full h-11 rounded-xl px-4 text-base"
          style="background: #f3f4f6; color: #1f2937"
          type="text"
          v-model="form.username"
          :placeholder="t('register.usernamePlaceholder')"
          placeholder-class="placeholder-gray"
        />
      </view>

      <view class="mb-5">
        <text class="block text-sm font-medium mb-3" style="color: #1f2937">{{ t('register.email') }}</text>
        <input
          class="w-full h-11 rounded-xl px-4 text-base"
          style="background: #f3f4f6; color: #1f2937"
          type="email"
          v-model="form.email"
          :placeholder="t('register.emailPlaceholder')"
          placeholder-class="placeholder-gray"
        />
      </view>

      <view class="mb-5">
        <text class="block text-sm font-medium mb-3" style="color: #1f2937">{{ t('register.password') }}</text>
        <input
          class="w-full h-11 rounded-xl px-4 text-base"
          style="background: #f3f4f6; color: #1f2937"
          :type="showPassword ? 'text' : 'password'"
          v-model="form.password"
          :placeholder="t('register.passwordPlaceholder')"
          placeholder-class="placeholder-gray"
        />
      </view>

      <view class="mb-5">
        <text class="block text-sm font-medium mb-3" style="color: #1f2937">{{ t('register.confirmPassword') }}</text>
        <input
          class="w-full h-11 rounded-xl px-4 text-base"
          style="background: #f3f4f6; color: #1f2937"
          :type="showPassword ? 'text' : 'password'"
          v-model="form.confirmPassword"
          :placeholder="t('register.confirmPlaceholder')"
          placeholder-class="placeholder-gray"
        />
      </view>

      <view class="mb-6" @click="showPassword = !showPassword">
        <text class="text-sm" style="color: #6b7280">{{ showPassword ? '👁️ ' : '👁️‍🗨️ ' }}{{ showPassword ? t('register.showPassword') : t('register.hidePassword') }}</text>
      </view>

      <button
        class="w-full h-12 rounded-full flex items-center justify-center font-semibold text-lg"
        :style="{ background: 'linear-gradient(to right, #667eea, #764ba2)', color: '#fff' }"
        :loading="loading"
        :disabled="loading"
        @click="handleRegister"
      >
        {{ loading ? t('register.registering') : t('register.registerBtn') }}
      </button>

      <view class="mt-4 p-4 rounded-xl text-center" style="background: #fef2f2" v-if="errorMsg">
        <text class="text-sm" style="color: #ef4444">{{ errorMsg }}</text>
      </view>
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
  email: '',
  password: '',
  confirmPassword: ''
})

const showPassword = ref(false)
const loading = ref(false)
const errorMsg = ref('')

const switchLang = (lang) => {
  setLocale(lang)
  locale.value = lang
}

const goBack = () => {
  uni.navigateBack()
}

const validate = () => {
  if (!form.username.trim()) {
    errorMsg.value = t('register.usernameRequired')
    return false
  }
  if (!form.email.trim()) {
    errorMsg.value = t('register.emailRequired')
    return false
  }
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

const handleRegister = () => {
  if (!validate()) return

  loading.value = true
  errorMsg.value = ''

  uni.request({
    url: getApp().globalData.BASE_URL + '/users/register',
    method: 'POST',
    data: {
      userCookieValue: '',
      data: {
        username: form.username,
        useremail: form.email,
        userpass: form.password
      }
    },
    header: {
      'Content-Type': 'application/json'
    },
    success: (res) => {
      if (res.data.err_code === 0) {
        uni.showToast({
          title: t('register.registerSuccess'),
          icon: 'success',
          duration: 1500
        })
        setTimeout(() => {
          uni.navigateBack()
        }, 1500)
      } else {
        // 根据 err_code 显示错误信息
        const errCode = res.data.err_code
        const msgMap = {
          '-4': t('register.usernameExists'),
          '-43': t('register.emailInvalid'),
          '-3': t('register.paramError'),
          '-2': t('register.requestFailed')
        }
        errorMsg.value = msgMap[errCode] || res.data.err_msg || t('register.registerFailed')
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

onMounted(() => {
  const systemInfo = uni.getSystemInfoSync()
  statusBarHeight.value = systemInfo.statusBarHeight || 20
  
  locale.value = getCurrentLocale()

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
