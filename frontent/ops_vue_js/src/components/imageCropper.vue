<script setup>
import { Modal } from "@tabler/core";
import { onMounted, ref } from "vue";

import "cropperjs";

const cro_sele = ref();
const cro_canv = ref();
const cro_imag = ref();

var cor_size_width = 300;
var cor_size_height = 300;

const is_have_URL = ref(false);

const reader = new FileReader();
reader.onload = () => {
  initCropper(reader.result);
};

onMounted(() => {
  cro_sele.value.$change(0, 0, cor_size_width, cor_size_height);
  cro_canv.value.style.width = cor_size_width.toString() + "px";
  cro_canv.value.style.height = cor_size_height.toString() + "px";

  cro_imag.value.src = "https://wnfed.com/usr/uploads/2020/07/736937178.jpg";

  //console.log(cro_canv.value.clientHeight)
});

function initCropper(imageSrc){

  is_have_URL.value = true;
  cro_imag.value.src=imageSrc;
}

function inputfile(e){
    const file = e.target.files[0];
    if (!file){
      e.target.value=""
      is_have_URL.value = false;
      return;
    }

    if (!file.type.startsWith('image/')) {
      e.target.value=""
      is_have_URL.value = false;
      return;
    }

    reader.readAsDataURL(file);
 
  }

function getsele() {
  console.log(cro_canv.value.$toCanvas())
}
</script>

<template>
  <div>
    <div class="col-6 col-sm-4 col-md-2 col-xl py-3">
      <input class="btn btn-outline-primary" type="file" accept="image/*" @change="inputfile">
    </div>

    <cropper-canvas
      ref="cro_canv"
      class="cropper-container"
      :hidden="!is_have_URL"
      background
      scale-step="0.1"
    >
      <cropper-image
        ref="cro_imag"
        src=""
        alt="Picture"
        initialCenterSize="cover"
        rotatable
        scalable
        skewable
        translatable
      ></cropper-image>
      <cropper-shade hidden></cropper-shade>
      <cropper-handle action="move" plain></cropper-handle>
      <cropper-selection ref="cro_sele">
        <cropper-grid role="grid" covered></cropper-grid>
        <cropper-crosshair centered></cropper-crosshair>
        <cropper-handle
          action="move"
          theme-color="rgba(255, 255, 255, 0)"
        ></cropper-handle>
      </cropper-selection>
    </cropper-canvas>
  </div>

  <div class="col-6 col-sm-4 col-md-2 col-xl py-3">
    <button class="btn btn-outline-primary" @click="getsele">getsele</button>
  </div>
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
