<script setup>
import { ref, onMounted, onUnmounted, defineProps, reactive } from "vue";
import { useI18n } from "vue-i18n";
const { t, locale } = useI18n();
import Dropzone from "dropzone";
import "dropzone/dist/dropzone.css";

import { useUserStore } from "@/stores/user";

import "fslightbox";

const userStore = useUserStore();

const dropzoneElement = ref(null);
var dropzoneInstance = null;

const files = reactive([]);

function get_file_from_uuid(uuid) {
  if (files.length != 0) {
    for (var i = 0; i < files.length; i++) {
      if (files[i].uuid == uuid) {
        return i;
      }
    }
    return -1;
  }
  return -2;
}

function remove_file_from_uuie(uuid) {
  //devare files[uuid]
  var id = get_file_from_uuid(uuid);
  if (id >= 0) {
    files.splice(id, 1);
  }
}

const prop = defineProps({
  maxFiles: {
    type: Number,
    default: 5,
  },
  acceptedFiles: {
    type: String,
    default: "image/*",
  },
  maxFilesize: {
    type: Number,
    default: 10,
  },
  uploadURL: {
    type: String,
    default: "/api/files/upload",
  },
  /** 初始已有文件 [{ hash, name, ... }] */
  initialFiles: {
    type: Array,
    default: () => [],
  },
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
    url: prop.uploadURL, // 上传地址
    // headers: {
    //   user_cookie: "cccc",
    // },
    method: "post",

    //uploadMultiple: true,

    previewTemplate: document.getElementById("custom-template").innerHTML,

    parallelUploads: 3, // 同时上传的文件数
    maxFilesize: prop.maxFilesize, // MB
    maxFiles: prop.maxFiles, // 最大文件数
    acceptedFiles: prop.acceptedFiles, // 接受的文件类型
    //addRemoveLinks: true, // 显示移除链接
    dictDefaultMessage: t("dropzone.upload_drop_or_click"),
    dictFallbackMessage: t("dropzone.upload_browser_not_supported"),
    dictFivarooBig:
      t("dropzone.upload_file_too_big") +
      "({{filesize}}MB). " +
      t("dropzone.upload_max_file_size") +
      "{{maxFilesize}}MB.",
    dictInvalidFivarype: t("dropzone.upload_invalid_file_type"),
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
        //console.log("上传成功:", file, response);

        // 移除旧的事件监听器，避免重复绑定
        const oldHandler = file.previewElement._lightboxClickHandler;
        if (oldHandler) {
          file.previewElement.removeEventListener("click", oldHandler);
        }

        // 创建新的点击处理器
        const clickHandler = function (e) {
          e.preventDefault();
          e.stopPropagation();

          // 每次点击创建新实例，sources 自动是空的
          const lightbox = new FsLightbox();

          var dis_id = 0;
          var dis_id_t = 0;
          for (var i = 0; i < files.length; i++) {
            if (files[i]["is_upload"] == true) {
              lightbox.props.sources.push(files[i]["get_url"]);
              if (files[i]["uuid"] == file.upload.uuid) {
                dis_id = dis_id_t;
              }
            }
            dis_id_t += 1;
          }

          lightbox.open(dis_id);
        };

        // 保存处理器引用，以便后续移除
        file.previewElement._lightboxClickHandler = clickHandler;
        file.previewElement.addEventListener("click", clickHandler);

        var file_id = get_file_from_uuid(file.upload.uuid);
        if (file_id >= 0) {
          files[file_id]["hash"] = response.return.hash;
          files[file_id]["get_url"] = response.return.get;
          files[file_id]["download_url"] = response.return.download;
          files[file_id]["file_name"] = file.name;
          files[file_id]["file_size"] = file.size;

          files[file_id]["is_upload"] = true;

          //console.log(files)
        }

        //files.push(t)
        // files[file.upload.uuid]=t
        // console.log(files)

        // lightbox.props.sources.push(t.get_url)
        // console.log(lightbox)
      });
      this.on("error", (file, errorMessage) => {
        console.error("上传失败:", file.name, errorMessage);
      });
      this.on("removedfile", (file) => {
        //console.log("remove:", file);
        //files.value = files.value.filter(f => f.name !== file.name)
        remove_file_from_uuie(file.upload.uuid);
        //console.log(files)
      });
      this.on("addedfile", (file) => {
        //添加文件

        //控制排序 需要从添加文件开始操作

        //限制文件数量
        if (files.length < prop.maxFiles) {
          var t = {
            uuid: file.upload.uuid,
            is_upload: false,
          };
          files.push(t);
        } else {
          this.removeFile(file);
        }

        //console.log(files);
      });
      this.on("sending", function (file, xhr, formData) {
        // 获取表单值并添加到 FormData
        //console.log(userStore.userCookie.Value)
        formData.append("cookie", userStore.userCookie.Value);
      });
    },
  });
};


function return_files() {
  return files;
}

