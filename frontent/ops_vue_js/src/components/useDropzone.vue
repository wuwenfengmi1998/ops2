<script setup>
import { ref, onMounted, onUnmounted,defineProps } from "vue";
import { useI18n } from "vue-i18n";
const { t, locale } = useI18n();
import Dropzone from "dropzone";
import "dropzone/dist/dropzone.css";

import { useUserStore } from "@/stores/user";
const userStore = useUserStore();

const dropzoneElement = ref(null);
let dropzoneInstance = null;
const files = ref([]);

const prop= defineProps({
  maxFiles: {
    type: Number,
    default: 5,
  },
  acceptedFiles:{
    type: String,
    default: "image/*",
  },
  maxFilesize:{
    type: Number,
    default: 10,
  }
});

// 初始化 Dropzone
const initDropzone = () => {
  if (!dropzoneElement.value) return;

  // 禁用自动发现
  Dropzone.autoDiscover = false;

  // 移除任何现有的 Dropzone 实例
  if (dropzoneInstance) {
    dropzoneInstance.destroy();
  }

  // 初始化新的实例
  dropzoneInstance = new Dropzone(dropzoneElement.value, {
    url: "/api/files/upload", // 上传地址
    // headers: {
    //   user_cookie: "cccc",
    // },
    method: "post",

    //uploadMultiple: true,

    previewTemplate: document.getElementById("custom-template").innerHTML,

    parallelUploads: 3, // 同时上传的文件数
    maxFilesize: prop.maxFilesize, // MB
    maxFiles: prop.maxFiles, // 最大文件数
    acceptedFiles:prop.acceptedFiles, // 接受的文件类型
    //addRemoveLinks: true, // 显示移除链接
    dictDefaultMessage: t("dropzone.upload_drop_or_click"),
    dictFallbackMessage: t("dropzone.upload_browser_not_supported"),
    dictFileTooBig:
      t("dropzone.upload_file_too_big") +
      "({{filesize}}MB). " +
      t("dropzone.upload_max_file_size") +
      "{{maxFilesize}}MB.",
    dictInvalidFileType: t("dropzone.upload_invalid_file_type"),
    dictResponseError: t("dropzone.upload_server_error") + "{{statusCode}}",
    //dictCancelUpload: t('dropzone.upload_cancel'),
    //dictUploadCanceled: t('dropzone.upload_canceled'),
    //dictCancelUploadConfirmation: t('dropzone.upload_cancel_confirmation'),
    dictRemoveFile: t("dropzone.upload_remove_file"),
    dictMaxFilesExceeded:
      t("dropzone.upload_max_files") +
      "{{maxFiles}}" +
      t("dropzone.upload_max_files_unit"),

    // 事件处理
    init: function () {
      this.on("success", (file, response) => {
        console.log("上传成功:", file, response);
        file.previewElement.addEventListener("click", function (e) {
          e.preventDefault();
          e.stopPropagation();

          // 处理点击事件
          console.log("缩略图被点击", file);

          // 可以在这里实现：
          // 1. 预览大图
          // 2. 显示文件详情
          // 3. 触发自定义操作
        });
      });
      this.on("error", (file, errorMessage) => {
        console.error("上传失败:", file.name, errorMessage);
      });
      this.on("removedfile", (file) => {
        console.log("remove:", file);
        //files.value = files.value.filter(f => f.name !== file.name)
      });
      this.on("addedfile", (file) => {
        console.log("addfile", file);
      });
      this.on("sending", function (file, xhr, formData) {
        // 获取表单值并添加到 FormData
        //console.log(userStore.userCookie.Value)
        formData.append("cookie", userStore.userCookie.Value);

      });
    },
  });
};

// 自定义方法
// const formatBytes = (bytes, decimals = 2) => {
//   if (bytes === 0) return '0 Bytes'
//   const k = 1024
//   const dm = decimals < 0 ? 0 : decimals
//   const sizes = ['Bytes', 'KB', 'MB', 'GB']
//   const i = Math.floor(Math.log(bytes) / Math.log(k))
//   return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
// }

// // 手动添加文件的方法
// const addFiles = (fileList) => {
//   if (dropzoneInstance) {
//     Array.from(fileList).forEach(file => {
//       dropzoneInstance.addFile(file)
//     })
//   }
// }

// // 获取所有已添加的文件
// const getAllFiles = () => {
//   return dropzoneInstance ? dropzoneInstance.files : []
// }

// // 清除所有文件
// const removeAllFiles = () => {
//   if (dropzoneInstance) {
//     dropzoneInstance.removeAllFiles(true)
//   }
// }

