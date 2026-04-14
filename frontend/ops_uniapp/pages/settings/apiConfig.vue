<template>
  <view class="min-h-screen" :style="{ paddingTop: statusBarHeight + 'px' }">
    <!-- 页面标题 -->
    <view class="mx-4 mt-4 mb-6">
      <text class="block text-2xl font-bold" style="color: #1f2937">{{ t('apiConfig.title') }}</text>
    </view>

    <!-- #ifndef H5 -->
    <!-- API 地址配置 - 非H5端显示 -->
    <view class="mx-4 mb-4 bg-white rounded-2xl p-6 shadow-sm">
      <text class="block text-base font-semibold mb-4" style="color: #1f2937">{{ t('apiConfig.apiUrl') }}</text>
      
      <view class="mb-4">
        <input
          class="w-full h-11 rounded-xl px-4 text-base"
          style="background: #f3f4f6; color: #1f2937"
          type="text"
          v-model="apiUrl"
          :placeholder="t('apiConfig.apiUrlPlaceholder')"
          placeholder-class="placeholder-gray"
        />
      </view>

      <view class="mb-3">
        <text class="block text-sm" style="color: #9ca3af">{{ t('apiConfig.format') }}: http://192.168.1.100/api/</text>
      </view>

      <view class="mb-4">
        <text class="block text-sm" style="color: #9ca3af">{{ t('apiConfig.current') }}: {{ currentApiUrl }}</text>
      </view>

      <button
        class="w-full h-11 rounded-full flex items-center justify-center font-semibold text-base"
        :style="{ background: 'linear-gradient(to right, #667eea, #764ba2)', color: '#fff' }"
        :loading="saving"
        @click="saveApiUrl"
      >
        {{ t('apiConfig.save') }}
      </button>
    </view>

   
    
    <!-- #endif -->
	<view class="mx-4 mb-4 bg-white rounded-2xl p-6 shadow-sm">
	  <text class="block text-base font-semibold mb-4" style="color: #1f2937">{{ t('apiConfig.testConnection') }}</text>
	  
	  <button
	    class="w-full h-10 rounded-full flex items-center justify-center text-base"
	    :style="{ background: testResult === true ? '#dcfce7' : testResult === false ? '#fef2f2' : '#f3f4f6', color: testResult === true ? '#22c55e' : testResult === false ? '#ef4444' : '#1f2937' }"
	    :loading="testing"
	    :disabled="!apiUrl"
	    @click="testConnection"
	  >
	    {{ testing ? t('apiConfig.testing') : t('apiConfig.testBtn') }}
	  </button>
	
	  <view class="mt-4 text-center" v-if="testResult !== null">
	    <text class="text-base font-medium" :style="{ color: testResult ? '#22c55e' : '#ef4444' }">
	      {{ testResult ? t('apiConfig.connectionSuccess') : t('apiConfig.connectionFailed') }}
	    </text>
	  </view>
	</view>

    <!-- 语言切换 -->
    <view class="mx-4 mb-8 bg-white rounded-2xl p-6 shadow-sm">
      <text class="block text-base font-semibold mb-4" style="color: #1f2937">{{ t('apiConfig.language') }}</text>
      
      <view class="flex gap-4">
        <view 
          class="flex-1 h-11 rounded-xl flex items-center justify-center border-2"
          :style="currentLang === 'zh' ? { borderColor: '#667eea', backgroundColor: '#f0f5ff' } : { borderColor: 'transparent', backgroundColor: '#f3f4f6' }"
          @click="switchLang('zh')"
        >
          <text class="text-base" :style="{ color: currentLang === 'zh' ? '#667eea' : '#1f2937' }">🇨🇳 中文</text>
        </view>
        <view 
          class="flex-1 h-11 rounded-xl flex items-center justify-center border-2"
          :style="currentLang === 'en' ? { borderColor: '#667eea', backgroundColor: '#f0f5ff' } : { borderColor: 'transparent', backgroundColor: '#f3f4f6' }"
          @click="switchLang('en')"
        >
          <text class="text-base" :style="{ color: currentLang === 'en' ? '#667eea' : '#1f2937' }">🇺🇸 English</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()

const statusBarHeight = ref(20)
const apiUrl = ref('')
const saving = ref(false)
const testing = ref(false)
const testResult = ref(null)
const currentLang = ref('zh')

const currentApiUrl = computed(() => {
  return getApp().globalData.BASE_URL || t('apiConfig.notSet')
})

const switchLang = (lang) => {
  locale.value = lang
  currentLang.value = lang
  uni.setStorageSync('locale', lang)
  uni.$emit('localeChanged', lang)
}

const saveApiUrl = () => {
  if (!apiUrl.value) {
    uni.showToast({
      title: t('apiConfig.pleaseInput'),
      icon: 'none'
    })
    return
  }

  let url = apiUrl.value.trim()
  if (!url.startsWith('http://') && !url.startsWith('https://')) {
    url = 'http://' + url
  }
  if (!url.endsWith('/')) {
    url = url + '/'
  }

  saving.value = true
  
  // 保存原始 URL 到本地存储
  uni.setStorageSync('apiUrl', url)
  
  // H5 端使用相对路径走代理，其他端使用完整 URL
  // #ifdef H5
  getApp().globalData.BASE_URL = '/api/'
  // #endif
  // #ifndef H5
  getApp().globalData.BASE_URL = url
  // #endif

  setTimeout(() => {
    saving.value = false
    uni.showToast({
      title: t('apiConfig.saveSuccess'),
      icon: 'success'
    })
  }, 500)
}

const testConnection = () => {
  testing.value = true
  testResult.value = null

  uni.request({
    url: getApp().globalData.BASE_URL,
    method: 'GET',
    timeout: 5000,
    success: (res) => {
      if (res.data && res.data.err_code === 0) {
        testResult.value = true
      } else {
        testResult.value = false
      }
    },
    fail: () => {
      testResult.value = false
    },
    complete: () => {
      testing.value = false
    }
  })
}

onMounted(() => {
  const systemInfo = uni.getSystemInfoSync()
  statusBarHeight.value = systemInfo.statusBarHeight || 20
  
  const savedApiUrl = uni.getStorageSync('apiUrl')
  if (savedApiUrl) {
    apiUrl.value = savedApiUrl
  } else {
    apiUrl.value = getApp().globalData.BASE_URL || ''
  }
  
  currentLang.value = locale.value || 'zh'
})
</script>

<style scoped>
.min-h-screen {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.placeholder-gray {
  color: #9ca3af;
}
</style>
