// stores/user.js
import { defineStore } from "pinia";
import { ref, computed } from "vue";

// 组合式 API 写法 (推荐)
export const useUserStore = defineStore("user", () => {
  // 状态 (State)
  const userInfo = ref(null);
  const cookieValue = ref("");
  const isLoggedIn = ref(false);

  const logout = () => {
    isLoggedIn.value = false;
  };
  const login = () => {
    isLoggedIn.value = true;
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
    login,
    loginUpdata,
  };
});
