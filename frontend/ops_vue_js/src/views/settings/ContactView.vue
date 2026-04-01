<script setup>
// Vue 核心响应式 API
import { reactive, ref ,onMounted} from 'vue'
// 国际化 hook，用于翻译文本
import { useI18n } from 'vue-i18n'
// 用户状态管理
import { useUserStore } from '@/stores/user'
// Toast 通知状态管理
import { useToastStore } from '@/stores/toast'
// 页面标题 composable
import { usePageTitle } from '@/composables/usePageTitle'
// 表单验证 composable 和邮箱验证函数
import { useValidation, isValidEmail } from '@/composables'
// 认证 API
import { authApi } from '@/api/auth'
// 设置页面导航侧边栏组件
import SettingNav from '@/components/SettingNav.vue'

// 设置页面标题
usePageTitle('settings.contact_information')

// 初始化国际化翻译函数
const { t } = useI18n()

// 获取用户 store 实例
const userStore = useUserStore()

// 获取 toast store 实例，用于显示操作结果通知
const toast = useToastStore()

// 解构验证函数和错误状态
const { validate, errors, clearErrors } = useValidation()

// 表单数据：仅包含邮箱字段
const form = reactive({ email: '' })

// 保存按钮的加载状态
const loading = ref(false)

/**
 * 修改邮箱地址
 * 验证邮箱格式后调用 API 更新
 */
async function handleChangeEmail() {
  // 清空之前的验证错误
  clearErrors()

  // 验证邮箱：必填 + 格式校验
  const err = validate('email', form.email, t('message.please_enter_your_email'), isValidEmail)
  // 验证失败则阻止提交
  if (!err) return

  // 设置加载状态
  loading.value = true
  try {
    // 调用 API 修改邮箱
    const { errCode } = await authApi.changeEmail(form.email)

    // 根据错误码处理不同结果
    switch (errCode) {
      case 0:
        // 成功：显示成功提示并刷新用户信息
        toast.success(t('message.change_ok'))
        await userStore.fetchUserInfo()
        break
      case -43:
        // 邮箱格式无效：在输入框显示提示并弹出错误
        //form.email = t('message.this_not_email')
        toast.error(t('message.this_not_email'))
        break
      default:
        // 其他错误：显示通用服务器错误
        toast.error(t('message.server_error'))
    }
  } catch {
    // 网络错误等异常已被 axios 拦截器统一处理
  } finally {
    // 关闭加载状态
    loading.value = false
  }
}

onMounted(()=>{
form.email=userStore.user.Email
})

</script>

<template>
  
  <!-- 页面容器，最大宽度限制并居中 -->
  <div class="mx-auto max-w-5xl px-6 py-6">
    <!-- 页面标题 -->
    <h2 class="mb-6 text-2xl font-bold text-gray-900 dark:text-white">{{ t('settings.my_account') }}</h2>

    <!-- 主内容区：左侧导航 + 右侧表单 -->
    <div class="flex flex-col gap-6 lg:flex-row lg:gap-8">
      <!-- 左侧设置导航菜单 -->
      <SettingNav />

      <!-- 右侧内容区域 -->
      <div class="flex-1 space-y-6">
        <!-- 区块标题 -->
        <h3 class="mb-4 text-sm font-semibold uppercase text-gray-400 tracking-wider dark:text-gray-500">{{ t('settings.email') }}</h3>

        <!-- 邮箱输入和修改按钮区域 -->
        <div class="flex flex-col gap-4 sm:flex-row sm:items-start">
          <!-- 邮箱输入框 -->
          <div class="flex-1">
            <input
              v-model="form.email"
              type="email"
              class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
              :class="errors.email ? 'border-red-500' : 'border-gray-300'"
              :placeholder="t('message.your_email_address')"
              @keydown.enter="handleChangeEmail"
            />
            <!-- 验证错误提示 -->
            <span v-if="errors.email" class="mt-1 block text-xs text-red-500">{{ errors.email }}</span>
          </div>

          <!-- 修改邮箱按钮 -->
          <button
            class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500/20 focus:outline-none disabled:active:scale-100"
            :disabled="loading"
            @click="handleChangeEmail"
          >
            {{ t('settings.change_email') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
