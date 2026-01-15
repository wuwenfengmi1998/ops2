
<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import Dropzone from 'dropzone'
import 'dropzone/dist/dropzone.css'

const dropzoneElement = ref(null)
let dropzoneInstance = null
const files = ref([])

// 初始化 Dropzone
const initDropzone = () => {
  if (!dropzoneElement.value) return

  // 禁用自动发现
  Dropzone.autoDiscover = false

  // 移除任何现有的 Dropzone 实例
  if (dropzoneInstance) {
    dropzoneInstance.destroy()
  }

  // 初始化新的实例
  dropzoneInstance = new Dropzone(dropzoneElement.value, {
    url: '/api/files/upload', // 上传地址
    method: 'post',
     // 确保不启用分片上传
    // chunking: false, // 明确禁用分片
    // forceChunking: false, // 强制不分片
    // chunkSize: false, // 不分片大小
    // retryChunks: false, // 不重试分片
    // parallelChunkUploads: false, // 不并行上传分片

    parallelUploads: 1, // 同时上传的文件数
    maxFilesize: 10, // MB
    maxFiles: 5, // 最大文件数
    acceptedFiles: 'image/*,.pdf,.doc,.docx', // 接受的文件类型
    addRemoveLinks: true, // 显示移除链接
    dictDefaultMessage: '拖放文件到这里或点击上传',
    dictFallbackMessage: '您的浏览器不支持拖放文件上传',
    dictFileTooBig: '文件太大 ({{filesize}}MB). 最大文件大小: {{maxFilesize}}MB.',
    dictInvalidFileType: '不支持此文件类型',
    dictResponseError: '服务器响应错误 {{statusCode}}',
    dictCancelUpload: '取消上传',
    dictUploadCanceled: '上传已取消',
    dictCancelUploadConfirmation: '确定要取消上传吗?',
    dictRemoveFile: '移除文件',
    dictMaxFilesExceeded: '您最多只能上传 {{maxFiles}} 个文件',
    
    // 事件处理
    init: function() {
      this.on('success', (file, response) => {
        console.log('上传成功:', file.name, response)
      }),
      this.on('error', (file, errorMessage) => {
        console.error('上传失败:', file.name, errorMessage)
      }),
      this.on('removedfile', (file) => {
        console.log('remove:')
        files.value = files.value.filter(f => f.name !== file.name)
      })
    }
  })
}

// 自定义方法
const formatBytes = (bytes, decimals = 2) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}

// 手动添加文件的方法
const addFiles = (fileList) => {
  if (dropzoneInstance) {
    Array.from(fileList).forEach(file => {
      dropzoneInstance.addFile(file)
    })
  }
}

// 获取所有已添加的文件
const getAllFiles = () => {
  return dropzoneInstance ? dropzoneInstance.files : []
}

// 清除所有文件
const removeAllFiles = () => {
  if (dropzoneInstance) {
    dropzoneInstance.removeAllFiles(true)
  }
}

// 组件挂载时初始化
onMounted(() => {
  initDropzone()
})

// 组件卸载时销毁
onUnmounted(() => {
  if (dropzoneInstance) {
    dropzoneInstance.destroy()
  }
})
</script>

<template>
  <div>
    <div ref="dropzoneElement" class="dropzone"></div>
    <div v-if="files.length > 0" class="mt-4">
      <h3>已选择的文件：</h3>
      <ul>
        <li v-for="file in files" :key="file.name">
          {{ file.name }} ({{ formatBytes(file.size) }})
        </li>
      </ul>
    </div>
  </div>
</template>


<style scoped>
.dropzone {
  border: 2px dashed #cccccc;
  border-radius: 5px;
  background: white;
  padding: 20px;
  min-height: 150px;
}
</style>