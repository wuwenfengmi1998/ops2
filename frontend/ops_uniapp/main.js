import App from './App'
import { i18n } from './locales/index.js'

// 默认 API 地址
const DEFAULT_API_URL = 'http://192.168.13.105/api/'

// #ifndef VUE3
import Vue from 'vue'
import './uni.promisify.adaptor'
Vue.config.productionTip = false

// 从本地存储读取用户配置的 API 地址
Vue.prototype.$BASE_URL = uni.getStorageSync('apiUrl') || DEFAULT_API_URL
Vue.prototype.$i18n = i18n

const app = new Vue({
  ...App,
  i18n
})
app.$mount()
// #endif

// #ifdef VUE3
import { createSSRApp } from 'vue'
export function createApp() {
  const app = createSSRApp(App)
  app.use(i18n)
  // 从本地存储读取用户配置的 API 地址
  app.config.globalProperties.$BASE_URL = uni.getStorageSync('apiUrl') || DEFAULT_API_URL
  app.config.globalProperties.BASE_URL = uni.getStorageSync('apiUrl') || DEFAULT_API_URL
  return {
    app
  }
}
// #endif
