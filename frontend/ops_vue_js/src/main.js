import { createApp } from 'vue'
import { createI18n } from 'vue-i18n'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

import './assets/main.css'

// Initialize theme and locale from localStorage
const savedTheme = localStorage.getItem('theme') // 改用 'theme' key
const savedLocale = localStorage.getItem('locale')
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches

if (savedTheme === 'dark' || (!savedTheme && prefersDark)) {
  document.documentElement.classList.add('dark')
} else {
  document.documentElement.classList.remove('dark')
}

import en from './i18n/en.json'
import zhCN from './i18n/zh-CN.json'

const i18n = createI18n({
  legacy: false,
  locale: savedLocale || 'en', // 使用已保存的语言或默认英文
  fallbackLocale: 'en',
  messages: {
    en,
    'zh-CN': zhCN,
  },
})

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(router)
app.use(i18n)

app.mount('#app')