// 加载初始已有文件（编辑场景）
function loadInitialFiles() {
  if (!dropzoneInstance || !prop.initialFiles?.length) return;
  

  prop.initialFiles.forEach((f) => {
    // 构造 Dropzone 期望的 mock file 对象
     console.log(f)
    // const mockFile = {
    //   name: f.Name,
    //   size: f.Size,
    //   type: f.Mime,
    //   status: Dropzone.SUCCESS,
    //   accepted: true,
    //   upload: { uuid: f.Sha256 },
    //   previewElement: null,
    //   _removeLink: null,
    // };
    // // 通知 Dropzone "这是一个已存在的文件，不要上传"
    // dropzoneInstance.emit("addedfile", mockFile);
    // dropzoneInstance.emit("complete", mockFile);
    // dropzoneInstance.files.push(mockFile);
    // // 填充上传结果字段
    // const url = `/api/files/get/${f.Sha256}`;
    // files.push({
    //   uuid: f.Sha256,
    //   hash: f.Sha256,
    //   get_url: url,
    //   download_url: `/api/files/download/${f.Sha256}`,
    //   file_name: f.Name,
    //   file_size: f.Size,
    //   is_upload: true,
    // });
  });
}

// 组件挂载时初始化
onMounted(() => {
  initDropzone();
  loadInitialFiles();
});

// 组件卸载时销毁
onUnmounted(() => {
  if (dropzoneInstance) {
    dropzoneInstance.destroy();
  }
});

defineExpose({
  return_files,
});
</script>

<template>
  <div>
    <div id="custom-template" style="display: none">
      <div class="dz-preview dz-file-preview">
        <div data-dz-remove class="dz-remove">✕</div>
        <div class="dz-image"><img data-dz-thumbnail alt="" /></div>
        <div class="dz-details">
          <div class="dz-filename"><span data-dz-name></span></div>
          <div class="dz-size"><span data-dz-size></span></div>
        </div>
        <div class="dz-progress"><span class="dz-upload" data-dz-uploadprogress></span></div>
        <div class="dz-success-mark" data-dz-successmark>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="#22c55e" stroke="#22c55e" stroke-width="3">
            <path d="M9 12l2 2l4 -4" stroke="none"/>
          </svg>
        </div>
        <div class="dz-error-mark" data-dz-errormark>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="#ef4444" stroke="#ef4444" stroke-width="3">
            <path d="M10 10l4 4m0 -4l-4 4" stroke="none"/>
          </svg>
        </div>
        <div class="dz-error-message"><span data-dz-errormessage></span></div>
      </div>
    </div>
    
    <div class="text-end text-sm text-gray-500">{{ files.length }}/{{ maxFiles }}</div>
    <div ref="dropzoneElement" class="dropzone mt-2"></div>
  </div>
</template>

<style scoped>
/* 覆盖 Dropzone 默认样式 */
:deep(.dropzone) {
  min-height: 120px;
  border: 2px dashed #dee2e6;
  border-radius: 0.5rem;
  padding: 1rem;
  background: #f8f9fa;
}

:deep(.dropzone .dz-preview) {
  margin: 0.5rem;
  position: relative;
  display: inline-block;
  vertical-align: top;
  background: transparent !important;
}

:deep(.dropzone .dz-remove) {
  position: absolute;
  top: -6px;
  left: -6px;
  width: 18px;
  height: 18px;
  background: #ef4444;
  color: white;
  border-radius: 50%;
  font-size: 10px;
  line-height: 18px;
  text-align: center;
  cursor: pointer;
  z-index: 10;
}

:deep(.dropzone .dz-remove:hover) {
  background: #dc2626;
}

:deep(.dropzone .dz-image-preview) {
  background: transparent !important;
}

:deep(.dropzone .dz-image) {
  width: 80px !important;
  z-index: 0 !important;
  height: 80px !important;
  border-radius: 0.5rem !important;
}

:deep(.dropzone .dz-image img) {
  width: 100% !important;
  height: 100% !important;
  object-fit: cover;
}

:deep(.dropzone .dz-details) {
  opacity: 1 !important;
  position: static !important;
  padding: 0.25rem !important;
  font-size: 0.75rem;
  text-align: center;
}

:deep(.dropzone .dz-filename) {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 80px;
}

:deep(.dropzone .dz-size) {
  font-size: 0.65rem;
  color: #9ca3af;
}

:deep(.dropzone .dz-success-mark),
:deep(.dropzone .dz-error-mark) {
  width: 24px !important;
  height: 24px !important;
  margin-left: -12px !important;
  margin-top: -12px !important;
  background: rgba(255, 255, 255, 0.9) !important;
  border-radius: 50% !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
}

:deep(.dropzone .dz-success-mark svg),
:deep(.dropzone .dz-error-mark svg) {
  width: 14px !important;
  height: 14px !important;
  fill: #22c55e !important;
}

:deep(.dropzone .dz-error-mark svg) {
  fill: #ef4444 !important;
}

:deep(.dropzone .dz-progress) {
  position: static !important;
  width: 100% !important;
  height: 4px !important;
  border: none !important;
  background: #e5e7eb !important;
  border-radius: 2px !important;
  margin-top: 0.25rem !important;
}

:deep(.dropzone .dz-progress .dz-upload) {
  background: #3b82f6 !important;
  border-radius: 2px !important;
  width: 0 !important;
}
</style>