// 组件挂载时初始化
onMounted(() => {
  initDropzone();
});

// 组件卸载时销毁
onUnmounted(() => {
  if (dropzoneInstance) {
    dropzoneInstance.destroy();
  }
});
</script>

<template>
  <div>
    <div id="custom-template" style="display: none">
      <div class="dz-preview dz-file-preview my-custom-style">
        <div class="remove-btn" data-dz-remove>
          <!-- <i class="bi bi-x"></i> -->
          X
        </div>
        <div class="dz-image">
          <img data-dz-thumbnail alt="File preview" />
          <!-- 缩略图 -->
        </div>
        <div class="dz-details">
          <div class="dz-filename"><span data-dz-name></span></div>
          <!-- 文件名 -->
          <div class="dz-size"><span data-dz-size></span></div>
          <!-- 文件大小 -->
        </div>
        <div class="dz-progress">
          <span class="dz-upload" data-dz-uploadprogress></span>
          <!-- 进度条 -->
        </div>
        <div class="dz-success-mark" data-dz-successmark>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="icon icon-tabler icons-tabler-outline icon-tabler-circle-check"
          >
            <path stroke="none" d="M0 0h24v24H0z" fill="none" />
            <path d="M3 12a9 9 0 1 0 18 0a9 9 0 1 0 -18 0" />
            <path d="M9 12l2 2l4 -4" />
          </svg>
        </div>
        <!-- 成功标记 -->
        <div class="dz-error-mark" data-dz-errormark>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="240"
            height="240"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="icon icon-tabler icons-tabler-outline icon-tabler-circle-x"
          >
            <path stroke="none" fill="none" d="M0 0h24v24H0z" />
            <path d="M3 12a9 9 0 1 0 18 0a9 9 0 1 0 -18 0" />
            <path d="M10 10l4 4m0 -4l-4 4" />
          </svg>
        </div>
        <!-- 错误标记 -->
        <div class="dz-error-message"><span data-dz-errormessage></span></div>
        <!-- 错误信息 -->

        <!-- 移除按钮 -->
      </div>
    </div>
    <div ref="dropzoneElement" class="dropzone"></div>
    <!-- <div v-if="files.length > 0" class="mt-4">
      <h3>已选择的文件：</h3>
      <ul>
        <li v-for="file in files" :key="file.name">
          {{ file.name }} ({{ formatBytes(file.size) }})
        </li>
      </ul>
    </div> -->
  </div>
</template>

<style scoped>
.dz_mark {
  height: 60px;
  width: 60px;
}

.thumbnail-container {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  justify-content: center;
  padding: 20px;
  background-color: white;
  border-radius: 15px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
}

/* 缩略图样式 */
.thumbnail {
  width: var(--thumbnail-size);
  height: var(--thumbnail-size);
  border-radius: var(--border-radius);
  object-fit: cover;
  border: 2px solid #e9ecef;
  transition: all 0.3s ease;
}

.thumbnail:hover {
  transform: scale(1.05);
  border-color: #6c757d;
}

/* 缩略图包装器 */
.thumbnail-wrapper {
  position: relative;
  width: var(--thumbnail-size);
  height: var(--thumbnail-size);
  margin-bottom: 10px;
}

/* 移除按钮 */
.remove-btn {
  position: absolute;
  top: -12px;
  right: -12px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background-color: #dc3545;
  color: white;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  z-index: 10;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
}

.remove-btn:hover {
  background-color: #bb2d3b;
  transform: scale(1.1);
}

/* 文件名称 */
.file-name {
  font-size: 12px;
  text-align: center;
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #495057;
}

/* 上传区域 */
.upload-area {
  border: 2px dashed #dee2e6;
  border-radius: 15px;
  padding: 30px;
  text-align: center;
  background-color: #f8f9fa;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 20px;
}

.upload-area:hover {
  border-color: #6c757d;
  background-color: #e9ecef;
}

.upload-icon {
  font-size: 48px;
  color: #6c757d;
  margin-bottom: 10px;
}

.preview-title {
  color: #343a40;
  border-bottom: 2px solid #e9ecef;
  padding-bottom: 10px;
  margin-bottom: 20px;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: #6c757d;
}

.empty-state i {
  font-size: 48px;
  margin-bottom: 15px;
  color: #adb5bd;
}

.counter-badge {
  position: absolute;
  top: -5px;
  right: -5px;
  background-color: #0d6efd;
  color: white;
  border-radius: 50%;
  width: 20px;
  height: 20px;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.thumbnail-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 20px;
}

.file-input {
  display: none;
}
</style>
