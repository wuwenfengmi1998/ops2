import { api } from './index'

export const authApi = {
  /** 登录 */
  login(username, password, remember) {
    return api.post('/users/login', { username, password, remember })
  },

  /** 注册 */
  register(username, email, password) {
    return api.post('/users/register', { username, useremail: email, userpass: password })
  },

  /** 通过 cookie 获取用户信息 */
  getUserInfo() {
    return api.post('/users/getinfo', {})
  },

  /** 修改密码 */
  changePassword(oldPass, newPass) {
    return api.post('/users/changePassword', { oldpass: oldPass, newpass: newPass })
  },

  /** 修改邮箱 */
  changeEmail(newEmail) {
    return api.post('/users/changeEmail', { newemail: newEmail })
  },

  /** 修改用户信息 */
  updateInfo(data) {
    return api.post('/users/updateInfo', data)
  },

  /** 更新头像 */
  updateAvatar(file) {
    return api.upload('/users/updateAvatar', file)
  },
}
