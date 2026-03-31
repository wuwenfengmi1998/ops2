<script setup lang="ts">
import MyOffcanvas from '@/components/MyOffcanvas.vue'
import { onMounted, ref } from 'vue'
//import { my_network_func } from '@/my_network_func'
const mos = ref<InstanceType<typeof MyOffcanvas> | null>(null)

onMounted(() => {
  // 可以在这里调用mos.value的方法
  //console.log('HomeView mounted', mos);
  document.title = 'Operations.'
})

function c() {
  mos.value?.showAlert('success', '111', 1000)
}

function d() {
  interface User {
    id: number
    name: string
  }

  const Users: Partial<User> = {}
  Users.id = 2
  Users.name = 'Bob'

     const jsonData = ref<User[]>([])
   jsonData.value = [{ id: 1, name: 'Alice' }]


  console.log('Users:', Users)

  //my_network_func.post_json('/1',jsonData.value[0],()=>{})
}
</script>

<template>
  <div>
    <MyOffcanvas ref="mos" />
    <button @click="c">test1</button>
    <button @click="d">test2</button>
  </div>
</template>
