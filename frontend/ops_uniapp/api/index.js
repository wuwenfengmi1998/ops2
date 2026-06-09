/**
 * API 请求封装
 * 基于 uni.request，统一处理 Cookie 认证和错误
 */

import { useConfigStore } from '../stores/config'

// Storage Keys
const STORAGE_KEY_COOKIE_SESSION = 'userCookieSession'

// 错误码定义
const ERR_COOKIE_EXPIRED = 'userCookieError'

/**
 * 获取 API 基础地址
 */
function getBaseUrl() {
  const configStore = useConfigStore()
  return configStore.getApiBaseUrl()
}

/**
 * 获取保存的 Cookie 值
 */
function getCookie() {
  try {
    // 优先从 storage 读取（避免 Pinia 初始化时序问题）
    const cookieStr = uni.getStorageSync(STORAGE_KEY_COOKIE_SESSION)
    if (cookieStr) {
      const cookie = JSON.parse(cookieStr)
      return cookie?.Value ?? ''
    }
  } catch {
    // ignore
  }
  return ''
}

/**
 * 请求封装 - 统一错误处理
 */
function request(options) {
  return new Promise((resolve, reject) => {
    const baseUrl = getBaseUrl()
    if (!baseUrl) {
      reject({ errCode: -1, message: 'API地址未配置' })
      return
    }

    // GET 请求不注入 cookie，POST 请求注入 cookie
    let data = options.data || {}
    if (options.method === 'POST') {
      const cookie = getCookie()
      // 后端期望格式: { data: {...业务数据}, userCookieValue: "xxx" }
      data = {
        data: data,  // 业务数据包装在 data 字段
        userCookieValue: cookie || undefined
      }
    }

    uni.request({
      url: baseUrl + options.path,
      method: options.method || 'GET',
      data,
      header: {
        'Content-Type': 'application/json'
      },
      timeout: options.timeout || 10000,
      success: (res) => {
        const data = res.data

        // 检查 Cookie 过期
        if (data?.err_code === ERR_COOKIE_EXPIRED || data?.err_code === -44) {
          // 清除存储的 Cookie
          uni.removeStorageSync(STORAGE_KEY_COOKIE_SESSION)
          uni.showToast({ title: '登录已过期，请重新登录', icon: 'none' })
          reject({ errCode: -44, message: '登录已过期' })
          return
        }

        if (res.statusCode === 200) {
          resolve({
            errCode: data?.err_code === 0 ? 0 : -1,
            data: data?.return || null,
            raw: data
          })
        } else {
          uni.showToast({ title: '服务异常', icon: 'none' })
          reject({ errCode: -1, message: '服务异常' })
        }
      },
      fail: (err) => {
        uni.showToast({ title: '网络错误', icon: 'none' })
        reject({ errCode: -2, message: '网络错误', detail: err })
      }
    })
  })
}

/**
 * API 接口汇总
 */
export const api = {
  /**
   * GET 请求（不需要认证）
   */
  get(path, data = {}) {
    return request({
      path,
      method: 'GET',
      data
    })
  },

  /**
   * POST JSON 请求（需要认证）
   */
  post(path, data = {}) {
    return request({
      path,
      method: 'POST',
      data  // 业务数据，会被包装成 { data, userCookieValue }
    })
  },

  /**
   * POST FormData（文件上传）
   */
  upload(path, fileData) {
    return new Promise((resolve, reject) => {
      const baseUrl = getBaseUrl()
      const cookie = getCookie()

      uni.uploadFile({
        url: baseUrl + path,
        filePath: fileData.uri,
        name: fileData.name || 'file',
        formData: {
          cookie: cookie || ''
        },
        success: (res) => {
          try {
            const data = JSON.parse(res.data)
            resolve({
              errCode: data?.err_code === 0 ? 0 : -1,
              data: data?.return || null,
              raw: data
            })
          } catch {
            reject({ errCode: -1, message: '解析响应失败' })
          }
        },
        fail: (err) => {
          uni.showToast({ title: '上传失败', icon: 'none' })
          reject({ errCode: -2, message: '上传失败', detail: err })
        }
      })
    })
  }
}

export default api
