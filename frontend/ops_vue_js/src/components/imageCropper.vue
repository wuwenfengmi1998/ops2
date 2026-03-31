<script setup>
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
reader.onload = () => { initCropper(reader.result); };
const emit = defineEmits(['crop-data-url', 'crop-error'])
onMounted(() => {
  cro_sele.value.$change(0, 0, cor_size_width, cor_size_height);
  cro_canv.value.style.width = cor_size_width.toString() + "px";
  cro_canv.value.style.height = cor_size_height.toString() + "px";
  cro_imag.value.src = "";
});
function initCropper(imageSrc) { is_have_URL.value = true; cro_imag.value.src = imageSrc; }
function cancel() { is_have_URL.value = false; }
function inputfile(e) {
  const file = e.target.files[0]; if (!file) { e.target.value = ""; is_have_URL.value = false; return; }
  if (!file.type.startsWith("image/")) { e.target.value = ""; is_have_URL.value = false; return; }
  reader.readAsDataURL(file);
}
function openFilePicker() {
  const fileInput = document.createElement("input");
  fileInput.type = "file"; fileInput.accept = "image/*"; fileInput.multiple = false;
  fileInput.onchange = (e) => { inputfile(e); };
  fileInput.click();
}
function getsele() {
  console.log('getsele called, checking elements:', {
    cro_canv: cro_canv.value,
    cro_sele: cro_sele.value,
    cro_imag: cro_imag.value,
  })
  
  // 使用async函数处理异步裁剪逻辑
  const cropImage = async () => {
    try {
      console.log('Starting crop process')
      
      // 方法1: 尝试调用原生的toCanvas方法（如果可用）
      if (cro_canv.value && typeof cro_canv.value.$toCanvas === 'function') {
        console.log('Using $toCanvas method')
        try {
          const result = await cro_canv.value.$toCanvas()
          console.log('$toCanvas result type:', typeof result, 'value:', result)
          
          // 检查结果类型，如果是Canvas元素则转换为data URL
          if (result) {
            let dataUrl
            if (result instanceof HTMLCanvasElement) {
              // 如果是Canvas对象，转换为data URL
              console.log('Converting HTMLCanvasElement to data URL')
              dataUrl = result.toDataURL('image/jpeg', 0.9)
            } else if (typeof result === 'string' && result.startsWith('data:image/')) {
              // 如果已经是data URL，直接使用
              dataUrl = result
            } else {
              // 其他情况，不直接发送
              console.warn('$toCanvas returned unexpected type:', typeof result)
              throw new Error('Unexpected return type from $toCanvas')
            }
            
            if (dataUrl) {
              console.log('Generated data URL from $toCanvas, length:', dataUrl.length)
              emit('crop-data-url', dataUrl)
              return  // 成功，结束函数
            }
          }
          console.log('$toCanvas returned empty or invalid result, falling back to manual crop')
        } catch (error) {
          console.warn('$toCanvas failed, using manual crop:', error.message)
        }
      }
      
      // 方法2: 使用手动裁剪方法
      console.log('Falling back to manual crop method')
      
      // 获取图像元素
      const img = cro_imag.value
      if (!img) {
        console.error('No image element found')
        emit('crop-error', '未找到图像')
        return
      }
      
      // 等待图像加载完成
      if (!img.complete) {
        console.log('Waiting for image to load...')
        await new Promise((resolve) => {
          img.onload = resolve
          img.onerror = resolve  // 即使加载失败也继续
          // 设置超时，防止无限等待
          setTimeout(resolve, 3000)
        })
      }
      
      if (!img.complete || img.naturalWidth === 0) {
        console.error('Image failed to load or has 0 dimensions')
        emit('crop-error', '图像加载失败')
        return
      }
      
      console.log('Image loaded successfully:', img.naturalWidth, 'x', img.naturalHeight)
      
      // 获取canvas的尺寸
      const canvasRect = cro_canv.value?.getBoundingClientRect?.() || {
        width: cor_size_width,
        height: cor_size_height,
        left: 0,
        top: 0
      }
      
      console.log('Canvas rectangle:', canvasRect)
      
      // 获取选择区域
      let selectionRect = null
      if (cro_sele.value && typeof cro_sele.value.$getRect === 'function') {
        selectionRect = cro_sele.value.$getRect()
        console.log('Selection rect:', selectionRect)
      }
      
      // 验证选择区域
      if (!selectionRect || !selectionRect.width || !selectionRect.height) {
        selectionRect = {
          left: 0,
          top: 0,
          width: canvasRect.width,
          height: canvasRect.height
        }
        console.log('Using entire canvas as selection:', selectionRect)
      }
      
      // 创建输出canvas
      const outputCanvas = document.createElement('canvas')
      outputCanvas.width = selectionRect.width
      outputCanvas.height = selectionRect.height
      const ctx = outputCanvas.getContext('2d')
      
      // 设置白色背景
      ctx.fillStyle = '#ffffff'
      ctx.fillRect(0, 0, outputCanvas.width, outputCanvas.height)
      
      // 计算图像在canvas中的显示方式
      const imgAspect = img.naturalWidth / img.naturalHeight
      const canvasAspect = canvasRect.width / canvasRect.height
      
      let drawWidth, drawHeight, drawX, drawY
      
      if (imgAspect > canvasAspect) {
        // 图像更宽，高度适配
        drawHeight = canvasRect.height
        drawWidth = canvasRect.height * imgAspect
        drawX = (canvasRect.width - drawWidth) / 2
        drawY = 0
      } else {
        // 图像更高，宽度适配
        drawWidth = canvasRect.width
        drawHeight = canvasRect.width / imgAspect
        drawX = 0
        drawY = (canvasRect.height - drawHeight) / 2
      }
      
      console.log('Image display info:', { drawX, drawY, drawWidth, drawHeight })
      
      // 计算裁剪区域（将选择区域转换到图像坐标）
      const cropX = Math.max(0, (selectionRect.left - drawX) / drawWidth * img.naturalWidth)
      const cropY = Math.max(0, (selectionRect.top - drawY) / drawHeight * img.naturalHeight)
      const cropWidth = Math.min(
        (selectionRect.width / drawWidth) * img.naturalWidth,
        img.naturalWidth - cropX
      )
      const cropHeight = Math.min(
        (selectionRect.height / drawHeight) * img.naturalHeight,
        img.naturalHeight - cropY
      )
      
      console.log('Final crop coordinates (image space):', {
        cropX, cropY, cropWidth, cropHeight
      })
      
      // 执行裁剪
      ctx.drawImage(
        img,
        cropX, cropY, cropWidth, cropHeight,      // 源图像裁剪区域
        0, 0, outputCanvas.width, outputCanvas.height  // 目标canvas区域
      )
      
      // 生成data URL（JPEG格式，质量0.9）
      const dataUrl = outputCanvas.toDataURL('image/jpeg', 0.9)
      console.log('Generated crop data URL, length:', dataUrl.length)
      
      // 确保我们发送的是有效的data URL字符串
      if (typeof dataUrl === 'string' && dataUrl.startsWith('data:image/')) {
        emit('crop-data-url', dataUrl)
      } else {
        console.error('Invalid data URL generated:', typeof dataUrl)
        emit('crop-error', '生成的数据URL无效')
      }
      
    } catch (error) {
      console.error('Crop process error:', error)
      emit('crop-error', '裁剪过程中发生错误：' + error.message)
    }
  }
  
  // 执行裁剪
  cropImage()
}
</script>

