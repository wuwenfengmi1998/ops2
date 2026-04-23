//用于保存其他用户的信息

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { usersApi } from '@/api/users';

const usersInfo = ref([]);

export const useUsersStore = defineStore('users', () => {

    function getUserFromUserID(userID) {
        return usersInfo.value?.find(item => item.UserID === userID) ?? null
    }

    function fetchUser(userID) {
        // 缓存命中则不再请求（依赖 usersInfo 的响应式）
        if (getUserFromUserID(userID)) return
        usersApi.getUserInfoFromUserID(userID).then((r) => {
            if (r.errCode == 0 && r.raw.err_code == 0 && r.raw.return?.userinfo) {
                // 防止并发写入重复数据
                if (!usersInfo.value.find(item => item.UserID === userID)) {
                    usersInfo.value.push(r.raw.return.userinfo)
                }
            }
        })
    }

    function getUsernameFromUserID(userID) {
        const target = getUserFromUserID(userID)
        if (target) {
            return target.Username
        }
        fetchUser(userID)
        return "..."
    }

    function getAvatarUrlFromUserID(userID) {
        const target = getUserFromUserID(userID)
        if (target?.AvatarPath) {
            return `/api/static/avatar/${target.AvatarPath}`
        }
        // 触发加载（如果还没加载过）
        fetchUser(userID)
        return `/ava.svg`
    }

    return {
        usersInfo, getUsernameFromUserID, getAvatarUrlFromUserID,
    }

})
