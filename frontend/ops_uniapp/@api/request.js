/**
 * 统一网络请求封装
 * 对标 PC 前端 src/api/index.js 的设计
 * - 自动注入 cookie
 * - 统一解析 err_code / return 字段
 * - Cookie 过期自动清理并跳转登录
 */

import { userStore } from '../store/user.js'

const API_BASE = '/api'

/**
 * 底层 POST 请求
 * @param {string} path  - 接口路径，如 /users/login
 * @param {object} data  - 业务数据（会被包在 data 字段下）
 * @returns {Promise<{errCode, data, raw}>}
 */
function post(path, data = {}) {
  const body = { data }

  // 自动注入 cookie
  const cookieValue = userStore.getCookieValue()
  if (cookieValue) {
    body.userCookieValue = cookieValue
  }

  return new Promise((resolve, reject) => {
    uni.request({
      url: API_BASE + path,
      method: 'POST',
      header: { 'Content-Type': 'application/json' },
      data: body,
      timeout: 15000,
      success(res) {
        const raw = res.data
        const errCode = raw?.err_code ?? -1

        // Cookie 过期（err_code === -44），自动登出并跳转登录
        if (errCode === -44) {
          userStore.logout()
          uni.reLaunch({ url: '/pages/signin' })
          reject(new Error('Cookie expired'))
          return
        }

        resolve({
          errCode,
          data: raw?.return ?? null,
          raw,
        })
      },
      fail(err) {
        uni.showToast({ title: '网络连接失败', icon: 'none' })
        reject(err)
      },
    })
  })
}

/**
 * 上传文件（FormData）
 * @param {string} path    - 接口路径
 * @param {string} filePath - 本地文件路径（uni.chooseImage 返回的 tempFilePaths[0]）
 * @param {string} name    - 文件字段名，默认 'file'
 */
function upload(path, filePath, name = 'file') {
  const cookieValue = userStore.getCookieValue()

  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: API_BASE + path,
      filePath,
      name,
      formData: cookieValue ? { cookie: cookieValue } : {},
      timeout: 30000,
      success(res) {
        try {
          const raw = JSON.parse(res.data)
          resolve({
            errCode: raw?.err_code ?? -1,
            data: raw?.return ?? null,
            raw,
          })
        } catch {
          reject(new Error('JSON parse error'))
        }
      },
      fail(err) {
        uni.showToast({ title: '上传失败', icon: 'none' })
        reject(err)
      },
    })
  })
}

export const request = { post, upload }
