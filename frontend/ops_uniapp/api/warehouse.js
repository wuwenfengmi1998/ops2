import api from './index.js'

export const warehouseApi = {
  // 获取容器列表
  listContainer(params = {}) {
    return api.post('/warehouse/list_container', params)
  },
  
  // 获取容器详情
  getContainer(id) {
    return api.post('/warehouse/get_container', { id })
  },
  
  // 新增容器
  addContainer(data) {
    return api.post('/warehouse/add_container', data)
  },
  
  // 更新容器
  updateContainer(data) {
    return api.post('/warehouse/update_container', data)
  },
  
  // 删除容器
  deleteContainer(id) {
    return api.post('/warehouse/delete_container', { id })
  },
  
  // 获取物品列表
  listItem(params = {}) {
    return api.post('/warehouse/list_item', params)
  },
  
  // 获取物品详情
  getItem(id) {
    return api.post('/warehouse/get_item', { id })
  },
  
  // 新增物品
  addItem(data) {
    return api.post('/warehouse/add_item', data)
  },
  
  // 更新物品
  updateItem(data) {
    return api.post('/warehouse/update_item', data)
  },
  
  // 删除物品
  deleteItem(id) {
    return api.post('/warehouse/delete_item', { id })
  },
  
  // 移动物品
  moveItem(id, containerId) {
    return api.post('/warehouse/move_item', { 
      item_id: id, 
      new_container: containerId 
    })
  },
  
  // 获取仓库统计
  getStats() {
    return api.post('/warehouse/get_stats', {})
  }
}
