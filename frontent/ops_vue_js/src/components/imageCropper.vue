<template>
  <div class="image-cropper">
    <div class="cropper-actions">
      <input
        type="file"
        ref="fileInput"
        @change="handleFileSelect"
        accept="image/*"
        style="display: none"
      />
      <button @click="$refs.fileInput.click()">选择图片</button>

      <div v-if="haveImageSrc" class="action-buttons">
        <button @click="setAspectRatio(1)">1:1</button>
        <button @click="setAspectRatio(16 / 9)">16:9</button>
        <button @click="setAspectRatio(4 / 3)">4:3</button>
        <button @click="rotate(-90)">↺</button>
        <button @click="rotate(90)">↻</button>
        <button @click="crop">裁剪</button>
        <button @click="reset">重置</button>
      </div>
    </div>

    <div v-if="haveImageSrc" class="cropper-container">
      <img ref="imageEl" alt="裁剪图片" />
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from "vue";
import Cropper from "cropperjs";

const emit = defineEmits(["cropped", "error"]);

const imageEl = ref(null);
const fileInput = ref(null);
const haveImageSrc = ref(false);
const resultImage = ref("");
let cropperInstance = null;

// 读取文件为DataURL
const readFileAsDataURL = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => resolve(e.target.result);
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
};

// 选择文件
const handleFileSelect = async (event) => {
  const file = event.target.files[0];
  if (!file) return;

  if (!file.type.startsWith("image/")) {
    emit("error", "请选择图片文件");
    return;
  }
  console.log("选择了文件:", file.name);
  try {
    const dataUrl = await readFileAsDataURL(file);
    haveImageSrc.value = true;
    console.log("图片加载完成，初始化裁剪器");
    await nextTick();
    if (cropperInstance) {
      cropperInstance.destroy();
      console.log("销毁旧的Cropper实例");
    }
    console.log("加载图片进行裁剪");
    imageEl.value.src = dataUrl;
    
    cropperInstance = new Cropper(imageEl.value, {
      aspectRatio: 1,
      viewMode: 2,
      autoCropArea: 0.8,
      zoomable: true,
      zoomOnWheel: true,
      zoomOnTouch: true,
      wheelZoomRatio: 0.1,
      ready: function () {
        console.log("Cropper 初始化成功");
      },
      crop: function (event) {
        // 可以在这里获取裁剪区域的信息
        console.log(event.detail);
      },
    });
  } catch (error) {
    emit("error", "图片加载失败");
  }
};
</script>

<style scoped>
.image-cropper {
  max-width: 600px;
  margin: 0 auto;
}

.cropper-actions {
  margin-bottom: 20px;
}

.action-buttons {
  margin-top: 10px;
}

.action-buttons button {
  margin: 0 5px 5px 0;
  padding: 5px 10px;
}

.cropper-container {
  max-height: 400px;
  overflow: hidden;
}

.result-container {
  margin-top: 20px;
  text-align: center;
}

.result-container img {
  max-width: 300px;
  max-height: 300px;
  border: 1px solid #ddd;
  margin: 10px 0;
}
</style>
