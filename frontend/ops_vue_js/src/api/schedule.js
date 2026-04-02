import { api } from './index'

export const scheduleApi = {
  
  getEvents(params = {}) {
    return api.post('/schedule/getevents', params)
  },

  
  addEvent(data) {
    return api.post('/schedule/addevent', data)
  },
}