//用于保存其他用户的信息

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { usersApi } from '@/api/users';

export const useUsersStore = defineStore('users', () => {
    const usersInfo =ref([]);

    function getUsernameFromUserID(userID){
        //console.log(userID)
        //先在usersInfo找找有没有
        const target = usersInfo.value?.find(item => item.UserID === userID)
        if(target){
            return target.Username //有的话直接返回
        }else{
            //没有的话 询问后端
            usersApi.getUserInfoFromUserID(userID).then((r)=>{
                //console.log(r)
                if(r.errCode==0)
                {
                    switch(r.raw.err_code){
                        case 0:
                            if(r.raw.return.userinfo){
                                usersInfo.value.push(r.raw.return.userinfo)
                            }
                            break;
                    }
                }
            })

            return "..." // 第一次返回这个，不会空白/报错
        }        
    }

    return{
        usersInfo,getUsernameFromUserID,
    }

})