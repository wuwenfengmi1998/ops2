<script setup lang="ts">
  import MyOffcanvas from '@/components/MyOffcanvas.vue';
  import { onMounted, ref } from 'vue';
  const mos = ref<InstanceType<typeof MyOffcanvas> | null>(null);

  onMounted(() => {
    // 可以在这里调用mos.value的方法
    //console.log('HomeView mounted', mos);
    document.title = 'Operations.';
  });

  function c(){
    mos.value?.showAlert("success","111",1000);
  }

</script>

<template>

<div>
  <MyOffcanvas ref="mos"/>
  <button @click="c">ttt</button>
  111
</div>

</template>
