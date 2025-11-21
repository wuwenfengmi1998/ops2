<template>
  <div class="cropper-demo">
    <div class="container">
      <header>
        <h1>CropperJS 组件演示</h1>
        <p class="description">
          此演示展示了CropperJS各个组件的功能和使用方法。您可以通过下方控制面板调整裁剪区域、查看预览效果，并了解每个组件的作用。
        </p>
      </header>
      
      <div class="demo-area">
        <div class="cropper-container">
          <h3>裁剪画布</h3>
          <div class="canvas-wrapper">
            <canvas 
              ref="canvas" 
              class="cropper-canvas"
              @mousedown="onCanvasMouseDown"
              @mousemove="onCanvasMouseMove"
              @mouseup="onCanvasMouseUp"
              @mouseleave="onCanvasMouseLeave"
            ></canvas>
          </div>
          
          <div class="viewer-container">
            <div class="viewer-label">实时预览</div>
            <canvas ref="previewCanvas" class="preview-viewer"></canvas>
          </div>
          
          <div class="status-bar">
            <div>选区位置: <span>{{ positionInfo }}</span></div>
            <div>选区尺寸: <span>{{ sizeInfo }}</span></div>
            <div>宽高比: <span>{{ ratioInfo }}</span></div>
          </div>
        </div>
        
        <div class="controls">
          <h3>控制面板</h3>
          
          <div class="control-group">
            <h4>选区操作</h4>
            <div class="btn-group">
              <button 
                :class="{ active: interactionMode === 'move' }"
                @click="setInteractionMode('move')"
              >
                移动选区
              </button>
              <button 
                :class="{ active: interactionMode === 'resize' }"
                @click="setInteractionMode('resize')"
              >
                调整大小
              </button>
              <button @click="rotateSelection">旋转</button>
            </div>
          </div>
          
          <div class="control-group">
            <h4>宽高比设置</h4>
            <div class="btn-group">
              <button 
                v-for="ratio in aspectRatios" 
                :key="ratio.value"
                @click="setAspectRatio(ratio.value)"
              >
                {{ ratio.label }}
              </button>
            </div>
          </div>
          
          <div class="control-group">
            <h4>遮罩设置</h4>
            <div class="slider-container">
              <label for="opacitySlider">遮罩透明度: <span>{{ shade.opacity }}</span></label>
              <input 
                type="range" 
                id="opacitySlider" 
                min="0" 
                max="1" 
                step="0.1" 
                v-model="shade.opacity"
              >
            </div>
            <div class="btn-group">
              <button 
                v-for="color in shadeColors" 
                :key="color.value"
                @click="setShadeColor(color.value)"
              >
                {{ color.label }}
              </button>
            </div>
          </div>
          
          <div class="control-group">
            <h4>十字准星</h4>
            <div class="crosshair-controls">
              <button @click="showCrosshair">显示准星</button>
              <button @click="hideCrosshair">隐藏准星</button>
              <button @click="moveCrosshair">移动准星</button>
            </div>
          </div>
          
          <div class="control-group">
            <h4>其他操作</h4>
            <div class="btn-group">
              <button @click="resetSelection">重置选区</button>
              <button @click="cropImage">裁剪图片</button>
              <button @click="downloadResult">下载结果</button>
            </div>
          </div>
        </div>
      </div>
      
      <div class="component-info">
        <h3>组件说明</h3>
        <div class="info-grid">
          <div class="info-card" v-for="component in components" :key="component.name">
            <h4>{{ component.name }}</h4>
            <p>{{ component.description }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'

// Refs
const canvas = ref(null)
const previewCanvas = ref(null)

// 响应式状态
const selection = reactive({
  x: 100,
  y: 100,
  width: 300,
  height: 200,
  aspectRatio: 0,
  isMoving: false,
  isResizing: false,
  startX: 0,
  startY: 0
})

const shade = reactive({
  opacity: 0.5,
  color: '#000000'
})

const crosshair = reactive({
  visible: false,
  x: 200,
  y: 150
})

const interactionMode = ref('move')
const rotation = ref(0)

// 计算属性
const positionInfo = computed(() => {
  return `x: ${selection.x}, y: ${selection.y}`
})

const sizeInfo = computed(() => {
  return `${selection.width} × ${selection.height}`
})

const ratioInfo = computed(() => {
  const ratio = (selection.width / selection.height).toFixed(2)
  return selection.aspectRatio ? 
    `${selection.aspectRatio}:1` : 
    `${ratio}:1 (自由)`
})

// 常量数据
const aspectRatios = [
  { value: 0, label: '自由' },
  { value: 1, label: '1:1 (方形)' },
  { value: 1.777, label: '16:9 (宽屏)' },
  { value: 0.75, label: '3:4 (竖屏)' }
]

const shadeColors = [
  { value: '#000000', label: '黑色遮罩' },
  { value: '#3498db', label: '蓝色遮罩' },
  { value: '#e74c3c', label: '红色遮罩' }
]

const components = [
  { 
    name: 'CropperImage', 
    description: '承载原始图片的基础组件，负责加载和显示图片，支持跨域处理和图片大小限制。' 
  },
  { 
    name: 'CropperShade', 
    description: '创建遮罩层，突出显示裁剪区域，可调整透明度和颜色。' 
  },
  { 
    name: 'CropperHandle', 
    description: '裁剪框的控制手柄，用于调整选区大小和旋转，支持八个方向的手柄。' 
  },
  { 
    name: 'CropperCrosshair', 
    description: '十字准星，用于精确定位，可显示/隐藏和移动到指定位置。' 
  },
  { 
    name: 'CropperSelection', 
    description: '管理裁剪选区，包括位置、大小和约束条件，支持宽高比锁定。' 
  },
  { 
    name: 'CropperCanvas', 
    description: '画布容器，承载所有裁剪组件，提供绘制上下文和事件处理。' 
  },
  { 
    name: 'CropperViewer', 
    description: '实时预览裁剪结果的视图，同步显示裁剪区域的内容。' 
  }
]

// 图片引用
const img = new Image()
img.src = 'https://images.unsplash.com/photo-1506744038136-46273834b3fb'

// 方法
const setCanvasSize = () => {
  if (!canvas.value) return
  
  const container = canvas.value.parentElement
  canvas.value.width = container.clientWidth
  canvas.value.height = container.clientHeight
  drawCanvas()
}

const drawCanvas = () => {
  if (!canvas.value || !previewCanvas.value) return
  
  const ctx = canvas.value.getContext('2d')
  const previewCtx = previewCanvas.value.getContext('2d')
  
  // 清空画布
  ctx.clearRect(0, 0, canvas.value.width, canvas.value.height)
  
  // 绘制背景
  ctx.fillStyle = '#ecf0f1'
  ctx.fillRect(0, 0, canvas.value.width, canvas.value.height)
  
  // 绘制图片
  if (img.complete) {
    const scale = Math.min(
      canvas.value.width / img.width, 
      canvas.value.height / img.height
    )
    const x = (canvas.value.width - img.width * scale) / 2
    const y = (canvas.value.height - img.height * scale) / 2
    
    ctx.save()
    if (rotation.value !== 0) {
      ctx.translate(canvas.value.width / 2, canvas.value.height / 2)
      ctx.rotate(rotation.value * Math.PI / 180)
      ctx.translate(-canvas.value.width / 2, -canvas.value.height / 2)
    }
    ctx.drawImage(img, x, y, img.width * scale, img.height * scale)
    ctx.restore()
  }
  
  // 绘制遮罩
  ctx.fillStyle = shade.color
  ctx.globalAlpha = shade.opacity
  ctx.fillRect(0, 0, canvas.value.width, canvas.value.height)
  
  // 清除选区区域的遮罩
  ctx.globalCompositeOperation = 'destination-out'
  ctx.fillRect(
    selection.x, 
    selection.y, 
    selection.width, 
    selection.height
  )
  ctx.globalCompositeOperation = 'source-over'
  ctx.globalAlpha = 1
  
  // 绘制选区边框
  ctx.strokeStyle = '#3498db'
  ctx.lineWidth = 2
  ctx.strokeRect(
    selection.x, 
    selection.y, 
    selection.width, 
    selection.height
  )
  
  // 绘制控制手柄
  const handleSize = 8
  ctx.fillStyle = '#3498db'
  
  // 四角手柄
  ctx.fillRect(
    selection.x - handleSize/2, 
    selection.y - handleSize/2, 
    handleSize, 
    handleSize
  )
  ctx.fillRect(
    selection.x + selection.width - handleSize/2, 
    selection.y - handleSize/2, 
    handleSize, 
    handleSize
  )
  ctx.fillRect(
    selection.x - handleSize/2, 
    selection.y + selection.height - handleSize/2, 
    handleSize, 
    handleSize
  )
  ctx.fillRect(
    selection.x + selection.width - handleSize/2, 
    selection.y + selection.height - handleSize/2, 
    handleSize, 
    handleSize
  )
  
  // 四边手柄
  ctx.fillRect(
    selection.x + selection.width/2 - handleSize/2, 
    selection.y - handleSize/2, 
    handleSize, 
    handleSize
  )
  ctx.fillRect(
    selection.x + selection.width/2 - handleSize/2, 
    selection.y + selection.height - handleSize/2, 
    handleSize, 
    handleSize
  )
  ctx.fillRect(
    selection.x - handleSize/2, 
    selection.y + selection.height/2 - handleSize/2, 
    handleSize, 
    handleSize
  )
  ctx.fillRect(
    selection.x + selection.width - handleSize/2, 
    selection.y + selection.height/2 - handleSize/2, 
    handleSize, 
    handleSize
  )
  
  // 绘制十字准星
  if (crosshair.visible) {
    ctx.strokeStyle = '#e74c3c'
    ctx.lineWidth = 1
    ctx.setLineDash([5, 5])
    
    // 横线
    ctx.beginPath()
    ctx.moveTo(0, crosshair.y)
    ctx.lineTo(canvas.value.width, crosshair.y)
    ctx.stroke()
    
    // 竖线
    ctx.beginPath()
    ctx.moveTo(crosshair.x, 0)
    ctx.lineTo(crosshair.x, canvas.value.height)
    ctx.stroke()
    
    ctx.setLineDash([])
    
    // 中心点
    ctx.fillStyle = '#e74c3c'
    ctx.beginPath()
    ctx.arc(crosshair.x, crosshair.y, 4, 0, Math.PI * 2)
    ctx.fill()
  }
  
  // 更新预览
  updatePreview()
}

const updatePreview = () => {
  if (!previewCanvas.value || !canvas.value) return
  
  const previewCtx = previewCanvas.value.getContext('2d')
  previewCtx.clearRect(0, 0, previewCanvas.value.width, previewCanvas.value.height)
  
  if (img.complete) {
    const scale = Math.min(
      canvas.value.width / img.width, 
      canvas.value.height / img.height
    )
    const imgX = (canvas.value.width - img.width * scale) / 2
    const imgY = (canvas.value.height - img.height * scale) / 2
    
    // 计算源图像中的对应区域
    const srcX = (selection.x - imgX) / scale
    const srcY = (selection.y - imgY) / scale
    const srcWidth = selection.width / scale
    const srcHeight = selection.height / scale
    
    // 绘制到预览画布
    previewCtx.drawImage(
      img, 
      srcX, srcY, srcWidth, srcHeight,
      0, 0, previewCanvas.value.width, previewCanvas.value.height
    )
    
    // 绘制预览边框
    previewCtx.strokeStyle = '#3498db'
    previewCtx.lineWidth = 2
    previewCtx.strokeRect(0, 0, previewCanvas.value.width, previewCanvas.value.height)
  }
}

const setAspectRatio = (ratio) => {
  selection.aspectRatio = ratio
  if (ratio > 0) {
    selection.height = selection.width / ratio
  }
  drawCanvas()
}

const setShadeColor = (color) => {
  shade.color = color
  drawCanvas()
}

const setInteractionMode = (mode) => {
  interactionMode.value = mode
}

const showCrosshair = () => {
  crosshair.visible = true
  drawCanvas()
}

const hideCrosshair = () => {
  crosshair.visible = false
  drawCanvas()
}

const moveCrosshair = () => {
  crosshair.x = Math.random() * canvas.value.width
  crosshair.y = Math.random() * canvas.value.height
  drawCanvas()
}

const resetSelection = () => {
  selection.x = 100
  selection.y = 100
  selection.width = 300
  selection.height = 200
  selection.aspectRatio = 0
  rotation.value = 0
  drawCanvas()
}

const rotateSelection = () => {
  rotation.value = (rotation.value + 90) % 360
  drawCanvas()
}

const cropImage = () => {
  alert('裁剪功能已触发！在实际应用中，这里会执行裁剪操作。')
}

const downloadResult = () => {
  alert('下载功能已触发！在实际应用中，这里会下载裁剪后的图片。')
}

// 画布事件处理
const onCanvasMouseDown = (e) => {
  if (!canvas.value) return
  
  const rect = canvas.value.getBoundingClientRect()
  const x = e.clientX - rect.left
  const y = e.clientY - rect.top
  
  // 检查是否点击在选区内
  if (x >= selection.x && x <= selection.x + selection.width &&
      y >= selection.y && y <= selection.y + selection.height) {
    selection.isMoving = true
    selection.startX = x - selection.x
    selection.startY = y - selection.y
  }
}

const onCanvasMouseMove = (e) => {
  if (selection.isMoving && canvas.value) {
    const rect = canvas.value.getBoundingClientRect()
    const x = e.clientX - rect.left
    const y = e.clientY - rect.top
    
    selection.x = x - selection.startX
    selection.y = y - selection.startY
    
    // 限制选区不超出画布
    selection.x = Math.max(
      0, 
      Math.min(canvas.value.width - selection.width, selection.x)
    )
    selection.y = Math.max(
      0, 
      Math.min(canvas.value.height - selection.height, selection.y)
    )
    
    drawCanvas()
  }
}

const onCanvasMouseUp = () => {
  selection.isMoving = false
}

const onCanvasMouseLeave = () => {
  selection.isMoving = false
}

// 生命周期
onMounted(() => {
  setCanvasSize()
  window.addEventListener('resize', setCanvasSize)
  
  img.onload = () => {
    drawCanvas()
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', setCanvasSize)
})

// 监听器
watch([shade, selection], () => {
  drawCanvas()
}, { deep: true })

watch(() => shade.opacity, () => {
  drawCanvas()
})
</script>

<style scoped>
.cropper-demo {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 20px;
  color: #333;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

header {
  text-align: center;
  margin-bottom: 30px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 10px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

h1 {
  color: #2c3e50;
  margin-bottom: 10px;
}

.description {
  color: #7f8c8d;
  max-width: 800px;
  margin: 0 auto;
  line-height: 1.6;
}

.demo-area {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  margin-bottom: 30px;
}

.cropper-container {
  flex: 1;
  min-width: 500px;
  background: white;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.controls {
  flex: 0 0 300px;
  background: white;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.canvas-wrapper {
  position: relative;
  width: 100%;
  height: 400px;
  border: 2px dashed #bdc3c7;
  border-radius: 5px;
  overflow: hidden;
  margin-bottom: 20px;
  background: #ecf0f1;
}

.cropper-canvas {
  width: 100%;
  height: 100%;
}

.viewer-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 20px;
}

.viewer-label {
  font-weight: bold;
  margin-bottom: 10px;
  color: #2c3e50;
}

.preview-viewer {
  width: 200px;
  height: 150px;
  border: 2px solid #3498db;
  border-radius: 5px;
  overflow: hidden;
  background: white;
}

.control-group {
  margin-bottom: 20px;
}

h3 {
  color: #2c3e50;
  margin-bottom: 15px;
  padding-bottom: 5px;
  border-bottom: 1px solid #ecf0f1;
}

.btn-group {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

button {
  padding: 10px 15px;
  border: none;
  border-radius: 5px;
  background: #3498db;
  color: white;
  cursor: pointer;
  transition: all 0.3s;
  font-weight: 600;
  flex: 1;
  min-width: 80px;
}

button:hover {
  background: #2980b9;
  transform: translateY(-2px);
}

button.active {
  background: #e74c3c;
}

.slider-container {
  margin: 15px 0;
}

label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
}

input[type="range"] {
  width: 100%;
}

.component-info {
  background: white;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  margin-top: 20px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 15px;
}

.info-card {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 8px;
  border-left: 4px solid #3498db;
}

.info-card h4 {
  color: #2c3e50;
  margin-bottom: 10px;
}

.info-card p {
  color: #7f8c8d;
  font-size: 0.9rem;
  line-height: 1.5;
}

.status-bar {
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
  font-size: 0.9rem;
  color: #7f8c8d;
}

.crosshair-controls {
  display: flex;
  gap: 10px;
  align-items: center;
}

.crosshair-controls button {
  padding: 8px 12px;
  font-size: 0.9rem;
}

@media (max-width: 900px) {
  .demo-area {
    flex-direction: column;
  }
  
  .cropper-container, .controls {
    min-width: 100%;
  }
  
  .status-bar {
    flex-direction: column;
    gap: 5px;
  }
}
</style>