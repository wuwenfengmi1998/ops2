<script setup>
import { ref } from "vue";
import { my_network_func } from "@/my_network_func";
import { useUserStore } from "@/stores/user";
import MyOffcanvas from "@/components/MyOffcanvas.vue";
import imageCropper from "@/components/imageCropper.vue";
const user = useUserStore();
const mos = ref();


</script>

<template>
  

  <MyOffcanvas ref="mos" />
</template>
