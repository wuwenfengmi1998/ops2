<script setup lang="ts">

import {myfuncs} from './myfunc.js';

import HeardMain from './components/HeardMain.vue';



myfuncs.test();

onload = () => {
  console.log("Page loaded");
}



</script>

<template>

  <HeardMain />

  <RouterView />
</template>

<style scoped></style>
