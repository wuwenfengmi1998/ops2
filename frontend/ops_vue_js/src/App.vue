<script setup>
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import { useUserStore } from '@/stores/user'

// Restore login state from localStorage/sessionStorage
onMounted(() => {
  const userStore = useUserStore()
  userStore.restoreSession()
})
</script>

<template>
  <RouterView />
</template>
