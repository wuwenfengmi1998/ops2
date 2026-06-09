/**
 * 用户相关 API
 */
import api from './index'

export const userApi = {
  /**
   * 用户登录
   * @param {string} username - 用户名
   * @param {string} password - 密码
   * @param {boolean} remember - 是否记住登录
   * @returns {Promise<{errCode, data: {cookie}}>}
   */
  login(username, password, remember = true) {
    return api.post('/users/login', {
      username,
      password,
      remember
    })
  },

  /**
   * 用户注册
   * @param {string} username - 用户名
   * @param {string} email - 邮箱
   * @param {string} password - 密码
   */
  register(username, email, password) {
    return api.post('/users/register', {
      username,
      useremail: email,
      userpass: password
    })
  },

  /**
   * 获取当前用户信息
   * @returns {Promise<{errCode, data: {user, userInfo}}>}
   */
  getUserInfo() {
    return api.post('/users/getinfo', {})
  },

  /**
   * 修改密码
   * @param {string} oldPass - 旧密码
   * @param {string} newPass - 新密码
   */
  changePassword(oldPass, newPass) {
    return api.post('/users/changePassword', {
      oldpass: oldPass,
      newpass: newPass
    })
  },

  /**
   * 修改邮箱
   * @param {string} newEmail - 新邮箱
   */
  changeEmail(newEmail) {
    return api.post('/users/changeEmail', {
      newemail: newEmail
    })
  },

  /**
   * 更新用户信息
   * @param {Object} data - 用户信息 { username, remark, birthday }
   */
  updateInfo(data) {
    return api.post('/users/updateInfo', data)
  },

  /**
   * 上传头像
   * @param {string} fileUri - 文件本地路径
   */
  updateAvatar(fileUri) {
    return api.upload('/users/updateAvatar', {
      name: 'file',
      uri: fileUri
    })
  }
}

export default userApi
