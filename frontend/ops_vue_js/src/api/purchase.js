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
}
