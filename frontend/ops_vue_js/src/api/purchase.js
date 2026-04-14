import { api } from './index'

export const purchaseApi = {
  /** 获取采购订单列表 */
  getOrders(params = {}) {
    return api.post('/purchase/getorders', params)
  },

  /** 获取订单数量统计 */
  getOrderCount() {
    return api.post('/purchase/getordercount', {})
  },

  /** 新增采购订单 */
  addOrder(data) {
    return api.post('/purchase/addorder', data)
  },

  /** 获取单个订单详情（包含费用明细、图片、状态变更记录） */
  getOrder(id) {
    return api.post('/purchase/getorder', { id })
  },

  /** 更新订单状态（可附带评论和图片） */
  updateOrderStatus(id, status, comment = '', photos = []) {
    return api.post('/purchase/updatestatus', { id, status, comment, photos })
  },

  /** 编辑订单基本信息（标题/备注/链接/款式/费用/图片） */
  updateOrder(id, data) {
    return api.post('/purchase/updateorder', { id, ...data })
  },

  /** 删除订单 */
  deleteOrder(id) {
    return api.post('/purchase/deleteorder', { id })
  },
}
