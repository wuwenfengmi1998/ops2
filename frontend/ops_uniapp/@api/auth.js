/**
 * 认证相关 API
 * 对标 PC 前端 src/api/auth.js
 */
import { request } from './request.js'

export const authApi = {
  /** 登录 */
  login(username, password, remember = false) {
    return request.post('/users/login', { username, password, remember })
  },

  /** 注册 */
  register(username, email, password) {
    return request.post('/users/register', { username, useremail: email, userpass: password })
  },

  /** 通过 cookie 获取用户信息 */
  getUserInfo() {
    return request.post('/users/getinfo', {})
  },

  /** 修改密码 */
  changePassword(oldPass, newPass) {
    return request.post('/users/changePassword', { oldpass: oldPass, newpass: newPass })
  },

  /** 修改邮箱 */
  changeEmail(newEmail) {
    return request.post('/users/changeEmail', { newemail: newEmail })
  },

  /** 修改用户信息 */
  updateInfo(data) {
    return request.post('/users/updateInfo', data)
  },

  /** 更新头像（文件上传） */
  updateAvatar(filePath) {
    return request.upload('/users/updateAvatar', filePath)
  },
}
