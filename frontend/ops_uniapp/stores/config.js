import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useConfigStore = defineStore('config', () => {
  // API 配置
  const apiBaseUrl = ref('')
  
  // 应用配置
  const appName = ref('OPS')
  const version = ref('1.0.0')
  
  // 主题配置
  const theme = ref('light')
  
  // 设置 API 地址
  const setApiBaseUrl = (url) => {
    apiBaseUrl.value = url
	uni.setStorageSync('baseUrl', url)
  }
  
  const getApiBaseUrl=()=>{
	  if(apiBaseUrl.value==='')
	  {
		  apiBaseUrl.value=uni.getStorageSync('baseUrl')
	  }
	  return apiBaseUrl.value
  }
  
  // 获取图片基础 URL（去掉末尾的 /api 部分）
  const getFileBaseUrl = () => {
    const base = getApiBaseUrl()
    if (base) {
      return base.replace(/\/api$/, '')
    }
    return base
  }

  
  // 设置主题
  const setTheme = (newTheme) => {
    theme.value = newTheme
  }
  
  return {
    apiBaseUrl,
    appName,
    version,
    theme,
    setApiBaseUrl,
	getApiBaseUrl,
    getFileBaseUrl,
    setTheme
  }
})
