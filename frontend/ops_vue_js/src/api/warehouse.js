import { api } from './index'

export const warehouseApi = {
  /** 获取容器列表 */
  getContainers(params = {}) {
    return api.post('/warehouse/list_container', params)
  },

  /** 获取容器详情（含图片） */
  getContainer(id) {
    return api.post('/warehouse/get_container', { id })
  },

  /** 新增容器 */
  addContainer(data) {
    return api.post('/warehouse/add_container', data)
  },

  /** 编辑容器 */
  updateContainer(data) {
    return api.post('/warehouse/update_container', data)
  },

  /** 删除容器 */
  deleteContainer(id) {
    return api.post('/warehouse/delete_container', { id })
  },

  /** 获取仓库统计 */
  getCount() {
    return api.post('/warehouse/count', {})
  },

  /** 获取物品列表 */
  getItems(params = {}) {
    return api.post('/warehouse/list_item', params)
  },

  /** 新增物品 */
  addItem(data) {
    return api.post('/warehouse/add_item', data)
  },

  /** 获取物品详情 */
  getItem(id) {
    return api.post('/warehouse/get_item', { id })
  },

  /** 编辑物品 */
  updateItem(data) {
    return api.post('/warehouse/update_item', data)
  },

  /** 删除物品 */
  deleteItem(id) {
    return api.post('/warehouse/delete_item', { id })
  },

  /** 移动物品 */
  moveItem(data) {
    return api.post('/warehouse/move_item', data)
  },

  /** 获取容器列表（用于移动目标选择） */
  getContainers(params = {}) {
    return api.post('/warehouse/list_container', params)
  },
}
