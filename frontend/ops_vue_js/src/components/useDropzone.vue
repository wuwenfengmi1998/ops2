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
  if (files.length != 0) { for (var i = 0; i < files.length; i++) { if (files[i].uuid == uuid) { return i; } } return -1; } return -2;
}
function remove_file_from_uuie(uuid) { var id = get_file_from_uuid(uuid); if (id >= 0) { files.splice(id, 1); } }
const prop = defineProps({
  maxFiles: { type: Number, default: 5 },
  acceptedFiles: { type: String, default: "image/*" },
  maxFilesize: { type: Number, default: 10 },
  uploadURL: { type: String, default: "/api/files/upload" },
});
const initDropzone = () => {
  if (!dropzoneElement.value) return;
  Dropzone.autoDiscover = false;
  if (dropzoneInstance) { dropzoneInstance.destroy(); }
  dropzoneInstance = new Dropzone(dropzoneElement.value, {
    url: prop.uploadURL,
    method: "post",
    previewTemplate: document.getElementById("custom-template").innerHTML,
    parallelUploads: 3,
    maxFilesize: prop.maxFilesize,
    maxFiles: prop.maxFiles,
    acceptedFiles: prop.acceptedFiles,
    dictDefaultMessage: t("dropzone.upload_drop_or_click"),
    dictFallbackMessage: t("dropzone.upload_browser_not_supported"),
    dictFileTooBig: t("dropzone.upload_file_too_big") + "({{filesize}}MB). " + t("dropzone.upload_max_file_size") + "{{maxFilesize}}MB.",
    dictInvalidFileType: t("dropzone.upload_invalid_file_type"),
    dictResponseError: t("dropzone.upload_server_error") + "{{statusCode}}",
    dictRemoveFile: t("dropzone.upload_remove_file"),
    dictMaxFilesExceeded: t("dropzone.upload_max_files") + " " + prop.maxFiles + t("dropzone.upload_max_files_unit"),
    init: function() {
      this.on("sending", function(file, xhr, formData) {
        formData.append("cookie", userStore.cookieValue);
      });
      this.on("success", function(file, serverResponse) {
        const data = JSON.parse(serverResponse);
        if (data.return && data.return.uuid) {
          file.uuid = data.return.uuid;
          files.push({ uuid: data.return.uuid, name: file.name, url: data.return.url || "" });
        }
      });
      this.on("removedfile", function(file) {
        if (file.uuid) { remove_file_from_uuie(file.uuid); }
      });
    },
  });
};
const emit = defineEmits(['files-updated', 'uuid-updated'])
onMounted(() => { initDropzone(); });
onUnmounted(() => { if (dropzoneInstance) { dropzoneInstance.destroy(); } });
defineExpose({ getFiles: () => files, getDropzone: () => dropzoneInstance });
</script>

<template>
  <div class="w-full">
    <div id="custom-template" class="dz-preview dz-file-preview hidden">
      <div class="relative inline-block rounded-lg border border-gray-200 bg-white p-2 dark:border-dk-muted dark:bg-dk-card">
        <img data-dz-thumbnail class="h-20 w-20 rounded object-cover" />
        <button data-dz-remove class="absolute -right-2 -top-2 flex h-6 w-6 items-center justify-center rounded-full bg-red-500 text-xs font-bold text-white shadow hover:bg-red-600">×</button>
      </div>
      <div class="mt-1 max-w-[5rem] truncate text-xs text-gray-600 dark:text-dk-subtle" data-dz-name></div>
      <div class="dz-progress mt-1 h-1 w-full rounded-full bg-gray-200 dark:bg-dk-muted"><span class="dz-upload block h-full rounded-full bg-blue-500" data-dz-uploadprogress></span></div>
      <div class="dz-error-message mt-1 text-xs text-red-500"><span data-dz-errormessage></span></div>
    </div>
    <div ref="dropzoneElement" class="dropzone cursor-pointer rounded-xl border-2 border-dashed border-gray-300 bg-gray-50 p-8 text-center transition-colors hover:border-blue-400 hover:bg-blue-50 dark:border-dk-muted dark:bg-dk-base dark:hover:border-blue-500 dark:hover:bg-dk-card">
      <div class="mb-2 text-4xl">📁</div>
      <div class="text-sm font-medium text-gray-600 dark:text-dk-subtle">{{ t('dropzone.upload_drop_or_click') }}</div>
    </div>
  </div>
</template>

<style scoped>
.dropzone { min-height: 150px; }
.dz-progress .dz-upload { transition: width 0.3s ease; }
</style>
