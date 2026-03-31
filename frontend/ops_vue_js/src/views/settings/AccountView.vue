<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { useValidation } from '@/composables'
import { authApi } from '@/api/auth'
import SettingNav from '@/components/SettingNav.vue'
import ImageCropper from '@/components/imageCropper.vue'

usePageTitle('settings.account_settings')
const { t } = useI18n()
const userStore = useUserStore()
const toast = useToastStore()
const { validate, errors, clearErrors } = useValidation()

const form = reactive({
  username: '',
  remark: '',
  birthday: '',
})

const avatarHasChanged = ref(false)
const avatarDataUrl = ref('')
const loading = ref(false)

onMounted(() => {
  if (userStore.user) {
    form.username = userStore.user.Username || ''
    form.remark = userStore.user.FirstName || ''
    form.birthday = userStore.birthday
  }
})

function handleCrop(dataUrl) {
  avatarHasChanged.value = true
  avatarDataUrl.value = dataUrl
}

function cancelAvatar() {
  avatarHasChanged.value = false
  avatarDataUrl.value = ''
}

function base64ToFile(base64) {
  const [info, data] = base64.split(',')
  const mime = info.match(/:(.*?);/)[1]
  const bytes = atob(data)
  const arr = new Uint8Array(bytes.length)
  for (let i = 0; i < bytes.length; i++) arr[i] = bytes.charCodeAt(i)
  return new File([arr], 'avatar.png', { type: mime })
}

