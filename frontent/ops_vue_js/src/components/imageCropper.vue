<script setup>
import { Modal } from "@tabler/core";
import { onMounted, ref } from "vue";
import "cropperjs";

import { useI18n } from "vue-i18n";
const { t, locale } = useI18n();

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

const emit = defineEmits(['crop'])

onMounted(() => {
  cro_sele.value.$change(0, 0, cor_size_width, cor_size_height);
  cro_canv.value.style.width = cor_size_width.toString() + "px";
  cro_canv.value.style.height = cor_size_height.toString() + "px";

  cro_imag.value.src = "";

  //console.log(cro_canv.value.clientHeight)
});

function initCropper(imageSrc) {
  is_have_URL.value = true;
  cro_imag.value.src = imageSrc;
}

function cancel() {
  is_have_URL.value = false;
}

function inputfile(e) {
  const file = e.target.files[0];
  if (!file) {
    e.target.value = "";
    is_have_URL.value = false;
    return;
  }

  if (!file.type.startsWith("image/")) {
    e.target.value = "";
    is_have_URL.value = false;
    return;
  }

  reader.readAsDataURL(file);
}

function openFilePicker() {
  const fileInput = document.createElement("input");
  fileInput.type = "file";
  fileInput.accept = "image/*"; // 可选：限制文件类型
  fileInput.multiple = false; // 可选：是否允许多选

  fileInput.onchange = (e) => {
    inputfile(e);
  };

  fileInput.click(); // 触发文件选择
}

function getsele() {
  const canvas = cro_canv.value.$toCanvas().then((a) => {
    //console.log(a);
    const imageData = a.toDataURL("image/png");

    emit('crop',imageData)
    //console.log(imageData);
  });
}
</script>

<template>
  <div class="d-flex flex-column flex-md-row">
    <div v-show="!is_have_URL" class="col-6 col-sm-4 col-md-2 col-xl py-3">
      <button class="btn btn-outline-primary" @click="openFilePicker">
        {{ t("cropper.select_image") }}
      </button>
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

    <div v-show="is_have_URL" class="thisbutton">
      <button class="btn btn-outline-primary" @click="getsele">
        {{ t("cropper.crop_image") }}
      </button>
      <button class="btn btn-outline-danger" @click="cancel">{{ t("cropper.closs") }}</button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.container {
  display: flex;
  /* 默认就是水平排列 */
  /* flex-direction: row; */
}

.thisbutton {
  display: flex;
  flex-direction: column;
  margin-left: 20px;
  margin-top: 20px;
  gap: 20px; /* 所有子元素之间的间距 */
}

.box {
  margin: 10px;
  flex-direction: column; /* 关键：改为纵向排列 */
}

.cropper-container {
  /* 四个角相同圆角 */
  border-radius: 10px;
  /* 基本描边 */
  //border: 2px solid #333;

  /* 基本阴影：x偏移 y偏移 模糊半径 扩展半径 颜色 */
  box-shadow: 5px 5px 10px rgba(0, 0, 0, 0.3);
}
</style>
