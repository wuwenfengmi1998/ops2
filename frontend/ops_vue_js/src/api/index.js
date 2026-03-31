import axios from 'axios'
import { useUserStore } from '@/stores/user'
import { useToastStore } from '@/stores/toast'

const API_BASE = '/api'

// 通用错误码
const ERR_COOKIE_EXPIRED = -44

/**
 * 创建 axios 实例，统一拦截
 */
const http = axios.create({
  baseURL: API_BASE,
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' },
})

// 请求拦截器：自动注入 userCookieValue
http.interceptors.request.use((config) => {
  const userStore = useUserStore()
  if (userStore.cookieValue) {
    if (config.data instanceof FormData) {
      config.data.append('cookie', userStore.cookieValue)
    } else if (typeof config.data === 'object') {
      config.data.userCookieValue = userStore.cookieValue
    }
  }
  return config
})

// 响应拦截器：统一处理 err_code 和错误
http.interceptors.response.use(
  (response) => {
    const data = response.data

    // Cookie 过期，自动登出
    if (data?.err_code === ERR_COOKIE_EXPIRED) {
      const userStore = useUserStore()
      userStore.logout()
      const toast = useToastStore()
      toast.warning('登录已过期，请重新登录')
      // 这里返回一个 rejected promise，让调用方知道请求失败了
      return Promise.reject(new Error('Cookie expired'))
    }

    return response
  },
  (error) => {
    const toast = useToastStore()

    if (error.message === 'Cookie expired') {
      // 已在上面处理
      return Promise.reject(error)
    }

    if (!error.response) {
      toast.error('网络错误')
    } else {
      toast.error('服务端错误')
    }

    return Promise.reject(error)
  },
)

/**
 * 封装请求方法：返回 { errCode, data } 格式
 * 成功时 errCode === 0，data 为服务端 return 字段
 * 失败时抛出异常
 */
function unwrapResponse(response) {
  const body = response.data
  return {
    errCode: body.err_code ?? -1,
    data: body.return ?? null,
    raw: body,
  }
}

export const api = {
  /**
   * GET 请求（一般不需要认证）
   */
  async get(path) {
    const res = await http.get(path)
    return unwrapResponse(res)
  },

  /**
   * POST JSON
   */
  async post(path, data = {}) {
    const res = await http.post(path, { data })
    return unwrapResponse(res)
  },

  /**
   * POST FormData（文件上传）
   */
  async upload(path, file) {
    const formData = new FormData()
    formData.append('file', file)
    const res = await http.post(path, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    return unwrapResponse(res)
  },
}

export default api
