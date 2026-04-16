/**
 * 采购订单相关 API
 * 对标 PC 前端 src/api/purchase.js
 */
import { request } from './request.js'

export const purchaseApi = {
  /** 获取订单列表（分页+搜索+状态过滤） */
  getOrders(params = {}) {
    return request.post('/purchase/getorders', params)
  },

  /** 获取单个订单详情 */
  getOrder(id) {
    return request.post('/purchase/getorder', { id })
  },

  /** 获取各状态订单数量统计 */
  getOrderCount() {
    return request.post('/purchase/getordercount', {})
  },

  /** 新增订单 */
  addOrder(data) {
    return request.post('/purchase/addorder', data)
  },

  /** 更新订单 */
  updateOrder(data) {
    return request.post('/purchase/updateorder', data)
  },

  /** 更新订单状态 */
  updateStatus(data) {
    return request.post('/purchase/updatestatus', data)
  },

  /** 删除订单 */
  deleteOrder(id) {
    return request.post('/purchase/deleteorder', { id })
  },
}
