import api from './index.js'

export const customerApi = {
  // 客户列表
  list(params = {}) {
    return api.post('/customer/list', params)
  },
  
  // 客户详情
  get(id) {
    return api.post('/customer/get', { id })
  },
  
  // 新增客户
  add(data) {
    return api.post('/customer/add', data)
  },
  
  // 编辑客户
  update(data) {
    return api.post('/customer/update', data)
  },
  
  // 删除客户
  delete(id) {
    return api.post('/customer/delete', { id })
  }
}
