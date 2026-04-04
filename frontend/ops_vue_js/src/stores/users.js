//用于保存其他用户的信息

import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUsersStore = defineStore('users', () => {
    const usersInfo =ref([]);

    function getUsernameFromUserID(userID){
        
        return "123"
    }

    return{
        usersInfo,getUsernameFromUserID,
    }

})