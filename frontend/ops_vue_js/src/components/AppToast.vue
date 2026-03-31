<script setup>
import { useToastStore } from '@/stores/toast'
import { IconAlertCircle, IconAlertTriangle, IconCheck, IconInfoCircle, IconX } from '@tabler/icons-vue'

const toastStore = useToastStore()

const icons = {
  success: IconCheck,
  warning: IconAlertTriangle,
  error: IconAlertCircle,
  info: IconInfoCircle,
}
</script>

<template>
  <Transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="translate-y-2 opacity-0"
    enter-to-class="translate-y-0 opacity-100"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="translate-y-0 opacity-100"
    leave-to-class="translate-y-2 opacity-0"
  >
    <div
      v-if="toastStore.visible"
      class="fixed left-1/2 top-5 z-[9999] flex max-w-sm -translate-x-1/2 transform items-start gap-3 rounded-lg border-0 px-4 py-3 shadow-lg text-white"
      :class="{
        'bg-green-600': toastStore.type === 'success',
        'bg-yellow-600': toastStore.type === 'warning',
        'bg-red-600': toastStore.type === 'danger',
        'bg-slate-700': toastStore.type === 'info',
      }"
      role="alert"
    >
      <component :is="icons[toastStore.type] || IconInfoCircle" :size="20" class="mt-0.5 shrink-0" />
      <span class="flex-1 text-sm">{{ toastStore.message }}</span>
      <button class="ml-1 opacity-70 hover:opacity-100" @click="toastStore.hide()">×</button>
    </div>
  </Transition>
</template>
