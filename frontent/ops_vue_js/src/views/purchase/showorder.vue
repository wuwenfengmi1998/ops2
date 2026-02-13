<script setup>
import { onMounted, watch, ref, reactive } from "vue";
import { useRouter } from "vue-router";
const router = useRouter();

const dynamicParam = router

onMounted(() => {
  console.log(dynamicParam);
});
</script>
<template>
  
</template>
