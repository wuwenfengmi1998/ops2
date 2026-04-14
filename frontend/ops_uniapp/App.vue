<script>
export default {
  onLaunch: function() {
    console.log('App Launch')
    
    // 初始化全局 API 地址
    const DEFAULT_API_URL = 'http://192.168.13.105/api/'
    const apiUrl = uni.getStorageSync('apiUrl') || DEFAULT_API_URL
    if (!getApp().globalData) {
      getApp().globalData = {}
    }
    getApp().globalData.BASE_URL = apiUrl
    
    this.checkLoginStatus()
  },
  onShow: function() {
    console.log('App Show')
  },
  onHide: function() {
    console.log('App Hide')
  },
  methods: {
    checkLoginStatus() {
      // 检查登录状态
      const sessionCookie = uni.getStorageSync('sessionCookie')
      if (!sessionCookie) {
        // 未登录，跳转到登录页
        uni.reLaunch({
          url: '/pages/login/login'
        })
      }
    }
  }
}
</script>

<style>
/*每个页面公共css */
page {
  background-color: #f5f5f5;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
</style>
