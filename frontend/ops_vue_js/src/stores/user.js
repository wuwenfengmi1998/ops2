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
  const groups = ref([])        // 用户加入的群组列表
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

  // 是否系统管理员（后端直接返回）
  const isSysAdmin = ref(false)

  // 是否为日历管理员（在 calendar_admin 群组中）
  const isCalendarAdmin = computed(() =>
    groups.value.some(g => g.name === 'calendar_admin')
  )

  // 用户加入的群组名称列表（计算属性）
  const groupNames = computed(() => groups.value.map(g => g.name))

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
    isSysAdmin.value = false
    groups.value = []
    isLoggedIn.value = false
    removeStorage(STORAGE_KEY_COOKIE)
  }

  async function fetchUserInfo() {
    try {
      const { errCode, data } = await authApi.getUserInfo()
      if (errCode === 0) {
        user.value = data.user ?? null
        userInfo.value = data.userInfo ?? null
        // 存储系统管理员状态
        isSysAdmin.value = data.isSysAdmin === true
        // 存储用户群组列表
        groups.value = data.groups ?? []
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
    isSysAdmin,
    isCalendarAdmin,
    groups,
    groupNames,
    cookieValue,
    avatarUrl,
    birthday,
    login,
    logout,
    fetchUserInfo,
    restoreSession,
  }
})
