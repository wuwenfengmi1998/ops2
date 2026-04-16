/**
 * 用户状态管理（轻量 Store，对标 PC 前端 src/stores/user.js）
 * 
 * uni-app 不支持 Pinia，这里用模块化单例替代，
 * 提供与 PC 前端完全一致的接口语义。
 */

import { authApi } from '../@api/auth.js'

const STORAGE_KEY_COOKIE = 'userCookie'
const STORAGE_KEY_USER   = 'userData'
const STORAGE_KEY_INFO   = 'userInfo'

// ── 内部状态 ──────────────────────────────────────────────
const _state = {
  user:       null,   // TabUser_ 基本信息
  userInfo:   null,   // TabUserInfo_ 详情
  userCookie: null,   // Cookie 对象
  isLoggedIn: false,
}

// ── 辅助方法 ──────────────────────────────────────────────
function _loadJson(key) {
  try {
    const raw = uni.getStorageSync(key)
    return raw ? JSON.parse(raw) : null
  } catch { return null }
}

function _saveJson(key, data) {
  uni.setStorageSync(key, JSON.stringify(data))
}

function _remove(key) {
  uni.removeStorageSync(key)
}

// ── 对外接口（与 PC 端 useUserStore 语义一致） ─────────────
export const userStore = {

  // ── Getters ──

  get isLoggedIn()  { return _state.isLoggedIn },
  get user()        { return _state.user },
  get userInfo()    { return _state.userInfo },
  get userCookie()  { return _state.userCookie },

  /** 获取 Cookie Value 字符串（供 request.js 自动注入） */
  getCookieValue() {
    return _state.userCookie?.Value ?? ''
  },

  /** 头像 URL */
  getAvatarUrl() {
    if (_state.userInfo?.AvatarPath) {
      return `/api/static/avatar/${_state.userInfo.AvatarPath}`
    }
    return '/static/ava.svg'
  },

  /** 生日（YYYY-MM-DD） */
  getBirthday() {
    if (!_state.userInfo?.Birthdate) return ''
    return String(_state.userInfo.Birthdate).substring(0, 10)
  },

  /** 用户名（显示名 > 账号名） */
  getDisplayName() {
    return _state.userInfo?.Username || _state.user?.Name || ''
  },

  // ── Actions ──

  /**
   * 登录成功后调用，保存 cookie 并拉取用户信息
   * @param {object} cookie - 后端返回的 Cookie 对象
   */
  login(cookie) {
    _state.userCookie = cookie
    _state.isLoggedIn = true
    _saveJson(STORAGE_KEY_COOKIE, cookie)
    // 验证是否过期
    if (cookie.ExpiresAt && new Date(cookie.ExpiresAt) < new Date()) {
      this.logout()
      return
    }
    this.fetchUserInfo()
  },

  /** 登出，清理所有本地状态 */
  logout() {
    _state.user       = null
    _state.userInfo   = null
    _state.userCookie = null
    _state.isLoggedIn = false
    _remove(STORAGE_KEY_COOKIE)
    _remove(STORAGE_KEY_USER)
    _remove(STORAGE_KEY_INFO)
  },

  /** 拉取用户信息并缓存 */
  async fetchUserInfo() {
    try {
      const { errCode, data } = await authApi.getUserInfo()
      if (errCode === 0) {
        _state.user     = data?.user     ?? null
        _state.userInfo = data?.userInfo ?? null
        if (_state.user)     _saveJson(STORAGE_KEY_USER, _state.user)
        if (_state.userInfo) _saveJson(STORAGE_KEY_INFO, _state.userInfo)
      }
    } catch { /* request.js 已处理提示 */ }
  },

  /**
   * 应用启动时恢复登录状态（在 App.vue 的 onLaunch 中调用）
   */
  restoreSession() {
    const cookie = _loadJson(STORAGE_KEY_COOKIE)
    if (cookie) {
      // 验证 cookie 是否还在有效期
      if (cookie.ExpiresAt && new Date(cookie.ExpiresAt) < new Date()) {
        this.logout()
        return
      }
      _state.userCookie = cookie
      _state.isLoggedIn = true
      // 恢复缓存的用户信息（避免冷启动白屏）
      _state.user     = _loadJson(STORAGE_KEY_USER)
      _state.userInfo = _loadJson(STORAGE_KEY_INFO)
      // 后台静默刷新
      this.fetchUserInfo()
    } else {
      this.logout()
    }
  },
}
