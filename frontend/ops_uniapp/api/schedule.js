/**
 * 日程管理 API
 */
import api from './index.js'

export const scheduleApi = {
  // 获取日程列表
  getEvents(data = {}) {
    return api.post('/schedule/getevents', data)
  },

  // 新增日程
  addEvent(data = {}) {
    return api.post('/schedule/addevent', data)
  },

  // 编辑日程
  editEvent(data = {}) {
    return api.post('/schedule/editevent', data)
  },

  // 删除日程
  deleteEvent(data = {}) {
    return api.post('/schedule/deleevent', data)
  }
}

export default scheduleApi
