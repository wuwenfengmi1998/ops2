<template>
  <div class="avatar-cropper">
    <!-- 上传区域 -->
    <div v-if="!imageSrc" class="upload-area" @click="triggerFileInput">
      <div class="upload-content">
        <div class="upload-icon">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M14 2H6C4.9 2 4 2.9 4 4V20C4 21.1 4.9 22 6 22H18C19.1 22 20 21.1 20 20V8L14 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M14 2V8H20" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M16 13H8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M16 17H8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M10 9H8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <h3>点击上传图片</h3>
        <p>支持 JPG、PNG 格式，建议尺寸 200×200 像素以上</p>
      </div>
      <input 
        ref="fileInput"
        type="file" 
        accept="image/*" 
        @change="handleFileSelect"
        style="display: none"
      >
    </div>

    <!-- 裁剪区域 -->
    <div v-else class="cropper-container">
      <div class="cropper-header">
        <h3>调整头像</h3>
        <button class="btn-close" @click="reset">×</button>
      </div>
      
      <div class="cropper-content">
        <!-- 裁剪画布 -->
        <div class="canvas-wrapper">
            <div class="canvas-container" ref="canvasContainer">
              <canvas ref="canvas" class="cropper-canvas"></canvas>
              <div 
                class="crop-box"
                :style="cropBoxStyle"
                @mousedown="startDrag"
                @touchstart="startDrag"
              >
                <div class="crop-box-handle" v-for="handle in handles" :key="handle.position"
                  :class="`handle-${handle.position}`"
                  @mousedown="startResize(handle.position, $event)"
                  @touchstart="startResize(handle.position, $event)"
                ></div>
              </div>
            </div>
        </div>

        <!-- 预览区域 -->
        <div class="preview-section">
          <div class="preview-container">
            <div class="preview-item">
              <h4>圆形预览</h4>
              <div class="preview-circle">
                <img v-if="previewUrl" :src="previewUrl" alt="头像预览" class="preview-image">
                <div v-else class="preview-placeholder">预览</div>
              </div>
            </div>
            <div class="preview-item">
              <h4>方形预览</h4>
              <div class="preview-square">
                <img v-if="previewUrl" :src="previewUrl" alt="头像预览" class="preview-image">
                <div v-else class="preview-placeholder">预览</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 控制面板 -->
      <div class="control-panel">
        <div class="control-group">
          <label>缩放: {{ Math.round(scale * 100) }}%</label>
          <input 
            type="range" 
            min="10" 
            max="200" 
            step="1" 
            v-model.number="scale"
            class="zoom-slider"
          >
        </div>

        <div class="control-group">
          <label>旋转: {{ rotation }}°</label>
          <input 
            type="range" 
            min="0" 
            max="360" 
            step="1" 
            v-model.number="rotation"
            class="rotate-slider"
          >
        </div>

        <div class="button-group">
          <button class="btn-secondary" @click="reset">重新选择</button>
          <button class="btn-primary" @click="cropAvatar">确认裁剪</button>
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-overlay">
      <div class="loading-spinner"></div>
      <p>处理中...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'

// 响应式数据
const fileInput = ref(null)
const canvas = ref(null)
const canvasContainer = ref(null)
const imageSrc = ref('')
const previewUrl = ref('')
const loading = ref(false)

// 裁剪状态
const state = reactive({
  scale: 100,
  rotation: 0,
  position: { x: 0, y: 0 },
  cropSize: 200,
  isDragging: false,
  isResizing: false,
  resizeDirection: null,
  dragStart: { x: 0, y: 0 }
})

// 计算属性
const scale = computed({
  get: () => state.scale,
  set: (value) => {
    state.scale = Math.max(10, Math.min(200, value))
    redrawCanvas()
  }
})

const rotation = computed({
  get: () => state.rotation,
  set: (value) => {
    state.rotation = value % 360
    redrawCanvas()
  }
})

const cropBoxStyle = computed(() => ({
  width: `${state.cropSize}px`,
  height: `${state.cropSize}px`,
  left: `${state.position.x}px`,
  top: `${state.position.y}px`,
  transform: `rotate(${state.rotation}deg)`
}))

// 控制点配置
const handles = [
  { position: 'nw' }, { position: 'n' }, { position: 'ne' },
  { position: 'w' }, { position: 'e' },
  { position: 'sw' }, { position: 's' }, { position: 'se' }
]

