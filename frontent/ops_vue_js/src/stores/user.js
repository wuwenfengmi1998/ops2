// stores/user.js
import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { myfuncs } from "@/myfunc.js";
import { my_network_func } from "@/my_network_func";

// 组合式 API 写法 (推荐)
export const useUserStore = defineStore("user", () => {
  // 状态 (State)
  const userInfo = ref(null);
  const user =ref(null)
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

  const getUserInfoFromCookie = () => {
    my_network_func.postJson("/users/getinfo", {}, (r) => {
      //console.log(r);
      switch (r.statusCode) {
        case 200:
          switch (r.data.err_code) {
            case 0:
              user.value=r.data.return.user
              if(r.data.return.userInfo){
                userInfo.value=r.data.return.userInfo
              }else{
                userInfo.value=null
              }
              break;
            default:
              break;
          }
          break;
        default:
          break;
      }
    });
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
    //到这里cookie应该是有效的，尝试获取用户info,因为有的info可能是隐藏的 所以用post携带当前cookie去请求用户info
    getUserInfoFromCookie();
  };

  const cookieUpdata = (cookie) => {
    userCookie.value = cookie;
    myfuncs.saveJsonT("userCookie", cookie);
    if (cookie.Remember) {
      //长期保存cookie
      myfuncs.saveJson("userCookie", cookie);
    }
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
    user,
    userInfo,
    userCookie,
    isLoggedIn,
    logout,
    login,
    loginFromStoreCookie,
    cookieUpdata,
  };
});
