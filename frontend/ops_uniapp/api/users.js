/**
 * 用户信息 API
 */
import api from './index.js'

export const usersApi = {
  // 通过用户ID获取用户信息（GET 请求）
  getUserInfoFromUserID(userID) {
    return api.get('/users/getuserinfo/' + userID)
  }
}

export default usersApi
