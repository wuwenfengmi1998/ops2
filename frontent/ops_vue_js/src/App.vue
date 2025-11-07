<script setup>
import { RouterLink, RouterView } from 'vue-router'
import HeardMain from './components/HeardMain.vue';
import FooterMain from './components/FooterMain.vue';


</script>

<template>

 <div class="page">

<HeardMain />

  <RouterView />

  <FooterMain />

  </div>
</template>
