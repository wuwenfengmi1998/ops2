<script setup>
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import AppToast from '@/components/AppToast.vue'
</script>

<template>
  <div class="flex min-h-screen flex-col">
    <AppHeader />
    <main class="flex-1">
      <RouterView />
    </main>
    <AppFooter />
    <AppToast />
  </div>
</template>
