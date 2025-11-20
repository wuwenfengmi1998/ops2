<template>
  <div class="basic_container">
    <div class="tool_wrap">
      <button :class="{ active: isCropperMove }" @click="handleMove">
        移动
      </button>
      <button @click="handleRotate">旋转</button>
      <button :class="{ active: isCropperSelection }" @click="handleCropper">
        {{ isCropperSelection ? '重置选区' : '裁剪' }}
      </button>
    </div>
    <div class="dialog_wrap">
      <div class="image_wrap" ref="imageWrap">
        <cropper-canvas ref="croppercanvas" background>
          <cropper-image
            :src="fileObj.fileShow"
            alt="Picture"
            ref="cropperimage"
            rotatable
            scalable
            skewable
            translatable
          ></cropper-image>
          <cropper-shade hidden ref="cropperShade"></cropper-shade>
          <cropper-handle :action="currentType" plain></cropper-handle>
          <cropper-selection
            id="cropperSelection"
            ref="cropperselection"
            movable
            resizable
            hidden
            outlined
            @change="onCropperSelectionChange"
          >
            <cropper-crosshair centered />
            <cropper-handle
              action="move"
              theme-color="rgba(255, 255, 255, 0.35)"
            />
            <cropper-handle action="n-resize" />
            <cropper-handle action="e-resize" />
            <cropper-handle action="s-resize" />
            <cropper-handle action="w-resize" />
            <cropper-handle action="ne-resize" />
            <cropper-handle action="nw-resize" />
            <cropper-handle action="se-resize" />
            <cropper-handle action="sw-resize" />
          </cropper-selection>
        </cropper-canvas>
      </div>
      <div class="info_wrap">
        <div class="cropper_preview">
          <cropper-viewer
            selection="#cropperSelection"
            style="width: 200px"
          ></cropper-viewer>
        </div>
        <div class="btn_wrap">
          <input type="file" ref="input_form" @change="handleUploadSuccess" />
          <button type="primary" @click="handleConfirm">确 认</button>
        </div>
        点击确认后，看控制台，有信息
      </div>
    </div>
  </div>
</template>

<script setup>
import 'cropperjs';
import { computed, ref } from 'vue';

const fileObj = ref({});

const croppercanvas = ref();
const cropperimage = ref();
const cropperselection = ref();

/**
 * 选区逻辑
 */
// 是否正在开始选区
const isCropperSelection = ref(false);
const isCropperMove = ref(true);

//  判断当前是移动还是选区
const currentType = computed(() => (isCropperMove.value ? 'move' : 'select'));

/**
 * 按钮方法
 */
// 旋转
function handleRotate() {
  cropperimage.value.$rotate('90deg');
  cropperimage.value.$center('contain');
}
// 裁剪
function handleCropper() {
  isCropperMove.value = false;
  if (isCropperMove.value) {
    cropperselection.value.$clear();
  } else {
    const cropperCanvas = croppercanvas.value;
    const cropperCanvasRect = cropperCanvas.getBoundingClientRect();

    const cropperImage = cropperimage.value;
    const cropperImageRect = cropperImage.getBoundingClientRect();
    const maxSelection = {
      x: cropperImageRect.left - cropperCanvasRect.left,
      y: cropperImageRect.top - cropperCanvasRect.top,
      width: cropperImageRect.width,
      height: cropperImageRect.height,
    };
    cropperselection.value.$change(
      maxSelection.x,
      maxSelection.y,
      maxSelection.width,
      maxSelection.height,
    );
  }
}
// 移动
function handleMove() {
  if (!isCropperMove.value) {
    isCropperMove.value = true;
    // 如果想要点击移动，清除选区，可以打开下面的代码注释
    // cropperselection.value.$clear();
  }
}

/**
 * 监听选择区变化
 * @param event
 */
function onCropperSelectionChange(event) {
  if (event.detail.width && event.detail.height) {
    isCropperSelection.value = true;
  } else {
    isCropperSelection.value = false;
  }
}

/**
 * 确认裁剪
 */
const emit = defineEmits(['success']);
async function handleConfirm() {
  if (isCropperSelection.value) {
    const res = await cropperselection.value.$toCanvas();

    const dataImage = res.toDataURL('image/png');
    const file = dataURLtoFile(dataImage, fileObj.value.name);
    emit('success', {
      ...fileObj.value,
      file: file,
      fileShow: dataImage,
    });
  }
}
// 将data:image转成新的file
function dataURLtoFile(dataurl, filename) {
  var arr = dataurl.split(','),
    mime = arr[0].match(/:(.*?);/)[1],
    bstr = atob(arr[1]),
    n = bstr.length,
    u8arr = new Uint8Array(n);
  while (n--) {
    u8arr[n] = bstr.charCodeAt(n);
  }
  const blob = new Blob([u8arr], { type: mime });
  const file = new File([blob], filename, { type: mime });
  return file;
}

/**
 * 文件上传
 */
const input_form = ref();
function handleUploadSuccess() {
  const files = input_form.value.files;

  if (files.length) {
    fileObj.value = {
      name: files[0].name,
      file: files[0],
      fileShow: URL.createObjectURL(files[0]),
    };
  }
}
</script>

<style scoped>
.dialog_wrap {
  display: flex;
  .image_wrap {
    width: 400px;
    height: 300px;
    flex-shrink: 0;

    cropper-canvas {
      width: 100%;
      height: 100%;
    }
  }
  .info_wrap {
    margin-left: 20px;
  }
}
button {
  & + button {
    margin-left: 20px;
  }
}
button.active {
  background-color: #c6dff8;
  border-color: #409eff;
}
</style>
