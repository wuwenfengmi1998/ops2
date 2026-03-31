<script setup>
import AppToast from '@/components/AppToast.vue'
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-gradient-to-br from-gray-50 to-gray-100 dark:from-dk-base dark:to-dk-card">
    <div class="w-full">
      <RouterView />
    </div>
    <AppToast />
  </div>
</template>
