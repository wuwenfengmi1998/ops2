<script setup>


import { Modal } from "@tabler/core";
import { onMounted, ref } from "vue";

import Cropper from 'cropperjs';



const image = new Image();
image.src = "https://wnfed.com/usr/uploads/2020/07/736937178.jpg";
image.alt = 'Picture';

const imagecroppert=ref()

const cropper = ref(null);


onMounted(()=>{

      cropper.value = new Cropper(image, {
      container: imagecroppert.value,
     
      

    });

    //cropper.value.
    console.log(cropper.value)



})

function getsele(){
  console.log(cropper.value.getCropperCanvas())

}


</script>

<template>

<div class="cropper-container" ref="imagecroppert"></div>


<button @click="getsele">getsele</button>



</template>

<style lang="scss" scoped>
.cropper-container {
  border: 1px solid var(--vp-c-divider);
  border-radius: 0.375rem;
  margin-bottom: 1rem;
  margin-top: 1rem;
  padding: 1.25rem 1.5rem;

  :deep(cropper-canvas) {
    height: 360px;
  }
}
</style>