<script setup lang="ts">
  import MyOffcanvas from '@/components/MyOffcanvas.vue';
  import { onMounted, ref } from 'vue';
  const mos = ref<InstanceType<typeof MyOffcanvas> | null>(null);

  onMounted(() => {
    // 可以在这里调用mos.value的方法
    console.log('HomeView mounted', mos);
  });


</script>

<template>

<div>
  <MyOffcanvas ref="mos"/>
  <button >ttt</button>
  111
</div>

</template>
