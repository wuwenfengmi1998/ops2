<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import CustomerFormModal from './CustomerFormModal.vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const customerId = ref(null)
const loading = ref(true)

onMounted(() => {
  const id = route.params.id
  if (!id) {
    router.push('/customer')
    return
  }
  customerId.value = parseInt(id)
  loading.value = false
})

function onClose() {
  router.push('/customer')
}

function onSubmitted() {
  router.push(`/customer/detail/${customerId.value}`)
}
</script>

<template>
  <div v-if="!loading && customerId">
    <CustomerFormModal
      :title="t('customer.edit_title')"
      :customer="{ id: customerId }"
      @close="onClose"
      @submit="onSubmitted"
    />
  </div>
</template>
