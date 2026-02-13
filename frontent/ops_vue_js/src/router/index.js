import {
  createRouter,
  createWebHistory,
  createWebHashHistory,
} from "vue-router";
import HomeView from "../views/HomeView.vue";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
    },
    
    {
      path: "/settings/account",
      name: "settings account",
      component: () => import("@/views/settings/account.vue"),
    },
    {
      path: "/settings/contact",
      name: "settings contact",
      component: () => import("@/views/settings/contact.vue"),
    },
    {
      path: "/settings/security",
      name: "settings security",
      component: () => import("@/views/settings/security.vue"),
    },
    {
      path: "/about",
      name: "about",
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import("@/views/AboutView.vue"),
    },
    {
      path: "/test",
      name: "test",
      component: () => import("@/views/test.vue"),
    },
    {
      path: "/login",
      name: "login",
      component: () => import("@/views/loginView.vue"),
    },
    {
      path: "/forgot_password",
      name: "forgot password",
      component: () => import("@/views/forgotPassword.vue"),
    },
    {
      path: "/register",
      name: "Register",
      component: () => import("@/views/registerView.vue"),
    },
    {
      path: "/admin",
      name: "admin",
      component: () => import("@/views/adminView.vue"),
    },

    {
      path: "/schedule",
      name: "schedule",
      component: () => import("@/views/scheduleView.vue"),
    },
    {
      path: "/purchase",
      name: "purchase",
      component: () => import("@/views/purchase/purchase.vue"),
    },
    {
      path: "/purchase/addorder",
      name: "purchase/addorder",
      component: () => import("@/views/purchase/addorder.vue"),
    },
        {
      path: "/purchase/showorder/:id",
      name: "purchase/showorder",
      component: () => import("@/views/purchase/showorder.vue"),
    },
    {
      path: "/warehouse",
      name: "warehouse",
      component: () => import("@/views/warehouse.vue"),
    },
    {
      path: "/404",
      name: "404",
      component: () => import("@/views/404.vue"),
    },
    // 404 页面 - 放在最后
    {
      path: "/:pathMatch(.*)*", // 通配符，匹配所有路由
      name: "NotFound",
      component: () => import("@/views/404.vue"),
    },
  ],
});

export default router;
