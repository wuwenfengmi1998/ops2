<script setup>
// Vue 核心响应式 API
import { reactive, ref, onMounted } from 'vue'
// 国际化 hook，用于翻译文本
import { useI18n } from 'vue-i18n'
// 用户状态管理
import { useUserStore } from '@/stores/user'
// Toast 通知状态管理
import { useToastStore } from '@/stores/toast'
// 页面标题 composable
import { usePageTitle } from '@/composables/usePageTitle'
// 表单验证 composable
import { useValidation } from '@/composables'
// 认证 API
import { authApi } from '@/api/auth'
// 设置页面导航侧边栏组件
import SettingNav from '@/components/SettingNav.vue'
// 图片裁剪组件
import ImageCropper from '@/components/imageCropper.vue'

// 设置页面标题，使用国际化 key
usePageTitle('settings.account_settings')

// 初始化国际化翻译函数
const { t } = useI18n()

// 获取用户 store 实例，用于访问和更新用户信息
const userStore = useUserStore()

// 获取 toast store 实例，用于显示操作结果通知
const toast = useToastStore()

// 解构验证函数和错误状态，用于表单字段验证
const { validate, errors, clearErrors } = useValidation()

// 表单数据，使用 reactive 创建响应式对象
const form = reactive({
  username: '',  // 用户姓名
  remark: '',     // 备注/昵称
  birthday: '',   // 生日日期
})

// 头像是否已被修改的标记
const avatarHasChanged = ref(false)

// 裁剪后的头像 Base64 数据
const avatarDataUrl = ref('')

// 保存按钮的加载状态，防止重复提交
const loading = ref(false)

// 生日输入框的 DOM 引用，用于调用原生日期选择器
const birthdayInput = ref(null)

// 组件挂载时执行的初始化逻辑
onMounted(() => {
  // 如果用户信息已加载，则填充表单
  if (userStore.user || userStore.userInfo) {
    // 姓名从 userInfo.Username 获取
    form.username = userStore.userInfo?.Username || userStore.user?.Name || ''
    // 备注从 userInfo.FirstName 获取
    form.remark = userStore.userInfo?.FirstName || ''
    // 生日从 userStore.birthday getter 获取（已转换格式）
    form.birthday = userStore.birthday
  }
})

/**
 * 打开原生日期选择器
 * 使用 HTML5 showPicker API（在支持的浏览器中）或回退到 focus 事件
 */
function openDatePicker() {
  // 使用showPicker API打开日期选择器（现代浏览器支持）
  if (birthdayInput.value && birthdayInput.value.showPicker) {
    birthdayInput.value.showPicker()
  } else {
    // 对于不支持showPicker的老浏览器，聚焦输入框触发选择器
    birthdayInput.value?.focus()
  }
}

/**
 * 处理裁剪后的头像数据
 * @param {string} dataUrl - 裁剪后的图片 Base64 URL
 */
function handleCrop(dataUrl) {
  avatarHasChanged.value = true    // 标记头像已修改
  avatarDataUrl.value = dataUrl    // 保存裁剪后的数据
}

/**
 * 取消头像修改，恢复原状
 */
function cancelAvatar() {
  avatarHasChanged.value = false   // 重置修改标记
  avatarDataUrl.value = ''         // 清空裁剪数据
}

/**
 * 将 Base64 字符串转换为 File 对象
 * @param {string} base64 - 包含 MIME 类型和数据的前缀的 Base64 字符串
 * @returns {File} 转换后的文件对象
 */
function base64ToFile(base64) {
  // 分离 MIME 信息和实际数据部分
  const [info, data] = base64.split(',')
  // 从 MIME 信息中提取文件类型
  const mime = info.match(/:(.*?);/)[1]
  // 解码 Base64 数据为二进制字符串
  const bytes = atob(data)
  // 转换为字节数组
  const arr = new Uint8Array(bytes.length)
  for (let i = 0; i < bytes.length; i++) arr[i] = bytes.charCodeAt(i)
  // 创建并返回 File 对象
  return new File([arr], 'avatar.png', { type: mime })
}

/**
 * 保存账户信息
 * 包含头像上传和个人资料更新
 */
