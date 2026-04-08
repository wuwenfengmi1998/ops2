import { createRouter, createWebHashHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    // ── 需要登录的页面（带 DefaultLayout） ──
    {
      path: '/',
      component: () => import('@/layouts/DefaultLayout.vue'),
      children: [
        {
          path: '',
          name: 'home',
          component: () => import('@/views/HomeView.vue'),
        },
        {
          path: 'settings/account',
          name: 'settings-account',
          component: () => import('@/views/settings/AccountView.vue'),
        },
        {
          path: 'settings/contact',
          name: 'settings-contact',
          component: () => import('@/views/settings/ContactView.vue'),
        },
        {
          path: 'settings/security',
          name: 'settings-security',
          component: () => import('@/views/settings/SecurityView.vue'),
        },
        {
          path: 'schedule',
          name: 'schedule',
          component: () => import('@/views/ScheduleView.vue'),
        },
        {
          path: 'purchase',
          name: 'purchase',
          component: () => import('@/views/purchase/PurchaseList.vue'),
        },
        {
          path: 'purchase/addorder',
          name: 'purchase-add',
          component: () => import('@/views/purchase/addorder.vue'),
        },
        {
          path: 'purchase/showorder/:id',
          name: 'purchase-show',
          component: () => import('@/views/purchase/ShowOrder.vue'),
        },
        {
          path: 'warehouse',
          name: 'warehouse',
          component: () => import('@/views/WarehouseView.vue'),
        },
        {
          path: 'admin',
          name: 'admin',
          component: () => import('@/views/AdminView.vue'),
        },
      ],
    },

    // ── 认证页面（AuthLayout，全屏居中） ──
    {
      path: '/login',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        {
          path: '',
          name: 'login',
          component: () => import('@/views/LoginView.vue'),
        },
      ],
    },
    {
      path: '/register',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        {
          path: '',
          name: 'register',
          component: () => import('@/views/RegisterView.vue'),
        },
      ],
    },
    {
      path: '/forgot_password',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        {
          path: '',
          name: 'forgot-password',
          component: () => import('@/views/ForgotPasswordView.vue'),
        },
      ],
    },

    // ── 404 ──
    {
      path: '/404',
      name: '404',
      component: () => import('@/views/NotFoundView.vue'),
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/404',
    },
  ],
})

// ── 全局前置守卫 ──
router.beforeEach((to) => {
  const userStore = useUserStore()

  // 不需要登录的页面
  const publicPages = ['/', '/login', '/register', '/forgot_password', '/schedule','/warehouse', '/404']
  if (publicPages.includes(to.path)) return true

  // 未登录 → 跳转登录
  if (!userStore.isLoggedIn) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  return true
})

export default router
