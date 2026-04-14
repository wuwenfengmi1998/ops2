<template>
  <view class="settings-container">
    <!-- 页面标题 -->
    <view class="page-header">
      <text class="page-title">{{ t('apiConfig.title') }}</text>
    </view>

    <!-- API 地址配置 -->
    <view class="config-section">
      <text class="section-title">{{ t('apiConfig.apiUrl') }}</text>
      
      <view class="input-group">
        <input
          class="input-field"
          type="text"
          v-model="apiUrl"
          :placeholder="t('apiConfig.apiUrlPlaceholder')"
          placeholder-class="placeholder"
        />
      </view>

      <view class="info-text">
        <text>{{ t('apiConfig.format') }}: http://192.168.1.100/api/</text>
      </view>

      <view class="info-text">
        <text>{{ t('apiConfig.current') }}: {{ currentApiUrl }}</text>
      </view>

      <button
        class="save-btn"
        :loading="saving"
        @click="saveApiUrl"
      >
        {{ t('apiConfig.save') }}
      </button>
    </view>

    <!-- 测试连接 -->
    <view class="config-section">
      <text class="section-title">{{ t('apiConfig.testConnection') }}</text>
      
      <button
        class="test-btn"
        :loading="testing"
        :disabled="!apiUrl"
        @click="testConnection"
      >
        {{ testing ? t('apiConfig.testing') : t('apiConfig.testBtn') }}
      </button>

      <view class="test-result" v-if="testResult !== null">
        <text 
          class="result-text"
          :class="testResult ? 'success' : 'failed'"
        >
          {{ testResult ? t('apiConfig.connectionSuccess') : t('apiConfig.connectionFailed') }}
        </text>
      </view>
    </view>

    <!-- 语言切换 -->
    <view class="config-section">
      <text class="section-title">{{ t('apiConfig.language') }}</text>
      
      <view class="lang-options">
        <view 
          class="lang-item"
          :class="{ active: currentLang === 'zh' }"
          @click="switchLang('zh')"
        >
          <text>🇨🇳 {{ t('apiConfig.zh') }}</text>
        </view>
        <view 
          class="lang-item"
          :class="{ active: currentLang === 'en' }"
          @click="switchLang('en')"
        >
          <text>🇺🇸 {{ t('apiConfig.en') }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()

// 响应式数据
const apiUrl = ref('')
const saving = ref(false)
const testing = ref(false)
const testResult = ref(null)
const currentLang = ref('zh')

// 计算属性
const currentApiUrl = computed(() => {
  return getApp().globalData.BASE_URL || t('apiConfig.notSet')
})

// 切换语言
const switchLang = (lang) => {
  locale.value = lang
  currentLang.value = lang
  uni.setStorageSync('locale', lang)
  uni.$emit('localeChanged', lang)
}

// 保存 API 地址
const saveApiUrl = () => {
  if (!apiUrl.value) {
    uni.showToast({
      title: t('apiConfig.pleaseInput'),
      icon: 'none'
    })
    return
  }

  // 验证 URL 格式
  let url = apiUrl.value.trim()
  if (!url.startsWith('http://') && !url.startsWith('https://')) {
    url = 'http://' + url
  }
  if (!url.endsWith('/')) {
    url = url + '/'
  }

  saving.value = true
  
  // 保存到本地
  uni.setStorageSync('apiUrl', url)
  
  // 更新全局配置
  getApp().globalData.BASE_URL = url

  setTimeout(() => {
    saving.value = false
    uni.showToast({
      title: t('apiConfig.saveSuccess'),
      icon: 'success'
    })
  }, 500)
}

// 测试连接
const testConnection = () => {
  testing.value = true
  testResult.value = null

  // 使用相对路径，走代理
  uni.request({
    url: '/api/',
    method: 'GET',
    timeout: 5000,
    success: (res) => {
      // 成功返回 {"err_code":0,"err_msg":"apiOK"}
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

// 生命周期
onMounted(() => {
  // 读取当前配置
  const savedApiUrl = uni.getStorageSync('apiUrl')
  if (savedApiUrl) {
    apiUrl.value = savedApiUrl
  } else {
    apiUrl.value = getApp().globalData.BASE_URL || ''
  }
  
  currentLang.value = locale.value || 'zh'
})
</script>

<style lang="scss" scoped>
.settings-container {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 40rpx 30rpx;
  padding-bottom: 140rpx;
}

.page-header {
  margin-bottom: 40rpx;
}

.page-title {
  font-size: 40rpx;
  font-weight: 700;
  color: #333;
}

.config-section {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 32rpx;
  margin-bottom: 30rpx;
}

.section-title {
  display: block;
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 24rpx;
}

.input-group {
  margin-bottom: 20rpx;
}

.input-field {
  width: 100%;
  height: 88rpx;
  background-color: #f5f5f5;
  border-radius: 12rpx;
  padding: 0 24rpx;
  font-size: 30rpx;
  color: #333;
}

.placeholder {
  color: #999;
}

.info-text {
  margin-top: 16rpx;
  
  text {
    font-size: 24rpx;
    color: #999;
  }
}

.save-btn {
  width: 100%;
  height: 88rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-size: 32rpx;
  font-weight: 600;
  border-radius: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 30rpx;
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

.test-btn {
  width: 100%;
  height: 80rpx;
  background-color: #f0f0f0;
  color: #333;
  font-size: 30rpx;
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;

  &::after {
    border: none;
  }

  &:active {
    background-color: #e0e0e0;
  }

  &[disabled] {
    opacity: 0.5;
  }
}

.test-result {
  margin-top: 24rpx;
  text-align: center;
}

.result-text {
  font-size: 28rpx;
  font-weight: 500;

  &.success {
    color: #52c41a;
  }

  &.failed {
    color: #ff4d4f;
  }
}

.lang-options {
  display: flex;
  gap: 20rpx;
}

.lang-item {
  flex: 1;
  height: 88rpx;
  background-color: #f5f5f5;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx solid transparent;

  &.active {
    background-color: #f0f5ff;
    border-color: #667eea;
  }

  text {
    font-size: 28rpx;
    color: #333;
  }
}
</style>
