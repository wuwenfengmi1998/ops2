import { createI18n } from 'vue-i18n'
import zh from './zh.js'
import en from './en.js'

// 获取本地语言设置，默认中文
function getLocale() {
  // 从本地存储读取
  const stored = uni.getStorageSync('locale')
  if (stored) return stored

  // 获取系统语言
  const sysInfo = uni.getSystemInfoSync()
  const sysLang = sysInfo.language || 'zh-CN'
  // 简单判断：是否以 zh 开头
  return sysLang.toLowerCase().startsWith('zh') ? 'zh' : 'en'
}

// 创建 i18n 实例
export const i18n = createI18n({
  legacy: false,  // uni-app 必须用 composition API 模式
  locale: getLocale(),
  fallbackLocale: 'zh',
  messages: {
    zh,
    en
  }
})

// 切换语言
export function setLocale(locale) {
  if (i18n.global.locale && typeof i18n.global.locale.value !== 'undefined') {
    i18n.global.locale.value = locale
  } else {
    i18n.global.locale = locale
  }
  uni.setStorageSync('locale', locale)
  // 触发页面更新
  uni.$emit('localeChanged', locale)
}

// 获取当前语言
export function getCurrentLocale() {
  if (i18n.global.locale && typeof i18n.global.locale.value !== 'undefined') {
    return i18n.global.locale.value
  }
  return i18n.global.locale
}
