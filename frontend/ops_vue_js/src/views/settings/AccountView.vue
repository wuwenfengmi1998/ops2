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
        <!-- Avatar -->
        <div class="mb-6 flex items-center gap-4">
          <div>
            <img
              :src="avatarHasChanged ? avatarDataUrl : userStore.avatarUrl"
              alt="Avatar"
              class="h-16 w-16 rounded-full border-2 border-gray-200 object-cover dark:border-dk-muted"
            />
          </div>
          <div>
            <ImageCropper @crop-data-url="handleCrop" />
            <button v-if="avatarHasChanged" class="mt-2 rounded-lg border border-gray-300 px-3 py-1 text-xs text-gray-600 hover:bg-gray-50 dark:border-dk-muted dark:text-gray-400 dark:hover:bg-dk-card" @click="cancelAvatar">
              {{ t('settings.cancel') }}
            </button>
          </div>
        </div>

        <h3 class="mb-4 text-sm font-semibold uppercase text-gray-400 tracking-wider dark:text-gray-500">Profile</h3>

        <!-- Form -->
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('settings.name') }}</label>
            <input
              v-model="form.username"
              type="text"
              class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
              :class="errors.username ? 'border-red-500' : 'border-gray-300'"
            />
            <span v-if="errors.username" class="mt-1 block text-xs text-red-500">{{ errors.username }}</span>
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('settings.remark') }}</label>
            <input
              v-model="form.remark"
              type="text"
              class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
              :class="errors.remark ? 'border-red-500' : 'border-gray-300'"
            />
            <span v-if="errors.remark" class="mt-1 block text-xs text-red-500">{{ errors.remark }}</span>
          </div>
          <div>
            <label class="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('settings.birthday') }}</label>
            <input
              v-model="form.birthday"
              type="date"
              class="w-full rounded-lg border border-gray-300 bg-white px-3.5 py-2 text-sm outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-dk-muted dark:bg-dk-base dark:text-white"
              :class="errors.birthday ? 'border-red-500' : 'border-gray-300'"
            />
            <span v-if="errors.birthday" class="mt-1 block text-xs text-red-500">{{ errors.birthday }}</span>
          </div>
        </div>

        <div class="mt-6">
          <button
            class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500/20 focus:outline-none disabled:active:scale-100"
            :disabled="loading"
            @click="handleSave"
          >
            {{ t('settings.save_changes') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
