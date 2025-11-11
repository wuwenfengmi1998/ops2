// stores/user.js
import { defineStore } from "pinia";
import { ref, computed } from "vue";

// 组合式 API 写法 (推荐)
export const useUserStore = defineStore("user", () => {
  // 状态 (State)
  const userInfo = ref(null);
  const token = ref("");
  const isLoggedIn = ref(false);

  const logout = () => {
    isLoggedIn.value = false;
  };
  const login = () => {
    isLoggedIn.value = true;
  };

  return {
    userInfo,
    token,
    isLoggedIn,
    logout,
    login,
  };
});
