import { api } from './index'

export const scheduleApi = {
  
  getEvents(params = {}) {
    return api.post('/schedule/getevents', params)
  },

  
  addEvent(data) {
    return api.post('/schedule/addevent', data)
  },

  editEvent(data) {
    return api.post('/schedule/editevent', data)
  },
  deleEvent(data) {
    return api.post('/schedule/deleevent', data)
  },
}