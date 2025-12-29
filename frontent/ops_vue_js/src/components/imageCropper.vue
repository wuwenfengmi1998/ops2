<script setup>


import { Modal } from "@tabler/core";
import { onMounted, ref } from "vue";



import 'cropperjs';



const image = new Image();
image.src = "https://wnfed.com/usr/uploads/2020/07/736937178.jpg";
image.alt = 'Picture';

const imagecroppert=ref()

const cropper = ref(null);


onMounted(()=>{


})

function getsele(){

}


</script>

<template>

<cropper-canvas class="cropper-container" background>
  <cropper-image src="https://wnfed.com/usr/uploads/2020/07/736937178.jpg" alt="Picture" rotatable scalable skewable translatable></cropper-image>
  <cropper-shade hidden></cropper-shade>
  <cropper-handle action="select" plain></cropper-handle>
  <cropper-selection initial-coverage="0.5" movable resizable>
    <cropper-grid role="grid" covered></cropper-grid>
    <cropper-crosshair centered></cropper-crosshair>
    <cropper-handle action="move" theme-color="rgba(255, 255, 255, 0.35)"></cropper-handle>
    <cropper-handle action="n-resize"></cropper-handle>
    <cropper-handle action="e-resize"></cropper-handle>
    <cropper-handle action="s-resize"></cropper-handle>
    <cropper-handle action="w-resize"></cropper-handle>
    <cropper-handle action="ne-resize"></cropper-handle>
    <cropper-handle action="nw-resize"></cropper-handle>
    <cropper-handle action="se-resize"></cropper-handle>
    <cropper-handle action="sw-resize"></cropper-handle>
  </cropper-selection>
</cropper-canvas>
<button @click="getsele">getsele</button>



</template>

<style lang="scss" scoped>
.cropper-container {


  width: 400px;
  height: 400px;

 
}
</style>