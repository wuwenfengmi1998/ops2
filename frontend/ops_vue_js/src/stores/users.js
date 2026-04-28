//用于保存其他用户的信息

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { usersApi } from '@/api/users';

const usersInfo = ref([]);
// 正在请求中的 promiseMap，同一 userID 只发一次请求
const inflightRequests = new Map();

export const useUsersStore = defineStore('users', () => {

    function getUserFromUserID(userID) {
        return usersInfo.value?.find(item => item.UserID === userID) ?? null
    }

    function fetchUser(userID) {
        // 已缓存则直接返回
        if (getUserFromUserID(userID)) return

        // 同一请求已在飞中，等待它完成后再更新缓存引用
        if (inflightRequests.has(userID)) return

        // 立即占位：同步标记为"请求中"，同一帧内后续调用能命中
        const placeholder = { UserID: userID, _loading: true }
        usersInfo.value.push(placeholder)

        const promise = usersApi.getUserInfoFromUserID(userID).then((r) => {
            if (r.errCode == 0 && r.raw.err_code == 0 && r.raw.return?.userinfo) {
                const info = r.raw.return.userinfo
                // 替换占位对象为真实数据
                const idx = usersInfo.value.findIndex(item => item.UserID === userID)
                if (idx !== -1) {
                    usersInfo.value.splice(idx, 1, info)
                }
            } else {
                // 请求失败，移除占位
                usersInfo.value = usersInfo.value.filter(item => item.UserID !== userID)
            }
        }).finally(() => {
            inflightRequests.delete(userID)
        })

        inflightRequests.set(userID, promise)
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
        usersInfo, getUsernameFromUserID, getAvatarUrlFromUserID, fetchUser,
    }

})