<template>
  <div class="space-y-4">
    <!-- Header and Instruction -->
    <div class="flex items-center justify-between">
      <div class="space-y-1">
        <h4 class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('cropper.upload_image') }}</h4>
        <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('cropper.supported_formats') }}</p>
      </div>
      <button 
        v-show="!is_have_URL"
        class="group flex items-center gap-2 rounded-lg bg-gradient-to-r from-blue-600 to-blue-500 px-4 py-2.5 text-sm font-medium text-white shadow-sm transition-all hover:from-blue-700 hover:to-blue-600 hover:shadow-md focus:outline-none focus:ring-2 focus:ring-blue-500/30"
        @click="openFilePicker"
      >
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
        </svg>
        {{ t("cropper.select_image") }}
      </button>
    </div>

    <!-- Cropper Area -->
    <div v-show="is_have_URL" class="rounded-xl border border-gray-200 bg-gray-50/50 p-4 dark:border-gray-800 dark:bg-gray-900/30">
      <cropper-canvas 
        ref="cro_canv" 
        class="cropper-container mx-auto"
        :hidden="!is_have_URL" 
        background 
        scale-step="0.1"
      >
        <cropper-image ref="cro_imag" src="" alt="Picture" initialCenterSize="cover" rotatable scalable skewable translatable></cropper-image>
        <cropper-shade hidden></cropper-shade>
        <cropper-handle action="move" plain></cropper-handle>
        <cropper-selection ref="cro_sele">
          <cropper-grid role="grid" covered></cropper-grid>
          <cropper-crosshair centered></cropper-crosshair>
          <cropper-handle action="move" theme-color="rgba(255, 255, 255, 0)"></cropper-handle>
        </cropper-selection>
      </cropper-canvas>
      
      <!-- Cropper Actions -->
      <div class="mt-4 flex flex-wrap items-center justify-between gap-3">
        <div class="text-xs text-gray-500 dark:text-gray-400">
          <span class="inline-flex items-center gap-1">
            <svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v3.586L7.707 9.293a1 1 0 00-1.414 1.414l3 3a1 1 0 001.414 0l3-3a1 1 0 00-1.414-1.414L11 10.586V7z" clip-rule="evenodd" />
            </svg>
            {{ t('cropper.drag_to_resize') }}
          </span>
        </div>
        
        <div class="flex gap-2">
          <button 
            class="group flex items-center gap-2 rounded-lg bg-gradient-to-r from-blue-600 to-blue-500 px-4 py-2 text-sm font-medium text-white shadow-sm transition-all hover:from-blue-700 hover:to-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500/30"
            @click="getsele"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            {{ t("cropper.crop_image") }}
          </button>
          <button 
            class="rounded-lg border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 hover:border-gray-400 dark:border-gray-700 dark:text-gray-300 dark:hover:bg-gray-800 dark:hover:border-gray-600"
            @click="cancel"
          >
            {{ t("cropper.closs") }}
          </button>
        </div>
      </div>
    </div>

    <!-- Preview Area (when image selected but not cropped) -->
    <div v-show="!is_have_URL" class="rounded-lg border-2 border-dashed border-gray-300 bg-gray-50/50 p-8 text-center dark:border-gray-700 dark:bg-gray-900/20">
      <div class="mx-auto max-w-xs">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
        <h4 class="mt-4 text-sm font-medium text-gray-900 dark:text-white">{{ t('cropper.ready_to_upload') }}</h4>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
          {{ t('cropper.click_button_or_drag') }}
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.cropper-container {
  width: 100%;
  max-width: 400px;
  height: 300px;
  border-radius: 12px;
  box-shadow: 0 8px 25px -5px rgba(0, 0, 0, 0.08), 0 4px 10px -2px rgba(0, 0, 0, 0.04);
  overflow: hidden;
  border: 1px solid rgba(0, 0, 0, 0.05);
  transition: box-shadow 0.2s ease;
}

.cropper-container:hover {
  box-shadow: 0 12px 32px -8px rgba(0, 0, 0, 0.12), 0 6px 14px -4px rgba(0, 0, 0, 0.06);
}

.dark .cropper-container {
  box-shadow: 0 8px 25px -5px rgba(0, 0, 0, 0.25), 0 4px 10px -2px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.dark .cropper-container:hover {
  box-shadow: 0 12px 32px -8px rgba(0, 0, 0, 0.35), 0 6px 14px -4px rgba(0, 0, 0, 0.2);
}
</style>