async function handleSave() {
  // 清空之前的验证错误
  clearErrors()

  // 验证必填字段：用户名
  const err1 = validate('username', form.username, t('settings.name'))
  // 验证必填字段：备注
  const err2 = validate('remark', form.remark, t('settings.remark'))
  // 验证必填字段：生日
  const err3 = validate('birthday', form.birthday, t('settings.birthday'))

  // 如果任一验证失败，则阻止提交
  if (!err1 || !err2 || !err3) return

  // 设置加载状态，防止重复提交
  loading.value = true
  try {
    // 如果头像有修改，先上传新头像
    if (avatarHasChanged.value) {
      // 将 Base64 转换为 File 对象
      const file = base64ToFile(avatarDataUrl.value)
      // 调用 API 上传头像
      await authApi.updateAvatar(file)
      // 重置头像修改状态
      avatarHasChanged.value = false
    }

    // 更新用户基本信息
    const { errCode } = await authApi.updateInfo({
      username: form.username,   // 用户名
      remark: form.remark,        // 备注
      birthday: form.birthday,    // 生日
    })

    // 根据返回的错误码判断操作结果
    if (errCode === 0) {
      // 保存成功，显示成功提示
      toast.success(t('message.save_ok'))
      // 重新获取用户信息以更新本地状态
      await userStore.fetchUserInfo()
    } else {
      // 服务器错误，显示错误提示
      toast.error(t('message.server_error'))
    }


    
  } catch {
    // 网络错误等异常已被 axios 拦截器统一处理
  } finally {
    // 无论成功或失败，都要关闭加载状态
    loading.value = false
  }
}
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

      <!-- 右侧表单区域 -->
      <div class="flex-1 space-y-6">

        <!-- ========== 头像上传区域 ========== -->
        <div class="mb-8 rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-800 dark:bg-gray-900/50">
          <!-- 区域标题 -->
          <h3 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">{{ t('settings.profile_picture') }}</h3>

          <!-- 头像和操作横向布局 -->
          <div class="flex flex-col items-start gap-6 md:flex-row md:items-center">
            <!-- 头像预览区域 -->
            <div class="relative">
              <!-- 头像图片：根据是否有修改显示新头像或原头像 -->
              <img
                :src="avatarHasChanged ? avatarDataUrl : userStore.avatarUrl"
                alt="Avatar"
                class="h-24 w-24 rounded-full border-4 border-white shadow-lg dark:border-gray-800"
              />
            </div>

            <!-- 头像操作区域 -->
            <div class="flex-1 space-y-4">
              <div class="space-y-3">
                <!-- 头像上传说明文字 -->
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  {{ t('settings.avatar_description') }}
                </p>

                <!-- 图片裁剪组件：用户选择并裁剪新头像 -->
                <ImageCropper @crop-data-url="handleCrop" />

                <!-- 取消按钮（仅在头像有修改时显示） -->
                <div v-if="avatarHasChanged" class="flex items-center gap-3 pt-2">
                  <!-- 未保存提示指示器 -->
                  <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                    <div class="h-2 w-2 rounded-full bg-blue-500 animate-pulse"></div>
                    {{ t('settings.avatar_unsaved') }}
                  </div>
                  <!-- 取消修改按钮 -->
                  <button
                    class="ml-auto rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 hover:border-gray-400 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-gray-800 dark:hover:border-gray-600"
                    @click="cancelAvatar"
                  >
                    {{ t('settings.cancel_changes') }}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- ========== 个人资料表单区域 ========== -->
        <div class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-800 dark:bg-gray-900/50">
          <!-- 表单区域标题 -->
          <div class="mb-6">
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">{{ t('settings.profile_information') }}</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              {{ t('settings.basic_information') }}
            </p>
          </div>

          <!-- 表单字段区域 -->
          <div class="space-y-6">
            <!-- 用户名和备注行（桌面端双列布局） -->
            <div class="grid gap-6 md:grid-cols-2">

              <!-- 用户名字段 -->
              <div class="space-y-2">
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('settings.name') }}
                  <!-- 必填标记 -->
                  <span class="ml-1 text-red-500">*</span>
                  <span class="ml-1 text-gray-400">{{ t('settings.name_hint') }}</span>
                </label>
                <div class="relative">
                  <!-- 用户名输入框 -->
                  <input
                    v-model="form.username"
                    type="text"
                    :placeholder="t('settings.placeholder_name')"
                    class="w-full rounded-lg border bg-white px-4 py-3 text-sm outline-none transition-all focus:ring-2 dark:bg-gray-900 dark:text-white"
                    :class="errors.username ? 'border-red-500 focus:border-red-500 focus:ring-red-500/20' : 'border-gray-300 focus:border-blue-500 focus:ring-blue-500/20 dark:border-gray-700'"
                  />
                  <!-- 用户图标（右侧装饰） -->
                  <div class="absolute right-3 top-3 pointer-events-none">
                    <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                  </div>
                </div>
                <!-- 验证错误提示 -->
                <span v-if="errors.username" class="block text-xs text-red-500">{{ errors.username }}</span>
              </div>

              <!-- 备注/昵称字段 -->
              <div class="space-y-2">
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('settings.remark') }}
                  <!-- 必填标记 -->
                  <span class="ml-1 text-red-500">*</span>
                  <span class="ml-1 text-gray-400">{{ t('settings.remark_hint') }}</span>
                </label>
                <div class="relative">
                  <!-- 备注输入框 -->
                  <input
                    v-model="form.remark"
                    type="text"
                    :placeholder="t('settings.placeholder_remark')"
                    class="w-full rounded-lg border bg-white px-4 py-3 text-sm outline-none transition-all focus:ring-2 dark:bg-gray-900 dark:text-white"
                    :class="errors.remark ? 'border-red-500 focus:border-red-500 focus:ring-red-500/20' : 'border-gray-300 focus:border-blue-500 focus:ring-blue-500/20 dark:border-gray-700'"
                  />
                  <!-- 备注图标（右侧装饰） -->
                  <div class="absolute right-3 top-3 pointer-events-none">
                    <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />
                    </svg>
                  </div>
                </div>
                <!-- 验证错误提示 -->
                <span v-if="errors.remark" class="block text-xs text-red-500">{{ errors.remark }}</span>
              </div>
            </div>

            <!-- 生日字段 -->
            <div class="space-y-2">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('settings.birthday') }}
                <!-- 可选标记 -->
                <span class="ml-1 text-gray-400">({{ t('settings.optional') }})</span>
              </label>
              <div class="relative">
                <!-- 日期选择输入框 -->
                <input
                  v-model="form.birthday"
                  type="date"
                  class="w-full cursor-pointer rounded-lg border bg-white px-4 py-3 text-sm outline-none transition-all focus:ring-2 dark:bg-gray-900 dark:text-white"
                  :class="errors.birthday ? 'border-red-500 focus:border-red-500 focus:ring-red-500/20' : 'border-gray-300 focus:border-blue-500 focus:ring-blue-500/20 dark:border-gray-700'"
                  @click="openDatePicker"
                  ref="birthdayInput"
                />
              </div>
              <!-- 验证错误提示 -->
              <span v-if="errors.birthday" class="block text-xs text-red-500">{{ errors.birthday }}</span>
              <!-- 帮助提示文字 -->
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                {{ t('settings.birthday_help') }}
              </p>
            </div>
          </div>

          <!-- 保存按钮区域 -->
          <div class="mt-8 border-t border-gray-100 pt-6 dark:border-gray-800">
            <div class="flex items-center justify-between">
              <!-- 保存提示文字 -->
              <div class="text-sm text-gray-600 dark:text-gray-400">
                {{ t('settings.save_notice') }}
              </div>
              <!-- 保存按钮 -->
              <button
                class="group flex items-center gap-2 rounded-lg bg-gradient-to-r from-blue-600 to-blue-500 px-6 py-3 text-sm font-semibold text-white shadow-sm transition-all hover:from-blue-700 hover:to-blue-600 hover:shadow-md focus:outline-none focus:ring-2 focus:ring-blue-500/30"
                :disabled="loading"
                @click="handleSave"
              >
                <!-- 保存成功图标（未加载时显示） -->
                <svg v-if="!loading" class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                <!-- 加载中旋转图标 -->
                <svg v-else class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                <!-- 按钮文字：加载中显示"保存中"，否则显示"保存更改" -->
                {{ loading ? t('settings.saving') : t('settings.save_changes') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
