//用于保存其他用户的信息

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { usersApi } from '@/api/users';

export const useUsersStore = defineStore('users', () => {
    const usersInfo =ref([]);

    function getUserFromUserID(userID){
        return usersInfo.value?.find(item => item.UserID === userID) ?? null
    }

    function getUsernameFromUserID(userID){
        const target = getUserFromUserID(userID)
        if (target) {
            return target.Username
        }
        usersApi.getUserInfoFromUserID(userID).then((r) => {
            if (r.errCode == 0 && r.raw.err_code == 0 && r.raw.return.userinfo) {
                usersInfo.value.push(r.raw.return.userinfo)
            }
        })
        return "..."
    }

    function getAvatarUrlFromUserID(userID) {
        const target = getUserFromUserID(userID)
        if (target?.AvatarPath) {
            return `/api/static/avatar/${target.AvatarPath}`
        }
        return `/ava.svg`
    }

    return{
        usersInfo,getUsernameFromUserID,getAvatarUrlFromUserID,
    }

})