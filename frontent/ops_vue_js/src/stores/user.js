// stores/user.js
import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { myfuncs } from "@/myfunc.js";

// 组合式 API 写法 (推荐)
export const useUserStore = defineStore("user", () => {
  // 状态 (State)
  const userInfo = ref(null);
  const userCookie = ref(null);
  const isLoggedIn = ref(false);

  const cookiesQualified = () => {
    //返回一个合格的cookie 就是没过期的cookie
    //如果cookie没过期直接返回，如果过期 顺便logout
    var cookieTimeout = userCookie.value.ExpiresAt;
    if (new Date(cookieTimeout) < new Date()) {
      //过期了
      logout();
    }
    return userCookie.value;
  };
  const logout = () => {
    userCookie.value = null;
    isLoggedIn.value = false;
    myfuncs.dele("userCookie");
    myfuncs.deleT("userCookie");
  };
  const login = (cookie) => {
    userCookie.value = cookie;
    isLoggedIn.value = true;
    //这里应该判读cookie的实效性
    userCookie.value = cookiesQualified();
  };

  const loginFromStoreCookie = () => {
    //从store获取cookie

    var cookie = myfuncs.loadJsonT("userCookie");
    if (cookie) {
      login(cookie);
    } else {
      cookie = myfuncs.loadJson("userCookie");
      if (cookie) {
        login(cookie);
      } else {
        logout();
      }
    }
  };


  return {
    userInfo,
    userCookie,
    isLoggedIn,
    logout,
    login,
    loginFromStoreCookie,
    
  };
});
