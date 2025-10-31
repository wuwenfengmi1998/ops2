import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createI18n } from 'vue-i18n'

import App from './App.vue'
import router from './router'

import '@tabler/core/dist/css/tabler.min.css'
import '@tabler/core/dist/css/tabler-vendors.min.css'
import '@tabler/core/dist/js/tabler.min.js'

// import 'bootstrap/dist/css/bootstrap.min.css'
// import 'bootstrap/dist/js/bootstrap.bundle.min.js'


const app = createApp(App)

// 添加全局变量
app.config.globalProperties.$appName = 'My Vue App'
app.config.globalProperties.$apiUrl = 'https://api.example.com'
app.config.globalProperties.$currentUser = {
  name: 'John Doe',
  role: 'admin'
}


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

app.use(createPinia())
app.use(router)
app.use(i18n)

app.mount('#app')
