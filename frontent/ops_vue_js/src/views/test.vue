<script setup>
import { useUserStore } from '@/stores/user'

import imageCropper from '@/components/imageCropper.vue';

import croppertest from '@/components/croppertest.vue';

const user = useUserStore()



</script>

<template>
    test

    <imageCropper></imageCropper>
   
    
</template>