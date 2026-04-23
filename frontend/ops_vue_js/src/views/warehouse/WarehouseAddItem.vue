<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToastStore } from '@/stores/toast'
import { usePageTitle } from '@/composables/usePageTitle'
import { warehouseApi } from '@/api/warehouse'
import useDropzone from '@/components/useDropzone.vue'

usePageTitle('warehouse.add_item')
const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const toast = useToastStore()

const containerId = ref(parseInt(route.params.id))

const form = reactive({
  name: '',
  serial_number: '',
  remark: '',
  quantity: 1,
})

const submitting = ref(false)
const loadingContainer = ref(true)
const containerName = ref('')

const dropzoneRef = ref(null)

function getPhotoHashes() {
  return dropzoneRef.value?.return_files().map((f) => f.hash) ?? []
}

onMounted(async () => {
  try {
    const { errCode, data } = await warehouseApi.getContainer(containerId.value)
    if (errCode === 0 && data?.container) {
      containerName.value = data.container.Title
    }
  } catch {
    // 找不到容器时仍允许提交
  } finally {
    loadingContainer.value = false
  }
})

async function submit() {
  if (!form.name.trim()) {
    toast.error(t('warehouse.item_name_required'))
    return
  }

  // 等待图片上传完成
  await new Promise((r) => setTimeout(r, 200))
  const hashes = getPhotoHashes()

  submitting.value = true
  try {
    const { errCode } = await warehouseApi.addItem({
      container_id: containerId.value,
      name: form.name.trim(),
      serial_number: form.serial_number.trim(),
      remark: form.remark.trim(),
      quantity: form.quantity > 0 ? form.quantity : 1,
      photos: hashes,
    })
    if (errCode === 0) {
      toast.success(t('message.save_success'))
      router.push(`/warehouse/container/${containerId.value}`)
    } else {
      toast.error(t('message.server_error'))
    }
  } catch {
    toast.error(t('message.server_error'))
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="p-4 max-w-2xl mx-auto space-y-4">

    <!-- 面包屑 -->
    <div class="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
      <RouterLink to="/warehouse/container" class="text-blue-500 hover:underline">
        {{ t('warehouse.container_list') }}
      </RouterLink>
      <span>/</span>
      <RouterLink
        v-if="containerName"
        :to="`/warehouse/container/${containerId}`"
        class="text-blue-500 hover:underline"
      >
        {{ containerName }}
      </RouterLink>
      <svg v-if="loadingContainer" class="h-3.5 w-3.5 animate-spin text-gray-400" viewBox="0 0 24 24" fill="none">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
      </svg>
      <span v-else>/</span>
      <span class="text-gray-700 dark:text-gray-200">{{ t('warehouse.add_item') }}</span>
    </div>

    <!-- 表单卡片 -->
    <div class="rounded-xl border border-gray-200 bg-white px-5 py-5 shadow dark:border-dk-muted dark:bg-dk-card">
      <h2 class="mb-5 text-base font-semibold text-gray-900 dark:text-white">
        {{ t('warehouse.add_item') }}
        <span v-if="containerName" class="ml-2 text-sm font-normal text-gray-500">
          → {{ containerName }}
        </span>
      </h2>

      <div class="space-y-4">

        <!-- 物品名称 -->
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('warehouse.item_name') }} *
          </label>
          <input
            v-model="form.name"
            type="text"
            :placeholder="t('warehouse.item_name_placeholder')"
            class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
            @keyup.enter="submit"
          />
        </div>

        <!-- 序列号 -->
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('warehouse.serial_number') }}
          </label>
          <input
            v-model="form.serial_number"
            type="text"
            :placeholder="t('warehouse.serial_number_placeholder')"
            class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
          />
        </div>

        <!-- 数量 -->
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('warehouse.quantity') }}
          </label>
          <input
            v-model.number="form.quantity"
            type="number"
            min="1"
            class="w-28 rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
          />
        </div>

        <!-- 备注 -->
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('warehouse.remark') }}
          </label>
          <textarea
            v-model="form.remark"
            :placeholder="t('warehouse.remark_placeholder')"
            rows="3"
            class="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white"
          ></textarea>
        </div>

        <!-- 图片上传 -->
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ t('purchase_addorder.upload_photos') }}
          </label>
          <useDropzone
            ref="dropzoneRef"
            uploadURL="/api/files/upload/image"
            :max-files="9"
          />
        </div>

      </div>

      <!-- 操作按钮 -->
      <div class="mt-6 flex justify-end gap-2">
        <button
          class="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm dark:border-dk-muted dark:bg-dk-base dark:text-white hover:bg-gray-50 dark:hover:bg-dk-muted"
          @click="router.push(`/warehouse/container/${containerId}`)"
        >
          {{ t('message.cancel') }}
        </button>
        <button
          class="inline-flex items-center gap-1.5 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
          :disabled="submitting"
          @click="submit"
        >
          <svg v-if="submitting" class="h-4 w-4 animate-spin" viewBox="0 0 24 24" fill="none">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
          </svg>
          {{ t('message.save') }}
        </button>
      </div>
    </div>
  </div>
</template>