async function handleSave() {
  clearErrors()

  const err1 = validate('username', form.username, t('settings.name'))
  const err2 = validate('remark', form.remark, t('settings.remark'))
  const err3 = validate('birthday', form.birthday, t('settings.birthday'))

  if (!err1 || !err2 || !err3) return

  loading.value = true
  try {
    if (avatarHasChanged.value) {
      const file = base64ToFile(avatarDataUrl.value)
      await authApi.updateAvatar(file)
      avatarHasChanged.value = false
    }

    const { errCode } = await authApi.updateInfo({
      username: form.username,
      remark: form.remark,
      birthday: form.birthday,
    })

    if (errCode === 0) {
      toast.success(t('message.save_ok'))
      await userStore.fetchUserInfo()
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    // 拦截器已处理
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-5xl px-6 py-6">
    <h2 class="mb-6 text-2xl font-bold text-gray-900 dark:text-white">{{ t('settings.my_account') }}</h2>
    <div class="flex flex-col gap-6 lg:flex-row lg:gap-8">
      <SettingNav />
      <div class="flex-1 space-y-6">
        <!-- Avatar Section -->
        <div class="mb-8 rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-800 dark:bg-gray-900/50">
          <h3 class="mb-4 text-lg font-semibold text-gray-900 dark:text-white">{{ t('settings.profile_picture') }}</h3>
          <div class="flex flex-col items-start gap-6 md:flex-row md:items-center">
            <!-- Avatar Preview -->
            <div class="relative">
              <img
                :src="avatarHasChanged ? avatarDataUrl : userStore.avatarUrl"
                alt="Avatar"
                class="h-24 w-24 rounded-full border-4 border-white shadow-lg dark:border-gray-800"
              />
              <div class="absolute -right-1 -top-1 h-6 w-6 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 p-0.5">
                <div class="h-full w-full rounded-full bg-white dark:bg-gray-900"></div>
              </div>
            </div>
            
            <!-- Avatar Actions -->
            <div class="flex-1 space-y-4">
              <div class="space-y-3">
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  {{ t('settings.avatar_description') }}
                </p>
                
                <!-- Image Cropper Component -->
                <ImageCropper @crop-data-url="handleCrop" />
                
                <!-- Cancel Button (when avatar changed) -->
                <div v-if="avatarHasChanged" class="flex items-center gap-3 pt-2">
                  <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                    <div class="h-2 w-2 rounded-full bg-blue-500 animate-pulse"></div>
                    {{ t('settings.avatar_unsaved') }}
                  </div>
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

        <!-- Profile Information Form -->
        <div class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-800 dark:bg-gray-900/50">
          <div class="mb-6">
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">{{ t('settings.profile_information') }}</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              {{ t('settings.basic_information') }}
            </p>
          </div>

          <!-- Form Grid -->
          <div class="space-y-6">
            <!-- Name and Remark Row -->
            <div class="grid gap-6 md:grid-cols-2">
              <div class="space-y-2">
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('settings.name') }}
                  <span class="ml-1 text-red-500">*</span>
                </label>
                <div class="relative">
                  <input
                    v-model="form.username"
                    type="text"
                    placeholder="请输入您的姓名"
                    class="w-full rounded-lg border bg-white px-4 py-3 text-sm outline-none transition-all focus:ring-2 dark:bg-gray-900 dark:text-white"
                    :class="errors.username ? 'border-red-500 focus:border-red-500 focus:ring-red-500/20' : 'border-gray-300 focus:border-blue-500 focus:ring-blue-500/20 dark:border-gray-700'"
                  />
                  <div class="absolute right-3 top-3">
                    <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                  </div>
                </div>
                <span v-if="errors.username" class="block text-xs text-red-500">{{ errors.username }}</span>
              </div>

              <div class="space-y-2">
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {{ t('settings.remark') }}
                  <span class="ml-1 text-gray-400">({{ t('settings.optional') }})</span>
                </label>
                <div class="relative">
                  <input
                    v-model="form.remark"
                    type="text"
                    placeholder="个人简介或备注"
                    class="w-full rounded-lg border bg-white px-4 py-3 text-sm outline-none transition-all focus:ring-2 dark:bg-gray-900 dark:text-white"
                    :class="errors.remark ? 'border-red-500 focus:border-red-500 focus:ring-red-500/20' : 'border-gray-300 focus:border-blue-500 focus:ring-blue-500/20 dark:border-gray-700'"
                  />
                  <div class="absolute right-3 top-3">
                    <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />
                    </svg>
                  </div>
                </div>
                <span v-if="errors.remark" class="block text-xs text-red-500">{{ errors.remark }}</span>
              </div>
            </div>

            <!-- Birthday Row -->
            <div class="space-y-2">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('settings.birthday') }}
                <span class="ml-1 text-gray-400">({{ t('settings.optional') }})</span>
              </label>
              <div class="relative">
                <input
                  v-model="form.birthday"
                  type="date"
                  class="w-full rounded-lg border bg-white px-4 py-3 text-sm outline-none transition-all focus:ring-2 dark:bg-gray-900 dark:text-white"
                  :class="errors.birthday ? 'border-red-500 focus:border-red-500 focus:ring-red-500/20' : 'border-gray-300 focus:border-blue-500 focus:ring-blue-500/20 dark:border-gray-700'"
                />
                <div class="absolute right-3 top-3">
                  <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                </div>
              </div>
              <span v-if="errors.birthday" class="block text-xs text-red-500">{{ errors.birthday }}</span>
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                {{ t('settings.birthday_help') }}
              </p>
            </div>
          </div>

          <!-- Save Button -->
          <div class="mt-8 border-t border-gray-100 pt-6 dark:border-gray-800">
            <div class="flex items-center justify-between">
              <div class="text-sm text-gray-600 dark:text-gray-400">
                {{ t('settings.save_notice') }}
              </div>
              <button
                class="group flex items-center gap-2 rounded-lg bg-gradient-to-r from-blue-600 to-blue-500 px-6 py-3 text-sm font-semibold text-white shadow-sm transition-all hover:from-blue-700 hover:to-blue-600 hover:shadow-md focus:outline-none focus:ring-2 focus:ring-blue-500/30"
                :disabled="loading"
                @click="handleSave"
              >
                <svg v-if="!loading" class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                <svg v-else class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                {{ loading ? t('settings.saving') : t('settings.save_changes') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
