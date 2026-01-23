<script setup>
import { onMounted, ref, watch, defineProps, reactive } from "vue";

import flatpickr from "flatpickr";
import "flatpickr/dist/flatpickr.css";
import "flatpickr/dist/l10n/zh.js";

import { useI18n } from "vue-i18n";
const { t, locale } = useI18n();

const datatimepack = ref();

const prop = defineProps({
  setdef: {
    type: String,
    default: "",
  },
  max_date: {
    type: [String, Date, Function],
    default: () => new Date(), // 默认值为当前时间
  },
});

const datatimepack_config = reactive({
  enableTime: true,
  dateFormat: "Y-m-d H:i",
  minuteIncrement: 1,
  time_24hr: true,
  maxDate: prop.max_date, // 只能选择当前时间之前的时间
  //locale:"zh"
});

const sele_data = reactive();

function getCurrentDateTime() {
  const now = new Date();

  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, "0"); // 月份从0开始
  const day = String(now.getDate()).padStart(2, "0");
  const hours = String(now.getHours()).padStart(2, "0");
  const minutes = String(now.getMinutes()).padStart(2, "0");

  return `${year}-${month}-${day} ${hours}:${minutes}`;
}

watch(locale, () => {
  if (locale.value == "zh-CN") {
    datatimepack_config.locale = "zh";
  } else {
    datatimepack_config.locale = "en";
  }
  //console.log(locale.value=="zh-CN"?"zh":"en")
});


onMounted(() => {
  // @formatter:off
  //console.log(getCurrentDateTime())
  //sele_data=getCurrentDateTime();

  // console.log(prop.setdef)

  if (prop.setdef == "") {
    datatimepack_config.defaultDate = getCurrentDateTime();
  } else {
    datatimepack_config.defaultDate = prop.setdef;
  }

  datatimepack_config.locale = locale.value == "zh-CN" ? "zh" : "en";
  flatpickr(datatimepack.value, datatimepack_config);
});

defineExpose({});
</script>

<template>
  <div></div>
  <input
    ref="datatimepack"
    type="datetime-local"
    class="form-control"
    v-model="sele_data"
  />
</template>
