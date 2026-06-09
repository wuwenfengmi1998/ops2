import api from './index.js'

export const workOrderApi = {
  // 获取工单列表
  list(params = {}) {
    return api.post('/work_order/list', params)
  },
  
  // 获取工单详情
  get(id) {
    return api.post('/work_order/get', { id })
  },
  
  // 获取工单统计
  getCount() {
    return api.post('/work_order/count', {})
  },
  
  // 新增工单
  add(data) {
    return api.post('/work_order/add', data)
  },
  
  // 更新工单
  update(data) {
    return api.post('/work_order/update', data)
  },
  
  // 删除工单
  delete(id) {
    return api.post('/work_order/delete', { id })
  },
  
  // 提交进度
  commit(data) {
    return api.post('/work_order/commit', data)
  },

  // 搜索采购订单
  searchPurchaseOrders(query, limit = 10) {
    return api.post('/work_order/search_purchase_orders', { search: query, limit })
  }
}
