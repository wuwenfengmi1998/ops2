/**
 * 用户信息缓存
 */
import { ref } from 'vue'
import { usersApi } from '../api/users'

// 全局缓存
const usersInfo = ref({})
// 请求中的 promiseMap，防止重复请求
const inflightRequests = {}

/**
 * 根据用户ID获取用户信息（带缓存）
 * @param {number} userID
 * @returns {Promise<{Username, UserEmail} | null>}
 */
export async function fetchUserInfo(userID) {
  if (!userID) return null

  // 已有缓存
  if (usersInfo.value[userID]) {
    return usersInfo.value[userID]
  }

  // 请求中，等待完成
  if (inflightRequests[userID]) {
    return inflightRequests[userID]
  }

  // 发起请求
  const promise = usersApi.getUserInfoFromUserID(userID).then((res) => {
    if (res.errCode === 0 && res.raw?.return?.userinfo) {
      const info = res.raw.return.userinfo
      usersInfo.value[userID] = info
      return info
    }
    return null
  }).finally(() => {
    delete inflightRequests[userID]
  })

  inflightRequests[userID] = promise
  return promise
}

/**
 * 根据用户ID获取用户名（同步，需先调用 fetchUserInfo）
 * @param {number} userID
 * @returns {string}
 */
export function getUsername(userID) {
  if (!userID) return ''
  const user = usersInfo.value[userID]
  return user?.Username || `用户${userID}`
}

export default {
  fetchUserInfo,
  getUsername
}
