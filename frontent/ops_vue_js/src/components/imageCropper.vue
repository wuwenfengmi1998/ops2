<script setup>

import { Modal } from "@tabler/core";
import { onMounted, ref } from "vue";

import  "@/libs/cropper.min.js"


const avata_toolt=ref()

const avatar_toolt_moda=ref()

const cropper=ref()

 const image = document.getElementById('cropper-image');
  // 初始化Cropper
  function initCropper(imageSrc) {
    if (cropper.value) {
      cropper.value.destroy();
    }

    image.src = imageSrc;
    cropper.value = new Cropper(image, {
      aspectRatio: 1,
      viewMode: 2,
      autoCropArea: 0.8,
      zoomable: true,
      zoomOnWheel: true,
      zoomOnTouch: true,
      wheelZoomRatio: 0.1,
      //minCanvasWidth: 400,
      //minCanvasHeight: 400,
      crop: updatePreview,
      ready() {
       
       
      }
    });
  }

  function inputfile(e){
    const file = e.target.files[0];
    if (!file) return;

    if (!file.type.startsWith('image/')) {
      showMessage('⚠️ 请选择有效的图片文件', 'error');
      return;
    }

    const reader = new FileReader();
    reader.onload = () => {
      initCropper(reader.result);
      currentScale = 1;
    };
    reader.readAsDataURL(file);
  }

onMounted(()=>{
  //console.log(avata_toolt)
  avatar_toolt_moda.value=new Modal(avata_toolt.value)
  avatar_toolt_moda.value.show()
})


</script>

<template>

  <div class="modal modal-blur fade" ref="avata_toolt" tabindex="-1" role="dialog" aria-hidden="true">
  <div class="modal-dialog modal-lg modal-dialog-centered" role="document">
    <div class="modal-content">
      <div class="container">
        <h1>头像裁剪工具</h1>

        <div class="flex-wrapper">
          <!-- 左侧裁剪区 -->
          <div class="crop-section">
            <div id="image-wrapper">
              <img id="cropper-image"
                src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=">
            </div>
            <!-- 上传进度 -->
            <div class="progress-container">
              <div class="progress-bar"></div>
            </div>

            <!-- 消息提示 -->
            <div id="message" class="alertavater"></div>

            <div class="preview-stats">
              <!-- <p>当前缩放: <span id="zoomValue">100%</span></p> -->
              <p>图片尺寸: <span id="imageSize">0 x 0</span></p>
            </div>
          </div>

          <!-- 右侧预览区 -->
          <div class="preview-section">
            <h3>实时预览</h3>
            <div class="preview-box">
              <img id="preview-img"
                src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=">
            </div>

            <!-- 控制按钮 -->
            <div class="controls">
              <label class="btn btn-primary">
                📁 选择图片
                <input type="file" accept="image/*" @change="inputfile">
              </label>
              <button class="btn btn-secondary " onclick="rotateImage(-90)">↩️ 左旋</button>
              <button class="btn btn-success " id="uploadBtn">✂️ 裁剪头像</button>

              <button class="btn btn-danger " style="margin-top: 150px;" onclick="avatar_toolt.hide()">❌ 取消</button>



            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
</template>

<style scoped>
/* 头像裁剪器样式*/

  .container {
    width: 95%;
    /* 改为百分比宽度 */
    margin: 20px auto;
    /* 增加上下边距 */
    max-width: 1200px;
    /* 保留最大宽度 */
    background: white;
    padding: 30px;
    border-radius: 12px;

  }

  .flex-wrapper {
    display: flex;
    gap: 30px;
    margin-top: 20px;
    flex-wrap: wrap;
    /* 添加换行支持 */
  }

  /* 裁剪区域 */
  .crop-section {
    flex: 1 1 60%;
    /* 弹性布局基础宽度 */
    min-width: 300px;
    /* 降低最小宽度 */
    height: auto;
    /* 移除固定高度 */
    min-height: 400px;
    /* 设置最小高度 */
  }

  #image-wrapper {
    width: 100%;
    height: 60vh;
    /* 改用视窗单位 */
    max-height: 600px;
    /* 设置最大高度 */
    background: #f8f9fa;
    border: 2px dashed #ddd;
    border-radius: 8px;
    overflow: hidden;
  }

  /* 预览区域自适应 */
  .preview-section {
    flex: 1 1 35%;
    /* 弹性布局基础宽度 */
    min-width: 250px;
    /* 设置合理最小宽度 */
  }

  /* 移动端适配 */
  @media (max-width: 768px) {
    .container {
      padding: 15px;
      /* 减少内边距 */
    }

    .flex-wrapper {
      flex-direction: column;
      /* 垂直排列 */
    }

    .crop-section,
    .preview-section {
      width: 100% !important;
      /* 强制全宽 */
      min-width: unset;
      /* 移除最小宽度 */
    }

    #image-wrapper {
      height: 50vh;
      /* 调整移动端高度 */
    }

    .preview-box {
      width: 120px;
      /* 缩小预览区域 */
      height: 120px;
    }

    .controls {
      flex-direction: column;
      /* 垂直排列按钮 */
      margin-top: 10px;
    }
  }


  #cropper-image {
    max-width: none !important;
    max-height: none !important;
  }

  /* 控制区域 */
  .controls {
    margin-top: 20px;
    display: flex;
    flex-direction: column;
    grid-template-columns: repeat(3, 1fr);
    gap: 10px;
  }

  /* 预览区域 */
  .preview-section {
    flex: 1;
    min-width: 50px;
  }

  .preview-box {
    width: 150px;
    height: 150px;
    /* border-radius: 50%; */
    border: 3px solid var(--primary-color);
    overflow: hidden;
    margin: 0 auto 20px;
  }

  #preview-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }



  /* 上传进度 */
  .progress-container {
    height: 8px;
    background: #eee;
    border-radius: 4px;
    margin-top: 20px;
    overflow: hidden;
    display: none;
  }

  .progress-bar {
    width: 0%;
    height: 100%;
    background: var(--primary-color);
    transition: width 0.3s ease;
  }

  /* 消息提示 */
  .alertavater {
    padding: 15px;
    border-radius: 6px;
    margin-top: 20px;
    display: none;
  }

  .alertavater-success {
    background: #dff0d8;
    color: #3c763d;
  }

  .alertavater-error {
    background: #f2dede;
    color: var(--error-color);
  }
  input[type="file"] {
    display: none;
  }


</style>