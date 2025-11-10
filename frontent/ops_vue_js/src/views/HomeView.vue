<script setup>
import { my_network_func } from '@/my_network_func';

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
    111
    <button @click="t">222</button>
  </main>
</template>
