import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useConfigStore } from './config'

// Storage Keys
const STORAGE_KEY_COOKIE = 'userCookie'
const STORAGE_KEY_COOKIE_SESSION = 'userCookieSession'

/**
 * 加载 JSON（永久存储）
 */
function loadJson(key) {
  try {
    const raw = uni.getStorageSync(key)
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

/**
 * 清除存储
 */
function removeStorage() {
  uni.removeStorageSync(STORAGE_KEY_COOKIE)
  uni.removeStorageSync(STORAGE_KEY_COOKIE_SESSION)
}

export const useUserStore = defineStore('user', () => {
  // ── State ──
  const userCookie = ref(null)   // Cookie 对象 { Value, Name, ExpiresAt, Remember }
  const user = ref(null)         // 用户基本信息 { ID, Name, ... }
  const userInfo = ref(null)    // 用户详细信息 { Username, FirstName, ... }

  // ── Getters ──
  /** Cookie 值字符串 */
  const cookieValue = computed(() => userCookie.value?.Value ?? '')

  /** 是否已登录 */
  const isLoggedIn = computed(() => !!userCookie.value)

  /** 用户名 */
  const username = computed(() => {
    return userInfo.value?.Username || user.value?.Name || ''
  })

  /** 头像 URL */
  const avatarUrl = computed(() => {
    if (userInfo.value?.AvatarPath) {
      const configStore = useConfigStore()
      return configStore.getApiBaseUrl() + '/static/avatar/' + userInfo.value.AvatarPath
    }
    return null
  })

  // ── Actions ──

  /**
   * 登录 - 保存 Cookie 并获取用户信息
   * @param {Object} cookie - 后端返回的 cookie 对象
   */
  function login(cookie) {
    userCookie.value = cookie

    // 持久化存储
    if (cookie.Remember) {
      uni.setStorageSync(STORAGE_KEY_COOKIE, JSON.stringify(cookie))
    }
    // 会话存储（始终保存）
    uni.setStorageSync(STORAGE_KEY_COOKIE_SESSION, JSON.stringify(cookie))

    // 检查 cookie 是否过期
    if (cookie.ExpiresAt && new Date(cookie.ExpiresAt) < new Date()) {
      logout()
      return
    }

    // 获取用户信息
    fetchUserInfo()
  }

  /**
   * 退出登录
   */
  function logout() {
    userCookie.value = null
    user.value = null
    userInfo.value = null
    removeStorage()
  }

  /**
   * 获取用户信息
   */
  async function fetchUserInfo() {
    try {
      const { api } = await import('../api/index.js')
      const res = await api.post('/users/getinfo', {})
      if (res.errCode === 0 && res.data) {
        user.value = res.data.user ?? null
        userInfo.value = res.data.userInfo ?? null
      }
    } catch (e) {
      console.error('获取用户信息失败', e)
    }
  }

  /**
   * 应用启动时恢复登录状态
   */
  function restoreSession() {
    // 优先使用会话存储，否则用永久存储
    let cookieStr = uni.getStorageSync(STORAGE_KEY_COOKIE_SESSION)
    if (!cookieStr) {
      cookieStr = uni.getStorageSync(STORAGE_KEY_COOKIE)
    }

    if (cookieStr) {
      try {
        const cookie = JSON.parse(cookieStr)
        // 检查是否过期
        if (cookie.ExpiresAt && new Date(cookie.ExpiresAt) < new Date()) {
          logout()
          return
        }
        // 直接设置状态并获取用户信息
        userCookie.value = cookie
        fetchUserInfo()
      } catch {
        logout()
      }
    }
  }

  return {
    // State
    userCookie,
    user,
    userInfo,
    // Getters
    cookieValue,
    isLoggedIn,
    username,
    avatarUrl,
    // Actions
    login,
    logout,
    fetchUserInfo,
    restoreSession
  }
})
