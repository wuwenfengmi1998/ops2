import api from './index'

export const customerApi = {
  // 新增客户
  add(data) {
    return api.post('/customer/add', data)
  },

  // 编辑客户
  update(data) {
    return api.post('/customer/update', data)
  },

  // 删除客户
  delete(data) {
    return api.post('/customer/delete', data)
  },

  // 客户列表
  list(data) {
    return api.post('/customer/list', data)
  },

  // 客户详情
  get(data) {
    return api.post('/customer/get', data)
  },
}
