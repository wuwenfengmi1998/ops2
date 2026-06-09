/**
 * 采购订单 API
 */
import api from './index.js'

export const purchaseApi = {
  // 获取订单列表
  getOrders(data = {}) {
    return api.post('/purchase/getorders', data)
  },

  // 获取订单详情
  getOrder(data = {}) {
    return api.post('/purchase/getorder', data)
  },

  // 获取订单数量统计
  getOrderCount() {
    return api.post('/purchase/getordercount', {})
  },

  // 更新订单状态
  updateOrderStatus(data = {}) {
    return api.post('/purchase/updatestatus', data)
  },

  // 新增订单
  addOrder(data = {}) {
    return api.post('/purchase/addorder', data)
  },

  // 更新订单
  updateOrder(data = {}) {
    return api.post('/purchase/updateorder', data)
  }
}

export default purchaseApi
