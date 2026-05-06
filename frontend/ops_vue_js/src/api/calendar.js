import { api } from './index'

export const calendarApi = {
  // 创建日历
  createCalendar(data) {
    return api.post('/calendar/calendar/create', data)
  },

  // 获取日历列表
  getCalendars() {
    return api.post('/calendar/calendar/list', {})
  },

  // 更新日历
  updateCalendar(data) {
    return api.post('/calendar/calendar/update', data)
  },

  // 删除日历
  deleteCalendar(id) {
    return api.post('/calendar/calendar/delete', { id })
  },

  // 获取日历事件
  getEvents(data) {
    return api.post('/calendar/calendar/events', data)
  },

  // 添加日历事件
  addEvent(data) {
    return api.post('/calendar/calendar/addevent', data)
  },

  // 更新日历事件
  updateEvent(data) {
    return api.post('/calendar/calendar/updateevent', data)
  },

  // 删除日历事件
  deleteEvent(id) {
    return api.post('/calendar/calendar/deleteevent', { id })
  }
}
