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
const emit = defineEmits(['crop_to_canvas'])
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
  cro_canv.value.$toCanvas().then((a) => { emit('crop_to_canvas',a) });
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
