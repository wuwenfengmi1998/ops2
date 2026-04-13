import { api } from './index'

export const purchaseApi = {
  /** 获取采购订单列表 */
  getOrders(params = {}) {
    return api.post('/purchase/getorders', params)
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
}
