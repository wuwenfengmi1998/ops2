<script setup>
import { onMounted, ref, watch, defineProps, reactive } from "vue";
import flatpickr from "flatpickr";
import "flatpickr/dist/flatpickr.css";
import "flatpickr/dist/l10n/zh.js";
import { useI18n } from "vue-i18n";
const { t, locale } = useI18n();
const datatimepack = ref();
const prop = defineProps({
  setdef: { type: String, default: "" },
  max_date: { type: [String, Date, Function], default: () => new Date() },
});
const datatimepack_config = reactive({
  enableTime: true,
  dateFormat: "Y-m-d H:i",
  minuteIncrement: 1,
  time_24hr: true,
  maxDate: prop.max_date,
});
const emit = defineEmits(['update:modelValue'])
const handleChange = (e) => { emit("update:modelValue", e.target.value); };
function getCurrentDateTime() {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, "0");
  const day = String(now.getDate()).padStart(2, "0");
  const hours = String(now.getHours()).padStart(2, "0");
  const minutes = String(now.getMinutes()).padStart(2, "0");
  return `${year}-${month}-${day} ${hours}:${minutes}`;
}
watch(locale, () => {
  if (locale.value == "zh-CN") { datatimepack_config.locale = "zh"; }
  else { datatimepack_config.locale = "en"; }
});
onMounted(() => {
  if (prop.setdef == "") { datatimepack_config.defaultDate = getCurrentDateTime(); }
  else { datatimepack_config.defaultDate = prop.setdef; }
  datatimepack_config.locale = locale.value == "zh-CN" ? "zh" : "en";
  flatpickr(datatimepack.value, datatimepack_config);
  emit("update:modelValue", datatimepack_config.defaultDate);
});
defineExpose({});
</script>

<template>
  <input ref="datatimepack" type="datetime-local" class="w-full rounded-lg border border-gray-300 bg-white py-2 pr-3 text-sm focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 focus:outline-none dark:border-dk-muted dark:bg-dk-card dark:text-dk-text" @input="handleChange" />
</template>
