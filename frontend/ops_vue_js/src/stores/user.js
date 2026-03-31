import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'

const STORAGE_KEY_COOKIE = 'userCookie'

function loadJson(key) {
  try {
    const raw = localStorage.getItem(key)
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

function loadJsonSession(key) {
  try {
    const raw = sessionStorage.getItem(key)
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

function saveJson(key, data) {
  localStorage.setItem(key, JSON.stringify(data))
}

function saveJsonSession(key, data) {
  sessionStorage.setItem(key, JSON.stringify(data))
}

function removeStorage(key) {
  localStorage.removeItem(key)
  sessionStorage.removeItem(key)
}

export const useUserStore = defineStore('user', () => {
  // ── State ──
  const user = ref(null)        // TabUser_ 基本信息
  const userInfo = ref(null)    // TabUserInfo_ 详情
  const userCookie = ref(null)  // Cookie session
  const isLoggedIn = ref(false)

  // ── Getters ──
  const cookieValue = computed(() => userCookie.value?.Value ?? '')

  const avatarUrl = computed(() => {
    if (userInfo.value?.AvatarPath) {
      return `/api/static/avatar/${userInfo.value.AvatarPath}`
    }
    return '/ava.svg'
  })

  const birthday = computed(() => {
    if (!userInfo.value?.Birthdate) return ''
    const d = new Date(userInfo.value.Birthdate)
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    return `${y}-${m}-${day}`
  })

  // ── Actions ──
  function login(cookie) {
    userCookie.value = cookie
    isLoggedIn.value = true
    // 保存 cookie
    saveJsonSession(STORAGE_KEY_COOKIE, cookie)
    if (cookie.Remember) {
      saveJson(STORAGE_KEY_COOKIE, cookie)
    }
    // 检查 cookie 是否过期
    if (cookie.ExpiresAt && new Date(cookie.ExpiresAt) < new Date()) {
      logout()
      return
    }
    // 获取用户信息
    fetchUserInfo()
  }

  function logout() {
    userCookie.value = null
    user.value = null
    userInfo.value = null
    isLoggedIn.value = false
    removeStorage(STORAGE_KEY_COOKIE)
  }

  async function fetchUserInfo() {
    try {
      const { errCode, data } = await authApi.getUserInfo()
      if (errCode === 0) {
        user.value = data.user ?? null
        userInfo.value = data.userInfo ?? null
      }
    } catch {
      // 拦截器已处理错误提示
    }
  }

  /** 应用启动时尝试从存储恢复登录状态 */
  function restoreSession() {
    let cookie = loadJsonSession(STORAGE_KEY_COOKIE)
    if (!cookie) {
      cookie = loadJson(STORAGE_KEY_COOKIE)
    }
    if (cookie) {
      login(cookie)
    } else {
      logout()
    }
  }

  return {
    user,
    userInfo,
    userCookie,
    isLoggedIn,
    cookieValue,
    avatarUrl,
    birthday,
    login,
    logout,
    fetchUserInfo,
    restoreSession,
  }
})
