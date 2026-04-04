import { api } from './index'

export const usersApi = {
  
  getUserInfoFromUserID(UserID) {
    return api.get('/users/getuserinfo/'+UserID)
  },
}