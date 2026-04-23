import { api } from './index'

export const workOrderApi = {
  /** 获取工单列表 */
  getList(params = {}) {
    return api.post('/work_order/list', params)
  },

  /** 获取工单数量统计 */
  getCount() {
    return api.post('/work_order/count', {})
  },

  /** 新增工单 */
  add(data) {
    return api.post('/work_order/add', data)
  },

  /** 编辑工单 */
  update(data) {
    return api.post('/work_order/update', data)
  },

  /** 获取工单详情（含图片、commits） */
  get(id) {
    return api.post('/work_order/get', { id })
  },

  /** 新增进度（状态变更） */
  commit(id, status, comment = '', photos = [], purchaseOrderIds = []) {
    return api.post('/work_order/commit', { id, status, comment, photos, purchaseOrderIds })
  },

  /** 删除工单 */
  delete(id) {
    return api.post('/work_order/delete', { id })
  },

  /** 搜索采购订单 */
  searchPurchaseOrders(search = '', limit = 5) {
    return api.post('/work_order/search_purchase_orders', { search, limit })
  },
}
