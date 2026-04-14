<script>
export default {
  onLaunch: function() {
    console.log('App Launch')
    
    // 初始化全局 API 地址
    const apiUrl = uni.getStorageSync('apiUrl')
    if (!getApp().globalData) {
      getApp().globalData = {}
    }
    
    // H5 端使用相对路径走代理，其他端使用完整 URL
    // #ifdef H5
    getApp().globalData.BASE_URL = '/api/'
    // #endif
    // #ifndef H5
    getApp().globalData.BASE_URL = apiUrl
    // #endif
    
    // 暂时禁用自动跳转登录
    // this.checkLoginStatus()
  },
  onShow: function() {
    console.log('App Show')
  },
  onHide: function() {
    console.log('App Hide')
  },
  methods: {
    checkLoginStatus() {
      const sessionCookie = uni.getStorageSync('sessionCookie')
      if (!sessionCookie) {
        uni.reLaunch({
          url: '/pages/login/login'
        })
      }
    }
  }
}
</script>

<style>
/* 全局样式 */
page {
  background-color: #f5f5f5;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  font-size: 16px;
  line-height: 1.5;
  color: #1f2937;
}

/* 移除默认按钮样式 */
button::after {
  border: none;
}

/* Flex 布局类 */
.flex { display: flex; }
.flex-col { flex-direction: column; }
.flex-1 { flex: 1; }
.items-center { align-items: center; }
.justify-center { justify-content: center; }
.justify-between { justify-content: space-between; }
.items-end { align-items: flex-end; }

/* Grid */
.grid { display: flex; flex-wrap: wrap; }
.grid-cols-4 > view, .grid-cols-4 > .grid-item { width: 25%; }

/* Gap */
.gap-3 { gap: 0.75rem; }
.gap-4 { gap: 1rem; }

/* Spacing */
.mx-4 { margin-left: 1rem; margin-right: 1rem; }
.my-4 { margin-top: 1rem; margin-bottom: 1rem; }
.mt-2 { margin-top: 0.5rem; }
.mt-4 { margin-top: 1rem; }
.mb-3 { margin-bottom: 0.75rem; }
.mb-4 { margin-bottom: 1rem; }
.mb-6 { margin-bottom: 1.5rem; }
.mb-8 { margin-bottom: 2rem; }
.ml-2 { margin-left: 0.5rem; }
.ml-4 { margin-left: 1rem; }
.mr-4 { margin-right: 1rem; }
.py-3 { padding-top: 0.75rem; padding-bottom: 0.75rem; }
.py-4 { padding-top: 1rem; padding-bottom: 1rem; }
.px-2 { padding-left: 0.5rem; padding-right: 0.5rem; }
.px-4 { padding-left: 1rem; padding-right: 1rem; }
.px-6 { padding-left: 1.5rem; padding-right: 1.5rem; }
.p-4 { padding: 1rem; }
.p-6 { padding: 1.5rem; }
.pb-8 { padding-bottom: 2rem; }

/* 文字 */
.text-sm { font-size: 0.875rem; }
.text-base { font-size: 1rem; }
.text-lg { font-size: 1.125rem; }
.text-xl { font-size: 1.25rem; }
.text-2xl { font-size: 1.5rem; }
.text-3xl { font-size: 1.875rem; }
.text-4xl { font-size: 2.25rem; }
.font-medium { font-weight: 500; }
.font-semibold { font-weight: 600; }
.font-bold { font-weight: 700; }
.text-center { text-align: center; }
.tracking-wider { letter-spacing: 0.05em; }

/* 圆角 */
.rounded-xl { border-radius: 0.75rem; }
.rounded-2xl { border-radius: 1rem; }
.rounded-3xl { border-radius: 1.5rem; }
.rounded-full { border-radius: 9999px; }

/* 阴影 */
.shadow-sm { box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05); }
.shadow-lg { box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05); }
.shadow-xl { box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04); }

/* 背景 */
.bg-white { background-color: #ffffff; }
.bg-gray-100 { background-color: #f3f4f6; }
.bg-gradient-to-br { background: linear-gradient(to bottom right, var(--tw-gradient-stops)); }
.bg-gradient-to-r { background: linear-gradient(to right, var(--tw-gradient-stops)); }

/* 渐变颜色 - 需要在元素上直接指定 */
.from-purple-500 { --tw-gradient-from: #8b5cf6; --tw-gradient-stops: var(--tw-gradient-from), var(--tw-gradient-to, rgba(139, 92, 246, 0)); }
.to-purple-700 { --tw-gradient-to: #7c3aed; }
.to-purple-800 { --tw-gradient-stops: var(--tw-gradient-from), #6b21a8; }

/* 文字颜色 */
.text-white { color: #ffffff; }
.text-white\/60 { color: rgba(255, 255, 255, 0.6); }
.text-white\/80 { color: rgba(255, 255, 255, 0.8); }
.text-white\/40 { color: rgba(255, 255, 255, 0.4); }
.text-gray-100 { color: #f3f4f6; }
.text-gray-800 { color: #1f2937; }
.text-gray-600 { color: #4b5563; }
.text-gray-500 { color: #6b7280; }
.text-gray-400 { color: #9ca3af; }

/* 宽度高度 */
.w-full { width: 100%; }
.w-20 { width: 5rem; }
.h-10 { height: 2.5rem; }
.h-11 { height: 2.75rem; }
.h-12 { height: 3rem; }
.h-20 { height: 5rem; }
.min-h-screen { min-height: 100vh; }

/* 边框 */
.border-b { border-bottom-width: 1px; border-bottom-style: solid; }
.border-gray-100 { border-color: #f3f4f6; }
.border-2 { border-width: 2px; border-style: solid; }

/* Overflow */
.overflow-hidden { overflow: hidden; }

/* 位置 */
.relative { position: relative; }
.absolute { position: absolute; }
.right-3 { right: 0.75rem; }
.bottom-2\.5 { bottom: 0.625rem; }

/* 其他 */
.cursor-pointer { cursor: pointer; }
.border { border-width: 1px; border-style: solid; }
</style>