// 方法
const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (!file) return

  if (!file.type.startsWith('image/')) {
    alert('请选择图片文件')
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    imageSrc.value = e.target.result
    nextTick(() => {
      initializeCropper()
    })
  }
  reader.readAsDataURL(file)
}

const initializeCropper = () => {
  if (!canvas.value || !canvasContainer.value) return

  const containerRect = canvasContainer.value.getBoundingClientRect()
  const containerSize = Math.min(containerRect.width, containerRect.height) - 40
  
  // 设置画布尺寸
  canvas.value.width = containerSize
  canvas.value.height = containerSize
  
  // 初始化裁剪框位置
  state.cropSize = Math.min(containerSize * 0.6, 300)
  state.position.x = (containerSize - state.cropSize) / 2
  state.position.y = (containerSize - state.cropSize) / 2
  state.scale = 100
  state.rotation = 0

  redrawCanvas()
}

const redrawCanvas = () => {
  if (!canvas.value || !imageSrc.value) return

  const ctx = canvas.value.getContext('2d')
  const size = canvas.value.width
  
  // 清空画布
  ctx.clearRect(0, 0, size, size)
  
  // 绘制背景
  ctx.fillStyle = '#f5f5f5'
  ctx.fillRect(0, 0, size, size)
  
  // 创建离屏canvas用于图片变换
  const offscreen = document.createElement('canvas')
  const offscreenCtx = offscreen.getContext('2d')
  
  const img = new Image()
  img.onload = () => {
    // 计算缩放后的尺寸
    const scaleFactor = state.scale / 100
    const imgWidth = img.width * scaleFactor
    const imgHeight = img.height * scaleFactor
    
    // 设置离屏canvas尺寸
    offscreen.width = Math.max(imgWidth, imgHeight) * 2
    offscreen.height = Math.max(imgWidth, imgHeight) * 2
    
    // 在离屏canvas上绘制并旋转图片
    offscreenCtx.save()
    offscreenCtx.translate(offscreen.width / 2, offscreen.height / 2)
    offscreenCtx.rotate(state.rotation * Math.PI / 180)
    offscreenCtx.drawImage(img, -imgWidth / 2, -imgHeight / 2, imgWidth, imgHeight)
    offscreenCtx.restore()
    
    // 在主画布上绘制
    ctx.drawImage(offscreen, 
      (size - offscreen.width) / 2, 
      (size - offscreen.height) / 2
    )
    
    updatePreview()
  }
  img.src = imageSrc.value
}

const updatePreview = () => {
  if (!canvas.value) return

  const previewCanvas = document.createElement('canvas')
  previewCanvas.width = state.cropSize
  previewCanvas.height = state.cropSize
  const previewCtx = previewCanvas.getContext('2d')
  
  // 绘制裁剪区域到预览canvas
  previewCtx.drawImage(
    canvas.value,
    state.position.x, state.position.y, state.cropSize, state.cropSize,
    0, 0, state.cropSize, state.cropSize
  )
  
  previewUrl.value = previewCanvas.toDataURL('image/png')
}

const startDrag = (e) => {
  console.log('startDrag', e)
  e.preventDefault()
  state.isDragging = true
  const clientX = e.type.includes('touch') ? e.touches[0].clientX : e.clientX
  const clientY = e.type.includes('touch') ? e.touches[0].clientY : e.clientY
  
  state.dragStart.x = clientX - state.position.x
  state.dragStart.y = clientY - state.position.y
  
  document.addEventListener('mousemove', onDrag)
  document.addEventListener('mouseup', stopDrag)
  document.addEventListener('touchmove', onDrag)
  document.addEventListener('touchend', stopDrag)
}

const onDrag = (e) => {
  if (!state.isDragging || !canvas.value) return
  
  e.preventDefault()
  const clientX = e.type.includes('touch') ? e.touches[0].clientX : e.clientX
  const clientY = e.type.includes('touch') ? e.touches[0].clientY : e.clientY
  
  const containerRect = canvas.value.getBoundingClientRect()
  const newX = clientX - containerRect.left - state.dragStart.x
  const newY = clientY - containerRect.top - state.dragStart.y
  
  // 限制在画布范围内
  state.position.x = Math.max(0, Math.min(canvas.value.width - state.cropSize, newX))
  state.position.y = Math.max(0, Math.min(canvas.value.height - state.cropSize, newY))
  
  updatePreview()
}

const stopDrag = () => {
  state.isDragging = false
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
  document.removeEventListener('touchmove', onDrag)
  document.removeEventListener('touchend', stopDrag)
}

