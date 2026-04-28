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

  /** 获取系统管理员列表（仅管理员可访问） */
  sysAdmins() {
    return api.post('/admin/sysadmins', {})
  },

  /** 获取用户列表（仅管理员可访问） */
  getUsers(params = {}) {
    return api.post('/admin/users', params)
  },

  /** 获取用户组列表（仅管理员可访问） */
  getGroups() {
    return api.post('/admin/groups', {})
  },

  /** 获取用户组成员列表（仅管理员可访问） */
  getGroupMembers(groupId, params = {}) {
    return api.post('/admin/group_members', { group_id: groupId, ...params })
  },

  /** 获取用户详细信息（仅管理员可访问） */
  getUserDetail(userId) {
    return api.post('/admin/user_detail', { user_id: userId })
  },

  /** 重置用户密码（仅管理员可访问） */
  resetUserPassword(userId, password) {
    return api.post('/admin/reset_user_password', { user_id: userId, password })
  },

  /** 获取登录失败日志（仅管理员可访问） */
  getLoginFailLogs(params = {}) {
    return api.post('/admin/login_fail_logs', params)
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
