/**
 * 日程相关 API
 * 对标 PC 前端 src/api/schedule.js
 */
import { request } from './request.js'

export const scheduleApi = {
  /** 获取日程列表 */
  getEvents(params = {}) {
    return request.post('/schedule/getevents', params)
  },

  /** 新增日程 */
  addEvent(data) {
    return request.post('/schedule/addevent', data)
  },

  /** 编辑日程 */
  editEvent(data) {
    return request.post('/schedule/editevent', data)
  },

  /** 删除日程 */
  deleEvent(data) {
    return request.post('/schedule/deleevent', data)
  },
}
