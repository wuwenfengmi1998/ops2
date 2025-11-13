<script setup>
import { ref } from "vue";
import { my_network_func } from "@/my_network_func";
import { useUserStore } from "@/stores/user";
import MyOffcanvas from "@/components/MyOffcanvas.vue";
const user = useUserStore();
const mos = ref();

function t() {
  console.log("test");

  my_network_func.postJson(
    "/users/test",
    {
      a: "1",
    },
    (a) => {
      console.log(a);
    }
  );
}
function tt(){
  mos.value?.showAlert("danger", "123123", 5000);
}
</script>

<template>
  <main>
    1112
    <button @click="t">222</button>
    {{ user.userCookie }}
    <button @click="tt">333</button>
  </main>
  <MyOffcanvas ref="mos" />
</template>
