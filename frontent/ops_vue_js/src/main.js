import './assets/main.css'

import { createApp } from 'vue'
import { createI18n } from 'vue-i18n'
import { createPinia } from 'pinia' // 1. 导入 createPinia
import App from './App.vue'
import router from './router'

import '@tabler/core/dist/css/tabler.min.css'


import en from './i18n/en.json'
import zhCN from './i18n/zh-CN.json'

const i18n = createI18n({
  legacy: false, // 使用 Composition API 模式
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en,
    'zh-CN': zhCN
  }
})

const pinia = createPinia()
const app = createApp(App)

app.use(router)
app.use(i18n)
app.use(pinia)

app.mount('#app')
