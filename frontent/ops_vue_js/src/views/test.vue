<script setup>
import { useUserStore } from '@/stores/user'

const user = useUserStore()

function t(){
    user.login();
}

</script>

<template>
    test
    <button @click="t">222</button>
    {{ user.isLoggedIn }}
</template>