const startResize = (direction, e) => {
  e.preventDefault()
  e.stopPropagation()
  
  state.isResizing = true
  state.resizeDirection = direction
  
  const clientX = e.type.includes('touch') ? e.touches[0].clientX : e.clientX
  const clientY = e.type.includes('touch') ? e.touches[0].clientY : e.clientY
  
  state.dragStart.x = clientX
  state.dragStart.y = clientY
  state.dragStart.width = state.cropSize
  state.dragStart.position = { ...state.position }
  
  document.addEventListener('mousemove', onResize)
  document.addEventListener('mouseup', stopResize)
  document.addEventListener('touchmove', onResize)
  document.addEventListener('touchend', stopResize)
}

const onResize = (e) => {
  if (!state.isResizing || !canvas.value) return
  
  e.preventDefault()
  const clientX = e.type.includes('touch') ? e.touches[0].clientX : e.clientX
  const clientY = e.type.includes('touch') ? e.touches[0].clientY : e.clientY
  
  const deltaX = clientX - state.dragStart.x
  const deltaY = clientY - state.dragStart.y
  const minSize = 50
  const maxSize = Math.min(canvas.value.width, canvas.value.height)
  
  let newSize = state.dragStart.width
  let newX = state.dragStart.position.x
  let newY = state.dragStart.position.y
  
  switch (state.resizeDirection) {
    case 'e':
      newSize = Math.max(minSize, Math.min(maxSize, state.dragStart.width + deltaX))
      break
    case 'w':
      newSize = Math.max(minSize, Math.min(maxSize, state.dragStart.width - deltaX))
      newX = state.dragStart.position.x + deltaX
      break
    case 's':
      newSize = Math.max(minSize, Math.min(maxSize, state.dragStart.width + deltaY))
      break
    case 'n':
      newSize = Math.max(minSize, Math.min(maxSize, state.dragStart.width - deltaY))
      newY = state.dragStart.position.y + deltaY
      break
    case 'se':
      newSize = Math.max(minSize, Math.min(maxSize, state.dragStart.width + Math.max(deltaX, deltaY)))
      break
    case 'sw':
      newSize = Math.max(minSize, Math.min(maxSize, state.dragStart.width + Math.max(-deltaX, deltaY)))
      newX = state.dragStart.position.x + deltaX
      break
    case 'ne':
      newSize = Math.max(minSize, Math.min(maxSize, state.dragStart.width + Math.max(deltaX, -deltaY)))
      newY = state.dragStart.position.y + deltaY
      break
    case 'nw':
      newSize = Math.max(minSize, Math.min(maxSize, state.dragStart.width + Math.max(-deltaX, -deltaY)))
      newX = state.dragStart.position.x + deltaX
      newY = state.dragStart.position.y + deltaY
      break
  }
  
  // 限制位置在画布范围内
  state.cropSize = newSize
  state.position.x = Math.max(0, Math.min(canvas.value.width - newSize, newX))
  state.position.y = Math.max(0, Math.min(canvas.value.height - newSize, newY))
  
  updatePreview()
}

const stopResize = () => {
  state.isResizing = false
  state.resizeDirection = null
  document.removeEventListener('mousemove', onResize)
  document.removeEventListener('mouseup', stopResize)
  document.removeEventListener('touchmove', onResize)
  document.removeEventListener('touchend', stopResize)
}

const cropAvatar = async () => {
  loading.value = true
  
  try {
    // 模拟处理时间
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 创建最终裁剪结果
    const resultCanvas = document.createElement('canvas')
    resultCanvas.width = 400 // 输出尺寸
    resultCanvas.height = 400
    const resultCtx = resultCanvas.getContext('2d')
    
    // 绘制圆形遮罩
    resultCtx.beginPath()
    resultCtx.arc(200, 200, 200, 0, Math.PI * 2)
    resultCtx.closePath()
    resultCtx.clip()
    
    // 绘制图片
    const img = new Image()
    img.onload = () => {
      resultCtx.drawImage(img, 0, 0, 400, 400)
      
      // 触发裁剪完成事件
      const dataUrl = resultCanvas.toDataURL('image/png')
      emit('cropped', dataUrl)
      
      loading.value = false
    }
    img.src = previewUrl.value
    
  } catch (error) {
    console.error('裁剪失败:', error)
    loading.value = false
  }
}

