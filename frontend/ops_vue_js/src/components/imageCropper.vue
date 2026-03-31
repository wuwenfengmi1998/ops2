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
  <div class="flex flex-col md:flex-row">
    <div v-show="!is_have_URL" class="w-full py-3 md:w-auto md:px-3">
      <button class="rounded-lg border border-blue-600 px-4 py-2 text-sm font-medium text-blue-600 transition-colors hover:bg-blue-50 dark:border-blue-400 dark:text-blue-400 dark:hover:bg-blue-900/20" @click="openFilePicker">{{ t("cropper.select_image") }}</button>
    </div>
    <cropper-canvas ref="cro_canv" class="cropper-container" :hidden="!is_have_URL" background scale-step="0.1">
      <cropper-image ref="cro_imag" src="" alt="Picture" initialCenterSize="cover" rotatable scalable skewable translatable></cropper-image>
      <cropper-shade hidden></cropper-shade>
      <cropper-handle action="move" plain></cropper-handle>
      <cropper-selection ref="cro_sele">
        <cropper-grid role="grid" covered></cropper-grid>
        <cropper-crosshair centered></cropper-crosshair>
        <cropper-handle action="move" theme-color="rgba(255, 255, 255, 0)"></cropper-handle>
      </cropper-selection>
    </cropper-canvas>
    <div v-show="is_have_URL" class="mt-3 flex gap-2 md:ml-3 md:mt-0">
      <button class="rounded-lg border border-blue-600 px-4 py-2 text-sm font-medium text-blue-600 transition-colors hover:bg-blue-50 dark:border-blue-400 dark:text-blue-400 dark:hover:bg-blue-900/20" @click="getsele">{{ t("cropper.crop_image") }}</button>
      <button class="rounded-lg border border-red-600 px-4 py-2 text-sm font-medium text-red-600 transition-colors hover:bg-red-50 dark:border-red-400 dark:text-red-400 dark:hover:bg-red-900/20" @click="cancel">{{ t("cropper.closs") }}</button>
    </div>
  </div>
</template>

<style scoped>
.cropper-container { box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06); border-radius: 0.5rem; }
</style>
