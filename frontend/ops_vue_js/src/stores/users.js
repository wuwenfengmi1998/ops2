//用于保存其他用户的信息

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { usersApi } from '@/api/users';

const usersInfo = ref([]);
// 正在请求中的 userID 集合，避免重复发请求
const pendingFetch = new Set();

export const useUsersStore = defineStore('users', () => {

    function getUserFromUserID(userID) {
        return usersInfo.value?.find(item => item.UserID === userID) ?? null
    }

    function fetchUser(userID) {
        if (pendingFetch.has(userID)) return
        pendingFetch.add(userID)
        usersApi.getUserInfoFromUserID(userID).then((r) => {
            if (r.errCode == 0 && r.raw.err_code == 0 && r.raw.return?.userinfo) {
                // 防止并发写入重复数据
                if (!usersInfo.value.find(item => item.UserID === userID)) {
                    usersInfo.value.push(r.raw.return.userinfo)
                }
            }
        }).finally(() => {
            pendingFetch.delete(userID)
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
