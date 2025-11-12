// stores/user.js
import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { myfuncs } from '@/myfunc.js'

// 组合式 API 写法 (推荐)
export const useUserStore = defineStore("user", () => {
  // 状态 (State)
  const userInfo = ref(null);
  const cookieValue = ref("");
  const isLoggedIn = ref(false);

  const logout = () => {
    isLoggedIn.value = false;
  };
  const loginFromStoreCookie = () => {
    //从store获取cookie
    var cookie=myfuncs.loadJson("userCookie")
    console.log(cookie)
    //isLoggedIn.value = true;
  };
  const loginUpdata = (cookie) => {
    console.log(cookie)
    cookieValue.value=cookie.value
  };

  return {
    userInfo,
    cookieValue,
    isLoggedIn,
    logout,
    loginFromStoreCookie,
    loginUpdata,
  };
});