const reset = () => {
  imageSrc.value = ''
  previewUrl.value = ''
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

// 组件挂载
onMounted(() => {
  window.addEventListener('resize', initializeCropper)
})

onUnmounted(() => {
  window.removeEventListener('resize', initializeCropper)
})

// 监听图片源变化
watch(imageSrc, () => {
  if (imageSrc.value) {
    nextTick(() => {
      initializeCropper()
    })
  }
})

// 定义组件事件
const emit = defineEmits(['cropped'])
</script>

<style scoped>
.avatar-cropper {
  max-width: 800px;
  margin: 0 auto;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

/* 上传区域 */
.upload-area {
  border: 2px dashed #d1d5db;
  border-radius: 12px;
  padding: 60px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #fafafa;
}

.upload-area:hover {
  border-color: #3b82f6;
  background: #f8fafc;
}

.upload-content .upload-icon {
  color: #9ca3af;
  margin-bottom: 16px;
}

.upload-content h3 {
  font-size: 18px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 8px;
}

.upload-content p {
  color: #6b7280;
  font-size: 14px;
}

/* 裁剪容器 */
.cropper-container {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

.cropper-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e5e7eb;
  background: #f9fafb;
}

.cropper-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #111827;
}

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #6b7280;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.2s;
}

.btn-close:hover {
  background: #f3f4f6;
  color: #374151;
}

.cropper-content {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 24px;
  padding: 24px;
}

/* 画布区域 */
.canvas-wrapper {
  position: relative;
}

.canvas-container {
  position: relative;
  width: 100%;
  height: 400px;
  background: #f9fafb;
  border-radius: 8px;
  overflow: hidden;
}

.cropper-canvas {
  width: 100%;
  height: 100%;
  display: block;
}

.crop-box {
  position: absolute;
  border: 2px solid #3b82f6;
  box-shadow: 0 0 0 9999px rgba(0, 0, 0, 0.3);
  cursor: move;
  z-index: 10;
}

.crop-box-handle {
  position: absolute;
  width: 12px;
  height: 12px;
  background: #3b82f6;
  border: 2px solid white;
  border-radius: 2px;
  z-index: 20;
}

.handle-nw { top: -6px; left: -6px; cursor: nw-resize; }
.handle-n { top: -6px; left: 50%; margin-left: -6px; cursor: n-resize; }
.handle-ne { top: -6px; right: -6px; cursor: ne-resize; }
.handle-w { top: 50%; left: -6px; margin-top: -6px; cursor: w-resize; }
.handle-e { top: 50%; right: -6px; margin-top: -6px; cursor: e-resize; }
.handle-sw { bottom: -6px; left: -6px; cursor: sw-resize; }
.handle-s { bottom: -6px; left: 50%; margin-left: -6px; cursor: s-resize; }
.handle-se { bottom: -6px; right: -6px; cursor: se-resize; }

/* 预览区域 */
.preview-section {
  display: flex;
  flex-direction: column;
}

.preview-container {
  display: flex;
  gap: 20px;
  flex-direction: column;
}

.preview-item h4 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  text-align: center;
}

.preview-circle {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  overflow: hidden;
  background: #f3f4f6;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
}

.preview-square {
  width: 120px;
  height: 120px;
  border-radius: 8px;
  overflow: hidden;
  background: #f3f4f6;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-placeholder {
  color: #9ca3af;
  font-size: 14px;
}

/* 控制面板 */
.control-panel {
  grid-column: 1 / -1;
  padding: 24px;
  border-top: 1px solid #e5e7eb;
  background: #f9fafb;
}

.control-group {
  margin-bottom: 20px;
}

.control-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.zoom-slider, .rotate-slider {
  width: 100%;
  height: 6px;
  border-radius: 3px;
  background: #e5e7eb;
  outline: none;
  -webkit-appearance: none;
}

.zoom-slider::-webkit-slider-thumb,
.rotate-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #3b82f6;
  cursor: pointer;
}

.button-group {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.btn-primary, .btn-secondary {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background: #2563eb;
}

.btn-secondary {
  background: white;
  color: #374151;
  border: 1px solid #d1d5db;
}

.btn-secondary:hover {
  background: #f9fafb;
  border-color: #9ca3af;
}

/* 加载状态 */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: white;
  z-index: 1000;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(255, 255, 255, 0.3);
  border-left: 4px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .cropper-content {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .preview-container {
    flex-direction: row;
    justify-content: center;
  }
  
  .canvas-container {
    height: 300px;
  }
  
  .button-group {
    flex-direction: column;
  }
}
</style>