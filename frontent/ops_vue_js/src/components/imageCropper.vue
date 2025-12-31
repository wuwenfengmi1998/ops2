<script setup>
import { Modal } from "@tabler/core";
import { onMounted, ref } from "vue";

import "cropperjs";

const cro_sele=ref()
const cro_canv=ref()

var cor_size_width=400
var cor_size_height=400

onMounted(() => {
  cro_sele.value.$change(0,0,cor_size_width,cor_size_height)
  cro_canv.value.style.width=cor_size_width.toString()+"px"
  cro_canv.value.style.height=cor_size_height.toString()+"px"

  //console.log(cro_canv.value.clientHeight)
});

function getsele() {}
</script>

<template>
  <cropper-canvas ref="cro_canv" class="cropper-container" background scale-step="0.2">
    <cropper-image
      src="https://wnfed.com/usr/uploads/2020/07/736937178.jpg"
      alt="Picture"
      initialCenterSize="cover"
      rotatable
      scalable
      skewable
      translatable
    ></cropper-image>
    <cropper-shade hidden></cropper-shade>
    <cropper-handle action="move" plain></cropper-handle>
    <cropper-selection
      ref="cro_sele"
    >
      <cropper-grid role="grid" covered></cropper-grid>
      <cropper-crosshair centered></cropper-crosshair>
      <cropper-handle
        action="move"
        theme-color="rgba(255, 255, 255, 0)"
      ></cropper-handle>

    </cropper-selection>
  </cropper-canvas>
  <button @click="getsele">getsele</button>
</template>

<style lang="scss" scoped>
.cropper-container {

  
   /* 四个角相同圆角 */
  border-radius: 10px;
  /* 基本描边 */
  //border: 2px solid #333;

   /* 基本阴影：x偏移 y偏移 模糊半径 扩展半径 颜色 */
  box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3);
}
</style>
