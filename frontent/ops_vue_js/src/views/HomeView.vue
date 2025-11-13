<script setup>
import { my_network_func } from '@/my_network_func';
import { useUserStore } from '@/stores/user'

const user = useUserStore()


  function t(){
    console.log("test")

    my_network_func.postJson("/users/test",{
      a:"1"
    },(a)=>{
      console.log(a)
    });
  }
</script>

<template>
  <main>
    1112
    <button @click="t">222</button>
    {{ user.userCookie }}
  </main>
</template>
